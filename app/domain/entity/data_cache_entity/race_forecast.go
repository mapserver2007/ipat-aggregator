package data_cache_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type RaceForecast struct {
	raceId    types.RaceId
	raceDate  types.RaceDate
	forecasts []*Forecast
}

func NewRaceForecast(
	rawRaceId string,
	rawRaceDate int,
	forecasts []*Forecast,
) *RaceForecast {
	return &RaceForecast{
		raceId:    types.RaceId(rawRaceId),
		raceDate:  types.RaceDate(rawRaceDate),
		forecasts: forecasts,
	}
}

func (r *RaceForecast) RaceId() types.RaceId {
	return r.raceId
}

func (r *RaceForecast) RaceDate() types.RaceDate {
	return r.raceDate
}

func (r *RaceForecast) Forecasts() []*Forecast {
	return r.forecasts
}

type Forecast struct {
	horseNumber             types.HorseNumber
	trainingComment         string
	previousTrainingComment string
	highlyRecommended       bool
	favoriteNum             int
	rivalNum                int
	markerNum               int
}

func NewForecast(
	horseNumber int,
	trainingComment string,
	previousTrainingComment string,
	highlyRecommended bool,
	favoriteNum int,
	rivalNum int,
	markerNum int,
) *Forecast {
	return &Forecast{
		horseNumber:             types.HorseNumber(horseNumber),
		trainingComment:         trainingComment,
		previousTrainingComment: previousTrainingComment,
		highlyRecommended:       highlyRecommended,
		favoriteNum:             favoriteNum,
		rivalNum:                rivalNum,
		markerNum:               markerNum,
	}
}

func (f *Forecast) HorseNumber() types.HorseNumber {
	return f.horseNumber
}

func (f *Forecast) TrainingComment() string {
	return f.trainingComment
}

func (f *Forecast) PreviousTrainingComment() string {
	return f.previousTrainingComment
}

func (f *Forecast) HighlyRecommended() bool {
	return f.highlyRecommended
}

func (f *Forecast) FavoriteNum() int {
	return f.favoriteNum
}

func (f *Forecast) RivalNum() int {
	return f.rivalNum
}

func (f *Forecast) MarkerNum() int {
	return f.markerNum
}
