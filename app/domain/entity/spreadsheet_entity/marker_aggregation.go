package spreadsheet_entity

import (
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"strconv"
)

type MarkerAggregation struct {
	markerCombinationId types.MarkerCombinationId
	raceCount           int
	hitCount            int
	hitRate             string
	payoutRate          string
	averagePayout       int
	maxPayout           int
	minPayout           int
	medianPayout        int
	filterId            filter.Id
}

func NewMarkerAggregation(
	markerCombinationId types.MarkerCombinationId,
	raceCount types.RaceCount,
	hitCount types.HitCount,
	averagePayout types.Payout,
	maxPayout types.Payout,
	minPayout types.Payout,
	medianPayout types.Payout,
	filterId filter.Id,
) *MarkerAggregation {
	hitRate := "0%"
	if raceCount > 0 {
		hitRate = fmt.Sprintf("%s%s", strconv.FormatFloat((float64(hitCount)*float64(100))/float64(raceCount), 'f', 2, 64), "%")
	}
	//payoutRate := "0%"
	//if averagePayout > 0 {
	//	payoutRate = averagePayout.Rate()
	//}
	return &MarkerAggregation{
		markerCombinationId: markerCombinationId,
		raceCount:           raceCount.Value(),
		hitCount:            hitCount.Value(),
		hitRate:             hitRate,
		averagePayout:       averagePayout.Value(),
		maxPayout:           maxPayout.Value(),
		minPayout:           minPayout.Value(),
		medianPayout:        medianPayout.Value(),
		filterId:            filterId,
	}
}
