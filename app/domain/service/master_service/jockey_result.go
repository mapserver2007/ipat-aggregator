package master_service

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/sirupsen/logrus"
)

type JockeyResult interface {
	CreateOrUpdate(ctx context.Context, jockeyResults []*data_cache_entity.JockeyResult) error
}

type jockeyResultService struct {
	jockeyResultRepository repository.JockeyResultRepository
	logger                 *logrus.Logger
}

func NewJockeyResult(
	jockeyResultRepository repository.JockeyResultRepository,
	logger *logrus.Logger,
) JockeyResult {
	return &jockeyResultService{
		jockeyResultRepository: jockeyResultRepository,
		logger:                 logger,
	}
}

func (j *jockeyResultService) CreateOrUpdate(
	ctx context.Context,
	jockeyResults []*data_cache_entity.JockeyResult,
) error {
	return nil
}
