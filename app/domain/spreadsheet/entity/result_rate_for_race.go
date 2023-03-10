package entity

import (
	"github.com/mapserver2007/tools/baken/app/domain/race/entity"
	"github.com/mapserver2007/tools/baken/app/domain/race/value_object"
)

type ResultRateForRace struct {
	Date value_object.RaceDate
	Race entity.Race
}
