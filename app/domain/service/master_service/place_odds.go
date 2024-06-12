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
	placeOddsUrl      = "https://race.netkeiba.com/api/api_get_jra_odds.html?race_id=%s&type=2&action=update"
	placeOddsFileName = "odds_%d.json"
)

type PlaceOdds interface {
	Get(ctx context.Context) ([]*data_cache_entity.Odds, error)
	CreateOrUpdate(ctx context.Context, odds []*data_cache_entity.Odds, markers []*marker_csv_entity.AnalysisMarker) error
}

type placeOddsService struct {
	oddsRepository      repository.OddsRepository
	oddsEntityConverter converter.OddsEntityConverter
}

func NewPlaceOdds(
	oddsRepository repository.OddsRepository,
	oddsEntityConverter converter.OddsEntityConverter,
) PlaceOdds {
	return &placeOddsService{
		oddsRepository:      oddsRepository,
		oddsEntityConverter: oddsEntityConverter,
	}
}

func (p *placeOddsService) Get(ctx context.Context) ([]*data_cache_entity.Odds, error) {
	files, err := p.oddsRepository.List(ctx, fmt.Sprintf("%s/odds/place", config.CacheDir))
	if err != nil {
		return nil, err
	}

	var odds []*data_cache_entity.Odds
	for _, file := range files {
		rawRaceOddsList, err := p.oddsRepository.Read(ctx, fmt.Sprintf("%s/odds/place/%s", config.CacheDir, file))
		if err != nil {
			return nil, err
		}
		for _, rawRaceOdds := range rawRaceOddsList {
			raceId := types.RaceId(rawRaceOdds.RaceId)
			raceDate := types.RaceDate(rawRaceOdds.RaceDate)
			for _, rawOdds := range rawRaceOdds.Odds {
				odds = append(odds, p.oddsEntityConverter.RawToDataCache(rawOdds, raceId, raceDate))
			}
		}
	}

	return odds, nil
}

func (p *placeOddsService) CreateOrUpdate(
	ctx context.Context,
	odds []*data_cache_entity.Odds,
	markers []*marker_csv_entity.AnalysisMarker,
) error {
	urls := p.createOddsUrls(odds, markers)
	if len(urls) == 0 {
		return nil
	}

	oddsMap := p.createOddsMap(odds)

	for _, url := range urls {
		var newOdds []*raw_entity.Odds
		time.Sleep(time.Millisecond)
		fetchOdds, err := p.oddsRepository.Fetch(ctx, url)
		if err != nil {
			return err
		}

		raceId, err := p.parseUrl(url)
		if err != nil {
			return err
		}

		var raceDate types.RaceDate
		if len(fetchOdds) > 0 {
			raceDate = fetchOdds[0].RaceDate()
		}

		for _, netKeibaFetchOdds := range fetchOdds {
			newOdds = append(newOdds, p.oddsEntityConverter.NetKeibaToRaw(netKeibaFetchOdds))
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
		err := p.oddsRepository.Write(ctx, fmt.Sprintf("%s/odds/place/%s", config.CacheDir, fmt.Sprintf(placeOddsFileName, raceDate.Value())), &raceOddsInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *placeOddsService) createOddsUrls(
	oddsList []*data_cache_entity.Odds,
	markers []*marker_csv_entity.AnalysisMarker,
) []string {
	var placeOddsUrls []string
	oddsMap := map[types.RaceId]bool{}

	for _, odds := range oddsList {
		if _, ok := oddsMap[odds.RaceId()]; !ok {
			oddsMap[odds.RaceId()] = true
		}
	}

	for _, marker := range markers {
		if _, ok := oddsMap[marker.RaceId()]; !ok {
			placeOddsUrls = append(placeOddsUrls, fmt.Sprintf(placeOddsUrl, marker.RaceId()))
		}
	}

	return placeOddsUrls
}

func (p *placeOddsService) createOddsMap(
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
			rawOddsList = append(rawOddsList, p.oddsEntityConverter.DataCacheToRaw(odds))
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

func (p *placeOddsService) parseUrl(
	url string,
) (types.RaceId, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return "", err
	}
	raceId := u.Query().Get("race_id")

	return types.RaceId(raceId), nil
}
