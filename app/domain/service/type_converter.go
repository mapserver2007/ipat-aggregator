package service

func ConvertToMap[T comparable, V comparable](elms []T, fn func(T) V) map[V]T {
	outputMap := map[V]T{}
	for _, elm := range elms {
		outputMap[fn(elm)] = elm
	}
	return outputMap
}

func ConvertToSliceMap[T comparable, V comparable](elms []T, fn func(T) V) map[V][]T {
	outputSliceMap := map[V][]T{}
	for _, elm := range elms {
		outputSliceMap[fn(elm)] = append(outputSliceMap[fn(elm)], elm)
	}
	return outputSliceMap
}
