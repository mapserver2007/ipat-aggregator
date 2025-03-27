package data_cache_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
)

type RaceTime struct {
	raceId     types.RaceId
	raceDate   types.RaceDate
	time       string
	timeIndex  int
	trackIndex int
	rapTimes   []decimal.Decimal
	first3f    decimal.Decimal
	first4f    decimal.Decimal
	last3f     decimal.Decimal
	last4f     decimal.Decimal
	rap5f      decimal.Decimal
}

func NewRaceTime(
	raceId types.RaceId,
	raceDate types.RaceDate,
	time string,
	timeIndex int,
	trackIndex int,
	rapTimes []decimal.Decimal,
	first3f decimal.Decimal,
	first4f decimal.Decimal,
	last3f decimal.Decimal,
	last4f decimal.Decimal,
	rap5f decimal.Decimal,
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

func (r *RaceTime) RapTimes() []decimal.Decimal {
	return r.rapTimes
}

func (r *RaceTime) First3f() decimal.Decimal {
	return r.first3f
}

func (r *RaceTime) First4f() decimal.Decimal {
	return r.first4f
}

func (r *RaceTime) Last3f() decimal.Decimal {
	return r.last3f
}

func (r *RaceTime) Last4f() decimal.Decimal {
	return r.last4f
}

func (r *RaceTime) Rap5f() decimal.Decimal {
	return r.rap5f
}
