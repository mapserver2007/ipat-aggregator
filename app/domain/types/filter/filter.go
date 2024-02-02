package filter

type Id byte

const (
	All                           Id = 0xFF // 全検索条件に引っ掛けるためのフィルタ
	TurfShortDistance             Id = 0xe5
	TurfMiddleDistance            Id = 0xe9
	TurfLongDistance              Id = 0xf1
	DirtShortDistance             Id = 0xe6
	DirtLongDistance              Id = 0xfa
	TurfShortDistanceJockeyTop1   Id = 0x25
	TurfMiddleDistanceJockeyTop1  Id = 0x29
	TurfLongDistanceJockeyTop1    Id = 0x31
	DirtShortDistanceJockeyTop1   Id = 0x26
	DirtLongDistanceJockeyTop1    Id = 0x3a
	TurfShortDistanceJockeyTop2   Id = 0x45
	TurfMiddleDistanceJockeyTop2  Id = 0x49
	TurfLongDistanceJockeyTop2    Id = 0x51
	DirtShortDistanceJockeyTop2   Id = 0x46
	DirtLongDistanceJockeyTop2    Id = 0x5a
	TurfShortDistanceJockeyOther  Id = 0x85
	TurfMiddleDistanceJockeyOther Id = 0x89
	TurfLongDistanceJockeyOther   Id = 0x91
	DirtShortDistanceJockeyOther  Id = 0x86
	DirtLongDistanceJockeyOther   Id = 0x9a
)

const (
	Turf           Id = 0x01
	Dirt           Id = 0x02
	ShortDistance  Id = 0x04
	MiddleDistance Id = 0x08
	LongDistance   Id = 0x10
	JokeyTop1      Id = 0x20
	JokeyTop2      Id = 0x40
	JokeyOther     Id = 0x80
)

var filterIdMap = map[Id]string{
	All: "条件なし",
	// 以下基本フィルタ
	Turf:           "芝",
	Dirt:           "ダート",
	ShortDistance:  "短距離",
	MiddleDistance: "中距離",
	LongDistance:   "長距離",
	JokeyTop1:      "ルメール",
	JokeyTop2:      "川田将雅",
	JokeyOther:     "その他騎手",
	// 以下組み合わせフィルタ
	TurfShortDistance:             "芝~1600m",
	TurfMiddleDistance:            "芝~2000m",
	TurfLongDistance:              "芝2001m~",
	DirtShortDistance:             "ダ~1600m",
	DirtLongDistance:              "ダ1601m~",
	TurfShortDistanceJockeyTop1:   "芝~1600m,ルメール",
	TurfMiddleDistanceJockeyTop1:  "芝~2000m,ルメール",
	TurfLongDistanceJockeyTop1:    "芝2001m~,ルメール",
	DirtShortDistanceJockeyTop1:   "ダ~1600m,ルメール",
	DirtLongDistanceJockeyTop1:    "ダ1601m~,ルメール",
	TurfShortDistanceJockeyTop2:   "芝~1600m,川田将雅",
	TurfMiddleDistanceJockeyTop2:  "芝~2000m,川田将雅",
	TurfLongDistanceJockeyTop2:    "芝2001m~,川田将雅",
	DirtShortDistanceJockeyTop2:   "ダ~1600m,川田将雅",
	DirtLongDistanceJockeyTop2:    "ダ1601m~,川田将雅",
	TurfShortDistanceJockeyOther:  "芝~1600m,その他騎手",
	TurfMiddleDistanceJockeyOther: "芝~2000m,その他騎手",
	TurfLongDistanceJockeyOther:   "芝2001m~,その他騎手",
	DirtShortDistanceJockeyOther:  "ダ~1600m,その他騎手",
	DirtLongDistanceJockeyOther:   "ダ1601m~,その他騎手",
}

func (i Id) Value() int {
	return int(i)
}

func (i Id) String() string {
	id, _ := filterIdMap[i]
	return id
}
