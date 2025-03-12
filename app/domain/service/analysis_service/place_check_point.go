package analysis_service

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
)

type PlaceCheckPoint interface {
	GetPositivePoint(ctx context.Context, race *analysis_entity.Race) error
	GetNegativePoint(ctx context.Context, race *analysis_entity.Race) error
}

type placeCheckPointService struct {
	placeNegativeCheckService PlaceNegativeCheck
}

func NewPlaceCheckPoint(
	placeNegativeCheckService PlaceNegativeCheck,
) PlaceCheckPoint {
	return &placeCheckPointService{
		placeNegativeCheckService: placeNegativeCheckService,
	}
}

func (s *placeCheckPointService) GetPositivePoint(
	ctx context.Context,
	race *analysis_entity.Race,
) error {

	// 実装を追加
	return nil
}

func (s *placeCheckPointService) GetNegativePoint(
	ctx context.Context,
	race *analysis_entity.Race,
) error {
	input := &PlaceNegativeCheckInput{
		Race: race,
	}

	s.placeNegativeCheckService.DirtShortDistanceRace(ctx, input)

	// 実装を追加
	return nil
}
