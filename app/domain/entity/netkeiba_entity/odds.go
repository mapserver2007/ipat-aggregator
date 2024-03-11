package netkeiba_entity

type Odds struct {
	odds          string
	popularNumber int
	horseNumber   int
}

func NewOdds(
	odds string,
	popularNumber int,
	horseNumber int,
) *Odds {
	return &Odds{
		odds:          odds,
		popularNumber: popularNumber,
		horseNumber:   horseNumber,
	}
}

func (o *Odds) Odds() string {
	return o.odds
}

func (o *Odds) PopularNumber() int {
	return o.popularNumber
}

func (o *Odds) HorseNumber() int {
	return o.horseNumber
}
