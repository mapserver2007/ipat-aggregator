package gateway

import "sort"

func SortedIntKeys[T any](m map[int]T) []int {
	keys := make([]int, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	values := make([]int, 0, len(keys))
	for _, key := range keys {
		values = append(values, key)
	}

	return values
}
