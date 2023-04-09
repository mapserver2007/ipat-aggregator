package entity

import vo "github.com/mapserver2007/ipat-aggregator/app/domain/jockey/value_object"

type Jockey struct {
	jockeyId   int
	jockeyName string
}

func NewJockey(
	jockeyId int,
	jockeyName string,
) *Jockey {
	return &Jockey{
		jockeyId:   jockeyId,
		jockeyName: jockeyName,
	}
}

func (j *Jockey) JockeyId() vo.JockeyId {
	return vo.JockeyId(j.jockeyId)
}

func (j *Jockey) JockeyName() string {
	return j.jockeyName
}
