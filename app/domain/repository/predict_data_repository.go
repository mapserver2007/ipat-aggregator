package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
)

type PredictDataRepository interface {
	Read(ctx context.Context, filePath string) ([]*raw_entity.Predict, error)
}
