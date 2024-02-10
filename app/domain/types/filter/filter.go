package filter

type Id uint64

const (
	All                         Id = 0xFFFF // 全検索条件に引っ掛けるためのフィルタ
	TurfSprintDistance          Id = 0x1F85
	TurfMileDistance            Id = 0x1F89
	TurfMiddleDistance          Id = 0x1FB1
	TurfLongDistance            Id = 0x1FC1
	DirtSprintDistance          Id = 0x1F86
	DirtMileDistance            Id = 0x1F8A
	DirtMiddleDistance          Id = 0x1F92
	DirtLongDistance            Id = 0x1FE2
	TurfSprintDistanceTopJockey Id = 0x1E85
	TurfMileDistanceTopJockey   Id = 0x1E89
	TurfMiddleDistanceTopJockey Id = 0x1EB1
	TurfLongDistanceTopJockey   Id = 0x1EC1
	DirtSprintDistanceTopJockey Id = 0x1E86
	DirtMileDistanceTopJockey   Id = 0x1E8A
	DirtMiddleDistanceTopJockey Id = 0x1E92
	DirtLongDistanceTopJockey   Id = 0x1EE2
	TurfClass1                  Id = 0x3FD
	TurfClass2                  Id = 0x5FD
	TurfClass3                  Id = 0x9FD
	TurfClass4                  Id = 0x11FD
	DirtClass1                  Id = 0x3FE
	DirtClass2                  Id = 0x5FE
	DirtClass3                  Id = 0x9FE
	DirtClass4                  Id = 0x11FE
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

	// 以下組み合わせフィルタ
	TurfSprintDistance:          "芝~1200m",
	TurfMileDistance:            "芝~1600m",
	TurfMiddleDistance:          "芝~2000m",
	TurfLongDistance:            "芝2001m~",
	DirtSprintDistance:          "ダ~1200m",
	DirtMileDistance:            "ダ~1600m",
	DirtMiddleDistance:          "ダ~1800m",
	DirtLongDistance:            "ダ1801m~",
	TurfSprintDistanceTopJockey: "芝~1200m,上位騎手",
	TurfMileDistanceTopJockey:   "芝~1600m,上位騎手",
	TurfMiddleDistanceTopJockey: "芝~2000m,上位騎手",
	TurfLongDistanceTopJockey:   "芝2001m~,上位騎手",
	DirtSprintDistanceTopJockey: "ダ~1200m,上位騎手",
	DirtMileDistanceTopJockey:   "ダ~1600m,上位騎手",
	DirtMiddleDistanceTopJockey: "ダ~1800m,上位騎手",
	DirtLongDistanceTopJockey:   "ダ1801m~,上位騎手",
	TurfClass1:                  "芝新馬・未勝利",
	TurfClass2:                  "芝1~3勝クラス",
	TurfClass3:                  "芝OP・L",
	TurfClass4:                  "芝重賞",
	DirtClass1:                  "ダ新馬・未勝利",
	DirtClass2:                  "ダ1~3勝クラス",
	DirtClass3:                  "ダOP・L",
	DirtClass4:                  "ダ重賞",
}

func (i Id) Value() int {
	return int(i)
}

func (i Id) String() string {
	id, _ := filterIdMap[i]
	return id
}
