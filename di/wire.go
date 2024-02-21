//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/analysis_usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/data_cache_usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/list_usecase"
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
		infrastructure.NewRaceDataRepository,
		infrastructure.NewRacingNumberDataRepository,
		infrastructure.NewJockeyDataRepository,
		infrastructure.NewRaceIdDataRepository,
		infrastructure.NewMarkerDataRepository,
	)
	return nil
}

func InitializeMarkerAnalysisUseCase() *analysis_usecase.AnalysisUseCase {
	wire.Build(
		analysis_usecase.NewAnalysisUseCase,
		service.NewAnalysisService,
		service.NewRaceConverter,
		service.NewTicketConverter,
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
