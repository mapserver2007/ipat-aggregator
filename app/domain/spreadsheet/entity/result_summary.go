package entity

func NewResultSummary(
	payments,
	repayments int,
) ResultSummary {
	return ResultSummary{
		Payments:   payments,
		Repayments: repayments,
	}
}
