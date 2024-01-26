package analysis_entity

type Result struct {
	calculable *Calculable
	filters    []IFilter
}

func NewResult(
	calculable *Calculable,
	filters []IFilter,
) *Result {
	return &Result{
		calculable: calculable,
		filters:    filters,
	}
}

func (r *Result) Calculable() *Calculable {
	return r.calculable
}
func (r *Result) Filters() []IFilter {
	return r.filters
}
