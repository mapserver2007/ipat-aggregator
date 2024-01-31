package analysis_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"github.com/shopspring/decimal"
)

type Calculable struct {
	payment   types.Payment
	payout    types.Payout
	voteCount int
	odds      decimal.Decimal
	number    types.BetNumber
	popular   int
	orderNo   int
	filters   []filter.Id
}

func NewCalculable(
	payment types.Payment,
	payout types.Payout,
	odds string,
	number types.BetNumber,
	popular int,
	orderNo int,
	filters []filter.Id,
) *Calculable {
	decimalOdds, _ := decimal.NewFromString(odds)
	return &Calculable{
		payment: payment,
		payout:  payout,
		odds:    decimalOdds,
		number:  number,
		popular: popular,
		orderNo: orderNo,
		filters: filters,
	}
}

func (n *Calculable) Payment() types.Payment {
	return n.payment
}

func (n *Calculable) Payout() types.Payout {
	return n.payout
}

func (n *Calculable) Odds() decimal.Decimal {
	return n.odds
}

func (n *Calculable) Number() types.BetNumber {
	return n.number
}

func (n *Calculable) Popular() int {
	return n.popular
}

func (n *Calculable) OrderNo() int {
	return n.orderNo
}

func (n *Calculable) Filters() []filter.Id {
	return n.filters
}
