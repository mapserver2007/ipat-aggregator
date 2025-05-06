package types

type RaceSeason int

const (
	UnkownSeason RaceSeason = iota
	Spring
	Summer
	Autumn
	Winter
)

var raceSeasonMap = map[RaceSeason]string{
	UnkownSeason: "未定義",
	Spring:       "春",
	Summer:       "夏",
	Autumn:       "秋",
	Winter:       "冬",
}

func (r RaceSeason) Value() int {
	return int(r)
}

func (r RaceSeason) String() string {
	return raceSeasonMap[r]
}
