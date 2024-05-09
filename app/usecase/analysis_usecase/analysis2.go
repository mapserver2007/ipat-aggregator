package analysis_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/analysis_service"
)

type Analysis2 interface {
	WinAndPlace(ctx context.Context, input *WinAndPlaceInput) error
	Trio(ctx context.Context, input *TrioInput) error
}

type WinAndPlaceInput struct {
	Markers []*marker_csv_entity.AnalysisMarker
	Races   []*data_cache_entity.Race
}

type TrioInput struct {
	Markers []*marker_csv_entity.AnalysisMarker
	Races   []*data_cache_entity.Race
	Odds    []*data_cache_entity.Odds
}

type analysis struct {
	trioService analysis_service.Trio
}

func NewAnalysis2(
	trioService analysis_service.Trio,
) Analysis2 {
	return &analysis{
		trioService: trioService,
	}
}

func (a *analysis) WinAndPlace(ctx context.Context, input *WinAndPlaceInput) error {
	//TODO implement me
	panic("implement me")
}

func (a *analysis) Trio(ctx context.Context, input *TrioInput) error {
	//TODO implement me
	panic("implement me")
}
