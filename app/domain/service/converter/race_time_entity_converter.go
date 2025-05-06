package converter

import (
	"fmt"
	"time"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type RaceTimeEntityConverter interface {
	DataCacheToRaw(input *data_cache_entity.RaceTime) *raw_entity.RaceTime
	NetKeibaToRaw(input *netkeiba_entity.RaceTime) *raw_entity.RaceTime
	RawToDataCache(input *raw_entity.RaceTime) *data_cache_entity.RaceTime
}

type raceTimeEntityConverter struct{}

func NewRaceTimeEntityConverter() RaceTimeEntityConverter {
	return &raceTimeEntityConverter{}
}

func (r *raceTimeEntityConverter) DataCacheToRaw(input *data_cache_entity.RaceTime) *raw_entity.RaceTime {
	rapTimes := make([]string, 0, len(input.RapTimes()))
	for _, rapTime := range input.RapTimes() {
		rapTimes = append(rapTimes, fmt.Sprintf("%.1f", rapTime.Seconds()))
	}
	return &raw_entity.RaceTime{
		RaceId:     input.RaceId().String(),
		RaceDate:   input.RaceDate().Value(),
		Time:       input.Time(),
		TimeIndex:  input.TimeIndex(),
		TrackIndex: input.TrackIndex(),
		RapTimes:   rapTimes,
		First3f:    fmt.Sprintf("%.1f", input.First3f().Seconds()),
		First4f:    fmt.Sprintf("%.1f", input.First4f().Seconds()),
		Last3f:     fmt.Sprintf("%.1f", input.Last3f().Seconds()),
		Last4f:     fmt.Sprintf("%.1f", input.Last4f().Seconds()),
		Rap5f:      fmt.Sprintf("%.1f", input.Rap5f().Seconds()),
	}
}

func (r *raceTimeEntityConverter) NetKeibaToRaw(input *netkeiba_entity.RaceTime) *raw_entity.RaceTime {
	rawRapTimes := make([]string, 0, len(input.RapTimes()))
	rapTimes := make([]time.Duration, 0, len(input.RapTimes()))

	for _, rapTime := range input.RapTimes() {
		rawRapTimes = append(rawRapTimes, fmt.Sprintf("%.1f", rapTime.Seconds()))
		rapTimes = append(rapTimes, rapTime)
	}

	var first3f, first4f, last3f, last4f, rap5f time.Duration
	if rapTimes[0] < 10*time.Second {
		first3f, first4f, rap5f = calcOddDistanceTime(rapTimes)
	} else {
		first3f, first4f, rap5f = calcEvenDistanceTime(rapTimes)
	}

	last3f, last4f = calcLastDistanceTime(rapTimes)

	return &raw_entity.RaceTime{
		RaceId:     input.RaceId(),
		RaceDate:   input.RaceDate(),
		Time:       input.Time(),
		TimeIndex:  input.TimeIndex(),
		TrackIndex: input.TrackIndex(),
		RapTimes:   rawRapTimes,
		First3f:    fmt.Sprintf("%.1f", first3f.Seconds()),
		First4f:    fmt.Sprintf("%.1f", first4f.Seconds()),
		Last3f:     fmt.Sprintf("%.1f", last3f.Seconds()),
		Last4f:     fmt.Sprintf("%.1f", last4f.Seconds()),
		Rap5f:      fmt.Sprintf("%.1f", rap5f.Seconds()),
	}
}

func calcOddDistanceTime(
	rapTimeValues []time.Duration,
) (time.Duration, time.Duration, time.Duration) {
	var first3f, first4f, rap5f time.Duration
	for i, rapTime := range rapTimeValues {
		if i >= 6 {
			break
		}
		if i < 3 {
			first3f += rapTime
			first4f += rapTime
			rap5f += rapTime
		}
		switch i {
		case 3:
			first3f += rapTime / 2
			first4f += rapTime
			rap5f += rapTime
		case 4:
			first4f += rapTime / 2
			rap5f += rapTime
		case 5:
			rap5f += rapTime / 2
		}
	}

	return first3f, first4f, rap5f
}

func calcEvenDistanceTime(
	rapTimeValues []time.Duration,
) (time.Duration, time.Duration, time.Duration) {
	var first3f, first4f, rap5f time.Duration
	for i, rapTime := range rapTimeValues {
		if i >= 5 {
			break
		}
		if i < 2 {
			first3f += rapTime
			first4f += rapTime
			rap5f += rapTime
		}
		switch i {
		case 2:
			first3f += rapTime
			first4f += rapTime
			rap5f += rapTime
		case 3:
			first4f += rapTime
			rap5f += rapTime
		case 4:
			rap5f += rapTime
		}
	}

	return first3f, first4f, rap5f
}

func calcLastDistanceTime(rapTimes []time.Duration) (time.Duration, time.Duration) {
	var last3f, last4f time.Duration
	reversedRapTimes := ReverseSlice(rapTimes)

	for i, rapTime := range reversedRapTimes {
		if i >= 4 {
			break
		}
		if i < 3 {
			last3f += rapTime
		}
		if i < 4 {
			last4f += rapTime
		}
	}

	return last3f, last4f
}

func (r *raceTimeEntityConverter) RawToDataCache(input *raw_entity.RaceTime) *data_cache_entity.RaceTime {
	rapTimes := make([]time.Duration, 0, len(input.RapTimes))
	for _, rawRapTime := range input.RapTimes {
		rapTime, _ := time.ParseDuration(rawRapTime)
		rapTimes = append(rapTimes, rapTime)
	}

	first3f, _ := time.ParseDuration(input.First3f + "s")
	first4f, _ := time.ParseDuration(input.First4f + "s")
	last3f, _ := time.ParseDuration(input.Last3f + "s")
	last4f, _ := time.ParseDuration(input.Last4f + "s")
	rap5f, _ := time.ParseDuration(input.Rap5f + "s")

	return data_cache_entity.NewRaceTime(
		types.RaceId(input.RaceId),
		types.RaceDate(input.RaceDate),
		input.Time,
		input.TimeIndex,
		input.TrackIndex,
		rapTimes,
		first3f,
		first4f,
		last3f,
		last4f,
		rap5f,
	)
}
