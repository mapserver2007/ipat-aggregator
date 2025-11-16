package types

type JockeyResultId string

func (j JockeyResultId) Value() string {
	return string(j)
}
