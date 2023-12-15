package spreadsheet_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type ShortSummary struct {
	payment int
	payout  int
}

func NewShortSummary(
	payment types.Payment,
	payout types.Payout,
) *ShortSummary {
	return &ShortSummary{
		payment: payment.Value(),
		payout:  payout.Value(),
	}
}

func (s *ShortSummary) GetPayment() int {
	return s.payment
}

func (s *ShortSummary) GetPayout() int {
	return s.payout
}
