package entity

import (
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/result/types"
	"math"
	"strconv"
)

type DetailSummary struct {
	betCount      types.BetCount
	hitCount      types.HitCount
	payment       types.Payment
	payout        types.Payout
	averagePayout types.Payout
	maxPayout     types.Payout
	minPayout     types.Payout
}

func NewDetailSummary(
	voteCount types.BetCount,
	hitCount types.HitCount,
	payment types.Payment,
	payout types.Payout,
	averagePayout types.Payout,
	maxPayout types.Payout,
	minPayout types.Payout,
) DetailSummary {
	return DetailSummary{
		betCount:      voteCount,
		hitCount:      hitCount,
		payment:       payment,
		payout:        payout,
		averagePayout: averagePayout,
		maxPayout:     maxPayout,
		minPayout:     minPayout,
	}
}

func (s *DetailSummary) GetBetCount() types.BetCount {
	return s.betCount
}

func (s *DetailSummary) GetHitCount() types.HitCount {
	return s.hitCount
}

func (s *DetailSummary) GetHitRate() string {
	if s.betCount == 0 {
		return fmt.Sprintf("%d%s", 0, "%")
	}
	hitRate := math.Round((float64(s.hitCount)/float64(s.betCount))*100) / 100
	return fmt.Sprintf("%s%s", strconv.FormatFloat(hitRate*100, 'f', 2, 64), "%")
}

func (s *DetailSummary) GetPayment() types.Payment {
	return s.payment
}

func (s *DetailSummary) GetPayout() types.Payout {
	return s.payout
}

func (s *DetailSummary) GetRecoveryRate() string {
	if s.payment == 0 {
		return fmt.Sprintf("%d%s", 0, "%")
	}
	return fmt.Sprintf("%s%s", strconv.FormatFloat((float64(s.payout)*float64(100))/float64(s.payment), 'f', 2, 64), "%")
}

func (s *DetailSummary) GetAveragePayout() types.Payout {
	return s.averagePayout
}
