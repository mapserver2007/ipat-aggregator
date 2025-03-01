package filter

import "sort"

type AttributeId uint64

const (
	All           AttributeId = 0x1FFFFFFFFF
	Turf          AttributeId = 0x1000000000
	Dirt          AttributeId = 0x800000000
	Distance1000m AttributeId = 0x400000000
	Distance1150m AttributeId = 0x200000000
	Distance1200m AttributeId = 0x100000000
	Distance1300m AttributeId = 0x80000000
	Distance1400m AttributeId = 0x40000000
	Distance1500m AttributeId = 0x20000000
	Distance1600m AttributeId = 0x10000000
	Distance1700m AttributeId = 0x8000000
	Distance1800m AttributeId = 0x4000000
	Distance1900m AttributeId = 0x2000000
	Distance2000m AttributeId = 0x1000000
	Distance2100m AttributeId = 0x800000
	Distance2200m AttributeId = 0x400000
	Distance2300m AttributeId = 0x200000
	Distance2400m AttributeId = 0x100000
	Distance2500m AttributeId = 0x80000
	Distance2600m AttributeId = 0x40000
	Distance3000m AttributeId = 0x20000
	Distance3200m AttributeId = 0x10000
	Distance3400m AttributeId = 0x8000
	Distance3600m AttributeId = 0x4000
	Tokyo         AttributeId = 0x2000
	Nakayama      AttributeId = 0x1000
	Kyoto         AttributeId = 0x800
	Hanshin       AttributeId = 0x400
	Niigata       AttributeId = 0x200
	Chukyo        AttributeId = 0x100
	Sapporo       AttributeId = 0x80
	Hakodate      AttributeId = 0x40
	Fukushima     AttributeId = 0x20
	Kokura        AttributeId = 0x10
	GoodToFirm    AttributeId = 0x8
	Good          AttributeId = 0x4
	Yielding      AttributeId = 0x2
	Soft          AttributeId = 0x1
)

var originAttributeIdMap = map[AttributeId]string{
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
}

func (a AttributeId) Value() uint64 {
	return uint64(a)
}

func (a AttributeId) String() string {
	id, _ := originAttributeIdMap[a]
	return id
}

func (a AttributeId) OriginFilters() []AttributeId {
	var ids []AttributeId

	if a == All {
		return []AttributeId{All}
	}

	for id := range originAttributeIdMap {
		if a&id == id {
			ids = append(ids, id)
		}
	}

	sort.Slice(ids, func(i, j int) bool {
		return ids[i] > ids[j]
	})

	return ids
}
