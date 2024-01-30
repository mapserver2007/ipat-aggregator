package spreadsheet_entity

import (
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
	"strconv"
)

type AnalysisData struct {
	hitMarkerCombinationAnalysisMap   map[types.MarkerCombinationId]*MarkerCombinationAnalysis
	unHitMarkerCombinationAnalysisMap map[types.MarkerCombinationId]*MarkerCombinationAnalysis
	raceCount                         int
	ticketAnalysis                    *TicketAnalysis
	allMarkerCombinationIds           []types.MarkerCombinationId
}

func NewAnalysisData(
	hitMarkerCombinationAnalysisMap map[types.MarkerCombinationId]*MarkerCombinationAnalysis,
	unHitMarkerCombinationAnalysisMap map[types.MarkerCombinationId]*MarkerCombinationAnalysis,
	raceCount int,
	ticketAnalysis *TicketAnalysis,
	allMarkerCombinationIds []types.MarkerCombinationId,
) *AnalysisData {
	return &AnalysisData{
		hitMarkerCombinationAnalysisMap:   hitMarkerCombinationAnalysisMap,
		unHitMarkerCombinationAnalysisMap: unHitMarkerCombinationAnalysisMap,
		raceCount:                         raceCount,
		ticketAnalysis:                    ticketAnalysis,
		allMarkerCombinationIds:           allMarkerCombinationIds,
	}
}

func (a *AnalysisData) HitMarkerCombinationAnalysisMap() map[types.MarkerCombinationId]*MarkerCombinationAnalysis {
	return a.hitMarkerCombinationAnalysisMap
}

func (a *AnalysisData) UnHitMarkerCombinationAnalysisMap() map[types.MarkerCombinationId]*MarkerCombinationAnalysis {
	return a.unHitMarkerCombinationAnalysisMap
}

func (a *AnalysisData) RaceCount() int {
	return a.raceCount
}

func (a *AnalysisData) TicketAnalysis() *TicketAnalysis {
	return a.ticketAnalysis
}

func (a *AnalysisData) AllMarkerCombinationIds() []types.MarkerCombinationId {
	return a.allMarkerCombinationIds
}

type MarkerCombinationAnalysis struct {
	raceCount int // 予想したレース数
	populars  []int
	odds      []decimal.Decimal
}

func NewMarkerCombinationAnalysis(
	raceCount int,
) *MarkerCombinationAnalysis {
	return &MarkerCombinationAnalysis{
		raceCount: raceCount,
		populars:  make([]int, 0),
		odds:      make([]decimal.Decimal, 0),
	}
}

func (m *MarkerCombinationAnalysis) MatchRate() float64 {
	return (float64(len(m.odds)) * float64(100)) / float64(m.raceCount)
}

func (m *MarkerCombinationAnalysis) MatchRateFormat() string {
	return rateFormat(m.MatchRate())
}

func (m *MarkerCombinationAnalysis) MatchCount() int {
	return len(m.odds)
}

func (m *MarkerCombinationAnalysis) AddPopular(popular int) {
	m.populars = append(m.populars, popular)
}

func (m *MarkerCombinationAnalysis) AveragePopular() float64 {
	if len(m.populars) == 0 {
		return 0
	}
	return float64(sum(m.populars)) / float64(len(m.populars))
}

func (m *MarkerCombinationAnalysis) AveragePopularFormat() string {
	return rateFormat(m.AveragePopular())
}

func (m *MarkerCombinationAnalysis) MedianPopular() int {
	if len(m.populars) == 0 {
		return 0
	}
	return m.populars[len(m.populars)/2]
}

func (m *MarkerCombinationAnalysis) MaxPopular() int {
	maxPopular := 0
	for _, popular := range m.populars {
		if maxPopular < popular {
			maxPopular = popular
		}
	}
	return maxPopular
}

func (m *MarkerCombinationAnalysis) MinPopular() int {
	minPopular := 0
	for _, popular := range m.populars {
		if minPopular > popular {
			minPopular = popular
		}
	}
	return minPopular
}

