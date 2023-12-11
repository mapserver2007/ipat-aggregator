package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
)

type JockeyDataRepository interface {
	Read(ctx context.Context, fileName string) ([]*raw_entity.Jockey, []int, error)
	Write(ctx context.Context, fileName string, jockeyInfo *raw_entity.JockeyInfo) error
	Fetch(ctx context.Context, url string) (*netkeiba_entity.Jockey, error)
}
