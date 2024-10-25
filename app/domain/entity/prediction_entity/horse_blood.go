package prediction_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type HorseBlood struct {
	sireId          types.HorseId
	broodmareSireId types.HorseId
}

func NewHorseBlood(
	rawSireId string,
	rawBroodmareSireId string,
) *HorseBlood {
	return &HorseBlood{
		sireId:          types.HorseId(rawSireId),
		broodmareSireId: types.HorseId(rawBroodmareSireId),
	}
}

func (h *HorseBlood) SireId() types.HorseId {
	return h.sireId
}

func (h *HorseBlood) BroodmareSireId() types.HorseId {
	return h.broodmareSireId
}
