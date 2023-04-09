package repository

import (
	"context"
	raw_jockey_entity "github.com/mapserver2007/ipat-aggregator/app/domain/jockey/raw_entity"
	raw_race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/raw_entity"
)

type RaceDB interface {
	ReadRaceInfo(ctx context.Context) (*raw_race_entity.RaceInfo, error)
	ReadRacingNumberInfo(ctx context.Context) (*raw_race_entity.RacingNumberInfo, error)
	ReadJockeyInfo(ctx context.Context) (*raw_jockey_entity.JockeyInfo, error)
	WriteRaceInfo(ctx context.Context, raceInfo *raw_race_entity.RaceInfo) error
	WriteRacingNumberInfo(ctx context.Context, racingNumberInfo *raw_race_entity.RacingNumberInfo) error
	WriteJockeyInfo(ctx context.Context, jockeyInfo *raw_jockey_entity.JockeyInfo) error
}
