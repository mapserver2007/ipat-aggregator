package main

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/spreadsheet_usecase"
	"github.com/mapserver2007/ipat-aggregator/di"
	"log"
)

const (
	analysisRaceStartDate = "20230819"
	analysisRaceEndDate   = "20240309"
	enableAnalysis        = true
	enablePrediction      = false
	enableAggregate       = false
)

func main() {
	ctx := context.Background()
	log.Println(ctx, "start")

	tickets, racingNumbers, ticketRaces, jockeys, analysisRaces, markers, err := masterFile(ctx)
	if err != nil {
		panic(err)
	}

	if enablePrediction {
		err = prediction(ctx, markers, analysisRaces)
		if err != nil {
			panic(err)
		}
	}

	if enableAnalysis {
		err = analysis(ctx, markers, analysisRaces)
		if err != nil {
			panic(err)
		}
	}

	if enableAggregate {
		err = list(ctx, tickets, racingNumbers, ticketRaces, jockeys)
		if err != nil {
			panic(err)
		}

		err = summary(ctx, tickets, racingNumbers, ticketRaces)
		if err != nil {
			panic(err)
		}
	}

	log.Println(ctx, "end")
}

func prediction(
	ctx context.Context,
	markers []*marker_csv_entity.AnalysisMarker,
	races []*data_cache_entity.Race,
) error {
	predictionUseCase := di.InitializePredictionUseCase()
	analysisUseCase := di.InitializeMarkerAnalysisUseCase()
	spreadSheetService := service.NewSpreadSheetService()
	filterService := service.NewFilterService()
	spreadSheetRepository, err := infrastructure.NewSpreadSheetPredictionRepository(spreadSheetService)
	if err != nil {
		return err
	}

	analysisData, err := analysisUseCase.CreateAnalysisData(ctx, markers, races)
	if err != nil {
		return err
	}

	predictionMarkers, err := predictionUseCase.Read(ctx)
	if err != nil {
		return err
	}
	var predictionRaceIds []types.RaceId
	predictionMarkerMap := map[types.RaceId]*marker_csv_entity.PredictionMarker{}
	for _, predictionMarker := range predictionMarkers {
		predictionMarkerMap[predictionMarker.RaceId()] = predictionMarker
		predictionRaceIds = append(predictionRaceIds, predictionMarker.RaceId())
	}

	predictionRaces, err := predictionUseCase.Fetch(ctx, predictionRaceIds)
	if err != nil {
		return err
	}

	spreadSheetUseCase := spreadsheet_usecase.NewPredictionUseCase(spreadSheetRepository, filterService, spreadSheetService)
	err = spreadSheetUseCase.Write(ctx, predictionRaces, predictionMarkerMap, analysisData)
	if err != nil {
		return err
	}

	return nil
}

func masterFile(
	ctx context.Context,
) (
	[]*ticket_csv_entity.Ticket,
	[]*data_cache_entity.RacingNumber,
	[]*data_cache_entity.Race,
	[]*data_cache_entity.Jockey,
	[]*data_cache_entity.Race,
	[]*marker_csv_entity.AnalysisMarker,
	error,
) {
	ticketUseCase := di.InitializeTicketUseCase()
	analysisUseCase := di.InitializeMarkerAnalysisUseCase()

	tickets, err := ticketUseCase.Read(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dataCacheUseCase := di.InitializeDataCacheUseCase()

	racingNumbers, ticketRaces, jockeys, excludeJockeyIds, raceIdMap, excludeDates, analysisRaces, err := dataCacheUseCase.Read(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	err = dataCacheUseCase.Write(ctx, tickets, racingNumbers, ticketRaces, jockeys, excludeJockeyIds, raceIdMap, excludeDates, analysisRaces, analysisRaceStartDate, analysisRaceEndDate)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	racingNumbers, ticketRaces, jockeys, _, _, _, analysisRaces, err = dataCacheUseCase.Read(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	markers, err := analysisUseCase.Read(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	return tickets, racingNumbers, ticketRaces, jockeys, analysisRaces, markers, nil
}

func analysis(
	ctx context.Context,
	markers []*marker_csv_entity.AnalysisMarker,
	races []*data_cache_entity.Race,
) error {
	spreadSheetService := service.NewSpreadSheetService()
	analysisService := service.NewAnalysisService(spreadSheetService)
	filterService := service.NewFilterService()
	analysisUseCase := di.InitializeMarkerAnalysisUseCase()
	spreadSheetRepository, err := infrastructure.NewSpreadSheetMarkerAnalysisRepository(spreadSheetService)
	if err != nil {
		return err
	}

	analysisData, err := analysisUseCase.CreateAnalysisData(ctx, markers, races)
	if err != nil {
		return err
	}

	spreadSheetUseCase := spreadsheet_usecase.NewMarkerAnalysisUseCase(spreadSheetRepository, analysisService, filterService)
	err = spreadSheetUseCase.Write(ctx, analysisData)
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

func ticketSummary(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	racingNumbers []*data_cache_entity.RacingNumber,
	races []*data_cache_entity.Race,
) error {

	return nil
}
