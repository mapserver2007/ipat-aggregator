package converter

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"strconv"
)

type HorseEntityConverter interface {
	NetKeibaToPrediction(input *netkeiba_entity.Horse) (*prediction_entity.Horse, error)
	NetKeibaToAnalysis(input *netkeiba_entity.Horse, raceDate types.RaceDate) (*analysis_entity.Horse, error)
	NetKeibaToDataCache(input *netkeiba_entity.Horse, raceDate types.RaceDate) (*data_cache_entity.Horse, error)
	RawToDataCache(input *raw_entity.Horse) (*data_cache_entity.Horse, error)
	DataCacheToAnalysis(input *data_cache_entity.Horse) (*analysis_entity.Horse, error)
	DataCacheToRaw(input *data_cache_entity.Horse) *raw_entity.Horse
	PredictionToAnalysis(input *prediction_entity.Horse) (*analysis_entity.Horse, error)
}

type horseEntityConverter struct{}

func NewHorseEntityConverter() HorseEntityConverter {
	return &horseEntityConverter{}
}

func (h *horseEntityConverter) NetKeibaToPrediction(input *netkeiba_entity.Horse) (*prediction_entity.Horse, error) {
	horseResults := make([]*prediction_entity.HorseResult, 0, len(input.HorseResults()))
	for _, rawHorseResult := range input.HorseResults() {
		horseResult, err := prediction_entity.NewHorseResult(
			rawHorseResult.RaceId(),
			rawHorseResult.RaceDate(),
			rawHorseResult.RaceName(),
			rawHorseResult.JockeyId(),
			rawHorseResult.OrderNo(),
			rawHorseResult.PopularNumber(),
			rawHorseResult.HorseNumber(),
			rawHorseResult.Odds(),
			rawHorseResult.Class(),
			rawHorseResult.Entries(),
			rawHorseResult.Distance(),
			rawHorseResult.RaceCourseId(),
			rawHorseResult.CourseCategoryId(),
			rawHorseResult.TrackConditionId(),
			rawHorseResult.HorseWeight(),
			rawHorseResult.RaceWeight(),
			rawHorseResult.Comment(),
		)
		if err != nil {
			return nil, err
		}
		horseResults = append(horseResults, horseResult)
	}

	horseBlood := prediction_entity.NewHorseBlood(input.HorseBlood().SireId(), input.HorseBlood().BroodmareSireId())

	return prediction_entity.NewHorse(
		input.HorseId(),
		input.HorseName(),
		input.HorseBirthDay(),
		input.TrainerId(),
		input.OwnerId(),
		input.BreederId(),
		horseBlood,
		horseResults,
	), nil
}

func (h *horseEntityConverter) NetKeibaToAnalysis(input *netkeiba_entity.Horse, raceDate types.RaceDate) (*analysis_entity.Horse, error) {
	horseResults := make([]*analysis_entity.HorseResult, 0, len(input.HorseResults()))
	for _, rawHorseResult := range input.HorseResults() {
		horseResult, err := analysis_entity.NewHorseResult(
			rawHorseResult.RaceId(),
			rawHorseResult.RaceDate(),
			rawHorseResult.RaceName(),
			rawHorseResult.JockeyId(),
			rawHorseResult.OrderNo(),
			rawHorseResult.PopularNumber(),
			rawHorseResult.HorseNumber(),
			rawHorseResult.Odds(),
			rawHorseResult.Class(),
			rawHorseResult.Entries(),
			rawHorseResult.Distance(),
			rawHorseResult.RaceCourseId(),
			rawHorseResult.CourseCategoryId(),
			rawHorseResult.TrackConditionId(),
			rawHorseResult.HorseWeight(),
			rawHorseResult.RaceWeight(),
			rawHorseResult.Comment(),
		)
		if err != nil {
			return nil, err
		}
		horseResults = append(horseResults, horseResult)
	}

	horseBlood := analysis_entity.NewHorseBlood(input.HorseBlood().SireId(), input.HorseBlood().BroodmareSireId())

	return analysis_entity.NewHorse(
		input.HorseId(),
		input.HorseName(),
		input.HorseBirthDay(),
		input.TrainerId(),
		input.OwnerId(),
		input.BreederId(),
		horseBlood,
		horseResults,
		raceDate,
	), nil
}

