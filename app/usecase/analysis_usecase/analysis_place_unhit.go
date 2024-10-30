package analysis_usecase

import "context"

func (a *analysis) PlaceUnHit(ctx context.Context, input *AnalysisInput) error {
	a.placeUnHitService.Create(ctx, input.Markers, input.Races)

	return nil
}
