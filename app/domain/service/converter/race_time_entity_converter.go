package converter

import (
	"fmt"
	"time"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/vo/data_cache_vo"
)

type RaceTimeEntityConverter interface {
	NetKeibaToRawV2(input *netkeiba_entity.RaceTime) *raw_entity.RaceTimeV2
	RawToDataCacheV2(input *raw_entity.RaceTimeV2) *data_cache_entity.RaceTimeV2
}

type raceTimeEntityConverter struct{}

func NewRaceTimeEntityConverter() RaceTimeEntityConverter {
	return &raceTimeEntityConverter{}
}

func (r *raceTimeEntityConverter) NetKeibaToRawV2(input *netkeiba_entity.RaceTime) *raw_entity.RaceTimeV2 {
	rawRapTimes := make([]string, 0, len(input.RapTimes()))
	rapTimes := make([]time.Duration, 0, len(input.RapTimes()))

	for _, rapTime := range input.RapTimes() {
		rawRapTimes = append(rawRapTimes, fmt.Sprintf("%.1f", rapTime.Seconds()))
		rapTimes = append(rapTimes, rapTime)
	}

	return &raw_entity.RaceTimeV2{
		RaceTimeId: input.RaceTimeId(),
		RaceId:     input.RaceId(),
		RaceDate:   input.RaceDate(),
		Time:       input.Time(),
		TimeIndex:  input.TimeIndex(),
		TrackIndex: input.TrackIndex(),
		RapTimes:   rawRapTimes,
	}
}

func (r *raceTimeEntityConverter) RawToDataCacheV2(input *raw_entity.RaceTimeV2) *data_cache_entity.RaceTimeV2 {
	return data_cache_entity.NewRaceTimeV2(
		types.RaceTime(input.RaceTimeId),
		types.RaceId(input.RaceId),
		types.RaceDate(input.RaceDate),
		data_cache_vo.NewRaceTime(input.Time),
		input.TimeIndex,
		input.TrackIndex,
		data_cache_vo.NewRaceTimeRap(input.RapTimes),
	)
}
