package analysis_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"github.com/shopspring/decimal"
)

type PlaceAllInCalculable struct {
	raceId              types.RaceId
	raceDate            types.RaceDate
	marker              types.Marker
	markerCombinationId types.MarkerCombinationId
	winOdds             decimal.Decimal
	placeOddsRange      []decimal.Decimal
	fixedPlaceOdds      decimal.Decimal
	popular             int
	orderNo             int
	entries             int
	trackCondition      types.TrackCondition
	raceCourse          types.RaceCourse
	distance            int
	jockeyId            types.JockeyId
	filters             []filter.Id
}

func NewPlaceAllInCalculable(
	raceId types.RaceId,
	raceDate types.RaceDate,
	markerCombinationId types.MarkerCombinationId,
	winOdds string,
	placeOddsRange []string,
	fixedPlaceOdds string,
	popular int,
	orderNo int,
	entries int,
	trackCondition types.TrackCondition,
	raceCourse types.RaceCourse,
	distance int,
	jockeyId types.JockeyId,
	filters []filter.Id,
) *PlaceAllInCalculable {
	decimalWinOdds, _ := decimal.NewFromString(winOdds)
	decimalPlaceOddsRange := make([]decimal.Decimal, 0, 2)
	for _, odds := range placeOddsRange {
		decimalOdds, _ := decimal.NewFromString(odds)
		decimalPlaceOddsRange = append(decimalPlaceOddsRange, decimalOdds)
	}
	decimalFixedPlaceOdds, _ := decimal.NewFromString(fixedPlaceOdds)
	marker, _ := types.NewMarker(markerCombinationId.Value() % 10)

	return &PlaceAllInCalculable{
		raceId:              raceId,
		raceDate:            raceDate,
		marker:              marker,
		markerCombinationId: markerCombinationId,
		winOdds:             decimalWinOdds,
		placeOddsRange:      decimalPlaceOddsRange,
		fixedPlaceOdds:      decimalFixedPlaceOdds,
		popular:             popular,
		orderNo:             orderNo,
		entries:             entries,
		trackCondition:      trackCondition,
		raceCourse:          raceCourse,
		distance:            distance,
		jockeyId:            jockeyId,
		filters:             filters,
	}
}

func (n *PlaceAllInCalculable) RaceId() types.RaceId {
	return n.raceId
}

func (n *PlaceAllInCalculable) RaceDate() types.RaceDate {
	return n.raceDate
}

func (n *PlaceAllInCalculable) Marker() types.Marker {
	return n.marker
}

func (n *PlaceAllInCalculable) MarkerCombinationId() types.MarkerCombinationId {
	return n.markerCombinationId
}

func (n *PlaceAllInCalculable) WinOdds() decimal.Decimal {
	return n.winOdds
}

func (n *PlaceAllInCalculable) PlaceOddsRange() []decimal.Decimal {
	return n.placeOddsRange
}

func (n *PlaceAllInCalculable) FixedPlaceOdds() decimal.Decimal {
	return n.fixedPlaceOdds
}

func (n *PlaceAllInCalculable) Popular() int {
	return n.popular
}

func (n *PlaceAllInCalculable) OrderNo() int {
	return n.orderNo
}

func (n *PlaceAllInCalculable) Entries() int {
	return n.entries
}

func (n *PlaceAllInCalculable) TrackCondition() types.TrackCondition {
	return n.trackCondition
}

func (n *PlaceAllInCalculable) RaceCourse() types.RaceCourse {
	return n.raceCourse
}

func (n *PlaceAllInCalculable) Distance() int {
	return n.distance
}

func (n *PlaceAllInCalculable) JockeyId() types.JockeyId {
	return n.jockeyId
}

func (n *PlaceAllInCalculable) Filters() []filter.Id {
	return n.filters
}
