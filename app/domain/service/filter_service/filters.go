package filter_service

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

func CourseCategoryFilters(courseCategory types.CourseCategory) []filter.AttributeId {
	var filterIds []filter.AttributeId
	switch courseCategory {
	case types.Turf:
		filterIds = append(filterIds, filter.Turf)
	case types.Dirt:
		filterIds = append(filterIds, filter.Dirt)
	}
	return filterIds
}

func RaceCourseFilters(raceCourseId types.RaceCourse) []filter.AttributeId {
	var filterIds []filter.AttributeId
	switch raceCourseId {
	case types.Tokyo:
		filterIds = append(filterIds, filter.Tokyo)
	case types.Nakayama:
		filterIds = append(filterIds, filter.Nakayama)
	case types.Kyoto:
		filterIds = append(filterIds, filter.Kyoto)
	case types.Hanshin:
		filterIds = append(filterIds, filter.Hanshin)
	case types.Niigata:
		filterIds = append(filterIds, filter.Niigata)
	case types.Chukyo:
		filterIds = append(filterIds, filter.Chukyo)
	case types.Sapporo:
		filterIds = append(filterIds, filter.Sapporo)
	case types.Hakodate:
		filterIds = append(filterIds, filter.Hakodate)
	case types.Fukushima:
		filterIds = append(filterIds, filter.Fukushima)
	case types.Kokura:
		filterIds = append(filterIds, filter.Kokura)
	}
	return filterIds
}

func DistanceFilters(distance int) []filter.AttributeId {
	var filterIds []filter.AttributeId
	switch distance {
	case 1000:
		filterIds = append(filterIds, filter.Distance1000m)
	case 1150:
		filterIds = append(filterIds, filter.Distance1150m)
	case 1200:
		filterIds = append(filterIds, filter.Distance1200m)
	case 1300:
		filterIds = append(filterIds, filter.Distance1300m)
	case 1400:
		filterIds = append(filterIds, filter.Distance1400m)
	case 1500:
		filterIds = append(filterIds, filter.Distance1500m)
	case 1600:
		filterIds = append(filterIds, filter.Distance1600m)
	case 1700:
		filterIds = append(filterIds, filter.Distance1700m)
	case 1800:
		filterIds = append(filterIds, filter.Distance1800m)
	case 1900:
		filterIds = append(filterIds, filter.Distance1900m)
	case 2000:
		filterIds = append(filterIds, filter.Distance2000m)
	case 2100:
		filterIds = append(filterIds, filter.Distance2100m)
	case 2200:
		filterIds = append(filterIds, filter.Distance2200m)
	case 2300:
		filterIds = append(filterIds, filter.Distance2300m)
	case 2400:
		filterIds = append(filterIds, filter.Distance2400m)
	case 2500:
		filterIds = append(filterIds, filter.Distance2500m)
	case 2600:
		filterIds = append(filterIds, filter.Distance2600m)
	case 3000:
		filterIds = append(filterIds, filter.Distance3000m)
	case 3200:
		filterIds = append(filterIds, filter.Distance3200m)
	case 3400:
		filterIds = append(filterIds, filter.Distance3400m)
	case 3600:
		filterIds = append(filterIds, filter.Distance3600m)
	}
	return filterIds
}

func TrackConditionFilters(trackCondition types.TrackCondition) []filter.AttributeId {
	var filterIds []filter.AttributeId
	switch trackCondition {
	case types.GoodToFirm:
		filterIds = append(filterIds, filter.GoodToFirm)
	case types.Good:
		filterIds = append(filterIds, filter.Good)
	case types.Yielding:
		filterIds = append(filterIds, filter.Yielding)
	case types.Soft:
		filterIds = append(filterIds, filter.Soft)
	}
	return filterIds
}

func MarkerCombinationFilter(
	race *data_cache_entity.Race,
	markerCombinationId types.MarkerCombinationId,
) []filter.MarkerCombinationId {
	var filterIds []filter.MarkerCombinationId

	switch race.CourseCategory() {
	case types.Turf:
		filterIds = append(filterIds, filter.MarkerCombinationTurf)
	case types.Dirt:
		filterIds = append(filterIds, filter.MarkerCombinationDirt)
	}

	for _, result := range race.PayoutResults() {
		switch result.TicketType() {
		case types.Win:
			filterIds = append(filterIds, filter.MarkerCombinationWin)
			marker, _ := types.NewMarker(markerCombinationId.Value() % 10)
			switch marker {
			case types.Favorite:
				filterIds = append(filterIds, filter.MarkerCombinationFavorite)
			case types.Rival:
				filterIds = append(filterIds, filter.MarkerCombinationRival)
			case types.BrackTriangle:
				filterIds = append(filterIds, filter.MarkerCombinationBrackTriangle)
			case types.WhiteTriangle:
				filterIds = append(filterIds, filter.MarkerCombinationWhiteTriangle)
			case types.Star:
				filterIds = append(filterIds, filter.MarkerCombinationStar)
			case types.Check:
				filterIds = append(filterIds, filter.MarkerCombinationCheck)
			}
		case types.Place:
			filterIds = append(filterIds, filter.MarkerCombinationPlace)
			marker, _ := types.NewMarker(markerCombinationId.Value() % 10)
			switch marker {
			case types.Favorite:
				filterIds = append(filterIds, filter.MarkerCombinationFavorite)
			case types.Rival:
				filterIds = append(filterIds, filter.MarkerCombinationRival)
			case types.BrackTriangle:
				filterIds = append(filterIds, filter.MarkerCombinationBrackTriangle)
			case types.WhiteTriangle:
				filterIds = append(filterIds, filter.MarkerCombinationWhiteTriangle)
			case types.Star:
				filterIds = append(filterIds, filter.MarkerCombinationStar)
			case types.Check:
				filterIds = append(filterIds, filter.MarkerCombinationCheck)
			}
		}
	}

	return filterIds
}
