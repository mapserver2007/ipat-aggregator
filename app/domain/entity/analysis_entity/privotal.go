package analysis_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
)

type Pivotal struct {
	marker        types.Marker
	odds          decimal.Decimal
	horseNumber   int
	popularNumber int
	orderNo       int
}

func NewPivotal(
	marker types.Marker,
	odds string,
	horseNumber int,
	popularNumber int,
	orderNo int,
) *Pivotal {
	decimalOdds, _ := decimal.NewFromString(odds)
	return &Pivotal{
		marker:        marker,
		odds:          decimalOdds,
		horseNumber:   horseNumber,
		popularNumber: popularNumber,
		orderNo:       orderNo,
	}
}

func (p *Pivotal) Marker() types.Marker {
	return p.marker
}

func (p *Pivotal) Odds() decimal.Decimal {
	return p.odds
}

func (p *Pivotal) HorseNumber() int {
	return p.horseNumber
}

func (p *Pivotal) PopularNumber() int {
	return p.popularNumber

}
func (p *Pivotal) OrderNo() int {
	return p.orderNo
}
