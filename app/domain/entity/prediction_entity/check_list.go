package prediction_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
)

type CheckList struct {
	raceId            types.RaceId
	raceDate          types.RaceDate
	raceName          string
	horseId           types.HorseId
	horseName         string
	winOdds           decimal.Decimal
	marker            types.Marker
	hitRate           float64
	checkList         []bool
	favoriteNum       int
	rivalNum          int
	highlyRecommended bool
	trainingComment   string
}

func NewCheckList(
	raceId types.RaceId,
	raceDate types.RaceDate,
	raceName string,
	horseId types.HorseId,
	horseName string,
	winOdds decimal.Decimal,
	marker types.Marker,
	hitRate float64,
	checkList []bool,
	favoriteNum int,
	rivalNum int,
	highlyRecommended bool,
	trainingComment string,
) (*CheckList, error) {
	return &CheckList{
		raceId:            raceId,
		raceDate:          raceDate,
		raceName:          raceName,
		horseId:           horseId,
		horseName:         horseName,
		winOdds:           winOdds,
		checkList:         checkList,
		favoriteNum:       favoriteNum,
		rivalNum:          rivalNum,
		highlyRecommended: highlyRecommended,
		trainingComment:   trainingComment,
	}, nil
}

func (c *CheckList) RaceId() types.RaceId {
	return c.raceId
}

func (c *CheckList) RaceDate() types.RaceDate {
	return c.raceDate
}

func (c *CheckList) RaceName() string {
	return c.raceName
}

func (c *CheckList) HorseId() types.HorseId {
	return c.horseId
}

func (c *CheckList) HorseName() string {
	return c.horseName
}

func (c *CheckList) WinOdds() decimal.Decimal {
	return c.winOdds
}

func (c *CheckList) Marker() types.Marker {
	return c.marker
}

func (c *CheckList) HitRate() float64 {
	return c.hitRate
}

func (c *CheckList) CheckList() []bool {
	return c.checkList
}

func (c *CheckList) FavoriteNum() int {
	return c.favoriteNum
}

func (c *CheckList) RivalNum() int {
	return c.rivalNum
}

func (c *CheckList) HighlyRecommended() bool {
	return c.highlyRecommended
}

func (c *CheckList) TrainingComment() string {
	return c.trainingComment
}
