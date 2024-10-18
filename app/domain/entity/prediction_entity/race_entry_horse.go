package prediction_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type RaceEntryHorse struct {
	horseId       types.HorseId
	horseName     string
	bracketNumber int
	horseNumber   types.HorseNumber
	jockeyId      types.JockeyId
}

func NewRaceEntryHorse(
	rawHorseId string,
	horseName string,
	bracketNumber int,
	rawHorseNumber int,
	rawJockeyId int,
) *RaceEntryHorse {
	return &RaceEntryHorse{
		horseId:       types.HorseId(rawHorseId),
		horseName:     horseName,
		bracketNumber: bracketNumber,
		horseNumber:   types.HorseNumber(rawHorseNumber),
		jockeyId:      types.JockeyId(rawJockeyId),
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
