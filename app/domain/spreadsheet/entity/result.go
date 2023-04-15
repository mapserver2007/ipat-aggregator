package entity

func NewResult(
	totalResultSummary ResultSummary,
	latestMonthResultSummary ResultSummary,
	latestYearResultSummary ResultSummary,
	bettingTicketSummary BettingTicketSummary,
	raceClassSummary RaceClassSummary,
	monthlySummary MonthlySummary,
	yearlySummary YearlySummary,
	courseCategorySummary CourseCategorySummary,
	distanceCategorySummary DistanceCategorySummary,
	raceCourseSummary RaceCourseSummary,
) *Summary {
	return &Summary{
		TotalResultSummary:       totalResultSummary,
		LatestMonthResultSummary: latestMonthResultSummary,
		LatestYearResultSummary:  latestYearResultSummary,
		BettingTicketSummary:     bettingTicketSummary,
		RaceClassSummary:         raceClassSummary,
		MonthlySummary:           monthlySummary,
		YearlySummary:            yearlySummary,
		CourseCategorySummary:    courseCategorySummary,
		DistanceCategorySummary:  distanceCategorySummary,
		RaceCourseSummary:        raceCourseSummary,
	}
}
