package entity

type RaceResult struct {
	orderNo       int
	horseName     string
	bracketNumber int
	horseNumber   int
	odds          string
	popularNumber int
}

func NewRaceResult(
	orderNo int,
	horseName string,
	bracketNumber int,
	horseNumber int,
	odds string,
	popularNumber int,
) *RaceResult {
	return &RaceResult{
		orderNo:       orderNo,
		horseName:     horseName,
		bracketNumber: bracketNumber,
		horseNumber:   horseNumber,
		odds:          odds,
		popularNumber: popularNumber,
	}
}

func (r *RaceResult) OrderNo() int {
	return r.orderNo
}

func (r *RaceResult) HorseName() string {
	return r.horseName
}

func (r *RaceResult) BracketNumber() int {
	return r.bracketNumber
}

func (r *RaceResult) HorseNumber() int {
	return r.horseNumber
}

func (r *RaceResult) Odds() string {
	return r.odds
}

func (r *RaceResult) PopularNumber() int {
	return r.popularNumber
}
