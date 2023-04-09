package repository

import (
	"context"
	raw_jockey_entity "github.com/mapserver2007/ipat-aggregator/app/domain/jockey/raw_entity"
	raw_race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/raw_entity"
)

type RaceClient interface {
	GetRacingNumbers(ctx context.Context, url string) ([]*raw_race_entity.RawRacingNumberNetkeiba, error)
	GetRaceResult(ctx context.Context, url string) (*raw_race_entity.RawRaceNetkeiba, error)
	GetJockey(ctx context.Context, url string) (*raw_jockey_entity.RawJockeyNetkeiba, error)
}
