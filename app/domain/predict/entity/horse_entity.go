package entity

type HorseEntity struct {
	horseName     string
	odds          string
	popularNumber int
}

func NewHorseEntity(
	horseName string,
	odds string,
	popularNumber int,
) *HorseEntity {
	return &HorseEntity{
		horseName:     horseName,
		odds:          odds,
		popularNumber: popularNumber,
	}
}

func (h *HorseEntity) HorseName() string {
	return h.horseName
}

func (h *HorseEntity) Odds() string {
	return h.odds
}

func (h *HorseEntity) PopularNumber() int {
	return h.popularNumber
}
