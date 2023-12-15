package spreadsheet_entity

type Summary struct {
	allShortSummary   *ShortSummary
	monthShortSummary *ShortSummary
	yearShortSummary  *ShortSummary
}

func NewSummary(
	allShortSummary *ShortSummary,
	monthShortSummary *ShortSummary,
	yearShortSummary *ShortSummary,
) *Summary {
	return &Summary{
		allShortSummary:   allShortSummary,
		monthShortSummary: monthShortSummary,
		yearShortSummary:  yearShortSummary,
	}
}

func (s *Summary) GetAllShortSummary() *ShortSummary {
	return s.allShortSummary
}

func (s *Summary) GetMonthShortSummary() *ShortSummary {
	return s.monthShortSummary
}

func (s *Summary) GetYearShortSummary() *ShortSummary {
	return s.yearShortSummary
}
