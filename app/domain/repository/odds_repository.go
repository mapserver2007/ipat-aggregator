package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
)

type OddsRepository interface {
	List(ctx context.Context, path string) ([]string, error)
	Read(ctx context.Context, path string) ([]*raw_entity.RaceOdds, error)
	Write(ctx context.Context, path string, data *raw_entity.RaceOddsInfo) error
	Fetch(ctx context.Context, url string) ([]*netkeiba_entity.Odds, error)
}
