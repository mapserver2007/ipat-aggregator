package filter

import "sort"

type Id uint64

const (
	All           Id = 0x7FFFFFFFFFFF
	Turf          Id = 0x400000000000
	Dirt          Id = 0x200000000000
	Distance1000m Id = 0x100000000000
	Distance1150m Id = 0x80000000000
	Distance1200m Id = 0x40000000000
	Distance1300m Id = 0x20000000000
	Distance1400m Id = 0x10000000000
	Distance1500m Id = 0x8000000000
	Distance1600m Id = 0x4000000000
	Distance1700m Id = 0x2000000000
	Distance1800m Id = 0x1000000000
	Distance1900m Id = 0x800000000
	Distance2000m Id = 0x400000000
	Distance2100m Id = 0x200000000
	Distance2200m Id = 0x100000000
	Distance2300m Id = 0x80000000
	Distance2400m Id = 0x40000000
	Distance2500m Id = 0x20000000
	Distance2600m Id = 0x10000000
	Distance3000m Id = 0x8000000
	Distance3200m Id = 0x4000000
	Distance3400m Id = 0x2000000
	Distance3600m Id = 0x1000000
	Tokyo         Id = 0x800000
	Nakayama      Id = 0x400000
	Kyoto         Id = 0x200000
	Hanshin       Id = 0x100000
	Niigata       Id = 0x80000
	Chukyo        Id = 0x40000
	Sapporo       Id = 0x20000
	Hakodate      Id = 0x10000
	Fukushima     Id = 0x8000
	Kokura        Id = 0x4000
	GoodToFirm    Id = 0x2000
	Good          Id = 0x1000
	Yielding      Id = 0x800
	Soft          Id = 0x400
	Win           Id = 0x80
	Place         Id = 0x100
	Quinella      Id = 0x180
	Exacta        Id = 0x200
	QuinellaPlace Id = 0x280
	Trio          Id = 0x300
	Trifecta      Id = 0x380
	Favorite      Id = 0x40
	Rival         Id = 0x20
	BrackTriangle Id = 0x10
	WhiteTriangle Id = 0x8
	Star          Id = 0x4
	Check         Id = 0x2
	NoMarker      Id = 0x1
)

var originFilterIdMap = map[Id]string{
	All:           "全レース",
	Turf:          "芝",
	Dirt:          "ダート",
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
	Win:           "単勝",
	Place:         "複勝",
	Quinella:      "馬連",
	Exacta:        "馬単",
	QuinellaPlace: "ワイド",
	Trio:          "3連複",
	Trifecta:      "3連単",
	Favorite:      "◎",
	Rival:         "◯",
	BrackTriangle: "▲",
	WhiteTriangle: "△",
	Star:          "☆",
	Check:         "✓",
	NoMarker:      "無",
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

	if i == All {
		return []Id{All}
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
