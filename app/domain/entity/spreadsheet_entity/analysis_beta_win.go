package spreadsheet_entity

import (
	"fmt"
	"math"

	"github.com/shopspring/decimal"
)

type AnalysisBetaRate struct {
	hitRate    float64
	payoutRate float64
}

func NewAnalysisBetaRate(
	odds []decimal.Decimal,
	raceCount int,
) *AnalysisBetaRate {
	hitRate := float64(len(odds)) * 100 / float64(raceCount)
	var payouts decimal.Decimal
	for _, o := range odds {
		payouts = payouts.Add(o)
	}
	payoutRate := payouts.InexactFloat64() * 100 / float64(raceCount)

	return &AnalysisBetaRate{
		hitRate:    hitRate,
		payoutRate: payoutRate,
	}
}

func (a *AnalysisBetaRate) HitRate() float64 {
	return a.hitRate
}

func (a *AnalysisBetaRate) HitRateFormat() string {
	return a.rateFormat(a.payoutRate)
}

func (a *AnalysisBetaRate) PayoutRate() float64 {
	return a.payoutRate
}

func (a *AnalysisBetaRate) PayoutRateFormat() string {
	return a.rateFormat(a.payoutRate)
}

func (a *AnalysisBetaRate) rateFormat(rate float64) string {
	if math.IsNaN(rate) {
		return "-"
	}
	return fmt.Sprintf("%.2f%%", rate)
}
