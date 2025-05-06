package data_cache_entity

import (
	"time"

	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type RaceTime struct {
	raceId     types.RaceId
	raceDate   types.RaceDate
	time       string
	timeIndex  int
	trackIndex int
	rapTimes   []time.Duration
	first3f    time.Duration
	first4f    time.Duration
	last3f     time.Duration
	last4f     time.Duration
	rap5f      time.Duration
}

func NewRaceTime(
	raceId types.RaceId,
	raceDate types.RaceDate,
	time string,
	timeIndex int,
	trackIndex int,
	rapTimes []time.Duration,
	first3f time.Duration,
	first4f time.Duration,
	last3f time.Duration,
	last4f time.Duration,
	rap5f time.Duration,
) *RaceTime {
	return &RaceTime{
		raceId:     raceId,
		raceDate:   raceDate,
		time:       time,
		timeIndex:  timeIndex,
		trackIndex: trackIndex,
		rapTimes:   rapTimes,
		first3f:    first3f,
		first4f:    first4f,
		last3f:     last3f,
		last4f:     last4f,
		rap5f:      rap5f,
	}
}

func (r *RaceTime) RaceId() types.RaceId {
	return r.raceId
}

func (r *RaceTime) RaceDate() types.RaceDate {
	return r.raceDate
}

func (r *RaceTime) Time() string {
	return r.time
}

func (r *RaceTime) TimeIndex() int {
	return r.timeIndex
}

func (r *RaceTime) TrackIndex() int {
	return r.trackIndex
}

func (r *RaceTime) RapTimes() []time.Duration {
	return r.rapTimes
}

func (r *RaceTime) First3f() time.Duration {
	return r.first3f
}

func (r *RaceTime) First4f() time.Duration {
	return r.first4f
}

func (r *RaceTime) Last3f() time.Duration {
	return r.last3f
}

func (r *RaceTime) Last4f() time.Duration {
	return r.last4f
}

func (r *RaceTime) Rap5f() time.Duration {
	return r.rap5f
}
