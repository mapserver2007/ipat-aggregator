package infrastructure

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"google.golang.org/api/sheets/v4"
	"log"
	"sort"
)

const (
	spreadSheetPredictionFileName = "spreadsheet_prediction.json"
)

type spreadSheetPredictionRepository struct {
	client             *sheets.Service
	spreadSheetConfigs []*spreadsheet_entity.SpreadSheetConfig
	spreadSheetService service.SpreadSheetService
}

func NewSpreadSheetPredictionRepository(
	spreadSheetService service.SpreadSheetService,
) (repository.SpreadSheetPredictionRepository, error) {
	ctx := context.Background()
	client, spreadSheetConfigs, err := getSpreadSheetConfigs(ctx, spreadSheetMarkerAnalysisFileName)
	if err != nil {
		return nil, err
	}

	return &spreadSheetPredictionRepository{
		client:             client,
		spreadSheetConfigs: spreadSheetConfigs,
		spreadSheetService: spreadSheetService,
	}, nil
}

func (s *spreadSheetPredictionRepository) Write(
	ctx context.Context,
	strictPredictionData *spreadsheet_entity.PredictionData,
	simplePredictionData *spreadsheet_entity.PredictionData,
	markerOddsRangeMap map[types.Marker]types.OddsRangeType,
	race *prediction_entity.Race,
) error {
	log.Println(ctx, fmt.Sprintf("write prediction %v start", race.RaceId()))

	strictOddsRangeRaceCountMap := strictPredictionData.OddsRangeRaceCountMap()
	strictPredictionMarkerCombinationData := strictPredictionData.PredictionMarkerCombinationData()
	strictValuesList, err := s.createOddsRangeValues(ctx, strictOddsRangeRaceCountMap, strictPredictionMarkerCombinationData)
	if err != nil {
		return err
	}

	simpleOddsRangeRaceCountMap := simplePredictionData.OddsRangeRaceCountMap()
	simplePredictionMarkerCombinationData := simplePredictionData.PredictionMarkerCombinationData()
	simpleValuesList, err := s.createOddsRangeValues(ctx, simpleOddsRangeRaceCountMap, simplePredictionMarkerCombinationData)
	if err != nil {
		return err
	}

	fmt.Println(strictValuesList)
	fmt.Println(simpleValuesList)

	log.Println(ctx, fmt.Sprintf("write prediction %v end", race.RaceId()))

	return nil
}

