package prediction_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type RaceEntryHorse struct {
	horseId       types.HorseId
	horseName     string
	bracketNumber int
	horseNumber   types.HorseNumber
	jockeyId      types.JockeyId
	raceWeight    float64
}

func NewRaceEntryHorse(
	rawHorseId string,
	horseName string,
	bracketNumber int,
	rawHorseNumber int,
	rawJockeyId string,
	raceWeight float64,
) *RaceEntryHorse {
	return &RaceEntryHorse{
		horseId:       types.HorseId(rawHorseId),
		horseName:     horseName,
		bracketNumber: bracketNumber,
		horseNumber:   types.HorseNumber(rawHorseNumber),
		jockeyId:      types.JockeyId(rawJockeyId),
		raceWeight:    raceWeight,
	}
}

func (r *RaceEntryHorse) HorseId() types.HorseId {
	return r.horseId
}

func (r *RaceEntryHorse) HorseName() string {
	return r.horseName
}

func (r *RaceEntryHorse) BracketNumber() int {
	return r.bracketNumber
}

func (r *RaceEntryHorse) HorseNumber() types.HorseNumber {
	return r.horseNumber
}

func (r *RaceEntryHorse) JockeyId() types.JockeyId {
	return r.jockeyId
}

func (r *RaceEntryHorse) RaceWeight() float64 {
	return r.raceWeight
}
