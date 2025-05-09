package netkeiba_entity

import (
	"time"
)

type RaceTime struct {
	raceId     string
	raceDate   int
	time       string
	timeIndex  int
	trackIndex int
	rapTimes   []time.Duration
}

func NewRaceTime(
	raceId string,
	raceDate int,
	time string,
	timeIndex int,
	trackIndex int,
	rapTimes []time.Duration,
) *RaceTime {
	return &RaceTime{
		raceId:     raceId,
		raceDate:   raceDate,
		time:       time,
		timeIndex:  timeIndex,
		trackIndex: trackIndex,
		rapTimes:   rapTimes,
	}
}

func (r *RaceTime) RaceId() string {
	return r.raceId
}

func (r *RaceTime) RaceDate() int {
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
