package entity

import (
	analyze_vo "github.com/mapserver2007/ipat-aggregator/app/domain/analyze/value_object"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	"strconv"
)

type WinAnalyze struct {
	popularNumber int
	payment       int
	payout        int
	odds          string
	isHit         bool
	class         race_vo.GradeClass
}

func NewWinAnalyze(
	popularNumber int,
	payment int,
	payout int,
	odds string,
	isHit bool,
	class race_vo.GradeClass,
) *WinAnalyze {
	return &WinAnalyze{
		popularNumber: popularNumber,
		payment:       payment,
		payout:        payout,
		odds:          odds,
		isHit:         isHit,
		class:         class,
	}
}

func (p *WinAnalyze) PopularNumber() int {
	return p.popularNumber
}

func (p *WinAnalyze) Payment() int {
	return p.payment
}

func (p *WinAnalyze) Payout() int {
	return p.payout
}

func (p *WinAnalyze) Odds() analyze_vo.WinOdds {
	odds, _ := strconv.ParseFloat(p.odds, 64)
	return analyze_vo.WinOdds(odds)
}

func (p *WinAnalyze) IsHit() bool {
	return p.isHit
}

func (p *WinAnalyze) Class() race_vo.GradeClass {
	return p.class
}
