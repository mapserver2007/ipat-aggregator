package service

import (
	"fmt"
	betting_ticket_entity "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"
	race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/entity"
	raw_race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/raw_entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
)

type RaceConverter struct{}

func NewRaceConverter() RaceConverter {
	return RaceConverter{}
}

func (r *RaceConverter) GetRaceId(
	record *betting_ticket_entity.CsvEntity,
	racingNumber *race_entity.RacingNumber,
) *race_vo.RaceId {
	var raceId race_vo.RaceId
	organizer := record.RaceCourse().Organizer()

	switch organizer {
	case race_vo.JRA:
		rawRaceId := fmt.Sprintf("%d%02d%02d%02d%02d", record.RaceDate().Year(), racingNumber.RaceCourseId(), racingNumber.Round(), racingNumber.Day(), record.RaceNo())
		raceId = race_vo.RaceId(rawRaceId)
	case race_vo.NAR:
		rawRaceId := fmt.Sprintf("%d%02d%02d%02d%02d", record.RaceDate().Year(), record.RaceCourse().Value(), record.RaceDate().Month(), record.RaceDate().Day(), record.RaceNo())
		raceId = race_vo.RaceId(rawRaceId)
	case race_vo.OverseaOrganizer:
		raceCourseIdForOversea := race_vo.ConvertToOverseaRaceCourseId(record.RaceCourse())
		rawRaceId := fmt.Sprintf("%d%s%02d%02d%02d", record.RaceDate().Year(), raceCourseIdForOversea, record.RaceDate().Month(), record.RaceDate().Day(), record.RaceNo())
		raceId = race_vo.RaceId(rawRaceId)
	}

	return &raceId
}

func (r *RaceConverter) ConvertToRaceMapByRacingNumberId(races []*race_entity.Race) map[race_vo.RacingNumberId]*race_entity.Race {
	return ConvertToMap(races, func(race *race_entity.Race) race_vo.RacingNumberId {
		return race_vo.NewRacingNumberId(race.RaceDate(), race.RaceCourseId())
	})
}

func (r *RaceConverter) ConvertToRaceMapByRaceId(races []*race_entity.Race) map[race_vo.RaceId]*race_entity.Race {
	return ConvertToMap(races, func(race *race_entity.Race) race_vo.RaceId {
		return race.RaceId()
	})
}

func (r *RaceConverter) ConvertToRawRaceMap(races []*raw_race_entity.Race) map[string]*raw_race_entity.Race {
	return ConvertToMap(races, func(race *raw_race_entity.Race) string {
		return race.RaceId
	})
}

func (r *RaceConverter) ConvertToRawRacingNumberMap(races []*raw_race_entity.RacingNumber) map[race_vo.RacingNumberId]*raw_race_entity.RacingNumber {
	return ConvertToMap(races, func(racingNumber *raw_race_entity.RacingNumber) race_vo.RacingNumberId {
		return race_vo.NewRacingNumberId(
			race_vo.RaceDate(racingNumber.Date),
			race_vo.RaceCourse(racingNumber.RaceCourseId),
		)
	})
}

func (r *RaceConverter) ConvertToRacingNumberMap(races []*race_entity.RacingNumber) map[race_vo.RacingNumberId]*race_entity.RacingNumber {
	return ConvertToMap(races, func(racingNumber *race_entity.RacingNumber) race_vo.RacingNumberId {
		return race_vo.NewRacingNumberId(
			racingNumber.Date(),
			racingNumber.RaceCourseId(),
		)
	})
}

func (r *RaceConverter) ConvertFromRawRacingNumberNetkeibaToRawRacingNumberCsv(rawRacingNumber *raw_race_entity.RawRacingNumberNetkeiba) *raw_race_entity.RacingNumber {
	return &raw_race_entity.RacingNumber{
		Date:         rawRacingNumber.Date(),
		Round:        rawRacingNumber.Round(),
		Day:          rawRacingNumber.Day(),
		RaceCourseId: rawRacingNumber.RaceCourseId(),
	}
}

func (r *RaceConverter) ConvertFromRawRaceNetkeibaToRawRaceCsv(rawRace *raw_race_entity.RawRaceNetkeiba, raceId *race_vo.RaceId, record *betting_ticket_entity.CsvEntity) *raw_race_entity.Race {
	return &raw_race_entity.Race{
		RaceId:         string(*raceId),
		RaceDate:       int(record.RaceDate()),
		RaceNumber:     record.RaceNo(),
		RaceCourseId:   record.RaceCourse().Value(),
		RaceName:       rawRace.RaceName(),
		Url:            rawRace.Url(),
		Time:           rawRace.Time(),
		Entries:        rawRace.Entries(),
		Distance:       rawRace.Distance(),
		Class:          rawRace.Class(),
		CourseCategory: rawRace.CourseCategory(),
		TrackCondition: rawRace.TrackCondition(),
		RaceResults:    r.ConvertFromRawRaceResultsNetkeibaToRawRaceResultsCsv(rawRace.RaceResults()),
		PayoutResults:  r.ConvertFromRawPayoutResultsNetkeibaToRawPayoutResultsCsv(rawRace.PayoutResults()),
	}
}

