package controller

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/prediction_usecase"
)

type Prediction struct {
	predictionUseCase prediction_usecase.Prediction
}

type PredictionInput struct {
	Master *MasterOutput
}

func NewPrediction(
	predictionUseCase prediction_usecase.Prediction,
) *Prediction {
	return &Prediction{
		predictionUseCase: predictionUseCase,
	}
}

func (p *Prediction) Execute(ctx context.Context, input *PredictionInput) error {
	return p.predictionUseCase.Execute(ctx, &prediction_usecase.PredictionInput{
		AnalysisMarkers:   input.Master.AnalysisMarkers,
		PredictionMarkers: input.Master.PredictionMarkers,
		Races:             input.Master.Races,
	})
}
