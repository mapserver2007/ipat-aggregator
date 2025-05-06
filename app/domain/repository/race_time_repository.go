package repository

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
)

type RaceTimeRepository interface {
	List(ctx context.Context, path string) ([]string, error)
	Read(ctx context.Context, path string) ([]*raw_entity.RaceTime, error)
	Write(ctx context.Context, path string, raceTimeInfo *raw_entity.RaceTimeInfo) error
	Fetch(ctx context.Context, url string) (*netkeiba_entity.RaceTime, error)
}
