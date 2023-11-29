package types

import "fmt"

type RaceId string

func NewRaceIdForJRA(
	year int,
	day int,
	raceCourse int,
	raceRound int,
	raceNo int,
) RaceId {
	rawRaceId := fmt.Sprintf("%d%02d%02d%02d%02d", year, raceCourse, raceRound, day, raceNo)
	return RaceId(rawRaceId)
}

func NewRaceIdForNAR(
	year int,
	month int,
	day int,
	raceCourse int,
	raceNo int,
) RaceId {
	rawRaceId := fmt.Sprintf("%d%02d%02d%02d%02d", year, raceCourse, month, day, raceNo)
	return RaceId(rawRaceId)
}

func NewRaceIdForOverseas(
	year int,
	month int,
	day int,
	raceCourse int,
	raceNo int,
) RaceId {
	rawRaceId := fmt.Sprintf("%d%02d%02d%02d%02d", year, raceCourse, month, day, raceNo)
	return RaceId(rawRaceId)
}
