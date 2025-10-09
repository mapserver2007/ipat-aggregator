package types

import "fmt"

type RaceId string

func (r RaceId) String() string {
	return string(r)
}

func NewRaceIdForJRA(
	year int,
	day int,
	raceCourse RaceCourse,
	raceRound int,
	raceNo int,
) RaceId {
	rawRaceId := fmt.Sprintf("%d%s%02d%02d%02d", year, raceCourse.Value(), raceRound, day, raceNo)
	return RaceId(rawRaceId)
}

func NewRaceIdForNAR(
	year int,
	month int,
	day int,
	raceCourse RaceCourse,
	raceNo int,
) RaceId {
	rawRaceId := fmt.Sprintf("%d%s%02d%02d%02d", year, raceCourse.Value(), month, day, raceNo)
	return RaceId(rawRaceId)
}

func NewRaceIdForOverseas(
	year int,
	month int,
	day int,
	raceCourse RaceCourse,
	raceNo int,
) RaceId {
	rawRaceId := fmt.Sprintf("%d%s%02d%02d%02d", year, raceCourse.Value(), month, day, raceNo)

	switch raceCourse {
	// 2025以降netkeibaの海外レースIDが0101から始まるようになっているので特殊対応
	case Meydan:
		if year == 2025 {
			rawRaceId = fmt.Sprintf("%d%s%02d%02d%02d", year, raceCourse.Value(), 1, 1, raceNo)
		} else {
			// 日付を-1してraceIdを設定する特殊対応
			// 月をまたぐわけではないのでtimeパッケージで厳密にはやらない
			rawRaceId = fmt.Sprintf("%d%s%02d%02d%02d", year, raceCourse.Value(), month, day-1, raceNo)
		}
	case Shatin, Deauville:
		if year == 2025 {
			rawRaceId = fmt.Sprintf("%d%s%02d%02d%02d", year, raceCourse.Value(), 1, 1, raceNo)
		}
	case Longchamp:
		if year == 2025 {
			// 2025年10月開催はmonth=2
			rawRaceId = fmt.Sprintf("%d%s%02d%02d%02d", year, raceCourse.Value(), 2, 1, raceNo)
		}
	case KingAbdulaziz, SantaAnitaPark, Delmar:
		// 日付を-1してraceIdを設定する特殊対応
		// 月をまたぐわけではないのでtimeパッケージで厳密にはやらない
		rawRaceId = fmt.Sprintf("%d%s%02d%02d%02d", year, raceCourse.Value(), month, day-1, raceNo)
	}

	return RaceId(rawRaceId)
}
