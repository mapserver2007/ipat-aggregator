package value_object

import (
	"fmt"
	"strconv"
	"strings"
)

type RacingNumberId string

func NewRacingNumberId(date RaceDate, raceCourse RaceCourse) RacingNumberId {
	return RacingNumberId(fmt.Sprintf("%d_%d", date, raceCourse.Value()))
}

func (r *RacingNumberId) Date() RaceDate {
	date, _ := strconv.Atoi(strings.Split(string(*r), "_")[0])
	return RaceDate(date)
}

func (r *RacingNumberId) RaceCourse() RaceCourse {
	raceCourse, _ := strconv.Atoi(strings.Split(string(*r), "_")[1])
	return RaceCourse(raceCourse)
}
