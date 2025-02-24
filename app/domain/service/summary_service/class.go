package summary_service

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type Class interface {
	Create(ctx context.Context, input *ClassInput) *ClassOutput
}

type ClassInput struct {
	Tickets []*ticket_csv_entity.RaceTicket
	Races   []*data_cache_entity.Race
	Classes []types.GradeClass
}

type ClassOutput struct {
	RaceCount     types.RaceCount
	BetCount      types.BetCount
	HitCount      types.HitCount
	Payment       types.Payment
	Payout        types.Payout
	AveragePayout types.Payout
	MaxPayout     types.Payout
	MinPayout     types.Payout
}

type classService struct{}

func NewClass() Class {
	return &classService{}
}

func (c *classService) Create(ctx context.Context, input *ClassInput) *ClassOutput {
	var classTickets []*ticket_csv_entity.Ticket
	raceIdTicketsMap := map[types.RaceId][]*ticket_csv_entity.Ticket{}
	raceMap := map[types.RaceId]types.GradeClass{}
	for _, race := range input.Races {
		raceMap[race.RaceId()] = race.Class()
	}

	for _, raceTicket := range input.Tickets {
		class, ok := raceMap[raceTicket.RaceId()]
		if ok && c.containsInSlice(input.Classes, class) {
			classTickets = append(classTickets, raceTicket.Ticket())
			if _, ok := raceIdTicketsMap[raceTicket.RaceId()]; !ok {
				raceIdTicketsMap[raceTicket.RaceId()] = make([]*ticket_csv_entity.Ticket, 0)
			}
			raceIdTicketsMap[raceTicket.RaceId()] = append(raceIdTicketsMap[raceTicket.RaceId()], raceTicket.Ticket())
		}
	}

	var hitTickets []*ticket_csv_entity.Ticket
	for _, ticket := range classTickets {
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
	for _, ticket := range classTickets {
		sumPayment += ticket.Payment().Value()
		sumPayout += ticket.Payout().Value()
	}

	betCount := len(classTickets)
	raceCount := len(raceIdTicketsMap)
	if len(hitTickets) > 0 {
		averagePayout = int(float64(sumPayout) / float64(len(hitTickets)))
	}

	return &ClassOutput{
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

func (c *classService) containsInSlice(slice []types.GradeClass, class types.GradeClass) bool {
	for _, c := range slice {
		if c == class {
			return true
		}
	}
	return false
}
