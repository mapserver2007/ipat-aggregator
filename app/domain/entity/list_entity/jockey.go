package list_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type Jockey struct {
	jockeyId types.JockeyId
}

func NewJockey(
	rawJockeyId int,
) *Jockey {
	return &Jockey{
		jockeyId: types.JockeyId(rawJockeyId),
	}
}

func (j *Jockey) JockeyId() types.JockeyId {
	return j.jockeyId
}
