package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
)

type PredictionDataRepository interface {
	Read(ctx context.Context, filePath string) ([]*marker_csv_entity.PredictionMarker, error)
	Fetch(ctx context.Context, raceUrl, oddsUrl string) (*netkeiba_entity.Race, []*netkeiba_entity.Odds, error)
}
