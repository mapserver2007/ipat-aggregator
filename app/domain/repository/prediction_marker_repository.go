package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
)

type PredictionMarkerRepository interface {
	Read(ctx context.Context, path string) ([]*marker_csv_entity.PredictionMarker, error)
	Push(ctx context.Context) error
}
