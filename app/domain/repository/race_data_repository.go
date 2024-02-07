package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
)

type RaceDataRepository interface {
	Read(ctx context.Context, filePath string) ([]*raw_entity.Race, error)
	Write(ctx context.Context, filePath string, raceInfo *raw_entity.RaceInfo) error
	Fetch(ctx context.Context, url string) (*netkeiba_entity.Race, error)
}
