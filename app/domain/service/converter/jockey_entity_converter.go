package converter

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
)

type JockeyEntityConverter interface {
	DataCacheToRaw(input *data_cache_entity.Jockey) *raw_entity.Jockey
	RawToDataCache(input *raw_entity.Jockey) *data_cache_entity.Jockey
	NetKeibaToRaw(input *netkeiba_entity.Jockey) *raw_entity.Jockey
}

type jockeyEntityConverter struct{}

func NewJockeyEntityConverter() JockeyEntityConverter {
	return &jockeyEntityConverter{}
}

func (j *jockeyEntityConverter) DataCacheToRaw(input *data_cache_entity.Jockey) *raw_entity.Jockey {
	return &raw_entity.Jockey{
		JockeyId:   input.JockeyId().Value(),
		JockeyName: input.JockeyName(),
	}
}

func (j *jockeyEntityConverter) RawToDataCache(input *raw_entity.Jockey) *data_cache_entity.Jockey {
	return data_cache_entity.NewJockey(
		input.JockeyId,
		input.JockeyName,
	)
}

func (j *jockeyEntityConverter) NetKeibaToRaw(input *netkeiba_entity.Jockey) *raw_entity.Jockey {
	return &raw_entity.Jockey{
		JockeyId:   input.Id(),
		JockeyName: input.Name(),
	}
}
