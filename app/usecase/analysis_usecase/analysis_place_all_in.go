package analysis_usecase

import "context"

func (a *analysis) PlaceAllIn(ctx context.Context, input *AnalysisInput) error {
	placeAllInCalculables, err := a.placeAllInService.Create(ctx, input.Markers, input.Races, input.Odds.Win, input.Odds.Place)
	if err != nil {
		return err
	}
	placeAllInMap, filters := a.placeAllInService.Convert(ctx, placeAllInCalculables)
	if err != nil {
		return err
	}
	err = a.placeAllInService.Write(ctx, placeAllInMap, filters)
	if err != nil {
		return err
	}

	return nil
}
