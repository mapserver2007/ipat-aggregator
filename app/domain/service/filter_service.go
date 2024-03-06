package service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"strings"
)

type FilterService interface {
	CreateAnalysisFilters(ctx context.Context, race *data_cache_entity.Race, raceResultByMarker *data_cache_entity.RaceResult) []filter.Id
	CreatePredictionFilters(ctx context.Context, race *prediction_entity.Race) (filter.Id, filter.Id)
	GetAnalysisFilters() []filter.Id
}

type filterService struct {
}

func NewFilterService() FilterService {
	return &filterService{}
}

func (f *filterService) CreateAnalysisFilters(
	ctx context.Context,
	race *data_cache_entity.Race,
	raceResultByMarker *data_cache_entity.RaceResult,
) []filter.Id {
	var filterIds []filter.Id
	filterIds = append(filterIds, f.createCourseCategoryFilter(ctx, race.CourseCategory())...)
	filterIds = append(filterIds, f.createDistanceFilter(ctx, race.Distance())...)
	filterIds = append(filterIds, f.createGradeClassFilter(ctx, race.Class())...)
	filterIds = append(filterIds, f.createTrackConditionFilter(ctx, race.TrackCondition())...)
	filterIds = append(filterIds, f.createRaceCourseFilter(ctx, race.RaceCourseId())...)
	filterIds = append(filterIds, f.createEntriesFilter(ctx, race.Entries())...)

	switch raceResultByMarker.JockeyId() {
	case 5339, 1088, 5366, 5509, 5585: // C.ルメール, 川田将雅, R.ムーア, J.モレイラ, D.レーン
		filterIds = append(filterIds, filter.TopJockey)
	default:
		filterIds = append(filterIds, filter.OtherJockey)
	}

	return filterIds
}

func (f *filterService) CreatePredictionFilters(
	ctx context.Context,
	race *prediction_entity.Race,
) (filter.Id, filter.Id) {
	var (
		strictFilterIds, simpleFilterIds     []filter.Id
		strictFilterId, simpleFilterId       filter.Id
		strictFilterNames, simpleFilterNames []string
	)
	strictFilterIds = append(strictFilterIds, f.createCourseCategoryFilter(ctx, race.CourseCategory())...)
	strictFilterIds = append(strictFilterIds, f.createDistanceFilter(ctx, race.Distance())...)
	strictFilterIds = append(strictFilterIds, f.createGradeClassFilter(ctx, race.Class())...)
	strictFilterIds = append(strictFilterIds, f.createTrackConditionFilter(ctx, race.TrackCondition())...)
	strictFilterIds = append(strictFilterIds, f.createRaceCourseFilter(ctx, race.RaceCourseId())...)
	strictFilterIds = append(strictFilterIds, f.createEntriesFilter(ctx, race.Entries())...)
	for _, filterId := range strictFilterIds {
		strictFilterNames = append(strictFilterNames, filterId.String())
	}
	strictFilterIds = append(strictFilterIds, []filter.Id{filter.TopJockey, filter.OtherJockey}...)
	for _, filterId := range strictFilterIds {
		strictFilterId = strictFilterId | filterId
	}

	simpleFilterIds = append(simpleFilterIds, f.createCourseCategoryFilter(ctx, race.CourseCategory())...)
	simpleFilterIds = append(simpleFilterIds, f.createDistanceSimpleFilter(ctx, race.Distance())...)
	simpleFilterIds = append(simpleFilterIds, f.createGradeClassSimpleFilter(ctx, race.Class())...)
	for _, filterId := range simpleFilterIds {
		simpleFilterNames = append(simpleFilterNames, filterId.String())
	}
	simpleFilterIds = append(strictFilterIds,
		[]filter.Id{filter.GoodTrack, filter.BadTrack, filter.CentralCourse, filter.LocalCourse, filter.TopJockey, filter.OtherJockey, filter.SmallNumberOfHorses, filter.LargeNumberOfHorses}...)
	for _, filterId := range simpleFilterIds {
		simpleFilterId = simpleFilterId | filterId
	}

	return filter.NewFilterId(strictFilterId.Value(), strings.Join(strictFilterNames, "/")),
		filter.NewFilterId(simpleFilterId.Value(), strings.Join(simpleFilterNames, "/"))
}

