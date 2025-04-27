package types

import "fmt"

type RaceId string

func (r RaceId) String() string {
	return string(r)
}

func NewRaceIdForJRA(
	year int,
	day int,
	raceCourse string,
	raceRound int,
	raceNo int,
) RaceId {
	rawRaceId := fmt.Sprintf("%d%s%02d%02d%02d", year, raceCourse, raceRound, day, raceNo)
	return RaceId(rawRaceId)
}

func NewRaceIdForNAR(
	year int,
	month int,
	day int,
	raceCourse string,
	raceNo int,
) RaceId {
	rawRaceId := fmt.Sprintf("%d%s%02d%02d%02d", year, raceCourse, month, day, raceNo)
	return RaceId(rawRaceId)
}

func NewRaceIdForOverseas(
	year int,
	month int,
	day int,
	raceCourse string,
	raceNo int,
) RaceId {
	rawRaceId := fmt.Sprintf("%d%s%02d%02d%02d", year, raceCourse, month, day, raceNo)

	switch raceCourse {
	case Meydan:
		// 2025ドバイ開催のIDはなぜか0101になっているので特殊対応
		if year == 2025 {
			rawRaceId = fmt.Sprintf("%d%s%02d%02d%02d", year, raceCourse, 1, 1, raceNo)
		} else {
			// 日付を-1してraceIdを設定する特殊対応
			// 月をまたぐわけではないのでtimeパッケージで厳密にはやらない
			rawRaceId = fmt.Sprintf("%d%s%02d%02d%02d", year, raceCourse, month, day-1, raceNo)
		}
	case Shatin:
		// 2025シャティン開催のIDはなぜか0101になっているので特殊対応
		if year == 2025 {
			rawRaceId = fmt.Sprintf("%d%s%02d%02d%02d", year, raceCourse, 1, 1, raceNo)
		}
	case KingAbdulaziz, SantaAnitaPark, Delmar:
		// 日付を-1してraceIdを設定する特殊対応
		// 月をまたぐわけではないのでtimeパッケージで厳密にはやらない
		rawRaceId = fmt.Sprintf("%d%s%02d%02d%02d", year, raceCourse, month, day-1, raceNo)
	}

	return RaceId(rawRaceId)
}
