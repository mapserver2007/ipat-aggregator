package types

type TicketResult int

const (
	TicketNoBet TicketResult = iota
	TicketHit
	TicketUnHit
)
