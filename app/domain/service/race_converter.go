package service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type RaceConverter interface {
	RawRacingNumberToRawRacingNumberMap(ctx context.Context, racingNumbers []*raw_entity.RacingNumber) map[types.RacingNumberId]*raw_entity.RacingNumber
}

type raceConverter struct{}

func NewRaceConverter() RaceConverter {
	return &raceConverter{}
}

func (r *raceConverter) RawRacingNumberToRawRacingNumberMap(ctx context.Context, racingNumbers []*raw_entity.RacingNumber) map[types.RacingNumberId]*raw_entity.RacingNumber {
	return ConvertToMap(racingNumbers, func(racingNumber *raw_entity.RacingNumber) types.RacingNumberId {
		return types.NewRacingNumberId(
			types.RaceDate(racingNumber.Date),
			types.RaceCourse(racingNumber.RaceCourseId),
		)
	})
}
