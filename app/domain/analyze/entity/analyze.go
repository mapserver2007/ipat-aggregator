package entity

type AnalyzeSummary struct {
	winPopularitySummary *WinAnalyzeSummary
}

type WinAnalyzeSummary struct {
	allSummaries            []*WinPopularAnalyzeSummary
	grade1Summaries         []*WinPopularAnalyzeSummary
	grade2Summaries         []*WinPopularAnalyzeSummary
	grade3Summaries         []*WinPopularAnalyzeSummary
	allowanceClassSummaries []*WinPopularAnalyzeSummary
}

func NewAnalyzeSummary(
	winPopularitySummary *WinAnalyzeSummary,
) *AnalyzeSummary {
	return &AnalyzeSummary{
		winPopularitySummary: winPopularitySummary,
	}
}

func (a *AnalyzeSummary) WinPopularitySummary() *WinAnalyzeSummary {
	return a.winPopularitySummary
}

func NewWinAnalyzeSummary(
	allSummaries []*WinPopularAnalyzeSummary,
	grade1Summaries []*WinPopularAnalyzeSummary,
	grade2Summaries []*WinPopularAnalyzeSummary,
	grade3Summaries []*WinPopularAnalyzeSummary,
	allowanceClassSummaries []*WinPopularAnalyzeSummary,
) *WinAnalyzeSummary {
	return &WinAnalyzeSummary{
		allSummaries:            allSummaries,
		grade1Summaries:         grade1Summaries,
		grade2Summaries:         grade2Summaries,
		grade3Summaries:         grade3Summaries,
		allowanceClassSummaries: allowanceClassSummaries,
	}
}

func (w *WinAnalyzeSummary) AllSummaries() []*WinPopularAnalyzeSummary {
	return w.allSummaries
}

func (w *WinAnalyzeSummary) Grade1Summaries() []*WinPopularAnalyzeSummary {
	return w.grade1Summaries
}

func (w *WinAnalyzeSummary) Grade2Summaries() []*WinPopularAnalyzeSummary {
	return w.grade2Summaries
}

func (w *WinAnalyzeSummary) Grade3Summaries() []*WinPopularAnalyzeSummary {
	return w.grade3Summaries
}

func (w *WinAnalyzeSummary) AllowanceClassSummaries() []*WinPopularAnalyzeSummary {
	return w.allowanceClassSummaries
}
