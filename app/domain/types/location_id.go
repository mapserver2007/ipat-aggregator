package types

type LocationId int

const (
	UnknownLocation LocationId = iota
	Miho
	Ritto
	LocalLocation
	OverseasLocation
)

var locationMap = map[LocationId]string{
	UnknownLocation:  "不明",
	Miho:             "美浦",
	Ritto:            "栗東",
	LocalLocation:    "地方",
	OverseasLocation: "海外",
}

func NewLocationId(name string) LocationId {
	for key, value := range locationMap {
		if value == name {
			return key
		}
	}
	return UnknownLocation
}

func (l LocationId) Value() int {
	return int(l)
}

func (l LocationId) Name() string {
	name, _ := locationMap[l]
	return name
}
