package entity

type PopularAnalyse struct {
	popularNumber       int
	betCount            int
	hitCount            int
	hitRate             float64
	averageOddsAtVote   float64
	averagePayoutAtVote float64
	averageOddsAtHit    float64
	averagePayoutAtHit  float64
}

func NewPopularAnalyse(
	popularNumber int,
	betCount int,
	hitCount int,
	hitRate float64,
	averageOddsAtVote float64,
	averagePayoutAtVote float64,
	averageOddsAtHit float64,
	averagePayoutAtHit float64,
) *PopularAnalyse {
	return &PopularAnalyse{
		popularNumber:       popularNumber,
		betCount:            betCount,
		hitCount:            hitCount,
		hitRate:             hitRate,
		averageOddsAtVote:   averageOddsAtVote,
		averagePayoutAtVote: averagePayoutAtVote,
		averageOddsAtHit:    averageOddsAtHit,
		averagePayoutAtHit:  averagePayoutAtHit,
	}
}

func (p *PopularAnalyse) PopularNumber() int {
	return p.popularNumber
}

func (p *PopularAnalyse) BetCount() int {
	return p.betCount
}

func (p *PopularAnalyse) HitCount() int {
	return p.hitCount
}

func (p *PopularAnalyse) HitRate() float64 {
	return p.hitRate
}

func (p *PopularAnalyse) AverageOddsAtVote() float64 {
	return p.averageOddsAtVote
}

func (p *PopularAnalyse) AveragePayoutAtVote() float64 {
	return p.averagePayoutAtVote
}

func (p *PopularAnalyse) AverageOddsAtHit() float64 {
	return p.averageOddsAtHit
}

func (p *PopularAnalyse) AveragePayoutAtHit() float64 {
	return p.averagePayoutAtHit
}
