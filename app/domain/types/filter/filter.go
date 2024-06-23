package filter

import "sort"

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
	All2          Id = 0x1FFFFFFFFF
	Turf2         Id = 0x1000000000
	Dirt2         Id = 0x800000000
	Distance1000m Id = 0x400000000
	Distance1150m Id = 0x200000000
	Distance1200m Id = 0x100000000
	Distance1300m Id = 0x80000000
	Distance1400m Id = 0x40000000
	Distance1500m Id = 0x20000000
	Distance1600m Id = 0x10000000
	Distance1700m Id = 0x8000000
	Distance1800m Id = 0x4000000
	Distance1900m Id = 0x2000000
	Distance2000m Id = 0x1000000
	Distance2100m Id = 0x800000
	Distance2200m Id = 0x400000
	Distance2300m Id = 0x200000
	Distance2400m Id = 0x100000
	Distance2500m Id = 0x80000
	Distance2600m Id = 0x40000
	Distance3000m Id = 0x20000
	Distance3200m Id = 0x10000
	Distance3400m Id = 0x8000
	Distance3600m Id = 0x4000
	Tokyo         Id = 0x2000
	Nakayama      Id = 0x1000
	Kyoto         Id = 0x800
	Hanshin       Id = 0x400
	Niigata       Id = 0x200
	Chukyo        Id = 0x100
	Sapporo       Id = 0x80
	Hakodate      Id = 0x40
	Fukushima     Id = 0x20
	Kokura        Id = 0x10
	GoodToFirm    Id = 0x8
	Good          Id = 0x4
	Yielding      Id = 0x2
	Soft          Id = 0x1
)

