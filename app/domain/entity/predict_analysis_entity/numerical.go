package predict_analysis_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
)

type Numerical struct {
	payment   types.Payment
	payout    types.Payout
	voteCount int
	odds      decimal.Decimal
	number    types.BetNumber
	popular   int
}

func NewNumerical(
	payment types.Payment,
	payout types.Payout,
	odds string,
	number types.BetNumber,
	popular int,
) *Numerical {
	decimalOdds, _ := decimal.NewFromString(odds)
	return &Numerical{
		payment: payment,
		payout:  payout,
		odds:    decimalOdds,
		number:  number,
		popular: popular,
	}
}

func (n *Numerical) Payment() types.Payment {
	return n.payment
}

func (n *Numerical) Payout() types.Payout {
	return n.payout
}

func (n *Numerical) Odds() decimal.Decimal {
	return n.odds
}

func (n *Numerical) Number() types.BetNumber {
	return n.number
}

func (n *Numerical) Popular() int {
	return n.popular
}
