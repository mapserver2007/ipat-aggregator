package master_service

import (
	"context"
	"fmt"
	neturl "net/url"
	"sort"
	"sync"
	"time"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"

	"github.com/mapserver2007/ipat-aggregator/config"
	"github.com/sirupsen/logrus"
)

const (
	trioOddsUrl      = "https://race.netkeiba.com/api/api_get_jra_odds.html?race_id=%s&type=7&sort=ninki&action=update"
	trioOddsSpUrl    = "https://race.sp.netkeiba.com/?pid=api_get_jra_odds&race_id=%s&type=7&sort=ninki&action=update"
	trioOddsFileName = "odds_%d.json"
)

type TrioOdds interface {
	Get(ctx context.Context) ([]*data_cache_entity.Odds, error)
	CreateOrUpdateV2(ctx context.Context, odds []*data_cache_entity.Odds, races []*data_cache_entity.Race) error
	CreateOrUpdate(ctx context.Context, odds []*data_cache_entity.Odds, markers []*marker_csv_entity.AnalysisMarker) error
}

type trioOddsService struct {
	oddsRepository      repository.OddsRepository
	oddsEntityConverter converter.OddsEntityConverter
	logger              *logrus.Logger
}

func NewTrioOdds(
	oddsRepository repository.OddsRepository,
	oddsEntityConverter converter.OddsEntityConverter,
	logger *logrus.Logger,
) TrioOdds {
	return &trioOddsService{
		oddsRepository:      oddsRepository,
		oddsEntityConverter: oddsEntityConverter,
		logger:              logger,
	}
}

func (o *trioOddsService) Get(ctx context.Context) ([]*data_cache_entity.Odds, error) {
	files, err := o.oddsRepository.List(ctx, fmt.Sprintf("%s/odds/trio", config.CacheDir))
	if err != nil {
		return nil, err
	}

	var odds []*data_cache_entity.Odds
	for _, file := range files {
		rawRaceOddsList, err := o.oddsRepository.Read(ctx, fmt.Sprintf("%s/odds/trio/%s", config.CacheDir, file))
		if err != nil {
			return nil, err
		}
		for _, rawRaceOdds := range rawRaceOddsList {
			raceId := types.RaceId(rawRaceOdds.RaceId)
			raceDate := types.RaceDate(rawRaceOdds.RaceDate)
			for _, rawOdds := range rawRaceOdds.Odds {
				odds = append(odds, o.oddsEntityConverter.RawToDataCache(rawOdds, raceId, raceDate))
			}
		}
	}

	return odds, nil
}

