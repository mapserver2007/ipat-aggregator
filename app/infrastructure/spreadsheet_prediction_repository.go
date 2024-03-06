package infrastructure

import (
	"context"
	"fmt"
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
	spreadSheetConfig  *spreadsheet_entity.SpreadSheetConfig
	spreadSheetService service.SpreadSheetService
}

func NewSpreadSheetPredictionRepository(
	spreadSheetService service.SpreadSheetService,
) (repository.SpreadSheetPredictionRepository, error) {
	ctx := context.Background()
	client, spreadSheetConfig, err := getSpreadSheetConfig(ctx, spreadSheetPredictionFileName)
	if err != nil {
		return nil, err
	}

	return &spreadSheetPredictionRepository{
		client:             client,
		spreadSheetConfig:  spreadSheetConfig,
		spreadSheetService: spreadSheetService,
	}, nil
}

func (s *spreadSheetPredictionRepository) Write(
	ctx context.Context,
	strictPredictionDataList []*spreadsheet_entity.PredictionData,
	simplePredictionDataList []*spreadsheet_entity.PredictionData,
) error {
	log.Println(ctx, fmt.Sprintf("write prediction start"))
	predictionDataSize := len(strictPredictionDataList)

	for idx := 0; idx < predictionDataSize; idx++ {
		strictPredictionData := strictPredictionDataList[idx]
		simplePredictionData := simplePredictionDataList[idx]

		strictValuesList, err := s.createOddsRangeValues(
			ctx,
			strictPredictionData.OddsRangeRaceCountMap(),
			strictPredictionData.PredictionMarkerCombinationData(),
			strictPredictionData.PredictionTitle(),
			strictPredictionData.RaceUrl(),
		)
		if err != nil {
			return err
		}

		simpleValuesList, err := s.createOddsRangeValues(
			ctx,
			simplePredictionData.OddsRangeRaceCountMap(),
			simplePredictionData.PredictionMarkerCombinationData(),
			simplePredictionData.PredictionTitle(),
			simplePredictionData.RaceUrl(),
		)
		if err != nil {
			return err
		}

		var strictValueList [][]interface{}
		for _, value := range strictValuesList {
			strictValueList = append(strictValueList, value...)
		}

		writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName(), fmt.Sprintf("E%d", 1+(idx*22)))
		_, err = s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
			Values: strictValueList,
		}).ValueInputOption("USER_ENTERED").Do()
		if err != nil {
			return err
		}

		var simpleValueList [][]interface{}
		for _, value := range simpleValuesList {
			simpleValueList = append(simpleValueList, value...)
		}

		writeRange = fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName(), fmt.Sprintf("O%d", 1+(idx*22)))
		_, err = s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
			Values: simpleValueList,
		}).ValueInputOption("USER_ENTERED").Do()
		if err != nil {
			return err
		}
	}

	log.Println(ctx, fmt.Sprintf("write prediction end"))

	return nil
}

