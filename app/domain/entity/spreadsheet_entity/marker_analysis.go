package spreadsheet_entity

import (
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"github.com/shopspring/decimal"
	"strconv"
)

type AnalysisData struct {
	hitDataMapByFilter      map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis
	unHitDataMapByFilter    map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis
	raceCountByFilter       map[filter.Id]int
	allMarkerCombinationIds []types.MarkerCombinationId
}

func NewAnalysisData(
	hitDataMapByFilter map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis,
	unHitDataMapByFilter map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis,
	raceCountByFilter map[filter.Id]int,
	allMarkerCombinationIds []types.MarkerCombinationId,
) *AnalysisData {
	return &AnalysisData{
		hitDataMapByFilter:      hitDataMapByFilter,
		unHitDataMapByFilter:    unHitDataMapByFilter,
		raceCountByFilter:       raceCountByFilter,
		allMarkerCombinationIds: allMarkerCombinationIds,
	}
}

func (a *AnalysisData) HitDataMapByFilter() map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis {
	return a.hitDataMapByFilter
}

func (a *AnalysisData) UnHitDataMapByFilter() map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis {
	return a.unHitDataMapByFilter
}

func (a *AnalysisData) RaceCountByFilter() map[filter.Id]int {
	return a.raceCountByFilter
}

func (a *AnalysisData) AllMarkerCombinationIds() []types.MarkerCombinationId {
	return a.allMarkerCombinationIds
}

type MarkerCombinationAnalysis struct {
	raceCount   int
	calculables []*analysis_entity.Calculable
}

func NewMarkerCombinationAnalysis(raceCount int) *MarkerCombinationAnalysis {
	return &MarkerCombinationAnalysis{
		raceCount:   raceCount,
		calculables: make([]*analysis_entity.Calculable, 0),
	}
}

func (m *MarkerCombinationAnalysis) AddCalculable(calculable *analysis_entity.Calculable) {
	m.calculables = append(m.calculables, calculable)
}

func (m *MarkerCombinationAnalysis) MatchRate() float64 {
	return (float64(m.MatchCount()) * float64(100)) / float64(m.raceCount)
}

func (m *MarkerCombinationAnalysis) MatchRateFormat() string {
	return rateFormat(m.MatchRate())
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

func rateFormat(rate float64) string {
	return fmt.Sprintf("%s%s", strconv.FormatFloat(rate, 'f', 2, 64), "%")
}
