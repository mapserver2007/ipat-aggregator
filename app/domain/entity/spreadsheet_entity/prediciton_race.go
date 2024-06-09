package spreadsheet_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type PredictionRace struct {
	raceId         types.RaceId
	raceName       string
	raceNumber     int
	raceCourseId   types.RaceCourse
	courseCategory types.CourseCategory
	url            string
	filter         filter.Id
}

func NewPredictionRace(
	raceId types.RaceId,
	raceName string,
	raceNumber int,
	raceCourseId types.RaceCourse,
	courseCategory types.CourseCategory,
	url string,
	filter filter.Id,
) *PredictionRace {
	return &PredictionRace{
		raceId:         raceId,
		raceName:       raceName,
		raceNumber:     raceNumber,
		raceCourseId:   raceCourseId,
		courseCategory: courseCategory,
		url:            url,
		filter:         filter,
	}
}

func (p *PredictionRace) RaceId() types.RaceId {
	return p.raceId
}

func (p *PredictionRace) RaceName() string {
	return p.raceName
}

func (p *PredictionRace) RaceNumber() int {
	return p.raceNumber
}

func (p *PredictionRace) RaceCourseId() types.RaceCourse {
	return p.raceCourseId
}

func (p *PredictionRace) CourseCategory() types.CourseCategory {
	return p.courseCategory
}

func (p *PredictionRace) Url() string {
	return p.url
}

func (p *PredictionRace) Filter() filter.Id {
	return p.filter
}