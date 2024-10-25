package prediction_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type RaceForecast struct {
	horseNumber         types.HorseNumber
	favoriteNum         int
	rivalNum            int
	markerNum           int
	trainingComment     string
	isHighlyRecommended bool
}

func NewRaceForecast(
	horseNumber types.HorseNumber,
	favoriteNum int,
	rivalNum int,
	markerNum int,
	trainingComment string,
	isHighlyRecommended bool,
) *RaceForecast {
	return &RaceForecast{
		horseNumber:         horseNumber,
		favoriteNum:         favoriteNum,
		rivalNum:            rivalNum,
		markerNum:           markerNum,
		trainingComment:     trainingComment,
		isHighlyRecommended: isHighlyRecommended,
	}
}

func (r *RaceForecast) HorseNumber() types.HorseNumber {
	return r.horseNumber
}

func (r *RaceForecast) FavoriteNum() int {
	return r.favoriteNum
}

func (r *RaceForecast) RivalNum() int {
	return r.rivalNum
}

func (r *RaceForecast) MarkerNum() int {
	return r.markerNum
}

func (r *RaceForecast) TrainingComment() string {
	return r.trainingComment
}

func (r *RaceForecast) IsHighlyRecommended() bool {
	return r.isHighlyRecommended
}
