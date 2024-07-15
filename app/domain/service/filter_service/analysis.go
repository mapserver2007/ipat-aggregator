package filter_service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type AnalysisFilter interface {
	CreatePlaceFilters(ctx context.Context, race *data_cache_entity.Race) []filter.Id
	CreatePlaceAllInFilters(ctx context.Context, race *data_cache_entity.Race, markerCombinationId types.MarkerCombinationId) []filter.Id
}

type filterService struct{}

func NewAnalysisFilter() AnalysisFilter {
	return &filterService{}
}

func (f *filterService) CreatePlaceFilters(ctx context.Context, race *data_cache_entity.Race) []filter.Id {
	var filterIds []filter.Id
	filterIds = append(filterIds, RaceCourseFilters(race.RaceCourseId())...)
	filterIds = append(filterIds, CourseCategoryFilters(race.CourseCategory())...)
	filterIds = append(filterIds, DistanceFilters(race.Distance())...)

	return filterIds
}

func (f *filterService) CreatePlaceAllInFilters(
	ctx context.Context,
	race *data_cache_entity.Race,
	markerCombinationId types.MarkerCombinationId,
) []filter.Id {
	var filterIds []filter.Id
	filterIds = append(filterIds, CourseCategoryFilters(race.CourseCategory())...)
	filterIds = append(filterIds, DistanceFilters(race.Distance())...)
	filterIds = append(filterIds, RaceCourseFilters(race.RaceCourseId())...)
	filterIds = append(filterIds, TrackConditionFilters(race.TrackCondition())...)
	filterIds = append(filterIds, MarkerFilters(markerCombinationId)...)
	return filterIds
}
