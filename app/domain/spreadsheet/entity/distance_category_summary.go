package entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
)

func NewDistanceCategorySummary(distanceCategoryRates map[value_object.DistanceCategory]ResultRate) DistanceCategorySummary {
	return DistanceCategorySummary{
		DistanceCategoryRates: distanceCategoryRates,
	}
}
