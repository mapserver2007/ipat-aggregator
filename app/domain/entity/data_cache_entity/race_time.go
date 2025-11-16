package data_cache_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/vo/data_cache_vo"
)

type RaceTimeV2 struct {
	raceTimeId  types.RaceTime
	raceId      types.RaceId
	raceDate    types.RaceDate
	time        *data_cache_vo.RaceTime
	timeIndex   int
	trackIndex  int
	raceTimeRap *data_cache_vo.RaceTimeRap
}

func NewRaceTimeV2(
	raceTimeId types.RaceTime,
	raceId types.RaceId,
	raceDate types.RaceDate,
	time *data_cache_vo.RaceTime,
	timeIndex int,
	trackIndex int,
	raceTimeRap *data_cache_vo.RaceTimeRap,
) *RaceTimeV2 {
	return &RaceTimeV2{
		raceTimeId:  raceTimeId,
		raceId:      raceId,
		raceDate:    raceDate,
		time:        time,
		timeIndex:   timeIndex,
		trackIndex:  trackIndex,
		raceTimeRap: raceTimeRap,
	}
}

func (r *RaceTimeV2) RaceTimeId() types.RaceTime {
	return r.raceTimeId
}

func (r *RaceTimeV2) RaceId() types.RaceId {
	return r.raceId
}

func (r *RaceTimeV2) RaceDate() types.RaceDate {
	return r.raceDate
}

func (r *RaceTimeV2) Time() *data_cache_vo.RaceTime {
	return r.time
}

func (r *RaceTimeV2) TimeIndex() int {
	return r.timeIndex
}

func (r *RaceTimeV2) TrackIndex() int {
	return r.trackIndex
}

func (r *RaceTimeV2) RaceTimeRap() *data_cache_vo.RaceTimeRap {
	return r.raceTimeRap
}
