package types

import (
	"errors"
)

type RaceTime string

func NewRaceTime(s string) (RaceTime, error) {
	if len(s) != 22 {
		return "", errors.New("Invalid raceTimeId format: " + s)
	}
	return RaceTime(s), nil
}

func (r RaceTime) Value() string {
	return string(r)
}
