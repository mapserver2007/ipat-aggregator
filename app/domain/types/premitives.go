package types

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Payment int

type Payout int

type BetCount int

type HitCount int

type RaceCount int

func (p Payment) Value() int {
	return int(p)
}

func (p Payout) Value() int {
	return int(p)
}

func (b BetCount) Value() int {
	return int(b)
}

func (h HitCount) Value() int {
	return int(h)
}

func (r RaceCount) Value() int {
	return int(r)
}

type RaceId string

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
	return RaceId(rawRaceId)
}

type RacingNumberId string

func NewRacingNumberId(
	date RaceDate,
	raceCourse RaceCourse,
) RacingNumberId {
	return RacingNumberId(fmt.Sprintf("%d_%s", date, raceCourse.Value()))
}

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
	return toDate(d).Year()
}

func (d RaceDate) Month() int {
	return int(toDate(d).Month())
}

func (d RaceDate) Day() int {
	return toDate(d).Day()
}

func (d RaceDate) Format(layout string) string {
	return toDate(d).Format(layout)
}

func toDate(d RaceDate) time.Time {
	date, err := time.Parse("20060102", strconv.Itoa(int(d)))
	if err != nil {
		panic(err)
	}
	return date
}

type RaceCourse string

const (
	UnknownPlace   = "00"
	Sapporo        = "01"
	Hakodate       = "02"
	Fukushima      = "03"
	Niigata        = "04"
	Tokyo          = "05"
	Nakayama       = "06"
	Chukyo         = "07"
	Kyoto          = "08"
	Hanshin        = "09"
	Kokura         = "10"
	Monbetsu       = "30"
	Morioka        = "35"
	Urawa          = "42"
	Hunabashi      = "43"
	Ooi            = "44"
	Kawasaki       = "45"
	Kanazawa       = "46"
	Nagoya         = "48"
	Sonoda         = "50"
	Kouchi         = "54"
	Saga           = "55"
	Longchamp      = "C8"
	Deauville      = "C4"
	Shatin         = "H1"
	Meydan         = "J0"
	SantaAnitaPark = "F3"
	Overseas       = "99" // その他海外
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

	return RaceCourse(raceCourse)
}

func (r RaceCourse) Name() string {
	if v, ok := raceCourseMap[r]; ok {
		return v
	}
	return ""
}

func (r RaceCourse) Value() string {
	return string(r)
}

func (r RaceCourse) JRA() bool {
	switch r {
	case Tokyo, Nakayama, Hanshin, Kyoto, Chukyo, Kokura, Niigata, Hakodate, Sapporo, Fukushima:
		return true
	}
	return false
}

func (r RaceCourse) NAR() bool {
	switch r {
	case Monbetsu, Morioka, Ooi, Kawasaki, Urawa, Hunabashi, Kanazawa, Nagoya, Sonoda, Kouchi, Saga:
		return true
	}
	return false
}

func (r RaceCourse) Oversea() bool {
	switch r {
	case Longchamp, Deauville, Shatin, Meydan, SantaAnitaPark:
		return true
	}
	return false
}

type Organizer int

const (
	UnknownOrganizer Organizer = iota
	JRA
	NAR
	OverseaOrganizer
)

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

type GradeClass int

const (
	NonGrade       GradeClass = 0
	Grade1         GradeClass = 1
	Grade2         GradeClass = 2
	Grade3         GradeClass = 3
	LocalGrade     GradeClass = 4
	OpenClass      GradeClass = 5
	JumpGrade1     GradeClass = 10
	JumpGrade2     GradeClass = 11
	JumpGrade3     GradeClass = 12
	ListedClass    GradeClass = 15
	Jpn1           GradeClass = 19
	Jpn2           GradeClass = 20
	Jpn3           GradeClass = 21
	Maiden         GradeClass = 31 // 未勝利
	OneWinClass    GradeClass = 32 // 1勝クラス
	TwoWinClass    GradeClass = 33 // 2勝クラス
	ThreeWinClass  GradeClass = 34 // 3勝クラス
	JumpMaiden     GradeClass = 35 // 障害未勝利
	JumpOpenClass  GradeClass = 36 // 障害オープン
	MakeDebut      GradeClass = 37 // 新馬
	AllowanceClass GradeClass = 98 // Class1-3は特別戦、AllowanceClassは非特別戦の条件戦
	NonGradeClass  GradeClass = 99 // リステッド,OP,条件戦をまとめるためのクラス
)

var gradeClassMap = map[GradeClass]string{
	NonGrade:       "なし",
	Grade1:         "G1",
	Grade2:         "G2",
	Grade3:         "G3",
	LocalGrade:     "地方重賞",
	OpenClass:      "OP/L/地方重賞",
	JumpGrade1:     "JG1",
	JumpGrade2:     "JG2",
	JumpGrade3:     "JG3",
	ListedClass:    "L",
	Jpn1:           "Jpn1",
	Jpn2:           "Jpn2",
	Jpn3:           "Jpn3",
	Maiden:         "未勝利",
	MakeDebut:      "新馬",
	OneWinClass:    "1勝クラス",
	TwoWinClass:    "2勝クラス",
	ThreeWinClass:  "3勝クラス",
	JumpMaiden:     "障害未勝利",
	JumpOpenClass:  "障害オープン",
	AllowanceClass: "条件戦",
	NonGradeClass:  "平場",
}

func (g GradeClass) String() string {
	gradeClassName, _ := gradeClassMap[g]
	return gradeClassName
}

type CourseCategory int

const (
	NonCourseCategory CourseCategory = iota
	Turf
	Dirt
	Jump
)

var courseCategoryMap = map[CourseCategory]string{
	NonCourseCategory: "不明",
	Turf:              "芝",
	Dirt:              "ダート",
	Jump:              "障害",
}

func NewCourseCategory(name string) CourseCategory {
	var courseCategory CourseCategory
	switch name {
	case "芝":
		courseCategory = Turf
	case "ダ":
		courseCategory = Dirt
	case "障":
		courseCategory = Jump
	}

	return courseCategory
}

func (c CourseCategory) String() string {
	courseCategoryName, _ := courseCategoryMap[c]
	return courseCategoryName
}

type JockeyId int

func (j JockeyId) Format() string {
	return fmt.Sprintf("%05d", j)
}

func (j JockeyId) Value() int {
	return int(j)
}
