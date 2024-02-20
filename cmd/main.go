package main

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/analysis_usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/spreadsheet_usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/ticket_usecase"
	"github.com/mapserver2007/ipat-aggregator/di"
	"log"
)

const (
	predictRaceStartDate = "20230916"
	predictRaceEndDate   = "20240210"
)

func main() {
	ctx := context.Background()
	log.Println(ctx, "start")

	tickets2, racingNumbers2, races2, jockeys2, predictRaces, markers, err := masterFile(ctx)
	if err != nil {
		panic(err)
	}

	err = list(ctx, tickets2, racingNumbers2, races2, jockeys2)
	if err != nil {
		panic(err)
	}

	err = analysis(ctx, markers, predictRaces)
	if err != nil {
		panic(err)
	}

	err = summary(ctx, tickets2, racingNumbers2, races2, jockeys2)
	if err != nil {
		panic(err)
	}

	log.Println(ctx, "end")
}

func masterFile(
	ctx context.Context,
) (
	[]*ticket_csv_entity.Ticket,
	[]*data_cache_entity.RacingNumber,
	[]*data_cache_entity.Race,
	[]*data_cache_entity.Jockey,
	[]*data_cache_entity.Race,
	[]*marker_csv_entity.Yamato,
	error,
) {
	betNumberConverter := service.NewBetNumberConverter()
	raceConverter := service.NewRaceConverter()
	ticketConverter := service.NewTicketConverter(raceConverter)
	analysisService := service.NewAnalysisService()
	ticketCsvRepository := infrastructure.NewTicketCsvRepository(betNumberConverter)
	markerDataRepository := infrastructure.NewMarkerDataRepository()
	ticketUseCase := ticket_usecase.NewTicket(ticketCsvRepository)
	analysisUseCase := analysis_usecase.NewAnalysis(markerDataRepository, analysisService, ticketConverter)

	tickets, err := ticketUseCase.Read(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dataCacheUseCase := di.InitializeDataCacheUseCase()

	racingNumbers, races, jockeys, excludeJockeyIds, raceIdMap, excludeDates, predictRaces, err := dataCacheUseCase.Read(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	err = dataCacheUseCase.Write(ctx, tickets, racingNumbers, races, jockeys, excludeJockeyIds, raceIdMap, excludeDates, predictRaces, predictRaceStartDate, predictRaceEndDate)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	racingNumbers, races, jockeys, _, _, _, predictRaces, err = dataCacheUseCase.Read(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	markers, err := analysisUseCase.Read(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	return tickets, racingNumbers, races, jockeys, predictRaces, markers, nil
}

func analysis(
	ctx context.Context,
	markers []*marker_csv_entity.Yamato,
	races []*data_cache_entity.Race,
) error {
	raceConverter := service.NewRaceConverter()
	ticketConverter := service.NewTicketConverter(raceConverter)
	analysisService := service.NewAnalysisService()
	spreadSheetService := service.NewSpreadSheetService()
	markerDataRepository := infrastructure.NewMarkerDataRepository()
	analysisUseCase := analysis_usecase.NewAnalysis(markerDataRepository, analysisService, ticketConverter)
	spreadSheetRepository, err := infrastructure.NewSpreadSheetMarkerAnalysisRepository(spreadSheetService)
	if err != nil {
		return err
	}

	analysisData, searchFilters, err := analysisUseCase.CreateAnalysisData(ctx, markers, races)
	if err != nil {
		return err
	}

	spreadSheetUseCase := spreadsheet_usecase.NewMarkerAnalysisUseCase(spreadSheetRepository, analysisService)
	err = spreadSheetUseCase.Write(ctx, analysisData, searchFilters)
	if err != nil {
		return err
	}

	return nil
}

func list(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	racingNumbers []*data_cache_entity.RacingNumber,
	races []*data_cache_entity.Race,
	jockeys []*data_cache_entity.Jockey,
) error {
	listUseCase := di.InitializeListUseCase()
	rows, err := listUseCase.Read(ctx, tickets, racingNumbers, races, jockeys)
	if err != nil {
		return err
	}

	raceConverter := service.NewRaceConverter()
	ticketConverter := service.NewTicketConverter(raceConverter)
	raceEntityConverter := service.NewRaceEntityConverter()
	spreadSheetService := service.NewSpreadSheetService()
	listService := service.NewListService(raceConverter, ticketConverter, raceEntityConverter)
	spreadSheetRepository, err := infrastructure.NewSpreadSheetListRepository(spreadSheetService)
	spreadSheetUseCase := spreadsheet_usecase.NewListUseCase(listService, spreadSheetRepository)
	spreadSheetUseCase.Write(ctx, rows, jockeys)

	return nil
}

func summary(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	racingNumbers []*data_cache_entity.RacingNumber,
	races []*data_cache_entity.Race,
	jockeys []*data_cache_entity.Jockey,
) error {
	raceConverter := service.NewRaceConverter()
	ticketConverter := service.NewTicketConverter(raceConverter)
	ticketAggregator := service.NewTicketAggregator(ticketConverter)
	summaryService := service.NewSummaryService(ticketAggregator)
	spreadSheetRepository, err := infrastructure.NewSpreadSheetSummaryRepository()
	if err != nil {
		return err
	}

	spreadSheetUseCase := spreadsheet_usecase.NewSummaryUseCase(summaryService, spreadSheetRepository)
	err = spreadSheetUseCase.Write(ctx, tickets, racingNumbers, races)
	if err != nil {
		return err
	}

	return nil
}
