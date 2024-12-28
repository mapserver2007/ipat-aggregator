package list_entity

import "github.com/shopspring/decimal"

type Horse struct {
	horseName     string
	odds          decimal.Decimal
	popularNumber int
}

func NewHorse(
	horseName string,
	odds decimal.Decimal,
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

func (h *Horse) Odds() decimal.Decimal {
	return h.odds
}

func (h *Horse) PopularNumber() int {
	return h.popularNumber
}
