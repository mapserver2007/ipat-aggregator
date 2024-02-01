package filter

type Id byte

const (
	All                     Id = 0xFF // 全検索条件に引っ掛けるためのフィルタ
	TurfShortDistance       Id = 0x35
	TurfLongDistance        Id = 0x39
	DirtShortDistance       Id = 0x36
	DirtLongDistance        Id = 0x3a
	TurfSmallNumberStarters Id = 0x1d
	TurfLargeNumberStarters Id = 0x2d
	DirtSmallNumberStarters Id = 0x1e
	DirtLargeNumberStarters Id = 0x2e
)

const (
	Turf                Id = 0x01
	Dirt                Id = 0x02
	ShortDistance       Id = 0x04
	LongDistance        Id = 0x08
	SmallNumberStarters Id = 0x10
	LargeNumberStarters Id = 0x20
)

var filterIdMap = map[Id]string{
	All:                     "条件なし",
	TurfShortDistance:       "芝短距離",
	TurfLongDistance:        "芝中長距離",
	DirtShortDistance:       "ダ短距離",
	DirtLongDistance:        "ダ中長距離",
	TurfSmallNumberStarters: "芝少頭数",
	TurfLargeNumberStarters: "芝多頭数",
	DirtSmallNumberStarters: "ダ少頭数",
	DirtLargeNumberStarters: "ダ多頭数",
	Turf:                    "芝",
	Dirt:                    "ダート",
	ShortDistance:           "短距離",
	LongDistance:            "長距離",
	SmallNumberStarters:     "少頭数",
	LargeNumberStarters:     "多頭数",
}

func (i Id) Value() int {
	return int(i)
}

func (i Id) String() string {
	id, _ := filterIdMap[i]
	return id
}
