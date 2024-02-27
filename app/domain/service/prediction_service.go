package service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
)

type PredictionService interface {
	CreateFilters(ctx context.Context, race *netkeiba_entity.Race) error
}

type predictionService struct {
}

func NewPredictionService() PredictionService {
	return &predictionService{}
}

func (p predictionService) CreateFilters(
	ctx context.Context,
	race *netkeiba_entity.Race,
) error {
	//TODO implement me
	return nil
}
