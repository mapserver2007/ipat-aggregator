package analysis_entity

type Result struct {
	calculable *Calculable
	filters    []IFilter
	hit        bool
}

func NewResult(
	calculable *Calculable,
	filters []IFilter,
	hit bool,
) *Result {
	return &Result{
		calculable: calculable,
		filters:    filters,
		hit:        hit,
	}
}

func (r *Result) Calculable() *Calculable {
	return r.calculable
}
func (r *Result) Filters() []IFilter {
	return r.filters
}

func (r *Result) Hit() bool {
	return r.hit
}
