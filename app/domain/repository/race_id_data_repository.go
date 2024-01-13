package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
)

type RaceIdDataRepository interface {
	Read(ctx context.Context, fileName string) ([]*raw_entity.RaceDate, []int, error)
	Write(ctx context.Context, fileName string, raceIdInfo *raw_entity.RaceIdInfo) error
	Fetch(ctx context.Context, url string) ([]string, error)
}
