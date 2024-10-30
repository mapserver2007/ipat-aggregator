package analysis_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
)

type PlaceUnHit interface {
	Create(ctx context.Context, markers []*marker_csv_entity.AnalysisMarker, races []*data_cache_entity.Race) error
}

type placeUnHitService struct {
}

func NewPlaceUnHit() PlaceUnHit {
	return &placeUnHitService{}
}

func (p *placeUnHitService) Create(
	ctx context.Context,
	markers []*marker_csv_entity.AnalysisMarker,
	races []*data_cache_entity.Race,
) error {
	//markerMap := converter.ConvertToMap(markers, func(marker *marker_csv_entity.AnalysisMarker) types.RaceId {
	//	return marker.RaceId()
	//})

	var unHitRaces []*data_cache_entity.Race

	for _, race := range races {
		//analysisMarker, ok := markerMap[race.RaceId()]
		//if !ok {
		//	continue
		//}
		//analysisMarker.

		for _, raceResult := range race.RaceResults() {
			// 除外になった馬は0倍
			if raceResult.Odds().IsZero() {
				continue
			}
			if raceResult.Odds().InexactFloat64() < 1.6 && raceResult.OrderNo() > 3 {
				unHitRaces = append(unHitRaces, race)
				break
			}
		}
	}

	fmt.Println(unHitRaces)

	return nil
}
