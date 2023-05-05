package entity

import (
	"fmt"
	"math"
	"strconv"
)

// https://db-keiba.com/popularity/
var generalWinRateSlice = []float64{
	0.0,
	0.331,
	0.191,
	0.138,
	0.092,
	0.071,
	0.050,
	0.039,
	0.028,
	0.020,
	0.018,
	0.011,
	0.010,
	0.008,
	0.006,
	0.003,
	0.003,
	0.001,
	0.001,
}

var generalPayoutRate = []float64{
	0.0,
	0.80,
	0.80,
	0.83,
	0.77,
	0.83,
	0.80,
	0.84,
	0.80,
	0.76,
	0.88,
	0.65,
	0.73,
	0.70,
	0.68,
	0.37,
	0.40,
	0.14,
	0.14,
}

type WinPopularAnalyzeSummary struct {
	popularNumber      int
	betCount           int
	hitCount           int
	hitRate            float64
	averageOddsAtVote  float64
	averageOddsAtHit   float64
	averageOddsAtUnHit float64
	totalPayment       int
	totalPayout        int
	averagePayment     int
	averagePayout      int
	medianPayment      int
	medianPayout       int
	maxPayout          int
	minPayout          int
	maxOddsAtHit       float64
	minOddsAtHit       float64
}

func DefaultWinPopularAnalyzeSummary(popularNumber int) *WinPopularAnalyzeSummary {
	return &WinPopularAnalyzeSummary{
		popularNumber: popularNumber,
	}
}

func NewWinPopularAnalyzeSummary(
	popularNumber int,
	betCount int,
	hitCount int,
	hitRate float64,
	averageOddsAtVote float64,
	averageOddsAtHit float64,
	averageOddsAtUnHit float64,
	totalPayment int,
	totalPayout int,
	averagePayment int,
	averagePayout int,
	medanPayment int,
	medianPayout int,
	maxPayout int,
	minPayout int,
	maxOddsAtHit float64,
	minOddsAtHit float64,
) *WinPopularAnalyzeSummary {
	return &WinPopularAnalyzeSummary{
		popularNumber:      popularNumber,
		betCount:           betCount,
		hitCount:           hitCount,
		hitRate:            hitRate,
		averageOddsAtVote:  averageOddsAtVote,
		averageOddsAtHit:   averageOddsAtHit,
		averageOddsAtUnHit: averageOddsAtUnHit,
		totalPayment:       totalPayment,
		totalPayout:        totalPayout,
		averagePayment:     averagePayment,
		averagePayout:      averagePayout,
		medianPayment:      medanPayment,
		medianPayout:       medianPayout,
		maxPayout:          maxPayout,
		minPayout:          minPayout,
		maxOddsAtHit:       maxOddsAtHit,
		minOddsAtHit:       minOddsAtHit,
	}
}

func (p *WinPopularAnalyzeSummary) PopularNumber() int {
	return p.popularNumber
}

func (p *WinPopularAnalyzeSummary) BetCount() int {
	return p.betCount
}

func (p *WinPopularAnalyzeSummary) HitCount() int {
	return p.hitCount
}

func (p *WinPopularAnalyzeSummary) HitRate() float64 {
	return p.hitRate
}

func (p *WinPopularAnalyzeSummary) FormattedHitRate() string {
	return fmt.Sprintf("%s%s", strconv.FormatFloat(p.hitRate*100, 'f', 2, 64), "%")
}

func (p *WinPopularAnalyzeSummary) AverageOddsAtVote() float64 {
	return p.averageOddsAtVote
}

func (p *WinPopularAnalyzeSummary) AverageOddsAtHit() float64 {
	return p.averageOddsAtHit
}

func (p *WinPopularAnalyzeSummary) AverageOddsAtUnHit() float64 {
	return p.averageOddsAtUnHit
}

func (p *WinPopularAnalyzeSummary) TotalPayment() int {
	return p.totalPayment
}

func (p *WinPopularAnalyzeSummary) TotalPayout() int {
	return p.totalPayout
}

func (p *WinPopularAnalyzeSummary) AveragePayment() int {
	return p.averagePayment
}

func (p *WinPopularAnalyzeSummary) AveragePayout() int {
	return p.averagePayout
}

func (p *WinPopularAnalyzeSummary) MedianPayment() int {
	return p.medianPayment
}

func (p *WinPopularAnalyzeSummary) MedianPayout() int {
	return p.medianPayout
}

func (p *WinPopularAnalyzeSummary) PayoutRate() float64 {
	if p.totalPayment == 0 {
		return 0
	}
	return math.Round((float64(p.totalPayout)/float64(p.totalPayment))*100) / 100
}

func (p *WinPopularAnalyzeSummary) FormattedPayoutRate() string {
	return fmt.Sprintf("%s%s", strconv.FormatFloat(p.PayoutRate()*100, 'f', 2, 64), "%")
}

func (p *WinPopularAnalyzeSummary) AveragePayoutRateAtHit() float64 {
	if p.averagePayment == 0 {
		return 0
	}
	return math.Round((float64(p.averagePayout)/float64(p.averagePayment))*100) / 100
}

func (p *WinPopularAnalyzeSummary) FormattedAveragePayoutRateAtHit() string {
	return fmt.Sprintf("%s%s", strconv.FormatFloat(p.AveragePayoutRateAtHit()*100, 'f', 2, 64), "%")
}

func (p *WinPopularAnalyzeSummary) MedianPayoutRateAtHit() float64 {
	if p.medianPayment == 0 {
		return 0
	}
	return math.Round((float64(p.medianPayout)/float64(p.medianPayment))*100) / 100
}

func (p *WinPopularAnalyzeSummary) FormattedMedianPayoutRateAtHit() string {
	return fmt.Sprintf("%s%s", strconv.FormatFloat(p.MedianPayoutRateAtHit()*100, 'f', 2, 64), "%")
}

func (p *WinPopularAnalyzeSummary) GeneralPayoutRate() float64 {
	if p.betCount == 0 {
		return 0
	}
	return generalPayoutRate[p.popularNumber]
}

func (p *WinPopularAnalyzeSummary) PayoutUpsideRate() float64 {
	return math.Round((p.PayoutRate()-p.GeneralPayoutRate())*100) / 100
}

func (p *WinPopularAnalyzeSummary) FormattedPayoutUpsideRate() string {
	return fmt.Sprintf("%s%s", strconv.FormatFloat(math.Round((p.PayoutRate()-p.GeneralPayoutRate())*100), 'f', 2, 64), "%")
}

func (p *WinPopularAnalyzeSummary) MaxPayout() int {
	return p.maxPayout
}

func (p *WinPopularAnalyzeSummary) MinPayout() int {
	return p.minPayout
}

func (p *WinPopularAnalyzeSummary) MaxOddsAtHit() float64 {
	return p.maxOddsAtHit
}

func (p *WinPopularAnalyzeSummary) MinOddsAtHit() float64 {
	return p.minOddsAtHit
}
