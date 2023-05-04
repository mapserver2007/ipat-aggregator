package value_object

import (
	"strings"
)

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
	Monbetsu  = 30
	Morioka   = 35
	Urawa     = 42
	Hunabashi = 43
	Ooi       = 44
	Kawasaki  = 45
	Nagoya    = 48
	Sonoda    = 50
	Kouchi    = 54
	Saga      = 55
	Longchamp = 90 // 値はダミー
	Deauville = 91 // 値はダミー
	Shatin    = 92 // 値はダミー
	Meydan    = 93 // 値はダミー
	Overseas  = 99 // その他海外
)

var raceCourseMap = map[RaceCourse]string{
	Tokyo:        "東京",
	Nakayama:     "中山",
	Hanshin:      "阪神",
	Kyoto:        "京都",
	Chukyo:       "中京",
	Kokura:       "小倉",
	Niigata:      "新潟",
	Hakodate:     "函館",
	Sapporo:      "札幌",
	Fukushima:    "福島",
	Monbetsu:     "門別",
	Morioka:      "盛岡",
	Ooi:          "大井",
	Kawasaki:     "川崎",
	Nagoya:       "名古屋",
	Sonoda:       "園田",
	Urawa:        "浦和",
	Hunabashi:    "船橋",
	Kouchi:       "高知",
	Saga:         "佐賀",
	Longchamp:    "パリロンシャン（フランス）",
	Deauville:    "ドーヴィル（フランス）",
	Shatin:       "シャティン（香港）",
	Meydan:       "メイダン（ＵＡＥ）",
	Overseas:     "海外",
	UnknownPlace: "不明",
}

// 海外だけは開催場所が文字列なので
var raceCourseOverseaIdMap = map[RaceCourse]string{
	Longchamp: "C8",
	Deauville: "C4",
	Shatin:    "H1",
	Meydan:    "J0",
}

func (r RaceCourse) Name() string {
	return convertToRaceCourseName(r)
}

func (r RaceCourse) Value() int {
	return int(r)
}

func (r RaceCourse) Organizer() Organizer {
	switch r {
	case Tokyo, Nakayama, Hanshin, Kyoto, Chukyo, Kokura, Niigata, Hakodate, Sapporo, Fukushima:
		return JRA
	case Monbetsu, Morioka, Ooi, Kawasaki, Urawa, Hunabashi, Nagoya, Sonoda, Kouchi, Saga:
		return NAR
	case Longchamp, Deauville, Shatin, Meydan:
		return OverseaOrganizer
	}

	return UnknownOrganizer
}

func convertToRaceCourseName(r RaceCourse) string {
	if v, ok := raceCourseMap[r]; ok {
		return v
	}
	return ""
}

func ConvertToOverseaRaceCourseId(r RaceCourse) string {
	if v, ok := raceCourseOverseaIdMap[r]; ok {
		return v
	}
	return ""
}

func ConvertToRaceCourse(name string) RaceCourse {
	for k, v := range raceCourseMap {
		if v == name {
			return k
		}
	}

	// その他海外
	if partialContains([]string{"イギリス", "フランス", "香港"}, name) {
		return Overseas
	}

	return UnknownPlace
}

func partialContains(elems []string, str string) bool {
	for _, elem := range elems {
		if strings.Contains(str, elem) {
			return true
		}
	}
	return false
}
