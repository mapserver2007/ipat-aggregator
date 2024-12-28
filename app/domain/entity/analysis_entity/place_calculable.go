package analysis_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"github.com/shopspring/decimal"
)

type PlaceCalculable struct {
	raceId              types.RaceId
	raceDate            types.RaceDate
	marker              types.Marker
	markerCombinationId types.MarkerCombinationId
	odds                decimal.Decimal
	number              types.BetNumber
	popular             int
	orderNo             int
	entries             int
	jockeyId            types.JockeyId
	filters             []filter.Id
}

func NewPlaceCalculable(
	raceId types.RaceId,
	raceDate types.RaceDate,
	markerCombinationId types.MarkerCombinationId,
	odds decimal.Decimal,
	number types.BetNumber,
	popular int,
	orderNo int,
	entries int,
	jockeyId types.JockeyId,
	filters []filter.Id,
) *PlaceCalculable {
	marker, _ := types.NewMarker(markerCombinationId.Value() % 10)
	return &PlaceCalculable{
		raceId:              raceId,
		raceDate:            raceDate,
		marker:              marker,
		markerCombinationId: markerCombinationId,
		odds:                odds,
		number:              number,
		popular:             popular,
		orderNo:             orderNo,
		entries:             entries,
		jockeyId:            jockeyId,
		filters:             filters,
	}
}

func (n *PlaceCalculable) RaceId() types.RaceId {
	return n.raceId
}

func (n *PlaceCalculable) RaceDate() types.RaceDate {
	return n.raceDate
}

func (n *PlaceCalculable) Marker() types.Marker {
	return n.marker
}

func (n *PlaceCalculable) MarkerCombinationId() types.MarkerCombinationId {
	return n.markerCombinationId
}

func (n *PlaceCalculable) Odds() decimal.Decimal {
	return n.odds
}

func (n *PlaceCalculable) Number() types.BetNumber {
	return n.number
}

func (n *PlaceCalculable) Popular() int {
	return n.popular
}

func (n *PlaceCalculable) OrderNo() int {
	return n.orderNo
}

func (n *PlaceCalculable) Entries() int {
	return n.entries
}

func (n *PlaceCalculable) JockeyId() types.JockeyId {
	return n.jockeyId
}

func (n *PlaceCalculable) Filters() []filter.Id {
	return n.filters
}
