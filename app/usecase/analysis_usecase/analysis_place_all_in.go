package analysis_usecase

import "context"

func (a *analysis) PlaceAllIn(ctx context.Context, input *AnalysisInput) error {
	placeAllInCalculables, err := a.placeAllInService.Create(ctx, input.Markers, input.Races, input.Odds.Win, input.Odds.Place)
	if err != nil {
		return err
	}
	placeAllInMap1, placeAllInMap2, attributeFilters, markerCombinationFilters := a.placeAllInService.Convert(ctx, placeAllInCalculables)
	err = a.placeAllInService.Write(ctx, placeAllInMap1, placeAllInMap2, attributeFilters, markerCombinationFilters)
	if err != nil {
		return err
	}

	return nil
}
