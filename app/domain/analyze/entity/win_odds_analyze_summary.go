package entity

//var oddsRangeSlice = []string{
//	"1.7",
//	"1.8",
//	"1.9",
//	"2.0-2.1",
//	"2.2-2.3",
//	"2.4-2.5",
//	"2.6-2.7",
//	"2.8-2.9",
//	"3.0-3.4",
//	"3.5-3.9",
//	"4.0-4.9",
//	"5.0-6.9",
//	"7.0-9.9",
//	"10-14.9",
//	"15-19.9",
//	"20-29.9",
//	"30-49.9",
//	"50-99.9",
//	"100-",
//}
//
//func OddsRangeSlice() []string {
//	return oddsRangeSlice
//}

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
