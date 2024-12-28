package types

type HorseNumber int

func (h HorseNumber) Value() int {
	return int(h)
}
