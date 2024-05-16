package list_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type Jockey struct {
	jockeyId   types.JockeyId
	jockeyName string
}

func NewJockey(
	jockeyId types.JockeyId,
	jockeyName string,
) *Jockey {
	return &Jockey{
		jockeyId:   jockeyId,
		jockeyName: jockeyName,
	}
}

func (j *Jockey) JockeyId() types.JockeyId {
	return j.jockeyId
}

func (j *Jockey) JockeyName() string {
	return j.jockeyName
}
