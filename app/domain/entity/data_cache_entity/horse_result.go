package data_cache_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
	"strconv"
)

type HorseResult struct {
	raceId         types.RaceId
	raceDate       types.RaceDate
	raceName       string
	jockeyId       types.JockeyId
	orderNo        int
	popularNumber  int
	horseNumber    types.HorseNumber
	odds           decimal.Decimal
	class          types.GradeClass
	entries        int
	distance       int
	raceCourse     types.RaceCourse
	courseCategory types.CourseCategory
	trackCondition types.TrackCondition
	horseWeight    int
	raceWeight     float64
	comment        string
}

func NewHorseResult(
	rawRaceId string,
	rawRaceDate int,
	raceName string,
	jockeyId string,
	orderNo int,
	popularNumber int,
	horseNumber int,
	rawOdds string,
	rawClass int,
	entries int,
	distance int,
	raceCourseId string,
	courseCategoryId int,
	trackConditionId int,
	horseWeight int,
	rawRaceWeight string,
	comment string,
) (*HorseResult, error) {
	var odds decimal.Decimal
	if rawOdds == "" {
		odds = decimal.Zero
	} else {
		tmp, err := decimal.NewFromString(rawOdds)
		if err != nil {
			return nil, err
		}
		odds = tmp
	}
	var raceWeight float64
	if rawRaceWeight == "" {
		raceWeight = 0
	} else {
		floatValue, err := strconv.ParseFloat(rawRaceWeight, 64)
		if err != nil {
			return nil, err
		}
		raceWeight = floatValue
	}

	return &HorseResult{
		raceId:         types.RaceId(rawRaceId),
		raceDate:       types.RaceDate(rawRaceDate),
		raceName:       raceName,
		jockeyId:       types.JockeyId(jockeyId),
		orderNo:        orderNo,
		popularNumber:  popularNumber,
		horseNumber:    types.HorseNumber(horseNumber),
		odds:           odds,
		class:          types.GradeClass(rawClass),
		entries:        entries,
		distance:       distance,
		raceCourse:     types.RaceCourse(raceCourseId),
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

func (h *HorseResult) HorseNumber() types.HorseNumber {
	return h.horseNumber
}

func (h *HorseResult) Odds() decimal.Decimal {
	return h.odds
}

func (h *HorseResult) Class() types.GradeClass {
	return h.class
}

func (h *HorseResult) Entries() int {
	return h.entries
}

func (h *HorseResult) Distance() int {
	return h.distance
}

func (h *HorseResult) RaceCourse() types.RaceCourse {
	return h.raceCourse
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

func (h *HorseResult) RaceWeight() float64 {
	return h.raceWeight
}

func (h *HorseResult) Comment() string {
	return h.comment
}
