package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/predict_csv_entity"
)

type PredictDataRepository interface {
	Read(ctx context.Context, filePath string) ([]*predict_csv_entity.Yamato, error)
}
