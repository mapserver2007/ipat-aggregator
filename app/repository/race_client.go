package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/race/raw_entity"
)

type RaceClient interface {
	GetRacingNumbers(ctx context.Context, url string) ([]*raw_entity.RawRacingNumberNetkeiba, error)
	GetRaceResult(ctx context.Context, url string) (*raw_entity.RawRaceNetkeiba, error)
}
