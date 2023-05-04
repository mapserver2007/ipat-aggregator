package entity

import (
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	"strconv"
)

type PopularAnalyse struct {
	popularNumber int
	payment       int
	payout        int
	odds          string
	isHit         bool
	class         race_vo.GradeClass
}

func NewPopularAnalyse(
	popularNumber int,
	payment int,
	payout int,
	odds string,
	isHit bool,
	class race_vo.GradeClass,
) *PopularAnalyse {
	return &PopularAnalyse{
		popularNumber: popularNumber,
		payment:       payment,
		payout:        payout,
		odds:          odds,
		isHit:         isHit,
		class:         class,
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

func (p *PopularAnalyse) Class() race_vo.GradeClass {
	return p.class
}
