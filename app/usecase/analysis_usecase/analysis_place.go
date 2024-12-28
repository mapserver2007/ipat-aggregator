package analysis_usecase

import "context"

func (a *analysis) Place(ctx context.Context, input *AnalysisInput) error {
	placeCalculables, err := a.placeService.Create(ctx, input.Markers, input.Races)
	if err != nil {
		return err
	}

	firstPlaceMap, secondPlaceMap, thirdPlaceMap, filters := a.placeService.Convert(ctx, placeCalculables)

	err = a.placeService.Write(ctx, firstPlaceMap, secondPlaceMap, thirdPlaceMap, filters)
	if err != nil {
		return err
	}

	return nil
}
