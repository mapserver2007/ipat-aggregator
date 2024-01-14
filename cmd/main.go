package main

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure"
	"github.com/mapserver2007/ipat-aggregator/app/usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/spreadsheet_usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/ticket_usecase"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/yamato_predict_usecase"
	"github.com/mapserver2007/ipat-aggregator/di"
	"log"
)

const newProc = true

func main() {
	ctx := context.Background()

	if newProc {
		tickets2, racingNumbers2, races2, jockeys2, predictRaces, predicts, err := masterFile(ctx)
		_ = predictRaces
		_ = predicts
		if err != nil {
			panic(err)
		}

		// 実験中
		predict(ctx)

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
	[]*data_cache_entity.Predict,
	error,
) {
	betNumberConverter := service.NewBetNumberConverter()
	ticketCsvRepository := infrastructure.NewTicketCsvRepository(betNumberConverter)
	ticketUseCase := ticket_usecase.NewTicket(ticketCsvRepository)
	tickets, err := ticketUseCase.Read(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dataCacheUseCase := di.InitializeDataCacheUseCase()

	racingNumbers, races, jockeys, excludeJockeyIds, raceIdMap, excludeDates, predictRaces, _, err := dataCacheUseCase.Read(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	err = dataCacheUseCase.Write(ctx, tickets, racingNumbers, races, jockeys, excludeJockeyIds, raceIdMap, excludeDates, predictRaces)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	racingNumbers, races, jockeys, _, _, _, predictRaces, predicts, err := dataCacheUseCase.Read(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	return tickets, racingNumbers, races, jockeys, predictRaces, predicts, nil
}

func predict(ctx context.Context) {
	raceConverter := service.NewRaceConverter()
	netKeibaService := service.NewNetKeibaService(raceConverter)
	raceIdRepository := infrastructure.NewRaceIdDataRepository()
	predictUseCase := yamato_predict_usecase.NewPredict(netKeibaService, raceIdRepository)
	_ = predictUseCase.Fetch(ctx)

	//dataCacheUseCase := di.InitializeDataCacheUseCase()
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

	summaryUseCase := spreadsheet_usecase.NewSummaryUseCase(summaryService, spreadSheetRepository)

	err = summaryUseCase.Write(ctx, tickets, racingNumbers, races)
	if err != nil {
		return err
	}

	return nil
}
