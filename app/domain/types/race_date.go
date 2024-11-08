package types

import (
	"strconv"
	"time"
)

type RaceDate int

func NewRaceDate(s string) (RaceDate, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	raceDate := RaceDate(i)
	return raceDate, nil
}

func (d RaceDate) Value() int {
	return int(d)
}

func (d RaceDate) Year() int {
	return d.Date().Year()
}

func (d RaceDate) Month() int {
	return int(d.Date().Month())
}

func (d RaceDate) Day() int {
	return d.Date().Day()
}

func (d RaceDate) Format(layout string) string {
	return d.Date().Format(layout)
}

func (d RaceDate) Date() time.Time {
	date, err := time.Parse("20060102", strconv.Itoa(int(d)))
	if err != nil {
		panic(err)
	}
	return date
}
