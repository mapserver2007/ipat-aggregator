package master_service

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/config"
	"github.com/sirupsen/logrus"
)

const (
	jockeyResultUrl      = "https://db.netkeiba.com/?pid=jockey_select&id=%s&year=%d&mode=%s&start_date=%d&end_date=%d"
	jockeyResultFileName = "jockey_result_%s.json"
)

type JockeyResult interface {
	CreateOrUpdate(ctx context.Context, jockeyResults []*data_cache_entity.JockeyResult) error
}

type jockeyResultService struct {
	jockeyResultRepository      repository.JockeyResultRepository
	jockeyResultEntityConverter converter.JockeyResultEntityConverter
	logger                      *logrus.Logger
}

func NewJockeyResult(
	jockeyResultRepository repository.JockeyResultRepository,
	jockeyResultEntityConverter converter.JockeyResultEntityConverter,
	logger *logrus.Logger,
) JockeyResult {
	return &jockeyResultService{
		jockeyResultRepository:      jockeyResultRepository,
		jockeyResultEntityConverter: jockeyResultEntityConverter,
		logger:                      logger,
	}
}

func (j *jockeyResultService) CreateOrUpdate(
	ctx context.Context,
	jockeyResults []*data_cache_entity.JockeyResult,
) error {
	taskCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	urls, err := j.createJockeyResultUrls(taskCtx)
	if err != nil {
		return err
	}

	if len(urls) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	const jockeyResultParallel = 5
	errorCh := make(chan error, jockeyResultParallel)
	resultCh := make(chan []*netkeiba_entity.JockeyResult, jockeyResultParallel)
	chunkSize := (len(urls) + jockeyResultParallel - 1) / jockeyResultParallel

	for i := 0; i < len(urls); i += chunkSize {
		end := min(i+chunkSize, len(urls))

		wg.Add(1)
		go func(splitUrls []string) {
			defer wg.Done()
			localJockeyResults := make([]*netkeiba_entity.JockeyResult, 0)
			j.logger.Infof("jockey result fetch processing: %v/%v", end, len(urls))
			for _, url := range splitUrls {
				time.Sleep(time.Millisecond)
				select {
				case <-taskCtx.Done():
					return
				default:
					jockeyResults, err := j.jockeyResultRepository.FetchJockeyResults(taskCtx, url)
					if err != nil {
						select {
						case errorCh <- err:
							cancel()
						default:
						}
						return
					}
					localJockeyResults = append(localJockeyResults, jockeyResults...)
				}
			}

			resultCh <- localJockeyResults
		}(urls[i:end])
	}

	wg.Wait()
	close(errorCh)
	close(resultCh)

	if err := <-errorCh; err != nil {
		return err
	}

	newJockeryResultMap := make(map[types.JockeyId][]*raw_entity.JockeyResult)
	for results := range resultCh {
		for _, jockeyResult := range results {
			newJockeryResultMap[types.JockeyId(jockeyResult.JockeyId())] = append(newJockeryResultMap[types.JockeyId(jockeyResult.JockeyId())], j.jockeyResultEntityConverter.NetKeibaToRaw(jockeyResult))
		}
	}

	for jockeyId, rawJockeyResults := range newJockeryResultMap {
		sort.Slice(rawJockeyResults, func(i, j int) bool {
			return rawJockeyResults[i].RaceId < rawJockeyResults[j].RaceId
		})
		jockeyResultInfo := raw_entity.JockeyResultInfo{
			JockeyResults: rawJockeyResults,
		}
		err := j.jockeyResultRepository.Write(ctx, fmt.Sprintf("%s/jockey_results/%s", config.CacheDir, fmt.Sprintf(jockeyResultFileName, jockeyId.Value())), &jockeyResultInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *jockeyResultService) createJockeyResultUrls(ctx context.Context) ([]string, error) {
	startDate, err := types.NewRaceDate(config.JockeyResultStartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := types.NewRaceDate(config.JockeyResultEndDate)
	if err != nil {
		return nil, err
	}

	years := make([]int, 0)
	for year := startDate.Year(); year <= endDate.Year(); year++ {
		years = append(years, year)
	}

	places := []string{
		"r1", "r2", "r3", "r4",
	}

	jockeyResultBaseUrls := make([]string, 0, len(config.JockeyResultTargetIds))
	for _, rawJockeyId := range config.JockeyResultTargetIds {
		for _, year := range years {
			for _, place := range places {
				url := fmt.Sprintf(jockeyResultUrl, rawJockeyId, year, place, startDate.Value(), endDate.Value())
				jockeyResultBaseUrls = append(jockeyResultBaseUrls, url)
			}
		}
	}

	taskCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	const raceIdParallel = 5
	errorCh := make(chan error, raceIdParallel)
	resultCh := make(chan []string, raceIdParallel)
	chunkSize := (len(jockeyResultBaseUrls) + raceIdParallel - 1) / raceIdParallel

	for i := 0; i < len(jockeyResultBaseUrls); i += chunkSize {
		end := min(i+chunkSize, len(jockeyResultBaseUrls))

		wg.Add(1)
		go func(splitUrls []string) {
			defer wg.Done()
			var localJockeyResultUrls []string
			j.logger.Infof("jockey result all pages fetch processing: %v/%v", end, len(jockeyResultBaseUrls))

			for _, url := range splitUrls {
				time.Sleep(time.Microsecond * 500)
				select {
				case <-taskCtx.Done():
					return
				default:
					urls, err := j.jockeyResultRepository.FetchUrls(taskCtx, url)
					if err != nil {
						select {
						case errorCh <- err:
							cancel()
						default:
						}
						return
					}
					localJockeyResultUrls = append(localJockeyResultUrls, urls...)
				}
			}

			resultCh <- localJockeyResultUrls
		}(jockeyResultBaseUrls[i:end])
	}

	wg.Wait()
	close(errorCh)
	close(resultCh)

	for err := range errorCh {
		if err != nil {
			return nil, err
		}
	}

	var jockeyResultUrls []string
	for results := range resultCh {
		jockeyResultUrls = append(jockeyResultUrls, results...)
	}

	return jockeyResultUrls, nil
}
