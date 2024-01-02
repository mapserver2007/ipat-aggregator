package spreadsheet_entity

import (
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"strconv"
)

type ShortSummary struct {
	payment    int
	payout     int
	payoutRate string
}

func NewShortSummary(
	payment types.Payment,
	payout types.Payout,
) *ShortSummary {
	payoutRate := "0%"
	if payment > 0 {
		payoutRate = fmt.Sprintf("%s%s", strconv.FormatFloat((float64(payout)*float64(100))/float64(payment), 'f', 2, 64), "%")
	}
	return &ShortSummary{
		payment:    payment.Value(),
		payout:     payout.Value(),
		payoutRate: payoutRate,
	}
}

func (s *ShortSummary) Payment() int {
	return s.payment
}

func (s *ShortSummary) Payout() int {
	return s.payout
}

func (s *ShortSummary) PayoutRate() string {
	return s.payoutRate
}
