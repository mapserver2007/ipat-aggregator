package summary_service

import (
	"context"
	"time"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type Term interface {
	Create(ctx context.Context, input *TermInput) *TermOutput
}

type TermInput struct {
	Tickets  []*ticket_csv_entity.RaceTicket
	DateFrom time.Time
	DateTo   time.Time
}

type TermOutput struct {
	RaceCount     types.RaceCount
	BetCount      types.BetCount
	HitCount      types.HitCount
	Payment       types.Payment
	Payout        types.Payout
	AveragePayout types.Payout
	MaxPayout     types.Payout
	MinPayout     types.Payout
}

type termService struct{}

func NewTerm() Term {
	return &termService{}
}

func (t *termService) Create(ctx context.Context, input *TermInput) *TermOutput {
	from := input.DateFrom
	to := input.DateTo

	var timeRangeTickets []*ticket_csv_entity.Ticket
	raceIdTicketsMap := map[types.RaceId][]*ticket_csv_entity.Ticket{}

	for _, raceTicket := range input.Tickets {
		if raceTicket.Ticket().RaceDate().Date().Unix() >= from.Unix() && raceTicket.Ticket().RaceDate().Date().Unix() < to.Unix() {
			timeRangeTickets = append(timeRangeTickets, raceTicket.Ticket())
			if _, ok := raceIdTicketsMap[raceTicket.RaceId()]; !ok {
				raceIdTicketsMap[raceTicket.RaceId()] = make([]*ticket_csv_entity.Ticket, 0)
			}

			raceIdTicketsMap[raceTicket.RaceId()] = append(raceIdTicketsMap[raceTicket.RaceId()], raceTicket.Ticket())
		}
	}

	var hitTickets []*ticket_csv_entity.Ticket
	for _, ticket := range timeRangeTickets {
		if ticket.TicketResult() == types.TicketHit {
			hitTickets = append(hitTickets, ticket)
		}
	}

	hitCount := len(hitTickets)
	var (
		sumPayment    int
		sumPayout     int
		maxPayout     int
		minPayout     int
		averagePayout int
	)
	for _, ticket := range hitTickets {
		if maxPayout < ticket.Payout().Value() {
			maxPayout = ticket.Payout().Value()
		}
		if minPayout == 0 || minPayout > ticket.Payout().Value() {
			minPayout = ticket.Payout().Value()
		}
	}
	for _, ticket := range timeRangeTickets {
		sumPayment += ticket.Payment().Value()
		sumPayout += ticket.Payout().Value()
	}

	betCount := len(timeRangeTickets)
	raceCount := len(raceIdTicketsMap)
	if len(hitTickets) > 0 {
		averagePayout = int(float64(sumPayout) / float64(len(hitTickets)))
	}

	return &TermOutput{
		RaceCount:     types.RaceCount(raceCount),
		BetCount:      types.BetCount(betCount),
		HitCount:      types.HitCount(hitCount),
		Payment:       types.Payment(sumPayment),
		Payout:        types.Payout(sumPayout),
		AveragePayout: types.Payout(averagePayout),
		MaxPayout:     types.Payout(maxPayout),
		MinPayout:     types.Payout(minPayout),
	}
}
