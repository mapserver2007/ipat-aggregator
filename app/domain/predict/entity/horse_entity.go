package entity

type HorseEntity struct {
	HorseName     string
	Odds          string
	PopularNumber int
}

func NewHorseEntity(
	horseName string,
	odds string,
	popularNumber int,
) *HorseEntity {
	return &HorseEntity{
		HorseName:     horseName,
		Odds:          odds,
		PopularNumber: popularNumber,
	}
}
