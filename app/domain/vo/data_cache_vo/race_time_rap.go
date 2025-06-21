package data_cache_vo

import "time"

type RaceTimeRap struct {
	rapTimes []time.Duration
	first3f  time.Duration
	first4f  time.Duration
	last3f   time.Duration
	last4f   time.Duration
	rap5f    time.Duration
}

func NewRaceTimeRap(rawRapTimes []string) *RaceTimeRap {
	rapTimes := make([]time.Duration, 0, len(rawRapTimes))
	for _, rawRapTime := range rawRapTimes {
		rapTime, _ := time.ParseDuration(rawRapTime + "s")
		rapTimes = append(rapTimes, rapTime)
	}

	var first3f, first4f, last3f, last4f, rap5f time.Duration
	if rapTimes[0] < 10*time.Second {
		first3f, first4f, rap5f = calcOddDistanceTime(rapTimes)
	} else {
		first3f, first4f, rap5f = calcEvenDistanceTime(rapTimes)
	}

	last3f, last4f = calcLastDistanceTime(rapTimes)

	return &RaceTimeRap{
		rapTimes: rapTimes,
		first3f:  first3f,
		first4f:  first4f,
		last3f:   last3f,
		last4f:   last4f,
		rap5f:    rap5f,
	}
}

func (r *RaceTimeRap) RapTimes() []time.Duration {
	return r.rapTimes
}

func (r *RaceTimeRap) First3f() time.Duration {
	return r.first3f
}

func (r *RaceTimeRap) First4f() time.Duration {
	return r.first4f
}

func (r *RaceTimeRap) Last3f() time.Duration {
	return r.last3f
}

func (r *RaceTimeRap) Last4f() time.Duration {
	return r.last4f
}

func (r *RaceTimeRap) Rap5f() time.Duration {
	return r.rap5f
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
	reversedRapTimes := reverseSlice(rapTimes)

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

func reverseSlice[T any](s []T) []T {
	result := make([]T, len(s))
	for i := range result {
		result[i] = s[len(s)-1-i]
	}
	return result
}
