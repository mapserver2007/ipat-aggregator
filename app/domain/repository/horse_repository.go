package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
)

type HorseRepository interface {
	Read(ctx context.Context, path string) (*raw_entity.HorseInfo, error)
	Write(ctx context.Context, path string, data *raw_entity.HorseInfo) error
	Fetch(ctx context.Context, url string) (*netkeiba_entity.Horse, error)
}
