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
	// 以下距離フィルタ
	Turf2         Id = 0x100000000
	Dirt2         Id = 0x80000000
	Distance1000m Id = 0x40000000
	Distance1150m Id = 0x20000000
	Distance1200m Id = 0x10000000
	Distance1300m Id = 0x8000000
	Distance1400m Id = 0x4000000
	Distance1500m Id = 0x2000000
	Distance1600m Id = 0x1000000
	Distance1700m Id = 0x800000
	Distance1800m Id = 0x400000
	Distance1900m Id = 0x200000
	Distance2000m Id = 0x100000
	Distance2100m Id = 0x80000
	Distance2200m Id = 0x40000
	Distance2300m Id = 0x20000
	Distance2400m Id = 0x10000
	Distance2500m Id = 0x8000
	Distance2600m Id = 0x4000
	Distance3000m Id = 0x2000
	Distance3200m Id = 0x1000
	Distance3400m Id = 0x800
	Distance3600m Id = 0x400
	Tokyo         Id = 0x200
	Nakayama      Id = 0x100
	Kyoto         Id = 0x80
	Hanshin       Id = 0x40
	Niigata       Id = 0x20
	Chukyo        Id = 0x10
	Sapporo       Id = 0x8
	Hakodate      Id = 0x4
	Fukushima     Id = 0x2
	Kokura        Id = 0x1
)

