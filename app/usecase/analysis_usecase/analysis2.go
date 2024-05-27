package analysis_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/analysis_service"
)

type Analysis2 interface {
	Execute(ctx context.Context, input *AnalysisInput) error
}

type AnalysisInput struct {
	Markers []*marker_csv_entity.AnalysisMarker
	Races   []*data_cache_entity.Race
}

type analysis struct {
	placeService analysis_service.Place
	trioService  analysis_service.Trio
}

func NewAnalysis2(
	placeService analysis_service.Place,
	trioService analysis_service.Trio,
) Analysis2 {
	return &analysis{
		placeService: placeService,
		trioService:  trioService,
	}
}

func (a *analysis) Execute(ctx context.Context, input *AnalysisInput) error {
	// TODO 3連複は後で実装
	//a.trioService.Create(ctx, input.Markers, input.Races)

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
