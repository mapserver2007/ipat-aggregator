package types

type BetCount int

func (b BetCount) Value() int {
	return int(b)
}
