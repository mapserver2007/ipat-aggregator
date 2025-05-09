package netkeiba_entity

type Race struct {
	raceId                string
	raceDate              int
	raceName              string
	raceCourseId          string
	raceNumber            int
	organizer             int
	url                   string
	time                  string
	startTime             string
	entries               int
	distance              int
	class                 int
	courseCategory        int
	trackCondition        int
	raceSexCondition      int
	raceWeightCondition   int
	raceCourseCornerIndex int
	raceAgeCondition      int
	raceEntryHorses       []*RaceEntryHorse
	raceResults           []*RaceResult
	payoutResults         []*PayoutResult
}

func NewRace(
	raceId string,
	raceCourseId string,
	raceNumber int,
	raceDate int,
	raceName string,
	organizer int,
	url string,
	time string,
	startTime string,
	entries int,
	distance int,
	class int,
	courseCategory int,
	trackCondition int,
	raceSexCondition int,
	raceWeightCondition int,
	raceCourseCornerIndex int,
	raceAgeCondition int,
	raceEntryHorses []*RaceEntryHorse,
	raceResults []*RaceResult,
	payoutResults []*PayoutResult,
) *Race {
	return &Race{
		raceId:                raceId,
		raceDate:              raceDate,
		raceName:              raceName,
		raceCourseId:          raceCourseId,
		raceNumber:            raceNumber,
		organizer:             organizer,
		url:                   url,
		time:                  time,
		startTime:             startTime,
		entries:               entries,
		distance:              distance,
		class:                 class,
		courseCategory:        courseCategory,
		trackCondition:        trackCondition,
		raceSexCondition:      raceSexCondition,
		raceWeightCondition:   raceWeightCondition,
		raceCourseCornerIndex: raceCourseCornerIndex,
		raceAgeCondition:      raceAgeCondition,
		raceEntryHorses:       raceEntryHorses,
		raceResults:           raceResults,
		payoutResults:         payoutResults,
	}
}

func (r *Race) RaceId() string {
	return r.raceId
}

func (r *Race) RaceDate() int {
	return r.raceDate
}

func (r *Race) RaceName() string {
	return r.raceName
}

func (r *Race) RaceCourseId() string {
	return r.raceCourseId
}

func (r *Race) RaceNumber() int {
	return r.raceNumber
}

func (r *Race) Organizer() int {
	return r.organizer
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

func (r *Race) TrackCondition() int {
	return r.trackCondition
}

func (r *Race) RaceSexCondition() int {
	return r.raceSexCondition
}

func (r *Race) RaceWeightCondition() int {
	return r.raceWeightCondition
}

func (r *Race) RaceCourseCornerIndex() int {
	return r.raceCourseCornerIndex
}

func (r *Race) RaceAgeCondition() int {
	return r.raceAgeCondition
}

func (r *Race) RaceEntryHorses() []*RaceEntryHorse {
	return r.raceEntryHorses
}

func (r *Race) RaceResults() []*RaceResult {
	return r.raceResults
}

func (r *Race) PayoutResults() []*PayoutResult {
	return r.payoutResults
}
