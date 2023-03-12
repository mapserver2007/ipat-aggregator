package entity

type RacingNumber struct {
	date         int
	round        int
	day          int
	raceCourseId int
}

func NewRacingNumber(
	date int,
	round int,
	day int,
	raceCourseId int,
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

func (r *RacingNumber) RaceCourseId() int {
	return r.raceCourseId
}
