package prediction_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/analysis_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/prediction_service"
)

type Prediction interface {
	Execute(ctx context.Context, input *PredictionInput) error
}

type PredictionInput struct {
	AnalysisMarkers   []*marker_csv_entity.AnalysisMarker
	PredictionMarkers []*marker_csv_entity.PredictionMarker
	Races             []*data_cache_entity.Race
}

type prediction struct {
	predictionOddsService prediction_service.Odds
	placeService          analysis_service.Place
}

func NewPrediction(
	predictionOddsService prediction_service.Odds,
	placeService analysis_service.Place,
) Prediction {
	return &prediction{
		predictionOddsService: predictionOddsService,
		placeService:          placeService,
	}
}

func (p *prediction) Execute(ctx context.Context, input *PredictionInput) error {
	predictionRaces := make([]*prediction_entity.Race, 0, len(input.PredictionMarkers))
	for _, marker := range input.PredictionMarkers {
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

	err = p.predictionOddsService.Convert(ctx, predictionRaces, calculables)
	if err != nil {
		return err
	}

	return nil
}
