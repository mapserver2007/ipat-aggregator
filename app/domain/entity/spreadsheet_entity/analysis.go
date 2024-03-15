package spreadsheet_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type AnalysisData struct {
	markerCombinationFilterMap  map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis
	oddsRangeRaceCountFilterMap map[filter.Id]map[types.MarkerCombinationId]map[types.OddsRangeType]int
	filters                     []filter.Id
	allMarkerCombinationIds     []types.MarkerCombinationId
}

func NewAnalysisData(
	markerCombinationFilterMap map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis,
	oddsRangeRaceCountFilterMap map[filter.Id]map[types.MarkerCombinationId]map[types.OddsRangeType]int,
	filters []filter.Id,
	allMarkerCombinationIds []types.MarkerCombinationId,
) *AnalysisData {
	return &AnalysisData{
		markerCombinationFilterMap:  markerCombinationFilterMap,
		oddsRangeRaceCountFilterMap: oddsRangeRaceCountFilterMap,
		filters:                     filters,
		allMarkerCombinationIds:     allMarkerCombinationIds,
	}
}

func (a *AnalysisData) MarkerCombinationFilterMap() map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis {
	return a.markerCombinationFilterMap
}

func (a *AnalysisData) OddsRangeRaceCountFilterMap() map[filter.Id]map[types.MarkerCombinationId]map[types.OddsRangeType]int {
	return a.oddsRangeRaceCountFilterMap
}

func (a *AnalysisData) Filters() []filter.Id {
	return a.filters
}

func (a *AnalysisData) AllMarkerCombinationIds() []types.MarkerCombinationId {
	return a.allMarkerCombinationIds
}

type MarkerCombinationAnalysis struct {
	raceCountOddsRangeMap map[types.OddsRangeType]int
	calculables           []*analysis_entity.Calculable
}

func NewMarkerCombinationAnalysis() *MarkerCombinationAnalysis {
	return &MarkerCombinationAnalysis{
		raceCountOddsRangeMap: map[types.OddsRangeType]int{},
		calculables:           make([]*analysis_entity.Calculable, 0),
	}
}

func (m *MarkerCombinationAnalysis) AddRaceCountOddsRangeMap(raceCountOddsRangeMap map[types.OddsRangeType]int) {
	m.raceCountOddsRangeMap = raceCountOddsRangeMap
}

func (m *MarkerCombinationAnalysis) AddCalculable(calculable *analysis_entity.Calculable) {
	m.calculables = append(m.calculables, calculable)
}

func (m *MarkerCombinationAnalysis) Calculables() []*analysis_entity.Calculable {
	return m.calculables
}
