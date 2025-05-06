package filter_service

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type PredictionFilter interface {
	CreateRaceConditionFilters(ctx context.Context, race *netkeiba_entity.Race) []filter.AttributeId
	CreateRaceTimeConditionFilters(ctx context.Context, race *netkeiba_entity.Race) []filter.AttributeId
}

type predictionFilter struct{}

func NewPredictionFilter() PredictionFilter {
	return &predictionFilter{}
}

func (p *predictionFilter) CreateRaceConditionFilters(
	ctx context.Context,
	race *netkeiba_entity.Race,
) []filter.AttributeId {
	var filterIds []filter.AttributeId
	filterIds = append(filterIds, CourseCategoryFilters(types.CourseCategory(race.CourseCategory()))...)
	filterIds = append(filterIds, DistanceFilters(race.Distance())...)
	filterIds = append(filterIds, RaceCourseFilters(types.RaceCourse(race.RaceCourseId()))...)

	return filterIds
}

func (p *predictionFilter) CreateRaceTimeConditionFilters(
	ctx context.Context,
	race *netkeiba_entity.Race,
) []filter.AttributeId {
	var filterIds []filter.AttributeId
	filterIds = append(filterIds, RaceCourseFilters(types.RaceCourse(race.RaceCourseId()))...)
	filterIds = append(filterIds, CourseCategoryFilters(types.CourseCategory(race.CourseCategory()))...)
	filterIds = append(filterIds, DistanceFilters(race.Distance())...)
	filterIds = append(filterIds, GradeClassFilters(types.GradeClass(race.Class()))...)
	filterIds = append(filterIds, TrackConditionFilters(types.TrackCondition(race.TrackCondition()))...)
	filterIds = append(filterIds, RaceAgeConditionFilters(types.RaceAgeCondition(race.RaceAgeCondition()))...)

	return filterIds
}
