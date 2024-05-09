package controller

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/analysis_usecase"
)

type Analysis struct {
	analysisUseCase analysis_usecase.Analysis2
}

func NewAnalysis(analysisUseCase analysis_usecase.Analysis2) *Analysis {
	return &Analysis{
		analysisUseCase: analysisUseCase,
	}
}

func (a *Analysis) Execute(ctx context.Context) error {
	// TODO
	return nil
}
