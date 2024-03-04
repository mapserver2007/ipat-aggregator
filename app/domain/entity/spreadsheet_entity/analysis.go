package spreadsheet_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type AnalysisData struct {
	markerCombinationMapByFilter map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis
	raceCountMapByFilter         map[filter.Id]map[types.MarkerCombinationId]map[types.OddsRangeType]int
	allMarkerCombinationIds      []types.MarkerCombinationId
}

func NewAnalysisData(
	markerCombinationMapByFilter map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis,
	raceCountMapByFilter map[filter.Id]map[types.MarkerCombinationId]map[types.OddsRangeType]int,
	allMarkerCombinationIds []types.MarkerCombinationId,
) *AnalysisData {
	return &AnalysisData{
		markerCombinationMapByFilter: markerCombinationMapByFilter,
		raceCountMapByFilter:         raceCountMapByFilter,
		allMarkerCombinationIds:      allMarkerCombinationIds,
	}
}

func (a *AnalysisData) MarkerCombinationMapByFilter() map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis {
	return a.markerCombinationMapByFilter
}

func (a *AnalysisData) RaceCountMapByFilter() map[filter.Id]map[types.MarkerCombinationId]map[types.OddsRangeType]int {
	return a.raceCountMapByFilter
}

func (a *AnalysisData) AllMarkerCombinationIds() []types.MarkerCombinationId {
	return a.allMarkerCombinationIds
}

type MarkerCombinationAnalysis struct {
	raceCountOddsRangeMap map[types.OddsRangeType]int
	calculables           []*analysis_entity.Calculable
}

func NewMarkerCombinationAnalysis(raceCountOddsRangeMap map[types.OddsRangeType]int) *MarkerCombinationAnalysis {
	return &MarkerCombinationAnalysis{
		raceCountOddsRangeMap: raceCountOddsRangeMap,
		calculables:           make([]*analysis_entity.Calculable, 0),
	}
}

func (m *MarkerCombinationAnalysis) AddCalculable(calculable *analysis_entity.Calculable) {
	m.calculables = append(m.calculables, calculable)
}

func (m *MarkerCombinationAnalysis) Calculables() []*analysis_entity.Calculable {
	return m.calculables
}
