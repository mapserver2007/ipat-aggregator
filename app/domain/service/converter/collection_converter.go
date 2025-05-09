package converter

import (
	"sort"

	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

func SortedRaceIdKeys[T any](m map[types.RaceId]T) []types.RaceId {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key.String())
	}
	sort.Strings(keys)

	raceIds := make([]types.RaceId, 0, len(keys))
	for _, key := range keys {
		raceIds = append(raceIds, types.RaceId(key))
	}

	return raceIds
}

func SortedRaceDateKeys[T any](m map[types.RaceDate]T) []types.RaceDate {
	keys := make([]int, 0, len(m))
	for key := range m {
		keys = append(keys, key.Value())
	}
	sort.Ints(keys)

	raceDates := make([]types.RaceDate, 0, len(keys))
	for _, key := range keys {
		raceDates = append(raceDates, types.RaceDate(key))
	}

	return raceDates
}

func ReverseSlice[T any](s []T) []T {
	n := len(s)
	result := make([]T, n)
	for i := 0; i < n; i++ {
		result[i] = s[n-1-i]
	}
	return result
}
