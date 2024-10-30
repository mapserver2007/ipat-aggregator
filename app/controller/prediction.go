package controller

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/prediction_usecase"
	"github.com/mapserver2007/ipat-aggregator/config"
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
	if config.EnablePredictionOdds {
		if err := p.predictionUseCase.Odds(ctx, &prediction_usecase.PredictionInput{
			AnalysisMarkers:   input.Master.AnalysisMarkers,
			PredictionMarkers: input.Master.PredictionMarkers,
			Races:             input.Master.Races,
		}); err != nil {
			return err
		}
	}

	if config.EnablePredictionCheckList {
		if err := p.predictionUseCase.CheckList(ctx, &prediction_usecase.PredictionInput{
			AnalysisMarkers:   input.Master.AnalysisMarkers,
			PredictionMarkers: input.Master.PredictionMarkers,
			Races:             input.Master.Races,
		}); err != nil {
			return err
		}
	}

	if config.EnablePredictionSync {
		if err := p.predictionUseCase.Sync(ctx); err != nil {
			return err
		}
	}

	return nil
}
