package filter

type Id uint64

const (
	All                         Id = 0x7FFF // 全検索条件に引っ掛けるためのフィルタ
	TurfSprintDistance          Id = 0x7F85
	TurfMileDistance            Id = 0x7F89
	TurfMiddleDistance          Id = 0x7FB1
	TurfLongDistance            Id = 0x7FC1
	DirtSprintDistance          Id = 0x7F86
	DirtMileDistance            Id = 0x7F8A
	DirtMiddleDistance          Id = 0x7F92
	DirtLongDistance            Id = 0x7FE2
	GoodTrackTurfSprintDistance Id = 0x3F85
	GoodTrackTurfMileDistance   Id = 0x3F89
	GoodTrackTurfMiddleDistance Id = 0x3FB1
	GoodTrackTurfLongDistance   Id = 0x3FC1
	GoodTrackDirtSprintDistance Id = 0x3F86
	GoodTrackDirtMileDistance   Id = 0x3F8A
	GoodTrackDirtMiddleDistance Id = 0x3F92
	GoodTrackDirtLongDistance   Id = 0x3FE2
	TurfClass1                  Id = 0x63FD
	TurfClass2                  Id = 0x65FD
	TurfClass3                  Id = 0x69FD
	TurfClass4                  Id = 0x71FD
	DirtClass1                  Id = 0x63FE
	DirtClass2                  Id = 0x65FE
	DirtClass3                  Id = 0x69FE
	DirtClass4                  Id = 0x71FE
	// 予想専用
	Predict1 Id = 0x4586
)

const (
	Turf            Id = 0x01
	Dirt            Id = 0x02
	Sprint          Id = 0x04
	Mile            Id = 0x08
	MiddleDistance1 Id = 0x10
	MiddleDistance2 Id = 0x20
	LongDistance    Id = 0x40
	TopJockey       Id = 0x80
	OtherJockey     Id = 0x100
	Class1          Id = 0x200
	Class2          Id = 0x400
	Class3          Id = 0x800
	Class4          Id = 0x1000
	GoodTrack       Id = 0x2000
	BadTrack        Id = 0x4000
)

var filterIdMap = map[Id]string{
	All: "条件なし",
	// 以下基本フィルタ
	Turf:            "芝",
	Dirt:            "ダート",
	Sprint:          "スプリント",
	Mile:            "マイル",
	MiddleDistance1: "中距離1",
	MiddleDistance2: "中距離2",
	LongDistance:    "長距離",
	TopJockey:       "上位騎手",
	OtherJockey:     "その他騎手",
	Class1:          "新馬・未勝利",
	Class2:          "1~3勝クラス",
	Class3:          "OP・L",
	Class4:          "重賞",
	GoodTrack:       "良馬場",
	BadTrack:        "良馬場以外",

	// 以下組み合わせフィルタ
	TurfSprintDistance:          "芝~1200m",
	TurfMileDistance:            "芝~1600m",
	TurfMiddleDistance:          "芝~2000m",
	TurfLongDistance:            "芝2001m~",
	DirtSprintDistance:          "ダ~1200m",
	DirtMileDistance:            "ダ~1600m",
	DirtMiddleDistance:          "ダ~1800m",
	DirtLongDistance:            "ダ1801m~",
	GoodTrackTurfSprintDistance: "芝良~1200m",
	GoodTrackTurfMileDistance:   "芝良~1600m",
	GoodTrackTurfMiddleDistance: "芝良~2000m",
	GoodTrackTurfLongDistance:   "芝良2001m~",
	GoodTrackDirtSprintDistance: "ダ良~1200m",
	GoodTrackDirtMileDistance:   "ダ良~1600m",
	GoodTrackDirtMiddleDistance: "ダ良~1800m",
	GoodTrackDirtLongDistance:   "ダ良1801m~",
	TurfClass1:                  "芝新馬・未勝利",
	TurfClass2:                  "芝1~3勝クラス",
	TurfClass3:                  "芝OP・L",
	TurfClass4:                  "芝重賞",
	DirtClass1:                  "ダ新馬・未勝利",
	DirtClass2:                  "ダ1~3勝クラス",
	DirtClass3:                  "ダOP・L",
	DirtClass4:                  "ダ重賞",
	Predict1:                    "ダ1400/2勝稍",
}

func (i Id) Value() int {
	return int(i)
}

func (i Id) String() string {
	id, _ := filterIdMap[i]
	return id
}
