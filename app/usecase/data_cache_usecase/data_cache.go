package data_cache_usecase

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"log"
	"time"
)

const (
	racingNumberFileName = "racing_number.json"
	raceResultFileName   = "race_result.json"
	jockeyFileName       = "jockey.json"
)

type dataCacheUseCase struct {
	racingNumberDataRepository  repository.RacingNumberDataRepository
	raceDataRepository          repository.RaceDataRepository
	jockeyDataRepository        repository.JockeyDataRepository
	netKeibaService             service.NetKeibaService
	raceConverter               service.RaceConverter
	racingNumberEntityConverter service.RacingNumberEntityConverter
	raceEntityConverter         service.RaceEntityConverter
	jockeyEntityConverter       service.JockeyEntityConverter
}

func NewDataCacheUseCase(
	racingNumberRepository repository.RacingNumberDataRepository,
	raceDataRepository repository.RaceDataRepository,
	jockeyDataRepository repository.JockeyDataRepository,
	netKeibaService service.NetKeibaService,
	raceConverter service.RaceConverter,
	racingNumberConverter service.RacingNumberEntityConverter,
	raceEntityConverter service.RaceEntityConverter,
	jockeyEntityConverter service.JockeyEntityConverter,
) *dataCacheUseCase {
	return &dataCacheUseCase{
		racingNumberDataRepository:  racingNumberRepository,
		raceDataRepository:          raceDataRepository,
		jockeyDataRepository:        jockeyDataRepository,
		netKeibaService:             netKeibaService,
		raceConverter:               raceConverter,
		racingNumberEntityConverter: racingNumberConverter,
		raceEntityConverter:         raceEntityConverter,
		jockeyEntityConverter:       jockeyEntityConverter,
	}
}

func (d *dataCacheUseCase) Read(ctx context.Context) ([]*data_cache_entity.RacingNumber, []*data_cache_entity.Race, []*data_cache_entity.Jockey, []int, error) {
	rawRacingNumbers, err := d.racingNumberDataRepository.Read(ctx, racingNumberFileName)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	racingNumbers := make([]*data_cache_entity.RacingNumber, 0, len(rawRacingNumbers))
	for _, rawRacingNumber := range rawRacingNumbers {
		racingNumbers = append(racingNumbers, d.racingNumberEntityConverter.RawToDataCache(rawRacingNumber))
	}

	rawRaces, err := d.raceDataRepository.Read(ctx, raceResultFileName)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	races := make([]*data_cache_entity.Race, 0, len(rawRaces))
	for _, rawRace := range rawRaces {
		races = append(races, d.raceEntityConverter.RawToDataCache(rawRace))
	}

	rawJockeys, excludeJockeyIds, err := d.jockeyDataRepository.Read(ctx, jockeyFileName)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	jockeys := make([]*data_cache_entity.Jockey, 0, len(rawJockeys))
	for _, rawJockey := range rawJockeys {
		jockeys = append(jockeys, d.jockeyEntityConverter.RawToDataCache(rawJockey))
	}

	return racingNumbers, races, jockeys, excludeJockeyIds, nil
}

func (d *dataCacheUseCase) Write(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	racingNumbers []*data_cache_entity.RacingNumber,
	races []*data_cache_entity.Race,
	jockeys []*data_cache_entity.Jockey,
	excludeJockeyIds []int,
) error {
	urls, _ := d.netKeibaService.CreateRacingNumberUrls(ctx, tickets, racingNumbers)
	newRawRacingNumbers := make([]*raw_entity.RacingNumber, 0, len(racingNumbers)+len(urls))
	for _, racingNumber := range racingNumbers {
		newRawRacingNumbers = append(newRawRacingNumbers, d.racingNumberEntityConverter.DataCacheToRaw(racingNumber))
	}
	for _, url := range urls {
		time.Sleep(time.Second * 1)
		log.Println(ctx, "fetch from "+url)
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

	urls, err = d.netKeibaService.CreateRaceUrls(ctx, tickets, races, racingNumbers)
	if err != nil {
		return err
	}
	newRaces := make([]*raw_entity.Race, 0, len(urls))
	for _, race := range races {
		newRaces = append(newRaces, d.raceEntityConverter.DataCacheToRaw(race))
	}
	ticketMap := d.raceConverter.ConvertToTicketMap(ctx, tickets, racingNumbers)
	for _, url := range urls {
		time.Sleep(time.Second * 1)
		log.Println(ctx, "fetch from "+url)
		fetchRace, err := d.raceDataRepository.Fetch(ctx, url)
		if err != nil {
			return err
		}

		ticket, ok := ticketMap[types.RaceId(fetchRace.RaceId())]
		if !ok {
			return fmt.Errorf("undefind raceId: %v", fetchRace.RaceId())
		}

		newRace := d.raceEntityConverter.NetKeibaToRaw(fetchRace, ticket)
		newRaces = append(newRaces, newRace)
		races = append(races, d.raceEntityConverter.RawToDataCache(newRace))
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
	newJockeys := make([]*raw_entity.Jockey, 0, len(urls))
	for _, jockey := range jockeys {
		newJockeys = append(newJockeys, d.jockeyEntityConverter.DataCacheToRaw(jockey))
	}
	var newExcludeJockeyIds []int
	for _, url := range urls {
		time.Sleep(time.Second * 1)
		log.Println(ctx, "fetch from "+url)
		jockey, err := d.jockeyDataRepository.Fetch(ctx, url)
		if err != nil {
			return err
		}
		if jockey.Name() == "" {
			newExcludeJockeyIds = append(newExcludeJockeyIds, jockey.Id())
		} else {
			newJockey := d.jockeyEntityConverter.NetKeibaToRaw(jockey)
			newJockeys = append(newJockeys, newJockey)
			jockeys = append(jockeys, d.jockeyEntityConverter.RawToDataCache(newJockey))
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

	return nil
}
