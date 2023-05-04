package entity

import "strconv"

type PopularAnalyse struct {
	popularNumber int
	payment       int
	payout        int
	odds          string
	isHit         bool
}

func NewPopularAnalyse(
	popularNumber int,
	payment int,
	payout int,
	odds string,
	isHit bool,
) *PopularAnalyse {
	return &PopularAnalyse{
		popularNumber: popularNumber,
		payment:       payment,
		payout:        payout,
		odds:          odds,
		isHit:         isHit,
	}
}

func (p *PopularAnalyse) PopularNumber() int {
	return p.popularNumber
}

func (p *PopularAnalyse) Payment() int {
	return p.payment
}

func (p *PopularAnalyse) Payout() int {
	return p.payout
}

func (p *PopularAnalyse) Odds() float64 {
	odds, _ := strconv.ParseFloat(p.odds, 64)
	return odds
}

func (p *PopularAnalyse) IsHit() bool {
	return p.isHit
}
