package entity

import (
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"strconv"
)

// ShortSummary 投資・回収・回収率のみの集計結果
type ShortSummary struct {
	payment types.Payment
	payout  types.Payout
}

func NewShortSummary(
	payment types.Payment,
	payout types.Payout,
) ShortSummary {
	return ShortSummary{
		payment: payment,
		payout:  payout,
	}
}

func (s *ShortSummary) GetPayment() types.Payment {
	return s.payment
}

func (s *ShortSummary) GetPayout() types.Payout {
	return s.payout
}

func (s *ShortSummary) GetRecoveryRate() string {
	if s.payment == 0 {
		return fmt.Sprintf("%d%s", 0, "%")
	}
	return fmt.Sprintf("%s%s", strconv.FormatFloat((float64(s.payout)*float64(100))/float64(s.payment), 'f', 2, 64), "%")
}
