package analysis_entity

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type RaceTimeCalculable struct {
	raceId             types.RaceId
	raceDate           types.RaceDate
	time               time.Duration
	timeIndex          int
	trackIndex         int
	rapTimes           []time.Duration
	first3f            time.Duration
	first4f            time.Duration
	last3f             time.Duration
	last4f             time.Duration
	rap5f              time.Duration
	attributeFilterIds []filter.AttributeId
}

func NewRaceTimeCalculable(
	raceId types.RaceId,
	raceDate types.RaceDate,
	rawTime string,
	timeIndex int,
	trackIndex int,
	rapTimes []time.Duration,
	first3f time.Duration,
	first4f time.Duration,
	last3f time.Duration,
	last4f time.Duration,
	rap5f time.Duration,
	attributeFilterIds []filter.AttributeId,
) (*RaceTimeCalculable, error) {
	time, err := parseToDuration(rawTime)
	if err != nil {
		return nil, err
	}

	return &RaceTimeCalculable{
		raceId:             raceId,
		raceDate:           raceDate,
		time:               time,
		timeIndex:          timeIndex,
		trackIndex:         trackIndex,
		rapTimes:           rapTimes,
		first3f:            first3f,
		first4f:            first4f,
		last3f:             last3f,
		last4f:             last4f,
		rap5f:              rap5f,
		attributeFilterIds: attributeFilterIds,
	}, nil
}

func parseToDuration(s string) (time.Duration, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid format: %s", s)
	}

	minutes, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes: %v", err)
	}

	seconds, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid seconds: %v", err)
	}

	totalSeconds := float64(minutes)*60 + seconds
	duration := time.Duration(totalSeconds * float64(time.Second))

	return duration, nil
}

func (r *RaceTimeCalculable) RaceId() types.RaceId {
	return r.raceId
}

func (r *RaceTimeCalculable) RaceDate() types.RaceDate {
	return r.raceDate
}

func (r *RaceTimeCalculable) Time() time.Duration {
	return r.time
}

func (r *RaceTimeCalculable) TimeIndex() int {
	return r.timeIndex
}

func (r *RaceTimeCalculable) TrackIndex() int {
	return r.trackIndex
}

func (r *RaceTimeCalculable) RapTimes() []time.Duration {
	return r.rapTimes
}

func (r *RaceTimeCalculable) First3f() time.Duration {
	return r.first3f
}

func (r *RaceTimeCalculable) First4f() time.Duration {
	return r.first4f
}

func (r *RaceTimeCalculable) Last3f() time.Duration {
	return r.last3f
}

func (r *RaceTimeCalculable) Last4f() time.Duration {
	return r.last4f
}

func (r *RaceTimeCalculable) Rap5f() time.Duration {
	return r.rap5f
}

func (r *RaceTimeCalculable) AttributeFilterIds() []filter.AttributeId {
	return r.attributeFilterIds
}