func (s *spreadSheetPredictionRepository) createOddsRangeValues(
	ctx context.Context,
	markerCombinationOddsRangeRaceCountMap map[types.MarkerCombinationId]map[types.OddsRangeType]int,
	predictionMarkerCombinationData map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis,
) ([][][]interface{}, error) {
	valuesList := make([][][]interface{}, 0)
	valuesList = append(valuesList, [][]interface{}{
		{
			"",
			"1着率",
			types.WinOddsRange1.String(),
			types.WinOddsRange2.String(),
			types.WinOddsRange3.String(),
			types.WinOddsRange4.String(),
			types.WinOddsRange5.String(),
			types.WinOddsRange6.String(),
			types.WinOddsRange7.String(),
			types.WinOddsRange8.String(),
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			"",
			"2着率",
			types.WinOddsRange1.String(),
			types.WinOddsRange2.String(),
			types.WinOddsRange3.String(),
			types.WinOddsRange4.String(),
			types.WinOddsRange5.String(),
			types.WinOddsRange6.String(),
			types.WinOddsRange7.String(),
			types.WinOddsRange8.String(),
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			"",
			"1着率",
			types.WinOddsRange1.String(),
			types.WinOddsRange2.String(),
			types.WinOddsRange3.String(),
			types.WinOddsRange4.String(),
			types.WinOddsRange5.String(),
			types.WinOddsRange6.String(),
			types.WinOddsRange7.String(),
			types.WinOddsRange8.String(),
		},
	})

	rateFormatFunc := func(matchCount int, raceCount int) string {
		if raceCount == 0 {
			return "-"
		}
		return fmt.Sprintf("%.2f%%", float64(matchCount)*100/float64(raceCount))
	}

	oddsRanges := []types.OddsRangeType{
		types.WinOddsRange1,
		types.WinOddsRange2,
		types.WinOddsRange3,
		types.WinOddsRange4,
		types.WinOddsRange5,
		types.WinOddsRange6,
		types.WinOddsRange7,
		types.WinOddsRange8,
	}

	rawMarkerCombinationIds := make([]int, 0, len(predictionMarkerCombinationData))
	for markerCombinationId := range predictionMarkerCombinationData {
		rawMarkerCombinationIds = append(rawMarkerCombinationIds, markerCombinationId.Value())
	}
	sort.Ints(rawMarkerCombinationIds)

	for _, rawMarkerCombinationId := range rawMarkerCombinationIds {
		markerCombinationId := types.MarkerCombinationId(rawMarkerCombinationId)
		markerCombinationAnalysisData := predictionMarkerCombinationData[markerCombinationId]
		switch markerCombinationId.TicketType() {
		case types.Win:
			marker, err := types.NewMarker(markerCombinationId.Value() % 10)
			if err != nil {
				return nil, err
			}

			oddsRangeMap := s.createOddsRangeMap(ctx, markerCombinationAnalysisData, 1)
			oddsRangeRaceCountMap := markerCombinationOddsRangeRaceCountMap[markerCombinationId]

			raceCount := 0
			for _, oddsRange := range oddsRanges {
				if n, ok := oddsRangeRaceCountMap[oddsRange]; ok {
					raceCount += n
				}
			}

			matchCount := 0
			for _, calculable := range markerCombinationAnalysisData.Calculables() {
				if calculable.OrderNo() == 1 {
					matchCount++
				}
			}

			valuesList[0] = append(valuesList[0], [][]interface{}{
				{
					marker.String(),
					rateFormatFunc(matchCount, raceCount),
					rateFormatFunc(oddsRangeMap[types.WinOddsRange1], oddsRangeRaceCountMap[types.WinOddsRange1]),
					rateFormatFunc(oddsRangeMap[types.WinOddsRange2], oddsRangeRaceCountMap[types.WinOddsRange2]),
					rateFormatFunc(oddsRangeMap[types.WinOddsRange3], oddsRangeRaceCountMap[types.WinOddsRange3]),
					rateFormatFunc(oddsRangeMap[types.WinOddsRange4], oddsRangeRaceCountMap[types.WinOddsRange4]),
					rateFormatFunc(oddsRangeMap[types.WinOddsRange5], oddsRangeRaceCountMap[types.WinOddsRange5]),
					rateFormatFunc(oddsRangeMap[types.WinOddsRange6], oddsRangeRaceCountMap[types.WinOddsRange6]),
					rateFormatFunc(oddsRangeMap[types.WinOddsRange7], oddsRangeRaceCountMap[types.WinOddsRange7]),
					rateFormatFunc(oddsRangeMap[types.WinOddsRange8], oddsRangeRaceCountMap[types.WinOddsRange8]),
				},
			}...)
		case types.Place:
			marker, err := types.NewMarker(markerCombinationId.Value() % 10)
			if err != nil {
				return nil, err
			}

			inOrder2oddsRangeMap := s.createOddsRangeMap(ctx, markerCombinationAnalysisData, 2)
			inOrder3oddsRangeMap := s.createOddsRangeMap(ctx, markerCombinationAnalysisData, 3)
			oddsRangeRaceCountMap := markerCombinationOddsRangeRaceCountMap[markerCombinationId]

			raceCount := 0
			for _, oddsRange := range oddsRanges {
				if n, ok := oddsRangeRaceCountMap[oddsRange]; ok {
					raceCount += n
				}
			}

			orderNo2MatchCount := 0
			orderNo3MatchCount := 0
			for _, calculable := range markerCombinationAnalysisData.Calculables() {
				if calculable.OrderNo() <= 2 {
					orderNo2MatchCount++
				}
				if calculable.OrderNo() <= 3 {
					orderNo3MatchCount++
				}
			}

			valuesList[1] = append(valuesList[1], [][]interface{}{
				{
					marker.String(),
					rateFormatFunc(orderNo2MatchCount, raceCount),
					rateFormatFunc(inOrder2oddsRangeMap[types.WinOddsRange1], oddsRangeRaceCountMap[types.WinOddsRange1]),
					rateFormatFunc(inOrder2oddsRangeMap[types.WinOddsRange2], oddsRangeRaceCountMap[types.WinOddsRange2]),
					rateFormatFunc(inOrder2oddsRangeMap[types.WinOddsRange3], oddsRangeRaceCountMap[types.WinOddsRange3]),
					rateFormatFunc(inOrder2oddsRangeMap[types.WinOddsRange4], oddsRangeRaceCountMap[types.WinOddsRange4]),
					rateFormatFunc(inOrder2oddsRangeMap[types.WinOddsRange5], oddsRangeRaceCountMap[types.WinOddsRange5]),
					rateFormatFunc(inOrder2oddsRangeMap[types.WinOddsRange6], oddsRangeRaceCountMap[types.WinOddsRange6]),
					rateFormatFunc(inOrder2oddsRangeMap[types.WinOddsRange7], oddsRangeRaceCountMap[types.WinOddsRange7]),
					rateFormatFunc(inOrder2oddsRangeMap[types.WinOddsRange8], oddsRangeRaceCountMap[types.WinOddsRange8]),
				},
			}...)
			valuesList[2] = append(valuesList[2], [][]interface{}{
				{
					marker.String(),
					rateFormatFunc(orderNo3MatchCount, raceCount),
					rateFormatFunc(inOrder3oddsRangeMap[types.WinOddsRange1], oddsRangeRaceCountMap[types.WinOddsRange1]),
					rateFormatFunc(inOrder3oddsRangeMap[types.WinOddsRange2], oddsRangeRaceCountMap[types.WinOddsRange2]),
					rateFormatFunc(inOrder3oddsRangeMap[types.WinOddsRange3], oddsRangeRaceCountMap[types.WinOddsRange3]),
					rateFormatFunc(inOrder3oddsRangeMap[types.WinOddsRange4], oddsRangeRaceCountMap[types.WinOddsRange4]),
					rateFormatFunc(inOrder3oddsRangeMap[types.WinOddsRange5], oddsRangeRaceCountMap[types.WinOddsRange5]),
					rateFormatFunc(inOrder3oddsRangeMap[types.WinOddsRange6], oddsRangeRaceCountMap[types.WinOddsRange6]),
					rateFormatFunc(inOrder3oddsRangeMap[types.WinOddsRange7], oddsRangeRaceCountMap[types.WinOddsRange7]),
					rateFormatFunc(inOrder3oddsRangeMap[types.WinOddsRange8], oddsRangeRaceCountMap[types.WinOddsRange8]),
				},
			}...)
		}
	}

	return valuesList, nil
}

