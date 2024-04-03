package data_cache_usecase

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"log"
	net_url "net/url"
	"time"
)

const (
	racingNumberFileName       = "racing_number.json"
	raceResultFileName         = "race_result.json"
	jockeyFileName             = "jockey.json"
	raceIdFileName             = "race_id.json"
	analysisRaceResultFilePath = "races/race_result_%d.json"
	analysisRaceOddsFilePath   = "odds/odds_%d.json"
)

type DataCacheUseCase struct {
	racingNumberDataRepository  repository.RacingNumberDataRepository
	raceDataRepository          repository.RaceDataRepository
	jockeyDataRepository        repository.JockeyDataRepository
	raceIdDataRepository        repository.RaceIdDataRepository
	markerDataRepository        repository.MarkerDataRepository
	oddsDataRepository          repository.OddsDataRepository
	netKeibaService             service.NetKeibaService
	raceConverter               service.RaceConverter
	racingNumberEntityConverter service.RacingNumberEntityConverter
	raceEntityConverter         service.RaceEntityConverter
	jockeyEntityConverter       service.JockeyEntityConverter
	oddsEntityConverter         service.OddsEntityConverter
}

func NewDataCacheUseCase(
	racingNumberRepository repository.RacingNumberDataRepository,
	raceDataRepository repository.RaceDataRepository,
	jockeyDataRepository repository.JockeyDataRepository,
	raceIdDataRepository repository.RaceIdDataRepository,
	markerDataRepository repository.MarkerDataRepository,
	oddsDataRepository repository.OddsDataRepository,
	netKeibaService service.NetKeibaService,
	raceConverter service.RaceConverter,
	racingNumberConverter service.RacingNumberEntityConverter,
	raceEntityConverter service.RaceEntityConverter,
	jockeyEntityConverter service.JockeyEntityConverter,
	oddsEntityConverter service.OddsEntityConverter,
) *DataCacheUseCase {
	return &DataCacheUseCase{
		racingNumberDataRepository:  racingNumberRepository,
		raceDataRepository:          raceDataRepository,
		jockeyDataRepository:        jockeyDataRepository,
		raceIdDataRepository:        raceIdDataRepository,
		markerDataRepository:        markerDataRepository,
		oddsDataRepository:          oddsDataRepository,
		netKeibaService:             netKeibaService,
		raceConverter:               raceConverter,
		racingNumberEntityConverter: racingNumberConverter,
		raceEntityConverter:         raceEntityConverter,
		jockeyEntityConverter:       jockeyEntityConverter,
		oddsEntityConverter:         oddsEntityConverter,
	}
}

