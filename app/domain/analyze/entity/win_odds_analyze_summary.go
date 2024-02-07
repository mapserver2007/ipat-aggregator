package entity

type WinOddsAnalyzeSummary struct {
	oddsRange          string
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

func DefaultWinOddsAnalyzeSummary() *WinOddsAnalyzeSummary {
	return &WinOddsAnalyzeSummary{}
}

func NewWinOddsAnalyzeSummary(
	oddsRange string,
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
) *WinOddsAnalyzeSummary {
	return &WinOddsAnalyzeSummary{
		oddsRange:          oddsRange,
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

func (s *WinOddsAnalyzeSummary) GetOddsRange() string {
	return s.oddsRange
}
