package main

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure"
	"github.com/mapserver2007/ipat-aggregator/app/usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/data_cache_usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/spreadsheet_usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/ticket_usecase"
	"github.com/mapserver2007/ipat-aggregator/di"
	"log"
)

func main() {
	ctx := context.Background()
	if sub(ctx) {
		return
	}
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

func sub(ctx context.Context) bool {
	// TODO DI

	betNumberConverter := service.NewBetNumberConverter()
	ticketCsvRepository := infrastructure.NewTicketCsvRepository(betNumberConverter)
	ticketUseCase := ticket_usecase.NewTicket(ticketCsvRepository)
	tickets, err := ticketUseCase.Read(ctx)
	if err != nil {
		panic(err)
	}

	raceConverter := service.NewRaceConverter()
	ticketConverter := service.NewTicketConverter()
	ticketAggregator := service.NewTicketAggregator(ticketConverter)
	netKeibaService := service.NewNetKeibaService(raceConverter)
	summaryService := service.NewSummaryService(ticketAggregator)
	racingNumberEntityConverter := service.NewRacingNumberEntityConverter()
	raceEntityConverter := service.NewRaceEntityConverter()
	jockeyEntityConverter := service.NewJockeyEntityConverter()
	racingNumberRepository := infrastructure.NewRacingNumberDataRepository()
	raceDataRepository := infrastructure.NewRaceDataRepository()
	jockeyDataRepository := infrastructure.NewJockeyDataRepository()
	spreadSheetRepository := infrastructure.NewSpreadSheetSummaryRepository()
	dataCacheUseCase := data_cache_usecase.NewDataCacheUseCase(racingNumberRepository, raceDataRepository, jockeyDataRepository, netKeibaService, raceConverter, racingNumberEntityConverter, raceEntityConverter, jockeyEntityConverter)
	summaryUseCase := spreadsheet_usecase.NewSummaryUseCase(summaryService, spreadSheetRepository)

	racingNumbers, races, jockeys, excludeJockeyIds, err := dataCacheUseCase.Read(ctx)
	if err != nil {
		panic(err)
	}

	err = dataCacheUseCase.Write(ctx, tickets, racingNumbers, races, jockeys, excludeJockeyIds)
	if err != nil {
		panic(err)
	}

	racingNumbers, races, jockeys, excludeJockeyIds, err = dataCacheUseCase.Read(ctx)
	if err != nil {
		panic(err)
	}

	err = summaryUseCase.Write(ctx, tickets)
	if err != nil {
		panic(err)
	}

	_ = racingNumbers
	_ = races
	_ = excludeJockeyIds
	_ = jockeys

	return true

}
