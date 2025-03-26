package analysis_service

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
)

type PlaceJockey interface {
	Create(ctx context.Context, races []*data_cache_entity.Race) error
}

type placeJockeyService struct {
}

func NewPlaceJockey() PlaceJockey {
	return &placeJockeyService{}
}

func (s *placeJockeyService) Create(
	ctx context.Context,
	races []*data_cache_entity.Race,
) error {
	return nil
}
