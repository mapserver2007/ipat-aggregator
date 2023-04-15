package entity

func NewYearlySummary(yearlyRates map[int]ResultRate) YearlySummary {
	return YearlySummary{
		YearlyRates: yearlyRates,
	}
}
