package converter

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
)

type JockeyResultEntityConverter interface {
	NetKeibaToRaw(input *netkeiba_entity.JockeyResult) *raw_entity.JockeyResult
	RawToDataCache(input *raw_entity.JockeyResult) *data_cache_entity.JockeyResult
}

type jockeyResultEntityConverter struct{}

func NewJockeyResultEntityConverter() JockeyResultEntityConverter {
	return &jockeyResultEntityConverter{}
}

func (j *jockeyResultEntityConverter) NetKeibaToRaw(input *netkeiba_entity.JockeyResult) *raw_entity.JockeyResult {
	return &raw_entity.JockeyResult{
		JockeyResultId:   input.JockeyResultId(),
		JockeyId:         input.JockeyId(),
		RaceId:           input.RaceId(),
		RaceDate:         input.RaceDate(),
		RaceCourseId:     input.RaceCourseId(),
		Odds:             input.Odds(),
		OrderNo:          input.OrderNo(),
		HorseId:          input.HorseId(),
		CourseCategoryId: input.CourseCategoryId(),
		Distance:         input.Distance(),
	}
}

func (j *jockeyResultEntityConverter) RawToDataCache(input *raw_entity.JockeyResult) *data_cache_entity.JockeyResult {
	return data_cache_entity.NewJockeyResult(
		input.JockeyResultId,
		input.JockeyId,
		input.RaceId,
		input.RaceDate,
		input.RaceCourseId,
		input.Odds,
		input.OrderNo,
		input.HorseId,
		input.CourseCategoryId,
		input.Distance,
	)
}