const (
	All2               Id = 0x1FFFFFFFF
	NiigataTurf1000m   Id = 0x140000020
	HakodateTurf1000m  Id = 0x140000004
	NakayamaTurf1200m  Id = 0x110000100
	KyotoTurf1200m     Id = 0x110000080
	HanshinTurf1200m   Id = 0x110000040
	NiigataTurf1200m   Id = 0x110000020
	ChukyoTurf1200m    Id = 0x110000010
	SapporoTurf1200m   Id = 0x110000008
	HakodateTurf1200m  Id = 0x110000004
	FukushimaTurf1200m Id = 0x110000002
	KokuraTurf1200m    Id = 0x110000001
	TokyoTurf1400m     Id = 0x104000200
	KyotoTurf1400m     Id = 0x104000080
	HanshinTurf1400m   Id = 0x104000040
	NiigataTurf1400m   Id = 0x104000020
	ChukyoTurf1400m    Id = 0x104000010
	SapporoTurf1500m   Id = 0x102000010
	NakayamaTurf1600m  Id = 0x101000100
	TokyoTurf1600m     Id = 0x101000200
	KyotoTurf1600m     Id = 0x101000080
	HanshinTurf1600m   Id = 0x101000040
	ChukyoTurf1600m    Id = 0x101000010
	NakayamaTurf1800m  Id = 0x100400100
	TokyoTurf1800m     Id = 0x100400200
	KyotoTurf1800m     Id = 0x100400080
	HanshinTurf1800m   Id = 0x100400040
	NiigataTurf1800m   Id = 0x100400020
	SapporoTurf1800m   Id = 0x100400008
	HakodateTurf1800m  Id = 0x100400004
	FukushimaTurf1800m Id = 0x100400002
	KokuraTurf1800m    Id = 0x100400001
	NakayamaTurf2000m  Id = 0x100100100
	TokyoTurf2000m     Id = 0x100100200
	KyotoTurf2000m     Id = 0x100100080
	NiigataTurf2000m   Id = 0x100100020
	ChukyoTurf2000m    Id = 0x100100010
	SapporoTurf2000m   Id = 0x100100008
	HakodateTurf2000m  Id = 0x100100004
	FukushimaTurf2000m Id = 0x100100002
	KokuraTurf2000m    Id = 0x100100001
	NakayamaTurf2200m  Id = 0x100040200
	KyotoTurf2200m     Id = 0x100040080
	HanshinTurf2200m   Id = 0x100040040
	NiigataTurf2200m   Id = 0x100040020
	ChukyoTurf2200m    Id = 0x100040010
	TokyoTurf2300m     Id = 0x100020200
	TokyoTurf2400m     Id = 0x100010200
	KyotoTurf2400m     Id = 0x100010080
	HanshinTurf2400m   Id = 0x100010040
	NiigataTurf2400m   Id = 0x100010020
	NakayamaTurf2500m  Id = 0x100008100
	TokyoTurf2500m     Id = 0x100008200
	HanshinTurf2600m   Id = 0x100004040
	SapporoTurf2600m   Id = 0x100004008
	HakodateTurf2600m  Id = 0x100004004
	FukushimaTurf2600m Id = 0x100004002
	KokuraTurf2600m    Id = 0x100004001
	HanshinTurf3000m   Id = 0x100002040
	ChukyoTurf3000m    Id = 0x100002010
	KyotoTurf3200m     Id = 0x100001080
	TokyoTurf3400m     Id = 0x100000A00
	NakayamaTurf3600m  Id = 0x100000500
	SapporoDirt1000m   Id = 0xC0000008
	HakodateDirt1000m  Id = 0xC0000004
	KokuraDirt1000m    Id = 0xC0000001
	FukushimaDirt1150m Id = 0xA0000002
	NakayamaDirt1200m  Id = 0x90000100
	KyotoDirt1200m     Id = 0x90000080
	NiigataDirt1200m   Id = 0x90000020
	ChukyoDirt1200m    Id = 0x90000010
	TokyoDirt1300m     Id = 0x88000200
	TokyoDirt1400m     Id = 0x84000200
	KyotoDirt1400m     Id = 0x84000080
	HanshinDirt1400m   Id = 0x84000040
	ChukyoDirt1400m    Id = 0x84000010
	TokyoDirt1600m     Id = 0x81000200
	SapporoDirt1700m   Id = 0x80800008
	HakodateDirt1700m  Id = 0x80800004
	FukushimaDirt1700m Id = 0x80800002
	KokuraDirt1700m    Id = 0x80800001
	NakayamaDirt1800m  Id = 0x80400100
	KyotoDirt1800m     Id = 0x80400080
	HanshinDirt1800m   Id = 0x80400040
	NiigataDirt1800m   Id = 0x80400020
	ChukyoDirt1800m    Id = 0x80400010
	KyotoDirt1900m     Id = 0x80200080
	ChukyoDirt1900m    Id = 0x80200010
	HanshinDirt2000m   Id = 0x80100040
	TokyoDirt2100m     Id = 0x80080200
	NakayamaDirt2400m  Id = 0x80010100
	SapporoDirt2400m   Id = 0x80010008
	HakodateDirt2400m  Id = 0x80010004
	FukushimaDirt2400m Id = 0x80010002
	KokuraDirt2400m    Id = 0x80010001
	NakayamaDirt2500m  Id = 0x80008100
	NiigataDirt2500m   Id = 0x80008020
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
	// 以下距離・場所フィルタ
	All2:               "全レース",
	NiigataTurf1000m:   "新潟芝1000m",
	HakodateTurf1000m:  "函館芝1000m",
	NakayamaTurf1200m:  "中山芝1200m",
	KyotoTurf1200m:     "京都芝1200m",
	HanshinTurf1200m:   "阪神芝1200m",
	NiigataTurf1200m:   "新潟芝1200m",
	ChukyoTurf1200m:    "中京芝1200m",
	SapporoTurf1200m:   "札幌芝1200m",
	HakodateTurf1200m:  "函館芝1200m",
	FukushimaTurf1200m: "福島芝1200m",
	KokuraTurf1200m:    "小倉芝1200m",
	TokyoTurf1400m:     "東京芝1400m",
	KyotoTurf1400m:     "京都芝1400m",
	HanshinTurf1400m:   "阪神芝1400m",
	NiigataTurf1400m:   "新潟芝1400m",
	ChukyoTurf1400m:    "中京芝1400m",
	SapporoTurf1500m:   "札幌芝1500m",
	NakayamaTurf1600m:  "中山芝1600m",
	TokyoTurf1600m:     "東京芝1600m",
	KyotoTurf1600m:     "京都芝1600m",
	HanshinTurf1600m:   "阪神芝1600m",
	ChukyoTurf1600m:    "中京芝1600m",
	NakayamaTurf1800m:  "中山芝1800m",
	TokyoTurf1800m:     "東京芝1800m",
	KyotoTurf1800m:     "京都芝1800m",
	HanshinTurf1800m:   "阪神芝1800m",
	NiigataTurf1800m:   "新潟芝1800m",
	SapporoTurf1800m:   "札幌芝1800m",
	HakodateTurf1800m:  "函館芝1800m",
	FukushimaTurf1800m: "福島芝1800m",
	KokuraTurf1800m:    "小倉芝1800m",
	NakayamaTurf2000m:  "中山芝2000m",
	TokyoTurf2000m:     "東京芝2000m",
	KyotoTurf2000m:     "京都芝2000m",
	NiigataTurf2000m:   "新潟芝2000m",
	ChukyoTurf2000m:    "中京芝2000m",
	SapporoTurf2000m:   "札幌芝2000m",
	HakodateTurf2000m:  "函館芝2000m",
	FukushimaTurf2000m: "福島芝2000m",
	KokuraTurf2000m:    "小倉芝2000m",
	NakayamaTurf2200m:  "中山芝2200m",
	KyotoTurf2200m:     "京都芝2200m",
	HanshinTurf2200m:   "阪神芝2200m",
	NiigataTurf2200m:   "新潟芝2200m",
	ChukyoTurf2200m:    "中京芝2200m",
	TokyoTurf2300m:     "東京芝2300m",
	TokyoTurf2400m:     "東京芝2400m",
	KyotoTurf2400m:     "京都芝2400m",
	HanshinTurf2400m:   "阪神芝2400m",
	NiigataTurf2400m:   "新潟芝2400m",
	NakayamaTurf2500m:  "中山芝2500m",
	TokyoTurf2500m:     "東京芝2500m",
	HanshinTurf2600m:   "阪神芝2600m",
	SapporoTurf2600m:   "札幌芝2600m",
	HakodateTurf2600m:  "函館芝2600m",
	FukushimaTurf2600m: "福島芝2600m",
	KokuraTurf2600m:    "小倉芝2600m",
	HanshinTurf3000m:   "阪神芝3000m",
	ChukyoTurf3000m:    "中京芝3000m",
	KyotoTurf3200m:     "京都芝3200m",
	TokyoTurf3400m:     "東京芝3400m",
	NakayamaTurf3600m:  "中山芝3600m",
	SapporoDirt1000m:   "札幌ダ1000m",
	HakodateDirt1000m:  "函館ダ1000m",
	KokuraDirt1000m:    "小倉ダ1000m",
	FukushimaDirt1150m: "福島ダ1150m",
	NakayamaDirt1200m:  "中山ダ1200m",
	KyotoDirt1200m:     "京都ダ1200m",
	NiigataDirt1200m:   "新潟ダ1200m",
	ChukyoDirt1200m:    "中京ダ1200m",
	TokyoDirt1300m:     "東京ダ1300m",
	TokyoDirt1400m:     "東京ダ1400m",
	KyotoDirt1400m:     "京都ダ1400m",
	HanshinDirt1400m:   "阪神ダ1400m",
	ChukyoDirt1400m:    "中京ダ1400m",
	TokyoDirt1600m:     "東京ダ1600m",
	SapporoDirt1700m:   "札幌ダ1700m",
	HakodateDirt1700m:  "函館ダ1700m",
	FukushimaDirt1700m: "福島ダ1700m",
	KokuraDirt1700m:    "小倉ダ1700m",
	NakayamaDirt1800m:  "中山ダ1800m",
	KyotoDirt1800m:     "京都ダ1800m",
	HanshinDirt1800m:   "阪神ダ1800m",
	NiigataDirt1800m:   "新潟ダ1800m",
	ChukyoDirt1800m:    "中京ダ1800m",
	KyotoDirt1900m:     "京都ダ1900m",
	ChukyoDirt1900m:    "中京ダ1900m",
	HanshinDirt2000m:   "阪神ダ2000m",
	TokyoDirt2100m:     "東京ダ2100m",
	NakayamaDirt2400m:  "中山ダ2400m",
	SapporoDirt2400m:   "札幌ダ2400m",
	HakodateDirt2400m:  "函館ダ2400m",
	FukushimaDirt2400m: "福島ダ2400m",
	KokuraDirt2400m:    "小倉ダ2400m",
	NakayamaDirt2500m:  "中山ダ2500m",
	NiigataDirt2500m:   "新潟ダ2500m",
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
