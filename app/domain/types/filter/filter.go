package filter

type Id uint64

const (
	All                                       Id = 0x7FFFFF // 全検索条件に引っ掛けるためのフィルタ
	TurfShortDistance1                        Id = 0x7FFE05
	TurfShortDistance2                        Id = 0x7FFE09
	TurfShortDistance3                        Id = 0x7FFE11
	TurfMiddleDistance1                       Id = 0x7FFE61
	TurfMiddleDistance2                       Id = 0x7FFE81
	TurfLongDistance                          Id = 0x7FFF01
	DirtShortDistance1                        Id = 0x7FFE06
	DirtShortDistance2                        Id = 0x7FFE0A
	DirtShortDistance3                        Id = 0x7FFE12
	DirtMiddleDistance1                       Id = 0x7FFE22
	DirtMiddleDistance2                       Id = 0x7FFE42
	DirtLongDistance                          Id = 0x7FFF82
	GoodTrackTurfShortDistance1CentralCourse  Id = 0x33FE05
	GoodTrackTurfShortDistance2CentralCourse  Id = 0x33FE09
	GoodTrackTurfShortDistance3CentralCourse  Id = 0x33FE11
	GoodTrackTurfMiddleDistance1CentralCourse Id = 0x33FE61
	GoodTrackTurfMiddleDistance2CentralCourse Id = 0x33FE81
	GoodTrackTurfLongDistanceCentralCourse    Id = 0x33FF01
	GoodTrackDirtShortDistance1CentralCourse  Id = 0x33FE06
	GoodTrackDirtShortDistance2CentralCourse  Id = 0x33FE0A
	GoodTrackDirtShortDistance3CentralCourse  Id = 0x33FE12
	GoodTrackDirtMiddleDistance2CentralCourse Id = 0x33FE62
	GoodTrackDirtLongDistanceCentralCourse    Id = 0x33FF82
	TurfClass1                                Id = 0x320FFD
	TurfClass2                                Id = 0x3217FD
	TurfClass3                                Id = 0x3227FD
	TurfClass4                                Id = 0x3247FD
	TurfClass5                                Id = 0x3287FD
	TurfClass6                                Id = 0x3307FD
	DirtClass1                                Id = 0x320FFE
	DirtClass2                                Id = 0x3217FE
	DirtClass3                                Id = 0x3227FE
	DirtClass4                                Id = 0x3247FE
	DirtClass5                                Id = 0x7287FE
	DirtClass6                                Id = 0x7307FE
	DirtBadConditionClass1                    Id = 0x340FFE
	DirtBadConditionClass2                    Id = 0x3417FE
	DirtBadConditionClass3                    Id = 0x3427FE
	DirtBadConditionClass4                    Id = 0x3447FE
	DirtBadConditionClass5                    Id = 0x7487FE
	DirtBadConditionClass6                    Id = 0x7507FE
	// 予想専用
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
	All: "条件なし",
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
	TurfShortDistance1:                        "芝~1200m",
	TurfShortDistance2:                        "芝~1400m",
	TurfShortDistance3:                        "芝~1600m",
	TurfMiddleDistance1:                       "芝~1800m",
	TurfMiddleDistance2:                       "芝~2000m",
	TurfLongDistance:                          "芝2001m~",
	DirtShortDistance1:                        "ダ~1200m",
	DirtShortDistance2:                        "ダ~1400m",
	DirtShortDistance3:                        "ダ~1600m",
	DirtMiddleDistance1:                       "ダ~1700m",
	DirtMiddleDistance2:                       "ダ~1800m",
	DirtLongDistance:                          "ダ1801m~",
	GoodTrackTurfShortDistance1CentralCourse:  "芝良多中央~1200m",
	GoodTrackTurfShortDistance2CentralCourse:  "芝良多中央~1400m",
	GoodTrackTurfShortDistance3CentralCourse:  "芝良多中央~1600m",
	GoodTrackTurfMiddleDistance1CentralCourse: "芝良多中央~1800m",
	GoodTrackTurfMiddleDistance2CentralCourse: "芝良多中央~2000m",
	GoodTrackTurfLongDistanceCentralCourse:    "芝良多中央2001m~",
	GoodTrackDirtShortDistance1CentralCourse:  "ダ良多中央~1200m",
	GoodTrackDirtShortDistance2CentralCourse:  "ダ良多中央~1400m",
	GoodTrackDirtShortDistance3CentralCourse:  "ダ良多中央~1600m",
	GoodTrackDirtMiddleDistance2CentralCourse: "ダ良多中央~1800m",
	GoodTrackDirtLongDistanceCentralCourse:    "ダ良多中央1801m~",
	TurfClass1:                                "芝良多中央未勝利",
	TurfClass2:                                "芝良多中央1勝",
	TurfClass3:                                "芝良多中央2勝",
	TurfClass4:                                "芝良多中央3勝",
	TurfClass5:                                "芝良多中央OP・L",
	TurfClass6:                                "芝良多中央重賞",
	DirtClass1:                                "ダ良多中央未勝利",
	DirtClass2:                                "ダ良多中央1勝",
	DirtClass3:                                "ダ良多中央2勝",
	DirtClass4:                                "ダ良多中央3勝",
	DirtClass5:                                "ダ良多全場OP・L",
	DirtClass6:                                "ダ良多全場重賞",
	DirtBadConditionClass1:                    "ダ稍重不多中央未勝利",
	DirtBadConditionClass2:                    "ダ稍重不多中央1勝",
	DirtBadConditionClass3:                    "ダ稍重不多中央2勝",
	DirtBadConditionClass4:                    "ダ稍重不多中央3勝",
	DirtBadConditionClass5:                    "ダ稍重不多全場OP・L",
	DirtBadConditionClass6:                    "ダ稍重不多全場重賞",
	// 予想専用
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
