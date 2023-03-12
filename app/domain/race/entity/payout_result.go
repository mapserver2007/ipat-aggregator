package entity

type PayoutResult struct {
	ticketType int
	numbers    []string
	odds       []string
}

func NewPayoutResult(
	ticketType int,
	numbers []string,
	odds []string,
) *PayoutResult {
	return &PayoutResult{
		ticketType: ticketType,
		numbers:    numbers,
		odds:       odds,
	}
}

func (p *PayoutResult) TicketType() int {
	return p.ticketType
}

func (p *PayoutResult) Numbers() []string {
	return p.numbers
}

func (p *PayoutResult) Odds() []string {
	return p.odds
}
