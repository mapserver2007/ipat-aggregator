package controller

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/analysis_usecase"
)

type Analysis struct {
	analysisUseCase analysis_usecase.Analysis
}

type AnalysisInput struct {
	Master *MasterOutput
}

func NewAnalysis(
	analysisUseCase analysis_usecase.Analysis,
) *Analysis {
	return &Analysis{
		analysisUseCase: analysisUseCase,
	}
}

func (a *Analysis) Execute(ctx context.Context, input *AnalysisInput) error {
	return a.analysisUseCase.Execute(ctx, &analysis_usecase.AnalysisInput{
		Markers: input.Master.AnalysisMarkers,
		Races:   input.Master.Races,
		Odds: &analysis_usecase.AnalysisOddsInput{
			Win:   input.Master.WinOdds,
			Place: input.Master.PlaceOdds,
		},
	})
}
