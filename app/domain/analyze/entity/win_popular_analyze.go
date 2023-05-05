package entity

import (
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	"strconv"
)

type WinPopularAnalyze struct {
	popularNumber int
	payment       int
	payout        int
	odds          string
	isHit         bool
	class         race_vo.GradeClass
}

func NewWinPopularAnalyze(
	popularNumber int,
	payment int,
	payout int,
	odds string,
	isHit bool,
	class race_vo.GradeClass,
) *WinPopularAnalyze {
	return &WinPopularAnalyze{
		popularNumber: popularNumber,
		payment:       payment,
		payout:        payout,
		odds:          odds,
		isHit:         isHit,
		class:         class,
	}
}

func (p *WinPopularAnalyze) PopularNumber() int {
	return p.popularNumber
}

func (p *WinPopularAnalyze) Payment() int {
	return p.payment
}

func (p *WinPopularAnalyze) Payout() int {
	return p.payout
}

func (p *WinPopularAnalyze) Odds() float64 {
	odds, _ := strconv.ParseFloat(p.odds, 64)
	return odds
}

func (p *WinPopularAnalyze) IsHit() bool {
	return p.isHit
}

func (p *WinPopularAnalyze) Class() race_vo.GradeClass {
	return p.class
}
