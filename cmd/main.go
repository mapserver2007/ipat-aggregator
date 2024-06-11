package main

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/controller"
	"github.com/mapserver2007/ipat-aggregator/config"
	"github.com/mapserver2007/ipat-aggregator/di"
	"log"
)

const (
	enableAggregation = false
	enableAnalysis    = true
	enablePrediction  = false
)

func main() {
	ctx := context.Background()
	log.Println(ctx, "start")

	masterCtrl := di.NewMaster()
	master, err := masterCtrl.Execute(ctx, &controller.MasterInput{
		StartDate: config.RaceStartDate,
		EndDate:   config.RaceEndDate,
	})
	if err != nil {
		log.Println("master error")
		panic(err)
	}

	if enableAggregation {
		aggregationCtrl := di.NewAggregation()
		err = aggregationCtrl.Execute(ctx, &controller.AggregationInput{
			Master: master,
		})
		if err != nil {
			log.Println("aggregation error")
			panic(err)
		}
	}

	if enableAnalysis {
		analysisCtrl := di.NewAnalysis()
		err = analysisCtrl.Execute(ctx, &controller.AnalysisInput{
			Master: master,
		})
		if err != nil {
			log.Println("analysis error")
			panic(err)
		}
	}

	if enablePrediction {
		predictionCtrl := di.NewPrediction()
		err = predictionCtrl.Execute(ctx, &controller.PredictionInput{
			Master: master,
		})
		if err != nil {
			log.Println("prediction error")
			panic(err)
		}
	}

	log.Println(ctx, "end")
}
