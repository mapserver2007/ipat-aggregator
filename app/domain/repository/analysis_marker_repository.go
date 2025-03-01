package repository

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
)

type AnalysisMarkerRepository interface {
	Read(ctx context.Context, path string) ([]*marker_csv_entity.AnalysisMarker, error)
}
