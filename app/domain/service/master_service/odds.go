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
	neturl "net/url"
	"sort"
	"time"
)

const (
	oddsUrl      = "https://race.netkeiba.com/api/api_get_jra_odds.html?race_id=%s&type=%d&action=update"
	oddsFileName = "odds_%d.json"
)

type Odds interface {
	Get(ctx context.Context) ([]*data_cache_entity.Odds, error)
	CreateOrUpdate(ctx context.Context, odds []*data_cache_entity.Odds, markers []*marker_csv_entity.AnalysisMarker) error
}

type oddsService struct {
	oddsRepository      repository.OddsRepository
	oddsEntityConverter converter.OddsEntityConverter
}

func NewOdds(
	oddsRepository repository.OddsRepository,
	oddsEntityConverter converter.OddsEntityConverter,
) Odds {
	return &oddsService{
		oddsRepository:      oddsRepository,
		oddsEntityConverter: oddsEntityConverter,
	}
}

func (o *oddsService) Get(ctx context.Context) ([]*data_cache_entity.Odds, error) {
	files, err := o.oddsRepository.List(ctx, fmt.Sprintf("%s/odds", config.CacheDir))
	if err != nil {
		return nil, err
	}

	var odds []*data_cache_entity.Odds
	for _, file := range files {
		rawRaceOddsList, err := o.oddsRepository.Read(ctx, fmt.Sprintf("%s/odds/%s", config.CacheDir, file))
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

func (o *oddsService) CreateOrUpdate(
	ctx context.Context,
	odds []*data_cache_entity.Odds,
	markers []*marker_csv_entity.AnalysisMarker,
) error {
	urls := o.createOddsUrls(odds, markers)
	if len(urls) == 0 {
		return nil
	}

	markerHorseNumberMap := map[types.RaceId][]int{}
	for _, marker := range markers {
		horseNumbers := make([]int, 0, len(marker.MarkerMap()))
		for _, horseNumber := range marker.MarkerMap() {
			horseNumbers = append(horseNumbers, horseNumber)
		}
		markerHorseNumberMap[marker.RaceId()] = horseNumbers
	}

	oddsMap := o.createOddsMap(odds)

	for _, url := range urls {
		newOdds := make([]*raw_entity.Odds, 0, 20) // 6頭BOXは20点
		time.Sleep(time.Millisecond)
		fetchOdds, err := o.oddsRepository.Fetch(ctx, url)
		if err != nil {
			return err
		}

		raceId, err := o.parseUrl(url)
		if err != nil {
			return err
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

	for _, raceDate := range service.SortedRaceDateKeys(oddsMap) {
		rawRaceOddsList := oddsMap[raceDate]
		raceOddsInfo := raw_entity.RaceOddsInfo{
			RaceOdds: rawRaceOddsList,
		}
		err := o.oddsRepository.Write(ctx, fmt.Sprintf("%s/odds/%s", config.CacheDir, fmt.Sprintf(oddsFileName, raceDate.Value())), &raceOddsInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *oddsService) createOddsUrls(
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
			trioOddsUrls = append(trioOddsUrls, fmt.Sprintf(oddsUrl, marker.RaceId(), 7))
		}
	}

	return trioOddsUrls
}

func (o *oddsService) createOddsMap(
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
			RaceId: raceId.String(),
			Odds:   rawOddsList,
		})
	}

	return oddsMap
}

func (o *oddsService) parseUrl(
	url string,
) (types.RaceId, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return "", err
	}
	raceId := u.Query().Get("race_id")

	return types.RaceId(raceId), nil
}

func (o *oddsService) containsInSliceAll(slice1, slice2 []int) bool {
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