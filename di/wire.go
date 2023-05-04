//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure"
	"github.com/mapserver2007/ipat-aggregator/app/service"
	"github.com/mapserver2007/ipat-aggregator/app/usecase"
)

func DataCacheInit() *usecase.DataCache {
	wire.Build(
		usecase.NewDataCache,
		service.NewCsvReader,
		service.NewRaceFetcher,
		service.NewRaceConverter,
		infrastructure.NewRaceDB,
		infrastructure.NewRaceClient,
	)
	return nil
}

func AggregatorInit() *service.Aggregator {
	wire.Build(
		service.NewAggregator,
		service.NewRaceConverter,
		service.NewBettingTicketConverter,
	)
	return nil
}

func PredictInit() *service.Predictor {
	wire.Build(
		service.NewPredictor,
		service.NewRaceConverter,
		service.NewBettingTicketConverter,
	)
	return nil
}

func AnalyserInit() *usecase.Analyser {
	wire.Build(
		usecase.NewAnalyser,
		service.NewAnalyser,
		service.NewRaceConverter,
	)
	return nil
}
