package spreadsheet_entity

import (
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"strconv"
)

type TicketResult struct {
	raceCount     int
	betCount      int
	hitCount      int
	payment       int
	payout        int
	averagePayout int
	maxPayout     int
	minPayout     int
	payoutRate    string
}

func NewTicketResult(
	raceCount types.RaceCount,
	betCount types.BetCount,
	hitCount types.HitCount,
	payment types.Payment,
	payout types.Payout,
	averagePayout types.Payout,
	maxPayout types.Payout,
	minPayout types.Payout,
) *TicketResult {
	payoutRate := "0%"
	if payment > 0 {
		payoutRate = fmt.Sprintf("%s%s", strconv.FormatFloat((float64(payout)*float64(100))/float64(payment), 'f', 2, 64), "%")
	}
	return &TicketResult{
		raceCount:     raceCount.Value(),
		betCount:      betCount.Value(),
		hitCount:      hitCount.Value(),
		payment:       payment.Value(),
		payout:        payout.Value(),
		averagePayout: averagePayout.Value(),
		maxPayout:     maxPayout.Value(),
		minPayout:     minPayout.Value(),
		payoutRate:    payoutRate,
	}
}

func (s *TicketResult) RaceCount() int {
	return s.raceCount
}

func (s *TicketResult) BetCount() int {
	return s.betCount
}

func (s *TicketResult) HitCount() int {
	return s.hitCount
}

func (s *TicketResult) Payment() int {
	return s.payment
}

func (s *TicketResult) Payout() int {
	return s.payout
}

func (s *TicketResult) AveragePayout() int {
	return s.averagePayout
}

func (s *TicketResult) MaxPayout() int {
	return s.maxPayout
}

func (s *TicketResult) MinPayout() int {
	return s.minPayout
}

func (s *TicketResult) PayoutRate() string {
	return s.payoutRate
}
