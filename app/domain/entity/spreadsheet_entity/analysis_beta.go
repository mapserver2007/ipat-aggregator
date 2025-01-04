package spreadsheet_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type AnalysisBeta struct {
	raceCount              int
	filterId               filter.AttributeId
	markerCombinationRates map[types.MarkerCombinationId]*AnalysisBetaRate
}

func NewAnalysisBeta(
	raceCount int,
	filterId filter.AttributeId,
	markerCombinationRates map[types.MarkerCombinationId]*AnalysisBetaRate,
) *AnalysisBeta {
	return &AnalysisBeta{
		raceCount:              raceCount,
		filterId:               filterId,
		markerCombinationRates: markerCombinationRates,
	}
}

func (a *AnalysisBeta) RaceCount() int {
	return a.raceCount
}

func (a *AnalysisBeta) FilterId() filter.AttributeId {
	return a.filterId
}

func (a *AnalysisBeta) MarkerCombinationRates() map[types.MarkerCombinationId]*AnalysisBetaRate {
	return a.markerCombinationRates
}
