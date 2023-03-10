package entity

func NewMonthlySummary(monthlyRates map[int]ResultRate) MonthlySummary {
	return MonthlySummary{
		MonthlyRates: monthlyRates,
	}
}
