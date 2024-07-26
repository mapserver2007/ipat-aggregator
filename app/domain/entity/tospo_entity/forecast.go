package tospo_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type Forecast struct {
	horseNumber types.HorseNumber
	favoriteNum int
	rivalNum    int
	markerNum   int
}

func NewForecast(
	horseNumber types.HorseNumber,
	favoriteNum int,
	rivalNum int,
	markerNum int,
) *Forecast {
	return &Forecast{
		horseNumber: horseNumber,
		favoriteNum: favoriteNum,
		rivalNum:    rivalNum,
		markerNum:   markerNum,
	}
}

func (f *Forecast) HorseNumber() types.HorseNumber {
	return f.horseNumber
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
