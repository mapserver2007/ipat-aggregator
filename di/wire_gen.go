// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure"
	service2 "github.com/mapserver2007/ipat-aggregator/app/service"
	"github.com/mapserver2007/ipat-aggregator/app/usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/data_cache_usecase"
)

// Injectors from wire.go:

func InitializeDataCacheUseCase() *data_cache_usecase.DataCacheUseCase {
	racingNumberDataRepository := infrastructure.NewRacingNumberDataRepository()
	raceDataRepository := infrastructure.NewRaceDataRepository()
	jockeyDataRepository := infrastructure.NewJockeyDataRepository()
	raceIdDataRepository := infrastructure.NewRaceIdDataRepository()
	markerDataRepository := infrastructure.NewMarkerDataRepository()
	raceConverter := service.NewRaceConverter()
	netKeibaService := service.NewNetKeibaService(raceConverter)
	racingNumberEntityConverter := service.NewRacingNumberEntityConverter()
	raceEntityConverter := service.NewRaceEntityConverter()
	jockeyEntityConverter := service.NewJockeyEntityConverter()
	dataCacheUseCase := data_cache_usecase.NewDataCacheUseCase(racingNumberDataRepository, raceDataRepository, jockeyDataRepository, raceIdDataRepository, markerDataRepository, netKeibaService, raceConverter, racingNumberEntityConverter, raceEntityConverter, jockeyEntityConverter)
	return dataCacheUseCase
}

func DataCacheInit() *usecase.DataCache {
	csvReader := service2.NewCsvReader()
	raceClient := infrastructure.NewRaceClient()
	raceDB := infrastructure.NewRaceDB(raceClient)
	raceFetcher := service2.NewRaceFetcher(raceClient)
	raceConverter := service2.NewRaceConverter()
	dataCache := usecase.NewDataCache(csvReader, raceDB, raceFetcher, raceConverter)
	return dataCache
}

func AggregatorInit() *service2.Aggregator {
	raceConverter := service2.NewRaceConverter()
	bettingTicketConverter := service2.NewBettingTicketConverter(raceConverter)
	summarizer := service2.NewSummarizer(raceConverter, bettingTicketConverter)
	aggregator := service2.NewAggregator(raceConverter, bettingTicketConverter, summarizer)
	return aggregator
}

func PredictInit() *service2.Predictor {
	raceConverter := service2.NewRaceConverter()
	bettingTicketConverter := service2.NewBettingTicketConverter(raceConverter)
	predictor := service2.NewPredictor(raceConverter, bettingTicketConverter)
	return predictor
}

func AnalyzerInit() *usecase.Analyzer {
	raceConverter := service2.NewRaceConverter()
	analyzer := service2.NewAnalyzer(raceConverter)
	usecaseAnalyzer := usecase.NewAnalyzer(analyzer)
	return usecaseAnalyzer
}
