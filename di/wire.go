//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure"
	service2 "github.com/mapserver2007/ipat-aggregator/app/service"
	"github.com/mapserver2007/ipat-aggregator/app/usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/data_cache_usecase"
)

func InitializeDataCacheUseCase() *data_cache_usecase.DataCacheUseCase {
	wire.Build(
		data_cache_usecase.NewDataCacheUseCase,
		service.NewRaceConverter,
		service.NewNetKeibaService,
		service.NewRacingNumberEntityConverter,
		service.NewRaceEntityConverter,
		service.NewJockeyEntityConverter,
		infrastructure.NewRaceDataRepository,
		infrastructure.NewRacingNumberDataRepository,
		infrastructure.NewJockeyDataRepository,
		infrastructure.NewRaceIdDataRepository,
		infrastructure.NewMarkerDataRepository,
	)
	return nil
}

//func InitializeSummaryUseCase() *spreadsheet_usecase.SummaryUseCase {
//	wire.Build(
//		spreadsheet_usecase.NewSummaryUseCase,
//		service.NewSummaryService,
//		service.NewTicketAggregator,
//		service.NewTicketConverter,
//		service.NewRaceConverter,
//		infrastructure.NewSpreadSheetSummaryRepository,
//	)
//	return nil
//}

// 以下古い

func DataCacheInit() *usecase.DataCache {
	wire.Build(
		usecase.NewDataCache,
		service2.NewCsvReader,
		service2.NewRaceFetcher,
		service2.NewRaceConverter,
		infrastructure.NewRaceDB,
		infrastructure.NewRaceClient,
		//service2.NewBettingTicketConverter,
	)
	return nil
}

func AggregatorInit() *service2.Aggregator {
	wire.Build(
		service2.NewAggregator,
		service2.NewRaceConverter,
		service2.NewBettingTicketConverter,
		service2.NewSummarizer,
	)
	return nil
}

func PredictInit() *service2.Predictor {
	wire.Build(
		service2.NewPredictor,
		service2.NewRaceConverter,
		service2.NewBettingTicketConverter,
	)
	return nil
}

func AnalyzerInit() *usecase.Analyzer {
	wire.Build(
		usecase.NewAnalyzer,
		service2.NewAnalyzer,
		service2.NewRaceConverter,
	)
	return nil
}
