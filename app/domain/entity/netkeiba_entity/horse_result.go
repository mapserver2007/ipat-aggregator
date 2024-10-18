package netkeiba_entity

type HorseResult struct {
	raceId           string
	raceDate         int
	raceName         string
	jockeyId         int
	orderNo          int
	popularNumber    int
	odds             string
	entries          int
	distance         int
	courseCategoryId int
	trackConditionId int
	horseWeight      int
	raceWeight       int
	comment          string
}

func NewHorseResult(
	raceId string,
	raceDate int,
	raceName string,
	jockeyId int,
	orderNo int,
	popularNumber int,
	odds string,
	entries int,
	distance int,
	courseCategoryId int,
	trackConditionId int,
	horseWeight int,
	raceWeight int,
	comment string,
) *HorseResult {
	return &HorseResult{
		raceId:           raceId,
		raceDate:         raceDate,
		raceName:         raceName,
		jockeyId:         jockeyId,
		orderNo:          orderNo,
		popularNumber:    popularNumber,
		odds:             odds,
		entries:          entries,
		distance:         distance,
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

func (h *HorseResult) JockeyId() int {
	return h.jockeyId
}

func (h *HorseResult) OrderNo() int {
	return h.orderNo
}

func (h *HorseResult) PopularNumber() int {
	return h.popularNumber
}

func (h *HorseResult) Odds() string {
	return h.odds
}

func (h *HorseResult) Entries() int {
	return h.entries
}

func (h *HorseResult) Distance() int {
	return h.distance
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

func (h *HorseResult) RaceWeight() int {
	return h.raceWeight
}

func (h *HorseResult) Comment() string {
	return h.comment
}
