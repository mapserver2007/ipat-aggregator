package types

type HorseBirthDay int

func (h HorseBirthDay) Value() int {
	return int(h)
}
