package entity

import vo "github.com/mapserver2007/ipat-aggregator/app/domain/jockey/value_object"

type Jockey struct {
	jockeyId   vo.JockeyId
	jockeyName string
}

func NewJockey(
	rawJockeyId int,
	jockeyName string,
) *Jockey {
	return &Jockey{
		jockeyId:   vo.JockeyId(rawJockeyId),
		jockeyName: jockeyName,
	}
}

func (j *Jockey) JockeyId() vo.JockeyId {
	return j.jockeyId
}

func (j *Jockey) JockeyName() string {
	return j.jockeyName
}