func (d *DataCacheUseCase) Read(
	ctx context.Context,
) (
	[]*data_cache_entity.RacingNumber,
	[]*data_cache_entity.Race,
	[]*data_cache_entity.Jockey,
	[]int,
	map[types.RaceDate][]types.RaceId,
	[]types.RaceDate,
	[]*data_cache_entity.Race,
	[]*data_cache_entity.Odds,
	error,
) {
	rawRacingNumbers, err := d.racingNumberDataRepository.Read(ctx, racingNumberFileName)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	racingNumbers := make([]*data_cache_entity.RacingNumber, 0, len(rawRacingNumbers))
	for _, rawRacingNumber := range rawRacingNumbers {
		racingNumbers = append(racingNumbers, d.racingNumberEntityConverter.RawToDataCache(rawRacingNumber))
	}

	rawRaces, err := d.raceDataRepository.Read(ctx, raceResultFileName)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	ticketRaces := make([]*data_cache_entity.Race, 0, len(rawRaces))
	for _, rawRace := range rawRaces {
		ticketRaces = append(ticketRaces, d.raceEntityConverter.RawToDataCache(rawRace))
	}

	rawJockeys, excludeJockeyIds, err := d.jockeyDataRepository.Read(ctx, jockeyFileName)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	jockeys := make([]*data_cache_entity.Jockey, 0, len(rawJockeys))
	for _, rawJockey := range rawJockeys {
		jockeys = append(jockeys, d.jockeyEntityConverter.RawToDataCache(rawJockey))
	}

	rawRaceDates, rawExcludeDates, err := d.raceIdDataRepository.Read(ctx, raceIdFileName)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	raceIdMap := map[types.RaceDate][]types.RaceId{}
	for _, rawRaceDate := range rawRaceDates {
		var raceIds []types.RaceId
		for _, rawRaceId := range rawRaceDate.RaceIds {
			raceIds = append(raceIds, types.RaceId(rawRaceId))
		}
		raceDate := types.RaceDate(rawRaceDate.RaceDate)
		raceIdMap[raceDate] = raceIds
	}
	excludeDates := make([]types.RaceDate, 0, len(rawExcludeDates))
	for _, rawExcludeDate := range rawExcludeDates {
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, err
		}
		excludeDates = append(excludeDates, types.RaceDate(rawExcludeDate))
	}

	analysisRaces := make([]*data_cache_entity.Race, 0)
	for raceDate := range raceIdMap {
		rawPredictRaces, err := d.raceDataRepository.Read(ctx, fmt.Sprintf(analysisRaceResultFilePath, raceDate.Value()))
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, err
		}
		for _, rawPredictRace := range rawPredictRaces {
			analysisRaces = append(analysisRaces, d.raceEntityConverter.RawToDataCache(rawPredictRace))
		}
	}

	analysisOdds := make([]*data_cache_entity.Odds, 0)
	for raceDate := range raceIdMap {
		rawRaceOddsList, err := d.oddsDataRepository.Read(ctx, fmt.Sprintf(analysisRaceOddsFilePath, raceDate.Value()))
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, err
		}
		for _, rawRaceOdds := range rawRaceOddsList {
			raceId := types.RaceId(rawRaceOdds.RaceId)
			for _, rawOdds := range rawRaceOdds.Odds {
				analysisOdds = append(analysisOdds, d.oddsEntityConverter.RawToDataCache(rawOdds, raceId, raceDate))
			}
		}
	}

	return racingNumbers, ticketRaces, jockeys, excludeJockeyIds, raceIdMap, excludeDates, analysisRaces, analysisOdds, nil
}

