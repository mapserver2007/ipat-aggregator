package umaca_csv_entity

import (
	"strconv"

	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type UmacaMaster struct {
	raceDate   types.RaceDate
	raceId     types.RaceId
	betNumber  types.BetNumber
	ticketType types.TicketType
	payment    types.Payment
}

func NewUmacaMaster(
	rawRaceDate,
	rawRaceId,
	rawBetNumber,
	rawTicketType,
	rawPayment string,
) (*UmacaMaster, error) {
	raceDate, err := types.NewRaceDate(rawRaceDate)
	if err != nil {
		return nil, err
	}

	payment, err := strconv.Atoi(rawPayment)
	if err != nil {
		return nil, err
	}

	return &UmacaMaster{
		raceDate:   raceDate,
		raceId:     types.RaceId(rawRaceId),
		betNumber:  types.NewBetNumber(rawBetNumber),
		ticketType: types.NewTicketType(rawTicketType),
		payment:    types.Payment(payment),
	}, nil
}

func (u *UmacaMaster) RaceDate() types.RaceDate {
	return u.raceDate
}

func (u *UmacaMaster) RaceId() types.RaceId {
	return u.raceId
}

func (u *UmacaMaster) BetNumber() types.BetNumber {
	return u.betNumber
}

func (u *UmacaMaster) TicketType() types.TicketType {
	return u.ticketType
}

func (u *UmacaMaster) Payment() types.Payment {
	return u.payment
}
