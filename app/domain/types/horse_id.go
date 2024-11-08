package types

type HorseId string

func (h HorseId) Value() string {
	return string(h)
}
