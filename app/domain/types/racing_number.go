package types

import "fmt"

type RacingNumberId string

func NewRacingNumberId(
	date RaceDate,
	raceCourse RaceCourse,
) RacingNumberId {
	return RacingNumberId(fmt.Sprintf("%d_%d", date, raceCourse.Value()))
}
