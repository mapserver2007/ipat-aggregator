package entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/result/types"
)

// ShortSummary 投資・回収・回収率のみの集計結果
type ShortSummary struct {
	payment      types.Payment
	payout       types.Payout
	recoveryRate string
}

func NewShortSummary(
	payment types.Payment,
	payout types.Payout,
	recoveryRate string,
) ShortSummary {
	return ShortSummary{
		payment:      payment,
		payout:       payout,
		recoveryRate: recoveryRate,
	}
}

func (s *ShortSummary) GetPayment() types.Payment {
	return s.payment
}

func (s *ShortSummary) GetPayout() types.Payout {
	return s.payout
}

func (s *ShortSummary) GetRecoveryRate() string {
	return s.recoveryRate
}
