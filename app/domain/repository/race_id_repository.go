package repository

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
)

type RaceIdRepository interface {
	Read(ctx context.Context, path string) (*raw_entity.RaceIdInfo, error)
	Write(ctx context.Context, path string, data *raw_entity.RaceIdInfo) error
	Fetch(ctx context.Context, url string) ([]string, error)
}
