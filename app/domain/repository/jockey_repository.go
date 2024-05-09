package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
)

type JockeyRepository interface {
	Read(ctx context.Context, path string) (*raw_entity.JockeyInfo, error)
	Write(ctx context.Context, path string, data *raw_entity.JockeyInfo) error
	Fetch(ctx context.Context, url string) (*netkeiba_entity.Jockey, error)
}
