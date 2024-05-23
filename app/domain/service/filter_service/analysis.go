package filter_service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type AnalysisFilter interface {
	Get(ctx context.Context) []filter.Id
	Create(ctx context.Context, race *data_cache_entity.Race) []filter.Id
}

type filterService struct{}

func NewAnalysisFilter() AnalysisFilter {
	return &filterService{}
}

func (f *filterService) Get(ctx context.Context) []filter.Id {
	return []filter.Id{
		filter.All,
		filter.TurfShortDistance1,
		filter.TurfShortDistance2,
		filter.TurfShortDistance3,
		filter.TurfMiddleDistance1,
		filter.TurfMiddleDistance2,
		filter.TurfLongDistance,
		filter.DirtShortDistance1,
		filter.DirtShortDistance2,
		filter.DirtShortDistance3,
		filter.DirtMiddleDistance1,
		filter.DirtMiddleDistance2,
		filter.DirtLongDistance,
		filter.TurfClass1,
		filter.DirtClass1,
		filter.TurfClass6,
		filter.DirtClass6,
		filter.TurfLargeNumberOfHorses,
		filter.TurfSmallNumberOfHorses,
		filter.DirtLargeNumberOfHorses,
		filter.DirtSmallNumberOfHorses,
	}
}

func (f *filterService) Create(ctx context.Context, race *data_cache_entity.Race) []filter.Id {
	var filterIds []filter.Id
	filterIds = append(filterIds, createCourseCategoryFilter(race.CourseCategory())...)
	filterIds = append(filterIds, createDistanceFilter(race.Distance())...)
	filterIds = append(filterIds, createGradeClassFilter(race.Class())...)
	filterIds = append(filterIds, createTrackConditionFilter(race.TrackCondition())...)
	filterIds = append(filterIds, createRaceCourseFilter(race.RaceCourseId())...)
	filterIds = append(filterIds, createEntriesFilter(race.Entries())...)

	return filterIds
}

func createCourseCategoryFilter(courseCategory types.CourseCategory) []filter.Id {
	var filterIds []filter.Id
	switch courseCategory {
	case types.Turf:
		filterIds = append(filterIds, filter.Turf)
	case types.Dirt:
		filterIds = append(filterIds, filter.Dirt)
	}
	return filterIds
}

func createGradeClassFilter(class types.GradeClass) []filter.Id {
	var filterIds []filter.Id
	switch class {
	case types.Grade1, types.Grade2, types.Grade3:
		filterIds = append(filterIds, filter.Class6)
	case types.OpenClass, types.ListedClass:
		filterIds = append(filterIds, filter.Class5)
	case types.ThreeWinClass:
		filterIds = append(filterIds, filter.Class4)
	case types.TwoWinClass:
		filterIds = append(filterIds, filter.Class3)
	case types.OneWinClass:
		filterIds = append(filterIds, filter.Class2)
	case types.Maiden:
		filterIds = append(filterIds, filter.Class1)
	}
	return filterIds
}

func createGradeClassSimpleFilter(class types.GradeClass) []filter.Id {
	var filterIds []filter.Id
	switch class {
	case types.Grade1, types.Grade2, types.Grade3, types.OpenClass, types.ListedClass:
		filterIds = append(filterIds, filter.Class56)
	case types.ThreeWinClass, types.TwoWinClass, types.OneWinClass:
		filterIds = append(filterIds, filter.Class234)
	case types.Maiden:
		filterIds = append(filterIds, filter.Class1)
	}
	return filterIds
}

func createDistanceFilter(distance int) []filter.Id {
	var filterIds []filter.Id
	if distance >= 1000 && distance <= 1200 {
		filterIds = append(filterIds, filter.ShortDistance1)
	} else if distance >= 1201 && distance <= 1400 {
		filterIds = append(filterIds, filter.ShortDistance2)
	} else if distance >= 1401 && distance <= 1600 {
		filterIds = append(filterIds, filter.ShortDistance3)
	} else if distance >= 1601 && distance <= 1700 {
		filterIds = append(filterIds, filter.MiddleDistance1)
	} else if distance >= 1701 && distance <= 1800 {
		filterIds = append(filterIds, filter.MiddleDistance2)
	} else if distance >= 1801 && distance <= 2000 {
		filterIds = append(filterIds, filter.MiddleDistance3)
	} else if distance >= 2001 {
		filterIds = append(filterIds, filter.LongDistance)
	}
	return filterIds
}

func createDistanceSimpleFilter(distance int) []filter.Id {
	var filterIds []filter.Id
	if distance >= 1000 && distance <= 1600 {
		filterIds = append(filterIds, filter.ShortDistance)
	} else if distance >= 1601 && distance <= 2000 {
		filterIds = append(filterIds, filter.MiddleDistance)
	} else if distance >= 2001 {
		filterIds = append(filterIds, filter.LongDistance)
	}
	return filterIds
}

func createTrackConditionFilter(trackCondition types.TrackCondition) []filter.Id {
	var filterIds []filter.Id
	switch trackCondition {
	case types.GoodToFirm:
		filterIds = append(filterIds, filter.GoodTrack)
	case types.Good, types.Yielding, types.Soft:
		filterIds = append(filterIds, filter.BadTrack)
	}
	return filterIds
}

func createRaceCourseFilter(raceCourseId types.RaceCourse) []filter.Id {
	var filterIds []filter.Id
	switch raceCourseId {
	case types.Tokyo, types.Nakayama, types.Hanshin, types.Kyoto:
		filterIds = append(filterIds, filter.CentralCourse)
	default:
		filterIds = append(filterIds, filter.LocalCourse)
	}
	return filterIds
}

func createEntriesFilter(entries int) []filter.Id {
	var filterIds []filter.Id
	if entries <= 10 {
		filterIds = append(filterIds, filter.SmallNumberOfHorses)
	} else {
		filterIds = append(filterIds, filter.LargeNumberOfHorses)
	}
	return filterIds
}
