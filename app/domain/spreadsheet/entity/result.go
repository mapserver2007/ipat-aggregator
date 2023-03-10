package entity

func NewResult(
	totalResultSummary ResultSummary,
	latestMonthResultSummary ResultSummary,
	bettingTicketSummary BettingTicketSummary,
	raceClassSummary RaceClassSummary,
	monthlySummary MonthlySummary,
	courseCategorySummary CourseCategorySummary,
	distanceCategorySummary DistanceCategorySummary,
	raceCourseSummary RaceCourseSummary,
) *Summary {
	return &Summary{
		TotalResultSummary:       totalResultSummary,
		LatestMonthResultSummary: latestMonthResultSummary,
		BettingTicketSummary:     bettingTicketSummary,
		RaceClassSummary:         raceClassSummary,
		MonthlySummary:           monthlySummary,
		CourseCategorySummary:    courseCategorySummary,
		DistanceCategorySummary:  distanceCategorySummary,
		RaceCourseSummary:        raceCourseSummary,
	}
}
