package filter

type Id byte

const (
	All               Id = 0xFF // 全検索条件に引っ掛けるためのフィルタ
	TurfShortDistance Id = 0x09
	TurfLongDistance  Id = 0x11
	DirtShortDistance Id = 0x0A
	DirtLongDistance  Id = 0x12
)

const (
	Turf          Id = 0x01
	Dirt          Id = 0x02
	Jump          Id = 0x04
	ShortDistance Id = 0x08
	LongDistance  Id = 0x10
)

var filterIdMap = map[Id]string{
	All:               "条件なし",
	TurfShortDistance: "芝短距離",
	TurfLongDistance:  "芝中長距離",
	DirtShortDistance: "ダ短距離",
	DirtLongDistance:  "ダ中長距離",
	Turf:              "芝",
	Dirt:              "ダート",
	Jump:              "障害",
	ShortDistance:     "短距離",
	LongDistance:      "長距離",
}

func (i Id) Value() int {
	return int(i)
}

func (i Id) String() string {
	id, _ := filterIdMap[i]
	return id
}
