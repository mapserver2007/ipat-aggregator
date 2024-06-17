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
	Turf2         Id = 0x100000000000
	Dirt2         Id = 0x80000000000
	Distance1000m Id = 0x40000000000
	Distance1150m Id = 0x20000000000
	Distance1200m Id = 0x10000000000
	Distance1300m Id = 0x8000000000
	Distance1400m Id = 0x4000000000
	Distance1500m Id = 0x2000000000
	Distance1600m Id = 0x1000000000
	Distance1700m Id = 0x800000000
	Distance1800m Id = 0x400000000
	Distance1900m Id = 0x200000000
	Distance2000m Id = 0x100000000
	Distance2100m Id = 0x80000000
	Distance2200m Id = 0x40000000
	Distance2300m Id = 0x20000000
	Distance2400m Id = 0x10000000
	Distance2500m Id = 0x8000000
	Distance2600m Id = 0x4000000
	Distance3000m Id = 0x2000000
	Distance3200m Id = 0x1000000
	Distance3400m Id = 0x800000
	Distance3600m Id = 0x400000
	Tokyo         Id = 0x200000
	Nakayama      Id = 0x100000
	Kyoto         Id = 0x80000
	Hanshin       Id = 0x40000
	Niigata       Id = 0x20000
	Chukyo        Id = 0x10000
	Sapporo       Id = 0x8000
	Hakodate      Id = 0x4000
	Fukushima     Id = 0x2000
	Kokura        Id = 0x1000
	GoodToFirm    Id = 0x800
	Good          Id = 0x400
	Yielding      Id = 0x200
	Soft          Id = 0x100
	Maiden        Id = 0x80
	OneWinClass   Id = 0x40
	TwoWinClass   Id = 0x20
	ThreeWinClass Id = 0x10
	OpenListed    Id = 0x8
	Grade3        Id = 0x4
	Grade2        Id = 0x2
	Grade1        Id = 0x1
)

