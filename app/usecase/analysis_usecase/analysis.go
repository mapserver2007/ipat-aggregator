package analysis_usecase

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/analysis_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/master_service"
)

type Analysis interface {
	Place(ctx context.Context, input *AnalysisInput) error
	PlaceAllIn(ctx context.Context, input *AnalysisInput) error
	PlaceUnHit(ctx context.Context, input *AnalysisInput) error
	PlaceJockey(ctx context.Context, input *AnalysisInput) error
	RaceTime(ctx context.Context, input *AnalysisInput) error
	Beta(ctx context.Context, input *AnalysisInput) error
}

type AnalysisInput struct {
	Markers   []*marker_csv_entity.AnalysisMarker
	Races     []*data_cache_entity.Race
	RaceTimes []*data_cache_entity.RaceTime
	Odds      *AnalysisOddsInput
	Jockeys   []*data_cache_entity.Jockey
}

type AnalysisOddsInput struct {
	Win      []*data_cache_entity.Odds
	Place    []*data_cache_entity.Odds
	Trio     []*data_cache_entity.Odds
	Quinella []*data_cache_entity.Odds
}

type analysis struct {
	placeService                analysis_service.Place
	placeAllInService           analysis_service.PlaceAllIn
	placeUnHitService           analysis_service.PlaceUnHit
	placeJockeyService          analysis_service.PlaceJockey
	betaWinService              analysis_service.BetaWin
	placeCheckPointService      analysis_service.PlaceCheckPoint
	raceTimeService             analysis_service.RaceTime
	horseMasterService          master_service.Horse
	raceForecastService         master_service.RaceForecast
	raceForecastEntityConverter converter.RaceForecastEntityConverter
	horseEntityConverter        converter.HorseEntityConverter
}

func NewAnalysis(
	placeService analysis_service.Place,
	placeAllInService analysis_service.PlaceAllIn,
	placeUnHitService analysis_service.PlaceUnHit,
	placeJockeyService analysis_service.PlaceJockey,
	betaWinService analysis_service.BetaWin,
	placeCheckPointService analysis_service.PlaceCheckPoint,
	raceTimeService analysis_service.RaceTime,
	horseMasterService master_service.Horse,
	raceForecastService master_service.RaceForecast,
	raceForecastEntityConverter converter.RaceForecastEntityConverter,
	horseEntityConverter converter.HorseEntityConverter,
) Analysis {
	return &analysis{
		placeService:                placeService,
		placeAllInService:           placeAllInService,
		placeUnHitService:           placeUnHitService,
		placeJockeyService:          placeJockeyService,
		betaWinService:              betaWinService,
		placeCheckPointService:      placeCheckPointService,
		horseMasterService:          horseMasterService,
		raceForecastService:         raceForecastService,
		raceTimeService:             raceTimeService,
		raceForecastEntityConverter: raceForecastEntityConverter,
		horseEntityConverter:        horseEntityConverter,
	}
}