const (
	TurfAll                      Id = 0x17FFFFFFFF
	DirtAll                      Id = 0xFFFFFFFFF
	TurfGoodToFirm               Id = 0x17FFFFFFF8
	TurfGood                     Id = 0x17FFFFFFF4
	TurfYielding                 Id = 0x17FFFFFFF2
	TurfSoft                     Id = 0x17FFFFFFF1
	DirtGoodToFirm               Id = 0xFFFFFFFF8
	DirtGood                     Id = 0xFFFFFFFF4
	DirtYielding                 Id = 0xFFFFFFFF2
	DirtSoft                     Id = 0xFFFFFFFF1
	NiigataTurf1000m             Id = 0x140000020F
	NiigataGoodToFirmTurf1000m   Id = 0x1400000208
	NiigataGoodTurf1000m         Id = 0x1400000204
	NiigataYieldingTurf1000m     Id = 0x1400000202
	NiigataSoftTurf1000m         Id = 0x1400000201
	HakodateTurf1000m            Id = 0x140000004F
	HakodateGoodToFirmTurf1000m  Id = 0x1400000048
	HakodateGoodTurf1000m        Id = 0x1400000044
	HakodateYieldingTurf1000m    Id = 0x1400000042
	HakodateSoftTurf1000m        Id = 0x1400000041
	NakayamaTurf1200m            Id = 0x110000100F
	NakayamaGoodToFirmTurf1200m  Id = 0x1100001008
	NakayamaGoodTurf1200m        Id = 0x1100001004
	NakayamaYieldingTurf1200m    Id = 0x1100001002
	NakayamaSoftTurf1200m        Id = 0x1100001001
	KyotoTurf1200m               Id = 0x110000080F
	KyotoGoodToFirmTurf1200m     Id = 0x1100000808
	KyotoGoodTurf1200m           Id = 0x1100000804
	KyotoYieldingTurf1200m       Id = 0x1100000802
	KyotoSoftTurf1200m           Id = 0x1100000801
	HanshinTurf1200m             Id = 0x110000040F
	HanshinGoodToFirmTurf1200m   Id = 0x1100000408
	HanshinGoodTurf1200m         Id = 0x1100000404
	HanshinYieldingTurf1200m     Id = 0x1100000402
	HanshinSoftTurf1200m         Id = 0x1100000401
	NiigataTurf1200m             Id = 0x110000020F
	NiigataGoodToFirmTurf1200m   Id = 0x1100000208
	NiigataGoodTurf1200m         Id = 0x1100000204
	NiigataYieldingTurf1200m     Id = 0x1100000202
	NiigataSoftTurf1200m         Id = 0x1100000201
	ChukyoTurf1200m              Id = 0x110000010F
	ChukyoGoodToFirmTurf1200m    Id = 0x1100000108
	ChukyoGoodTurf1200m          Id = 0x1100000104
	ChukyoYieldingTurf1200m      Id = 0x1100000102
	ChukyoSoftTurf1200m          Id = 0x1100000101
	SapporoTurf1200m             Id = 0x110000008F
	SapporoGoodToFirmTurf1200m   Id = 0x1100000088
	SapporoGoodTurf1200m         Id = 0x1100000084
	SapporoYieldingTurf1200m     Id = 0x1100000082
	SapporoSoftTurf1200m         Id = 0x1100000081
	HakodateTurf1200m            Id = 0x110000004F
	HakodateGoodToFirmTurf1200m  Id = 0x1100000048
	HakodateGoodTurf1200m        Id = 0x1100000044
	HakodateYieldingTurf1200m    Id = 0x1100000042
	HakodateSoftTurf1200m        Id = 0x1100000041
	FukushimaTurf1200m           Id = 0x110000002F
	FukushimaGoodToFirmTurf1200m Id = 0x1100000028
	FukushimaGoodTurf1200m       Id = 0x1100000024
	FukushimaYieldingTurf1200m   Id = 0x1100000022
	FukushimaSoftTurf1200m       Id = 0x1100000021
	KokuraTurf1200m              Id = 0x110000001F
	KokuraGoodToFirmTurf1200m    Id = 0x1100000018
	KokuraGoodTurf1200m          Id = 0x1100000014
	KokuraYieldingTurf1200m      Id = 0x1100000012
	KokuraSoftTurf1200m          Id = 0x1100000011
	TokyoTurf1400m               Id = 0x104000200F
	TokyoGoodToFirmTurf1400m     Id = 0x1040002008
	TokyoGoodTurf1400m           Id = 0x1040002004
	TokyoYieldingTurf1400m       Id = 0x1040002002
	TokyoSoftTurf1400m           Id = 0x1040002001
	KyotoTurf1400m               Id = 0x104000080F
	KyotoGoodToFirmTurf1400m     Id = 0x1040000808
	KyotoGoodTurf1400m           Id = 0x1040000804
	KyotoYieldingTurf1400m       Id = 0x1040000802
	KyotoSoftTurf1400m           Id = 0x1040000801
	HanshinTurf1400m             Id = 0x104000040F
	HanshinGoodToFirmTurf1400m   Id = 0x1040000408
	HanshinGoodTurf1400m         Id = 0x1040000404
	HanshinYieldingTurf1400m     Id = 0x1040000402
	HanshinSoftTurf1400m         Id = 0x1040000401
	NiigataTurf1400m             Id = 0x104000020F
	NiigataGoodToFirmTurf1400m   Id = 0x1040000208
	NiigataGoodTurf1400m         Id = 0x1040000204
	NiigataYieldingTurf1400m     Id = 0x1040000202
	NiigataSoftTurf1400m         Id = 0x1040000201
	ChukyoTurf1400m              Id = 0x104000010F
	ChukyoGoodToFirmTurf1400m    Id = 0x1040000108
	ChukyoGoodTurf1400m          Id = 0x1040000104
	ChukyoYieldingTurf1400m      Id = 0x1040000102
	ChukyoSoftTurf1400m          Id = 0x1040000101
	SapporoTurf1500m             Id = 0x102000010F
	SapporoGoodToFirmTurf1500m   Id = 0x1020000108
	SapporoGoodTurf1500m         Id = 0x1020000104
	SapporoYieldingTurf1500m     Id = 0x1020000102
	SapporoSoftTurf1500m         Id = 0x1020000101
	NakayamaTurf1600m            Id = 0x101000100F
	NakayamaGoodToFirmTurf1600m  Id = 0x1010001008
	NakayamaGoodTurf1600m        Id = 0x1010001004
	NakayamaYieldingTurf1600m    Id = 0x1010001002
	NakayamaSoftTurf1600m        Id = 0x1010001001
	TokyoTurf1600m               Id = 0x101000200F
	TokyoGoodToFirmTurf1600m     Id = 0x1010002008
	TokyoGoodTurf1600m           Id = 0x1010002004
	TokyoYieldingTurf1600m       Id = 0x1010002002
	TokyoSoftTurf1600m           Id = 0x1010002001
	KyotoTurf1600m               Id = 0x101000080F
	KyotoGoodToFirmTurf1600m     Id = 0x1010000808
	KyotoGoodTurf1600m           Id = 0x1010000804
	KyotoYieldingTurf1600m       Id = 0x1010000802
	KyotoSoftTurf1600m           Id = 0x1010000801
	HanshinTurf1600m             Id = 0x101000040F
	HanshinGoodToFirmTurf1600m   Id = 0x1010000408
	HanshinGoodTurf1600m         Id = 0x1010000404
	HanshinYieldingTurf1600m     Id = 0x1010000402
	HanshinSoftTurf1600m         Id = 0x1010000401
	ChukyoTurf1600m              Id = 0x101000010F
	ChukyoGoodToFirmTurf1600m    Id = 0x1010000108
	ChukyoGoodTurf1600m          Id = 0x1010000104
	ChukyoYieldingTurf1600m      Id = 0x1010000102
	ChukyoSoftTurf1600m          Id = 0x1010000101
	NakayamaTurf1800m            Id = 0x100400100F
	NakayamaGoodToFirmTurf1800m  Id = 0x1004001008
	NakayamaGoodTurf1800m        Id = 0x1004001004
	NakayamaYieldingTurf1800m    Id = 0x1004001002
	NakayamaSoftTurf1800m        Id = 0x1004001001
	TokyoTurf1800m               Id = 0x100400200F
	TokyoGoodToFirmTurf1800m     Id = 0x1004002008
	TokyoGoodTurf1800m           Id = 0x1004002004
	TokyoYieldingTurf1800m       Id = 0x1004002002
	TokyoSoftTurf1800m           Id = 0x1004002001
	KyotoTurf1800m               Id = 0x100400080F
	KyotoGoodToFirmTurf1800m     Id = 0x1004000808
	KyotoGoodTurf1800m           Id = 0x1004000804
	KyotoYieldingTurf1800m       Id = 0x1004000802
	KyotoSoftTurf1800m           Id = 0x1004000801
	HanshinTurf1800m             Id = 0x100400040F
	HanshinGoodToFirmTurf1800m   Id = 0x1004000408
	HanshinGoodTurf1800m         Id = 0x1004000404
	HanshinYieldingTurf1800m     Id = 0x1004000402
	HanshinSoftTurf1800m         Id = 0x1004000401
	NiigataTurf1800m             Id = 0x100400020F
	NiigataGoodToFirmTurf1800m   Id = 0x1004000208
	NiigataGoodTurf1800m         Id = 0x1004000204
	NiigataYieldingTurf1800m     Id = 0x1004000202
	NiigataSoftTurf1800m         Id = 0x1004000201
	SapporoTurf1800m             Id = 0x100400008F
	SapporoGoodToFirmTurf1800m   Id = 0x1004000088
	SapporoGoodTurf1800m         Id = 0x1004000084
	SapporoYieldingTurf1800m     Id = 0x1004000082
	SapporoSoftTurf1800m         Id = 0x1004000081
	HakodateTurf1800m            Id = 0x100400004F
	HakodateGoodToFirmTurf1800m  Id = 0x1004000048
	HakodateGoodTurf1800m        Id = 0x1004000044
	HakodateYieldingTurf1800m    Id = 0x1004000042
	HakodateSoftTurf1800m        Id = 0x1004000041
	FukushimaTurf1800m           Id = 0x100400002F
	FukushimaGoodToFirmTurf1800m Id = 0x1004000028
	FukushimaGoodTurf1800m       Id = 0x1004000024
	FukushimaYieldingTurf1800m   Id = 0x1004000022
	FukushimaSoftTurf1800m       Id = 0x1004000021
	KokuraTurf1800m              Id = 0x100400001F
	KokuraGoodToFirmTurf1800m    Id = 0x1004000018
	KokuraGoodTurf1800m          Id = 0x1004000014
	KokuraYieldingTurf1800m      Id = 0x1004000012
	KokuraSoftTurf1800m          Id = 0x1004000011
	NakayamaTurf2000m            Id = 0x100100100F
	NakayamaGoodToFirmTurf2000m  Id = 0x1001001008
	NakayamaGoodTurf2000m        Id = 0x1001001004
	NakayamaYieldingTurf2000m    Id = 0x1001001002
	NakayamaSoftTurf2000m        Id = 0x1001001001
	TokyoTurf2000m               Id = 0x100100200F
	TokyoGoodToFirmTurf2000m     Id = 0x1001002008
	TokyoGoodTurf2000m           Id = 0x1001002004
	TokyoYieldingTurf2000m       Id = 0x1001002002
	TokyoSoftTurf2000m           Id = 0x1001002001
	KyotoTurf2000m               Id = 0x100100080F
	KyotoGoodToFirmTurf2000m     Id = 0x1001000808
	KyotoGoodTurf2000m           Id = 0x1001000804
	KyotoYieldingTurf2000m       Id = 0x1001000802
	KyotoSoftTurf2000m           Id = 0x1001000801
	NiigataTurf2000m             Id = 0x100100020F
	NiigataGoodToFirmTurf2000m   Id = 0x1001000208
	NiigataGoodTurf2000m         Id = 0x1001000204
	NiigataYieldingTurf2000m     Id = 0x1001000202
	NiigataSoftTurf2000m         Id = 0x1001000201
	ChukyoTurf2000m              Id = 0x100100010F
	ChukyoGoodToFirmTurf2000m    Id = 0x1001000108
	ChukyoGoodTurf2000m          Id = 0x1001000104
	ChukyoYieldingTurf2000m      Id = 0x1001000102
	ChukyoSoftTurf2000m          Id = 0x1001000101
	SapporoTurf2000m             Id = 0x100100008F
	SapporoGoodToFirmTurf2000m   Id = 0x1001000088
	SapporoGoodTurf2000m         Id = 0x1001000084
	SapporoYieldingTurf2000m     Id = 0x1001000082
	SapporoSoftTurf2000m         Id = 0x1001000081
	HakodateTurf2000m            Id = 0x100100004F
	HakodateGoodToFirmTurf2000m  Id = 0x1001000048
	HakodateGoodTurf2000m        Id = 0x1001000044
	HakodateYieldingTurf2000m    Id = 0x1001000042
	HakodateSoftTurf2000m        Id = 0x1001000041
	FukushimaTurf2000m           Id = 0x100100002F
	FukushimaGoodToFirmTurf2000m Id = 0x1001000028
	FukushimaGoodTurf2000m       Id = 0x1001000024
	FukushimaYieldingTurf2000m   Id = 0x1001000022
	FukushimaSoftTurf2000m       Id = 0x1001000021
	KokuraTurf2000m              Id = 0x100100001F
	KokuraGoodToFirmTurf2000m    Id = 0x1001000018
	KokuraGoodTurf2000m          Id = 0x1001000014
	KokuraYieldingTurf2000m      Id = 0x1001000012
	KokuraSoftTurf2000m          Id = 0x1001000011
	NakayamaTurf2200m            Id = 0x100040100F
	NakayamaGoodToFirmTurf2200m  Id = 0x1000401008
	NakayamaGoodTurf2200m        Id = 0x1000401004
	NakayamaYieldingTurf2200m    Id = 0x1000401002
	NakayamaSoftTurf2200m        Id = 0x1000401001
	KyotoTurf2200m               Id = 0x100040080F
	KyotoGoodToFirmTurf2200m     Id = 0x1000400808
	KyotoGoodTurf2200m           Id = 0x1000400804
	KyotoYieldingTurf2200m       Id = 0x1000400802
	KyotoSoftTurf2200m           Id = 0x1000400801
	HanshinTurf2200m             Id = 0x100040040F
	HanshinGoodToFirmTurf2200m   Id = 0x1000400408
	HanshinGoodTurf2200m         Id = 0x1000400404
	HanshinYieldingTurf2200m     Id = 0x1000400402
	HanshinSoftTurf2200m         Id = 0x1000400401
	NiigataTurf2200m             Id = 0x100040020F
	NiigataGoodToFirmTurf2200m   Id = 0x1000400208
	NiigataGoodTurf2200m         Id = 0x1000400204
	NiigataYieldingTurf2200m     Id = 0x1000400202
	NiigataSoftTurf2200m         Id = 0x1000400201
	ChukyoTurf2200m              Id = 0x100040010F
	ChukyoGoodToFirmTurf2200m    Id = 0x1000400108
	ChukyoGoodTurf2200m          Id = 0x1000400104
	ChukyoYieldingTurf2200m      Id = 0x1000400102
	ChukyoSoftTurf2200m          Id = 0x1000400101
	TokyoTurf2300m               Id = 0x100020200F
	TokyoGoodToFirmTurf2300m     Id = 0x1000202008
	TokyoGoodTurf2300m           Id = 0x1000202004
	TokyoYieldingTurf2300m       Id = 0x1000202002
	TokyoSoftTurf2300m           Id = 0x1000202001
	TokyoTurf2400m               Id = 0x100010200F
	TokyoGoodToFirmTurf2400m     Id = 0x1000102008
	TokyoGoodTurf2400m           Id = 0x1000102004
	TokyoYieldingTurf2400m       Id = 0x1000102002
	TokyoSoftTurf2400m           Id = 0x1000102001
	KyotoTurf2400m               Id = 0x100010080F
	KyotoGoodToFirmTurf2400m     Id = 0x1000100808
	KyotoGoodTurf2400m           Id = 0x1000100804
	KyotoYieldingTurf2400m       Id = 0x1000100802
	KyotoSoftTurf2400m           Id = 0x1000100801
	HanshinTurf2400m             Id = 0x100010040F
	HanshinGoodToFirmTurf2400m   Id = 0x1000100408
	HanshinGoodTurf2400m         Id = 0x1000100404
	HanshinYieldingTurf2400m     Id = 0x1000100402
	HanshinSoftTurf2400m         Id = 0x1000100401
	NiigataTurf2400m             Id = 0x100010020F
	NiigataGoodToFirmTurf2400m   Id = 0x1000100208
	NiigataGoodTurf2400m         Id = 0x1000100204
	NiigataYieldingTurf2400m     Id = 0x1000100202
	NiigataSoftTurf2400m         Id = 0x1000100201
	NakayamaTurf2500m            Id = 0x100008100F
	NakayamaGoodToFirmTurf2500m  Id = 0x1000081008
	NakayamaGoodTurf2500m        Id = 0x1000081004
	NakayamaYieldingTurf2500m    Id = 0x1000081002
	NakayamaSoftTurf2500m        Id = 0x1000081001
	TokyoTurf2500m               Id = 0x100008200F
	TokyoGoodToFirmTurf2500m     Id = 0x1000082008
	TokyoGoodTurf2500m           Id = 0x1000082004
	TokyoYieldingTurf2500m       Id = 0x1000082002
	TokyoSoftTurf2500m           Id = 0x1000082001
	HanshinTurf2600m             Id = 0x100004040F
	HanshinGoodToFirmTurf2600m   Id = 0x1000040408
	HanshinGoodTurf2600m         Id = 0x1000040404
	HanshinYieldingTurf2600m     Id = 0x1000040402
	HanshinSoftTurf2600m         Id = 0x1000040401
	SapporoTurf2600m             Id = 0x100004008F
	SapporoGoodToFirmTurf2600m   Id = 0x1000040088
	SapporoGoodTurf2600m         Id = 0x1000040084
	SapporoYieldingTurf2600m     Id = 0x1000040082
	SapporoSoftTurf2600m         Id = 0x1000040081
	HakodateTurf2600m            Id = 0x100004004F
	HakodateGoodToFirmTurf2600m  Id = 0x1000040048
	HakodateGoodTurf2600m        Id = 0x1000040044
	HakodateYieldingTurf2600m    Id = 0x1000040042
	HakodateSoftTurf2600m        Id = 0x1000040041
	FukushimaTurf2600m           Id = 0x100004002F
	FukushimaGoodToFirmTurf2600m Id = 0x1000040028
	FukushimaGoodTurf2600m       Id = 0x1000040024
	FukushimaYieldingTurf2600m   Id = 0x1000040022
	FukushimaSoftTurf2600m       Id = 0x1000040021
	KokuraTurf2600m              Id = 0x100004001F
	KokuraGoodToFirmTurf2600m    Id = 0x1000040018
	KokuraGoodTurf2600m          Id = 0x1000040014
	KokuraYieldingTurf2600m      Id = 0x1000040012
	KokuraSoftTurf2600m          Id = 0x1000040011
	HanshinTurf3000m             Id = 0x100002040F
	HanshinGoodToFirmTurf3000m   Id = 0x1000020408
	HanshinGoodTurf3000m         Id = 0x1000020404
	HanshinYieldingTurf3000m     Id = 0x1000020402
	HanshinSoftTurf3000m         Id = 0x1000020401
	ChukyoTurf3000m              Id = 0x100002010F
	ChukyoGoodToFirmTurf3000m    Id = 0x1000020108
	ChukyoGoodTurf3000m          Id = 0x1000020104
	ChukyoYieldingTurf3000m      Id = 0x1000020102
	ChukyoSoftTurf3000m          Id = 0x1000020101
	KyotoTurf3200m               Id = 0x100001080F
	KyotoGoodToFirmTurf3200m     Id = 0x1000010808
	KyotoGoodTurf3200m           Id = 0x1000010804
	KyotoYieldingTurf3200m       Id = 0x1000010802
	KyotoSoftTurf3200m           Id = 0x1000010801
	TokyoTurf3400m               Id = 0x100000A00F
	TokyoGoodToFirmTurf3400m     Id = 0x100000A008
	TokyoGoodTurf3400m           Id = 0x100000A004
	TokyoYieldingTurf3400m       Id = 0x100000A002
	TokyoSoftTurf3400m           Id = 0x100000A001
	NakayamaTurf3600m            Id = 0x100000500F
	NakayamaGoodToFirmTurf3600m  Id = 0x1000005008
	NakayamaGoodTurf3600m        Id = 0x1000005004
	NakayamaYieldingTurf3600m    Id = 0x1000005002
	NakayamaSoftTurf3600m        Id = 0x1000005001
	SapporoDirt1000m             Id = 0xC0000008F
	SapporoGoodToFirmDirt1000m   Id = 0xC00000088
	SapporoGoodDirt1000m         Id = 0xC00000084
	SapporoYieldingDirt1000m     Id = 0xC00000082
	SapporoSoftDirt1000m         Id = 0xC00000081
	HakodateDirt1000m            Id = 0xC0000004F
	HakodateGoodToFirmDirt1000m  Id = 0xC00000048
	HakodateGoodDirt1000m        Id = 0xC00000044
	HakodateYieldingDirt1000m    Id = 0xC00000042
	HakodateSoftDirt1000m        Id = 0xC00000041
	KokuraDirt1000m              Id = 0xC0000001F
	KokuraGoodToFirmDirt1000m    Id = 0xC00000018
	KokuraGoodDirt1000m          Id = 0xC00000014
	KokuraYieldingDirt1000m      Id = 0xC00000012
	KokuraSoftDirt1000m          Id = 0xC00000011
	FukushimaDirt1150m           Id = 0xA0000002F
	FukushimaGoodToFirmDirt1150m Id = 0xA00000028
	FukushimaGoodDirt1150m       Id = 0xA00000024
	FukushimaYieldingDirt1150m   Id = 0xA00000022
	FukushimaSoftDirt1150m       Id = 0xA00000021
	NakayamaDirt1200m            Id = 0x90000100F
	NakayamaGoodToFirmDirt1200m  Id = 0x900001008
	NakayamaGoodDirt1200m        Id = 0x900001004
	NakayamaYieldingDirt1200m    Id = 0x900001002
	NakayamaSoftDirt1200m        Id = 0x900001001
	KyotoDirt1200m               Id = 0x90000080F
	KyotoGoodToFirmDirt1200m     Id = 0x900000808
	KyotoGoodDirt1200m           Id = 0x900000804
	KyotoYieldingDirt1200m       Id = 0x900000802
	KyotoSoftDirt1200m           Id = 0x900000801
	NiigataDirt1200m             Id = 0x90000020F
	NiigataGoodToFirmDirt1200m   Id = 0x900000208
	NiigataGoodDirt1200m         Id = 0x900000204
	NiigataYieldingDirt1200m     Id = 0x900000202
	NiigataSoftDirt1200m         Id = 0x900000201
	ChukyoDirt1200m              Id = 0x90000010F
	ChukyoGoodToFirmDirt1200m    Id = 0x900000108
	ChukyoGoodDirt1200m          Id = 0x900000104
	ChukyoYieldingDirt1200m      Id = 0x900000102
	ChukyoSoftDirt1200m          Id = 0x900000101
	TokyoDirt1300m               Id = 0x88000200F
	TokyoGoodToFirmDirt1300m     Id = 0x880002008
	TokyoGoodDirt1300m           Id = 0x880002004
	TokyoYieldingDirt1300m       Id = 0x880002002
	TokyoSoftDirt1300m           Id = 0x880002001
	TokyoDirt1400m               Id = 0x84000200F
	TokyoGoodToFirmDirt1400m     Id = 0x840002008
	TokyoGoodDirt1400m           Id = 0x840002004
	TokyoYieldingDirt1400m       Id = 0x840002002
	TokyoSoftDirt1400m           Id = 0x840002001
	KyotoDirt1400m               Id = 0x84000080F
	KyotoGoodToFirmDirt1400m     Id = 0x840000808
	KyotoGoodDirt1400m           Id = 0x840000804
	KyotoYieldingDirt1400m       Id = 0x840000802
	KyotoSoftDirt1400m           Id = 0x840000801
	HanshinDirt1400m             Id = 0x84000040F
	HanshinGoodToFirmDirt1400m   Id = 0x840000408
	HanshinGoodDirt1400m         Id = 0x840000404
	HanshinYieldingDirt1400m     Id = 0x840000402
	HanshinSoftDirt1400m         Id = 0x840000401
	ChukyoDirt1400m              Id = 0x84000010F
	ChukyoGoodToFirmDirt1400m    Id = 0x840000108
	ChukyoGoodDirt1400m          Id = 0x840000104
	ChukyoYieldingDirt1400m      Id = 0x840000102
	ChukyoSoftDirt1400m          Id = 0x840000101
	TokyoDirt1600m               Id = 0x81000200F
	TokyoGoodToFirmDirt1600m     Id = 0x810002008
	TokyoGoodDirt1600m           Id = 0x810002004
	TokyoYieldingDirt1600m       Id = 0x810002002
	TokyoSoftDirt1600m           Id = 0x810002001
	SapporoDirt1700m             Id = 0x80800008F
	SapporoGoodToFirmDirt1700m   Id = 0x808000088
	SapporoGoodDirt1700m         Id = 0x808000084
	SapporoYieldingDirt1700m     Id = 0x808000082
	SapporoSoftDirt1700m         Id = 0x808000081
	HakodateDirt1700m            Id = 0x80800004F
	HakodateGoodToFirmDirt1700m  Id = 0x808000048
	HakodateGoodDirt1700m        Id = 0x808000044
	HakodateYieldingDirt1700m    Id = 0x808000042
	HakodateSoftDirt1700m        Id = 0x808000041
	FukushimaDirt1700m           Id = 0x80800002F
	FukushimaGoodToFirmDirt1700m Id = 0x808000028
	FukushimaGoodDirt1700m       Id = 0x808000024
	FukushimaYieldingDirt1700m   Id = 0x808000022
	FukushimaSoftDirt1700m       Id = 0x808000021
	KokuraDirt1700m              Id = 0x80800001F
	KokuraGoodToFirmDirt1700m    Id = 0x808000018
	KokuraGoodDirt1700m          Id = 0x808000014
	KokuraYieldingDirt1700m      Id = 0x808000012
	KokuraSoftDirt1700m          Id = 0x808000011
	NakayamaDirt1800m            Id = 0x80400100F
	NakayamaGoodToFirmDirt1800m  Id = 0x804001008
	NakayamaGoodDirt1800m        Id = 0x804001004
	NakayamaYieldingDirt1800m    Id = 0x804001002
	NakayamaSoftDirt1800m        Id = 0x804001001
	KyotoDirt1800m               Id = 0x80400080F
	KyotoGoodToFirmDirt1800m     Id = 0x804000808
	KyotoGoodDirt1800m           Id = 0x804000804
	KyotoYieldingDirt1800m       Id = 0x804000802
	KyotoSoftDirt1800m           Id = 0x804000801
	HanshinDirt1800m             Id = 0x80400040F
	HanshinGoodToFirmDirt1800m   Id = 0x804000408
	HanshinGoodDirt1800m         Id = 0x804000404
	HanshinYieldingDirt1800m     Id = 0x804000402
	HanshinSoftDirt1800m         Id = 0x804000401
	NiigataDirt1800m             Id = 0x80400020F
	NiigataGoodToFirmDirt1800m   Id = 0x804000208
	NiigataGoodDirt1800m         Id = 0x804000204
	NiigataYieldingDirt1800m     Id = 0x804000202
	NiigataSoftDirt1800m         Id = 0x804000201
	ChukyoDirt1800m              Id = 0x80400010F
	ChukyoGoodToFirmDirt1800m    Id = 0x804000108
	ChukyoGoodDirt1800m          Id = 0x804000104
	ChukyoYieldingDirt1800m      Id = 0x804000102
	ChukyoSoftDirt1800m          Id = 0x804000101
	KyotoDirt1900m               Id = 0x80200080F
	KyotoGoodToFirmDirt1900m     Id = 0x802000808
	KyotoGoodDirt1900m           Id = 0x802000804
	KyotoYieldingDirt1900m       Id = 0x802000802
	KyotoSoftDirt1900m           Id = 0x802000801
	ChukyoDirt1900m              Id = 0x80200010F
	ChukyoGoodToFirmDirt1900m    Id = 0x802000108
	ChukyoGoodDirt1900m          Id = 0x802000104
	ChukyoYieldingDirt1900m      Id = 0x802000102
	ChukyoSoftDirt1900m          Id = 0x802000101
	HanshinDirt2000m             Id = 0x80100040F
	HanshinGoodToFirmDirt2000m   Id = 0x801000408
	HanshinGoodDirt2000m         Id = 0x801000404
	HanshinYieldingDirt2000m     Id = 0x801000402
	HanshinSoftDirt2000m         Id = 0x801000401
	TokyoDirt2100m               Id = 0x80080200F
	TokyoGoodToFirmDirt2100m     Id = 0x800802008
	TokyoGoodDirt2100m           Id = 0x800802004
	TokyoYieldingDirt2100m       Id = 0x800802002
	TokyoSoftDirt2100m           Id = 0x800802001
	NakayamaDirt2400m            Id = 0x80010100F
	NakayamaGoodToFirmDirt2400m  Id = 0x800101008
	NakayamaGoodDirt2400m        Id = 0x800101004
	NakayamaYieldingDirt2400m    Id = 0x800101002
	NakayamaSoftDirt2400m        Id = 0x800101001
	SapporoDirt2400m             Id = 0x80010008F
	SapporoGoodToFirmDirt2400m   Id = 0x800100088
	SapporoGoodDirt2400m         Id = 0x800100084
	SapporoYieldingDirt2400m     Id = 0x800100082
	SapporoSoftDirt2400m         Id = 0x800100081
	HakodateDirt2400m            Id = 0x80010004F
	HakodateGoodToFirmDirt2400m  Id = 0x800100048
	HakodateGoodDirt2400m        Id = 0x800100044
	HakodateYieldingDirt2400m    Id = 0x800100042
	HakodateSoftDirt2400m        Id = 0x800100041
	FukushimaDirt2400m           Id = 0x80010002F
	FukushimaGoodToFirmDirt2400m Id = 0x800100028
	FukushimaGoodDirt2400m       Id = 0x800100024
	FukushimaYieldingDirt2400m   Id = 0x800100022
	FukushimaSoftDirt2400m       Id = 0x800100021
	KokuraDirt2400m              Id = 0x80010001F
	KokuraGoodToFirmDirt2400m    Id = 0x800100018
	KokuraGoodDirt2400m          Id = 0x800100014
	KokuraYieldingDirt2400m      Id = 0x800100012
	KokuraSoftDirt2400m          Id = 0x800100011
	NakayamaDirt2500m            Id = 0x80008100F
	NakayamaGoodToFirmDirt2500m  Id = 0x800081008
	NakayamaGoodDirt2500m        Id = 0x800081004
	NakayamaYieldingDirt2500m    Id = 0x800081002
	NakayamaSoftDirt2500m        Id = 0x800081001
	NiigataDirt2500m             Id = 0x80008020F
	NiigataGoodToFirmDirt2500m   Id = 0x800080208
	NiigataGoodDirt2500m         Id = 0x800080204
	NiigataYieldingDirt2500m     Id = 0x800080202
	NiigataSoftDirt2500m         Id = 0x800080201
)

