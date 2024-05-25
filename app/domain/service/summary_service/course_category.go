package summary_service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type CourseCategory interface {
	Create(ctx context.Context, input *CourseCategoryInput) *CourseCategoryOutput
}

type CourseCategoryInput struct {
	Tickets          []*ticket_csv_entity.RaceTicket
	Races            []*data_cache_entity.Race
	CourseCategories []types.CourseCategory
}

type CourseCategoryOutput struct {
	RaceCount     types.RaceCount
	BetCount      types.BetCount
	HitCount      types.HitCount
	Payment       types.Payment
	Payout        types.Payout
	AveragePayout types.Payout
	MaxPayout     types.Payout
	MinPayout     types.Payout
}

type courseCategoryService struct{}

func NewCourseCategory() CourseCategory {
	return &courseCategoryService{}
}

func (c *courseCategoryService) Create(ctx context.Context, input *CourseCategoryInput) *CourseCategoryOutput {
	var courseCategoryTickets []*ticket_csv_entity.Ticket
	raceIdTicketsMap := map[types.RaceId][]*ticket_csv_entity.Ticket{}
	raceMap := map[types.RaceId]types.CourseCategory{}
	for _, race := range input.Races {
		raceMap[race.RaceId()] = race.CourseCategory()
	}

	for _, raceTicket := range input.Tickets {
		class, ok := raceMap[raceTicket.RaceId()]
		if ok {
			if c.containsInSlice(input.CourseCategories, class) {
				courseCategoryTickets = append(courseCategoryTickets, raceTicket.Ticket())
				if _, ok := raceIdTicketsMap[raceTicket.RaceId()]; !ok {
					raceIdTicketsMap[raceTicket.RaceId()] = make([]*ticket_csv_entity.Ticket, 0)
				}
				raceIdTicketsMap[raceTicket.RaceId()] = append(raceIdTicketsMap[raceTicket.RaceId()], raceTicket.Ticket())
			}
		}
	}

	var hitTickets []*ticket_csv_entity.Ticket
	for _, ticket := range courseCategoryTickets {
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
	for _, ticket := range courseCategoryTickets {
		sumPayment += ticket.Payment().Value()
		sumPayout += ticket.Payout().Value()
	}

	betCount := len(courseCategoryTickets)
	raceCount := len(raceIdTicketsMap)
	if len(hitTickets) > 0 {
		averagePayout = int(float64(sumPayout) / float64(len(hitTickets)))
	}

	return &CourseCategoryOutput{
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

func (c *courseCategoryService) containsInSlice(slice []types.CourseCategory, courseCategory types.CourseCategory) bool {
	for _, c := range slice {
		if c == courseCategory {
			return true
		}
	}
	return false
}
