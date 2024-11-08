package types

type Payout int

func (p Payout) Value() int {
	return int(p)
}
