package analysis_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"github.com/shopspring/decimal"
)

type Calculable struct {
	ticketType types.TicketType
	odds       decimal.Decimal
	number     types.BetNumber
	popular    int
	orderNo    int
	entries    int
	filters    []filter.Id
	isHit      bool
	pivotals   []*Pivotal
}

func NewCalculable(
	ticketType types.TicketType,
	odds string,
	number types.BetNumber,
	popular int,
	orderNo int,
	entries int,
	filters []filter.Id,
	isHit bool,
	pivotals []*Pivotal,
) *Calculable {
	decimalOdds, _ := decimal.NewFromString(odds)
	return &Calculable{
		ticketType: ticketType,
		odds:       decimalOdds,
		number:     number,
		popular:    popular,
		orderNo:    orderNo,
		entries:    entries,
		filters:    filters,
		isHit:      isHit,
		pivotals:   pivotals,
	}
}

func (n *Calculable) TicketType() types.TicketType {
	return n.ticketType
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

func (n *Calculable) Entries() int {
	return n.entries
}

func (n *Calculable) Filters() []filter.Id {
	return n.filters
}

func (n *Calculable) IsHit() bool {
	return n.isHit
}

func (n *Calculable) Pivotals() []*Pivotal {
	return n.pivotals
}
