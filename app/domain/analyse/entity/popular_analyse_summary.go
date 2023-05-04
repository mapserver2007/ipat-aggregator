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

type PopularAnalyseSummary struct {
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
	medianPayment      float64
	medianPayout       float64
	maxPayout          int
	minPayout          int
	maxOddsAtHit       float64
	minOddsAtHit       float64
}

func DefaultPopularAnalyseSummary(popularNumber int) *PopularAnalyseSummary {
	return &PopularAnalyseSummary{
		popularNumber: popularNumber,
	}
}

func NewPopularAnalyseSummary(
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
	medanPayment float64,
	medianPayout float64,
	maxPayout int,
	minPayout int,
	maxOddsAtHit float64,
	minOddsAtHit float64,
) *PopularAnalyseSummary {
	return &PopularAnalyseSummary{
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

func (p *PopularAnalyseSummary) PopularNumber() int {
	return p.popularNumber
}

func (p *PopularAnalyseSummary) BetCount() int {
	return p.betCount
}

func (p *PopularAnalyseSummary) HitCount() int {
	return p.hitCount
}

func (p *PopularAnalyseSummary) HitRate() float64 {
	return p.hitRate
}

func (p *PopularAnalyseSummary) FormattedHitRate() string {
	return fmt.Sprintf("%s%s", strconv.FormatFloat(p.hitRate*100, 'f', 2, 64), "%")
}

func (p *PopularAnalyseSummary) AverageOddsAtVote() float64 {
	return p.averageOddsAtVote
}

func (p *PopularAnalyseSummary) AverageOddsAtHit() float64 {
	return p.averageOddsAtHit
}

func (p *PopularAnalyseSummary) AverageOddsAtUnHit() float64 {
	return p.averageOddsAtUnHit
}

func (p *PopularAnalyseSummary) TotalPayment() int {
	return p.totalPayment
}

func (p *PopularAnalyseSummary) TotalPayout() int {
	return p.totalPayout
}

func (p *PopularAnalyseSummary) AveragePayment() int {
	return p.averagePayment
}

func (p *PopularAnalyseSummary) AveragePayout() int {
	return p.averagePayout
}

func (p *PopularAnalyseSummary) MedianPayment() float64 {
	return p.medianPayment
}

func (p *PopularAnalyseSummary) MedianPayout() float64 {
	return p.medianPayout
}

func (p *PopularAnalyseSummary) PayoutRate() float64 {
	if p.totalPayment == 0 {
		return 0
	}
	return math.Round((float64(p.totalPayout)/float64(p.totalPayment))*100) / 100
}

func (p *PopularAnalyseSummary) FormattedPayoutRate() string {
	return fmt.Sprintf("%s%s", strconv.FormatFloat(p.PayoutRate()*100, 'f', 2, 64), "%")
}

func (p *PopularAnalyseSummary) GeneralWinRate() float64 {
	if p.betCount == 0 {
		return 0
	}
	return generalWinRateSlice[p.popularNumber]
}

func (p *PopularAnalyseSummary) GeneralPayoutRate() float64 {
	if p.betCount == 0 {
		return 0
	}
	return generalPayoutRate[p.popularNumber]
}

func (p *PopularAnalyseSummary) PayoutUpsideRate() float64 {
	return math.Round((p.PayoutRate()-p.GeneralPayoutRate())*100) / 100
}

func (p *PopularAnalyseSummary) FormattedPayoutUpsideRate() string {
	return fmt.Sprintf("%s%s", strconv.FormatFloat(math.Round((p.PayoutRate()-p.GeneralPayoutRate())*100), 'f', 2, 64), "%")
}

func (p *PopularAnalyseSummary) MaxPayout() int {
	return p.maxPayout
}

func (p *PopularAnalyseSummary) MinPayout() int {
	return p.minPayout
}

func (p *PopularAnalyseSummary) MaxOddsAtHit() float64 {
	return p.maxOddsAtHit
}

func (p *PopularAnalyseSummary) MinOddsAtHit() float64 {
	return p.minOddsAtHit
}
