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
	All2          Id = 0x1FFFFFFFFF
	Turf2         Id = 0x1000000000
	Dirt2         Id = 0x800000000
	Distance1000m Id = 0x4000000000
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
	All2:       "全レース",
	Turf2:      "芝",
	Dirt2:      "ダート",
	Tokyo:      "東京",
	Nakayama:   "中山",
	Kyoto:      "京都",
	Hanshin:    "阪神",
	Niigata:    "新潟",
	Chukyo:     "中京",
	Sapporo:    "札幌",
	Hakodate:   "函館",
	Fukushima:  "福島",
	Kokura:     "小倉",
	GoodToFirm: "良",
	Good:       "稍重",
	Yielding:   "重",
	Soft:       "不良",

	TurfGoodToFirm:               "芝良",
	TurfGood:                     "芝稍重",
	TurfYielding:                 "芝重",
	TurfSoft:                     "芝不良",
	DirtGoodToFirm:               "ダート良",
	DirtGood:                     "ダート稍重",
	DirtYielding:                 "ダート重",
	DirtSoft:                     "ダート不良",
	NiigataTurf1000m:             "新潟芝1000m",
	NiigataGoodToFirmTurf1000m:   "新潟良芝1000m",
	NiigataGoodTurf1000m:         "新潟稍芝1000m",
	NiigataYieldingTurf1000m:     "新潟重芝1000m",
	NiigataSoftTurf1000m:         "新潟不芝1000m",
	HakodateTurf1000m:            "函館芝1000m",
	HakodateGoodToFirmTurf1000m:  "函館良芝1000m",
	HakodateGoodTurf1000m:        "函館稍芝1000m",
	HakodateYieldingTurf1000m:    "函館重芝1000m",
	HakodateSoftTurf1000m:        "函館不芝1000m",
	NakayamaTurf1200m:            "中山芝1200m",
	NakayamaGoodToFirmTurf1200m:  "中山良芝1200m",
	NakayamaGoodTurf1200m:        "中山稍芝1200m",
	NakayamaYieldingTurf1200m:    "中山重芝1200m",
	NakayamaSoftTurf1200m:        "中山不芝1200m",
	KyotoTurf1200m:               "京都芝1200m",
	KyotoGoodToFirmTurf1200m:     "京都良芝1200m",
	KyotoGoodTurf1200m:           "京都稍芝1200m",
	KyotoYieldingTurf1200m:       "京都重芝1200m",
	KyotoSoftTurf1200m:           "京都不芝1200m",
	HanshinTurf1200m:             "阪神芝1200m",
	HanshinGoodToFirmTurf1200m:   "阪神良芝1200m",
	HanshinGoodTurf1200m:         "阪神稍芝1200m",
	HanshinYieldingTurf1200m:     "阪神重芝1200m",
	HanshinSoftTurf1200m:         "阪神不芝1200m",
	NiigataTurf1200m:             "新潟芝1200m",
	NiigataGoodToFirmTurf1200m:   "新潟良芝1200m",
	NiigataGoodTurf1200m:         "新潟稍芝1200m",
	NiigataYieldingTurf1200m:     "新潟重芝1200m",
	NiigataSoftTurf1200m:         "新潟不芝1200m",
	ChukyoTurf1200m:              "中京芝1200m",
	ChukyoGoodToFirmTurf1200m:    "中京良芝1200m",
	ChukyoGoodTurf1200m:          "中京稍芝1200m",
	ChukyoYieldingTurf1200m:      "中京重芝1200m",
	ChukyoSoftTurf1200m:          "中京不芝1200m",
	SapporoTurf1200m:             "札幌芝1200m",
	SapporoGoodToFirmTurf1200m:   "札幌良芝1200m",
	SapporoGoodTurf1200m:         "札幌稍芝1200m",
	SapporoYieldingTurf1200m:     "札幌重芝1200m",
	SapporoSoftTurf1200m:         "札幌不芝1200m",
	HakodateTurf1200m:            "函館芝1200m",
	HakodateGoodToFirmTurf1200m:  "函館良芝1200m",
	HakodateGoodTurf1200m:        "函館稍芝1200m",
	HakodateYieldingTurf1200m:    "函館重芝1200m",
	HakodateSoftTurf1200m:        "函館不芝1200m",
	FukushimaTurf1200m:           "福島芝1200m",
	FukushimaGoodToFirmTurf1200m: "福島良芝1200m",
	FukushimaGoodTurf1200m:       "福島稍芝1200m",
	FukushimaYieldingTurf1200m:   "福島重芝1200m",
	FukushimaSoftTurf1200m:       "福島不芝1200m",
	KokuraTurf1200m:              "小倉芝1200m",
	KokuraGoodToFirmTurf1200m:    "小倉良芝1200m",
	KokuraGoodTurf1200m:          "小倉稍芝1200m",
	KokuraYieldingTurf1200m:      "小倉重芝1200m",
	KokuraSoftTurf1200m:          "小倉不芝1200m",
	TokyoTurf1400m:               "東京芝1400m",
	TokyoGoodToFirmTurf1400m:     "東京良芝1400m",
	TokyoGoodTurf1400m:           "東京稍芝1400m",
	TokyoYieldingTurf1400m:       "東京重芝1400m",
	TokyoSoftTurf1400m:           "東京不芝1400m",
	KyotoTurf1400m:               "京都芝1400m",
	KyotoGoodToFirmTurf1400m:     "京都良芝1400m",
	KyotoGoodTurf1400m:           "京都稍芝1400m",
	KyotoYieldingTurf1400m:       "京都重芝1400m",
	KyotoSoftTurf1400m:           "京都不芝1400m",
	HanshinTurf1400m:             "阪神芝1400m",
	HanshinGoodToFirmTurf1400m:   "阪神良芝1400m",
	HanshinGoodTurf1400m:         "阪神稍芝1400m",
	HanshinYieldingTurf1400m:     "阪神重芝1400m",
	HanshinSoftTurf1400m:         "阪神不芝1400m",
	NiigataTurf1400m:             "新潟芝1400m",
	NiigataGoodToFirmTurf1400m:   "新潟良芝1400m",
	NiigataGoodTurf1400m:         "新潟稍芝1400m",
	NiigataYieldingTurf1400m:     "新潟重芝1400m",
	NiigataSoftTurf1400m:         "新潟不芝1400m",
	ChukyoTurf1400m:              "中京芝1400m",
	ChukyoGoodToFirmTurf1400m:    "中京良芝1400m",
	ChukyoGoodTurf1400m:          "中京稍芝1400m",
	ChukyoYieldingTurf1400m:      "中京重芝1400m",
	ChukyoSoftTurf1400m:          "中京不芝1400m",
	SapporoTurf1500m:             "札幌芝1500m",
	SapporoGoodToFirmTurf1500m:   "札幌良芝1500m",
	SapporoGoodTurf1500m:         "札幌稍芝1500m",
	SapporoYieldingTurf1500m:     "札幌重芝1500m",
	SapporoSoftTurf1500m:         "札幌不芝1500m",
	NakayamaTurf1600m:            "中山芝1600m",
	NakayamaGoodToFirmTurf1600m:  "中山良芝1600m",
	NakayamaGoodTurf1600m:        "中山稍芝1600m",
	NakayamaYieldingTurf1600m:    "中山重芝1600m",
	NakayamaSoftTurf1600m:        "中山不芝1600m",
	TokyoTurf1600m:               "東京芝1600m",
	TokyoGoodToFirmTurf1600m:     "東京良芝1600m",
	TokyoGoodTurf1600m:           "東京稍芝1600m",
	TokyoYieldingTurf1600m:       "東京重芝1600m",
	TokyoSoftTurf1600m:           "東京不芝1600m",
	KyotoTurf1600m:               "京都芝1600m",
	KyotoGoodToFirmTurf1600m:     "京都良芝1600m",
	KyotoGoodTurf1600m:           "京都稍芝1600m",
	KyotoYieldingTurf1600m:       "京都重芝1600m",
	KyotoSoftTurf1600m:           "京都不芝1600m",
	HanshinTurf1600m:             "阪神芝1600m",
	HanshinGoodToFirmTurf1600m:   "阪神良芝1600m",
	HanshinGoodTurf1600m:         "阪神稍芝1600m",
	HanshinYieldingTurf1600m:     "阪神重芝1600m",
	HanshinSoftTurf1600m:         "阪神不芝1600m",
	ChukyoTurf1600m:              "中京芝1600m",
	ChukyoGoodToFirmTurf1600m:    "中京良芝1600m",
	ChukyoGoodTurf1600m:          "中京稍芝1600m",
	ChukyoYieldingTurf1600m:      "中京重芝1600m",
	ChukyoSoftTurf1600m:          "中京不芝1600m",
	NakayamaTurf1800m:            "中山芝1800m",
	NakayamaGoodToFirmTurf1800m:  "中山良芝1800m",
	NakayamaGoodTurf1800m:        "中山稍芝1800m",
	NakayamaYieldingTurf1800m:    "中山重芝1800m",
	NakayamaSoftTurf1800m:        "中山不芝1800m",
	TokyoTurf1800m:               "東京芝1800m",
	TokyoGoodToFirmTurf1800m:     "東京良芝1800m",
	TokyoGoodTurf1800m:           "東京稍芝1800m",
	TokyoYieldingTurf1800m:       "東京重芝1800m",
	TokyoSoftTurf1800m:           "東京不芝1800m",
	KyotoTurf1800m:               "京都芝1800m",
	KyotoGoodToFirmTurf1800m:     "京都良芝1800m",
	KyotoGoodTurf1800m:           "京都稍芝1800m",
	KyotoYieldingTurf1800m:       "京都重芝1800m",
	KyotoSoftTurf1800m:           "京都不芝1800m",
	HanshinTurf1800m:             "阪神芝1800m",
	HanshinGoodToFirmTurf1800m:   "阪神良芝1800m",
	HanshinGoodTurf1800m:         "阪神稍芝1800m",
	HanshinYieldingTurf1800m:     "阪神重芝1800m",
	HanshinSoftTurf1800m:         "阪神不芝1800m",
	NiigataTurf1800m:             "新潟芝1800m",
	NiigataGoodToFirmTurf1800m:   "新潟良芝1800m",
	NiigataGoodTurf1800m:         "新潟稍芝1800m",
	NiigataYieldingTurf1800m:     "新潟重芝1800m",
	NiigataSoftTurf1800m:         "新潟不芝1800m",
	SapporoTurf1800m:             "札幌芝1800m",
	SapporoGoodToFirmTurf1800m:   "札幌良芝1800m",
	SapporoGoodTurf1800m:         "札幌重芝1800m",
	SapporoYieldingTurf1800m:     "札幌稍芝1800m",
	SapporoSoftTurf1800m:         "札幌不芝1800m",
	HakodateTurf1800m:            "函館芝1800m",
	HakodateGoodToFirmTurf1800m:  "函館良芝1800m",
	HakodateGoodTurf1800m:        "函館稍芝1800m",
	HakodateYieldingTurf1800m:    "函館重芝1800m",
	HakodateSoftTurf1800m:        "函館不芝1800m",
	FukushimaTurf1800m:           "福島芝1800m",
	FukushimaGoodToFirmTurf1800m: "福島良芝1800m",
	FukushimaGoodTurf1800m:       "福島稍芝1800m",
	FukushimaYieldingTurf1800m:   "福島重芝1800m",
	FukushimaSoftTurf1800m:       "福島不芝1800m",
	KokuraTurf1800m:              "小倉芝1800m",
	KokuraGoodToFirmTurf1800m:    "小倉良芝1800m",
	KokuraGoodTurf1800m:          "小倉稍芝1800m",
	KokuraYieldingTurf1800m:      "小倉重芝1800m",
	KokuraSoftTurf1800m:          "小倉不芝1800m",
	NakayamaTurf2000m:            "中山芝2000m",
	NakayamaGoodToFirmTurf2000m:  "中山良芝2000m",
	NakayamaGoodTurf2000m:        "中山稍芝2000m",
	NakayamaYieldingTurf2000m:    "中山重芝2000m",
	NakayamaSoftTurf2000m:        "中山不芝2000m",
	TokyoTurf2000m:               "東京芝2000m",
	TokyoGoodToFirmTurf2000m:     "東京良芝2000m",
	TokyoGoodTurf2000m:           "東京稍芝2000m",
	TokyoYieldingTurf2000m:       "東京重芝2000m",
	TokyoSoftTurf2000m:           "東京不芝2000m",
	KyotoTurf2000m:               "京都芝2000m",
	KyotoGoodToFirmTurf2000m:     "京都良芝2000m",
	KyotoGoodTurf2000m:           "京都稍芝2000m",
	KyotoYieldingTurf2000m:       "京都重芝2000m",
	KyotoSoftTurf2000m:           "京都不芝2000m",
	NiigataTurf2000m:             "新潟芝2000m",
	NiigataGoodToFirmTurf2000m:   "新潟良芝2000m",
	NiigataGoodTurf2000m:         "新潟稍芝2000m",
	NiigataYieldingTurf2000m:     "新潟重芝2000m",
	NiigataSoftTurf2000m:         "新潟不芝2000m",
	ChukyoTurf2000m:              "中京芝2000m",
	ChukyoGoodToFirmTurf2000m:    "中京良芝2000m",
	ChukyoGoodTurf2000m:          "中京稍芝2000m",
	ChukyoYieldingTurf2000m:      "中京重芝2000m",
	ChukyoSoftTurf2000m:          "中京不芝2000m",
	SapporoTurf2000m:             "札幌芝2000m",
	SapporoGoodToFirmTurf2000m:   "札幌良芝2000m",
	SapporoGoodTurf2000m:         "札幌稍芝2000m",
	SapporoYieldingTurf2000m:     "札幌重芝2000m",
	SapporoSoftTurf2000m:         "札幌不芝2000m",
	HakodateTurf2000m:            "函館芝2000m",
	HakodateGoodToFirmTurf2000m:  "函館良芝2000m",
	HakodateGoodTurf2000m:        "函館稍芝2000m",
	HakodateYieldingTurf2000m:    "函館重芝2000m",
	HakodateSoftTurf2000m:        "函館不芝2000m",
	FukushimaTurf2000m:           "福島芝2000m",
	FukushimaGoodToFirmTurf2000m: "福島良芝2000m",
	FukushimaGoodTurf2000m:       "福島稍芝2000m",
	FukushimaYieldingTurf2000m:   "福島重芝2000m",
	FukushimaSoftTurf2000m:       "福島不芝2000m",
	KokuraTurf2000m:              "小倉芝2000m",
	KokuraGoodToFirmTurf2000m:    "小倉良芝2000m",
	KokuraGoodTurf2000m:          "小倉稍芝2000m",
	KokuraYieldingTurf2000m:      "小倉重芝2000m",
	KokuraSoftTurf2000m:          "小倉不芝2000m",
	NakayamaTurf2200m:            "中山芝2200m",
	NakayamaGoodToFirmTurf2200m:  "中山芝2200m",
	NakayamaGoodTurf2200m:        "中山芝2200m",
	NakayamaYieldingTurf2200m:    "中山芝2200m",
	NakayamaSoftTurf2200m:        "中山芝2200m",
	KyotoTurf2200m:               "京都芝2200m",
	KyotoGoodToFirmTurf2200m:     "京都良芝2200m",
	KyotoGoodTurf2200m:           "京都稍芝2200m",
	KyotoYieldingTurf2200m:       "京都重芝2200m",
	KyotoSoftTurf2200m:           "京都不芝2200m",
	HanshinTurf2200m:             "阪神芝2200m",
	HanshinGoodToFirmTurf2200m:   "阪神良芝2200m",
	HanshinGoodTurf2200m:         "阪神稍芝2200m",
	HanshinYieldingTurf2200m:     "阪神重芝2200m",
	HanshinSoftTurf2200m:         "阪神不芝2200m",
	NiigataTurf2200m:             "新潟芝2200m",
	NiigataGoodToFirmTurf2200m:   "新潟良芝2200m",
	NiigataGoodTurf2200m:         "新潟稍芝2200m",
	NiigataYieldingTurf2200m:     "新潟重芝2200m",
	NiigataSoftTurf2200m:         "新潟不芝2200m",
	ChukyoTurf2200m:              "中京芝2200m",
	ChukyoGoodToFirmTurf2200m:    "中京良芝2200m",
	ChukyoGoodTurf2200m:          "中京稍芝2200m",
	ChukyoYieldingTurf2200m:      "中京重芝2200m",
	ChukyoSoftTurf2200m:          "中京不芝2200m",
	TokyoTurf2300m:               "東京芝2300m",
	TokyoGoodToFirmTurf2300m:     "東京良芝2300m",
	TokyoGoodTurf2300m:           "東京稍芝2300m",
	TokyoYieldingTurf2300m:       "東京重芝2300m",
	TokyoSoftTurf2300m:           "東京不芝2300m",
	TokyoTurf2400m:               "東京芝2400m",
	TokyoGoodToFirmTurf2400m:     "東京良芝2400m",
	TokyoGoodTurf2400m:           "東京稍芝2400m",
	TokyoYieldingTurf2400m:       "東京重芝2400m",
	TokyoSoftTurf2400m:           "東京不芝2400m",
	KyotoTurf2400m:               "京都芝2400m",
	KyotoGoodToFirmTurf2400m:     "京都良芝2400m",
	KyotoGoodTurf2400m:           "京都稍芝2400m",
	KyotoYieldingTurf2400m:       "京都重芝2400m",
	KyotoSoftTurf2400m:           "京都不芝2400m",
	HanshinTurf2400m:             "阪神芝2400m",
	HanshinGoodToFirmTurf2400m:   "阪神良芝2400m",
	HanshinGoodTurf2400m:         "阪神稍芝2400m",
	HanshinYieldingTurf2400m:     "阪神重芝2400m",
	HanshinSoftTurf2400m:         "阪神不芝2400m",
	NiigataTurf2400m:             "新潟芝2400m",
	NiigataGoodToFirmTurf2400m:   "新潟良芝2400m",
	NiigataGoodTurf2400m:         "新潟稍芝2400m",
	NiigataYieldingTurf2400m:     "新潟重芝2400m",
	NiigataSoftTurf2400m:         "新潟不芝2400m",
	NakayamaTurf2500m:            "中山芝2500m",
	NakayamaGoodToFirmTurf2500m:  "中山良芝2500m",
	NakayamaGoodTurf2500m:        "中山稍芝2500m",
	NakayamaYieldingTurf2500m:    "中山重芝2500m",
	NakayamaSoftTurf2500m:        "中山不芝2500m",
	TokyoTurf2500m:               "東京芝2500m",
	TokyoGoodToFirmTurf2500m:     "東京良芝2500m",
	TokyoGoodTurf2500m:           "東京稍芝2500m",
	TokyoYieldingTurf2500m:       "東京重芝2500m",
	TokyoSoftTurf2500m:           "東京不芝2500m",
	HanshinTurf2600m:             "阪神芝2600m",
	HanshinGoodToFirmTurf2600m:   "阪神良芝2600m",
	HanshinGoodTurf2600m:         "阪神稍芝2600m",
	HanshinYieldingTurf2600m:     "阪神重芝2600m",
	HanshinSoftTurf2600m:         "阪神不芝2600m",
	SapporoTurf2600m:             "札幌芝2600m",
	SapporoGoodToFirmTurf2600m:   "札幌良芝2600m",
	SapporoGoodTurf2600m:         "札幌稍芝2600m",
	SapporoYieldingTurf2600m:     "札幌重芝2600m",
	SapporoSoftTurf2600m:         "札幌不芝2600m",
	HakodateTurf2600m:            "函館芝2600m",
	HakodateGoodToFirmTurf2600m:  "函館良芝2600m",
	HakodateGoodTurf2600m:        "函館稍芝2600m",
	HakodateYieldingTurf2600m:    "函館重芝2600m",
	HakodateSoftTurf2600m:        "函館不芝2600m",
	FukushimaTurf2600m:           "福島芝2600m",
	FukushimaGoodToFirmTurf2600m: "福島良芝2600m",
	FukushimaGoodTurf2600m:       "福島稍芝2600m",
	FukushimaYieldingTurf2600m:   "福島重芝2600m",
	FukushimaSoftTurf2600m:       "福島不芝2600m",
	KokuraTurf2600m:              "小倉芝2600m",
	KokuraGoodToFirmTurf2600m:    "小倉良芝2600m",
	KokuraGoodTurf2600m:          "小倉稍芝2600m",
	KokuraYieldingTurf2600m:      "小倉重芝2600m",
	KokuraSoftTurf2600m:          "小倉不芝2600m",
	HanshinTurf3000m:             "阪神芝3000m",
	HanshinGoodToFirmTurf3000m:   "阪神良芝3000m",
	HanshinGoodTurf3000m:         "阪神稍芝3000m",
	HanshinYieldingTurf3000m:     "阪神重芝3000m",
	HanshinSoftTurf3000m:         "阪神不芝3000m",
	ChukyoTurf3000m:              "中京芝3000m",
	ChukyoGoodToFirmTurf3000m:    "中京良芝3000m",
	ChukyoGoodTurf3000m:          "中京稍芝3000m",
	ChukyoYieldingTurf3000m:      "中京重芝3000m",
	ChukyoSoftTurf3000m:          "中京不芝3000m",
	KyotoTurf3200m:               "京都芝3200m",
	KyotoGoodToFirmTurf3200m:     "京都良芝3200m",
	KyotoGoodTurf3200m:           "京都稍芝3200m",
	KyotoYieldingTurf3200m:       "京都重芝3200m",
	KyotoSoftTurf3200m:           "京都不芝3200m",
	TokyoTurf3400m:               "東京芝3400m",
	TokyoGoodToFirmTurf3400m:     "東京良芝3400m",
	TokyoGoodTurf3400m:           "東京稍芝3400m",
	TokyoYieldingTurf3400m:       "東京重芝3400m",
	TokyoSoftTurf3400m:           "東京不芝3400m",
	NakayamaTurf3600m:            "中山芝3600m",
	NakayamaGoodToFirmTurf3600m:  "中山良芝3600m",
	NakayamaGoodTurf3600m:        "中山稍芝3600m",
	NakayamaYieldingTurf3600m:    "中山重芝3600m",
	NakayamaSoftTurf3600m:        "中山不芝3600m",
	SapporoDirt1000m:             "札幌ダ1000m",
	SapporoGoodToFirmDirt1000m:   "札幌良ダ1000m",
	SapporoGoodDirt1000m:         "札幌稍ダ1000m",
	SapporoYieldingDirt1000m:     "札幌重ダ1000m",
	SapporoSoftDirt1000m:         "札幌不ダ1000m",
	HakodateDirt1000m:            "函館ダ1000m",
	HakodateGoodToFirmDirt1000m:  "函館良ダ1000m",
	HakodateGoodDirt1000m:        "函館稍ダ1000m",
	HakodateYieldingDirt1000m:    "函館重ダ1000m",
	HakodateSoftDirt1000m:        "函館不ダ1000m",
	KokuraDirt1000m:              "小倉ダ1000m",
	KokuraGoodToFirmDirt1000m:    "小倉良ダ1000m",
	KokuraGoodDirt1000m:          "小倉稍ダ1000m",
	KokuraYieldingDirt1000m:      "小倉重ダ1000m",
	KokuraSoftDirt1000m:          "小倉不ダ1000m",
	FukushimaDirt1150m:           "福島ダ1150m",
	FukushimaGoodToFirmDirt1150m: "福島良ダ1150m",
	FukushimaGoodDirt1150m:       "福島稍ダ1150m",
	FukushimaYieldingDirt1150m:   "福島重ダ1150m",
	FukushimaSoftDirt1150m:       "福島不ダ1150m",
	NakayamaDirt1200m:            "中山ダ1200m",
	NakayamaGoodToFirmDirt1200m:  "中山良ダ1200m",
	NakayamaGoodDirt1200m:        "中山稍ダ1200m",
	NakayamaYieldingDirt1200m:    "中山重ダ1200m",
	NakayamaSoftDirt1200m:        "中山不ダ1200m",
	KyotoDirt1200m:               "京都ダ1200m",
	KyotoGoodToFirmDirt1200m:     "京都良ダ1200m",
	KyotoGoodDirt1200m:           "京都稍ダ1200m",
	KyotoYieldingDirt1200m:       "京都重ダ1200m",
	KyotoSoftDirt1200m:           "京都不ダ1200m",
	NiigataDirt1200m:             "新潟ダ1200m",
	NiigataGoodToFirmDirt1200m:   "新潟良ダ1200m",
	NiigataGoodDirt1200m:         "新潟稍ダ1200m",
	NiigataYieldingDirt1200m:     "新潟重ダ1200m",
	NiigataSoftDirt1200m:         "新潟不ダ1200m",
	ChukyoDirt1200m:              "中京ダ1200m",
	ChukyoGoodToFirmDirt1200m:    "中京良ダ1200m",
	ChukyoGoodDirt1200m:          "中京稍ダ1200m",
	ChukyoYieldingDirt1200m:      "中京重ダ1200m",
	ChukyoSoftDirt1200m:          "中京不ダ1200m",
	TokyoDirt1300m:               "東京ダ1300m",
	TokyoGoodToFirmDirt1300m:     "東京良ダ1300m",
	TokyoGoodDirt1300m:           "東京稍ダ1300m",
	TokyoYieldingDirt1300m:       "東京重ダ1300m",
	TokyoSoftDirt1300m:           "東京不ダ1300m",
	TokyoDirt1400m:               "東京ダ1400m",
	TokyoGoodToFirmDirt1400m:     "東京良ダ1400m",
	TokyoGoodDirt1400m:           "東京稍ダ1400m",
	TokyoYieldingDirt1400m:       "東京重ダ1400m",
	TokyoSoftDirt1400m:           "東京不ダ1400m",
	KyotoDirt1400m:               "京都ダ1400m",
	KyotoGoodToFirmDirt1400m:     "京都良ダ1400m",
	KyotoGoodDirt1400m:           "京都稍ダ1400m",
	KyotoYieldingDirt1400m:       "京都重ダ1400m",
	KyotoSoftDirt1400m:           "京都不ダ1400m",
	HanshinDirt1400m:             "阪神ダ1400m",
	HanshinGoodToFirmDirt1400m:   "阪神良ダ1400m",
	HanshinGoodDirt1400m:         "阪神稍ダ1400m",
	HanshinYieldingDirt1400m:     "阪神重ダ1400m",
	HanshinSoftDirt1400m:         "阪神不ダ1400m",
	ChukyoDirt1400m:              "中京ダ1400m",
	ChukyoGoodToFirmDirt1400m:    "中京良ダ1400m",
	ChukyoGoodDirt1400m:          "中京稍ダ1400m",
	ChukyoYieldingDirt1400m:      "中京重ダ1400m",
	ChukyoSoftDirt1400m:          "中京不ダ1400m",
	TokyoDirt1600m:               "東京ダ1600m",
	TokyoGoodToFirmDirt1600m:     "東京良ダ1600m",
	TokyoGoodDirt1600m:           "東京稍ダ1600m",
	TokyoYieldingDirt1600m:       "東京重ダ1600m",
	TokyoSoftDirt1600m:           "東京不ダ1600m",
	SapporoDirt1700m:             "札幌ダ1700m",
	SapporoGoodToFirmDirt1700m:   "札幌良ダ1700m",
	SapporoGoodDirt1700m:         "札幌稍ダ1700m",
	SapporoYieldingDirt1700m:     "札幌重ダ1700m",
	SapporoSoftDirt1700m:         "札幌不ダ1700m",
	HakodateDirt1700m:            "函館ダ1700m",
	HakodateGoodToFirmDirt1700m:  "函館良ダ1700m",
	HakodateGoodDirt1700m:        "函館稍ダ1700m",
	HakodateYieldingDirt1700m:    "函館重ダ1700m",
	HakodateSoftDirt1700m:        "函館不ダ1700m",
	FukushimaDirt1700m:           "福島ダ1700m",
	FukushimaGoodToFirmDirt1700m: "福島良ダ1700m",
	FukushimaGoodDirt1700m:       "福島稍ダ1700m",
	FukushimaYieldingDirt1700m:   "福島重ダ1700m",
	FukushimaSoftDirt1700m:       "福島不ダ1700m",
	KokuraDirt1700m:              "小倉ダ1700m",
	KokuraGoodToFirmDirt1700m:    "小倉良ダ1700m",
	KokuraGoodDirt1700m:          "小倉稍ダ1700m",
	KokuraYieldingDirt1700m:      "小倉重ダ1700m",
	KokuraSoftDirt1700m:          "小倉不ダ1700m",
	NakayamaDirt1800m:            "中山ダ1800m",
	NakayamaGoodToFirmDirt1800m:  "中山良ダ1800m",
	NakayamaGoodDirt1800m:        "中山稍ダ1800m",
	NakayamaYieldingDirt1800m:    "中山重ダ1800m",
	NakayamaSoftDirt1800m:        "中山不ダ1800m",
	KyotoDirt1800m:               "京都ダ1800m",
	KyotoGoodToFirmDirt1800m:     "京都良ダ1800m",
	KyotoGoodDirt1800m:           "京都稍ダ1800m",
	KyotoYieldingDirt1800m:       "京都重ダ1800m",
	KyotoSoftDirt1800m:           "京都不ダ1800m",
	HanshinDirt1800m:             "阪神ダ1800m",
	HanshinGoodToFirmDirt1800m:   "阪神良ダ1800m",
	HanshinGoodDirt1800m:         "阪神稍ダ1800m",
	HanshinYieldingDirt1800m:     "阪神重ダ1800m",
	HanshinSoftDirt1800m:         "阪神不ダ1800m",
	NiigataDirt1800m:             "新潟ダ1800m",
	NiigataGoodToFirmDirt1800m:   "新潟良ダ1800m",
	NiigataGoodDirt1800m:         "新潟稍ダ1800m",
	NiigataYieldingDirt1800m:     "新潟重ダ1800m",
	NiigataSoftDirt1800m:         "新潟不ダ1800m",
	ChukyoDirt1800m:              "中京ダ1800m",
	ChukyoGoodToFirmDirt1800m:    "中京良ダ1800m",
	ChukyoGoodDirt1800m:          "中京稍ダ1800m",
	ChukyoYieldingDirt1800m:      "中京重ダ1800m",
	ChukyoSoftDirt1800m:          "中京不ダ1800m",
	KyotoDirt1900m:               "京都ダ1900m",
	KyotoGoodToFirmDirt1900m:     "京都良ダ1900m",
	KyotoGoodDirt1900m:           "京都稍ダ1900m",
	KyotoYieldingDirt1900m:       "京都重ダ1900m",
	KyotoSoftDirt1900m:           "京都不ダ1900m",
	ChukyoDirt1900m:              "中京ダ1900m",
	ChukyoGoodToFirmDirt1900m:    "中京良ダ1900m",
	ChukyoGoodDirt1900m:          "中京稍ダ1900m",
	ChukyoYieldingDirt1900m:      "中京重ダ1900m",
	ChukyoSoftDirt1900m:          "中京不ダ1900m",
	HanshinDirt2000m:             "阪神ダ2000m",
	HanshinGoodToFirmDirt2000m:   "阪神良ダ2000m",
	HanshinGoodDirt2000m:         "阪神稍ダ2000m",
	HanshinYieldingDirt2000m:     "阪神重ダ2000m",
	HanshinSoftDirt2000m:         "阪神不ダ2000m",
	TokyoDirt2100m:               "東京ダ2100m",
	TokyoGoodToFirmDirt2100m:     "東京良ダ2100m",
	TokyoGoodDirt2100m:           "東京稍ダ2100m",
	TokyoYieldingDirt2100m:       "東京重ダ2100m",
	TokyoSoftDirt2100m:           "東京不ダ2100m",
	NakayamaDirt2400m:            "中山ダ2400m",
	NakayamaGoodToFirmDirt2400m:  "中山良ダ2400m",
	NakayamaGoodDirt2400m:        "中山稍ダ2400m",
	NakayamaYieldingDirt2400m:    "中山重ダ2400m",
	NakayamaSoftDirt2400m:        "中山不ダ2400m",
	SapporoDirt2400m:             "札幌ダ2400m",
	SapporoGoodToFirmDirt2400m:   "札幌良ダ2400m",
	SapporoGoodDirt2400m:         "札幌稍ダ2400m",
	SapporoYieldingDirt2400m:     "札幌重ダ2400m",
	SapporoSoftDirt2400m:         "札幌不ダ2400m",
	HakodateDirt2400m:            "函館ダ2400m",
	HakodateGoodToFirmDirt2400m:  "函館良ダ2400m",
	HakodateGoodDirt2400m:        "函館稍ダ2400m",
	HakodateYieldingDirt2400m:    "函館重ダ2400m",
	HakodateSoftDirt2400m:        "函館不ダ2400m",
	FukushimaDirt2400m:           "福島ダ2400m",
	FukushimaGoodToFirmDirt2400m: "福島良ダ2400m",
	FukushimaGoodDirt2400m:       "福島稍ダ2400m",
	FukushimaYieldingDirt2400m:   "福島重ダ2400m",
	FukushimaSoftDirt2400m:       "福島不ダ2400m",
	KokuraDirt2400m:              "小倉ダ2400m",
	KokuraGoodToFirmDirt2400m:    "小倉良ダ2400m",
	KokuraGoodDirt2400m:          "小倉稍ダ2400m",
	KokuraYieldingDirt2400m:      "小倉重ダ2400m",
	KokuraSoftDirt2400m:          "小倉不ダ2400m",
	NakayamaDirt2500m:            "中山ダ2500m",
	NakayamaGoodToFirmDirt2500m:  "中山良ダ2500m",
	NakayamaGoodDirt2500m:        "中山稍ダ2500m",
	NakayamaYieldingDirt2500m:    "中山重ダ2500m",
	NakayamaSoftDirt2500m:        "中山不ダ2500m",
	NiigataDirt2500m:             "新潟ダ2500m",
	NiigataGoodToFirmDirt2500m:   "新潟良ダ2500m",
	NiigataGoodDirt2500m:         "新潟稍ダ2500m",
	NiigataYieldingDirt2500m:     "新潟重ダ2500m",
	NiigataSoftDirt2500m:         "新潟不ダ2500m",
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
