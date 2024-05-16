//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/mapserver2007/ipat-aggregator/app/controller"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/aggregation_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/master_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/summary_service"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure/gateway"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/aggregation_usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/analysis_usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/data_cache_usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/list_usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/master_usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/prediction_usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/ticket_usecase"
)

func InitializeDataCacheUseCase() *data_cache_usecase.DataCacheUseCase {
	wire.Build(
		data_cache_usecase.NewDataCacheUseCase,
		service.NewRaceConverter,
		service.NewNetKeibaService,
		service.NewRacingNumberEntityConverter,
		service.NewRaceEntityConverter,
		service.NewJockeyEntityConverter,
		service.NewTicketConverter,
		service.NewOddsEntityConverter,
		infrastructure.NewRaceDataRepository,
		infrastructure.NewRacingNumberDataRepository,
		infrastructure.NewJockeyDataRepository,
		infrastructure.NewRaceIdDataRepository,
		infrastructure.NewMarkerDataRepository,
		infrastructure.NewOddsDataRepository,
	)
	return nil
}

func InitializeMarkerAnalysisUseCase() *analysis_usecase.AnalysisUseCase {
	wire.Build(
		analysis_usecase.NewAnalysisUseCase,
		service.NewAnalysisService,
		service.NewFilterService,
		service.NewRaceConverter,
		service.NewTicketConverter,
		service.NewSpreadSheetService,
		infrastructure.NewMarkerDataRepository,
	)
	return nil
}

func InitializeListUseCase() *list_usecase.ListUseCase {
	wire.Build(
		list_usecase.NewListUseCase,
		service.NewListService,
		service.NewRaceConverter,
		service.NewTicketConverter,
		service.NewRaceEntityConverter,
	)
	return nil
}

func InitializeTicketUseCase() *ticket_usecase.TicketUseCase {
	wire.Build(
		ticket_usecase.NewTicketUseCase,
		service.NewBetNumberConverter,
		infrastructure.NewTicketCsvRepository,
	)
	return nil
}

func InitializePredictionUseCase() *prediction_usecase.PredictionUseCase {
	wire.Build(
		prediction_usecase.NewPredictionUseCase,
		service.NewNetKeibaService,
		service.NewRaceConverter,
		service.NewTicketConverter,
		service.NewRaceEntityConverter,
		service.NewFilterService,
		infrastructure.NewRaceIdDataRepository,
		infrastructure.NewPredictionDataRepository,
	)
	return nil
}

// 以下リファクタリング後

var MasterSet = wire.NewSet(
	master_usecase.NewMaster,
	master_service.NewTicket,
	master_service.NewRaceId,
	master_service.NewRace,
	master_service.NewJockey,
	master_service.NewOdds,
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
	infrastructure.NewSpreadSummeryRepository,
	gateway.NewSpreadSheetSummaryGateway,
	gateway.NewSpreadSheetTicketSummaryGateway,
	gateway.NewSpreadSheetListGateway,
	converter.NewRaceEntityConverter,
	converter.NewJockeyEntityConverter,
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
		controller.NewAggregation,
	)
	return nil
}

//var AnalysisSet = wire.NewSet(
//	analysis_usecase.NewAnalysis2,
//	//analysis_service.NewTrio,
//
//)
//
//func NewAnalysis() *controller.NewAnalysis {
//	wire.Build(
//		AnalysisSet,
//		controller.NewAnalysis,
//	)
//	return nil
//}