func (s *spreadSheetPredictionRepository) Style(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s *spreadSheetPredictionRepository) Clear(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s *spreadSheetPredictionRepository) createOddsRangeMap(
	ctx context.Context,
	markerCombinationAnalysis *spreadsheet_entity.MarkerCombinationAnalysis,
	inOrderNo int,
) map[types.OddsRangeType]int {
	oddsRangeMap := map[types.OddsRangeType]int{
		types.WinOddsRange1: 0,
		types.WinOddsRange2: 0,
		types.WinOddsRange3: 0,
		types.WinOddsRange4: 0,
		types.WinOddsRange5: 0,
		types.WinOddsRange6: 0,
		types.WinOddsRange7: 0,
		types.WinOddsRange8: 0,
	}

	for _, calculable := range markerCombinationAnalysis.Calculables() {
		if calculable.OrderNo() <= inOrderNo {
			odds := calculable.Odds().InexactFloat64()
			if odds >= 1.0 && odds <= 1.5 {
				oddsRangeMap[types.WinOddsRange1]++
			} else if odds >= 1.6 && odds <= 2.0 {
				oddsRangeMap[types.WinOddsRange2]++
			} else if odds >= 2.1 && odds <= 2.9 {
				oddsRangeMap[types.WinOddsRange3]++
			} else if odds >= 3.0 && odds <= 4.9 {
				oddsRangeMap[types.WinOddsRange4]++
			} else if odds >= 5.0 && odds <= 9.9 {
				oddsRangeMap[types.WinOddsRange5]++
			} else if odds >= 10.0 && odds <= 19.9 {
				oddsRangeMap[types.WinOddsRange6]++
			} else if odds >= 20.0 && odds <= 49.9 {
				oddsRangeMap[types.WinOddsRange7]++
			} else if odds >= 50.0 {
				oddsRangeMap[types.WinOddsRange8]++
			}
		}
	}

	return oddsRangeMap
}

func (s *spreadSheetPredictionRepository) createValueList() {}
