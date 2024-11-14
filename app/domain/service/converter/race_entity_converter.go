package converter

import (
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/list_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/tospo_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"time"
)

type RaceEntityConverter interface {
	DataCacheToRaw(input *data_cache_entity.Race) *raw_entity.Race
	NetKeibaToRaw(input *netkeiba_entity.Race) *raw_entity.Race
	RawToDataCache(input *raw_entity.Race) *data_cache_entity.Race
	DataCacheToList(input *data_cache_entity.Race) *list_entity.Race
	NetKeibaToPrediction(input1 *netkeiba_entity.Race, input2 []*netkeiba_entity.Odds, filters []filter.Id) *prediction_entity.Race
	TospoToPrediction(input1 *tospo_entity.Forecast, input2 *tospo_entity.TrainingComment, input3 []*tospo_entity.Memo, input4 *tospo_entity.PaddockComment) *prediction_entity.RaceForecast
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
			HorseId:       raceResult.HorseId().Value(),
			HorseName:     raceResult.HorseName(),
			BracketNumber: raceResult.BracketNumber(),
			HorseNumber:   raceResult.HorseNumber().Value(),
			JockeyId:      raceResult.JockeyId().Value(),
			Odds:          raceResult.Odds().StringFixed(1),
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
		Organizer:           input.Organizer().Value(),
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
			HorseId:       raceResult.HorseId(),
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
			raceResult.HorseId,
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
		input.Organizer,
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

func (r *raceEntityConverter) NetKeibaToPrediction(
	input1 *netkeiba_entity.Race,
	input2 []*netkeiba_entity.Odds,
	filters []filter.Id,
) *prediction_entity.Race {
	raceEntryHorses := make([]*prediction_entity.RaceEntryHorse, 0, len(input1.RaceEntryHorses()))
	for _, rawRaceEntryHorse := range input1.RaceEntryHorses() {
		raceEntryHorses = append(raceEntryHorses, prediction_entity.NewRaceEntryHorse(
			rawRaceEntryHorse.HorseId(),
			rawRaceEntryHorse.HorseName(),
			rawRaceEntryHorse.BracketNumber(),
			rawRaceEntryHorse.HorseNumber(),
			rawRaceEntryHorse.JockeyId(),
			rawRaceEntryHorse.TrainerId(),
			rawRaceEntryHorse.RaceWeight(),
		))
	}

	var predictionOdds []*prediction_entity.Odds
	for _, nkOdds := range input2 {
		predictionOdds = append(predictionOdds, prediction_entity.NewOdds(
			nkOdds.Odds()[0],
			nkOdds.PopularNumber(),
			nkOdds.HorseNumbers()[0],
		))
	}

	// レース結果はここでは取らないが、予想処理のために空で設定しておく
	raceResultHorseNumbers := make([]int, 3)

	return prediction_entity.NewRace(
		input1.RaceId(),
		input1.RaceName(),
		input1.RaceDate(),
		input1.RaceNumber(),
		input1.Entries(),
		input1.Distance(),
		input1.Class(),
		input1.CourseCategory(),
		input1.TrackCondition(),
		input1.RaceSexCondition(),
		input1.RaceWeightCondition(),
		input1.RaceCourseId(),
		input1.Url(),
		raceEntryHorses,
		raceResultHorseNumbers,
		predictionOdds,
		filters,
	)
}

func (r *raceEntityConverter) TospoToPrediction(
	input1 *tospo_entity.Forecast,
	input2 *tospo_entity.TrainingComment,
	input3 []*tospo_entity.Memo,
	input4 *tospo_entity.PaddockComment,
) *prediction_entity.RaceForecast {
	reporterMemos := make([]string, 0, len(input3))
	if len(input3) > 0 {
		for _, memo := range input3 {
			// 2週間以内のコメントだけ使う
			twoWeeksAgo := time.Now().AddDate(0, 0, -14)
			if memo.Date().After(twoWeeksAgo) || memo.Date().Equal(twoWeeksAgo) {
				reporterMemos = append(reporterMemos, fmt.Sprintf("%s %s", memo.Date().Format("2006/01/02"), memo.Comment()))
			}
		}
	}

	var (
		paddockComment    string
		paddockEvaluation int
	)
	if input4 != nil {
		paddockComment = input4.Comment()
		paddockEvaluation = input4.Evaluation()
	}

	return prediction_entity.NewRaceForecast(
		input1.HorseNumber(),
		input1.FavoriteNum(),
		input1.RivalNum(),
		input1.MarkerNum(),
		input2.TrainingComment(),
		input2.IsHighlyRecommended(),
		reporterMemos,
		paddockComment,
		paddockEvaluation,
	)
}
