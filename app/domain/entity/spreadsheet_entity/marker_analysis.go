package spreadsheet_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"github.com/shopspring/decimal"
)

type AnalysisData struct {
	hitDataMapByFilter      map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis
	unHitDataMapByFilter    map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis
	raceCountMapByFilter    map[filter.Id]map[types.MarkerCombinationId]map[types.OddsRangeType]int
	allMarkerCombinationIds []types.MarkerCombinationId
}

func NewAnalysisData(
	hitDataMapByFilter map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis,
	unHitDataMapByFilter map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis,
	raceCountMapByFilter map[filter.Id]map[types.MarkerCombinationId]map[types.OddsRangeType]int,
	allMarkerCombinationIds []types.MarkerCombinationId,
) *AnalysisData {
	return &AnalysisData{
		hitDataMapByFilter:      hitDataMapByFilter,
		unHitDataMapByFilter:    unHitDataMapByFilter,
		raceCountMapByFilter:    raceCountMapByFilter,
		allMarkerCombinationIds: allMarkerCombinationIds,
	}
}

func (a *AnalysisData) HitDataMapByFilter() map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis {
	return a.hitDataMapByFilter
}

func (a *AnalysisData) UnHitDataMapByFilter() map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis {
	return a.unHitDataMapByFilter
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

func (m *MarkerCombinationAnalysis) MatchCount() int {
	var odds []decimal.Decimal
	for _, calculable := range m.calculables {
		odds = append(odds, calculable.Odds())
	}
	return len(odds)
}

func (m *MarkerCombinationAnalysis) Odds() []decimal.Decimal {
	var odds []decimal.Decimal
	for _, calculable := range m.calculables {
		odds = append(odds, calculable.Odds())
	}
	return odds
}
