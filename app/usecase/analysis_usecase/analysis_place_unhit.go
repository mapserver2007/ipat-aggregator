package analysis_usecase

import "context"

func (a *analysis) PlaceUnHit(ctx context.Context, input *AnalysisInput) error {
	unHitRaces := a.placeUnHitService.GetRaces(ctx, input.Markers, input.Races)

	_ = unHitRaces

	return nil
}
