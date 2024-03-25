package spreadsheet_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"github.com/shopspring/decimal"
)

type AnalysisData struct {
	markerCombinationFilterMap map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis
	oddsRangeCountFilterMap    map[filter.Id]map[types.MarkerCombinationId]map[types.OddsRangeType]int
	raceCountFilterMap         map[filter.Id]map[types.TicketType]int
	filters                    []filter.Id
	allMarkerCombinationIds    []types.MarkerCombinationId
}

func NewAnalysisData(
	markerCombinationFilterMap map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis,
	oddsRangeCountFilterMap map[filter.Id]map[types.MarkerCombinationId]map[types.OddsRangeType]int,
	raceCountFilterMap map[filter.Id]map[types.TicketType]int,
	filters []filter.Id,
	allMarkerCombinationIds []types.MarkerCombinationId,
) *AnalysisData {
	return &AnalysisData{
		markerCombinationFilterMap: markerCombinationFilterMap,
		oddsRangeCountFilterMap:    oddsRangeCountFilterMap,
		raceCountFilterMap:         raceCountFilterMap,
		filters:                    filters,
		allMarkerCombinationIds:    allMarkerCombinationIds,
	}
}

// MarkerCombinationFilterMap 打たれた印(的中・不的中・オッズ等の情報)に関する情報。的中回数、オッズはここから算出する
func (a *AnalysisData) MarkerCombinationFilterMap() map[filter.Id]map[types.MarkerCombinationId]*MarkerCombinationAnalysis {
	return a.markerCombinationFilterMap
}

// OddsRangeCountFilterMap オッズ幅に対する情報。対象レースすべての情報を保持している(的中・不的中の合算)
func (a *AnalysisData) OddsRangeCountFilterMap() map[filter.Id]map[types.MarkerCombinationId]map[types.OddsRangeType]int {
	return a.oddsRangeCountFilterMap
}

// RaceCountFilterMap 券種単位のレース数
func (a *AnalysisData) RaceCountFilterMap() map[filter.Id]map[types.TicketType]int {
	return a.raceCountFilterMap
}

func (a *AnalysisData) Filters() []filter.Id {
	return a.filters
}

func (a *AnalysisData) AllMarkerCombinationIds() []types.MarkerCombinationId {
	return a.allMarkerCombinationIds
}

type MarkerCombinationAnalysis struct {
	calculables []*analysis_entity.Calculable
}

func NewMarkerCombinationAnalysis() *MarkerCombinationAnalysis {
	return &MarkerCombinationAnalysis{
		calculables: make([]*analysis_entity.Calculable, 0),
	}
}

func (m *MarkerCombinationAnalysis) AddCalculable(calculable *analysis_entity.Calculable) {
	m.calculables = append(m.calculables, calculable)
}

func (m *MarkerCombinationAnalysis) Calculables() []*analysis_entity.Calculable {
	return m.calculables
}

type Odds struct {
	ticketType types.TicketType
	odds       decimal.Decimal
	number     types.BetNumber
}

func NewOdds(
	ticketType types.TicketType,
	odds string,
	number types.BetNumber,
) *Odds {
	return &Odds{
		ticketType: ticketType,
		odds:       decimal.RequireFromString(odds),
		number:     number,
	}
}

func (m *Odds) TicketType() types.TicketType {
	return m.ticketType
}

func (m *Odds) Odds() decimal.Decimal {
	return m.odds
}

func (m *Odds) Number() types.BetNumber {
	return m.number
}
