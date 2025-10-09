package analysis_service

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
)

type usmService struct {
}

type Usm interface {
	Create(ctx context.Context) ([]*analysis_entity.JockeyResult, error)
}

func NewUsm() Usm {
	return &usmService{}
}

func (u *usmService) Create(ctx context.Context) ([]*analysis_entity.JockeyResult, error) {
	return nil, nil
}
