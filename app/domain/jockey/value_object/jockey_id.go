package value_object

import "fmt"

type JockeyId int

func (j JockeyId) Format() string {
	return fmt.Sprintf("%05d", j)
}

func (j JockeyId) Value() int {
	return int(j)
}