package converter

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
)

type HorseEntityConverter interface {
	NetKeibaToPrediction(input *netkeiba_entity.Horse) (*prediction_entity.Horse, error)
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

	horse, err := prediction_entity.NewHorse(
		input.HorseId(),
		input.HorseName(),
		input.HorseBirthDay(),
		input.TrainerId(),
		input.OwnerId(),
		input.BreederId(),
		horseBlood,
		horseResults,
	)

	if err != nil {
		return nil, err
	}

	return horse, nil
}
