package types

type BreederId string

func (b BreederId) Value() string {
	return string(b)
}
