package prediction_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
)

func (p *prediction) Odds(ctx context.Context, input *PredictionInput) error {
	predictionMarkers := input.PredictionMarkers
	predictionRaces := make([]*prediction_entity.Race, 0, len(predictionMarkers))

	for _, marker := range predictionMarkers {
		predictionRace, err := p.predictionOddsService.Get(ctx, marker.RaceId())
		if err != nil {
			return err
		}
		predictionRaces = append(predictionRaces, predictionRace)
	}

	calculables, err := p.placeService.Create(ctx, input.AnalysisMarkers, input.Races)
	if err != nil {
		return err
	}

	firstPlaceMap, secondPlaceMap, thirdPlaceMap, raceCourseMap := p.predictionOddsService.ConvertAll(ctx, predictionRaces, predictionMarkers, calculables)
	err = p.predictionOddsService.Write(ctx, firstPlaceMap, secondPlaceMap, thirdPlaceMap, raceCourseMap)
	if err != nil {
		return err
	}

	return nil
}
