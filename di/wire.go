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
)

var MasterSet = wire.NewSet(
	master_usecase.NewMaster,
	master_service.NewTicket,
	master_service.NewRaceId,
	master_service.NewRace,
	master_service.NewJockey,
	master_service.NewTrioOdds,
	master_service.NewAnalysisMarker,
	master_service.NewPredictionMarker,
	master_service.NewBetNumberConverter,
	converter.NewRaceEntityConverter,
	converter.NewJockeyEntityConverter,
	converter.NewOddsEntityConverter,
	infrastructure.NewTicketRepository,
	infrastructure.NewRaceIdRepository,
	infrastructure.NewRaceRepository,
	infrastructure.NewJockeyRepository,
	infrastructure.NewOddsRepository,
	infrastructure.NewAnalysisMarkerRepository,
	infrastructure.NewPredictionMarkerRepository,
	gateway.NewNetKeibaGateway,
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
	filter_service.NewAnalysisFilter,
	infrastructure.NewSpreadSheetRepository,
)

var PredictionSet = wire.NewSet(
	prediction_usecase.NewPrediction,
	prediction_service.NewOdds,
	filter_service.NewPredictionFilter,
	infrastructure.NewOddsRepository,
	infrastructure.NewRaceRepository,
	gateway.NewNetKeibaGateway,
)

var SpreadSheetGatewaySet = wire.NewSet(
	gateway.NewSpreadSheetSummaryGateway,
	gateway.NewSpreadSheetTicketSummaryGateway,
	gateway.NewSpreadSheetListGateway,
	gateway.NewSpreadSheetAnalysisPlaceGateway,
	gateway.NewSpreadSheetPredictionGateway,
)

func NewMaster() *controller.Master {
	wire.Build(
		MasterSet,
		controller.NewMaster,
	)
	return nil
}

func NewAggregation() *controller.Aggregation {
	wire.Build(
		AggregationSet,
		SpreadSheetGatewaySet,
		controller.NewAggregation,
	)
	return nil
}

func NewAnalysis() *controller.Analysis {
	wire.Build(
		AnalysisSet,
		SpreadSheetGatewaySet,
		controller.NewAnalysis,
	)
	return nil
}

func NewPrediction() *controller.Prediction {
	wire.Build(
		PredictionSet,
		AnalysisSet,
		SpreadSheetGatewaySet,
		controller.NewPrediction,
	)
	return nil
}
