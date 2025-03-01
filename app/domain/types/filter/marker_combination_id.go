package filter

import "sort"

type MarkerCombinationId uint64

const (
	MarkerCombinationTurf          MarkerCombinationId = 0x2
	MarkerCombinationDirt          MarkerCombinationId = 0x1
	MarkerCombinationWin           MarkerCombinationId = 0x4000000
	MarkerCombinationPlace         MarkerCombinationId = 0x2000000
	MarkerCombinationQuinellaPlace MarkerCombinationId = 0x1000000
	MarkerCombinationQuinella      MarkerCombinationId = 0x800000
	MarkerCombinationExacta        MarkerCombinationId = 0x400000
	MarkerCombinationTrio          MarkerCombinationId = 0x200000
	MarkerCombinationTrifecta      MarkerCombinationId = 0x100000
	MarkerCombinationFavorite      MarkerCombinationId = 0x80000
	MarkerCombinationRival         MarkerCombinationId = 0x40000
	MarkerCombinationBrackTriangle MarkerCombinationId = 0x20000
	MarkerCombinationWhiteTriangle MarkerCombinationId = 0x10000
	MarkerCombinationStar          MarkerCombinationId = 0x8000
	MarkerCombinationCheck         MarkerCombinationId = 0x4000
)

var originMarkerCombinationIdMap = map[MarkerCombinationId]string{
	MarkerCombinationTurf:          "芝",
	MarkerCombinationDirt:          "ダート",
	MarkerCombinationWin:           "単勝",
	MarkerCombinationPlace:         "複勝",
	MarkerCombinationQuinellaPlace: "ワイド",
	MarkerCombinationQuinella:      "馬連",
	MarkerCombinationExacta:        "馬単",
	MarkerCombinationTrio:          "三連複",
	MarkerCombinationTrifecta:      "三連単",
	MarkerCombinationFavorite:      "◎",
	MarkerCombinationRival:         "◯",
	MarkerCombinationBrackTriangle: "▲",
	MarkerCombinationWhiteTriangle: "△",
	MarkerCombinationStar:          "☆",
	MarkerCombinationCheck:         "✓",
}

func (m MarkerCombinationId) Value() uint64 {
	return uint64(m)
}

func (m MarkerCombinationId) String() string {
	return originMarkerCombinationIdMap[m]
}

func (m MarkerCombinationId) OriginFilters() []MarkerCombinationId {
	var ids []MarkerCombinationId

	for id := range originMarkerCombinationIdMap {
		if m&id == id {
			ids = append(ids, id)
		}
	}

	sort.Slice(ids, func(i, j int) bool {
		return ids[i] > ids[j]
	})

	return ids
}
