// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from wire.go:

func InitializeDataCacheUseCase() *data_cache_usecase.DataCacheUseCase {
	racingNumberDataRepository := infrastructure.NewRacingNumberDataRepository()
	raceDataRepository := infrastructure.NewRaceDataRepository()
	jockeyDataRepository := infrastructure.NewJockeyDataRepository()
	raceIdDataRepository := infrastructure.NewRaceIdDataRepository()
	markerDataRepository := infrastructure.NewMarkerDataRepository()
	oddsDataRepository := infrastructure.NewOddsDataRepository()
	raceConverter := service.NewRaceConverter()
	ticketConverter := service.NewTicketConverter(raceConverter)
	netKeibaService := service.NewNetKeibaService(raceConverter, ticketConverter)
	racingNumberEntityConverter := service.NewRacingNumberEntityConverter()
	raceEntityConverter := service.NewRaceEntityConverter()
	jockeyEntityConverter := service.NewJockeyEntityConverter()
	oddsEntityConverter := service.NewOddsEntityConverter()
	dataCacheUseCase := data_cache_usecase.NewDataCacheUseCase(racingNumberDataRepository, raceDataRepository, jockeyDataRepository, raceIdDataRepository, markerDataRepository, oddsDataRepository, netKeibaService, raceConverter, racingNumberEntityConverter, raceEntityConverter, jockeyEntityConverter, oddsEntityConverter)
	return dataCacheUseCase
}

func InitializeMarkerAnalysisUseCase() *analysis_usecase.AnalysisUseCase {
	markerDataRepository := infrastructure.NewMarkerDataRepository()
	filterService := service.NewFilterService()
	spreadSheetService := service.NewSpreadSheetService()
	analysisService := service.NewAnalysisService(filterService, spreadSheetService)
	raceConverter := service.NewRaceConverter()
	ticketConverter := service.NewTicketConverter(raceConverter)
	analysisUseCase := analysis_usecase.NewAnalysisUseCase(markerDataRepository, analysisService, ticketConverter)
	return analysisUseCase
}

func InitializeListUseCase() *list_usecase.ListUseCase {
	raceConverter := service.NewRaceConverter()
	ticketConverter := service.NewTicketConverter(raceConverter)
	raceEntityConverter := service.NewRaceEntityConverter()
	listService := service.NewListService(raceConverter, ticketConverter, raceEntityConverter)
	listUseCase := list_usecase.NewListUseCase(listService)
	return listUseCase
}

func InitializeTicketUseCase() *ticket_usecase.TicketUseCase {
	betNumberConverter := service.NewBetNumberConverter()
	ticketCsvRepository := infrastructure.NewTicketCsvRepository(betNumberConverter)
	ticketUseCase := ticket_usecase.NewTicketUseCase(ticketCsvRepository)
	return ticketUseCase
}

func InitializePredictionUseCase() *prediction_usecase.PredictionUseCase {
	raceConverter := service.NewRaceConverter()
	ticketConverter := service.NewTicketConverter(raceConverter)
	netKeibaService := service.NewNetKeibaService(raceConverter, ticketConverter)
	raceIdDataRepository := infrastructure.NewRaceIdDataRepository()
	predictionDataRepository := infrastructure.NewPredictionDataRepository()
	raceEntityConverter := service.NewRaceEntityConverter()
	filterService := service.NewFilterService()
	predictionUseCase := prediction_usecase.NewPredictionUseCase(netKeibaService, raceIdDataRepository, predictionDataRepository, raceEntityConverter, filterService)
	return predictionUseCase
}

