package types

import "strings"

type RaceCourse string

const (
	UnknownPlace   RaceCourse = "00"
	Sapporo        RaceCourse = "01"
	Hakodate       RaceCourse = "02"
	Fukushima      RaceCourse = "03"
	Niigata        RaceCourse = "04"
	Tokyo          RaceCourse = "05"
	Nakayama       RaceCourse = "06"
	Chukyo         RaceCourse = "07"
	Kyoto          RaceCourse = "08"
	Hanshin        RaceCourse = "09"
	Kokura         RaceCourse = "10"
	Monbetsu       RaceCourse = "30"
	Morioka        RaceCourse = "35"
	Urawa          RaceCourse = "42"
	Hunabashi      RaceCourse = "43"
	Ooi            RaceCourse = "44"
	Kawasaki       RaceCourse = "45"
	Kanazawa       RaceCourse = "46"
	Nagoya         RaceCourse = "48"
	Sonoda         RaceCourse = "50"
	Kouchi         RaceCourse = "54"
	Saga           RaceCourse = "55"
	Longchamp      RaceCourse = "C8"
	Deauville      RaceCourse = "C4"
	Shatin         RaceCourse = "H1"
	Meydan         RaceCourse = "J0"
	SantaAnitaPark RaceCourse = "F3"
	KingAbdulaziz  RaceCourse = "P0"
	York           RaceCourse = "AH"
	Delmar         RaceCourse = "FP"
	Overseas       RaceCourse = "99" // その他海外
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
	KingAbdulaziz:  "Ｋアブドゥルアジーズ（サウジアラビア）",
	York:           "ヨーク（イギリス）",
	Delmar:         "デルマー（アメリカ）",
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
	case Longchamp, Deauville, Shatin, Meydan, SantaAnitaPark, KingAbdulaziz, York, Delmar:
		return true
	}
	return false
}