func (o *trioOddsService) CreateOrUpdateV2(
	ctx context.Context,
	odds []*data_cache_entity.Odds,
	races []*data_cache_entity.Race,
) error {
	taskCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	urls := o.createOddsUrlsV2(odds, races)
	if len(urls) == 0 {
		return nil
	}

	oddsMap := o.createOddsMap(odds)

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
			o.logger.Infof("trio odds fetch processing: %v/%v", end, len(urls))
			for _, url := range splitUrls {
				time.Sleep(time.Millisecond)
				select {
				case <-taskCtx.Done():
					return
				default:
					fetchOdds, err := o.oddsRepository.Fetch(taskCtx, url)
					if err != nil {
						select {
						case errorCh <- err:
							cancel()
						}
						return
					}

					raceId, err := o.parseUrl(url)
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
						newOdds = append(newOdds, o.oddsEntityConverter.NetKeibaToRaw(netKeibaFetchOdds))
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
		err := o.oddsRepository.Write(ctx, fmt.Sprintf("%s/odds/trio/%s", config.CacheDir, fmt.Sprintf(trioOddsFileName, raceDate.Value())), &raceOddsInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *trioOddsService) CreateOrUpdate(
	ctx context.Context,
	odds []*data_cache_entity.Odds,
	markers []*marker_csv_entity.AnalysisMarker,
) error {
	taskCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	urls := o.createOddsUrls(odds, markers)
	if len(urls) == 0 {
		return nil
	}

	markerHorseNumberMap := map[types.RaceId][]types.HorseNumber{}
	for _, marker := range markers {
		horseNumbers := make([]types.HorseNumber, 0, len(marker.MarkerMap()))
		for _, horseNumber := range marker.MarkerMap() {
			horseNumbers = append(horseNumbers, horseNumber)
		}
		markerHorseNumberMap[marker.RaceId()] = horseNumbers
	}

	oddsMap := o.createOddsMap(odds)

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
			o.logger.Infof("trio odds fetch processing: %v/%v", end, len(urls))
			for _, url := range splitUrls {
				time.Sleep(time.Millisecond)
				select {
				case <-taskCtx.Done():
					return
				default:
					newOdds := make([]*raw_entity.Odds, 0, 20) // 6頭BOXは20点
					fetchOdds, err := o.oddsRepository.Fetch(taskCtx, url)
					if err != nil {
						select {
						case errorCh <- err:
							cancel()
						}
						return
					}

					raceId, err := o.parseUrl(url)
					if err != nil {
						select {
						case errorCh <- err:
							cancel()
						}
						return
					}

					horseNumbers, ok := markerHorseNumberMap[raceId]
					if !ok {
						// 分析印がないレース(新馬、障害など)はスキップ
						continue
					}

					var raceDate types.RaceDate
					if len(fetchOdds) > 0 {
						raceDate = fetchOdds[0].RaceDate()
					}

					for _, netKeibaFetchOdds := range fetchOdds {
						if o.containsInSliceAll(horseNumbers, netKeibaFetchOdds.HorseNumbers()) {
							newOdds = append(newOdds, o.oddsEntityConverter.NetKeibaToRaw(netKeibaFetchOdds))
						}
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
		err := o.oddsRepository.Write(ctx, fmt.Sprintf("%s/odds/trio/%s", config.CacheDir, fmt.Sprintf(trioOddsFileName, raceDate.Value())), &raceOddsInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *trioOddsService) createOddsUrls(
	oddsList []*data_cache_entity.Odds,
	markers []*marker_csv_entity.AnalysisMarker,
) []string {
	var trioOddsUrls []string
	oddsMap := map[types.RaceId]bool{}

	for _, odds := range oddsList {
		if _, ok := oddsMap[odds.RaceId()]; !ok {
			oddsMap[odds.RaceId()] = true
		}
	}

	for _, marker := range markers {
		if _, ok := oddsMap[marker.RaceId()]; !ok {
			trioOddsUrls = append(trioOddsUrls, fmt.Sprintf(trioOddsUrl, marker.RaceId()))
		}
	}

	return trioOddsUrls
}

func (o *trioOddsService) createOddsUrlsV2(
	oddsList []*data_cache_entity.Odds,
	races []*data_cache_entity.Race,
) []string {
	var trioOddsUrls []string
	raceIdMap := map[types.RaceId]bool{}

	for _, odds := range oddsList {
		if _, ok := raceIdMap[odds.RaceId()]; !ok {
			raceIdMap[odds.RaceId()] = true
		}
	}

	fetchableRaceStartDate, err := types.NewRaceDate(config.RaceStartDate)
	if err != nil {
		o.logger.Errorf("failed to create race start date: %v", err)
		return nil
	}

	fetchableRaceEndDate, err := types.NewRaceDate(config.RaceEndDate)
	if err != nil {
		o.logger.Errorf("failed to create race end date: %v", err)
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
					trioOddsUrls = append(trioOddsUrls, fmt.Sprintf(trioOddsUrl, race.RaceId()))
				}
			}
		}
	}

	return trioOddsUrls
}

func (o *trioOddsService) createOddsMap(
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
			rawOddsList = append(rawOddsList, o.oddsEntityConverter.DataCacheToRaw(odds))
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

func (o *trioOddsService) parseUrl(
	url string,
) (types.RaceId, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return "", err
	}
	raceId := u.Query().Get("race_id")

	return types.RaceId(raceId), nil
}

func (o *trioOddsService) containsInSliceAll(slice1, slice2 []types.HorseNumber) bool {
	for _, val2 := range slice2 {
		found := false
		for _, val1 := range slice1 {
			if val1 == val2 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
