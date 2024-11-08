package types

import "fmt"

type Marker int

const (
	Favorite      Marker = iota + 1 // ◎
	Rival                           // ◯
	BrackTriangle                   // ▲
	WhiteTriangle                   // △
	Star                            // ☆
	Check                           // ✓
	NoMarker      Marker = 9        // 無
	AnyMarker     Marker = 0        // 印(any)
)

var markerMap = map[Marker]string{
	Favorite:      "◎",
	Rival:         "◯",
	BrackTriangle: "▲",
	WhiteTriangle: "△",
	Star:          "☆",
	Check:         "✓",
	NoMarker:      "無",
	AnyMarker:     "印",
}

func NewMarker(value int) (Marker, error) {
	for mark := range markerMap {
		if int(mark) == value {
			return Marker(value), nil
		}
	}
	return 0, fmt.Errorf("invalid marker value: %d", value)
}

func (m Marker) Value() int {
	return int(m)
}

func (m Marker) String() string {
	marker, _ := markerMap[m]
	return marker
}
