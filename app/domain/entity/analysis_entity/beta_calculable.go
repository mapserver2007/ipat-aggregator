package analysis_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"github.com/shopspring/decimal"
)

type BetaCalculable struct {
	raceId              types.RaceId
	raceDate            types.RaceDate
	marker              types.Marker
	markerCombinationId types.MarkerCombinationId
	odds                decimal.Decimal
	number              types.BetNumber
	popular             int
	entries             int
	jockeyId            types.JockeyId
	filters             []filter.AttributeId
}

func NewBetaCalculable(
	raceId types.RaceId,
	raceDate types.RaceDate,
	markerCombinationId types.MarkerCombinationId,
	odds decimal.Decimal,
	number types.BetNumber,
	popular int,
	entries int,
	jockeyId types.JockeyId,
	filters []filter.AttributeId,
) *BetaCalculable {
	marker, _ := types.NewMarker(markerCombinationId.Value() % 10)
	return &BetaCalculable{
		raceId:              raceId,
		raceDate:            raceDate,
		marker:              marker,
		markerCombinationId: markerCombinationId,
		odds:                odds,
		number:              number,
		popular:             popular,
		entries:             entries,
		jockeyId:            jockeyId,
		filters:             filters,
	}
}

func (n *BetaCalculable) RaceId() types.RaceId {
	return n.raceId
}

func (n *BetaCalculable) RaceDate() types.RaceDate {
	return n.raceDate
}

func (n *BetaCalculable) Marker() types.Marker {
	return n.marker
}

func (n *BetaCalculable) MarkerCombinationId() types.MarkerCombinationId {
	return n.markerCombinationId
}

func (n *BetaCalculable) Odds() decimal.Decimal {
	return n.odds
}

func (n *BetaCalculable) Number() types.BetNumber {
	return n.number
}

func (n *BetaCalculable) Popular() int {
	return n.popular
}

func (n *BetaCalculable) Entries() int {
	return n.entries
}

func (n *BetaCalculable) JockeyId() types.JockeyId {
	return n.jockeyId
}

func (n *BetaCalculable) Filters() []filter.AttributeId {
	return n.filters
}
