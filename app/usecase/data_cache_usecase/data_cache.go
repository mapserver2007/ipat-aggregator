package data_cache_usecase

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"log"
	neturl "net/url"
	"time"
)

const (
	racingNumberFileName = "racing_number.json"
	raceResultFileName   = "race_result.json"
	jockeyFileName       = "jockey.json"
)

type dataCacheUseCase struct {
	racingNumberDataRepository repository.RacingNumberDataRepository
	raceDataRepository         repository.RaceDataRepository
	jockeyDataRepository       repository.JockeyDataRepository
	netKeibaService            service.NetKeibaService
	raceConverter              service.RaceConverter
}

func NewDataCacheUseCase(
	racingNumberRepository repository.RacingNumberDataRepository,
	raceDataRepository repository.RaceDataRepository,
	jockeyDataRepository repository.JockeyDataRepository,
	netKeibaService service.NetKeibaService,
	raceConverter service.RaceConverter,
) *dataCacheUseCase {
	return &dataCacheUseCase{
		racingNumberDataRepository: racingNumberRepository,
		raceDataRepository:         raceDataRepository,
		jockeyDataRepository:       jockeyDataRepository,
		netKeibaService:            netKeibaService,
		raceConverter:              raceConverter,
	}
}

func (d *dataCacheUseCase) Read(ctx context.Context) ([]*raw_entity.RacingNumber, []*raw_entity.Race, []*raw_entity.Jockey, []int, error) {
	rawRacingNumbers, err := d.racingNumberDataRepository.Read(ctx, racingNumberFileName)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	rawRaces, err := d.raceDataRepository.Read(ctx, raceResultFileName)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	rawJockeys, excludeJockeyIds, err := d.jockeyDataRepository.Read(ctx, jockeyFileName)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return rawRacingNumbers, rawRaces, rawJockeys, excludeJockeyIds, nil
}

func (d *dataCacheUseCase) Write(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	rawRacingNumbers []*raw_entity.RacingNumber,
	rawRaces []*raw_entity.Race,
	rawJockeys []*raw_entity.Jockey,
	excludeJockeyIds []int,
) error {
	urls, _ := d.netKeibaService.CreateRacingNumberUrls(ctx, tickets, rawRacingNumbers)
	newRawRacingNumbers := make([]*raw_entity.RacingNumber, 0, len(urls))
	for _, url := range urls {
		time.Sleep(time.Second * 1)
		log.Println(ctx, "fetch from "+url)
		racingNumbers, err := d.racingNumberDataRepository.Fetch(ctx, url)
		if err != nil {
			return err
		}
		for _, racingNumber := range racingNumbers {
			newRawRacingNumbers = append(newRawRacingNumbers, &raw_entity.RacingNumber{
				Date:         racingNumber.Date(),
				Round:        racingNumber.Round(),
				Day:          racingNumber.Day(),
				RaceCourseId: racingNumber.RaceCourseId(),
			})
		}
	}
	rawRacingNumbers = append(rawRacingNumbers, newRawRacingNumbers...)

	racingNumberInfo := raw_entity.RacingNumberInfo{
		RacingNumbers: newRawRacingNumbers,
	}

	err := d.racingNumberDataRepository.Write(ctx, racingNumberFileName, &racingNumberInfo)
	if err != nil {
		return err
	}

	urls, err = d.netKeibaService.CreateRaceUrls(ctx, tickets, rawRaces, rawRacingNumbers)
	if err != nil {
		return err
	}
	newRaces := make([]*raw_entity.Race, 0, len(urls))
	ticketMap := d.raceConverter.ConvertToTicketMap(ctx, tickets, rawRacingNumbers)
	for _, url := range urls {
		time.Sleep(time.Second * 1)
		log.Println(ctx, "fetch from "+url)
		race, err := d.raceDataRepository.Fetch(ctx, url)
		if err != nil {
			return err
		}

		parsedUrl, err := neturl.Parse(url)
		if err != nil {
			return err
		}
		queryParams, err := neturl.ParseQuery(parsedUrl.RawQuery)
		if err != nil {
			return err
		}
		raceId := queryParams.Get("race_id")
		ticket, ok := ticketMap[types.RaceId(raceId)]
		if !ok {
			return fmt.Errorf("undefind raceId: %v", raceId)
		}

		raceResults := make([]*raw_entity.RaceResult, 0, len(race.RaceResults()))
		for _, raceResult := range race.RaceResults() {
			raceResults = append(raceResults, &raw_entity.RaceResult{
				OrderNo:       raceResult.OrderNo(),
				HorseName:     raceResult.HorseName(),
				BracketNumber: raceResult.BracketNumber(),
				HorseNumber:   raceResult.HorseNumber(),
				JockeyId:      raceResult.JockeyId(),
				Odds:          raceResult.Odds(),
				PopularNumber: raceResult.PopularNumber(),
			})
		}
		payoutResults := make([]*raw_entity.PayoutResult, 0, len(race.PayoutResults()))
		for _, payoutResult := range race.PayoutResults() {
			payoutResults = append(payoutResults, &raw_entity.PayoutResult{
				TicketType: payoutResult.TicketType(),
				Numbers:    payoutResult.Numbers(),
				Odds:       payoutResult.Odds(),
				Populars:   payoutResult.Populars(),
			})
		}

		newRaces = append(newRaces, &raw_entity.Race{
			RaceId:         raceId,
			RaceDate:       ticket.RaceDate().Value(),
			RaceNumber:     ticket.RaceNo(),
			RaceCourseId:   ticket.RaceCourse().Value(),
			RaceName:       race.RaceName(),
			Url:            race.Url(),
			Time:           race.Time(),
			StartTime:      race.StartTime(),
			Entries:        race.Entries(),
			Distance:       race.Distance(),
			Class:          race.Class(),
			CourseCategory: race.CourseCategory(),
			TrackCondition: race.TrackCondition(),
			RaceResults:    raceResults,
			PayoutResults:  payoutResults,
		})
	}
	rawRaces = append(rawRaces, newRaces...)

	raceInfo := raw_entity.RaceInfo{
		Races: rawRaces,
	}

	err = d.raceDataRepository.Write(ctx, raceResultFileName, &raceInfo)
	if err != nil {
		return err
	}

	urls, err = d.netKeibaService.CreateJockeyUrls(ctx, rawJockeys, excludeJockeyIds)
	if err != nil {
		return err
	}
	newJockeys := make([]*raw_entity.Jockey, 0, len(urls))
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
			newJockeys = append(newJockeys, &raw_entity.Jockey{
				JockeyId:   jockey.Id(),
				JockeyName: jockey.Name(),
			})
		}
	}
	rawJockeys = append(rawJockeys, newJockeys...)
	excludeJockeyIds = append(excludeJockeyIds, newExcludeJockeyIds...)

	jockeyInfo := raw_entity.JockeyInfo{
		Jockeys:          rawJockeys,
		ExcludeJockeyIds: excludeJockeyIds,
	}

	err = d.jockeyDataRepository.Write(ctx, jockeyFileName, &jockeyInfo)
	if err != nil {
		return err
	}

	return nil
}
