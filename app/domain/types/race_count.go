package types

type RaceCount int

func (r RaceCount) Value() int {
	return int(r)
}
