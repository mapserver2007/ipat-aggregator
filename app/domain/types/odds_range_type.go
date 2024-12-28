package types

type OddsRangeType int

const (
	UnknownOddsRangeType OddsRangeType = iota
	WinOddsRange1
	WinOddsRange2
	WinOddsRange3
	WinOddsRange4
	WinOddsRange5
	WinOddsRange6
	WinOddsRange7
	WinOddsRange8
	WinOddsRange9
	TrioOddsRange1
	TrioOddsRange2
	TrioOddsRange3
	TrioOddsRange4
	TrioOddsRange5
	TrioOddsRange6
	TrioOddsRange7
	TrioOddsRange8
)

var oddsRangeMap = map[OddsRangeType]string{
	WinOddsRange1:  "1.0-1.4",
	WinOddsRange2:  "1.5-1.9",
	WinOddsRange3:  "2.0-2.2",
	WinOddsRange4:  "2.3-3.0",
	WinOddsRange5:  "3.1-4.9",
	WinOddsRange6:  "5.0-9.9",
	WinOddsRange7:  "10.0-19.9",
	WinOddsRange8:  "20.0-49.9",
	WinOddsRange9:  "50.0-",
	TrioOddsRange1: "1.0-9.9",
	TrioOddsRange2: "10.0-19.9",
	TrioOddsRange3: "20.0-29.9",
	TrioOddsRange4: "30.0-49.9",
	TrioOddsRange5: "50.0-99.9",
	TrioOddsRange6: "100-299",
	TrioOddsRange7: "300-499",
	TrioOddsRange8: "500-",
}

func (m OddsRangeType) Value() int {
	return int(m)
}

func (m OddsRangeType) String() string {
	oddsRange, _ := oddsRangeMap[m]
	return oddsRange
}
