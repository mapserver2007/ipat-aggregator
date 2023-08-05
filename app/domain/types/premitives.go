package types

type Payment int

type Payout int

type BetCount int

type HitCount int

func (p Payment) Value() int {
	return int(p)
}

func (p Payout) Value() int {
	return int(p)
}

func (b BetCount) Value() int {
	return int(b)
}

func (h HitCount) Value() int {
	return int(h)
}
