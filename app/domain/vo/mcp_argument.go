package vo

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type MCPArgument struct {
	raceNames        []string
	raceDates        []types.RaceDate
	raceCourses      []types.RaceCourse
	courseCategories []types.CourseCategory
	distances        []int
	odds             []float64
	ticketTypes      []types.TicketType
}

func NewMCPArgument(
	raceNames []string,
	raceDates []types.RaceDate,
	raceCourses []types.RaceCourse,
	courseCategories []types.CourseCategory,
	distances []int,
	odds []float64,
	ticketTypes []types.TicketType,
) *MCPArgument {
	return &MCPArgument{
		raceNames:        raceNames,
		raceDates:        raceDates,
		raceCourses:      raceCourses,
		courseCategories: courseCategories,
		distances:        distances,
		odds:             odds,
		ticketTypes:      ticketTypes,
	}
}

func (m *MCPArgument) RaceNames() []string {
	return m.raceNames
}

func (m *MCPArgument) RaceDates() []types.RaceDate {
	return m.raceDates
}

func (m *MCPArgument) RaceCourses() []types.RaceCourse {
	return m.raceCourses
}

func (m *MCPArgument) CourseCategories() []types.CourseCategory {
	return m.courseCategories
}

func (m *MCPArgument) Distances() []int {
	return m.distances
}

func (m *MCPArgument) Odds() []float64 {
	return m.odds
}

func (m *MCPArgument) TicketTypes() []types.TicketType {
	return m.ticketTypes
}
