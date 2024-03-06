package spreadsheet_entity

import (
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type PredictionData struct {
	predictionMarkerCombinationData map[types.MarkerCombinationId]*MarkerCombinationAnalysis
	oddsRangeRaceCountMap           map[types.MarkerCombinationId]map[types.OddsRangeType]int
	predictionTitle                 string
	raceUrl                         string
}

func NewPredictionData(
	predictionMarkerCombinationData map[types.MarkerCombinationId]*MarkerCombinationAnalysis,
	oddsRangeRaceCountMap map[types.MarkerCombinationId]map[types.OddsRangeType]int,
	race *prediction_entity.Race,
	filter filter.Id,
) *PredictionData {
	predictionTitle := fmt.Sprintf("%s%dR %s %s", race.RaceCourseId().Name(), race.RaceNumber(), race.RaceName(), filter.String())
	return &PredictionData{
		predictionMarkerCombinationData: predictionMarkerCombinationData,
		oddsRangeRaceCountMap:           oddsRangeRaceCountMap,
		predictionTitle:                 predictionTitle,
		raceUrl:                         race.Url(),
	}
}

func (p *PredictionData) PredictionMarkerCombinationData() map[types.MarkerCombinationId]*MarkerCombinationAnalysis {
	return p.predictionMarkerCombinationData
}

func (p *PredictionData) OddsRangeRaceCountMap() map[types.MarkerCombinationId]map[types.OddsRangeType]int {
	return p.oddsRangeRaceCountMap
}

func (p *PredictionData) PredictionTitle() string {
	return p.predictionTitle
}

func (p *PredictionData) RaceUrl() string {
	return p.raceUrl
}
