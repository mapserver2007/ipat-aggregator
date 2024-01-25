package predict_analysis_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type IFilter interface {
	FilterId() types.AnalysisFilterId
	Value() any
}

type Filter[T any] struct {
	filterId types.AnalysisFilterId
	value    T
}

func NewFilter[T any](
	filterId types.AnalysisFilterId,
	value T,
) *Filter[T] {
	return &Filter[T]{
		filterId: filterId,
		value:    value,
	}
}

func (f *Filter[T]) FilterId() types.AnalysisFilterId {
	return f.filterId
}

func (f *Filter[T]) Value() any {
	return f.value
}
