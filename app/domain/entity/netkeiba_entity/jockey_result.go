package netkeiba_entity

type JockeyResult struct {
	id string
}

func NewJockeyResult(
	id string,
) *JockeyResult {
	return &JockeyResult{
		id: id,
	}
}

func (j *JockeyResult) Id() string {
	return j.id
}
