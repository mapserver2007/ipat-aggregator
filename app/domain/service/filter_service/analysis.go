package filter_service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
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
		filter.All2,
		filter.NiigataTurf1000m,
		filter.HakodateTurf1000m,
		filter.NakayamaTurf1200m,
		filter.KyotoTurf1200m,
		filter.HanshinTurf1200m,
		filter.NiigataTurf1200m,
		filter.ChukyoTurf1200m,
		filter.SapporoTurf1200m,
		filter.HakodateTurf1200m,
		filter.FukushimaTurf1200m,
		filter.KokuraTurf1200m,
		filter.TokyoTurf1400m,
		filter.KyotoTurf1400m,
		filter.HanshinTurf1400m,
		filter.NiigataTurf1400m,
		filter.ChukyoTurf1400m,
		filter.SapporoTurf1500m,
		filter.NakayamaTurf1600m,
		filter.TokyoTurf1600m,
		filter.KyotoTurf1600m,
		filter.HanshinTurf1600m,
		filter.ChukyoTurf1600m,
		filter.NakayamaTurf1800m,
		filter.TokyoTurf1800m,
		filter.KyotoTurf1800m,
		filter.HanshinTurf1800m,
		filter.NiigataTurf1800m,
		filter.SapporoTurf1800m,
		filter.HakodateTurf1800m,
		filter.FukushimaTurf1800m,
		filter.KokuraTurf1800m,
		filter.NakayamaTurf2000m,
		filter.TokyoTurf2000m,
		filter.KyotoTurf2000m,
		filter.NiigataTurf2000m,
		filter.ChukyoTurf2000m,
		filter.SapporoTurf2000m,
		filter.HakodateTurf2000m,
		filter.FukushimaTurf2000m,
		filter.KokuraTurf2000m,
		filter.NakayamaTurf2200m,
		filter.KyotoTurf2200m,
		filter.HanshinTurf2200m,
		filter.NiigataTurf2200m,
		filter.ChukyoTurf2200m,
		filter.TokyoTurf2300m,
		filter.TokyoTurf2400m,
		filter.KyotoTurf2400m,
		filter.HanshinTurf2400m,
		filter.NiigataTurf2400m,
		filter.NakayamaTurf2500m,
		filter.TokyoTurf2500m,
		filter.HanshinTurf2600m,
		filter.SapporoTurf2600m,
		filter.HakodateTurf2600m,
		filter.FukushimaTurf2600m,
		filter.KokuraTurf2600m,
		filter.HanshinTurf3000m,
		filter.ChukyoTurf3000m,
		filter.KyotoTurf3200m,
		filter.TokyoTurf3400m,
		filter.NakayamaTurf3600m,
		filter.SapporoDirt1000m,
		filter.HakodateDirt1000m,
		filter.KokuraDirt1000m,
		filter.FukushimaDirt1150m,
		filter.NakayamaDirt1200m,
		filter.KyotoDirt1200m,
		filter.NiigataDirt1200m,
		filter.ChukyoDirt1200m,
		filter.TokyoDirt1300m,
		filter.TokyoDirt1400m,
		filter.KyotoDirt1400m,
		filter.HanshinDirt1400m,
		filter.ChukyoDirt1400m,
		filter.TokyoDirt1600m,
		filter.SapporoDirt1700m,
		filter.HakodateDirt1700m,
		filter.FukushimaDirt1700m,
		filter.KokuraDirt1700m,
		filter.NakayamaDirt1800m,
		filter.KyotoDirt1800m,
		filter.HanshinDirt1800m,
		filter.NiigataDirt1800m,
		filter.ChukyoDirt1800m,
		filter.KyotoDirt1900m,
		filter.ChukyoDirt1900m,
		filter.HanshinDirt2000m,
		filter.TokyoDirt2100m,
		filter.NakayamaDirt2400m,
		filter.SapporoDirt2400m,
		filter.HakodateDirt2400m,
		filter.FukushimaDirt2400m,
		filter.KokuraDirt2400m,
		filter.NakayamaDirt2500m,
		filter.NiigataDirt2500m,
	}
}

func (f *filterService) Create(ctx context.Context, race *data_cache_entity.Race) []filter.Id {
	var filterIds []filter.Id
	filterIds = append(filterIds, CourseCategoryFilters(race.CourseCategory())...)
	filterIds = append(filterIds, DistanceFilters(race.Distance())...)
	filterIds = append(filterIds, RaceCourseFilters(race.RaceCourseId())...)

	return filterIds
}