func (h *horseEntityConverter) NetKeibaToDataCache(input *netkeiba_entity.Horse, raceDate types.RaceDate) (*data_cache_entity.Horse, error) {
	horseResults := make([]*data_cache_entity.HorseResult, 0, len(input.HorseResults()))
	for _, rawHorseResult := range input.HorseResults() {
		horseResult, err := data_cache_entity.NewHorseResult(
			rawHorseResult.RaceId(),
			rawHorseResult.RaceDate(),
			rawHorseResult.RaceName(),
			rawHorseResult.JockeyId(),
			rawHorseResult.OrderNo(),
			rawHorseResult.PopularNumber(),
			rawHorseResult.HorseNumber(),
			rawHorseResult.Odds(),
			rawHorseResult.Class(),
			rawHorseResult.Entries(),
			rawHorseResult.Distance(),
			rawHorseResult.RaceCourseId(),
			rawHorseResult.CourseCategoryId(),
			rawHorseResult.TrackConditionId(),
			rawHorseResult.HorseWeight(),
			strconv.FormatFloat(rawHorseResult.RaceWeight(), 'f', 1, 64),
			rawHorseResult.Comment(),
		)
		if err != nil {
			return nil, err
		}
		horseResults = append(horseResults, horseResult)
	}

	horseBlood := data_cache_entity.NewHorseBlood(input.HorseBlood().SireId(), input.HorseBlood().BroodmareSireId())

	return data_cache_entity.NewHorse(
		input.HorseId(),
		input.HorseName(),
		input.HorseBirthDay(),
		input.TrainerId(),
		input.OwnerId(),
		input.BreederId(),
		horseBlood,
		horseResults,
		raceDate.Value(),
	), nil
}

func (h *horseEntityConverter) RawToDataCache(input *raw_entity.Horse) (*data_cache_entity.Horse, error) {
	horseBlood := data_cache_entity.NewHorseBlood(
		input.HorseBlood.SireId,
		input.HorseBlood.BroodmareSireId,
	)

	horseResults := make([]*data_cache_entity.HorseResult, 0, len(input.HorseResults))
	for _, rawHorseResult := range input.HorseResults {
		horseResult, err := data_cache_entity.NewHorseResult(
			rawHorseResult.RaceId,
			rawHorseResult.RaceDate,
			rawHorseResult.RaceName,
			rawHorseResult.JockeyId,
			rawHorseResult.OrderNo,
			rawHorseResult.PopularNumber,
			rawHorseResult.HorseNumber,
			rawHorseResult.Odds,
			rawHorseResult.Class,
			rawHorseResult.Entries,
			rawHorseResult.Distance,
			rawHorseResult.RaceCourseId,
			rawHorseResult.CourseCategoryId,
			rawHorseResult.TrackConditionId,
			rawHorseResult.HorseWeight,
			rawHorseResult.RaceWeight,
			rawHorseResult.Comment,
		)
		if err != nil {
			return nil, err
		}

		horseResults = append(horseResults, horseResult)
	}

	return data_cache_entity.NewHorse(
		input.HorseId,
		input.HorseName,
		input.HorseBirthDay,
		input.TrainerId,
		input.OwnerId,
		input.BreederId,
		horseBlood,
		horseResults,
		input.LatestRaceDate,
	), nil
}

func (h *horseEntityConverter) DataCacheToAnalysis(input *data_cache_entity.Horse) (*analysis_entity.Horse, error) {
	horseBlood := analysis_entity.NewHorseBlood(
		input.HorseBlood().SireId().Value(),
		input.HorseBlood().BroodmareSireId().Value(),
	)

	horseResults := make([]*analysis_entity.HorseResult, 0, len(input.HorseResults()))
	for _, rawHorseResult := range input.HorseResults() {
		horseResult, err := analysis_entity.NewHorseResult(
			rawHorseResult.RaceId().String(),
			rawHorseResult.RaceDate().Value(),
			rawHorseResult.RaceName(),
			rawHorseResult.JockeyId().Value(),
			rawHorseResult.OrderNo(),
			rawHorseResult.PopularNumber(),
			rawHorseResult.HorseNumber().Value(),
			rawHorseResult.Odds().String(),
			rawHorseResult.Class().Value(),
			rawHorseResult.Entries(),
			rawHorseResult.Distance(),
			rawHorseResult.RaceCourse().Value(),
			rawHorseResult.CourseCategory().Value(),
			rawHorseResult.TrackCondition().Value(),
			rawHorseResult.HorseWeight(),
			rawHorseResult.RaceWeight(),
			rawHorseResult.Comment(),
		)
		if err != nil {
			return nil, err
		}

		horseResults = append(horseResults, horseResult)
	}

	return analysis_entity.NewHorse(
		input.HorseId().Value(),
		input.HorseName(),
		input.HorseBirthDay().Value(),
		input.TrainerId().Value(),
		input.OwnerId().Value(),
		input.BreederId().Value(),
		horseBlood,
		horseResults,
		input.LatestRaceDate(),
	), nil
}

