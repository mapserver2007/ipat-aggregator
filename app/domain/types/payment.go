package types

type Payment int

func (p Payment) Value() int {
	return int(p)
}
