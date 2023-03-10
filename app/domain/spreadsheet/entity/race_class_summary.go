package entity

import (
	"github.com/mapserver2007/tools/baken/app/domain/race/value_object"
)

func NewRaceClassSummary(raceClassSummary map[value_object.GradeClass]ResultRate) RaceClassSummary {
	return RaceClassSummary{
		RaceClassRates: raceClassSummary,
	}
}
