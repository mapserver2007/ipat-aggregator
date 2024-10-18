package prediction_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
)

type HorseResult struct {
	raceId         types.RaceId
	raceDate       types.RaceDate
	raceName       string
	jockeyId       types.JockeyId
	orderNo        int
	popularNumber  int
	odds           decimal.Decimal
	entries        int
	distance       int
	courseCategory types.CourseCategory
	trackCondition types.TrackCondition
	horseWeight    int
	raceWeight     int
	comment        string
}

func NewHorseResult(
	rawRaceId string,
	rawRaceDate int,
	raceName string,
	jockeyId int,
	orderNo int,
	popularNumber int,
	rawOdds string,
	entries int,
	distance int,
	courseCategoryId int,
	trackConditionId int,
	horseWeight int,
	raceWeight int,
	comment string,
) (*HorseResult, error) {
	odds, err := decimal.NewFromString(rawOdds)
	if err != nil {
		return nil, err
	}

	return &HorseResult{
		raceId:         types.RaceId(rawRaceId),
		raceDate:       types.RaceDate(rawRaceDate),
		raceName:       raceName,
		jockeyId:       types.JockeyId(jockeyId),
		orderNo:        orderNo,
		popularNumber:  popularNumber,
		odds:           odds,
		entries:        entries,
		distance:       distance,
		courseCategory: types.CourseCategory(courseCategoryId),
		trackCondition: types.TrackCondition(trackConditionId),
		horseWeight:    horseWeight,
		raceWeight:     raceWeight,
		comment:        comment,
	}, nil
}

func (h *HorseResult) RaceId() types.RaceId {
	return h.raceId
}

func (h *HorseResult) RaceDate() types.RaceDate {
	return h.raceDate
}

func (h *HorseResult) RaceName() string {
	return h.raceName
}

func (h *HorseResult) JockeyId() types.JockeyId {
	return h.jockeyId
}

func (h *HorseResult) OrderNo() int {
	return h.orderNo
}

func (h *HorseResult) PopularNumber() int {
	return h.popularNumber
}

func (h *HorseResult) Odds() decimal.Decimal {
	return h.odds
}

func (h *HorseResult) Entries() int {
	return h.entries
}

func (h *HorseResult) Distance() int {
	return h.distance
}

func (h *HorseResult) CourseCategory() types.CourseCategory {
	return h.courseCategory
}

func (h *HorseResult) TrackCondition() types.TrackCondition {
	return h.trackCondition
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
