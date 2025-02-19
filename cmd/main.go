package main

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/controller"
	"github.com/mapserver2007/ipat-aggregator/config"
	"github.com/mapserver2007/ipat-aggregator/di"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	timeout         = 20 * time.Minute
	scheduleSetting = "*/3 * * * *"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&SLF4JFormatter{})

	logFile, err := os.OpenFile("/tmp/ipat-aggregator.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
		return
	}
	defer logFile.Close()

	logger.SetOutput(io.MultiWriter(os.Stdout, logFile))

	masterCtrl := di.NewMaster(logger)
	master, err := masterCtrl.Execute(ctx, &controller.MasterInput{
		StartDate: config.RaceStartDate,
		EndDate:   config.RaceEndDate,
	})

	if err != nil {
		logger.Errorf("master raed error: %v", err)
		return
	}

	app := cli.NewApp()
	app.Name = "ipat-aggregator-cli"

	app.Commands = []cli.Command{
		{
			Name:    "aggregation",
			Aliases: []string{"g"},
			Usage:   "aggregation",
			Action: func(c *cli.Context) error {
				logger.Infof("aggregation start")
				aggregationCtrl := di.NewAggregation(logger)
				aggregationCtrl.Execute(ctx, &controller.AggregationInput{
					Master: master,
				})
				logger.Infof("aggregation end")
				return nil
			},
		},
		{
			Name:    "analysis-place",
			Aliases: []string{"ap1"},
			Usage:   "analysis-place",
			Action: func(c *cli.Context) error {
				logger.Infof("analysis place start")
				analysisCtrl := di.NewAnalysis(logger)
				analysisCtrl.Place(ctx, &controller.AnalysisInput{
					Master: master,
				})
				logger.Infof("analysis place end")
				return nil
			},
		},
		{
			Name:    "analysis-place-all-in",
			Aliases: []string{"ap2"},
			Usage:   "analysis-place-all-in",
			Action: func(c *cli.Context) error {
				logger.Infof("analysis place all in start")
				analysisCtrl := di.NewAnalysis(logger)
				analysisCtrl.PlaceAllIn(ctx, &controller.AnalysisInput{
					Master: master,
				})
				logger.Infof("analysis place all in end")
				return nil
			},
		},
		{
			Name:    "analysis-beta",
			Aliases: []string{"ap4"},
			Usage:   "analysis-beta",
			Action: func(c *cli.Context) error {
				logger.Infof("analysis beta in start")
				analysisCtrl := di.NewAnalysis(logger)
				analysisCtrl.Beta(ctx, &controller.AnalysisInput{
					Master: master,
				})
				logger.Infof("analysis beta in end")
				return nil
			},
		},
		{
			Name:    "prediction",
			Aliases: []string{"p1"},
			Usage:   "prediction",
			Action: func(c *cli.Context) error {
				logger.Infof("prediction start")
				predictionCtrl := di.NewPrediction(logger)
				predictionCtrl.Prediction(ctx, &controller.PredictionInput{
					Master: master,
				})
				logger.Infof("prediction end")
				return nil
			},
		},
		{
			Name:    "sync marker",
			Aliases: []string{"p2"},
			Usage:   "sync marker",
			Action: func(c *cli.Context) error {
				logger.Infof("sync marker start")
				predictionCtrl := di.NewPrediction(logger)
				predictionCtrl.SyncMarker(ctx)
				logger.Infof("sync marker end")
				return nil
			},
		},
	}

	app.Run(os.Args)

	//scheduler, err := func() (gocron.Scheduler, error) {
	//	jst, err := time.LoadLocation("Asia/Tokyo")
	//	if err != nil {
	//		logger.Errorf("location error: %v", err)
	//		return nil, err
	//	}
	//
	//	s, err := gocron.NewScheduler(gocron.WithLocation(jst))
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	return s, nil
	//}()
	//
	//if err != nil {
	//	logger.Errorf("failed to create scheduler, %v", err)
	//	return
	//}
	//
	//scheduler.Start()
	//defer func() {
	//	err = scheduler.Shutdown()
	//	logger.Infof("scheduler.Shutdown")
	//	if err != nil {
	//		logger.Errorf("scheduler.Shutdown error: %v", err)
	//		return
	//	}
	//}()
	//
	//_, err = scheduler.NewJob(
	//	gocron.CronJob(scheduleSetting, false),
	//	gocron.NewTask(func() {
	//		logger.Infof("scheduler processing start")
	//
	//		taskCtx, taskCancel := context.WithTimeout(ctx, timeout)
	//		defer taskCancel()
	//
	//		done := make(chan error, 1)
	//		go func() {
	//			err = app.Run(os.Args)
	//			done <- err
	//		}()
	//
	//		select {
	//		case err = <-done:
	//			if err != nil {
	//				logger.Errorf("app.Run error: %v", err)
	//			}
	//		case <-taskCtx.Done():
	//			logger.Warnf("task timed out and was canceled: %v", taskCtx.Err())
	//		}
	//		logger.Infof("scheduler processing end")
	//	}),
	//)
	//
	//if err != nil {
	//	logger.Errorf("scheduler.NewJob error: %v", err)
	//	return
	//}

	//<-ctx.Done()
	//logger.Infof("Interrupted cli: %v", ctx.Err())
	//err = scheduler.Shutdown()
	//if err != nil {
	//	logger.Errorf("scheduler.Shutdown error: %v", err)
	//}
}
