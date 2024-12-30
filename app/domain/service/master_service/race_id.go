package master_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/config"
	"github.com/sirupsen/logrus"
	net_url "net/url"
	"sort"
	"sync"
	"time"
)

const (
	raceListUrlForJRA = "https://race.netkeiba.com/top/race_list_sub.html?kaisai_date=%d"
	raceIdFileName    = "race_id.json"
)

// UMACA経由で海外レースを購入した場合の対応
// 数自体少ないので個別に管理する
var overseasRaceDates = map[types.RaceDate][]types.RaceId{
	20241208: {"2024H1120805", "2024H1120808"},
}

type RaceId interface {
	Get(ctx context.Context) (map[types.RaceDate][]types.RaceId, []types.RaceDate, error)
	CreateOrUpdate(ctx context.Context, startDate, endDate string) error
	Update(ctx context.Context, raceDateMapForNAROrOversea map[types.RaceDate][]types.RaceId) error
}

type raceIdService struct {
	raceIdRepository repository.RaceIdRepository
	logger           *logrus.Logger
}

func NewRaceId(
	raceIdRepository repository.RaceIdRepository,
	logger *logrus.Logger,
) RaceId {
	return &raceIdService{
		raceIdRepository: raceIdRepository,
		logger:           logger,
	}
}

func (r *raceIdService) Get(ctx context.Context) (map[types.RaceDate][]types.RaceId, []types.RaceDate, error) {
	var (
		raceDateMap  map[types.RaceDate][]types.RaceId
		excludeDates []types.RaceDate
	)

	rawRaceInfo, err := r.raceIdRepository.Read(ctx, fmt.Sprintf("%s/%s", config.CacheDir, raceIdFileName))
	if err != nil {
		return nil, nil, err
	}
	if rawRaceInfo != nil {
		raceDateMap = map[types.RaceDate][]types.RaceId{}
		for _, rawRaceDate := range rawRaceInfo.RaceDates {
			var raceIds []types.RaceId
			for _, rawRaceId := range rawRaceDate.RaceIds {
				raceIds = append(raceIds, types.RaceId(rawRaceId))
			}
			raceDate := types.RaceDate(rawRaceDate.RaceDate)
			overseasRaceIds, ok := overseasRaceDates[raceDate]
			if ok {
				raceIds = append(raceIds, overseasRaceIds...)
			}
			raceDateMap[raceDate] = raceIds
		}
		excludeDates = make([]types.RaceDate, 0, len(rawRaceInfo.ExcludeDates))
		for _, rawExcludeDate := range rawRaceInfo.ExcludeDates {
			excludeDates = append(excludeDates, types.RaceDate(rawExcludeDate))
		}
	}

	return raceDateMap, excludeDates, nil
}