func (h *horseEntityConverter) DataCacheToRaw(input *data_cache_entity.Horse) *raw_entity.Horse {
	horseBlood := &raw_entity.HorseBlood{
		SireId:          input.HorseBlood().SireId().Value(),
		BroodmareSireId: input.HorseBlood().BroodmareSireId().Value(),
	}

	rawHorseResults := make([]*raw_entity.HorseResult, 0, len(input.HorseResults()))
	for _, horseResult := range input.HorseResults() {
		rawHorseResult := &raw_entity.HorseResult{
			RaceId:           horseResult.RaceId().String(),
			RaceDate:         horseResult.RaceDate().Value(),
			RaceName:         horseResult.RaceName(),
			JockeyId:         horseResult.JockeyId().Value(),
			OrderNo:          horseResult.OrderNo(),
			PopularNumber:    horseResult.PopularNumber(),
			HorseNumber:      horseResult.HorseNumber().Value(),
			Odds:             horseResult.Odds().String(),
			Class:            horseResult.Class().Value(),
			Entries:          horseResult.Entries(),
			Distance:         horseResult.Distance(),
			RaceCourseId:     horseResult.RaceCourse().Value(),
			CourseCategoryId: horseResult.CourseCategory().Value(),
			TrackConditionId: horseResult.TrackCondition().Value(),
			HorseWeight:      horseResult.HorseWeight(),
			RaceWeight:       strconv.FormatFloat(horseResult.RaceWeight(), 'f', 1, 64),
			Comment:          horseResult.Comment(),
		}
		rawHorseResults = append(rawHorseResults, rawHorseResult)
	}

	return &raw_entity.Horse{
		HorseId:        input.HorseId().Value(),
		HorseName:      input.HorseName(),
		HorseBirthDay:  input.HorseBirthDay().Value(),
		TrainerId:      input.TrainerId().Value(),
		OwnerId:        input.OwnerId().Value(),
		BreederId:      input.BreederId().Value(),
		HorseBlood:     horseBlood,
		HorseResults:   rawHorseResults,
		LatestRaceDate: input.LatestRaceDate().Value(),
	}
}

func (h *horseEntityConverter) PredictionToAnalysis(
	input *prediction_entity.Horse,
) (*analysis_entity.Horse, error) {
	horseBlood := analysis_entity.NewHorseBlood(
		input.HorseBlood().SireId().Value(),
		input.HorseBlood().BroodmareSireId().Value(),
	)

	horseResults := make([]*analysis_entity.HorseResult, 0, len(input.HorseResults()))
	for _, rawHorseResult := range input.HorseResults() {
		horseResult, err := analysis_entity.NewHorseResult(
			rawHorseResult.RaceId().String(),
			rawHorseResult.RaceDate().Value(),
			rawHorseResult.RaceName(),
			rawHorseResult.JockeyId().Value(),
			rawHorseResult.OrderNo(),
			rawHorseResult.PopularNumber(),
			rawHorseResult.HorseNumber().Value(),
			rawHorseResult.Odds().String(),
			rawHorseResult.Class().Value(),
			rawHorseResult.Entries(),
			rawHorseResult.Distance(),
			rawHorseResult.RaceCourse().Value(),
			rawHorseResult.CourseCategory().Value(),
			rawHorseResult.TrackCondition().Value(),
			rawHorseResult.HorseWeight(),
			rawHorseResult.RaceWeight(),
			rawHorseResult.Comment(),
		)
		if err != nil {
			return nil, err
		}
		horseResults = append(horseResults, horseResult)
	}

	return analysis_entity.NewHorse(
		input.HorseId().Value(),
		input.HorseName(),
		input.HorseBirthDay().Value(),
		input.TrainerId().Value(),
		input.OwnerId().Value(),
		input.BreederId().Value(),
		horseBlood,
		horseResults,
		types.RaceDate(0),
	), nil
}
