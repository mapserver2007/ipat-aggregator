package analysis_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/analysis_service"
	"github.com/mapserver2007/ipat-aggregator/config"
)

type Analysis interface {
	Execute(ctx context.Context, input *AnalysisInput) error
}

type AnalysisInput struct {
	Markers []*marker_csv_entity.AnalysisMarker
	Races   []*data_cache_entity.Race
	Odds    *AnalysisOddsInput
}

type AnalysisOddsInput struct {
	Win   []*data_cache_entity.Odds
	Place []*data_cache_entity.Odds
}

type analysis struct {
	placeService      analysis_service.Place
	trioService       analysis_service.Trio
	placeAllInService analysis_service.PlaceAllIn
	placeUnHitService analysis_service.PlaceUnHit
}

func NewAnalysis(
	placeService analysis_service.Place,
	trioService analysis_service.Trio,
	placeAllInService analysis_service.PlaceAllIn,
	placeUnHitService analysis_service.PlaceUnHit,
) Analysis {
	return &analysis{
		placeService:      placeService,
		trioService:       trioService,
		placeAllInService: placeAllInService,
		placeUnHitService: placeUnHitService,
	}
}

func (a *analysis) Execute(ctx context.Context, input *AnalysisInput) error {
	if config.EnableAnalysisPlace {
		placeCalculables, err := a.placeService.Create(ctx, input.Markers, input.Races)
		if err != nil {
			return err
		}

		firstPlaceMap, secondPlaceMap, thirdPlaceMap, filters := a.placeService.Convert(ctx, placeCalculables)

		err = a.placeService.Write(ctx, firstPlaceMap, secondPlaceMap, thirdPlaceMap, filters)
		if err != nil {
			return err
		}
	}

	if config.EnableAnalysisPlaceAllIn {
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
	}

	if config.EnableAnalysisPlaceUnHit {
		calculables, err := a.placeAllInService.Create(ctx, input.Markers, input.Races, input.Odds.Win, input.Odds.Place)
		if err != nil {
			return err
		}

		err = a.placeUnHitService.Convert(ctx, calculables, config.AnalysisPlaceUnHitWinUpperOdds, config.AnalysisPlaceUnHitWinLowerOdds)
		if err != nil {
			return err
		}
	}

	return nil
}
