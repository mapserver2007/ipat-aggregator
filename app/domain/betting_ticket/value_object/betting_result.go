package value_object

type BettingResult int

const (
	Unknown BettingResult = iota
	Hit
	UnHit
)
