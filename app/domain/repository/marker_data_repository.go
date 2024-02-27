package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
)

type MarkerDataRepository interface {
	Read(ctx context.Context, filePath string) ([]*marker_csv_entity.AnalysisMarker, error)
}
