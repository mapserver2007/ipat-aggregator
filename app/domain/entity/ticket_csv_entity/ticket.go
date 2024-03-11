package ticket_csv_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"strconv"
)

type Ticket struct {
	raceDate     types.RaceDate
	entryNo      int
	raceCourse   types.RaceCourse
	raceNo       int
	betNumber    types.BetNumber
	ticketType   types.TicketType
	ticketResult types.TicketResult
	payment      types.Payment
	payout       types.Payout
}

func NewTicket(
	betNumber types.BetNumber,
	rawRaceDate,
	rawEntryNo,
	rawRaceCourse,
	rawRaceNo,
	rawTicketType string,
	rawTicketResult bool,
	rawPayment,
	rawPayout string,
) (*Ticket, error) {
	raceDate, err := types.NewRaceDate(rawRaceDate)
	if err != nil {
		return nil, err
	}

	entryNo, err := strconv.Atoi(rawEntryNo)
	if err != nil {
		return nil, err
	}

	raceCourse := types.NewRaceCourse(rawRaceCourse)

	raceNo, err := strconv.Atoi(rawRaceNo)
	if err != nil {
		return nil, err
	}

	ticketType := types.NewTicketType(rawTicketType)

	ticketResult := types.TicketUnHit
	if rawTicketResult {
		ticketResult = types.TicketHit
	}

	payment, err := strconv.Atoi(rawPayment)
	if err != nil {
		return nil, err
	}

	payout := 0
	if rawPayout != "" {
		payout, err = strconv.Atoi(rawPayout)
		if err != nil {
			return nil, err
		}
	}

	return &Ticket{
		raceDate:     raceDate,
		entryNo:      entryNo,
		raceCourse:   raceCourse,
		raceNo:       raceNo,
		betNumber:    betNumber,
		ticketType:   ticketType,
		ticketResult: ticketResult,
		payment:      types.Payment(payment),
		payout:       types.Payout(payout),
	}, nil
}

func (t *Ticket) RaceDate() types.RaceDate {
	return t.raceDate
}

func (t *Ticket) EntryNo() int {
	return t.entryNo
}

func (t *Ticket) RaceCourse() types.RaceCourse {
	return t.raceCourse
}

func (t *Ticket) RaceNo() int {
	return t.raceNo
}

func (t *Ticket) BetNumber() types.BetNumber {
	return t.betNumber
}

func (t *Ticket) TicketType() types.TicketType {
	return t.ticketType
}

func (t *Ticket) TicketResult() types.TicketResult {
	return t.ticketResult
}

func (t *Ticket) Payment() types.Payment {
	return t.payment
}

func (t *Ticket) Payout() types.Payout {
	return t.payout
}
