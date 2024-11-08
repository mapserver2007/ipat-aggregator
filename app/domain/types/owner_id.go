package types

type OwnerId string

func (o OwnerId) Value() string {
	return string(o)
}
