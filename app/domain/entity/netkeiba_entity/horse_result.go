package netkeiba_entity

type HorseResult struct {
	raceId           string
	raceDate         int
	raceName         string
	jockeyId         string
	orderNo          int
	popularNumber    int
	horseNumber      int
	odds             string
	class            int
	entries          int
	distance         int
	raceCourseId     string
	courseCategoryId int
	trackConditionId int
	horseWeight      int
	raceWeight       float64
	comment          string
}

func NewHorseResult(
	raceId string,
	raceDate int,
	raceName string,
	jockeyId string,
	orderNo int,
	popularNumber int,
	horseNumber int,
	odds string,
	class int,
	entries int,
	distance int,
	raceCourseId string,
	courseCategoryId int,
	trackConditionId int,
	horseWeight int,
	raceWeight float64,
	comment string,
) *HorseResult {
	return &HorseResult{
		raceId:           raceId,
		raceDate:         raceDate,
		raceName:         raceName,
		jockeyId:         jockeyId,
		orderNo:          orderNo,
		popularNumber:    popularNumber,
		horseNumber:      horseNumber,
		odds:             odds,
		class:            class,
		entries:          entries,
		distance:         distance,
		raceCourseId:     raceCourseId,
		courseCategoryId: courseCategoryId,
		trackConditionId: trackConditionId,
		horseWeight:      horseWeight,
		raceWeight:       raceWeight,
		comment:          comment,
	}
}

func (h *HorseResult) RaceId() string {
	return h.raceId
}

func (h *HorseResult) RaceDate() int {
	return h.raceDate
}

func (h *HorseResult) RaceName() string {
	return h.raceName
}

func (h *HorseResult) JockeyId() string {
	return h.jockeyId
}

func (h *HorseResult) OrderNo() int {
	return h.orderNo
}

func (h *HorseResult) PopularNumber() int {
	return h.popularNumber
}

func (h *HorseResult) HorseNumber() int {
	return h.horseNumber
}

func (h *HorseResult) Odds() string {
	return h.odds
}

func (h *HorseResult) Class() int {
	return h.class
}

func (h *HorseResult) Entries() int {
	return h.entries
}

func (h *HorseResult) Distance() int {
	return h.distance
}

func (h *HorseResult) RaceCourseId() string {
	return h.raceCourseId
}

func (h *HorseResult) CourseCategoryId() int {
	return h.courseCategoryId
}

func (h *HorseResult) TrackConditionId() int {
	return h.trackConditionId
}

func (h *HorseResult) HorseWeight() int {
	return h.horseWeight
}

func (h *HorseResult) RaceWeight() float64 {
	return h.raceWeight
}

func (h *HorseResult) Comment() string {
	return h.comment
}