var originFilterIdMap = map[Id]string{
	All2:          "全レース",
	Turf2:         "芝",
	Dirt2:         "ダート",
	Tokyo:         "東京",
	Nakayama:      "中山",
	Kyoto:         "京都",
	Hanshin:       "阪神",
	Niigata:       "新潟",
	Chukyo:        "中京",
	Sapporo:       "札幌",
	Hakodate:      "函館",
	Fukushima:     "福島",
	Kokura:        "小倉",
	GoodToFirm:    "良",
	Good:          "稍重",
	Yielding:      "重",
	Soft:          "不良",
	Distance1000m: "1000m",
	Distance1150m: "1150m",
	Distance1200m: "1200m",
	Distance1300m: "1300m",
	Distance1400m: "1400m",
	Distance1500m: "1500m",
	Distance1600m: "1600m",
	Distance1700m: "1700m",
	Distance1800m: "1800m",
	Distance1900m: "1900m",
	Distance2000m: "2000m",
	Distance2100m: "2100m",
	Distance2200m: "2200m",
	Distance2300m: "2300m",
	Distance2400m: "2400m",
	Distance2500m: "2500m",
	Distance2600m: "2600m",
	Distance3000m: "3000m",
	Distance3200m: "3200m",
	Distance3400m: "3400m",
	Distance3600m: "3600m",
}

func (i Id) Value() uint64 {
	return uint64(i)
}

func (i Id) String() string {
	id, _ := originFilterIdMap[i]
	return id
}

func (i Id) OriginFilters() []Id {
	var ids []Id

	if i == All2 {
		return []Id{All2}
	}

	for id := range originFilterIdMap {
		if i&id == id {
			ids = append(ids, id)
		}
	}

	sort.Slice(ids, func(i, j int) bool {
		return ids[i] > ids[j]
	})

	return ids
}