func NewMaster() *controller.Master {
	betNumberConverter := master_service.NewBetNumberConverter()
	ticketRepository := infrastructure.NewTicketRepository(betNumberConverter)
	ticket := master_service.NewTicket(ticketRepository)
	netKeibaGateway := gateway.NewNetKeibaGateway()
	raceIdRepository := infrastructure.NewRaceIdRepository(netKeibaGateway)
	raceId := master_service.NewRaceId(raceIdRepository)
	raceRepository := infrastructure.NewRaceRepository(netKeibaGateway)
	raceEntityConverter := converter.NewRaceEntityConverter()
	race := master_service.NewRace(raceRepository, raceEntityConverter)
	jockeyRepository := infrastructure.NewJockeyRepository(netKeibaGateway)
	jockeyEntityConverter := converter.NewJockeyEntityConverter()
	jockey := master_service.NewJockey(jockeyRepository, jockeyEntityConverter)
	oddsRepository := infrastructure.NewOddsRepository(netKeibaGateway)
	oddsEntityConverter := converter.NewOddsEntityConverter()
	odds := master_service.NewOdds(oddsRepository, oddsEntityConverter)
	analysisMarkerRepository := infrastructure.NewAnalysisMarkerRepository()
	analysisMarker := master_service.NewAnalysisMarker(analysisMarkerRepository)
	predictionMarkerRepository := infrastructure.NewPredictionMarkerRepository()
	predictionMarker := master_service.NewPredictionMarker(predictionMarkerRepository)
	master := master_usecase.NewMaster(ticket, raceId, race, jockey, odds, analysisMarker, predictionMarker)
	controllerMaster := controller.NewMaster(master)
	return controllerMaster
}

func NewAggregation() *controller.Aggregation {
	term := summary_service.NewTerm()
	ticket := summary_service.NewTicket()
	class := summary_service.NewClass()
	courseCategory := summary_service.NewCourseCategory()
	distanceCategory := summary_service.NewDistanceCategory()
	raceCourse := summary_service.NewRaceCourse()
	spreadSheetSummaryGateway := gateway.NewSpreadSheetSummaryGateway()
	spreadSheetTicketSummaryGateway := gateway.NewSpreadSheetTicketSummaryGateway()
	spreadSheetListGateway := gateway.NewSpreadSheetListGateway()
	spreadSheetRepository := infrastructure.NewSpreadSummeryRepository(spreadSheetSummaryGateway, spreadSheetTicketSummaryGateway, spreadSheetListGateway)
	summary := aggregation_service.NewSummary(term, ticket, class, courseCategory, distanceCategory, raceCourse, spreadSheetRepository)
	aggregation_usecaseSummary := aggregation_usecase.NewSummary(summary)
	ticketSummary := aggregation_service.NewTicketSummary(term, spreadSheetRepository)
	aggregation_usecaseTicketSummary := aggregation_usecase.NewTicketSummary(ticketSummary)
	raceEntityConverter := converter.NewRaceEntityConverter()
	jockeyEntityConverter := converter.NewJockeyEntityConverter()
	list := aggregation_service.NewList(raceEntityConverter, jockeyEntityConverter, spreadSheetRepository)
	aggregation_usecaseList := aggregation_usecase.NewList(list)
	aggregation := controller.NewAggregation(aggregation_usecaseSummary, aggregation_usecaseTicketSummary, aggregation_usecaseList)
	return aggregation
}

// wire.go:

var MasterSet = wire.NewSet(master_usecase.NewMaster, master_service.NewTicket, master_service.NewRaceId, master_service.NewRace, master_service.NewJockey, master_service.NewOdds, master_service.NewAnalysisMarker, master_service.NewPredictionMarker, master_service.NewBetNumberConverter, converter.NewRaceEntityConverter, converter.NewJockeyEntityConverter, converter.NewOddsEntityConverter, infrastructure.NewTicketRepository, infrastructure.NewRaceIdRepository, infrastructure.NewRaceRepository, infrastructure.NewJockeyRepository, infrastructure.NewOddsRepository, infrastructure.NewAnalysisMarkerRepository, infrastructure.NewPredictionMarkerRepository, gateway.NewNetKeibaGateway)

var AggregationSet = wire.NewSet(aggregation_usecase.NewSummary, aggregation_usecase.NewTicketSummary, aggregation_usecase.NewList, aggregation_service.NewSummary, aggregation_service.NewTicketSummary, aggregation_service.NewList, summary_service.NewTerm, summary_service.NewTicket, summary_service.NewClass, summary_service.NewCourseCategory, summary_service.NewDistanceCategory, summary_service.NewRaceCourse, infrastructure.NewSpreadSummeryRepository, gateway.NewSpreadSheetSummaryGateway, gateway.NewSpreadSheetTicketSummaryGateway, gateway.NewSpreadSheetListGateway, converter.NewRaceEntityConverter, converter.NewJockeyEntityConverter)
