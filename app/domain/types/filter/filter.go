package filter

type Id uint64

const (
	All                         Id = 0x7FFFF // 全検索条件に引っ掛けるためのフィルタ
	TurfSprintDistance          Id = 0x7FF85
	TurfMileDistance            Id = 0x7FF89
	TurfMiddleDistance          Id = 0x7FFB1
	TurfLongDistance            Id = 0x7FFC1
	DirtSprintDistance          Id = 0x7FF86
	DirtMileDistance            Id = 0x7FF8A
	DirtMiddleDistance          Id = 0x7FF92
	DirtLongDistance            Id = 0x7FFE2
	GoodTrackTurfSprintDistance Id = 0x4FF85
	GoodTrackTurfMileDistance   Id = 0x4FF89
	GoodTrackTurfMiddleDistance Id = 0x4FFB1
	GoodTrackTurfLongDistance   Id = 0x4FFC1
	GoodTrackDirtSprintDistance Id = 0x4FF86
	GoodTrackDirtMileDistance   Id = 0x4FF8A
	GoodTrackDirtMiddleDistance Id = 0x4FF92
	GoodTrackDirtLongDistance   Id = 0x4FFE2
	TurfClass1                  Id = 0x483FD
	TurfClass2                  Id = 0x485FD
	TurfClass3                  Id = 0x489FD
	TurfClass4                  Id = 0x491FD
	TurfClass5                  Id = 0x4A1FD
	TurfClass6                  Id = 0x4C1FD
	DirtClass1                  Id = 0x483FE
	DirtClass2                  Id = 0x485FE
	DirtClass3                  Id = 0x489FE
	DirtClass4                  Id = 0x491FE
	DirtClass5                  Id = 0x4A1FE
	DirtClass6                  Id = 0x4C1FE
	DirtBadConditionClass1      Id = 0x503FE
	DirtBadConditionClass2      Id = 0x505FE
	DirtBadConditionClass3      Id = 0x509FE
	DirtBadConditionClass4      Id = 0x511FE
	DirtBadConditionClass5      Id = 0x521FE
	DirtBadConditionClass6      Id = 0x541FE
	// 予想専用
	Predict1 Id = 0x4586
)

const (
	Turf                Id = 0x01
	Dirt                Id = 0x02
	Sprint              Id = 0x04
	Mile                Id = 0x08
	MiddleDistance1     Id = 0x10
	MiddleDistance2     Id = 0x20
	LongDistance        Id = 0x40
	TopJockey           Id = 0x80
	OtherJockey         Id = 0x100
	Class1              Id = 0x200
	Class2              Id = 0x400
	Class3              Id = 0x800
	Class4              Id = 0x1000
	Class5              Id = 0x2000
	Class6              Id = 0x4000
	GoodTrack           Id = 0x8000
	BadTrack            Id = 0x10000
	SmallNumberOfHorses Id = 0x20000
	LargeNumberOfHorses Id = 0x40000
)

var filterIdMap = map[Id]string{
	All: "条件なし",
	// 以下基本フィルタ
	Turf:                "芝",
	Dirt:                "ダート",
	Sprint:              "スプリント",
	Mile:                "マイル",
	MiddleDistance1:     "中距離1",
	MiddleDistance2:     "中距離2",
	LongDistance:        "長距離",
	TopJockey:           "上位騎手",
	OtherJockey:         "その他騎手",
	Class1:              "未勝利",
	Class2:              "1勝クラス",
	Class3:              "2勝クラス",
	Class4:              "3勝クラス",
	Class5:              "OP・L",
	Class6:              "重賞",
	GoodTrack:           "良馬場",
	BadTrack:            "良馬場以外",
	SmallNumberOfHorses: "少頭数",
	LargeNumberOfHorses: "多頭数",

	// 以下組み合わせフィルタ
	TurfSprintDistance:          "芝~1200m",
	TurfMileDistance:            "芝~1600m",
	TurfMiddleDistance:          "芝~2000m",
	TurfLongDistance:            "芝2001m~",
	DirtSprintDistance:          "ダ~1200m",
	DirtMileDistance:            "ダ~1600m",
	DirtMiddleDistance:          "ダ~1800m",
	DirtLongDistance:            "ダ1801m~",
	GoodTrackTurfSprintDistance: "芝良多~1200m",
	GoodTrackTurfMileDistance:   "芝良多~1600m",
	GoodTrackTurfMiddleDistance: "芝良多~2000m",
	GoodTrackTurfLongDistance:   "芝良多2001m~",
	GoodTrackDirtSprintDistance: "ダ良多~1200m",
	GoodTrackDirtMileDistance:   "ダ良多~1600m",
	GoodTrackDirtMiddleDistance: "ダ良多~1800m",
	GoodTrackDirtLongDistance:   "ダ良多1801m~",
	TurfClass1:                  "芝良多未勝利",
	TurfClass2:                  "芝良多1勝",
	TurfClass3:                  "芝良多2勝",
	TurfClass4:                  "芝良多3勝",
	TurfClass5:                  "芝良多OP・L",
	TurfClass6:                  "芝良多重賞",
	DirtClass1:                  "ダ良多未勝利",
	DirtClass2:                  "ダ良多1勝",
	DirtClass3:                  "ダ良多2勝",
	DirtClass4:                  "ダ良多3勝",
	DirtClass5:                  "ダ良多OP・L",
	DirtClass6:                  "ダ良多重賞",
	DirtBadConditionClass1:      "ダ稍重不多未勝利",
	DirtBadConditionClass2:      "ダ稍重不多1勝",
	DirtBadConditionClass3:      "ダ稍重不多2勝",
	DirtBadConditionClass4:      "ダ稍重不多3勝",
	DirtBadConditionClass5:      "ダ稍重不多OP・L",
	DirtBadConditionClass6:      "ダ稍重不多重賞",

	Predict1: "ダ1400/2勝稍",
}

func (i Id) Value() int {
	return int(i)
}

func (i Id) String() string {
	id, _ := filterIdMap[i]
	return id
}