const (
	All2               Id = 0x1FFFFFFFFFFF
	NiigataTurf1000m   Id = 0x140000020000
	HakodateTurf1000m  Id = 0x140000004000
	NakayamaTurf1200m  Id = 0x110000100000
	KyotoTurf1200m     Id = 0x110000080000
	HanshinTurf1200m   Id = 0x110000040000
	NiigataTurf1200m   Id = 0x110000020000
	ChukyoTurf1200m    Id = 0x110000010000
	SapporoTurf1200m   Id = 0x110000008000
	HakodateTurf1200m  Id = 0x110000004000
	FukushimaTurf1200m Id = 0x110000002000
	KokuraTurf1200m    Id = 0x110000001000
	TokyoTurf1400m     Id = 0x104000200000
	KyotoTurf1400m     Id = 0x104000080000
	HanshinTurf1400m   Id = 0x104000040000
	NiigataTurf1400m   Id = 0x104000020000
	ChukyoTurf1400m    Id = 0x104000010000
	SapporoTurf1500m   Id = 0x102000010000
	NakayamaTurf1600m  Id = 0x101000100000
	TokyoTurf1600m     Id = 0x101000200000
	KyotoTurf1600m     Id = 0x101000080000
	HanshinTurf1600m   Id = 0x101000040000
	ChukyoTurf1600m    Id = 0x101000010000
	NakayamaTurf1800m  Id = 0x100400100000
	TokyoTurf1800m     Id = 0x100400200000
	KyotoTurf1800m     Id = 0x100400080000
	HanshinTurf1800m   Id = 0x100400040000
	NiigataTurf1800m   Id = 0x100400020000
	SapporoTurf1800m   Id = 0x100400008000
	HakodateTurf1800m  Id = 0x100400004000
	FukushimaTurf1800m Id = 0x100400002000
	KokuraTurf1800m    Id = 0x100400001000
	NakayamaTurf2000m  Id = 0x100100100000
	TokyoTurf2000m     Id = 0x100100200000
	KyotoTurf2000m     Id = 0x100100080000
	NiigataTurf2000m   Id = 0x100100020000
	ChukyoTurf2000m    Id = 0x100100010000
	SapporoTurf2000m   Id = 0x100100008000
	HakodateTurf2000m  Id = 0x100100004000
	FukushimaTurf2000m Id = 0x100100002000
	KokuraTurf2000m    Id = 0x100100001000
	NakayamaTurf2200m  Id = 0x100040200000
	KyotoTurf2200m     Id = 0x100040080000
	HanshinTurf2200m   Id = 0x100040040000
	NiigataTurf2200m   Id = 0x100040020000
	ChukyoTurf2200m    Id = 0x100040010000
	TokyoTurf2300m     Id = 0x100020200000
	TokyoTurf2400m     Id = 0x100010200000
	KyotoTurf2400m     Id = 0x100010080000
	HanshinTurf2400m   Id = 0x100010040000
	NiigataTurf2400m   Id = 0x100010020000
	NakayamaTurf2500m  Id = 0x100008100000
	TokyoTurf2500m     Id = 0x100008200000
	HanshinTurf2600m   Id = 0x100004040000
	SapporoTurf2600m   Id = 0x100004008000
	HakodateTurf2600m  Id = 0x100004004000
	FukushimaTurf2600m Id = 0x100004002000
	KokuraTurf2600m    Id = 0x100004001000
	HanshinTurf3000m   Id = 0x100002040000
	ChukyoTurf3000m    Id = 0x100002010000
	KyotoTurf3200m     Id = 0x100001080000
	TokyoTurf3400m     Id = 0x100000A00000
	NakayamaTurf3600m  Id = 0x100000500000
	SapporoDirt1000m   Id = 0xC0000008000
	HakodateDirt1000m  Id = 0xC0000004000
	KokuraDirt1000m    Id = 0xC0000001000
	FukushimaDirt1150m Id = 0xA0000002000
	NakayamaDirt1200m  Id = 0x90000100000
	KyotoDirt1200m     Id = 0x90000080000
	NiigataDirt1200m   Id = 0x90000020000
	ChukyoDirt1200m    Id = 0x90000010000
	TokyoDirt1300m     Id = 0x88000200000
	TokyoDirt1400m     Id = 0x84000200000
	KyotoDirt1400m     Id = 0x84000080000
	HanshinDirt1400m   Id = 0x84000040000
	ChukyoDirt1400m    Id = 0x84000010000
	TokyoDirt1600m     Id = 0x81000200000
	SapporoDirt1700m   Id = 0x80800008000
	HakodateDirt1700m  Id = 0x80800004000
	FukushimaDirt1700m Id = 0x80800002000
	KokuraDirt1700m    Id = 0x80800001000
	NakayamaDirt1800m  Id = 0x80400100000
	KyotoDirt1800m     Id = 0x80400080000
	HanshinDirt1800m   Id = 0x80400040000
	NiigataDirt1800m   Id = 0x80400020000
	ChukyoDirt1800m    Id = 0x80400010000
	KyotoDirt1900m     Id = 0x80200080000
	ChukyoDirt1900m    Id = 0x80200010000
	HanshinDirt2000m   Id = 0x80100040000
	TokyoDirt2100m     Id = 0x80080200000
	NakayamaDirt2400m  Id = 0x80010100000
	SapporoDirt2400m   Id = 0x80010008000
	HakodateDirt2400m  Id = 0x80010004000
	FukushimaDirt2400m Id = 0x80010002000
	KokuraDirt2400m    Id = 0x80010001000
	NakayamaDirt2500m  Id = 0x80008100000
	NiigataDirt2500m   Id = 0x80008020000
)

var filterIdMap = map[Id]string{
	All: "全レース",
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
	Turf2:              "芝",
	Dirt2:              "ダート",
	Tokyo:              "東京",
	Nakayama:           "中山",
	Kyoto:              "京都",
	Hanshin:            "阪神",
	Niigata:            "新潟",
	Chukyo:             "中京",
	Sapporo:            "札幌",
	Hakodate:           "函館",
	Fukushima:          "福島",
	Kokura:             "小倉",
	GoodToFirm:         "良",
	Good:               "稍重",
	Yielding:           "重",
	Soft:               "不良",
	Maiden:             "未勝利",
	OneWinClass:        "1勝クラス",
	TwoWinClass:        "2勝クラス",
	ThreeWinClass:      "3勝クラス",
	OpenListed:         "OP・L",
	Grade3:             "G3",
	Grade2:             "G2",
	Grade1:             "G1",
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
