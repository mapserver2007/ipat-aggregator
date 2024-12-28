package types

type TrainerId string

func (t TrainerId) Value() string {
	return string(t)
}
