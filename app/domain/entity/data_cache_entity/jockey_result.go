package data_cache_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type JockeyResult struct {
	jockeyId types.JockeyId
}

func NewJockeyResult(
	jockeyId string,
) *JockeyResult {
	return &JockeyResult{
		jockeyId: types.JockeyId(jockeyId),
	}
}

func (j *JockeyResult) JockeyId() types.JockeyId {
	return j.jockeyId
}
