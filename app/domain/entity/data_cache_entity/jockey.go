package data_cache_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type Jockey struct {
	jockeyId   types.JockeyId
	jockeyName string
}

func NewJockey(
	jockeyId int,
	jockeyName string,
) *Jockey {
	return &Jockey{
		jockeyId:   types.JockeyId(jockeyId),
		jockeyName: jockeyName,
	}
}

func (j *Jockey) JockeyId() types.JockeyId {
	return j.jockeyId
}

func (j *Jockey) JockeyName() string {
	return j.jockeyName
}
