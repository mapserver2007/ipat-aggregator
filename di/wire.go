//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/mapserver2007/ipat-aggregator/app/controller"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/aggregation_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/analysis_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/filter_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/master_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/prediction_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/summary_service"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure/gateway"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/aggregation_usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/analysis_usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/master_usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/prediction_usecase"
	"github.com/sirupsen/logrus"
)

var MasterSet = wire.NewSet(
	master_usecase.NewMaster,
	master_service.NewTicket,
	master_service.NewRaceId,
	master_service.NewRace,
	master_service.NewJockey,
	master_service.NewWinOdds,
	master_service.NewPlaceOdds,
	master_service.NewTrioOdds,
	master_service.NewAnalysisMarker,
	master_service.NewPredictionMarker,
	master_service.NewBetNumberConverter,
	master_service.NewUmacaTicket,
	master_service.NewRaceForecast,
	converter.NewRaceEntityConverter,
	converter.NewJockeyEntityConverter,
	converter.NewOddsEntityConverter,
	converter.NewRaceForecastEntityConverter,
	infrastructure.NewTicketRepository,
	infrastructure.NewRaceIdRepository,
	infrastructure.NewRaceRepository,
	infrastructure.NewRaceForecastRepository,
	infrastructure.NewJockeyRepository,
	infrastructure.NewOddsRepository,
	infrastructure.NewAnalysisMarkerRepository,
	infrastructure.NewPredictionMarkerRepository,
	infrastructure.NewUmacaTicketRepository,
	gateway.NewNetKeibaGateway,
	gateway.NewNetKeibaCollector,
	gateway.NewTospoGateway,
)

var AggregationSet = wire.NewSet(
	aggregation_usecase.NewSummary,
	aggregation_usecase.NewTicketSummary,
	aggregation_usecase.NewList,
	aggregation_service.NewSummary,
	aggregation_service.NewTicketSummary,
	aggregation_service.NewList,
	summary_service.NewTerm,
	summary_service.NewTicket,
	summary_service.NewClass,
	summary_service.NewCourseCategory,
	summary_service.NewDistanceCategory,
	summary_service.NewRaceCourse,
	infrastructure.NewSpreadSheetRepository,
	converter.NewRaceEntityConverter,
	converter.NewJockeyEntityConverter,
)

var AnalysisSet = wire.NewSet(
	analysis_usecase.NewAnalysis,
	analysis_service.NewPlace,
	analysis_service.NewTrio,
	analysis_service.NewPlaceAllIn,
	analysis_service.NewPlaceUnHit,
	analysis_service.NewPlaceCheckList,
	master_service.NewHorse,
	master_service.NewRaceForecast,
	filter_service.NewAnalysisFilter,
	infrastructure.NewHorseRepository,
	infrastructure.NewRaceForecastRepository,
	infrastructure.NewSpreadSheetRepository,
	gateway.NewNetKeibaGateway,
	gateway.NewNetKeibaCollector,
	gateway.NewTospoGateway,
	converter.NewHorseEntityConverter,
	converter.NewRaceForecastEntityConverter,
)

var PredictionSet = wire.NewSet(
	prediction_usecase.NewPrediction,
	prediction_service.NewOdds,
	prediction_service.NewPlaceCandidate,
	prediction_service.NewMarkerSync,
	filter_service.NewPredictionFilter,
	infrastructure.NewOddsRepository,
	infrastructure.NewRaceRepository,
	infrastructure.NewJockeyRepository,
	infrastructure.NewTrainerRepository,
	infrastructure.NewRaceIdRepository,
	converter.NewRaceEntityConverter,
)

var SpreadSheetGatewaySet = wire.NewSet(
	gateway.NewSpreadSheetSummaryGateway,
	gateway.NewSpreadSheetTicketSummaryGateway,
	gateway.NewSpreadSheetListGateway,
	gateway.NewSpreadSheetAnalysisPlaceGateway,
	gateway.NewSpreadSheetAnalysisPlaceAllInGateway,
	gateway.NewSpreadSheetPredictionOddsGateway,
	gateway.NewSpreadSheetPredictionCheckListGateway,
	gateway.NewSpreadSheetPredictionMarkerGateway,
)

func NewMaster(
	logger *logrus.Logger,
) *controller.Master {
	wire.Build(
		MasterSet,
		controller.NewMaster,
	)
	return nil
}

func NewAggregation(
	logger *logrus.Logger,
) *controller.Aggregation {
	wire.Build(
		AggregationSet,
		SpreadSheetGatewaySet,
		controller.NewAggregation,
	)
	return nil
}

func NewAnalysis(
	logger *logrus.Logger,
) *controller.Analysis {
	wire.Build(
		AnalysisSet,
		SpreadSheetGatewaySet,
		controller.NewAnalysis,
	)
	return nil
}

func NewPrediction(
	logger *logrus.Logger,
) *controller.Prediction {
	wire.Build(
		PredictionSet,
		AnalysisSet,
		SpreadSheetGatewaySet,
		controller.NewPrediction,
	)
	return nil
}
