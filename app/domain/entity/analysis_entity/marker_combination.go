package analysis_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type Layer1 struct {
	MarkerCombination map[types.MarkerCombinationId]*Layer2
}

type Layer2 struct {
	RaceDate map[types.RaceDate]*Layer3
}

type Layer3 struct {
	// 同じレースでの印の組み合わせ(ワイド、複勝)がmarkerCombinationIdに対して複数出る場合があるのでsliceで保持
	RaceId map[types.RaceId][]*Calculable
}

type MarkerCombinationOrder struct {
}
