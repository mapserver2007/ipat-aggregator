package value_object

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	DefaultQuinellaSeparator = "―"
	QuinellaSeparator        = "-"
	ExactaSeparator          = "→"
)

type BetNumber string

func NewBetNumber(number string) BetNumber {
	number = strings.Replace(number, DefaultQuinellaSeparator, QuinellaSeparator, -1)
	return BetNumber(number)
}

func (b BetNumber) List() []int {
	separators := fmt.Sprintf("[%s,%s]", QuinellaSeparator, ExactaSeparator)
	list := regexp.MustCompile(separators).Split(string(b), -1)
	var betNumbers []int
	for _, s := range list {
		betNumber, _ := strconv.Atoi(s)
		betNumbers = append(betNumbers, betNumber)
	}

	return betNumbers
}

func (b BetNumber) String() string {
	// 三連複はダッシュなのでハイフンでつなぐ
	if strings.Contains(string(b), QuinellaSeparator) {
		return strings.Replace(string(b), QuinellaSeparator, "-", -1)
	}
	return string(b)
}
