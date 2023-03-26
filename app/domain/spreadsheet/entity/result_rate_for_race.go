package entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/race/entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
)

type ResultRateForRace struct {
	Date value_object.RaceDate
	Race entity.Race
}
