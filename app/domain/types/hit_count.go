package types

type HitCount int

func (h HitCount) Value() int {
	return int(h)
}