func (d *DataCacheUseCase) Write(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	racingNumbers []*data_cache_entity.RacingNumber,
	ticketRaces []*data_cache_entity.Race,
	jockeys []*data_cache_entity.Jockey,
	excludeJockeyIds []int,
	raceIdMap map[types.RaceDate][]types.RaceId,
	excludeDates []types.RaceDate,
	analysisRaces []*data_cache_entity.Race,
	analysisOdds []*data_cache_entity.Odds,
	startDate, endDate string,
	markers []*marker_csv_entity.AnalysisMarker,
) error {
	urls, _ := d.netKeibaService.CreateRacingNumberUrls(ctx, tickets, racingNumbers)
	newRawRacingNumbers := make([]*raw_entity.RacingNumber, 0, len(racingNumbers)+len(urls))
	for _, racingNumber := range racingNumbers {
		newRawRacingNumbers = append(newRawRacingNumbers, d.racingNumberEntityConverter.DataCacheToRaw(racingNumber))
	}
	for _, url := range urls {
		time.Sleep(time.Second * 1)
		log.Println(ctx, "fetch racingNumber from "+url)
		fetchRacingNumbers, err := d.racingNumberDataRepository.Fetch(ctx, url)
		if err != nil {
			return err
		}
		for _, fetchRacingNumber := range fetchRacingNumbers {
			newRawRacingNumber := d.racingNumberEntityConverter.NetKeibaToRaw(fetchRacingNumber)
			racingNumbers = append(racingNumbers, d.racingNumberEntityConverter.RawToDataCache(newRawRacingNumber))
			newRawRacingNumbers = append(newRawRacingNumbers, newRawRacingNumber)
		}
	}

	racingNumberInfo := raw_entity.RacingNumberInfo{
		RacingNumbers: newRawRacingNumbers,
	}

	err := d.racingNumberDataRepository.Write(ctx, racingNumberFileName, &racingNumberInfo)
	if err != nil {
		return err
	}

	urls, err = d.netKeibaService.CreateRaceUrls(ctx, tickets, ticketRaces, racingNumbers)
	if err != nil {
		return err
	}
	newRaces := make([]*raw_entity.Race, 0, len(ticketRaces)+len(urls))
	for _, race := range ticketRaces {
		newRaces = append(newRaces, d.raceEntityConverter.DataCacheToRaw(race))
	}
	for _, url := range urls {
		time.Sleep(time.Second * 1)
		log.Println(ctx, "fetch race from "+url)
		fetchRace, err := d.raceDataRepository.Fetch(ctx, url)
		if err != nil {
			return err
		}
		newRace := d.raceEntityConverter.NetKeibaToRaw(fetchRace)
		newRaces = append(newRaces, newRace)
	}

	raceInfo := raw_entity.RaceInfo{
		Races: newRaces,
	}

	err = d.raceDataRepository.Write(ctx, raceResultFileName, &raceInfo)
	if err != nil {
		return err
	}

	urls, err = d.netKeibaService.CreateJockeyUrls(ctx, jockeys, excludeJockeyIds)
	if err != nil {
		return err
	}
	newJockeys := make([]*raw_entity.Jockey, 0, len(jockeys)+len(urls))
	for _, jockey := range jockeys {
		newJockeys = append(newJockeys, d.jockeyEntityConverter.DataCacheToRaw(jockey))
	}
	var newExcludeJockeyIds []int
	for _, url := range urls {
		time.Sleep(time.Second * 1)
		log.Println(ctx, "fetch jockey from "+url)
		jockey, err := d.jockeyDataRepository.Fetch(ctx, url)
		if err != nil {
			return err
		}
		if jockey.Name() == "" {
			newExcludeJockeyIds = append(newExcludeJockeyIds, jockey.Id())
		} else {
			newJockey := d.jockeyEntityConverter.NetKeibaToRaw(jockey)
			newJockeys = append(newJockeys, newJockey)
		}
	}
	excludeJockeyIds = append(excludeJockeyIds, newExcludeJockeyIds...)

	jockeyInfo := raw_entity.JockeyInfo{
		Jockeys:          newJockeys,
		ExcludeJockeyIds: excludeJockeyIds,
	}

	err = d.jockeyDataRepository.Write(ctx, jockeyFileName, &jockeyInfo)
	if err != nil {
		return err
	}

	urls, err = d.netKeibaService.CreateRaceIdUrls(ctx, startDate, endDate, raceIdMap, excludeDates)
	if err != nil {
		return err
	}
	var newRawRaceDates []*raw_entity.RaceDate
	var newRawExcludeDates []int
	for raceDate, raceIds := range raceIdMap {
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
		time.Sleep(time.Second * 1)
		u, err := net_url.Parse(url)
		if err != nil {
			return err
		}
		date, err := types.NewRaceDate(u.Query().Get("kaisai_date"))
		if err != nil {
			return err
		}
		rawRaceIds, err := d.raceIdDataRepository.Fetch(ctx, url)
		if err != nil {
			return err
		}
		if len(rawRaceIds) == 0 {
			newRawExcludeDates = append(newRawExcludeDates, date.Value())
			continue
		}
		log.Println(ctx, "fetch raceId from "+url)

		rawRaceDate := raw_entity.RaceDate{
			RaceDate: date.Value(),
			RaceIds:  rawRaceIds,
		}
		newRawRaceDates = append(newRawRaceDates, &rawRaceDate)
	}

	raceIdInfo := raw_entity.RaceIdInfo{
		RaceDates:    newRawRaceDates,
		ExcludeDates: newRawExcludeDates,
	}

	err = d.raceIdDataRepository.Write(ctx, raceIdFileName, &raceIdInfo)
	if err != nil {
		return err
	}

	markerHorseNumberMap := map[types.RaceId][]int{}
	for _, marker := range markers {
		horseNumbers := make([]int, 0, len(marker.MarkerMap()))
		for _, horseNumber := range marker.MarkerMap() {
			horseNumbers = append(horseNumbers, horseNumber)
		}
		markerHorseNumberMap[marker.RaceId()] = horseNumbers
	}

	enableRaceIdMap := map[types.RaceId]types.RaceDate{}
	for _, rawRaceDate := range raceIdInfo.RaceDates {
		for _, rawRaceId := range rawRaceDate.RaceIds {
			raceId := types.RaceId(rawRaceId)
			if _, ok := markerHorseNumberMap[raceId]; ok {
				enableRaceIdMap[types.RaceId(rawRaceId)] = types.RaceDate(rawRaceDate.RaceDate)
			}
		}
	}

	raceUrls, err := d.netKeibaService.CreateAnalysisRaceUrls(ctx, analysisRaces, enableRaceIdMap)
	if err != nil {
		return err
	}

	oddsUrls, err := d.netKeibaService.CreateOddsUrls(ctx, analysisOdds, enableRaceIdMap)
	if err != nil {
		return err
	}

	raceMap := map[types.RaceDate][]*raw_entity.Race{}
	for _, race := range analysisRaces {
		if _, ok := raceMap[race.RaceDate()]; !ok {
			raceMap[race.RaceDate()] = make([]*raw_entity.Race, 0)
		}
		raceMap[race.RaceDate()] = append(raceMap[race.RaceDate()], d.raceEntityConverter.DataCacheToRaw(race))
	}

	oddsMap := map[types.RaceDate][]*raw_entity.RaceOdds{}

	raceIdOddsMap := map[types.RaceId][]*data_cache_entity.Odds{}
	for _, odds := range analysisOdds {
		if _, ok := raceIdOddsMap[odds.RaceId()]; !ok {
			raceIdOddsMap[odds.RaceId()] = make([]*data_cache_entity.Odds, 0)
		}
		raceIdOddsMap[odds.RaceId()] = append(raceIdOddsMap[odds.RaceId()], odds)
	}
	for raceId, oddsList := range raceIdOddsMap {
		raceDate := oddsList[0].RaceDate()
		rawOddsList := make([]*raw_entity.Odds, 0, len(oddsList))
		for _, odds := range oddsList {
			rawOddsList = append(rawOddsList, d.oddsEntityConverter.DataCacheToRaw(odds))
		}

		if _, ok := oddsMap[raceDate]; !ok {
			oddsMap[raceDate] = make([]*raw_entity.RaceOdds, 0)
		}

		oddsMap[raceDate] = append(oddsMap[raceDate], &raw_entity.RaceOdds{
			RaceId: raceId.String(),
			Odds:   rawOddsList,
		})
	}

	for _, url := range raceUrls {
		time.Sleep(time.Second * 1)
		log.Println(ctx, "fetch analysisRace from "+url)
		fetchRace, err := d.raceDataRepository.Fetch(ctx, url)
		if err != nil {
			return err
		}
		newRace := d.raceEntityConverter.NetKeibaToRaw(fetchRace)
		raceMap[types.RaceDate(newRace.RaceDate)] = append(raceMap[types.RaceDate(newRace.RaceDate)], newRace)
	}

	for raceDate, rawRaces := range raceMap {
		raceInfo = raw_entity.RaceInfo{
			Races: rawRaces,
		}
		filePath := fmt.Sprintf(analysisRaceResultFilePath, raceDate.Value())
		err = d.raceDataRepository.Write(ctx, filePath, &raceInfo)
		if err != nil {
			return err
		}
	}

	for _, url := range oddsUrls {
		newOdds := make([]*raw_entity.Odds, 0, 20) // 6頭BOXは20点
		time.Sleep(time.Second * 1)
		log.Println(ctx, "fetch analysisOdds from "+url)
		fetchOdds, err := d.oddsDataRepository.Fetch(ctx, url)
		if err != nil {
			return err
		}

		u, err := net_url.Parse(url)
		if err != nil {
			return err
		}
		raceId := u.Query().Get("race_id")
		horseNumbers, ok := markerHorseNumberMap[types.RaceId(raceId)]
		if !ok {
			// 分析印がないレース(新馬、障害など)はスキップ
			continue
		}

		var raceDate types.RaceDate
		if len(fetchOdds) > 0 {
			raceDate = fetchOdds[0].RaceDate()
		}

		for _, netKeibaFetchOdds := range fetchOdds {
			if containsInSliceAll(horseNumbers, netKeibaFetchOdds.HorseNumbers()) {
				newOdds = append(newOdds, d.oddsEntityConverter.NetKeibaToRaw(netKeibaFetchOdds))
			}
		}

		if _, ok := oddsMap[raceDate]; !ok {
			oddsMap[raceDate] = make([]*raw_entity.RaceOdds, 0)
		}

		oddsMap[raceDate] = append(oddsMap[raceDate], &raw_entity.RaceOdds{
			RaceId: raceId,
			Odds:   newOdds,
		})
	}

	for raceDate, rawRaceOddsList := range oddsMap {
		raceOddsInfo := raw_entity.RaceOddsInfo{
			RaceOdds: rawRaceOddsList,
		}
		filePath := fmt.Sprintf(analysisRaceOddsFilePath, raceDate.Value())
		err = d.oddsDataRepository.Write(ctx, filePath, &raceOddsInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func containsInSliceAll(slice1, slice2 []int) bool {
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
