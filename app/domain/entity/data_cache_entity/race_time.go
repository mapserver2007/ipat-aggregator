package data_cache_entity

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/vo/data_cache_vo"
)

type RaceTime struct {
	raceTimeId   types.RaceTime
	raceId       types.RaceId
	raceDate     types.RaceDate
	time         string
	durationTime time.Duration
	timeIndex    int
	trackIndex   int
	rapTimes     []time.Duration
	first3f      time.Duration
	first4f      time.Duration
	last3f       time.Duration
	last4f       time.Duration
	rap5f        time.Duration
}

type RaceTimeV2 struct {
	raceTimeId   types.RaceTime
	raceId       types.RaceId
	raceDate     types.RaceDate
	time         string
	durationTime time.Duration
	timeIndex    int
	trackIndex   int
	rapTimeRap   *data_cache_vo.RaceTimeRap
}

func NewRaceTime(
	raceTimeId types.RaceTime,
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
	durationTime, _ := timeToDuration(time)
	return &RaceTime{
		raceTimeId:   raceTimeId,
		raceId:       raceId,
		raceDate:     raceDate,
		time:         time,
		durationTime: durationTime,
		timeIndex:    timeIndex,
		trackIndex:   trackIndex,
		rapTimes:     rapTimes,
		first3f:      first3f,
		first4f:      first4f,
		last3f:       last3f,
		last4f:       last4f,
		rap5f:        rap5f,
	}
}

func (r *RaceTime) RaceTimeId() types.RaceTime {
	return r.raceTimeId
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

func (r *RaceTime) DurationTime() time.Duration {
	return r.durationTime
}

func (r *RaceTime) DurationTimeFormat() string {
	return formatRaceTime(r.durationTime)
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

func NewRaceTimeV2(
	raceTimeId types.RaceTime,
	raceId types.RaceId,
	raceDate types.RaceDate,
	time string,
	timeIndex int,
	trackIndex int,
	rapTimeRap *data_cache_vo.RaceTimeRap,
) *RaceTimeV2 {
	durationTime, _ := timeToDuration(time)
	return &RaceTimeV2{
		raceTimeId:   raceTimeId,
		raceId:       raceId,
		raceDate:     raceDate,
		time:         time,
		durationTime: durationTime,
		timeIndex:    timeIndex,
		trackIndex:   trackIndex,
		rapTimeRap:   rapTimeRap,
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

func (r *RaceTimeV2) Time() string {
	return r.time
}

func (r *RaceTimeV2) DurationTime() time.Duration {
	return r.durationTime
}

func (r *RaceTimeV2) DurationTimeFormat() string {
	return formatRaceTime(r.durationTime)
}

func (r *RaceTimeV2) TimeIndex() int {
	return r.timeIndex
}

func (r *RaceTimeV2) TrackIndex() int {
	return r.trackIndex
}

func (r *RaceTimeV2) RapTimeRap() *data_cache_vo.RaceTimeRap {
	return r.rapTimeRap
}

func timeToDuration(input string) (time.Duration, error) {
	var minutes, seconds float64

	if strings.Contains(input, ":") {
		parts := strings.Split(input, ":")
		if len(parts) != 2 {
			return 0, fmt.Errorf("invalid time format: %s", input)
		}
		min, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, fmt.Errorf("invalid minutes: %w", err)
		}
		sec, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid seconds: %w", err)
		}
		minutes = float64(min)
		seconds = sec
	} else {
		sec, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid seconds: %w", err)
		}
		seconds = sec
	}

	totalSeconds := minutes*60 + seconds
	return time.Duration(totalSeconds * float64(time.Second)), nil
}

func formatRaceTime(d time.Duration) string {
	minutes := int(d / time.Minute)
	seconds := int((d % time.Minute) / time.Second)
	tenths := int((d % time.Second) / (time.Millisecond * 100))
	if minutes > 0 {
		return fmt.Sprintf("%d分%d秒%d", minutes, seconds, tenths)
	}
	return fmt.Sprintf("%d秒%d", seconds, tenths)
}
