package entity

func NewResultRate(
	voteCount,
	hitCount,
	payments,
	repayments int,
) ResultRate {
	return ResultRate{
		VoteCount:  voteCount,
		HitCount:   hitCount,
		Payments:   payments,
		Repayments: repayments,
	}
}
