package netkeiba_entity

type PayoutResult struct {
	ticketType int
	numbers    []string
	odds       []string
	populars   []int
}

func NewPayoutResult(
	ticketType int,
	numbers []string,
	odds []string,
	populars []int,
) *PayoutResult {
	return &PayoutResult{
		ticketType: ticketType,
		numbers:    numbers,
		odds:       odds,
		populars:   populars,
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

func (p *PayoutResult) Populars() []int {
	return p.populars
}
