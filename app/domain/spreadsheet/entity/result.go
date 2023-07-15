package entity

func NewResult(
	bettingTicketSummary BettingTicketSummary,
	raceClassSummary RaceClassSummary,
	monthlySummary MonthlySummary,
	yearlySummary YearlySummary,
	courseCategorySummary CourseCategorySummary,
	distanceCategorySummary DistanceCategorySummary,
	raceCourseSummary RaceCourseSummary,
) *Summary {
	return &Summary{
		BettingTicketSummary:    bettingTicketSummary,
		RaceClassSummary:        raceClassSummary,
		MonthlySummary:          monthlySummary,
		YearlySummary:           yearlySummary,
		CourseCategorySummary:   courseCategorySummary,
		DistanceCategorySummary: distanceCategorySummary,
		RaceCourseSummary:       raceCourseSummary,
	}
}
