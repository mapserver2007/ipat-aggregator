package controller

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/prediction_usecase"
	"github.com/sirupsen/logrus"
	"sync"
)

const parallel = 3

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

func (p *Prediction) Execute(ctx context.Context, input *PredictionInput) error {
	var wg sync.WaitGroup
	errors := make(chan error, parallel)

	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			switch i {
			case 0:
				p.logger.Info("fetching prediction odds start")
				if err := p.predictionUseCase.Odds(ctx, &prediction_usecase.PredictionInput{
					AnalysisMarkers:   input.Master.AnalysisMarkers,
					PredictionMarkers: input.Master.PredictionMarkers,
					Races:             input.Master.Races,
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
			case 2:
				p.logger.Info("fetching prediction marker sync start")
				if err := p.predictionUseCase.Sync(ctx); err != nil {
					errors <- err
				}
				p.logger.Info("fetching prediction marker sync end")
			}
		}()
	}

	wg.Wait()
	close(errors)
	for err := range errors {
		p.logger.Errorf("prediction error: %v", err)
	}

	return nil
}
