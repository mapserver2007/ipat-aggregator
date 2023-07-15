package value_object

import "strings"

type RepaymentComments []string

func (r RepaymentComments) String() string {
	return strings.Join(r, "\n")
}
