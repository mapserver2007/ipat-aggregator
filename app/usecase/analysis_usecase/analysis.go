package analysis_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/analysis_service"
)

type Analysis interface {
	Place(ctx context.Context, input *AnalysisInput) error
	PlaceAllIn(ctx context.Context, input *AnalysisInput) error
	PlaceUnHit(ctx context.Context, input *AnalysisInput) error
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
