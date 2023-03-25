package repository

import (
	"context"
	raw_race_entity "github.com/mapserver2007/tools/baken/app/domain/race/raw_entity"
)

type RaceDB interface {
	ReadRaceInfo(ctx context.Context) (*raw_race_entity.RaceInfo, error)
	ReadRacingNumberInfo(ctx context.Context) (*raw_race_entity.RacingNumberInfo, error)
	WriteRaceInfo(ctx context.Context, raceInfo *raw_race_entity.RaceInfo) error
	WriteRacingNumberInfo(ctx context.Context, racingNumberInfo *raw_race_entity.RacingNumberInfo) error
}
