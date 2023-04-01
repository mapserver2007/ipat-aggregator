package raw_entity

type RawRaceNetkeiba struct {
	raceName       string
	url            string
	time           string
	startTime      string
	entries        int
	distance       int
	class          int
	courseCategory int
	trackCondition string
	raceResults    []*RawRaceResultNetkeiba
	payoutResults  []*RawPayoutResultNetkeiba
}

type RawRaceResultNetkeiba struct {
	orderNo       int
	horseName     string
	bracketNumber int
	horseNumber   int
	odds          string
	popularNumber int
}

type RawPayoutResultNetkeiba struct {
	ticketType int
	numbers    []string
	odds       []string
}

type RawRacingNumberNetkeiba struct {
	date         int
	round        int
	day          int
	raceCourseId int
}

func NewRawRaceNetkeiba(
	raceName string,
	url string,
	time string,
	startTime string,
	entries int,
	distance int,
	class int,
	courseCategory int,
	trackCondition string,
	raceResults []*RawRaceResultNetkeiba,
	payoutResults []*RawPayoutResultNetkeiba,
) *RawRaceNetkeiba {
	return &RawRaceNetkeiba{
		raceName:       raceName,
		url:            url,
		time:           time,
		startTime:      startTime,
		entries:        entries,
		distance:       distance,
		class:          class,
		courseCategory: courseCategory,
		trackCondition: trackCondition,
		raceResults:    raceResults,
		payoutResults:  payoutResults,
	}
}

func (r *RawRaceNetkeiba) RaceName() string {
	return r.raceName
}

func (r *RawRaceNetkeiba) Url() string {
	return r.url
}

func (r *RawRaceNetkeiba) Time() string {
	return r.time
}

func (r *RawRaceNetkeiba) StartTime() string {
	return r.startTime
}

func (r *RawRaceNetkeiba) Entries() int {
	return r.entries
}

func (r *RawRaceNetkeiba) Distance() int {
	return r.distance
}

func (r *RawRaceNetkeiba) Class() int {
	return r.class
}

func (r *RawRaceNetkeiba) CourseCategory() int {
	return r.courseCategory
}

func (r *RawRaceNetkeiba) TrackCondition() string {
	return r.trackCondition
}

func (r *RawRaceNetkeiba) RaceResults() []*RawRaceResultNetkeiba {
	return r.raceResults
}

func (r *RawRaceNetkeiba) PayoutResults() []*RawPayoutResultNetkeiba {
	return r.payoutResults
}

func NewRawRaceResultNetkeiba(
	orderNo int,
	horseName string,
	bracketNumber int,
	horseNumber int,
	odds string,
	popularNumber int,
) *RawRaceResultNetkeiba {
	return &RawRaceResultNetkeiba{
		orderNo:       orderNo,
		horseName:     horseName,
		bracketNumber: bracketNumber,
		horseNumber:   horseNumber,
		odds:          odds,
		popularNumber: popularNumber,
	}
}

func (r *RawRaceResultNetkeiba) OrderNo() int {
	return r.orderNo
}

func (r *RawRaceResultNetkeiba) HorseName() string {
	return r.horseName
}

func (r *RawRaceResultNetkeiba) BracketNumber() int {
	return r.bracketNumber
}

func (r *RawRaceResultNetkeiba) HorseNumber() int {
	return r.horseNumber
}

func (r *RawRaceResultNetkeiba) Odds() string {
	return r.odds
}

func (r *RawRaceResultNetkeiba) PopularNumber() int {
	return r.popularNumber
}

func NewRawPayoutResultNetkeiba(
	ticketType int,
	numbers []string,
	odds []string,
) *RawPayoutResultNetkeiba {
	return &RawPayoutResultNetkeiba{
		ticketType: ticketType,
		numbers:    numbers,
		odds:       odds,
	}
}

func (r *RawPayoutResultNetkeiba) TicketType() int {
	return r.ticketType
}

func (r *RawPayoutResultNetkeiba) Numbers() []string {
	return r.numbers
}

func (r *RawPayoutResultNetkeiba) Odds() []string {
	return r.odds
}

func NewRawRacingNumberNetkeiba(
	date int,
	round int,
	day int,
	raceCourseId int,
) *RawRacingNumberNetkeiba {
	return &RawRacingNumberNetkeiba{
		date:         date,
		round:        round,
		day:          day,
		raceCourseId: raceCourseId,
	}
}

func (r *RawRacingNumberNetkeiba) Date() int {
	return r.date
}

func (r *RawRacingNumberNetkeiba) Round() int {
	return r.round
}

func (r *RawRacingNumberNetkeiba) Day() int {
	return r.day
}

func (r *RawRacingNumberNetkeiba) RaceCourseId() int {
	return r.raceCourseId
}
