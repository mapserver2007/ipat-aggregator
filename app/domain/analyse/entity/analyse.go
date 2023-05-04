package entity

type AnalyseSummary struct {
	popularAnalyseSummaries []*PopularAnalyseSummary
}

func NewAnalyseSummary(
	popularAnalyseSummaries []*PopularAnalyseSummary,
) *AnalyseSummary {
	return &AnalyseSummary{
		popularAnalyseSummaries: popularAnalyseSummaries,
	}
}

func (a *AnalyseSummary) PopularAnalyseSummaries() []*PopularAnalyseSummary {
	return a.popularAnalyseSummaries
}
