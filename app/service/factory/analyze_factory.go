package factory

import analyze_entity "github.com/mapserver2007/ipat-aggregator/app/domain/analyze/entity"

var winOddsRangeSlice = []string{
	"1.1",
	"1.2",
	"1.3",
	"1.4",
	"1.5",
	"1.6",
	"1.7",
	"1.8",
	"1.9",
	"2.0-2.1",
	"2.2-2.3",
	"2.4-2.5",
	"2.6-2.7",
	"2.8-2.9",
	"3.0-3.4",
	"3.5-3.9",
	"4.0-4.9",
	"5.0-6.9",
	"7.0-9.9",
	"10-14.9",
	"15-19.9",
	"20-29.9",
	"30-49.9",
	"50-99.9",
	"100-",
}

func DefaultWinPopularAnalyzeSummarySlice() []*analyze_entity.WinPopularAnalyzeSummary {
	return make([]*analyze_entity.WinPopularAnalyzeSummary, 18)
}

func DefaultWinOddsAnalyzeSummaryMap() map[string]*analyze_entity.WinOddsAnalyzeSummary {
	winOddsAnalyzeSummaryMap := map[string]*analyze_entity.WinOddsAnalyzeSummary{}
	for _, oddsRange := range winOddsRangeSlice {
		winOddsAnalyzeSummaryMap[oddsRange] = analyze_entity.DefaultWinOddsAnalyzeSummary()
	}
	return winOddsAnalyzeSummaryMap
}