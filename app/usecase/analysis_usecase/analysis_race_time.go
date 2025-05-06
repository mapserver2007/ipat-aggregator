package analysis_usecase

import (
	"context"
)

func (a *analysis) RaceTime(ctx context.Context, input *AnalysisInput) error {
	calculables, err := a.raceTimeService.Create(ctx, input.Races, input.RaceTimes)
	if err != nil {
		return err
	}
	analysisRaceTimeMap, attributeFilters, conditionFilters := a.raceTimeService.Convert(ctx, calculables)
	err = a.raceTimeService.Write(ctx, analysisRaceTimeMap, attributeFilters, conditionFilters)
	if err != nil {
		return err
	}

	return nil
}
