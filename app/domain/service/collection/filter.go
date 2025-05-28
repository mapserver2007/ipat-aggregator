package collection

import "golang.org/x/exp/constraints"

func Filter[T any](src []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range src {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func Max[T constraints.Ordered](src []T) (max T) {
	if len(src) == 0 {
		return
	}
	max = src[0]
	for _, v := range src[1:] {
		if v > max {
			max = v
		}
	}
	return
}

func Min[T constraints.Ordered](src []T) (min T) {
	if len(src) == 0 {
		return
	}
	min = src[0]
	for _, v := range src[1:] {
		if v < min {
			min = v
		}
	}
	return
}
