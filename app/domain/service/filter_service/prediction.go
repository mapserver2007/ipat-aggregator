package filter_service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type PredictionFilter interface {
	Create(ctx context.Context, race *netkeiba_entity.Race) filter.Id
}

type predictionFilter struct{}

func NewPredictionFilter() PredictionFilter {
	return &predictionFilter{}
}

func (p *predictionFilter) Create(
	ctx context.Context,
	race *netkeiba_entity.Race,
) filter.Id {
	var filterIds []filter.Id
	filterIds = append(filterIds, CourseCategoryFilters(types.CourseCategory(race.CourseCategory()))...)
	filterIds = append(filterIds, DistanceFilters(race.Distance())...)
	filterIds = append(filterIds, RaceCourseFilters(types.RaceCourse(race.RaceCourseId()))...)

	var predictionFilterId filter.Id
	for _, f := range filterIds {
		predictionFilterId |= f
	}

	return predictionFilterId
}
