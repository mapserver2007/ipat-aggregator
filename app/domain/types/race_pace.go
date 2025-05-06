package types

type RacePace int

const (
	UnknownRacePace RacePace = iota
	High
	Middle
	Slow
)

var racePaceMap = map[RacePace]string{
	UnknownRacePace: "?",
	High:            "H",
	Middle:          "M",
	Slow:            "S",
}

func NewRacePace(s string) RacePace {
	for k, v := range racePaceMap {
		if v == s {
			return k
		}
	}
	return UnknownRacePace
}

func (r RacePace) Value() int {
	return int(r)
}

func (r RacePace) String() string {
	return racePaceMap[r]
}
