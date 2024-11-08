package types

type JockeyId string

func (j JockeyId) Value() string {
	return string(j)
}
