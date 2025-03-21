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
	raceResultUrlForJRA     = "https://race.netkeiba.com/race/result.html?race_id=%s&organizer=1&race_date=%d"
	raceResultUrlForNAR     = "https://nar.netkeiba.com/race/result.html?race_id=%s&organizer=2&race_date=%d"
	raceResultUrlForOversea = "https://race.netkeiba.com/race/result.html?race_id=%s&organizer=3&race_date=%d"
	raceFileName            = "race_%d.json"
)

type Race interface {
	Get(ctx context.Context) ([]*data_cache_entity.Race, error)
	CreateOrUpdate(
		ctx context.Context,
		races []*data_cache_entity.Race,
		raceDateMap map[types.RaceDate][]types.RaceId,
	) error
}

type raceService struct {
	raceRepository      repository.RaceRepository
	raceEntityConverter converter.RaceEntityConverter
	logger              *logrus.Logger
}

func NewRace(
	raceRepository repository.RaceRepository,
	raceEntityConverter converter.RaceEntityConverter,
	logger *logrus.Logger,
) Race {
	return &raceService{
		raceRepository:      raceRepository,
		raceEntityConverter: raceEntityConverter,
		logger:              logger,
	}
}

func (r *raceService) Get(ctx context.Context) ([]*data_cache_entity.Race, error) {
	files, err := r.raceRepository.List(ctx, fmt.Sprintf("%s/races", config.CacheDir))
	if err != nil {
		return nil, err
	}

	var races []*data_cache_entity.Race
	for _, file := range files {
		rawRaces, err := r.raceRepository.Read(ctx, fmt.Sprintf("%s/races/%s", config.CacheDir, file))
		if err != nil {
			return nil, err
		}
		for _, rawRace := range rawRaces {
			races = append(races, r.raceEntityConverter.RawToDataCache(rawRace))
		}
	}

	return races, nil
}

func (r *raceService) CreateOrUpdate(
	ctx context.Context,
	races []*data_cache_entity.Race,
	raceDateMap map[types.RaceDate][]types.RaceId,
) error {
	taskCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	urls := r.createRaceUrls(races, raceDateMap)
	if len(urls) == 0 {
		return nil
	}

	raceMap := map[types.RaceDate][]*raw_entity.Race{}
	for _, race := range races {
		_, ok := raceMap[race.RaceDate()]
		if !ok {
			raceMap[race.RaceDate()] = make([]*raw_entity.Race, 0)
		}
		raceMap[race.RaceDate()] = append(raceMap[race.RaceDate()], r.raceEntityConverter.DataCacheToRaw(race))
	}

	var wg sync.WaitGroup
	const raceIdParallel = 5
	errorCh := make(chan error, 1)
	resultCh := make(chan []*netkeiba_entity.Race, raceIdParallel)
	chunkSize := (len(urls) + raceIdParallel - 1) / raceIdParallel

	for i := 0; i < len(urls); i += chunkSize {
		end := i + chunkSize
		if end > len(urls) {
			end = len(urls)
		}

		wg.Add(1)
		go func(splitUrls []string) {
			defer wg.Done()
			localRaces := make([]*netkeiba_entity.Race, 0, len(splitUrls))
			r.logger.Infof("race fetch processing: %v/%v", end, len(urls))
			for _, url := range splitUrls {
				time.Sleep(time.Millisecond)
				select {
				case <-taskCtx.Done():
					return
				default:
					race, err := r.raceRepository.FetchRace(ctx, url)
					if err != nil {
						select {
						case errorCh <- err:
							cancel()
						}
						return
					}
					localRaces = append(localRaces, race)
				}
			}

			resultCh <- localRaces
		}(urls[i:end])
	}

	wg.Wait()
	close(errorCh)
	close(resultCh)

	if err := <-errorCh; err != nil {
		return err
	}

	for results := range resultCh {
		for _, race := range results {
			rawRace := r.raceEntityConverter.NetKeibaToRaw(race)
			raceMap[types.RaceDate(rawRace.RaceDate)] = append(raceMap[types.RaceDate(rawRace.RaceDate)], rawRace)
		}
	}

	for raceDate, rawRaces := range raceMap {
		sort.Slice(rawRaces, func(i, j int) bool {
			return rawRaces[i].RaceId < rawRaces[j].RaceId
		})
		raceInfo := raw_entity.RaceInfo{
			Races: rawRaces,
		}
		err := r.raceRepository.Write(ctx, fmt.Sprintf("%s/races/%s", config.CacheDir, fmt.Sprintf(raceFileName, raceDate.Value())), &raceInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *raceService) createRaceUrls(
	races []*data_cache_entity.Race,
	raceDateMap map[types.RaceDate][]types.RaceId,
) []string {
	var raceUrls []string

	raceMap := map[types.RaceId]*data_cache_entity.Race{}
	for _, race := range races {
		raceMap[race.RaceId()] = race
	}

	raceIdMap := map[types.RaceId]types.RaceDate{}
	for raceDate, raceIds := range raceDateMap {
		for _, raceId := range raceIds {
			raceIdMap[raceId] = raceDate
		}
	}

	for _, raceId := range converter.SortedRaceIdKeys(raceIdMap) {
		if _, ok := raceMap[raceId]; !ok {
			runes := []rune(raceId.String())
			rawRaceCourseId := string(runes[4:6])
			raceCourse := types.RaceCourse(rawRaceCourseId)
			if raceCourse.JRA() {
				raceUrls = append(raceUrls, fmt.Sprintf(raceResultUrlForJRA, raceId, raceIdMap[raceId]))
			} else if raceCourse.NAR() {
				raceUrls = append(raceUrls, fmt.Sprintf(raceResultUrlForNAR, raceId, raceIdMap[raceId]))
			} else if raceCourse.Oversea() {
				raceUrls = append(raceUrls, fmt.Sprintf(raceResultUrlForOversea, raceId, raceIdMap[raceId]))
			}
		}
	}

	return raceUrls
}
