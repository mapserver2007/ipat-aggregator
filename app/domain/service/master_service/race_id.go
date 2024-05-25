package master_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/config"
	net_url "net/url"
	"sort"
	"time"
)

const (
	raceListUrlForJRA = "https://race.netkeiba.com/top/race_list_sub.html?kaisai_date=%d"
	raceIdFileName    = "race_id.json"
)

type RaceId interface {
	Get(ctx context.Context) (map[types.RaceDate][]types.RaceId, []types.RaceDate, error)
	CreateOrUpdate(ctx context.Context, startDate, endDate string) error
	Update(ctx context.Context, raceDateMapForNAROrOversea map[types.RaceDate][]types.RaceId) error
}

type raceIdService struct {
	raceIdRepository repository.RaceIdRepository
}

func NewRaceId(
	raceIdRepository repository.RaceIdRepository,
) RaceId {
	return &raceIdService{
		raceIdRepository: raceIdRepository,
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

	for _, url := range urls {
		time.Sleep(time.Millisecond)
		u, err := net_url.Parse(url)
		if err != nil {
			return err
		}
		date, err := types.NewRaceDate(u.Query().Get("kaisai_date"))
		if err != nil {
			return err
		}

		rawRaceIds, err := r.raceIdRepository.Fetch(ctx, url)
		if err != nil {
			return err
		}
		if len(rawRaceIds) == 0 {
			newRawExcludeDates = append(newRawExcludeDates, date.Value())
			continue
		}

		rawRaceDate := raw_entity.RaceDate{
			RaceDate: date.Value(),
			RaceIds:  rawRaceIds,
		}
		newRawRaceDates = append(newRawRaceDates, &rawRaceDate)
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