func (r *RaceConverter) ConvertFromRawRaceResultsNetkeibaToRawRaceResultsCsv(rawRaceResults []*raw_race_entity.RawRaceResultNetkeiba) []*raw_race_entity.RaceResult {
	var raceResults []*raw_race_entity.RaceResult
	for _, rawRaceResult := range rawRaceResults {
		raceResult := &raw_race_entity.RaceResult{
			OrderNo:       rawRaceResult.OrderNo(),
			HorseName:     rawRaceResult.HorseName(),
			BracketNumber: rawRaceResult.BracketNumber(),
			HorseNumber:   rawRaceResult.HorseNumber(),
			Odds:          rawRaceResult.Odds(),
			PopularNumber: rawRaceResult.PopularNumber(),
		}
		raceResults = append(raceResults, raceResult)
	}

	return raceResults
}

func (r *RaceConverter) ConvertFromRawPayoutResultsNetkeibaToRawPayoutResultsCsv(rawPayoutResults []*raw_race_entity.RawPayoutResultNetkeiba) []*raw_race_entity.PayoutResult {
	var payoutResults []*raw_race_entity.PayoutResult
	for _, rawPayoutResult := range rawPayoutResults {
		payoutResult := &raw_race_entity.PayoutResult{
			TicketType: rawPayoutResult.TicketType(),
			Numbers:    rawPayoutResult.Numbers(),
			Odds:       rawPayoutResult.Odds(),
		}
		payoutResults = append(payoutResults, payoutResult)
	}

	return payoutResults
}

func (r *RaceConverter) ConvertFromRawRacesCsvToRaces(rawRaces []*raw_race_entity.Race) []*race_entity.Race {
	var races []*race_entity.Race
	for _, rawRace := range rawRaces {
		race := race_entity.NewRace(
			rawRace.RaceId,
			rawRace.RaceDate,
			rawRace.RaceNumber,
			rawRace.RaceCourseId,
			rawRace.RaceName,
			rawRace.Url,
			rawRace.Time,
			rawRace.Entries,
			rawRace.Distance,
			rawRace.Class,
			rawRace.CourseCategory,
			rawRace.TrackCondition,
			r.ConvertFromRawRaceResultsCsvToRaceResults(rawRace.RaceResults),
			r.ConvertFromRawPayoutResultsCsvToPayoutResults(rawRace.PayoutResults),
		)
		races = append(races, race)
	}

	return races
}

func (r *RaceConverter) ConvertFromRawRacingNumbersCsvToRacingNumbers(rawRacingNumbers []*raw_race_entity.RacingNumber) []*race_entity.RacingNumber {
	var racingNumbers []*race_entity.RacingNumber
	for _, rawRacingNumber := range rawRacingNumbers {
		racingNumber := race_entity.NewRacingNumber(
			rawRacingNumber.Date,
			rawRacingNumber.Round,
			rawRacingNumber.Day,
			rawRacingNumber.RaceCourseId,
		)
		racingNumbers = append(racingNumbers, racingNumber)
	}

	return racingNumbers
}

func (r *RaceConverter) ConvertFromRawRaceResultsCsvToRaceResults(rawRaceResults []*raw_race_entity.RaceResult) []*race_entity.RaceResult {
	var raceResults []*race_entity.RaceResult
	for _, rawRaceResult := range rawRaceResults {
		raceResult := race_entity.NewRaceResult(
			rawRaceResult.OrderNo,
			rawRaceResult.HorseName,
			rawRaceResult.BracketNumber,
			rawRaceResult.HorseNumber,
			rawRaceResult.Odds,
			rawRaceResult.PopularNumber,
		)
		raceResults = append(raceResults, raceResult)
	}

	return raceResults
}

func (r *RaceConverter) ConvertFromRawPayoutResultsCsvToPayoutResults(rawPayoutResults []*raw_race_entity.PayoutResult) []*race_entity.PayoutResult {
	var payoutResults []*race_entity.PayoutResult
	for _, rawPayoutResult := range rawPayoutResults {
		payoutResult := race_entity.NewPayoutResult(
			rawPayoutResult.TicketType,
			rawPayoutResult.Numbers,
			rawPayoutResult.Odds,
		)
		payoutResults = append(payoutResults, payoutResult)
	}

	return payoutResults
}
