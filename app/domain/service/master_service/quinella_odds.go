package master_service

import (
	"context"
	"fmt"
	neturl "net/url"
	"sort"
	"sync"
	"time"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/config"
	"github.com/sirupsen/logrus"
)

const (
	quinellaOddsUrl      = "https://race.netkeiba.com/api/api_get_jra_odds.html?race_id=%s&type=4&sort=ninki&action=update"
	quinellaOddsSpUrl    = "https://race.sp.netkeiba.com/?pid=api_get_jra_odds&race_id=%s&type=4&sort=ninki&action=update"
	quinellaOddsFileName = "odds_%d.json"
)

type QuinellaOdds interface {
	Get(ctx context.Context) ([]*data_cache_entity.Odds, error)
	CreateOrUpdateV2(ctx context.Context, odds []*data_cache_entity.Odds, races []*data_cache_entity.Race) error
}

type quinellaOddsService struct {
	oddsRepository      repository.OddsRepository
	oddsEntityConverter converter.OddsEntityConverter
	logger              *logrus.Logger
}

func NewQuinellaOdds(
	oddsRepository repository.OddsRepository,
	oddsEntityConverter converter.OddsEntityConverter,
	logger *logrus.Logger,
) QuinellaOdds {
	return &quinellaOddsService{
		oddsRepository:      oddsRepository,
		oddsEntityConverter: oddsEntityConverter,
		logger:              logger,
	}
}

func (q *quinellaOddsService) Get(ctx context.Context) ([]*data_cache_entity.Odds, error) {
	files, err := q.oddsRepository.List(ctx, fmt.Sprintf("%s/odds/quinella", config.CacheDir))
	if err != nil {
		return nil, err
	}

	var odds []*data_cache_entity.Odds
	for _, file := range files {
		rawRaceOddsList, err := q.oddsRepository.Read(ctx, fmt.Sprintf("%s/odds/quinella/%s", config.CacheDir, file))
		if err != nil {
			return nil, err
		}
		for _, rawRaceOdds := range rawRaceOddsList {
			raceId := types.RaceId(rawRaceOdds.RaceId)
			raceDate := types.RaceDate(rawRaceOdds.RaceDate)
			for _, rawOdds := range rawRaceOdds.Odds {
				odds = append(odds, q.oddsEntityConverter.RawToDataCache(rawOdds, raceId, raceDate))
			}
		}
	}

	return odds, nil
}