func (m *MarkerCombinationAnalysis) Odds() []decimal.Decimal {
	return m.odds
}

func (m *MarkerCombinationAnalysis) AddOdds(odds decimal.Decimal) {
	m.odds = append(m.odds, odds)
}

func (m *MarkerCombinationAnalysis) AverageOdds() float64 {
	if len(m.odds) == 0 {
		return 0
	}
	return decimal.Sum(decimal.Zero, m.odds...).InexactFloat64() / float64(len(m.odds))
}

func (m *MarkerCombinationAnalysis) AverageOddsFormat() string {
	return rateFormat(m.AverageOdds())
}

func (m *MarkerCombinationAnalysis) MedianOdds() decimal.Decimal {
	if len(m.odds) == 0 {
		return decimal.Zero
	}
	return m.odds[len(m.odds)/2]
}

func (m *MarkerCombinationAnalysis) MaxOdds() decimal.Decimal {
	maxOdds := decimal.Zero
	for _, odds := range m.odds {
		if maxOdds.LessThan(odds) {
			maxOdds = odds
		}
	}
	return maxOdds
}

func (m *MarkerCombinationAnalysis) MinOdds() decimal.Decimal {
	minOdds := decimal.Zero
	for _, odds := range m.odds {
		if minOdds.GreaterThan(odds) {
			minOdds = odds
		}
	}
	return minOdds
}

type TicketAnalysis struct {
	raceCount int // 投票したレース数
	payments  []int
	payouts   []int
}

func NewTicketAnalysis(
	raceCount int,
) *TicketAnalysis {
	return &TicketAnalysis{
		raceCount: raceCount,
		payments:  make([]int, 0),
		payouts:   make([]int, 0),
	}
}

func (t *TicketAnalysis) AddPayment(payment int) {
	t.payments = append(t.payments, payment)
}

func (t *TicketAnalysis) AddPayout(payout int) {
	t.payouts = append(t.payouts, payout)
}

func (t *TicketAnalysis) HitCount() int {
	return len(t.payouts)
}

func (t *TicketAnalysis) HitRate() float64 {
	if len(t.payments) == 0 {
		return 0
	}
	return (float64(len(t.payouts)) * float64(100)) / float64(len(t.payments))
}

func (t *TicketAnalysis) HitRateFormat() string {
	return rateFormat(t.HitRate())
}

func (t *TicketAnalysis) VoteCount() int {
	return len(t.payments)
}

func (t *TicketAnalysis) RaceCount() int {
	return t.raceCount
}

func (t *TicketAnalysis) PayoutRate() float64 {
	totalPayment := sum(t.payments)
	totalPayout := sum(t.payouts)
	if totalPayment == 0 {
		return 0
	}
	return (float64(totalPayout) * float64(100)) / float64(totalPayment)
}

func (t *TicketAnalysis) PayoutRateFormat() string {
	return rateFormat(t.PayoutRate())
}

func (t *TicketAnalysis) AveragePayoutRate() float64 {
	totalPayout := sum(t.payouts)
	if len(t.payments) == 0 {
		return 0
	}
	return (float64(totalPayout) * float64(100)) / float64(len(t.payments))
}

func (t *TicketAnalysis) AveragePayoutRateFormat() string {
	return rateFormat(t.AveragePayoutRate())
}

func (t *TicketAnalysis) MedianPayout() int {
	if len(t.payouts) == 0 {
		return 0
	}
	return t.payouts[len(t.payouts)/2]
}

func (t *TicketAnalysis) MaxPayout() int {
	maxPayout := 0
	for _, payout := range t.payouts {
		if maxPayout < payout {
			maxPayout = payout
		}
	}
	return maxPayout
}

func (t *TicketAnalysis) MinPayout() int {
	minPayout := 0
	for _, payout := range t.payouts {
		if minPayout > payout {
			minPayout = payout
		}
	}
	return minPayout
}

func rateFormat(rate float64) string {
	return fmt.Sprintf("%s%s", strconv.FormatFloat(rate, 'f', 2, 64), "%")
}

func sum(nums []int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
}
