package value_object

import (
	"strconv"
	"time"
)

type RaceDate int

func NewRaceDate(s string) RaceDate {
	i, _ := strconv.Atoi(s)
	return RaceDate(i)
}

func (d RaceDate) Year() int {
	return toDate(d).Year()
}

func (d RaceDate) Month() int {
	return int(toDate(d).Month())
}

func (d RaceDate) Day() int {
	return toDate(d).Day()
}

func (d RaceDate) DateFormat() string {
	return toDate(d).Format("2006/01/02")
}

func toDate(d RaceDate) time.Time {
	date, err := time.Parse("20060102", strconv.Itoa(int(d)))
	if err != nil {
		panic(err)
	}
	return date
}
