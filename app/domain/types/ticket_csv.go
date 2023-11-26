package types

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

func (d *RaceDate) Year() int {
	return toDate(d).Year()
}

func (d *RaceDate) Month() int {
	return int(toDate(d).Month())
}

func (d *RaceDate) Day() int {
	return toDate(d).Day()
}

func (d *RaceDate) Format(layout string) string {
	return toDate(d).Format(layout)
}

func toDate(d *RaceDate) time.Time {
	date, err := time.Parse("20060102", strconv.Itoa(int(*d)))
	if err != nil {
		panic(err)
	}
	return date
}

type RaceCourse int

const (
	UnknownPlace RaceCourse = iota
	Sapporo
	Hakodate
	Fukushima
	Niigata
	Tokyo
	Nakayama
	Chukyo
	Kyoto
	Hanshin
	Kokura
	Monbetsu       RaceCourse = 30
	Morioka        RaceCourse = 35
	Urawa          RaceCourse = 42
	Hunabashi      RaceCourse = 43
	Ooi            RaceCourse = 44
	Kawasaki       RaceCourse = 45
	Kanazawa       RaceCourse = 46
	Nagoya         RaceCourse = 48
	Sonoda         RaceCourse = 50
	Kouchi         RaceCourse = 54
	Saga           RaceCourse = 55
	Longchamp      RaceCourse = 90 // 値はダミー
	Deauville      RaceCourse = 91 // 値はダミー
	Shatin         RaceCourse = 92 // 値はダミー
	Meydan         RaceCourse = 93 // 値はダミー
	SantaAnitaPark RaceCourse = 94 // 値はダミー
	Overseas       RaceCourse = 99 // その他海外
)

var raceCourseMap = map[RaceCourse]string{
	Tokyo:          "東京",
	Nakayama:       "中山",
	Hanshin:        "阪神",
	Kyoto:          "京都",
	Chukyo:         "中京",
	Kokura:         "小倉",
	Niigata:        "新潟",
	Hakodate:       "函館",
	Sapporo:        "札幌",
	Fukushima:      "福島",
	Monbetsu:       "門別",
	Morioka:        "盛岡",
	Ooi:            "大井",
	Kawasaki:       "川崎",
	Nagoya:         "名古屋",
	Sonoda:         "園田",
	Urawa:          "浦和",
	Hunabashi:      "船橋",
	Kanazawa:       "金沢",
	Kouchi:         "高知",
	Saga:           "佐賀",
	Longchamp:      "パリロンシャン（フランス）",
	Deauville:      "ドーヴィル（フランス）",
	Shatin:         "シャティン（香港）",
	Meydan:         "メイダン（ＵＡＥ）",
	SantaAnitaPark: "サンタアニタパーク（アメリカ）",
	Overseas:       "海外",
	UnknownPlace:   "不明",
}

// 海外だけは開催場所が文字列なので個別管理
var raceCourseOverseaIdMap = map[RaceCourse]string{
	Longchamp:      "C8",
	Deauville:      "C4",
	Shatin:         "H1",
	Meydan:         "J0",
	SantaAnitaPark: "F3",
}

func NewRaceCourse(s string) RaceCourse {
	raceCourse := UnknownPlace
	for k, v := range raceCourseMap {
		if v == s {
			return k
		}
	}
	// その他海外
	for _, elem := range []string{"イギリス", "フランス", "香港"} {
		if strings.Contains(s, elem) {
			raceCourse = Overseas
			break
		}
	}

	return raceCourse
}

type TicketType int

const (
	UnknownTicket TicketType = iota
	Win
	Place
	BracketQuinella
	Quinella
	Exacta
	ExactaWheelOfFirst
	QuinellaPlace
	QuinellaPlaceWheel
	Trio
	TrioFormation
	TrioWheelOfFirst
	Trifecta
	TrifectaFormation
	TrifectaWheelOfFirst
	TrifectaWheelOfSecondMulti
)

var ticketTypeMap = map[TicketType]string{
	Win:                        "単勝",
	Place:                      "複勝",
	BracketQuinella:            "枠連",
	Quinella:                   "馬連",
	Exacta:                     "馬単",
	ExactaWheelOfFirst:         "馬単1着ながし",
	QuinellaPlace:              "ワイド",
	QuinellaPlaceWheel:         "ワイドながし",
	Trio:                       "3連複",
	TrioFormation:              "3連複フォーメーション",
	TrioWheelOfFirst:           "3連複軸1頭ながし",
	Trifecta:                   "3連単",
	TrifectaFormation:          "3連単フォーメーション",
	TrifectaWheelOfFirst:       "3連単1着ながし",
	TrifectaWheelOfSecondMulti: "3連単軸2頭ながしマルチ",
	UnknownTicket:              "不明",
}

func NewTicketType(name string) TicketType {
	for key, value := range ticketTypeMap {
		if value == name {
			return key
		}
	}

	return UnknownTicket
}

func (b TicketType) Name() string {
	name, _ := ticketTypeMap[b]
	return name
}

func (b TicketType) OriginTicketType() TicketType {
	switch b {
	case ExactaWheelOfFirst:
		return Exacta
	case QuinellaPlaceWheel:
		return QuinellaPlace
	case TrioFormation, TrioWheelOfFirst:
		return Trio
	case TrifectaFormation, TrifectaWheelOfFirst, TrifectaWheelOfSecondMulti:
		return Trifecta
	}
	return b
}

func (b TicketType) Value() int {
	return int(b)
}

type TicketResult int

const (
	TicketHit   TicketResult = 1
	TicketUnHit TicketResult = 2
)

type BetNumber string

const (
	DefaultQuinellaSeparator = "―"
	QuinellaSeparator        = "-"
	ExactaSeparator          = "→"
)

func NewBetNumber(number string) BetNumber {
	number = strings.Replace(number, DefaultQuinellaSeparator, QuinellaSeparator, -1)
	return BetNumber(number)
}

func (b BetNumber) List() []int {
	separators := fmt.Sprintf("[%s,%s]", QuinellaSeparator, ExactaSeparator)
	list := regexp.MustCompile(separators).Split(string(b), -1)
	var betNumbers []int
	for _, s := range list {
		betNumber, _ := strconv.Atoi(s)
		betNumbers = append(betNumbers, betNumber)
	}

	return betNumbers
}

func (b BetNumber) String() string {
	// 三連複はダッシュなのでハイフンでつなぐ
	if strings.Contains(string(b), QuinellaSeparator) {
		return strings.Replace(string(b), QuinellaSeparator, "-", -1)
	}
	return string(b)
}
