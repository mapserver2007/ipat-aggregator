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
	// 海外の場合、日をまたぐケースがあり開催日時とrace_idが一致しない場合がある(例：3月のドバイ)
	if raceCourse == Meydan || raceCourse == KingAbdulaziz || raceCourse == SantaAnitaPark || raceCourse == Delmar {
		// 日付を-1してraceIdを設定する特殊対応
		// 月をまたぐわけではないのでtimeパッケージで厳密にはやらない
		rawRaceId = fmt.Sprintf("%d%s%02d%02d%02d", year, raceCourse, month, day-1, raceNo)
	}
	return RaceId(rawRaceId)
}
