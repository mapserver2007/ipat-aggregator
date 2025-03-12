package netkeiba_entity

type RaceResult struct {
	orderNo        int
	horseId        string
	horseName      string
	bracketNumber  int
	horseNumber    int
	jockeyId       string
	odds           string
	popularNumber  int
	jockeyWeight   string
	horseWeight    int
	horseWeightAdd int
}

func NewRaceResult(
	orderNo int,
	horseId string,
	horseName string,
	bracketNumber int,
	horseNumber int,
	jockeyId string,
	odds string,
	popularNumber int,
	jockeyWeight string,
	horseWeight int,
	horseWeightAdd int,
) *RaceResult {
	return &RaceResult{
		orderNo:        orderNo,
		horseId:        horseId,
		horseName:      horseName,
		bracketNumber:  bracketNumber,
		horseNumber:    horseNumber,
		jockeyId:       jockeyId,
		odds:           odds,
		popularNumber:  popularNumber,
		jockeyWeight:   jockeyWeight,
		horseWeight:    horseWeight,
		horseWeightAdd: horseWeightAdd,
	}
}

func (r *RaceResult) OrderNo() int {
	return r.orderNo
}

func (r *RaceResult) HorseId() string {
	return r.horseId
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

func (r *RaceResult) JockeyId() string {
	return r.jockeyId
}

func (r *RaceResult) Odds() string {
	return r.odds
}

func (r *RaceResult) PopularNumber() int {
	return r.popularNumber
}

func (r *RaceResult) JockeyWeight() string {
	return r.jockeyWeight
}

func (r *RaceResult) HorseWeight() int {
	return r.horseWeight
}

func (r *RaceResult) HorseWeightAdd() int {
	return r.horseWeightAdd
}
