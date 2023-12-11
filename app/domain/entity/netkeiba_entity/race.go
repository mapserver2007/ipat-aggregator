package netkeiba_entity

type Race struct {
	raceName       string
	url            string
	time           string
	startTime      string
	entries        int
	distance       int
	class          int
	courseCategory int
	trackCondition string
	raceResults    []*RaceResult
	payoutResults  []*PayoutResult
}

func NewRace(
	raceName string,
	url string,
	time string,
	startTime string,
	entries int,
	distance int,
	class int,
	courseCategory int,
	trackCondition string,
	raceResults []*RaceResult,
	payoutResults []*PayoutResult,
) *Race {
	return &Race{
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

func (r *Race) RaceName() string {
	return r.raceName
}

func (r *Race) Url() string {
	return r.url
}

func (r *Race) Time() string {
	return r.time
}

func (r *Race) StartTime() string {
	return r.startTime
}

func (r *Race) Entries() int {
	return r.entries
}

func (r *Race) Distance() int {
	return r.distance
}

func (r *Race) Class() int {
	return r.class
}

func (r *Race) CourseCategory() int {
	return r.courseCategory
}

func (r *Race) TrackCondition() string {
	return r.trackCondition
}

func (r *Race) RaceResults() []*RaceResult {
	return r.raceResults
}

func (r *Race) PayoutResults() []*PayoutResult {
	return r.payoutResults
}
