package data_cache_vo

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type RaceTime struct {
	time string
}

func NewRaceTime(time string) *RaceTime {
	return &RaceTime{
		time: time,
	}
}

func (r *RaceTime) String() string {
	return r.time
}

func (r *RaceTime) Duration() time.Duration {
	durationTime, _ := timeToDuration(r.time)
	return durationTime
}

func (r *RaceTime) DurationFormat() string {
	return formatRaceTime(r.Duration())
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
