package analysis_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/filter_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"github.com/mapserver2007/ipat-aggregator/config"
)

type PlaceUnHit interface {
	GetUnHitRaces(ctx context.Context, markers []*marker_csv_entity.AnalysisMarker, races []*data_cache_entity.Race) []*analysis_entity.Race
	GetUnHitRaceRate(ctx context.Context, race *analysis_entity.Race, calculables []*analysis_entity.PlaceCalculable) error
	FetchHorse(ctx context.Context, horseId types.HorseId) (*netkeiba_entity.Horse, error)
}

type placeUnHitService struct {
	horseRepository      repository.HorseRepository
	horseEntityConverter converter.HorseEntityConverter
	filterService        filter_service.AnalysisFilter
}

func NewPlaceUnHit(
	horseRepository repository.HorseRepository,
	horseEntityConverter converter.HorseEntityConverter,
	filterService filter_service.AnalysisFilter,
) PlaceUnHit {
	return &placeUnHitService{
		horseRepository:      horseRepository,
		horseEntityConverter: horseEntityConverter,
		filterService:        filterService,
	}
}

func (p *placeUnHitService) GetUnHitRaces(
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

			if raceResult.Odds().InexactFloat64() < config.AnalysisUnHitWinLowerOdds && raceResult.OrderNo() > 3 {
				analysisRaceResults = append(analysisRaceResults, analysis_entity.NewRaceResult(
					raceResult.OrderNo(),
					raceResult.HorseId(),
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
				p.filterService.CreatePlaceFilters(ctx, race),
			))
		}
	}

	return unHitRaces
}

func (p *placeUnHitService) GetUnHitRaceRate(
	ctx context.Context,
	race *analysis_entity.Race,
	calculables []*analysis_entity.PlaceCalculable,
) error {

	var analysisFilter filter.Id
	for _, f := range race.AnalysisFilters() {
		analysisFilter |= f
	}

	// 実用上は1レースで分析対象のオッズは1つ想定だが、仕様上は複数オッズも計算可能なのでループを回す
	for _, raceResult := range race.RaceResults() {
		for _, calculable := range calculables {
			match := true
			for _, f := range calculable.Filters() {
				if f&analysisFilter == 0 {
					match = false
					break
				}
			}
			if match {
				odds := raceResult.Odds().InexactFloat64()
				fmt.Println(odds)
			}
		}

	}

	return nil
}

func (p *placeUnHitService) FetchHorse(
	ctx context.Context,
	horseId types.HorseId,
) (*netkeiba_entity.Horse, error) {
	horse, err := p.horseRepository.Fetch(ctx, fmt.Sprintf(horseUrl, horseId))
	if err != nil {
		return nil, err
	}

	return horse, nil
}