func (s *spreadSheetPredictionRepository) createOddsRangeValues(
	ctx context.Context,
	markerCombinationOddsRangeRaceCountMap map[types.MarkerCombinationId]map[types.OddsRangeType]int,
	predictionMarkerCombinationData map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis,
	predictionTitle string,
	raceUrl string,
) ([][][]interface{}, error) {
	valuesList := make([][][]interface{}, 0)
	valuesList = append(valuesList, [][]interface{}{
		{
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
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

	var raceCount int
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

			raceCount = 0
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

			valuesList[1] = append(valuesList[1], [][]interface{}{
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

			valuesList[2] = append(valuesList[2], [][]interface{}{
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
			valuesList[3] = append(valuesList[3], [][]interface{}{
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

	valuesList[0][0][1] = fmt.Sprintf("=HYPERLINK(\"%s\",\"%s(%d)\")", raceUrl, predictionTitle, raceCount)
	return valuesList, nil
}

func (s *spreadSheetPredictionRepository) Style(
	ctx context.Context,
	markerOddsRangeMapList []map[types.Marker]types.OddsRangeType,
) error {
	var requests []*sheets.Request
	for dataIndex, markerOddsRangeMap := range markerOddsRangeMapList {
		sortedMarkers := make([]int, 0, len(markerOddsRangeMap))
		for marker := range markerOddsRangeMap {
			sortedMarkers = append(sortedMarkers, marker.Value())
		}
		sort.Ints(sortedMarkers)

		for filterKindNum := range []int{0, 1} {
			for orderIndex := range []int{0, 1, 2} {
				for rowIndex, rawMarker := range sortedMarkers {
					marker := types.Marker(rawMarker)
					oddsRangeType, ok := markerOddsRangeMap[marker]
					if !ok {
						return fmt.Errorf("marker %s not found in markerOddsRangeMapList", marker.String())
					}
					requests = append(requests, []*sheets.Request{
						{
							RepeatCell: &sheets.RepeatCellRequest{
								Fields: "userEnteredFormat.backgroundColor",
								Range: &sheets.GridRange{
									SheetId:          s.spreadSheetConfig.SheetId(),
									StartColumnIndex: int64(5 + oddsRangeType.Value() + (10 * filterKindNum)),
									StartRowIndex:    int64(2 + rowIndex + (7 * orderIndex) + (22 * dataIndex)),
									EndColumnIndex:   int64(6 + oddsRangeType.Value() + (10 * filterKindNum)),
									EndRowIndex:      int64(3 + rowIndex + (7 * orderIndex) + (22 * dataIndex)),
								},
								Cell: &sheets.CellData{
									UserEnteredFormat: &sheets.CellFormat{
										BackgroundColor: &sheets.Color{
											Red:   1.0,
											Blue:  0.0,
											Green: 1.0,
										},
									},
								},
							},
						},
						{
							RepeatCell: &sheets.RepeatCellRequest{
								Fields: "userEnteredFormat.textFormat.bold",
								Range: &sheets.GridRange{
									SheetId:          s.spreadSheetConfig.SheetId(),
									StartColumnIndex: int64(5 + oddsRangeType.Value() + (10 * filterKindNum)),
									StartRowIndex:    int64(2 + rowIndex + (7 * orderIndex) + (22 * dataIndex)),
									EndColumnIndex:   int64(6 + oddsRangeType.Value() + (10 * filterKindNum)),
									EndRowIndex:      int64(3 + rowIndex + (7 * orderIndex) + (22 * dataIndex)),
								},
								Cell: &sheets.CellData{
									UserEnteredFormat: &sheets.CellFormat{
										TextFormat: &sheets.TextFormat{
											Bold: true,
										},
									},
								},
							},
						},
					}...)
				}
			}
		}
	}

	for dataIndex := 0; dataIndex < len(markerOddsRangeMapList); dataIndex++ {
		for colIndex := range []int{0, 1} {
			requests = append(requests, []*sheets.Request{
				{
					RepeatCell: &sheets.RepeatCellRequest{
						Fields: "userEnteredFormat.backgroundColor",
						Range: &sheets.GridRange{
							SheetId:          s.spreadSheetConfig.SheetId(),
							StartColumnIndex: int64(5 + (10 * colIndex)),
							StartRowIndex:    int64(0 + (22 * dataIndex)),
							EndColumnIndex:   int64(14 + (10 * colIndex)),
							EndRowIndex:      int64(1 + (22 * dataIndex)),
						},
						Cell: &sheets.CellData{
							UserEnteredFormat: &sheets.CellFormat{
								BackgroundColor: &sheets.Color{
									Red:   0.0,
									Blue:  1.0,
									Green: 0.0,
								},
							},
						},
					},
				},
				{
					RepeatCell: &sheets.RepeatCellRequest{
						Fields: "userEnteredFormat.textFormat.bold",
						Range: &sheets.GridRange{
							SheetId:          s.spreadSheetConfig.SheetId(),
							StartColumnIndex: int64(5 + (10 * colIndex)),
							StartRowIndex:    int64(0 + (22 * dataIndex)),
							EndColumnIndex:   int64(14 + (10 * colIndex)),
							EndRowIndex:      int64(1 + (22 * dataIndex)),
						},
						Cell: &sheets.CellData{
							UserEnteredFormat: &sheets.CellFormat{
								TextFormat: &sheets.TextFormat{
									Bold: true,
								},
							},
						},
					},
				},
				{
					RepeatCell: &sheets.RepeatCellRequest{
						Fields: "userEnteredFormat.textFormat.foregroundColor",
						Range: &sheets.GridRange{
							SheetId:          s.spreadSheetConfig.SheetId(),
							StartColumnIndex: int64(5 + (10 * colIndex)),
							StartRowIndex:    int64(0 + (22 * dataIndex)),
							EndColumnIndex:   int64(14 + (10 * colIndex)),
							EndRowIndex:      int64(1 + (22 * dataIndex)),
						},
						Cell: &sheets.CellData{
							UserEnteredFormat: &sheets.CellFormat{
								TextFormat: &sheets.TextFormat{
									ForegroundColor: &sheets.Color{
										Red:   1.0,
										Green: 1.0,
										Blue:  1.0,
									},
								},
							},
						},
					},
				},
			}...)
			for rowIndex := range []int{0, 1, 2} {
				requests = append(requests, []*sheets.Request{
					{
						RepeatCell: &sheets.RepeatCellRequest{
							Fields: "userEnteredFormat.textFormat.bold",
							Range: &sheets.GridRange{
								SheetId:          s.spreadSheetConfig.SheetId(),
								StartColumnIndex: int64(5 + (10 * colIndex)),
								StartRowIndex:    int64(1 + (22 * dataIndex) + (7 * rowIndex)),
								EndColumnIndex:   int64(14 + (10 * colIndex)),
								EndRowIndex:      int64(2 + (22 * dataIndex) + (7 * rowIndex)),
							},
							Cell: &sheets.CellData{
								UserEnteredFormat: &sheets.CellFormat{
									TextFormat: &sheets.TextFormat{
										Bold: true,
									},
								},
							},
						},
					},
					{
						RepeatCell: &sheets.RepeatCellRequest{
							Fields: "userEnteredFormat.backgroundColor",
							Range: &sheets.GridRange{
								SheetId:          s.spreadSheetConfig.SheetId(),
								StartColumnIndex: int64(5 + (10 * colIndex)),
								StartRowIndex:    int64(1 + (22 * dataIndex) + (7 * rowIndex)),
								EndColumnIndex:   int64(6 + (10 * colIndex)),
								EndRowIndex:      int64(2 + (22 * dataIndex) + (7 * rowIndex)),
							},
							Cell: &sheets.CellData{
								UserEnteredFormat: &sheets.CellFormat{
									BackgroundColor: &sheets.Color{
										Red:   1.0,
										Blue:  0.0,
										Green: 1.0,
									},
								},
							},
						},
					},
					{
						RepeatCell: &sheets.RepeatCellRequest{
							Fields: "userEnteredFormat.backgroundColor",
							Range: &sheets.GridRange{
								SheetId:          s.spreadSheetConfig.SheetId(),
								StartColumnIndex: int64(6 + (10 * colIndex)),
								StartRowIndex:    int64(1 + (22 * dataIndex) + (7 * rowIndex)),
								EndColumnIndex:   int64(14 + (10 * colIndex)),
								EndRowIndex:      int64(2 + (22 * dataIndex) + (7 * rowIndex)),
							},
							Cell: &sheets.CellData{
								UserEnteredFormat: &sheets.CellFormat{
									BackgroundColor: &sheets.Color{
										Red:   1.0,
										Blue:  0.0,
										Green: 0.0,
									},
								},
							},
						},
					},
					{
						RepeatCell: &sheets.RepeatCellRequest{
							Fields: "userEnteredFormat.textFormat.foregroundColor",
							Range: &sheets.GridRange{
								SheetId:          s.spreadSheetConfig.SheetId(),
								StartColumnIndex: int64(6 + (10 * colIndex)),
								StartRowIndex:    int64(1 + (22 * dataIndex) + (7 * rowIndex)),
								EndColumnIndex:   int64(14 + (10 * colIndex)),
								EndRowIndex:      int64(2 + (22 * dataIndex) + (7 * rowIndex)),
							},
							Cell: &sheets.CellData{
								UserEnteredFormat: &sheets.CellFormat{
									TextFormat: &sheets.TextFormat{
										ForegroundColor: &sheets.Color{
											Red:   1.0,
											Green: 1.0,
											Blue:  1.0,
										},
									},
								},
							},
						},
					},
				}...)
			}
		}
	}

	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetPredictionRepository) Clear(ctx context.Context) error {
	requests := []*sheets.Request{
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Fields: "*",
				Range: &sheets.GridRange{
					SheetId:          s.spreadSheetConfig.SheetId(),
					StartColumnIndex: 4,
					StartRowIndex:    0,
					EndColumnIndex:   40,
					EndRowIndex:      9999,
				},
				Cell: &sheets.CellData{},
			},
		},
	}
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
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
