package main

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure"
	"github.com/mapserver2007/ipat-aggregator/app/usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/analysis_usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/spreadsheet_usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/ticket_usecase"
	"github.com/mapserver2007/ipat-aggregator/di"
	"log"
)

const newProc = true

func main() {
	ctx := context.Background()

	if newProc {
		tickets2, racingNumbers2, races2, jockeys2, predictRaces, markers, err := masterFile(ctx)
		if err != nil {
			panic(err)
		}

		// 実験中
		analysis(ctx, markers, predictRaces, tickets2, racingNumbers2)

		err = summary(ctx, tickets2, racingNumbers2, races2, jockeys2)
		if err != nil {
			panic(err)
		}

		return
	}

	// 以下旧処理
	spreadSheetClient := infrastructure.NewSpreadSheetClient(ctx)
	spreadSheetMonthlyBettingTicketClient := infrastructure.NewSpreadSheetMonthlyBettingTicketClient(ctx)
	spreadSheetListClient := infrastructure.NewSpreadSheetListClient(ctx)
	spreadSheetAnalyzeClient := infrastructure.NewSpreadSheetAnalyzeClient(ctx)

	log.Println(ctx, "start")

	dataCacheUseCase := di.DataCacheInit()
	records, raceNumbers, races, jockeys, err := dataCacheUseCase.ReadAndUpdate(ctx)
	if err != nil {
		panic(err)
	}

	aggregator := di.AggregatorInit()
	summary := aggregator.GetSummary(records, raceNumbers, races)
	monthlyBettingTicketSummary := aggregator.GetMonthlyBettingTicketSummary(records, raceNumbers, races)

	predictor := di.PredictInit()
	predictResults, err := predictor.Predict(records, raceNumbers, races, jockeys)
	if err != nil {
		panic(err)
	}

	//analyzer := di.AnalyzerInit()
	//analyzeSummary := analyzer.WinAnalyze(records, raceNumberInfo.RacingNumbers(), raceInfo.Races())

	//spreadSheetUseCase := di.SpreadSheetInit()
	spreadSheetUseCase := usecase.NewSpreadSheet(spreadSheetClient, spreadSheetMonthlyBettingTicketClient, spreadSheetListClient, spreadSheetAnalyzeClient)
	err = spreadSheetUseCase.WriteSummary(ctx, summary)
	if err != nil {
		panic(err)
	}
	err = spreadSheetUseCase.WriteMonthlyBettingTicketSummary(ctx, monthlyBettingTicketSummary)
	if err != nil {
		panic(err)
	}
	err = spreadSheetUseCase.WriteStyleMonthlyBettingTicketSummary(ctx, monthlyBettingTicketSummary)
	if err != nil {
		panic(err)
	}
	styleMap, err := spreadSheetUseCase.WriteList(ctx, predictResults)
	if err != nil {
		panic(err)
	}
	err = spreadSheetUseCase.WriteStyleSummary(ctx, summary)
	if err != nil {
		panic(err)
	}
	err = spreadSheetUseCase.WriteStyleList(ctx, predictResults, styleMap)
	if err != nil {
		panic(err)
	}
	//err = spreadSheetUseCase.WriteAnalyze(ctx, analyzeSummary)
	//if err != nil {
	//	panic(err)
	//}
	//err = spreadSheetUseCase.WriteStyleAnalyze(ctx, analyzeSummary)
	//if err != nil {
	//	panic(err)
	//}

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
	predictAnalysisService := service.NewAnalysisService()
	ticketCsvRepository := infrastructure.NewTicketCsvRepository(betNumberConverter)
	markerDataRepository := infrastructure.NewMarkerDataRepository()
	ticketUseCase := ticket_usecase.NewTicket(ticketCsvRepository)
	analysisUseCase := analysis_usecase.NewAnalysis(markerDataRepository, predictAnalysisService, ticketConverter)

	tickets, err := ticketUseCase.Read(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dataCacheUseCase := di.InitializeDataCacheUseCase()

	racingNumbers, races, jockeys, excludeJockeyIds, raceIdMap, excludeDates, predictRaces, err := dataCacheUseCase.Read(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	err = dataCacheUseCase.Write(ctx, tickets, racingNumbers, races, jockeys, excludeJockeyIds, raceIdMap, excludeDates, predictRaces)
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
	tickets []*ticket_csv_entity.Ticket,
	racingNumbers []*data_cache_entity.RacingNumber,
) error {
	raceConverter := service.NewRaceConverter()
	ticketConverter := service.NewTicketConverter(raceConverter)
	predictAnalysisService := service.NewAnalysisService()
	markerDataRepository := infrastructure.NewMarkerDataRepository()
	predictUseCase := analysis_usecase.NewAnalysis(markerDataRepository, predictAnalysisService, ticketConverter)
	spreadSheetRepository, err := infrastructure.NewSpreadSheetMarkerAnalysisRepository()
	if err != nil {
		return err
	}

	analysisData, err := predictUseCase.CreateAnalysisData(ctx, markers, races, tickets, racingNumbers)
	if err != nil {
		return err
	}

	spreadSheetUseCase := spreadsheet_usecase.NewMarkerAnalysisUseCase(spreadSheetRepository, predictAnalysisService)
	spreadSheetUseCase.Write(ctx, analysisData)

	return nil
}

func summary(ctx context.Context, tickets []*ticket_csv_entity.Ticket, racingNumbers []*data_cache_entity.RacingNumber, races []*data_cache_entity.Race, jockeys []*data_cache_entity.Jockey) error {
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
