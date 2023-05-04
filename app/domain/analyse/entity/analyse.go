package entity

type AnalyseSummary struct {
	winPopularitySummary *WinAnalyseSummary
}

type WinAnalyseSummary struct {
	allSummaries            []*PopularAnalyseSummary
	grade1Summaries         []*PopularAnalyseSummary
	grade2Summaries         []*PopularAnalyseSummary
	grade3Summaries         []*PopularAnalyseSummary
	allowanceClassSummaries []*PopularAnalyseSummary
}

func NewAnalyseSummary(
	winPopularitySummary *WinAnalyseSummary,
) *AnalyseSummary {
	return &AnalyseSummary{
		winPopularitySummary: winPopularitySummary,
	}
}

func (a *AnalyseSummary) WinPopularitySummary() *WinAnalyseSummary {
	return a.winPopularitySummary
}

func NewWinAnalyseSummary(
	allSummaries []*PopularAnalyseSummary,
	grade1Summaries []*PopularAnalyseSummary,
	grade2Summaries []*PopularAnalyseSummary,
	grade3Summaries []*PopularAnalyseSummary,
	allowanceClassSummaries []*PopularAnalyseSummary,
) *WinAnalyseSummary {
	return &WinAnalyseSummary{
		allSummaries:            allSummaries,
		grade1Summaries:         grade1Summaries,
		grade2Summaries:         grade2Summaries,
		grade3Summaries:         grade3Summaries,
		allowanceClassSummaries: allowanceClassSummaries,
	}
}

func (w *WinAnalyseSummary) AllSummaries() []*PopularAnalyseSummary {
	return w.allSummaries
}

func (w *WinAnalyseSummary) Grade1Summaries() []*PopularAnalyseSummary {
	return w.grade1Summaries
}

func (w *WinAnalyseSummary) Grade2Summaries() []*PopularAnalyseSummary {
	return w.grade2Summaries
}

func (w *WinAnalyseSummary) Grade3Summaries() []*PopularAnalyseSummary {
	return w.grade3Summaries
}

func (w *WinAnalyseSummary) AllowanceClassSummaries() []*PopularAnalyseSummary {
	return w.allowanceClassSummaries
}
