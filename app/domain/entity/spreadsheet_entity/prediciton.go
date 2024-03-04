package spreadsheet_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type PredictionData struct {
	predictionMarkerCombinationData map[types.MarkerCombinationId]*MarkerCombinationAnalysis
	oddsRangeRaceCountMap           map[types.MarkerCombinationId]map[types.OddsRangeType]int
	race                            *prediction_entity.Race
}

func NewPredictionData(
	predictionMarkerCombinationData map[types.MarkerCombinationId]*MarkerCombinationAnalysis,
	oddsRangeRaceCountMap map[types.MarkerCombinationId]map[types.OddsRangeType]int,
	race *prediction_entity.Race,
) *PredictionData {
	return &PredictionData{
		predictionMarkerCombinationData: predictionMarkerCombinationData,
		oddsRangeRaceCountMap:           oddsRangeRaceCountMap,
		race:                            race,
	}
}

func (p *PredictionData) PredictionMarkerCombinationData() map[types.MarkerCombinationId]*MarkerCombinationAnalysis {
	return p.predictionMarkerCombinationData
}

func (p *PredictionData) OddsRangeRaceCountMap() map[types.MarkerCombinationId]map[types.OddsRangeType]int {
	return p.oddsRangeRaceCountMap
}

func (p *PredictionData) Race() *prediction_entity.Race {
	return p.race
}
