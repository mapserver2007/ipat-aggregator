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
	payment      int
	payout       int
}

func NewTicket(
	rawRaceDate,
	rawEntryNo,
	rawRaceCourse,
	rawRaceNo,
	rawBetNumber,
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

	betNumber := types.NewBetNumber(rawBetNumber)

	ticketType := types.NewTicketType(rawTicketType)

	ticketResult := types.TicketUnHit
	if rawTicketResult {
		ticketResult = types.TicketHit
	}

	payment, err := strconv.Atoi(rawPayment)
	if err != nil {
		return nil, err
	}

	payout, err := strconv.Atoi(rawPayout)
	if err != nil {
		return nil, err
	}

	return &Ticket{
		raceDate:     raceDate,
		entryNo:      entryNo,
		raceCourse:   raceCourse,
		raceNo:       raceNo,
		betNumber:    betNumber,
		ticketType:   ticketType,
		ticketResult: ticketResult,
		payment:      payment,
		payout:       payout,
	}, nil
}
