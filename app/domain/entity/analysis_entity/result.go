package analysis_entity

type Result struct {
	calculable *Calculable
	hit        bool
}

func NewResult(
	calculable *Calculable,
	hit bool,
) *Result {
	return &Result{
		calculable: calculable,
		hit:        hit,
	}
}

func (r *Result) Calculable() *Calculable {
	return r.calculable
}

func (r *Result) Hit() bool {
	return r.hit
}
