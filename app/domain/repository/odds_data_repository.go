package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
)

type OddsDataRepository interface {
	Read(ctx context.Context, filePath string) ([]*raw_entity.RaceOdds, error)
	Write(ctx context.Context, filePath string, oddsInfo *raw_entity.RaceOddsInfo) error
	Fetch(ctx context.Context, url string) ([]*netkeiba_entity.Odds, error)
}
