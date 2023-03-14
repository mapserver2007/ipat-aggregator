package main

import (
	"context"
	"github.com/gocolly/colly"
	"github.com/mapserver2007/tools/baken/app/infrastructure"
	"github.com/mapserver2007/tools/baken/app/service"
	"github.com/mapserver2007/tools/baken/app/usecase"
	"log"
)

func main() {
	ctx := context.Background()
	csvReader := service.NewCsvReader()
	collector := colly.NewCollector()
	raceClient := infrastructure.NewRaceClient(collector)
	raceFetcher := service.NewRaceFetcher(raceClient)
	raceConverter := service.NewRaceConverter()

	raceDB := infrastructure.NewRaceDB(raceClient)
	spreadSheetClient := infrastructure.NewSpreadSheetClient(ctx, "secret.json", "spreadsheet_calc.json")
	spreadSheetListClient := infrastructure.NewSpreadSheetListClient(ctx, "secret.json", "spreadsheet_list.json")

	log.Println(ctx, "start")

	dataCacheUseCase := usecase.NewDataCache(csvReader, raceDB, raceFetcher, raceConverter)
	//raceNumberInfo, raceInfo, err := dataCacheUseCase.ReadCache(ctx)

	entities, raceNumberInfo, raceInfo, err := dataCacheUseCase.ReadAndUpdate(ctx)
	if err != nil {
		panic(err)
	}

	aggregator := service.NewAggregator(raceConverter, entities, raceNumberInfo, raceInfo)
	summary := aggregator.GetSummary()

	predictor := service.NewPredictor(raceConverter, entities, raceInfo.Races)
	predictResults, err := predictor.Predict()
	if err != nil {
		panic(err)
	}

	analyser := service.NewAnalyser(entities)
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
