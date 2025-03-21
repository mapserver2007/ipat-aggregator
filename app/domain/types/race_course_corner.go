package types

type RaceCourseCornerIndex int

const (
	HanshinTurfOuterCorner  RaceCourseCornerIndex = 219
	KyotoTurfOuterCorner    RaceCourseCornerIndex = 212
	NakayamaTurfOuterCorner RaceCourseCornerIndex = 208
	HanshinTurfInnerCorner  RaceCourseCornerIndex = 197
	TokyoTurfCorner         RaceCourseCornerIndex = 191
	HanshinDirtCorner       RaceCourseCornerIndex = 182
	SapporoTurfCorner       RaceCourseCornerIndex = 169
	KyotoTurfInnerCorner    RaceCourseCornerIndex = 168
	TokyoDirtCorner         RaceCourseCornerIndex = 161
	ChukyoTurfCorner        RaceCourseCornerIndex = 156
	HakodateTurfCorner      RaceCourseCornerIndex = 147
	NakayamaTurfInnerCorner RaceCourseCornerIndex = 146
	KyotoDirtCorner         RaceCourseCornerIndex = 144
	SapporoDirtCorner       RaceCourseCornerIndex = 141
	FukushimaTurfCorner     RaceCourseCornerIndex = 139
	KokuraTurfCorner        RaceCourseCornerIndex = 134
	NiigataTurfCorner       RaceCourseCornerIndex = 128
	NakayamaDirtCorner      RaceCourseCornerIndex = 122
	ChukyoDirtCorner        RaceCourseCornerIndex = 117
	HakodateDirtCorner      RaceCourseCornerIndex = 116
	KokuraDirtCorner        RaceCourseCornerIndex = 106
	FukushimaDirtCorner     RaceCourseCornerIndex = 103
	NiigataDirtCorner       RaceCourseCornerIndex = 102
	NiigataTurfStraight     RaceCourseCornerIndex = 0
	UnknownCorner           RaceCourseCornerIndex = -1
)

var raceCourseCornerIndexMap = map[RaceCourseCornerIndex]string{
	HanshinTurfOuterCorner:  "阪神芝外",
	KyotoTurfOuterCorner:    "京都芝外",
	NakayamaTurfOuterCorner: "中山芝外",
	HanshinTurfInnerCorner:  "阪神芝内",
	TokyoTurfCorner:         "東京芝",
	HanshinDirtCorner:       "阪神ダ",
	SapporoTurfCorner:       "札幌芝",
	KyotoTurfInnerCorner:    "京都芝内",
	TokyoDirtCorner:         "東京ダ",
	ChukyoTurfCorner:        "中京芝",
	HakodateTurfCorner:      "函館芝",
	NakayamaTurfInnerCorner: "中山芝内",
	KyotoDirtCorner:         "京都ダ",
	SapporoDirtCorner:       "札幌ダ",
	FukushimaTurfCorner:     "福島芝",
	KokuraTurfCorner:        "小倉芝",
	NiigataTurfCorner:       "新潟芝",
	NakayamaDirtCorner:      "中山ダ",
	ChukyoDirtCorner:        "中京ダ",
	HakodateDirtCorner:      "函館ダ",
	KokuraDirtCorner:        "小倉ダ",
	FukushimaDirtCorner:     "福島ダ",
	NiigataDirtCorner:       "新潟ダ",
	NiigataTurfStraight:     "新潟芝直線",
	UnknownCorner:           "不明",
}

func (r RaceCourseCornerIndex) Name() string {
	if v, ok := raceCourseCornerIndexMap[r]; ok {
		return v
	}
	return ""
}

func (r RaceCourseCornerIndex) Value() int {
	return int(r)
}
