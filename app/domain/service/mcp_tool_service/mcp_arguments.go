package mcp_tool_service

import (
	"fmt"

	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/vo"
	"github.com/mark3labs/mcp-go/mcp"
)

const (
	toolNameRaceName       = "raceNames"
	toolNameRaceDate       = "raceDate"
	toolNameRaceCourse     = "raceCourse"
	toolNameCourseCategory = "courseCategory"
	toolNameDistance       = "distance"
	toolOdds               = "odds"
	toolTicketType         = "ticketType"
)

func LoadMCPArgument(request mcp.CallToolRequest) (*vo.MCPArgument, error) {
	rawRaceNames := request.GetStringSlice(toolNameRaceDate, []string{})
	raceNames := make([]string, 0)
	for _, name := range rawRaceNames {
		raceNames = append(raceNames, name)
	}
	rawRaceDates := request.GetIntSlice(toolNameRaceDate, []int{})
	raceDates := make([]types.RaceDate, 0)
	for _, i := range rawRaceDates {
		raceDate, err := types.NewRaceDate(fmt.Sprintf("%06d", i))
		if err != nil {
			continue
		}
		raceDates = append(raceDates, raceDate)
	}

	rawRaceCourses := request.GetStringSlice(toolNameRaceCourse, []string{})
	raceCourses := make([]types.RaceCourse, 0)
	for _, name := range rawRaceCourses {
		raceCourses = append(raceCourses, types.NewRaceCourse(name))
	}

	rawCourseCategories := request.GetStringSlice(toolNameCourseCategory, []string{})
	courseCategories := make([]types.CourseCategory, 0)
	for _, name := range rawCourseCategories {
		courseCategories = append(courseCategories, types.NewCourseCategory(name))
	}

	distances := request.GetIntSlice(toolNameDistance, []int{})
	odds := request.GetFloatSlice(toolOdds, []float64{})
	rawTicketTypes := request.GetStringSlice(toolTicketType, []string{})
	ticketTypes := make([]types.TicketType, 0)
	for _, name := range rawTicketTypes {
		ticketTypes = append(ticketTypes, types.NewTicketType(name))
	}

	return vo.NewMCPArgument(
		raceNames,
		raceDates,
		raceCourses,
		courseCategories,
		distances,
		odds,
		ticketTypes,
	), nil
}

func GetToolOps() []mcp.ToolOption {
	return []mcp.ToolOption{
		getRaceNameToolOps(),
		getRaceDateToolOps(),
		getRaceCourseToolOps(),
		getCourseCategoryToolOps(),
		getDistanceToolOps(),
		getOddsToolOps(),
		getTicketTypeToolOps(),
	}
}

func getRaceNameToolOps() mcp.ToolOption {
	return mcp.WithArray(toolNameRaceName,
		mcp.Description("Name of race (multiple designations possible)"),
		mcp.Items(map[string]any{"type": "string"}),
	)
}

func getRaceDateToolOps() mcp.ToolOption {
	return mcp.WithArray(toolNameRaceDate,
		mcp.Description("Date of race (multiple designations possible)"),
		mcp.Items(map[string]any{"type": "integer"}),
	)
}

func getRaceCourseToolOps() mcp.ToolOption {
	return mcp.WithArray(toolNameRaceCourse,
		mcp.Description("Race course name (multiple designations possible)"),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Enum(
			types.Tokyo.Name(),
			types.Nakayama.Name(),
			types.Hanshin.Name(),
			types.Kyoto.Name(),
			types.Chukyo.Name(),
			types.Kokura.Name(),
			types.Niigata.Name(),
			types.Hakodate.Name(),
			types.Sapporo.Name(),
		),
	)
}

func getCourseCategoryToolOps() mcp.ToolOption {
	return mcp.WithString(toolNameCourseCategory,
		mcp.Description("Course category name"),
		mcp.Enum(types.Turf.String(), types.Dirt.String()),
	)
}

func getDistanceToolOps() mcp.ToolOption {
	return mcp.WithNumber(toolNameDistance,
		mcp.Description("distance"),
	)
}

func getOddsToolOps() mcp.ToolOption {
	return mcp.WithArray(toolOdds,
		mcp.Description("odds"),
		mcp.Items(map[string]any{"type": "number"}),
		mcp.Min(1.0),
	)
}

func getTicketTypeToolOps() mcp.ToolOption {
	return mcp.WithArray(toolTicketType,
		mcp.Description("ticket type"),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Enum(
			types.Win.Name(),
			types.Place.Name(),
			types.BracketQuinella.Name(),
			types.Quinella.Name(),
			types.QuinellaWheel.Name(),
			types.Exacta.Name(),
			types.ExactaWheelOfFirst.Name(),
			types.QuinellaPlace.Name(),
			types.QuinellaPlaceWheel.Name(),
			types.QuinellaPlaceFormation.Name(),
			types.Trio.Name(),
			types.TrioFormation.Name(),
			types.TrioWheelOfFirst.Name(),
			types.TrioWheelOfSecond.Name(),
			types.TrioBox.Name(),
			types.Trifecta.Name(),
			types.TrifectaFormation.Name(),
			types.TrifectaWheelOfFirst.Name(),
			types.TrifectaWheelOfSecond.Name(),
			types.TrifectaWheelOfFirstMulti.Name(),
			types.TrifectaWheelOfSecondMulti.Name(),
		),
	)
}
