package entity

type Horse struct {
	horseName     string
	odds          string
	popularNumber int
}

func NewHorse(
	horseName string,
	odds string,
	popularNumber int,
) *Horse {
	return &Horse{
		horseName:     horseName,
		odds:          odds,
		popularNumber: popularNumber,
	}
}

func (h *Horse) HorseName() string {
	return h.horseName
}

func (h *Horse) Odds() string {
	return h.odds
}

func (h *Horse) PopularNumber() int {
	return h.popularNumber
}
