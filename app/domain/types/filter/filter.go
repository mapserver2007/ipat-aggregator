package filter

type Id byte

const (
	All                         Id = 0xFF // 全検索条件に引っ掛けるためのフィルタ
	TurfSprintDistance          Id = 0xc5
	TurfMileDistance            Id = 0xc9
	TurfMiddleDistance          Id = 0xd1
	TurfLongDistance            Id = 0xe1
	DirtSprintDistance          Id = 0xc6
	DirtMileDistance            Id = 0xca
	DirtMiddleDistance          Id = 0xf2
	TurfSprintDistanceTopJockey Id = 0x45
	TurfMileDistanceTopJockey   Id = 0x49
	TurfMiddleDistanceTopJockey Id = 0x51
	TurfLongDistanceTopJockey   Id = 0x61
	DirtSprintDistanceTopJockey Id = 0x46
	DirtMileDistanceTopJockey   Id = 0x4a
	DirtMiddleDistanceTopJockey Id = 0x72
)

const (
	Turf           Id = 0x01
	Dirt           Id = 0x02
	Sprint         Id = 0x04
	Mile           Id = 0x08
	MiddleDistance Id = 0x10
	LongDistance   Id = 0x20
	TopJockey      Id = 0x40
	OtherJockey    Id = 0x80
)

var filterIdMap = map[Id]string{
	All: "条件なし",
	// 以下基本フィルタ
	Turf:           "芝",
	Dirt:           "ダート",
	Sprint:         "スプリント",
	Mile:           "マイル",
	MiddleDistance: "中距離",
	LongDistance:   "長距離",
	TopJockey:      "上位騎手",
	OtherJockey:    "その他騎手",
	// 以下組み合わせフィルタ
	TurfSprintDistance:          "芝~1200m",
	TurfMileDistance:            "芝~1600m",
	TurfMiddleDistance:          "芝~2000m",
	TurfLongDistance:            "芝2001m~",
	DirtSprintDistance:          "ダ~1200m",
	DirtMileDistance:            "ダ~1600m",
	DirtMiddleDistance:          "ダ1601m~",
	TurfSprintDistanceTopJockey: "芝~1200m,上位騎手",
	TurfMileDistanceTopJockey:   "芝~1600m,上位騎手",
	TurfMiddleDistanceTopJockey: "芝~2000m,上位騎手",
	TurfLongDistanceTopJockey:   "芝2001m~,上位騎手",
	DirtSprintDistanceTopJockey: "ダ~1200m,上位騎手",
	DirtMileDistanceTopJockey:   "ダ~1600m,上位騎手",
	DirtMiddleDistanceTopJockey: "ダ1601m~,上位騎手",
}

func (i Id) Value() int {
	return int(i)
}

func (i Id) String() string {
	id, _ := filterIdMap[i]
	return id
}
