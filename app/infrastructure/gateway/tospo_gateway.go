package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/tospo_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"io"
	"log"
	"net/http"
	"sort"
)

type TospoGateway interface {
	FetchForecast(ctx context.Context, url string) ([]*tospo_entity.Forecast, error)
	FetchTrainingComment(ctx context.Context, url string) ([]*tospo_entity.TrainingComment, error)
}

type tospoGateway struct{}

func NewTospoGateway() TospoGateway {
	return &tospoGateway{}
}

func (t *tospoGateway) FetchForecast(
	ctx context.Context,
	url string,
) ([]*tospo_entity.Forecast, error) {
	log.Println(ctx, fmt.Sprintf("fetching forecast from %s", url))
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var rawForecast *raw_entity.ForecastInfo
	if err := json.Unmarshal(body, &rawForecast); err != nil {
		return nil, err
	}

	var raceForecastDataMap map[string]map[string]raw_entity.RaceForecastData
	err = json.Unmarshal(rawForecast.Body.RaceForecastDataList, &raceForecastDataMap)
	if err != nil {
		return nil, err
	}

	horseNameMap := map[string]types.HorseNumber{}
	for _, raceEntry := range rawForecast.Body.RaceEntries {
		horseNameMap[raceEntry.HorseName] = types.HorseNumber(raceEntry.HorseNumber)
	}

	markerNum := len(raceForecastDataMap)
	favoriteMarkerMap := map[types.HorseNumber]int{}
	rivalMarkerMap := map[types.HorseNumber]int{}

	for _, raceForecastData := range raceForecastDataMap {
		for _, forecastData := range raceForecastData {
			horseNumber, ok := horseNameMap[forecastData.HorseName]
			if ok {
				switch forecastData.ReporterMarkType {
				case 2:
					favoriteMarkerMap[horseNumber]++
				case 3:
					rivalMarkerMap[horseNumber]++
				}
			}
		}
	}

	forecasts := make([]*tospo_entity.Forecast, 0, len(horseNameMap))
	for _, horseNumber := range horseNameMap {
		var favoriteMarkerNum, rivalMarkerNum int
		if num, ok := favoriteMarkerMap[horseNumber]; ok {
			favoriteMarkerNum = num
		}
		if num, ok := rivalMarkerMap[horseNumber]; ok {
			rivalMarkerNum = num
		}
		forecasts = append(forecasts, tospo_entity.NewForecast(
			horseNumber,
			favoriteMarkerNum,
			rivalMarkerNum,
			markerNum,
		))
	}

	sort.Slice(forecasts, func(i, j int) bool {
		return forecasts[i].HorseNumber() < forecasts[j].HorseNumber()
	})

	return forecasts, nil
}

func (t *tospoGateway) FetchTrainingComment(
	ctx context.Context,
	url string,
) ([]*tospo_entity.TrainingComment, error) {
	log.Println(ctx, fmt.Sprintf("fetching training comment from %s", url))
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var rawTrainingComment *raw_entity.TrainingComment
	if err := json.Unmarshal(body, &rawTrainingComment); err != nil {
		return nil, err
	}

	// tospo側の問題でBAD REQUESTとなりデータが取れない場合があり、エラーにはせずnilで返す
	if rawTrainingComment.Body == nil {
		return nil, nil
	}

	trainingComments := make([]*tospo_entity.TrainingComment, 0, len(rawTrainingComment.Body.RaceTrainingComments))
	for _, raceTrainingComment := range rawTrainingComment.Body.RaceTrainingComments {
		previousTrainingComment := ""
		if raceTrainingComment.RaceHistoryCommentInfo != nil {
			previousTrainingComment = raceTrainingComment.RaceHistoryCommentInfo.TrainingComment
		}

		trainingComments = append(trainingComments, tospo_entity.NewTrainingComment(
			raceTrainingComment.HorseNumber,
			raceTrainingComment.TrainingComment,
			previousTrainingComment,
			raceTrainingComment.Prediction,
		))
	}

	sort.Slice(trainingComments, func(i, j int) bool {
		return trainingComments[i].HorseNumber() < trainingComments[j].HorseNumber()
	})

	return trainingComments, nil
}
