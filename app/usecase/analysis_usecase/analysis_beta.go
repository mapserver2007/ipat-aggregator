package analysis_usecase

import "context"

func (a *analysis) Beta(ctx context.Context, input *AnalysisInput) error {
	winCalculables, err := a.betaWinService.Create(ctx, input.Markers, input.Races)
	if err != nil {
		return err
	}

	a.betaWinService.Convert(ctx, winCalculables)

	return nil
}
