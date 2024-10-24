package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
)

type HorseRepository interface {
	Fetch(ctx context.Context, url string) (*netkeiba_entity.Horse, error)
}