func (q *quinellaOddsService) CreateOrUpdateV2(
	ctx context.Context,
	odds []*data_cache_entity.Odds,
	races []*data_cache_entity.Race,
) error {
	taskCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	urls := q.createOddsUrlsV2(odds, races)
	if len(urls) == 0 {
		return nil
	}

	oddsMap := q.createOddsMap(odds)

	var wg sync.WaitGroup
	const workerParallel = 5
	errorCh := make(chan error, 1)
	chunkSize := (len(urls) + workerParallel - 1) / workerParallel

	for i := 0; i < len(urls); i = i + chunkSize {
		end := i + chunkSize
		if end > len(urls) {
			end = len(urls)
		}

		wg.Add(1)
		go func(splitUrls []string) {
			defer wg.Done()
			q.logger.Infof("quinella odds fetch processing: %v/%v", end, len(urls))
			for _, url := range splitUrls {
				time.Sleep(time.Millisecond)
				select {
				case <-taskCtx.Done():
					return
				default:
					fetchOdds, err := q.oddsRepository.Fetch(taskCtx, url)
					if err != nil {
						select {
						case errorCh <- err:
							cancel()
						}
						return
					}

					raceId, err := q.parseUrl(url)
					if err != nil {
						select {
						case errorCh <- err:
							cancel()
						}
						return
					}

					var raceDate types.RaceDate
					if len(fetchOdds) > 0 {
						raceDate = fetchOdds[0].RaceDate()
					}

					newOdds := make([]*raw_entity.Odds, 0, len(fetchOdds))
					for _, netKeibaFetchOdds := range fetchOdds {
						newOdds = append(newOdds, q.oddsEntityConverter.NetKeibaToRaw(netKeibaFetchOdds))
					}

					if _, ok := oddsMap[raceDate]; !ok {
						oddsMap[raceDate] = make([]*raw_entity.RaceOdds, 0)
					}

					sort.Slice(newOdds, func(i, j int) bool {
						return newOdds[i].Popular < newOdds[j].Popular
					})

					oddsMap[raceDate] = append(oddsMap[raceDate], &raw_entity.RaceOdds{
						RaceId:   raceId.String(),
						RaceDate: raceDate.Value(),
						Odds:     newOdds,
					})
				}
			}
		}(urls[i:end])
	}

	wg.Wait()
	close(errorCh)

	if err := <-errorCh; err != nil {
		return err
	}

	for _, raceDate := range service.SortedRaceDateKeys(oddsMap) {
		rawRaceOddsList := oddsMap[raceDate]
		sort.Slice(rawRaceOddsList, func(i, j int) bool {
			return rawRaceOddsList[i].RaceId < rawRaceOddsList[j].RaceId
		})
		raceOddsInfo := raw_entity.RaceOddsInfo{
			RaceOdds: rawRaceOddsList,
		}
		err := q.oddsRepository.Write(ctx, fmt.Sprintf("%s/odds/quinella/%s", config.CacheDir, fmt.Sprintf(quinellaOddsFileName, raceDate.Value())), &raceOddsInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (q *quinellaOddsService) createOddsUrlsV2(
	oddsList []*data_cache_entity.Odds,
	races []*data_cache_entity.Race,
) []string {
	var quinellaOddsUrls []string
	raceIdMap := map[types.RaceId]bool{}

	for _, odds := range oddsList {
		if _, ok := raceIdMap[odds.RaceId()]; !ok {
			raceIdMap[odds.RaceId()] = true
		}
	}

	fetchableRaceStartDate, err := types.NewRaceDate(config.RaceStartDate)
	if err != nil {
		q.logger.Errorf("failed to create race start date: %v", err)
		return nil
	}

	fetchableRaceEndDate, err := types.NewRaceDate(config.RaceEndDate)
	if err != nil {
		q.logger.Errorf("failed to create race end date: %v", err)
		return nil
	}

	for _, race := range races {
		if race.RaceDate() >= fetchableRaceStartDate && race.RaceDate() <= fetchableRaceEndDate {
			// JRA以外はオッズ取得できないためスキップ
			if race.Organizer() != types.JRA {
				continue
			}

			// 新馬、障害はスキップ
			switch race.Class() {
			case types.MakeDebut, types.JumpMaiden, types.JumpOpenClass, types.JumpGrade1, types.JumpGrade2, types.JumpGrade3:
				continue
			default:
				if _, ok := raceIdMap[race.RaceId()]; !ok {
					quinellaOddsUrls = append(quinellaOddsUrls, fmt.Sprintf(quinellaOddsUrl, race.RaceId()))
				}
			}
		}
	}

	return quinellaOddsUrls
}

func (q *quinellaOddsService) createOddsMap(
	analysisOdds []*data_cache_entity.Odds,
) map[types.RaceDate][]*raw_entity.RaceOdds {
	oddsMap := map[types.RaceDate][]*raw_entity.RaceOdds{}
	raceIdOddsMap := map[types.RaceId][]*data_cache_entity.Odds{}

	for _, odds := range analysisOdds {
		if _, ok := raceIdOddsMap[odds.RaceId()]; !ok {
			raceIdOddsMap[odds.RaceId()] = make([]*data_cache_entity.Odds, 0)
		}
		raceIdOddsMap[odds.RaceId()] = append(raceIdOddsMap[odds.RaceId()], odds)
	}

	for _, raceId := range service.SortedRaceIdKeys(raceIdOddsMap) {
		oddsList := raceIdOddsMap[raceId]
		raceDate := oddsList[0].RaceDate()
		rawOddsList := make([]*raw_entity.Odds, 0, len(oddsList))
		for _, odds := range oddsList {
			rawOddsList = append(rawOddsList, q.oddsEntityConverter.DataCacheToRaw(odds))
		}

		if _, ok := oddsMap[raceDate]; !ok {
			oddsMap[raceDate] = make([]*raw_entity.RaceOdds, 0)
		}

		oddsMap[raceDate] = append(oddsMap[raceDate], &raw_entity.RaceOdds{
			RaceId:   raceId.String(),
			RaceDate: raceDate.Value(),
			Odds:     rawOddsList,
		})
	}

	return oddsMap
}

func (q *quinellaOddsService) parseUrl(
	url string,
) (types.RaceId, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return "", err
	}
	raceId := u.Query().Get("race_id")

	return types.RaceId(raceId), nil
}
