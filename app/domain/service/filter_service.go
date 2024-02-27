package service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type FilterService interface {
	CreateAnalysisFilters(ctx context.Context, race *data_cache_entity.Race, raceResultByMarker *data_cache_entity.RaceResult) []filter.Id
	CreatePredictionFilters(ctx context.Context) error
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
	switch race.CourseCategory() {
	case types.Turf:
		filterIds = append(filterIds, filter.Turf)
	case types.Dirt:
		filterIds = append(filterIds, filter.Dirt)
	}
	if race.Distance() >= 1000 && race.Distance() <= 1200 {
		filterIds = append(filterIds, filter.ShortDistance1)
	} else if race.Distance() >= 1201 && race.Distance() <= 1400 {
		filterIds = append(filterIds, filter.ShortDistance2)
	} else if race.Distance() >= 1401 && race.Distance() <= 1600 {
		filterIds = append(filterIds, filter.ShortDistance3)
	} else if race.Distance() >= 1601 && race.Distance() <= 1700 {
		filterIds = append(filterIds, filter.MiddleDistance1)
	} else if race.Distance() >= 1701 && race.Distance() <= 1800 {
		filterIds = append(filterIds, filter.MiddleDistance2)
	} else if race.Distance() >= 1801 && race.Distance() <= 2000 {
		filterIds = append(filterIds, filter.MiddleDistance3)
	} else if race.Distance() >= 2001 {
		filterIds = append(filterIds, filter.LongDistance)
	}
	switch raceResultByMarker.JockeyId() {
	case 5339, 1088, 5366, 5509, 5585: // C.ルメール, 川田将雅, R.ムーア, J.モレイラ, D.レーン
		filterIds = append(filterIds, filter.TopJockey)
	default:
		filterIds = append(filterIds, filter.OtherJockey)
	}
	switch race.Class() {
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
	case types.Maiden, types.MakeDebut:
		filterIds = append(filterIds, filter.Class1)
	}
	switch race.TrackCondition() {
	case types.GoodToFirm:
		filterIds = append(filterIds, filter.GoodTrack)
	case types.Good, types.Yielding, types.Soft:
		filterIds = append(filterIds, filter.BadTrack)
	}

	switch race.RaceCourseId() {
	case types.Tokyo, types.Nakayama, types.Hanshin, types.Kyoto:
		filterIds = append(filterIds, filter.CentralCourse)
	default:
		filterIds = append(filterIds, filter.LocalCourse)
	}

	return filterIds
}

func (f *filterService) CreatePredictionFilters(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
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
