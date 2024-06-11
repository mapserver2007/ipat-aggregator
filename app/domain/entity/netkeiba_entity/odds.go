package netkeiba_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type Odds struct {
	ticketType    types.TicketType
	odds          string
	popularNumber int
	horseNumbers  []types.HorseNumber
	raceDate      types.RaceDate
}

func NewOdds(
	ticketType types.TicketType,
	odds string,
	popularNumber int,
	horseNumbers []types.HorseNumber,
	raceDate types.RaceDate,
) *Odds {
	return &Odds{
		ticketType:    ticketType,
		odds:          odds,
		popularNumber: popularNumber,
		horseNumbers:  horseNumbers,
		raceDate:      raceDate,
	}
}

func (o *Odds) TicketType() types.TicketType {
	return o.ticketType
}

func (o *Odds) Odds() string {
	return o.odds
}

func (o *Odds) PopularNumber() int {
	return o.popularNumber
}

func (o *Odds) HorseNumbers() []types.HorseNumber {
	return o.horseNumbers
}

func (o *Odds) RaceDate() types.RaceDate {
	return o.raceDate
}
