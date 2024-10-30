package analysis_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
)

type PlaceUnHitRace struct {
	raceId      types.RaceId
	raceDate    types.RaceDate
	horseNumber types.HorseNumber
	marker      types.Marker
	odds        decimal.Decimal
	orderNo     int
	entries     int
	jockeyId    types.JockeyId
	horseId     types.HorseId
}

func NewPlaceUnHitRace(
	raceId types.RaceId,
	raceDate types.RaceDate,
	horseNumber types.HorseNumber,
	marker types.Marker,
	odds decimal.Decimal,
	orderNo int,
	entries int,
	jockeyId types.JockeyId,
	horseId types.HorseId,
) *PlaceUnHitRace {
	return &PlaceUnHitRace{
		raceId:      raceId,
		raceDate:    raceDate,
		horseNumber: horseNumber,
		marker:      marker,
		odds:        odds,
		orderNo:     orderNo,
		entries:     entries,
		jockeyId:    jockeyId,
		horseId:     horseId,
	}
}

func (n *PlaceUnHitRace) RaceId() types.RaceId {
	return n.raceId
}

func (n *PlaceUnHitRace) RaceDate() types.RaceDate {
	return n.raceDate
}

func (n *PlaceUnHitRace) HorseNumber() types.HorseNumber {
	return n.horseNumber
}

func (n *PlaceUnHitRace) Marker() types.Marker {
	return n.marker
}

func (n *PlaceUnHitRace) Odds() decimal.Decimal {
	return n.odds
}

func (n *PlaceUnHitRace) OrderNo() int {
	return n.orderNo
}

func (n *PlaceUnHitRace) Entries() int {
	return n.entries
}

func (n *PlaceUnHitRace) JockeyId() types.JockeyId {
	return n.jockeyId
}

func (n *PlaceUnHitRace) HorseId() types.HorseId {
	return n.horseId
}
