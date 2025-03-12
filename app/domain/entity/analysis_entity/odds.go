package analysis_entity

import (
	"github.com/shopspring/decimal"

	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type Odds struct {
	raceId        types.RaceId
	raceDate      types.RaceDate
	ticketType    types.TicketType
	number        types.BetNumber
	popularNumber int
	odds          decimal.Decimal
}

func NewOdds(
	raceId types.RaceId,
	raceDate types.RaceDate,
	ticketType types.TicketType,
	number types.BetNumber,
	popularNumber int,
	odds []string,
) (*Odds, error) {
	oddsDecimal, err := decimal.NewFromString(odds[0])
	if err != nil {
		return nil, err
	}

	return &Odds{
		raceId:        raceId,
		raceDate:      raceDate,
		ticketType:    ticketType,
		number:        number,
		popularNumber: popularNumber,
		odds:          oddsDecimal,
	}, nil
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

func (o *Odds) Odds() decimal.Decimal {
	return o.odds
}
