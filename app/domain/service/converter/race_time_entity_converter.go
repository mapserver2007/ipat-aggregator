package converter

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
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
		rapTimes = append(rapTimes, rapTime.String())
	}
	return &raw_entity.RaceTime{
		RaceId:     input.RaceId().String(),
		RaceDate:   input.RaceDate().Value(),
		Time:       input.Time(),
		TimeIndex:  input.TimeIndex(),
		TrackIndex: input.TrackIndex(),
		RapTimes:   rapTimes,
		First3f:    input.First3f().String(),
		First4f:    input.First4f().String(),
		Last3f:     input.Last3f().String(),
		Last4f:     input.Last4f().String(),
		Rap5f:      input.Rap5f().String(),
	}
}

func (r *raceTimeEntityConverter) NetKeibaToRaw(input *netkeiba_entity.RaceTime) *raw_entity.RaceTime {
	rawRapTimes := make([]string, 0, len(input.RapTimes()))
	rapTimes := make([]decimal.Decimal, 0, len(input.RapTimes()))

	for _, rapTime := range input.RapTimes() {
		rawRapTimes = append(rawRapTimes, rapTime.String())
		rapTimes = append(rapTimes, rapTime)
	}

	var first3f, first4f, last3f, last4f, rap5f decimal.Decimal
	if rapTimes[0].LessThan(decimal.NewFromFloat(10)) {
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
		First3f:    first3f.String(),
		First4f:    first4f.String(),
		Last3f:     last3f.String(),
		Last4f:     last4f.String(),
		Rap5f:      rap5f.String(),
	}
}

func calcOddDistanceTime(
	rapTimeValues []decimal.Decimal,
) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	var first3f, first4f, rap5f decimal.Decimal
	for i, rapTime := range rapTimeValues {
		if i >= 6 {
			break
		}
		if i < 3 {
			first3f = first3f.Add(rapTime)
			first4f = first4f.Add(rapTime)
			rap5f = rap5f.Add(rapTime)
		}
		switch i {
		case 3:
			first3f = first3f.Add(rapTime.Mul(decimal.NewFromFloat(0.5)))
			first4f = first4f.Add(rapTime)
			rap5f = rap5f.Add(rapTime)
		case 4:
			first4f = first4f.Add(rapTime.Mul(decimal.NewFromFloat(0.5)))
			rap5f = rap5f.Add(rapTime)
		case 5:
			rap5f = rap5f.Add(rapTime.Mul(decimal.NewFromFloat(0.5)))
		}
	}

	return first3f, first4f, rap5f
}

func calcEvenDistanceTime(
	rapTimeValues []decimal.Decimal,
) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	var first3f, first4f, rap5f decimal.Decimal
	for i, rapTime := range rapTimeValues {
		if i >= 5 {
			break
		}
		if i < 2 {
			first3f = first3f.Add(rapTime)
			first4f = first4f.Add(rapTime)
			rap5f = rap5f.Add(rapTime)
		}
		switch i {
		case 2:
			first3f = first3f.Add(rapTime)
			first4f = first4f.Add(rapTime)
			rap5f = rap5f.Add(rapTime)
		case 3:
			first4f = first4f.Add(rapTime)
			rap5f = rap5f.Add(rapTime)
		case 4:
			rap5f = rap5f.Add(rapTime)
		}
	}

	return first3f, first4f, rap5f
}

func calcLastDistanceTime(rapTimes []decimal.Decimal) (decimal.Decimal, decimal.Decimal) {
	var last3f, last4f decimal.Decimal
	reversedRapTimes := ReverseSlice(rapTimes)

	for i, rapTime := range reversedRapTimes {
		if i >= 4 {
			break
		}
		if i < 3 {
			last3f = last3f.Add(rapTime)
		}
		if i < 4 {
			last4f = last4f.Add(rapTime)
		}
	}

	return last3f, last4f
}

func (r *raceTimeEntityConverter) RawToDataCache(input *raw_entity.RaceTime) *data_cache_entity.RaceTime {
	rapTimes := make([]decimal.Decimal, 0, len(input.RapTimes))
	for _, rawRapTime := range input.RapTimes {
		rapTime, _ := decimal.NewFromString(rawRapTime)
		rapTimes = append(rapTimes, rapTime)
	}

	first3f, _ := decimal.NewFromString(input.First3f)
	first4f, _ := decimal.NewFromString(input.First4f)
	last3f, _ := decimal.NewFromString(input.Last3f)
	last4f, _ := decimal.NewFromString(input.Last4f)
	rap5f, _ := decimal.NewFromString(input.Rap5f)

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
