package controller

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/analysis_usecase"
	"github.com/sirupsen/logrus"
	"sync"
)

const analysisParallel = 2

type Analysis struct {
	analysisUseCase analysis_usecase.Analysis
	logger          *logrus.Logger
}

type AnalysisInput struct {
	Master *MasterOutput
}

func NewAnalysis(
	analysisUseCase analysis_usecase.Analysis,
	logger *logrus.Logger,
) *Analysis {
	return &Analysis{
		analysisUseCase: analysisUseCase,
		logger:          logger,
	}
}

func (a *Analysis) Execute(ctx context.Context, input *AnalysisInput) error {
	var wg sync.WaitGroup
	errors := make(chan error, analysisParallel)

	for i := 0; i < analysisParallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			switch i {
			case 0:
				a.logger.Info("fetching analysis place start")
				if err := a.analysisUseCase.Place(ctx, &analysis_usecase.AnalysisInput{
					Markers: input.Master.AnalysisMarkers,
					Races:   input.Master.Races,
					Odds: &analysis_usecase.AnalysisOddsInput{
						Win:   input.Master.WinOdds,
						Place: input.Master.PlaceOdds,
					},
				}); err != nil {
					errors <- err
				}
				a.logger.Info("fetching analysis place end")
			case 1:
				a.logger.Info("fetching analysis place all in start")
				if err := a.analysisUseCase.PlaceAllIn(ctx, &analysis_usecase.AnalysisInput{
					Markers: input.Master.AnalysisMarkers,
					Races:   input.Master.Races,
					Odds: &analysis_usecase.AnalysisOddsInput{
						Win:   input.Master.WinOdds,
						Place: input.Master.PlaceOdds,
					},
				}); err != nil {
					errors <- err
				}
				a.logger.Info("fetching analysis place all in end")
			case 2:
				// TODO 開発止めてるので実行しない
				a.logger.Info("fetching analysis place un hit in start")
				if err := a.analysisUseCase.PlaceUnHit(ctx, &analysis_usecase.AnalysisInput{
					Markers: input.Master.AnalysisMarkers,
					Races:   input.Master.Races,
					Odds: &analysis_usecase.AnalysisOddsInput{
						Win:   input.Master.WinOdds,
						Place: input.Master.PlaceOdds,
					},
				}); err != nil {
					errors <- err
				}
				a.logger.Info("fetching analysis place un hit in end")
			}
		}()
	}

	wg.Wait()
	close(errors)
	for err := range errors {
		a.logger.Errorf("analysis error: %v", err)
	}

	return nil
}
