package filter

type Id uint64

const (
	All                     Id = 0x7FFFFF // 全検索条件に引っ掛けるためのフィルタ
	TurfShortDistance1      Id = 0x7FFE05
	TurfShortDistance2      Id = 0x7FFE09
	TurfShortDistance3      Id = 0x7FFE11
	TurfMiddleDistance1     Id = 0x7FFE61
	TurfMiddleDistance2     Id = 0x7FFE81
	TurfLongDistance        Id = 0x7FFF01
	DirtShortDistance1      Id = 0x7FFE06
	DirtShortDistance2      Id = 0x7FFE0A
	DirtShortDistance3      Id = 0x7FFE12
	DirtMiddleDistance1     Id = 0x7FFE22
	DirtMiddleDistance2     Id = 0x7FFE42
	DirtLongDistance        Id = 0x7FFF82
	TurfClass1                 = 0x7E0FFD
	DirtClass1                 = 0x7E0FFE
	TurfClass6                 = 0x7F07FD
	DirtClass6                 = 0x7F07FE
	TurfLargeNumberOfHorses    = 0x77FFFD
	TurfSmallNumberOfHorses    = 0x6FFFFD
	DirtLargeNumberOfHorses    = 0x77FFFE
	DirtSmallNumberOfHorses    = 0x6FFFFE
)

const (
	Turf                Id = 0x01
	Dirt                Id = 0x02
	ShortDistance1      Id = 0x04
	ShortDistance2      Id = 0x08
	ShortDistance3      Id = 0x10
	MiddleDistance1     Id = 0x20
	MiddleDistance2     Id = 0x40
	MiddleDistance3     Id = 0x80
	LongDistance        Id = 0x100
	TopJockey           Id = 0x200
	OtherJockey         Id = 0x400
	Class1              Id = 0x800
	Class2              Id = 0x1000
	Class3              Id = 0x2000
	Class4              Id = 0x4000
	Class5              Id = 0x8000
	Class6              Id = 0x10000
	GoodTrack           Id = 0x20000
	BadTrack            Id = 0x40000
	SmallNumberOfHorses Id = 0x80000
	LargeNumberOfHorses Id = 0x100000
	CentralCourse       Id = 0x200000
	LocalCourse         Id = 0x400000
	Class234            Id = 0x7000
	Class56             Id = 0x18000
	ShortDistance       Id = 0x1C
	MiddleDistance      Id = 0xE0
)

var filterIdMap = map[Id]string{
	All: "全レース",
	// 以下基本フィルタ
	Turf:                "芝",
	Dirt:                "ダート",
	ShortDistance1:      "1000~1200m",
	ShortDistance2:      "1201~1400m",
	ShortDistance3:      "1401~1600m",
	MiddleDistance1:     "1601~1700m",
	MiddleDistance2:     "1701~1800m",
	MiddleDistance3:     "1801~2000m",
	LongDistance:        "2001m~",
	TopJockey:           "上位騎手",
	OtherJockey:         "その他騎手",
	Class1:              "未勝利",
	Class2:              "1勝",
	Class3:              "2勝",
	Class4:              "3勝",
	Class5:              "OP・L",
	Class6:              "重賞",
	GoodTrack:           "良",
	BadTrack:            "稍重不",
	SmallNumberOfHorses: "少",
	LargeNumberOfHorses: "多",
	CentralCourse:       "中央",
	LocalCourse:         "ローカル",
	Class234:            "123勝",
	Class56:             "OP・重賞",
	ShortDistance:       "1000~1600m",
	MiddleDistance:      "1601~2000m",
	// 以下組み合わせフィルタ
	TurfShortDistance1:      "芝~1200m",
	TurfShortDistance2:      "芝~1400m",
	TurfShortDistance3:      "芝~1600m",
	TurfMiddleDistance1:     "芝~1800m",
	TurfMiddleDistance2:     "芝~2000m",
	TurfLongDistance:        "芝2001m~",
	DirtShortDistance1:      "ダ~1200m",
	DirtShortDistance2:      "ダ~1400m",
	DirtShortDistance3:      "ダ~1600m",
	DirtMiddleDistance1:     "ダ~1700m",
	DirtMiddleDistance2:     "ダ~1800m",
	DirtLongDistance:        "ダ1801m~",
	TurfClass1:              "芝未勝利",
	DirtClass1:              "ダ未勝利",
	TurfClass6:              "芝重賞",
	DirtClass6:              "ダ重賞",
	TurfLargeNumberOfHorses: "芝多頭数",
	TurfSmallNumberOfHorses: "芝少頭数",
	DirtLargeNumberOfHorses: "ダ多頭数",
	DirtSmallNumberOfHorses: "ダ少頭数",
}

func NewFilterId(rawId uint64, name string) Id {
	id := Id(rawId)
	filterIdMap[id] = name
	return id
}

func (i Id) Value() uint64 {
	return uint64(i)
}

func (i Id) String() string {
	id, _ := filterIdMap[i]
	return id
}
