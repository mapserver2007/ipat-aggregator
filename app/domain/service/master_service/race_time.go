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
	raceDBUrl        = "https://db.netkeiba.com/race/%s/"
	raceTimeFileName = "race_time_%d.json"
)

type RaceTime interface {
	Get(ctx context.Context) ([]*data_cache_entity.RaceTime, error)
	CreateOrUpdate(
		ctx context.Context,
		raceTimes []*data_cache_entity.RaceTime,
		races []*data_cache_entity.Race,
		raceDateMap map[types.RaceDate][]types.RaceId,
	) error
}

type raceTimeService struct {
	raceTimeRepository      repository.RaceTimeRepository
	raceTimeEntityConverter converter.RaceTimeEntityConverter
	logger                  *logrus.Logger
}

func NewRaceTime(
	raceTimeRepository repository.RaceTimeRepository,
	raceTimeEntityConverter converter.RaceTimeEntityConverter,
	logger *logrus.Logger,
) RaceTime {
	return &raceTimeService{
		raceTimeRepository:      raceTimeRepository,
		raceTimeEntityConverter: raceTimeEntityConverter,
		logger:                  logger,
	}
}

func (r *raceTimeService) Get(
	ctx context.Context,
) ([]*data_cache_entity.RaceTime, error) {
	files, err := r.raceTimeRepository.List(ctx, fmt.Sprintf("%s/race_times", config.CacheDir))
	if err != nil {
		return nil, err
	}

	var raceTimes []*data_cache_entity.RaceTime
	for _, file := range files {
		rawRaceTimes, err := r.raceTimeRepository.Read(ctx, fmt.Sprintf("%s/race_times/%s", config.CacheDir, file))
		if err != nil {
			return nil, err
		}
		for _, rawRaceTime := range rawRaceTimes {
			raceTimes = append(raceTimes, r.raceTimeEntityConverter.RawToDataCache(rawRaceTime))
		}
	}

	return raceTimes, nil
}

func (r *raceTimeService) CreateOrUpdate(
	ctx context.Context,
	raceTimes []*data_cache_entity.RaceTime,
	races []*data_cache_entity.Race,
	raceDateMap map[types.RaceDate][]types.RaceId,
) error {
	taskCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 除外レース
	excludeRaceIdMap := make(map[types.RaceId]struct{})
	for _, race := range races {
		switch race.Class() {
		case types.JumpGrade1, types.JumpGrade2, types.JumpGrade3, types.JumpMaiden, types.JumpOpenClass:
			excludeRaceIdMap[race.RaceId()] = struct{}{}
		}
	}

	urls := r.createRaceTimeUrls(raceTimes, raceDateMap, excludeRaceIdMap)
	if len(urls) == 0 {
		return nil
	}

	raceTimeMap := map[types.RaceDate][]*raw_entity.RaceTime{}
	for _, raceTime := range raceTimes {
		_, ok := raceTimeMap[raceTime.RaceDate()]
		if !ok {
			raceTimeMap[raceTime.RaceDate()] = make([]*raw_entity.RaceTime, 0)
		}
		raceTimeMap[raceTime.RaceDate()] = append(raceTimeMap[raceTime.RaceDate()], r.raceTimeEntityConverter.DataCacheToRaw(raceTime))
	}

	var wg sync.WaitGroup
	const raceIdParallel = 5
	errorCh := make(chan error, 1)
	resultCh := make(chan []*netkeiba_entity.RaceTime, raceIdParallel)
	chunkSize := (len(urls) + raceIdParallel - 1) / raceIdParallel

	for i := 0; i < len(urls); i += chunkSize {
		end := i + chunkSize
		if end > len(urls) {
			end = len(urls)
		}

		wg.Add(1)
		go func(splitUrls []string) {
			defer wg.Done()
			localRaceTimes := make([]*netkeiba_entity.RaceTime, 0, len(splitUrls))
			r.logger.Infof("race time fetch processing: %v/%v", end, len(urls))
			for _, url := range splitUrls {
				time.Sleep(time.Millisecond)
				select {
				case <-taskCtx.Done():
					return
				default:
					raceTime, err := r.raceTimeRepository.Fetch(taskCtx, url)
					if err != nil {
						select {
						case errorCh <- err:
							cancel()
						}
						return
					}
					localRaceTimes = append(localRaceTimes, raceTime)
				}
			}

			resultCh <- localRaceTimes
		}(urls[i:end])
	}

	wg.Wait()
	close(errorCh)
	close(resultCh)

	if err := <-errorCh; err != nil {
		return err
	}

	for results := range resultCh {
		for _, raceTime := range results {
			rawRaceTime := r.raceTimeEntityConverter.NetKeibaToRaw(raceTime)
			raceTimeMap[types.RaceDate(rawRaceTime.RaceDate)] = append(raceTimeMap[types.RaceDate(rawRaceTime.RaceDate)], rawRaceTime)
		}
	}

	for raceDate, rawRaceTimes := range raceTimeMap {
		sort.Slice(rawRaceTimes, func(i, j int) bool {
			return rawRaceTimes[i].RaceId < rawRaceTimes[j].RaceId
		})
		raceTimeInfo := raw_entity.RaceTimeInfo{
			RaceTimes: rawRaceTimes,
		}
		err := r.raceTimeRepository.Write(ctx, fmt.Sprintf("%s/race_times/%s", config.CacheDir, fmt.Sprintf(raceTimeFileName, raceDate.Value())), &raceTimeInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *raceTimeService) createRaceTimeUrls(
	raceTimes []*data_cache_entity.RaceTime,
	raceDateMap map[types.RaceDate][]types.RaceId,
	excludeRaceIdMap map[types.RaceId]struct{},
) []string {
	var raceTimeUrls []string

	raceTimeMap := map[types.RaceId]*data_cache_entity.RaceTime{}
	for _, raceTime := range raceTimes {
		raceTimeMap[raceTime.RaceId()] = raceTime
	}

	startDate, err := types.NewRaceDate(config.RaceTimeStartDate)
	if err != nil {
		r.logger.Errorf("failed to create race date: %v", err)
		return nil
	}

	endDate, err := types.NewRaceDate(config.RaceTimeEndDate)
	if err != nil {
		r.logger.Errorf("failed to create race date: %v", err)
		return nil
	}

	raceIdMap := map[types.RaceId]types.RaceDate{}
	for raceDate, raceIds := range raceDateMap {
		if raceDate < startDate || raceDate > endDate {
			continue
		}

		for _, raceId := range raceIds {
			raceIdMap[raceId] = raceDate
		}
	}

	for _, raceId := range converter.SortedRaceIdKeys(raceIdMap) {
		if _, ok := raceTimeMap[raceId]; !ok {
			runes := []rune(raceId.String())
			rawRaceCourseId := string(runes[4:6])
			raceCourse := types.RaceCourse(rawRaceCourseId)
			if _, ok := excludeRaceIdMap[raceId]; !ok && raceCourse.JRA() {
				raceTimeUrls = append(raceTimeUrls, fmt.Sprintf(raceDBUrl, raceId.String()))
			}
		}
	}

	return raceTimeUrls
}
