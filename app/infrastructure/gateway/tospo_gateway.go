package gateway

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"sync"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/tospo_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/sirupsen/logrus"
)

type TospoGateway interface {
	FetchForecast(ctx context.Context, url string) ([]*tospo_entity.Forecast, error)
	FetchTrainingComment(ctx context.Context, url string) ([]*tospo_entity.TrainingComment, error)
	FetchReporterMemo(ctx context.Context, url string) ([]*tospo_entity.ReporterMemo, error)
	FetchPaddockComment(ctx context.Context, url string) ([]*tospo_entity.PaddockComment, error)
}

type tospoGateway struct {
	logger *logrus.Logger
	mu     sync.Mutex
}

func NewTospoGateway(
	logger *logrus.Logger,
) TospoGateway {
	return &tospoGateway{
		logger: logger,
	}
}

func (t *tospoGateway) FetchForecast(
	ctx context.Context,
	url string,
) ([]*tospo_entity.Forecast, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.logger.Infof("fetching forecast from %s", url)
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
	if err = json.Unmarshal(body, &rawForecast); err != nil {
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
	t.mu.Lock()
	defer t.mu.Unlock()

	t.logger.Infof("fetching training comment from %s", url)
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
	if err = json.Unmarshal(body, &rawTrainingComment); err != nil {
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

func (t *tospoGateway) FetchReporterMemo(
	ctx context.Context,
	url string,
) ([]*tospo_entity.ReporterMemo, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.logger.Infof("fetching reporter memo from %s", url)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var rawReporterMemo *raw_entity.ReporterMemo
	if err = json.Unmarshal(body, &rawReporterMemo); err != nil {
		return nil, err
	}

	horseNumberMemoMap := map[int][]*tospo_entity.Memo{}
	for _, reporterMemo := range rawReporterMemo.Body.ReceivedMemoList {
		if _, ok := horseNumberMemoMap[reporterMemo.HorseNumber]; !ok {
			horseNumberMemoMap[reporterMemo.HorseNumber] = []*tospo_entity.Memo{}
		}
		memo, err := tospo_entity.NewMemo(
			reporterMemo.Comment,
			reporterMemo.Date,
		)
		if err != nil {
			return nil, err
		}
		horseNumberMemoMap[reporterMemo.HorseNumber] = append(horseNumberMemoMap[reporterMemo.HorseNumber], memo)
	}

	reporterMemos := make([]*tospo_entity.ReporterMemo, 0, len(horseNumberMemoMap))
	for horseNumber, memos := range horseNumberMemoMap {
		reporterMemos = append(reporterMemos, tospo_entity.NewReporterMemo(
			types.HorseNumber(horseNumber),
			memos,
		))
	}

	return reporterMemos, nil
}

func (t *tospoGateway) FetchPaddockComment(
	ctx context.Context,
	url string,
) ([]*tospo_entity.PaddockComment, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.logger.Infof("fetching paddock comment from %s", url)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var rawPaddockCommentInfo *raw_entity.PaddockCommentInfo
	if err = json.Unmarshal(body, &rawPaddockCommentInfo); err != nil {
		return nil, err
	}

	var paddockComments []*tospo_entity.PaddockComment
	for _, racePaddockComment := range rawPaddockCommentInfo.Body.RaceEntryList {
		if racePaddockComment.Comment != "" && racePaddockComment.Evaluation > 0 {
			paddockComments = append(paddockComments, tospo_entity.NewPaddockComment(
				types.HorseNumber(racePaddockComment.HorseNumber),
				racePaddockComment.Comment,
				racePaddockComment.Evaluation,
			))
		}
	}

	sort.Slice(paddockComments, func(i, j int) bool {
		return paddockComments[i].HorseNumber() < paddockComments[j].HorseNumber()
	})

	return paddockComments, nil
}
