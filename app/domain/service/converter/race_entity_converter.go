package converter

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/list_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
)

type RaceEntityConverter interface {
	DataCacheToRaw(input *data_cache_entity.Race) *raw_entity.Race
	NetKeibaToRaw(input *netkeiba_entity.Race) *raw_entity.Race
	RawToDataCache(input *raw_entity.Race) *data_cache_entity.Race
	DataCacheToList(input *data_cache_entity.Race) *list_entity.Race
}

type raceEntityConverter struct{}

func NewRaceEntityConverter() RaceEntityConverter {
	return &raceEntityConverter{}
}

func (r *raceEntityConverter) DataCacheToRaw(input *data_cache_entity.Race) *raw_entity.Race {
	raceResults := make([]*raw_entity.RaceResult, 0, len(input.RaceResults()))
	for _, raceResult := range input.RaceResults() {
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
	payoutResults := make([]*raw_entity.PayoutResult, 0, len(input.PayoutResults()))
	for _, payoutResult := range input.PayoutResults() {
		rawNumbers := make([]string, 0, len(payoutResult.Numbers()))
		for _, number := range payoutResult.Numbers() {
			rawNumbers = append(rawNumbers, number.String())
		}
		payoutResults = append(payoutResults, &raw_entity.PayoutResult{
			TicketType: payoutResult.TicketType().Value(),
			Numbers:    rawNumbers,
			Odds:       payoutResult.Odds(),
			Populars:   payoutResult.Populars(),
		})
	}

	return &raw_entity.Race{
		RaceId:              input.RaceId().String(),
		RaceDate:            input.RaceDate().Value(),
		RaceNumber:          input.RaceNumber(),
		RaceCourseId:        input.RaceCourseId().Value(),
		RaceName:            input.RaceName(),
		Url:                 input.Url(),
		Time:                input.Time(),
		StartTime:           input.StartTime(),
		Entries:             input.Entries(),
		Distance:            input.Distance(),
		Class:               input.Class().Value(),
		CourseCategory:      input.CourseCategory().Value(),
		TrackCondition:      input.TrackCondition().Value(),
		RaceSexCondition:    input.RaceSexCondition().Value(),
		RaceWeightCondition: input.RaceWeightCondition().Value(),
		RaceResults:         raceResults,
		PayoutResults:       payoutResults,
	}
}

func (r *raceEntityConverter) NetKeibaToRaw(input *netkeiba_entity.Race) *raw_entity.Race {
	raceResults := make([]*raw_entity.RaceResult, 0, len(input.RaceResults()))
	for _, raceResult := range input.RaceResults() {
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
	payoutResults := make([]*raw_entity.PayoutResult, 0, len(input.PayoutResults()))
	for _, payoutResult := range input.PayoutResults() {
		rawNumbers := make([]string, 0, len(payoutResult.Numbers()))
		for _, number := range payoutResult.Numbers() {
			rawNumbers = append(rawNumbers, number)
		}
		payoutResults = append(payoutResults, &raw_entity.PayoutResult{
			TicketType: payoutResult.TicketType(),
			Numbers:    rawNumbers,
			Odds:       payoutResult.Odds(),
			Populars:   payoutResult.Populars(),
		})
	}

	return &raw_entity.Race{
		RaceId:              input.RaceId(),
		RaceDate:            input.RaceDate(),
		RaceNumber:          input.RaceNumber(),
		RaceCourseId:        input.RaceCourseId(),
		RaceName:            input.RaceName(),
		Organizer:           input.Organizer(),
		Url:                 input.Url(),
		Time:                input.Time(),
		StartTime:           input.StartTime(),
		Entries:             input.Entries(),
		Distance:            input.Distance(),
		Class:               input.Class(),
		CourseCategory:      input.CourseCategory(),
		TrackCondition:      input.TrackCondition(),
		RaceSexCondition:    input.RaceSexCondition(),
		RaceWeightCondition: input.RaceWeightCondition(),
		RaceResults:         raceResults,
		PayoutResults:       payoutResults,
	}
}

func (r *raceEntityConverter) RawToDataCache(input *raw_entity.Race) *data_cache_entity.Race {
	raceResults := make([]*data_cache_entity.RaceResult, 0, len(input.RaceResults))
	for _, raceResult := range input.RaceResults {
		raceResults = append(raceResults, data_cache_entity.NewRaceResult(
			raceResult.OrderNo,
			raceResult.HorseName,
			raceResult.BracketNumber,
			raceResult.HorseNumber,
			raceResult.JockeyId,
			raceResult.Odds,
			raceResult.PopularNumber,
		))
	}
	payoutResults := make([]*data_cache_entity.PayoutResult, 0, len(input.PayoutResults))
	for _, payoutResult := range input.PayoutResults {
		payoutResults = append(payoutResults, data_cache_entity.NewPayoutResult(
			payoutResult.TicketType,
			payoutResult.Numbers,
			payoutResult.Odds,
			payoutResult.Populars,
		))
	}

	return data_cache_entity.NewRace(
		input.RaceId,
		input.RaceDate,
		input.RaceNumber,
		input.RaceCourseId,
		input.RaceName,
		input.Url,
		input.Time,
		input.StartTime,
		input.Entries,
		input.Distance,
		input.Class,
		input.CourseCategory,
		input.TrackCondition,
		input.RaceSexCondition,
		input.RaceWeightCondition,
		raceResults,
		payoutResults,
	)
}

func (r *raceEntityConverter) DataCacheToList(input *data_cache_entity.Race) *list_entity.Race {
	raceResults := make([]*list_entity.RaceResult, 0, len(input.RaceResults()))
	for _, raceResult := range input.RaceResults() {
		raceResults = append(raceResults, list_entity.NewRaceResult(
			raceResult.OrderNo(),
			raceResult.HorseName(),
			raceResult.BracketNumber(),
			raceResult.HorseNumber(),
			raceResult.JockeyId(),
			raceResult.Odds(),
			raceResult.PopularNumber(),
		))
	}

	return list_entity.NewRace(
		input.RaceId(),
		input.RaceNumber(),
		input.RaceName(),
		input.StartTime(),
		input.Class(),
		input.RaceCourseId(),
		input.CourseCategory(),
		input.RaceDate(),
		input.Distance(),
		input.TrackCondition(),
		input.Url(),
		raceResults,
	)
}
