package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
)

type RaceDataRepository interface {
	Read(ctx context.Context, fileName string) ([]*raw_entity.Race, error)
	Write(ctx context.Context) error
	Fetch(ctx context.Context) error
}