func (r *raceIdService) CreateOrUpdate(
	ctx context.Context,
	startDate, endDate string,
) error {
	taskCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	newRawRaceDates := make([]*raw_entity.RaceDate, 0)
	newRawExcludeDates := make([]int, 0)

	raceDateMap, excludeDates, err := r.Get(ctx)
	if err != nil {
		return err
	}

	urls, err := r.createRaceIdUrls(startDate, endDate, raceDateMap, excludeDates)
	if err != nil {
		return err
	}

	if len(urls) == 0 {
		return nil
	}

	for _, raceDate := range converter.SortedRaceDateKeys(raceDateMap) {
		raceIds := raceDateMap[raceDate]
		rawRaceIds := make([]string, 0, len(raceIds))
		for _, raceId := range raceIds {
			rawRaceIds = append(rawRaceIds, raceId.String())
		}
		newRawRaceDates = append(newRawRaceDates, &raw_entity.RaceDate{
			RaceDate: raceDate.Value(),
			RaceIds:  rawRaceIds,
		})
	}
	for _, excludeDate := range excludeDates {
		newRawExcludeDates = append(newRawExcludeDates, excludeDate.Value())
	}

	var wg sync.WaitGroup
	const raceIdParallel = 5
	errorCh := make(chan error, 1)
	resultCh1 := make(chan []*raw_entity.RaceDate, raceIdParallel)
	resultCh2 := make(chan []int, raceIdParallel)
	chunkSize := (len(urls) + raceIdParallel - 1) / raceIdParallel

	for i := 0; i < len(urls); i += chunkSize {
		end := i + chunkSize
		if end > len(urls) {
			end = len(urls)
		}

		wg.Add(1)
		go func(splitUrls []string) {
			defer wg.Done()
			localNewRawRaceDates := make([]*raw_entity.RaceDate, 0)
			localNewRawExcludeDates := make([]int, 0)

			r.logger.Infof("raceId fetch processing: %v/%v", end, len(urls))
			for _, url := range splitUrls {
				time.Sleep(time.Millisecond)
				select {
				case <-taskCtx.Done():
					return
				default:
					u, err := net_url.Parse(url)
					if err != nil {
						select {
						case errorCh <- err:
							cancel()
						}
						return
					}
					date, err := types.NewRaceDate(u.Query().Get("kaisai_date"))
					if err != nil {
						select {
						case errorCh <- err:
							cancel()
						}
						return
					}

					rawRaceIds, err := r.raceIdRepository.Fetch(taskCtx, url)
					if err != nil {
						select {
						case errorCh <- err:
							cancel()
						}
						return
					}

					if len(rawRaceIds) == 0 {
						localNewRawExcludeDates = append(localNewRawExcludeDates, date.Value())
					} else {
						rawRaceDate := raw_entity.RaceDate{
							RaceDate: date.Value(),
							RaceIds:  rawRaceIds,
						}
						localNewRawRaceDates = append(localNewRawRaceDates, &rawRaceDate)
					}
				}
			}

			resultCh1 <- localNewRawRaceDates
			resultCh2 <- localNewRawExcludeDates
		}(urls[i:end])
	}

	wg.Wait()
	close(errorCh)
	close(resultCh1)
	close(resultCh2)

	if err := <-errorCh; err != nil {
		return err
	}

	for results := range resultCh1 {
		newRawRaceDates = append(newRawRaceDates, results...)
	}

	for results := range resultCh2 {
		newRawExcludeDates = append(newRawExcludeDates, results...)
	}

	sort.Slice(newRawRaceDates, func(i, j int) bool {
		return newRawRaceDates[i].RaceDate < newRawRaceDates[j].RaceDate
	})

	sort.Ints(newRawExcludeDates)

	err = r.raceIdRepository.Write(ctx, fmt.Sprintf("%s/%s", config.CacheDir, raceIdFileName), &raw_entity.RaceIdInfo{
		RaceDates:    newRawRaceDates,
		ExcludeDates: newRawExcludeDates,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *raceIdService) Update(ctx context.Context, raceDateMapForNAROrOversea map[types.RaceDate][]types.RaceId) error {
	newRawRaceDates := make([]*raw_entity.RaceDate, 0)
	newRawExcludeDates := make([]int, 0)

	raceDateMap, excludeDates, err := r.Get(ctx)
	if err != nil {
		return err
	}

	for raceDate, narOrOverseaRaceIds := range raceDateMapForNAROrOversea {
		if _, ok := raceDateMap[raceDate]; !ok {
			raceDateMap[raceDate] = narOrOverseaRaceIds
		} else {
			raceDateMap[raceDate] = append(raceDateMap[raceDate], narOrOverseaRaceIds...)
		}
	}

	for _, raceDate := range converter.SortedRaceDateKeys(raceDateMap) {
		raceIdMap := map[types.RaceId]bool{}
		raceIds := raceDateMap[raceDate]
		for _, raceId := range raceIds {
			raceIdMap[raceId] = true
		}

		rawRaceIds := make([]string, 0, len(raceIdMap))
		for raceId := range raceIdMap {
			rawRaceIds = append(rawRaceIds, raceId.String())
		}

		sort.Strings(rawRaceIds)
		newRawRaceDates = append(newRawRaceDates, &raw_entity.RaceDate{
			RaceDate: raceDate.Value(),
			RaceIds:  rawRaceIds,
		})
	}

	for _, excludeDate := range excludeDates {
		// excludeDatesに含まれている日付はJRA開催のもののみなので、地方・海外の日付は除外しない
		newRawExcludeDates = append(newRawExcludeDates, excludeDate.Value())
	}

	sort.Slice(newRawRaceDates, func(i, j int) bool {
		return newRawRaceDates[i].RaceDate < newRawRaceDates[j].RaceDate
	})

	sort.Ints(newRawExcludeDates)

	err = r.raceIdRepository.Write(ctx, fmt.Sprintf("%s/%s", config.CacheDir, raceIdFileName), &raw_entity.RaceIdInfo{
		RaceDates:    newRawRaceDates,
		ExcludeDates: newRawExcludeDates,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *raceIdService) createRaceIdUrls(
	dateFrom, dateTo string,
	raceIdMap map[types.RaceDate][]types.RaceId,
	excludeDates []types.RaceDate,
) ([]string, error) {
	urls := make([]string, 0)
	excludeDateMap := map[types.RaceDate]bool{}
	for _, excludeDate := range excludeDates {
		excludeDateMap[excludeDate] = true
	}

	startTime, _ := time.Parse("20060102", dateFrom)
	endTime, _ := time.Parse("20060102", dateTo)
	for d := startTime; d.Before(endTime) || d.Equal(endTime); d = d.AddDate(0, 0, 1) {
		date, err := types.NewRaceDate(d.Format("20060102"))
		if err != nil {
			return nil, err
		}
		if excludeDateMap != nil {
			if _, ok := excludeDateMap[date]; ok {
				continue
			}
		}
		if raceIdMap == nil {
			urls = append(urls, fmt.Sprintf(raceListUrlForJRA, date))
		} else {
			if _, ok := raceIdMap[date]; !ok {
				urls = append(urls, fmt.Sprintf(raceListUrlForJRA, date))
			}
		}
	}

	return urls, nil
}
