package entity

import (
	race_entity "github.com/mapserver2007/tools/baken/app/domain/race/entity"
)

func NewRaceResultSummary(
	resultSummary ResultSummary,
	race race_entity.Race,
	raceHandicappingSummary RaceHandicappingSummary,
) RaceResultSummary {
	return RaceResultSummary{
		ResultSummary:           resultSummary,
		Race:                    race,
		RaceHandicappingSummary: raceHandicappingSummary,
	}
}
