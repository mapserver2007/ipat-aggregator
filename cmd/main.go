package main

import (
	"context"
	"github.com/gocolly/colly"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure"
	"github.com/mapserver2007/ipat-aggregator/app/service"
	"github.com/mapserver2007/ipat-aggregator/app/usecase"
	"log"
)

func main() {
	ctx := context.Background()
	csvReader := service.NewCsvReader()
	collector := colly.NewCollector()
	raceClient := infrastructure.NewRaceClient(collector)
	raceFetcher := service.NewRaceFetcher(raceClient)
	raceConverter := service.NewRaceConverter()
	bettingTicketConverter := service.NewBettingTicketConverter()

	raceDB := infrastructure.NewRaceDB(raceClient)
	spreadSheetClient := infrastructure.NewSpreadSheetClient(ctx, "secret.json", "spreadsheet_calc.json")
	spreadSheetListClient := infrastructure.NewSpreadSheetListClient(ctx, "secret.json", "spreadsheet_list.json")

	log.Println(ctx, "start")

	dataCacheUseCase := usecase.NewDataCache(csvReader, raceDB, raceFetcher, raceConverter)

	records, raceNumberInfo, raceInfo, err := dataCacheUseCase.ReadAndUpdate(ctx)
	if err != nil {
		panic(err)
	}

	aggregator := service.NewAggregator(raceConverter, bettingTicketConverter, records, raceNumberInfo, raceInfo)
	summary := aggregator.GetSummary()

	predictor := service.NewPredictor(raceConverter, bettingTicketConverter, records, raceNumberInfo, raceInfo)
	predictResults, err := predictor.Predict()
	if err != nil {
		panic(err)
	}

	analyser := service.NewAnalyser(records)
	analyser.Analyse()

	spreadSheetUseCase := usecase.NewSpreadSheet(spreadSheetClient, spreadSheetListClient)
	err = spreadSheetUseCase.WriteSummary(ctx, summary)
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

	log.Println(ctx, "end")
}
