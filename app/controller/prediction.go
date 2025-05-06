package controller

import (
	"context"
	"sync"

	"github.com/mapserver2007/ipat-aggregator/app/usecase/prediction_usecase"
	"github.com/sirupsen/logrus"
)

type Prediction struct {
	predictionUseCase prediction_usecase.Prediction
	logger            *logrus.Logger
}

type PredictionInput struct {
	Master *MasterOutput
}

func NewPrediction(
	predictionUseCase prediction_usecase.Prediction,
	logger *logrus.Logger,
) *Prediction {
	return &Prediction{
		predictionUseCase: predictionUseCase,
		logger:            logger,
	}
}

func (p *Prediction) Prediction(ctx context.Context, input *PredictionInput) {
	var wg sync.WaitGroup
	const predictionParallel = 2
	errors := make(chan error, predictionParallel)

	for i := range make([]struct{}, predictionParallel) {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			switch i {
			case 0:
				p.logger.Info("fetching prediction odds start")
				if err := p.predictionUseCase.Odds(ctx, &prediction_usecase.PredictionInput{
					AnalysisMarkers:   input.Master.AnalysisMarkers,
					PredictionMarkers: input.Master.PredictionMarkers,
					Races:             input.Master.Races,
					RaceTimes:         input.Master.RaceTimes,
				}); err != nil {
					errors <- err
				}
				p.logger.Info("fetching prediction odds end")
			case 1:
				p.logger.Info("fetching prediction checklist start")
				if err := p.predictionUseCase.CheckList(ctx, &prediction_usecase.PredictionInput{
					AnalysisMarkers:   input.Master.AnalysisMarkers,
					PredictionMarkers: input.Master.PredictionMarkers,
					Races:             input.Master.Races,
				}); err != nil {
					errors <- err
				}
				p.logger.Info("fetching prediction checklist end")
			}
		}(i)
	}

	wg.Wait()
	close(errors)
	for err := range errors {
		p.logger.Errorf("prediction error: %v", err)
	}
}

func (p *Prediction) SyncMarker(ctx context.Context) {
	p.logger.Info("fetching prediction marker sync start")
	if err := p.predictionUseCase.Sync(ctx); err != nil {
		p.logger.Errorf("sync marker error: %v", err)
	}
	p.logger.Info("fetching prediction marker sync end")
}
