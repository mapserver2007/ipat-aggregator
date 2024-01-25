package predict_analysis_entity

type Result struct {
	numerical *Numerical
	filters   []IFilter
}

func NewResult(
	numerical *Numerical,
	filters []IFilter,
) *Result {
	return &Result{
		numerical: numerical,
		filters:   filters,
	}
}

func (r *Result) Numerical() *Numerical {
	return r.numerical
}
func (r *Result) Filters() []IFilter {
	return r.filters
}
