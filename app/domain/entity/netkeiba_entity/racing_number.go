package netkeiba_entity

type RacingNumber struct {
	date         int
	round        int
	day          int
	raceCourseId string
}

func NewRacingNumber(
	date int,
	round int,
	day int,
	raceCourseId string,
) *RacingNumber {
	return &RacingNumber{
		date:         date,
		round:        round,
		day:          day,
		raceCourseId: raceCourseId,
	}
}

func (r *RacingNumber) Date() int {
	return r.date
}

func (r *RacingNumber) Round() int {
	return r.round
}

func (r *RacingNumber) Day() int {
	return r.day
}

func (r *RacingNumber) RaceCourseId() string {
	return r.raceCourseId
}
