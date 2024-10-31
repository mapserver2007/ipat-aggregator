package analysis_service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type PlaceUnHit interface {
	GetRaces(ctx context.Context, markers []*marker_csv_entity.AnalysisMarker, races []*data_cache_entity.Race) []*analysis_entity.Race
}

type placeUnHitService struct {
}

func NewPlaceUnHit() PlaceUnHit {
	return &placeUnHitService{}
}

func (p *placeUnHitService) GetRaces(
	ctx context.Context,
	markers []*marker_csv_entity.AnalysisMarker,
	races []*data_cache_entity.Race,
) []*analysis_entity.Race {
	markerMap := converter.ConvertToMap(markers, func(marker *marker_csv_entity.AnalysisMarker) types.RaceId {
		return marker.RaceId()
	})

	var unHitRaces []*analysis_entity.Race
	for _, race := range races {
		analysisMarker, ok := markerMap[race.RaceId()]
		if !ok {
			continue
		}

		var (
			analysisRaceResults []*analysis_entity.RaceResult
			analysisMarkers     []*analysis_entity.Marker
		)
		for _, raceResult := range race.RaceResults() {
			// 除外になった馬は0倍
			if raceResult.Odds().IsZero() {
				continue
			}
			// TODO race_cacheにhorseIdがないとまずい。取り直す

			if raceResult.Odds().InexactFloat64() < 1.6 && raceResult.OrderNo() > 3 {
				analysisRaceResults = append(analysisRaceResults, analysis_entity.NewRaceResult(
					raceResult.OrderNo(),
					raceResult.HorseName(),
					raceResult.HorseNumber(),
					raceResult.JockeyId(),
					raceResult.Odds(),
					raceResult.PopularNumber(),
				))

				marker := types.NoMarker
				switch raceResult.HorseNumber() {
				case analysisMarker.Favorite():
					marker = types.Favorite
				case analysisMarker.Rival():
					marker = types.Rival
				case analysisMarker.BrackTriangle():
					marker = types.BrackTriangle
				case analysisMarker.WhiteTriangle():
					marker = types.WhiteTriangle
				case analysisMarker.Star():
					marker = types.Star
				case analysisMarker.Check():
					marker = types.Check
				}

				analysisMarkers = append(analysisMarkers, analysis_entity.NewMarker(
					marker, raceResult.HorseNumber(),
				))
			}
		}

		if len(analysisRaceResults) > 0 {
			unHitRaces = append(unHitRaces, analysis_entity.NewRace(
				race.RaceId(),
				race.RaceDate(),
				race.RaceNumber(),
				race.RaceCourseId(),
				race.RaceName(),
				race.Url(),
				race.Entries(),
				race.Distance(),
				race.Class(),
				race.CourseCategory(),
				race.TrackCondition(),
				race.RaceWeightCondition(),
				analysisRaceResults,
				analysisMarkers,
			))
		}
	}

	return unHitRaces
}
