package filter_service

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type AnalysisFilter interface {
	CreatePlaceFilters(ctx context.Context, race *data_cache_entity.Race) []filter.AttributeId
	CreatePlaceAllInFilters(ctx context.Context, race *data_cache_entity.Race, markerCombinationId types.MarkerCombinationId) ([]filter.AttributeId, []filter.MarkerCombinationId)
	CreateBetaFilters(ctx context.Context, race *data_cache_entity.Race, markerCombinationIds []types.MarkerCombinationId) []filter.AttributeId
}

type filterService struct{}

func NewAnalysisFilter() AnalysisFilter {
	return &filterService{}
}

func (f *filterService) CreatePlaceFilters(ctx context.Context, race *data_cache_entity.Race) []filter.AttributeId {
	var filterIds []filter.AttributeId
	filterIds = append(filterIds, RaceCourseFilters(race.RaceCourseId())...)
	filterIds = append(filterIds, CourseCategoryFilters(race.CourseCategory())...)
	filterIds = append(filterIds, DistanceFilters(race.Distance())...)

	return filterIds
}

func (f *filterService) CreatePlaceAllInFilters(
	ctx context.Context,
	race *data_cache_entity.Race,
	markerCombinationId types.MarkerCombinationId,
) ([]filter.AttributeId, []filter.MarkerCombinationId) {
	var (
		attributeFilterIds         []filter.AttributeId
		markerCombinationFilterIds []filter.MarkerCombinationId
	)
	attributeFilterIds = append(attributeFilterIds, CourseCategoryFilters(race.CourseCategory())...)
	attributeFilterIds = append(attributeFilterIds, DistanceFilters(race.Distance())...)
	attributeFilterIds = append(attributeFilterIds, RaceCourseFilters(race.RaceCourseId())...)
	attributeFilterIds = append(attributeFilterIds, TrackConditionFilters(race.TrackCondition())...)
	attributeFilterIds = append(attributeFilterIds, GradeClassFilters(race.Class())...)
	attributeFilterIds = append(attributeFilterIds, SeasonFilters(race.RaceDate())...)

	markerCombinationFilterIds = append(markerCombinationFilterIds, MarkerCombinationFilter(race, markerCombinationId)...)

	return attributeFilterIds, markerCombinationFilterIds
}

func (f *filterService) CreateBetaFilters(
	ctx context.Context,
	race *data_cache_entity.Race,
	markerCombinationIds []types.MarkerCombinationId,
) []filter.AttributeId {
	var filterIds []filter.AttributeId
	filterIds = append(filterIds, CourseCategoryFilters(race.CourseCategory())...)
	filterIds = append(filterIds, DistanceFilters(race.Distance())...)
	filterIds = append(filterIds, RaceCourseFilters(race.RaceCourseId())...)
	filterIds = append(filterIds, TrackConditionFilters(race.TrackCondition())...)
	return filterIds
}
