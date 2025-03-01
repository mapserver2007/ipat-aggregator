package spreadsheet_entity

import (
	"fmt"

	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
)

type AnalysisPlaceUnhit struct {
	raceUrl             string
	raceDate            types.RaceDate
	raceNumber          int
	raceCourse          types.RaceCourse
	raceName            string
	class               types.GradeClass
	courseCategory      types.CourseCategory
	distance            int
	raceWeightCondition types.RaceWeightCondition
	trackCondition      types.TrackCondition
	entries             int
	horseNumber         types.HorseNumber
	horseUrl            string
	horseName           string
	jockeyUrl           string
	jockeyName          string
	popularNumber       int
	odds                decimal.Decimal
	orderNo             int
}

func NewAnalysisPlaceUnhit(
	raceId types.RaceId,
	raceDate types.RaceDate,
	raceNumber int,
	raceCourse types.RaceCourse,
	raceName string,
	class types.GradeClass,
	courseCategory types.CourseCategory,
	distance int,
	raceWeightCondition types.RaceWeightCondition,
	trackCondition types.TrackCondition,
	entries int,
	horseNumber types.HorseNumber,
	horseId types.HorseId,
	horseName string,
	jockeyId types.JockeyId,
	jockeyName string,
	popularNumber int,
	odds decimal.Decimal,
	orderNo int,
) *AnalysisPlaceUnhit {
	return &AnalysisPlaceUnhit{
		raceUrl:             fmt.Sprintf("https://race.netkeiba.com/race/shutuba.html?race_id=%s", raceId.String()),
		raceDate:            raceDate,
		raceNumber:          raceNumber,
		raceCourse:          raceCourse,
		raceName:            raceName,
		class:               class,
		courseCategory:      courseCategory,
		distance:            distance,
		raceWeightCondition: raceWeightCondition,
		trackCondition:      trackCondition,
		entries:             entries,
		horseNumber:         horseNumber,
		horseUrl:            fmt.Sprintf("https://db.netkeiba.com/horse/%s", horseId),
		horseName:           horseName,
		jockeyUrl:           fmt.Sprintf("https://db.netkeiba.com/jockey/%s", jockeyId),
		jockeyName:          jockeyName,
		popularNumber:       popularNumber,
		odds:                odds,
		orderNo:             orderNo,
	}
}

func (a *AnalysisPlaceUnhit) RaceUrl() string {
	return a.raceUrl
}

func (a *AnalysisPlaceUnhit) RaceDate() types.RaceDate {
	return a.raceDate
}

func (a *AnalysisPlaceUnhit) RaceNumber() int {
	return a.raceNumber
}

func (a *AnalysisPlaceUnhit) RaceCourse() types.RaceCourse {
	return a.raceCourse
}

func (a *AnalysisPlaceUnhit) RaceName() string {
	return a.raceName
}

func (a *AnalysisPlaceUnhit) Class() types.GradeClass {
	return a.class
}

func (a *AnalysisPlaceUnhit) CourseCategory() types.CourseCategory {
	return a.courseCategory
}

func (a *AnalysisPlaceUnhit) Distance() int {
	return a.distance
}

func (a *AnalysisPlaceUnhit) RaceWeightCondition() types.RaceWeightCondition {
	return a.raceWeightCondition
}

func (a *AnalysisPlaceUnhit) TrackCondition() types.TrackCondition {
	return a.trackCondition
}

func (a *AnalysisPlaceUnhit) Entries() int {
	return a.entries
}

func (a *AnalysisPlaceUnhit) HorseNumber() types.HorseNumber {
	return a.horseNumber
}

func (a *AnalysisPlaceUnhit) HorseUrl() string {
	return a.horseUrl
}

func (a *AnalysisPlaceUnhit) HorseName() string {
	return a.horseName
}

func (a *AnalysisPlaceUnhit) JockeyUrl() string {
	return a.jockeyUrl
}

func (a *AnalysisPlaceUnhit) JockeyName() string {
	return a.jockeyName
}

func (a *AnalysisPlaceUnhit) PopularNumber() int {
	return a.popularNumber
}

func (a *AnalysisPlaceUnhit) Odds() decimal.Decimal {
	return a.odds
}

func (a *AnalysisPlaceUnhit) OrderNo() int {
	return a.orderNo
}
