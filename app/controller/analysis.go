package controller

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/usecase/analysis_usecase"
	"github.com/sirupsen/logrus"
)

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

func (a *Analysis) Place(ctx context.Context, input *AnalysisInput) {
	a.logger.Info("fetching analysis place start")
	if err := a.analysisUseCase.Place(ctx, &analysis_usecase.AnalysisInput{
		Markers: input.Master.AnalysisMarkers,
		Races:   input.Master.Races,
		Odds: &analysis_usecase.AnalysisOddsInput{
			Win:   input.Master.WinOdds,
			Place: input.Master.PlaceOdds,
		},
	}); err != nil {
		a.logger.Errorf("analysis place error: %v", err)
	}
	a.logger.Info("fetching analysis place end")
}

func (a *Analysis) PlaceAllIn(ctx context.Context, input *AnalysisInput) {
	a.logger.Info("fetching analysis place all in start")
	if err := a.analysisUseCase.PlaceAllIn(ctx, &analysis_usecase.AnalysisInput{
		Markers: input.Master.AnalysisMarkers,
		Races:   input.Master.Races,
		Odds: &analysis_usecase.AnalysisOddsInput{
			Win:   input.Master.WinOdds,
			Place: input.Master.PlaceOdds,
		},
	}); err != nil {
		a.logger.Errorf("analysis place all in error: %v", err)
	}
	a.logger.Info("fetching analysis place all in end")
}

func (a *Analysis) PlaceUnHit(ctx context.Context, input *AnalysisInput) {
	a.logger.Info("fetching analysis place un hit in start")
	if err := a.analysisUseCase.PlaceUnHit(ctx, &analysis_usecase.AnalysisInput{
		Markers: input.Master.AnalysisMarkers,
		Races:   input.Master.Races,
		Jockeys: input.Master.Jockeys,
		Odds: &analysis_usecase.AnalysisOddsInput{
			Win:   input.Master.WinOdds,
			Place: input.Master.PlaceOdds,
			Trio:  input.Master.TrioOdds,
		},
	}); err != nil {
		a.logger.Errorf("analysis place un hit error: %v", err)
	}
	a.logger.Info("fetching analysis place un hit in end")
}

func (a *Analysis) Beta(ctx context.Context, input *AnalysisInput) {
	a.logger.Info("fetching analysis beta start")
	if err := a.analysisUseCase.Beta(ctx, &analysis_usecase.AnalysisInput{
		Markers: input.Master.AnalysisMarkers,
		Races:   input.Master.Races,
		Odds: &analysis_usecase.AnalysisOddsInput{
			Win:   input.Master.WinOdds,
			Place: input.Master.PlaceOdds,
		},
	}); err != nil {
		a.logger.Errorf("analysis beta error: %v", err)
	}
	a.logger.Info("fetching analysis beta end")
}
