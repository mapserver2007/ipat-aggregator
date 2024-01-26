package spreadsheet_entity

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
)

type MarkerAnalysis struct {
	hitCount  int
	voteCount int
	raceCount int
	payments  []int
	payouts   []int
	populars  []int
	odds      []decimal.Decimal
}

func NewMarkerAnalysis() *MarkerAnalysis {
	return &MarkerAnalysis{
		hitCount:  0,
		voteCount: 0,
		raceCount: 0,
		payments:  make([]int, 0),
		payouts:   make([]int, 0),
		populars:  make([]int, 0),
		odds:      make([]decimal.Decimal, 0),
	}
}

func (m *MarkerAnalysis) AddHitCount() {
	m.hitCount += 1
}

func (m *MarkerAnalysis) AddVoteCount() {
	m.voteCount += 1
}

func (m *MarkerAnalysis) AddRaceCount() {
	m.raceCount += 1
}

func (m *MarkerAnalysis) AddPayment(payment int) {
	m.payments = append(m.payments, payment)
}

func (m *MarkerAnalysis) AddPayout(payout int) {
	m.payouts = append(m.payouts, payout)
}

func (m *MarkerAnalysis) HitRate() float64 {
	if m.voteCount == 0 {
		return 0
	}
	return (float64(m.hitCount) * float64(100)) / float64(m.voteCount)
}

func (m *MarkerAnalysis) HitRateFormat() string {
	return fmt.Sprintf("%s%s", strconv.FormatFloat(m.HitRate(), 'f', 1, 64), "%")
}

func (m *MarkerAnalysis) HitCount() int {
	return m.hitCount
}

func (m *MarkerAnalysis) VoteCount() int {
	return m.voteCount
}

func (m *MarkerAnalysis) PayoutRate() float64 {
	totalPayment := sum(m.payments)
	totalPayout := sum(m.payouts)
	if totalPayment == 0 {
		return 0
	}
	return (float64(totalPayout) * float64(100)) / float64(totalPayment)
}

func (m *MarkerAnalysis) PayoutRateFormat() string {
	return rateFormat(m.PayoutRate())
}

func (m *MarkerAnalysis) AveragePayoutRate() float64 {
	totalPayout := sum(m.payouts)
	if m.voteCount == 0 {
		return 0
	}
	return (float64(totalPayout) * float64(100)) / float64(m.voteCount)
}

func (m *MarkerAnalysis) AveragePayoutRateFormat() string {
	return rateFormat(m.AveragePayoutRate())
}

func (m *MarkerAnalysis) MedianPayout() int {
	if len(m.payouts) == 0 {
		return 0
	}
	return m.payouts[len(m.payouts)/2]
}

func (m *MarkerAnalysis) MaxPayout() int {
	maxPayout := 0
	for _, payout := range m.payouts {
		if maxPayout < payout {
			maxPayout = payout
		}
	}
	return maxPayout
}

func (m *MarkerAnalysis) MinPayout() int {
	minPayout := 0
	for _, payout := range m.payouts {
		if minPayout > payout {
			minPayout = payout
		}
	}
	return minPayout
}

func (m *MarkerAnalysis) AddPopular(popular int) {
	m.populars = append(m.populars, popular)
}

func (m *MarkerAnalysis) AveragePopular() float64 {
	if m.voteCount == 0 {
		return 0
	}
	return float64(sum(m.populars)) / float64(m.voteCount)
}

func (m *MarkerAnalysis) AveragePopularFormat() string {
	return rateFormat(m.AveragePopular())
}

func (m *MarkerAnalysis) MaxPopular() int {
	maxPopular := 0
	for _, popular := range m.populars {
		if maxPopular < popular {
			maxPopular = popular
		}
	}
	return maxPopular
}

func (m *MarkerAnalysis) MinPopular() int {
	minPopular := 0
	for _, popular := range m.populars {
		if minPopular > popular {
			minPopular = popular
		}
	}
	return minPopular
}

func (m *MarkerAnalysis) AddOdds(odds decimal.Decimal) {
	m.odds = append(m.odds, odds)
}

func (m *MarkerAnalysis) AverageOdds() float64 {
	if m.voteCount == 0 {
		return 0
	}
	totalOdds := decimal.NewFromFloat(0)
	for _, odds := range m.odds {
		totalOdds = totalOdds.Add(odds)
	}
	return totalOdds.InexactFloat64() / float64(m.voteCount)
}

func (m *MarkerAnalysis) AverageOddsFormat() string {
	return rateFormat(m.AverageOdds())
}

func (m *MarkerAnalysis) MaxOdds() float64 {
	maxOdds := decimal.NewFromFloat(0)
	for _, odds := range m.odds {
		if maxOdds.LessThan(odds) {
			maxOdds = odds
		}
	}
	return maxOdds.InexactFloat64()
}

func (m *MarkerAnalysis) MinOdds() float64 {
	minOdds := decimal.NewFromFloat(0)
	for _, odds := range m.odds {
		if minOdds.GreaterThan(odds) {
			minOdds = odds
		}
	}
	return minOdds.InexactFloat64()
}

func rateFormat(rate float64) string {
	return fmt.Sprintf("%s%s", strconv.FormatFloat(rate, 'f', 1, 64), "%")
}

func sum(nums []int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
}
