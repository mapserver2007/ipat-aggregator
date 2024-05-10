package spreadsheet_entity

type TicketSummary struct {
	winTermResult           *TicketResult
	placeTermResult         *TicketResult
	quinellaTermResult      *TicketResult
	exactaTermResult        *TicketResult
	quinellaPlaceTermResult *TicketResult
	trioTermResult          *TicketResult
	trifectaTermResult      *TicketResult
}

func NewTicketSummary(
	winTermResult *TicketResult,
	placeTermResult *TicketResult,
	quinellaTermResult *TicketResult,
	exactaTermResult *TicketResult,
	quinellaPlaceTermResult *TicketResult,
	trioTermResult *TicketResult,
	trifectaTermResult *TicketResult,
) *TicketSummary {
	return &TicketSummary{
		winTermResult:           winTermResult,
		placeTermResult:         placeTermResult,
		quinellaTermResult:      quinellaTermResult,
		exactaTermResult:        exactaTermResult,
		quinellaPlaceTermResult: quinellaPlaceTermResult,
		trioTermResult:          trioTermResult,
		trifectaTermResult:      trifectaTermResult,
	}
}

func (t *TicketSummary) WinTermResult() *TicketResult {
	return t.winTermResult
}

func (t *TicketSummary) PlaceTermResult() *TicketResult {
	return t.placeTermResult
}

func (t *TicketSummary) QuinellaTermResult() *TicketResult {
	return t.quinellaTermResult
}

func (t *TicketSummary) ExactaTermResult() *TicketResult {
	return t.exactaTermResult
}

func (t *TicketSummary) QuinellaPlaceTermResult() *TicketResult {
	return t.quinellaPlaceTermResult
}

func (t *TicketSummary) TrioTermResult() *TicketResult {
	return t.trioTermResult
}

func (t *TicketSummary) TrifectaTermResult() *TicketResult {
	return t.trifectaTermResult
}
