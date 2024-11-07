// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from wire.go:

func NewMaster() *controller.Master {
	betNumberConverter := master_service.NewBetNumberConverter()
	ticketRepository := infrastructure.NewTicketRepository(betNumberConverter)
	ticket := master_service.NewTicket(ticketRepository)
	netKeibaCollector := gateway.NewNetKeibaCollector()
	netKeibaGateway := gateway.NewNetKeibaGateway(netKeibaCollector)
	raceIdRepository := infrastructure.NewRaceIdRepository(netKeibaGateway)
	raceId := master_service.NewRaceId(raceIdRepository)
	raceRepository := infrastructure.NewRaceRepository(netKeibaGateway)
	raceEntityConverter := converter.NewRaceEntityConverter()
	race := master_service.NewRace(raceRepository, raceEntityConverter)
	tospoGateway := gateway.NewTospoGateway()
	raceForecastRepository := infrastructure.NewRaceForecastRepository(tospoGateway)
	raceForecastEntityConverter := converter.NewRaceForecastEntityConverter()
	raceForecast := master_service.NewRaceForecast(raceForecastRepository, raceForecastEntityConverter)
	jockeyRepository := infrastructure.NewJockeyRepository(netKeibaGateway)
	jockeyEntityConverter := converter.NewJockeyEntityConverter()
	jockey := master_service.NewJockey(jockeyRepository, jockeyEntityConverter)
	oddsRepository := infrastructure.NewOddsRepository(netKeibaGateway)
	oddsEntityConverter := converter.NewOddsEntityConverter()
	winOdds := master_service.NewWinOdds(oddsRepository, oddsEntityConverter)
	placeOdds := master_service.NewPlaceOdds(oddsRepository, oddsEntityConverter)
	trioOdds := master_service.NewTrioOdds(oddsRepository, oddsEntityConverter)
	analysisMarkerRepository := infrastructure.NewAnalysisMarkerRepository()
	analysisMarker := master_service.NewAnalysisMarker(analysisMarkerRepository)
	predictionMarkerRepository := infrastructure.NewPredictionMarkerRepository(netKeibaGateway)
	predictionMarker := master_service.NewPredictionMarker(predictionMarkerRepository)
	umacaTicketRepository := infrastructure.NewUmacaTicketRepository()
	umacaTicket := master_service.NewUmacaTicket(umacaTicketRepository, ticketRepository)
	master := master_usecase.NewMaster(ticket, raceId, race, raceForecast, jockey, winOdds, placeOdds, trioOdds, analysisMarker, predictionMarker, umacaTicket)
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
	spreadSheetAnalysisPlaceGateway := gateway.NewSpreadSheetAnalysisPlaceGateway()
	spreadSheetAnalysisPlaceAllInGateway := gateway.NewSpreadSheetAnalysisPlaceAllInGateway()
	spreadSheetPredictionOddsGateway := gateway.NewSpreadSheetPredictionOddsGateway()
	spreadSheetPredictionCheckListGateway := gateway.NewSpreadSheetPredictionCheckListGateway()
	spreadSheetPredictionMarkerGateway := gateway.NewSpreadSheetPredictionMarkerGateway()
	spreadSheetRepository := infrastructure.NewSpreadSheetRepository(spreadSheetSummaryGateway, spreadSheetTicketSummaryGateway, spreadSheetListGateway, spreadSheetAnalysisPlaceGateway, spreadSheetAnalysisPlaceAllInGateway, spreadSheetPredictionOddsGateway, spreadSheetPredictionCheckListGateway, spreadSheetPredictionMarkerGateway)
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

func NewAnalysis() *controller.Analysis {
	analysisFilter := filter_service.NewAnalysisFilter()
	spreadSheetSummaryGateway := gateway.NewSpreadSheetSummaryGateway()
	spreadSheetTicketSummaryGateway := gateway.NewSpreadSheetTicketSummaryGateway()
	spreadSheetListGateway := gateway.NewSpreadSheetListGateway()
	spreadSheetAnalysisPlaceGateway := gateway.NewSpreadSheetAnalysisPlaceGateway()
	spreadSheetAnalysisPlaceAllInGateway := gateway.NewSpreadSheetAnalysisPlaceAllInGateway()
	spreadSheetPredictionOddsGateway := gateway.NewSpreadSheetPredictionOddsGateway()
	spreadSheetPredictionCheckListGateway := gateway.NewSpreadSheetPredictionCheckListGateway()
	spreadSheetPredictionMarkerGateway := gateway.NewSpreadSheetPredictionMarkerGateway()
	spreadSheetRepository := infrastructure.NewSpreadSheetRepository(spreadSheetSummaryGateway, spreadSheetTicketSummaryGateway, spreadSheetListGateway, spreadSheetAnalysisPlaceGateway, spreadSheetAnalysisPlaceAllInGateway, spreadSheetPredictionOddsGateway, spreadSheetPredictionCheckListGateway, spreadSheetPredictionMarkerGateway)
	place := analysis_service.NewPlace(analysisFilter, spreadSheetRepository)
	trio := analysis_service.NewTrio(analysisFilter)
	placeAllIn := analysis_service.NewPlaceAllIn(analysisFilter, spreadSheetRepository)
	netKeibaCollector := gateway.NewNetKeibaCollector()
	netKeibaGateway := gateway.NewNetKeibaGateway(netKeibaCollector)
	horseRepository := infrastructure.NewHorseRepository(netKeibaGateway)
	horseEntityConverter := converter.NewHorseEntityConverter()
	placeUnHit := analysis_service.NewPlaceUnHit(horseRepository, horseEntityConverter, analysisFilter)
	horse := master_service.NewHorse(horseRepository, horseEntityConverter)
	analysis := analysis_usecase.NewAnalysis(place, trio, placeAllIn, placeUnHit, horse, horseEntityConverter)
	controllerAnalysis := controller.NewAnalysis(analysis)
	return controllerAnalysis
}

func NewPrediction() *controller.Prediction {
	netKeibaCollector := gateway.NewNetKeibaCollector()
	netKeibaGateway := gateway.NewNetKeibaGateway(netKeibaCollector)
	oddsRepository := infrastructure.NewOddsRepository(netKeibaGateway)
	raceRepository := infrastructure.NewRaceRepository(netKeibaGateway)
	spreadSheetSummaryGateway := gateway.NewSpreadSheetSummaryGateway()
	spreadSheetTicketSummaryGateway := gateway.NewSpreadSheetTicketSummaryGateway()
	spreadSheetListGateway := gateway.NewSpreadSheetListGateway()
	spreadSheetAnalysisPlaceGateway := gateway.NewSpreadSheetAnalysisPlaceGateway()
	spreadSheetAnalysisPlaceAllInGateway := gateway.NewSpreadSheetAnalysisPlaceAllInGateway()
	spreadSheetPredictionOddsGateway := gateway.NewSpreadSheetPredictionOddsGateway()
	spreadSheetPredictionCheckListGateway := gateway.NewSpreadSheetPredictionCheckListGateway()
	spreadSheetPredictionMarkerGateway := gateway.NewSpreadSheetPredictionMarkerGateway()
	spreadSheetRepository := infrastructure.NewSpreadSheetRepository(spreadSheetSummaryGateway, spreadSheetTicketSummaryGateway, spreadSheetListGateway, spreadSheetAnalysisPlaceGateway, spreadSheetAnalysisPlaceAllInGateway, spreadSheetPredictionOddsGateway, spreadSheetPredictionCheckListGateway, spreadSheetPredictionMarkerGateway)
	predictionFilter := filter_service.NewPredictionFilter()
	odds := prediction_service.NewOdds(oddsRepository, raceRepository, spreadSheetRepository, predictionFilter)
	tospoGateway := gateway.NewTospoGateway()
	raceForecastRepository := infrastructure.NewRaceForecastRepository(tospoGateway)
	horseRepository := infrastructure.NewHorseRepository(netKeibaGateway)
	jockeyRepository := infrastructure.NewJockeyRepository(netKeibaGateway)
	trainerRepository := infrastructure.NewTrainerRepository(netKeibaGateway)
	raceEntityConverter := converter.NewRaceEntityConverter()
	horseEntityConverter := converter.NewHorseEntityConverter()
	placeCheckList := prediction_service.NewPlaceCheckList()
	placeCandidate := prediction_service.NewPlaceCandidate(raceRepository, raceForecastRepository, horseRepository, jockeyRepository, trainerRepository, oddsRepository, spreadSheetRepository, raceEntityConverter, horseEntityConverter, predictionFilter, placeCheckList, odds)
	raceIdRepository := infrastructure.NewRaceIdRepository(netKeibaGateway)
	markerSync := prediction_service.NewMarkerSync(raceIdRepository, raceRepository, spreadSheetRepository)
	analysisFilter := filter_service.NewAnalysisFilter()
	place := analysis_service.NewPlace(analysisFilter, spreadSheetRepository)
	prediction := prediction_usecase.NewPrediction(odds, placeCandidate, markerSync, place)
	controllerPrediction := controller.NewPrediction(prediction)
	return controllerPrediction
}

// wire.go:

var MasterSet = wire.NewSet(master_usecase.NewMaster, master_service.NewTicket, master_service.NewRaceId, master_service.NewRace, master_service.NewJockey, master_service.NewWinOdds, master_service.NewPlaceOdds, master_service.NewTrioOdds, master_service.NewAnalysisMarker, master_service.NewPredictionMarker, master_service.NewBetNumberConverter, master_service.NewUmacaTicket, master_service.NewRaceForecast, converter.NewRaceEntityConverter, converter.NewJockeyEntityConverter, converter.NewOddsEntityConverter, converter.NewRaceForecastEntityConverter, infrastructure.NewTicketRepository, infrastructure.NewRaceIdRepository, infrastructure.NewRaceRepository, infrastructure.NewRaceForecastRepository, infrastructure.NewJockeyRepository, infrastructure.NewOddsRepository, infrastructure.NewAnalysisMarkerRepository, infrastructure.NewPredictionMarkerRepository, infrastructure.NewUmacaTicketRepository, gateway.NewNetKeibaGateway, gateway.NewNetKeibaCollector, gateway.NewTospoGateway)

var AggregationSet = wire.NewSet(aggregation_usecase.NewSummary, aggregation_usecase.NewTicketSummary, aggregation_usecase.NewList, aggregation_service.NewSummary, aggregation_service.NewTicketSummary, aggregation_service.NewList, summary_service.NewTerm, summary_service.NewTicket, summary_service.NewClass, summary_service.NewCourseCategory, summary_service.NewDistanceCategory, summary_service.NewRaceCourse, infrastructure.NewSpreadSheetRepository, converter.NewRaceEntityConverter, converter.NewJockeyEntityConverter)

var AnalysisSet = wire.NewSet(analysis_usecase.NewAnalysis, analysis_service.NewPlace, analysis_service.NewTrio, analysis_service.NewPlaceAllIn, analysis_service.NewPlaceUnHit, master_service.NewHorse, filter_service.NewAnalysisFilter, infrastructure.NewHorseRepository, infrastructure.NewSpreadSheetRepository, gateway.NewNetKeibaGateway, gateway.NewNetKeibaCollector, converter.NewHorseEntityConverter)

var PredictionSet = wire.NewSet(prediction_usecase.NewPrediction, prediction_service.NewOdds, prediction_service.NewPlaceCandidate, prediction_service.NewPlaceCheckList, prediction_service.NewMarkerSync, filter_service.NewPredictionFilter, infrastructure.NewOddsRepository, infrastructure.NewRaceRepository, infrastructure.NewJockeyRepository, infrastructure.NewTrainerRepository, infrastructure.NewRaceForecastRepository, infrastructure.NewRaceIdRepository, gateway.NewTospoGateway, converter.NewRaceEntityConverter)

var SpreadSheetGatewaySet = wire.NewSet(gateway.NewSpreadSheetSummaryGateway, gateway.NewSpreadSheetTicketSummaryGateway, gateway.NewSpreadSheetListGateway, gateway.NewSpreadSheetAnalysisPlaceGateway, gateway.NewSpreadSheetAnalysisPlaceAllInGateway, gateway.NewSpreadSheetPredictionOddsGateway, gateway.NewSpreadSheetPredictionCheckListGateway, gateway.NewSpreadSheetPredictionMarkerGateway)
