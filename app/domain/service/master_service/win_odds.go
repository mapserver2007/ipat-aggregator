package master_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/config"
	"github.com/sirupsen/logrus"
	neturl "net/url"
	"sort"
	"sync"
	"time"
)

const (
	winOddsUrl      = "https://race.netkeiba.com/api/api_get_jra_odds.html?race_id=%s&type=1&action=update"
	winOddsSpUrl    = "https://race.sp.netkeiba.com/?pid=api_get_jra_odds&race_id=%s&type=1&action=update"
	winOddsFileName = "odds_%d.json"
)

type WinOdds interface {
	Get(ctx context.Context) ([]*data_cache_entity.Odds, error)
	CreateOrUpdate(ctx context.Context, odds []*data_cache_entity.Odds, markers []*marker_csv_entity.AnalysisMarker) error
}

type winOddsService struct {
	oddsRepository      repository.OddsRepository
	oddsEntityConverter converter.OddsEntityConverter
	logger              *logrus.Logger
}

func NewWinOdds(
	oddsRepository repository.OddsRepository,
	oddsEntityConverter converter.OddsEntityConverter,
	logger *logrus.Logger,
) WinOdds {
	return &winOddsService{
		oddsRepository:      oddsRepository,
		oddsEntityConverter: oddsEntityConverter,
		logger:              logger,
	}
}

func (w *winOddsService) Get(ctx context.Context) ([]*data_cache_entity.Odds, error) {
	files, err := w.oddsRepository.List(ctx, fmt.Sprintf("%s/odds/win", config.CacheDir))
	if err != nil {
		return nil, err
	}

	var odds []*data_cache_entity.Odds
	for _, file := range files {
		rawRaceOddsList, err := w.oddsRepository.Read(ctx, fmt.Sprintf("%s/odds/win/%s", config.CacheDir, file))
		if err != nil {
			return nil, err
		}
		for _, rawRaceOdds := range rawRaceOddsList {
			raceId := types.RaceId(rawRaceOdds.RaceId)
			raceDate := types.RaceDate(rawRaceOdds.RaceDate)
			for _, rawOdds := range rawRaceOdds.Odds {
				odds = append(odds, w.oddsEntityConverter.RawToDataCache(rawOdds, raceId, raceDate))
			}
		}
	}

	return odds, nil
}

func (w *winOddsService) CreateOrUpdate(
	ctx context.Context,
	odds []*data_cache_entity.Odds,
	markers []*marker_csv_entity.AnalysisMarker,
) error {
	taskCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	urls := w.createOddsUrls(odds, markers)
	if len(urls) == 0 {
		return nil
	}

	oddsMap := w.createOddsMap(odds)

	var wg sync.WaitGroup
	const workerParallel = 5
	errorCh := make(chan error, 1)
	chunkSize := (len(urls) + workerParallel - 1) / workerParallel

	var mu sync.Mutex

	for i := 0; i < len(urls); i = i + chunkSize {
		end := i + chunkSize
		if end > len(urls) {
			end = len(urls)
		}

		wg.Add(1)
		go func(splitUrls []string) {
			defer wg.Done()
			w.logger.Infof("win odds fetch processing: %v/%v", end, len(urls))
			for _, url := range splitUrls {
				time.Sleep(time.Millisecond)
				select {
				case <-taskCtx.Done():
					return
				default:
					var newOdds []*raw_entity.Odds
					fetchOdds, err := w.oddsRepository.Fetch(taskCtx, url)
					if err != nil {
						select {
						case errorCh <- err:
							cancel()
						}
						return
					}

					raceId, err := w.parseUrl(url)
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

					for _, netKeibaFetchOdds := range fetchOdds {
						newOdds = append(newOdds, w.oddsEntityConverter.NetKeibaToRaw(netKeibaFetchOdds))
					}

					mu.Lock()
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
					mu.Unlock()
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
		err := w.oddsRepository.Write(ctx, fmt.Sprintf("%s/odds/win/%s", config.CacheDir, fmt.Sprintf(winOddsFileName, raceDate.Value())), &raceOddsInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *winOddsService) createOddsUrls(
	oddsList []*data_cache_entity.Odds,
	markers []*marker_csv_entity.AnalysisMarker,
) []string {
	var winOddsUrls []string
	oddsMap := map[types.RaceId]bool{}

	for _, odds := range oddsList {
		if _, ok := oddsMap[odds.RaceId()]; !ok {
			oddsMap[odds.RaceId()] = true
		}
	}

	for _, marker := range markers {
		if _, ok := oddsMap[marker.RaceId()]; !ok {
			winOddsUrls = append(winOddsUrls, fmt.Sprintf(winOddsUrl, marker.RaceId()))
		}
	}

	return winOddsUrls
}

func (w *winOddsService) createOddsMap(
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
			rawOddsList = append(rawOddsList, w.oddsEntityConverter.DataCacheToRaw(odds))
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

func (w *winOddsService) parseUrl(
	url string,
) (types.RaceId, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return "", err
	}
	raceId := u.Query().Get("race_id")

	return types.RaceId(raceId), nil
}