func (f *filterService) createCourseCategoryFilter(ctx context.Context, courseCategory types.CourseCategory) []filter.Id {
	var filterIds []filter.Id
	switch courseCategory {
	case types.Turf:
		filterIds = append(filterIds, filter.Turf)
	case types.Dirt:
		filterIds = append(filterIds, filter.Dirt)
	}
	return filterIds
}

func (f *filterService) createGradeClassFilter(ctx context.Context, class types.GradeClass) []filter.Id {
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

func (f *filterService) createGradeClassSimpleFilter(ctx context.Context, class types.GradeClass) []filter.Id {
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

func (f *filterService) createDistanceFilter(ctx context.Context, distance int) []filter.Id {
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

func (f *filterService) createDistanceSimpleFilter(ctx context.Context, distance int) []filter.Id {
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

func (f *filterService) createTrackConditionFilter(ctx context.Context, trackCondition types.TrackCondition) []filter.Id {
	var filterIds []filter.Id
	switch trackCondition {
	case types.GoodToFirm:
		filterIds = append(filterIds, filter.GoodTrack)
	case types.Good, types.Yielding, types.Soft:
		filterIds = append(filterIds, filter.BadTrack)
	}
	return filterIds
}

func (f *filterService) createRaceCourseFilter(ctx context.Context, raceCourseId types.RaceCourse) []filter.Id {
	var filterIds []filter.Id
	switch raceCourseId {
	case types.Tokyo, types.Nakayama, types.Hanshin, types.Kyoto:
		filterIds = append(filterIds, filter.CentralCourse)
	default:
		filterIds = append(filterIds, filter.LocalCourse)
	}
	return filterIds
}

func (f *filterService) createEntriesFilter(ctx context.Context, entries int) []filter.Id {
	var filterIds []filter.Id
	if entries <= 10 {
		filterIds = append(filterIds, filter.SmallNumberOfHorses)
	} else {
		filterIds = append(filterIds, filter.LargeNumberOfHorses)
	}
	return filterIds
}

func (f *filterService) GetAnalysisFilters() []filter.Id {
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
		filter.GoodTrackTurfShortDistance1CentralCourse,
		filter.GoodTrackTurfShortDistance2CentralCourse,
		filter.GoodTrackTurfShortDistance3CentralCourse,
		filter.GoodTrackTurfMiddleDistance1CentralCourse,
		filter.GoodTrackTurfMiddleDistance2CentralCourse,
		filter.GoodTrackTurfLongDistanceCentralCourse,
		filter.GoodTrackDirtShortDistance1CentralCourse,
		filter.GoodTrackDirtShortDistance2CentralCourse,
		filter.GoodTrackDirtShortDistance3CentralCourse,
		filter.GoodTrackDirtMiddleDistance2CentralCourse,
		filter.GoodTrackDirtLongDistanceCentralCourse,
		filter.TurfClass1,
		filter.TurfClass2,
		filter.TurfClass3,
		filter.TurfClass4,
		filter.TurfClass5,
		filter.TurfClass6,
		filter.DirtClass1,
		filter.DirtClass2,
		filter.DirtClass3,
		filter.DirtClass4,
		filter.DirtClass5,
		filter.DirtClass6,
		filter.DirtBadConditionClass1,
		filter.DirtBadConditionClass2,
		filter.DirtBadConditionClass3,
		filter.DirtBadConditionClass4,
		filter.DirtBadConditionClass5,
		filter.DirtBadConditionClass6,
	}
}
