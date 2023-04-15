package main

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure"
	"github.com/mapserver2007/ipat-aggregator/app/usecase"
	"github.com/mapserver2007/ipat-aggregator/di"
	"log"
)

func main() {
	ctx := context.Background()
	spreadSheetClient := infrastructure.NewSpreadSheetClient(ctx)
	spreadSheetListClient := infrastructure.NewSpreadSheetListClient(ctx)

	log.Println(ctx, "start")

	dataCacheUseCase := di.DataCacheInit()
	records, raceNumberInfo, raceInfo, jockeyInfo, err := dataCacheUseCase.ReadAndUpdate(ctx)
	if err != nil {
		panic(err)
	}

	aggregator := di.AggregatorInit()
	summary := aggregator.GetSummary(records, raceNumberInfo.RacingNumbers(), raceInfo.Races())

	predictor := di.PredictInit()
	predictResults, err := predictor.Predict(records, raceNumberInfo.RacingNumbers(), raceInfo.Races(), jockeyInfo.Jockeys())
	if err != nil {
		panic(err)
	}

	//analyser := di.AnalyserInit()
	//analyser.Analyse(records, raceNumberInfo.RacingNumbers(), raceInfo.Races())

	//analyser := service.NewAnalyser()
	//analyser.Analyse(records)

	//spreadSheetUseCase := di.SpreadSheetInit()
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
