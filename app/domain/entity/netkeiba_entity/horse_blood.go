package netkeiba_entity

type HorseBlood struct {
	sireId          string
	broodmareSireId string
}

func NewHorseBlood(
	sireId string,
	broodmareSireId string,
) *HorseBlood {
	return &HorseBlood{
		sireId:          sireId,
		broodmareSireId: broodmareSireId,
	}
}

func (h *HorseBlood) SireId() string {
	return h.sireId
}

func (h *HorseBlood) BroodmareSireId() string {
	return h.broodmareSireId
}
