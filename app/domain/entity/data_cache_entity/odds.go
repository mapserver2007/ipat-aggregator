package data_cache_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type Odds struct {
	raceId        types.RaceId
	raceDate      types.RaceDate
	ticketType    types.TicketType
	number        types.BetNumber
	popularNumber int
	odds          string
}

func NewOdds(
	raceId types.RaceId,
	raceDate types.RaceDate,
	ticketType types.TicketType,
	number types.BetNumber,
	popularNumber int,
	odds string,
) *Odds {
	return &Odds{
		raceId:        raceId,
		raceDate:      raceDate,
		ticketType:    ticketType,
		number:        number,
		popularNumber: popularNumber,
		odds:          odds,
	}
}

func (o *Odds) RaceId() types.RaceId {
	return o.raceId
}

func (o *Odds) RaceDate() types.RaceDate {
	return o.raceDate
}

func (o *Odds) TicketType() types.TicketType {
	return o.ticketType
}

func (o *Odds) Number() types.BetNumber {
	return o.number
}

func (o *Odds) PopularNumber() int {
	return o.popularNumber

}

func (o *Odds) Odds() string {
	return o.odds
}
