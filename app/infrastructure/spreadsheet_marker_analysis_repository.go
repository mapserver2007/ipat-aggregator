package infrastructure

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"google.golang.org/api/sheets/v4"
	"log"
)

const (
	spreadSheetMarkerAnalysisFileName = "spreadsheet_marker_analysis.json"
)

type spreadSheetMarkerAnalysisRepository struct {
	client             *sheets.Service
	spreadSheetConfigs []*spreadsheet_entity.SpreadSheetConfig
}

func NewSpreadSheetMarkerAnalysisRepository() (repository.SpreadSheetMarkerAnalysisRepository, error) {
	ctx := context.Background()
	client, spreadSheetConfigs, err := getSpreadSheetConfigs(ctx, spreadSheetMarkerAnalysisFileName)
	if err != nil {
		return nil, err
	}

	return &spreadSheetMarkerAnalysisRepository{
		client:             client,
		spreadSheetConfigs: spreadSheetConfigs,
	}, nil
}

func (s *spreadSheetMarkerAnalysisRepository) Write(
	ctx context.Context,
	analysisData *spreadsheet_entity.AnalysisData,
	filters []filter.Id,
) error {
	for _, spreadSheetConfig := range s.spreadSheetConfigs {
		var sheetMarker types.Marker
		switch spreadSheetConfig.SheetName() {
		case types.Favorite.String():
			sheetMarker = types.Favorite
		case types.Rival.String():
			sheetMarker = types.Rival
		case types.BrackTriangle.String():
			sheetMarker = types.BrackTriangle
		case types.WhiteTriangle.String():
			sheetMarker = types.WhiteTriangle
		case types.Star.String():
			sheetMarker = types.Star
		case types.Check.String():
			sheetMarker = types.Check
		default:
			return fmt.Errorf("invalid sheet name: %s", spreadSheetConfig.SheetName())
		}

		log.Println(ctx, fmt.Sprintf("write marker %s analysis start", sheetMarker.String()))
		var valuesList [3][][]interface{}
		valuesList[0] = [][]interface{}{
			{
				"",
				"レース数",
				"1着率",
				types.WinOddsRange1.String(),
				types.WinOddsRange2.String(),
				types.WinOddsRange3.String(),
				types.WinOddsRange4.String(),
				types.WinOddsRange5.String(),
				types.WinOddsRange6.String(),
				types.WinOddsRange7.String(),
				types.WinOddsRange8.String(),
				"2着以内率",
				types.WinOddsRange1.String(),
				types.WinOddsRange2.String(),
				types.WinOddsRange3.String(),
				types.WinOddsRange4.String(),
				types.WinOddsRange5.String(),
				types.WinOddsRange6.String(),
				types.WinOddsRange7.String(),
				types.WinOddsRange8.String(),
				"3着以内率",
				types.WinOddsRange1.String(),
				types.WinOddsRange2.String(),
				types.WinOddsRange3.String(),
				types.WinOddsRange4.String(),
				types.WinOddsRange5.String(),
				types.WinOddsRange6.String(),
				types.WinOddsRange7.String(),
				types.WinOddsRange8.String(),
			},
		}
		valuesList[1] = [][]interface{}{
			{
				"",
				"レース数",
				"1着数",
				types.WinOddsRange1.String(),
				types.WinOddsRange2.String(),
				types.WinOddsRange3.String(),
				types.WinOddsRange4.String(),
				types.WinOddsRange5.String(),
				types.WinOddsRange6.String(),
				types.WinOddsRange7.String(),
				types.WinOddsRange8.String(),
				"2着以内率",
				types.WinOddsRange1.String(),
				types.WinOddsRange2.String(),
				types.WinOddsRange3.String(),
				types.WinOddsRange4.String(),
				types.WinOddsRange5.String(),
				types.WinOddsRange6.String(),
				types.WinOddsRange7.String(),
				types.WinOddsRange8.String(),
				"3着以内率",
				types.WinOddsRange1.String(),
				types.WinOddsRange2.String(),
				types.WinOddsRange3.String(),
				types.WinOddsRange4.String(),
				types.WinOddsRange5.String(),
				types.WinOddsRange6.String(),
				types.WinOddsRange7.String(),
				types.WinOddsRange8.String(),
			},
		}
		valuesList[2] = [][]interface{}{
			{
				"",
				"レース数",
				"2着以下数",
				types.WinOddsRange1.String(),
				types.WinOddsRange2.String(),
				types.WinOddsRange3.String(),
				types.WinOddsRange4.String(),
				types.WinOddsRange5.String(),
				types.WinOddsRange6.String(),
				types.WinOddsRange7.String(),
				types.WinOddsRange8.String(),
				"3着以下数",
				types.WinOddsRange1.String(),
				types.WinOddsRange2.String(),
				types.WinOddsRange3.String(),
				types.WinOddsRange4.String(),
				types.WinOddsRange5.String(),
				types.WinOddsRange6.String(),
				types.WinOddsRange7.String(),
				types.WinOddsRange8.String(),
				"4着以下数",
				types.WinOddsRange1.String(),
				types.WinOddsRange2.String(),
				types.WinOddsRange3.String(),
				types.WinOddsRange4.String(),
				types.WinOddsRange5.String(),
				types.WinOddsRange6.String(),
				types.WinOddsRange7.String(),
				types.WinOddsRange8.String(),
			},
		}

		rateFormatFunc := func(matchCount int, raceCount int) string {
			if raceCount == 0 {
				return "-"
			}
			return fmt.Sprintf("%.2f%%", float64(matchCount)*100/float64(raceCount))
		}

		allMarkerCombinationIds := analysisData.AllMarkerCombinationIds()
		markerCombinationMap := analysisData.MarkerCombinationMapByFilter()
		raceCountMap := analysisData.RaceCountMapByFilter()

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

		for idx, f := range filters {
			rowPosition := idx + 1
			for _, markerCombinationId := range allMarkerCombinationIds {
				data, ok := markerCombinationMap[f][markerCombinationId]
				if ok {
					switch markerCombinationId.TicketType() {
					case types.Win:
						marker, err := types.NewMarker(markerCombinationId.Value() % 10)
						if err != nil {
							return err
						}
						if marker != sheetMarker {
							continue
						}

						oddsRangeMap := s.createHitWinOddsRangeMap(ctx, data, 1)
						oddsRangeRaceCountMap := raceCountMap[f][markerCombinationId]
						raceCount := 0
						for _, oddsRange := range oddsRanges {
							if n, ok := oddsRangeRaceCountMap[oddsRange]; ok {
								raceCount += n
							}
						}

						matchCount := 0
						for _, calculable := range data.Calculables() {
							if calculable.OrderNo() == 1 {
								matchCount++
							}
						}

						valuesList[0] = append(valuesList[0], [][]interface{}{
							{
								f.String(),
								raceCount,
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
						valuesList[1] = append(valuesList[1], [][]interface{}{
							{
								f.String(),
								raceCount,
								matchCount,
								oddsRangeMap[types.WinOddsRange1],
								oddsRangeMap[types.WinOddsRange2],
								oddsRangeMap[types.WinOddsRange3],
								oddsRangeMap[types.WinOddsRange4],
								oddsRangeMap[types.WinOddsRange5],
								oddsRangeMap[types.WinOddsRange6],
								oddsRangeMap[types.WinOddsRange7],
								oddsRangeMap[types.WinOddsRange8],
							},
						}...)
					case types.Place:
						marker, err := types.NewMarker(markerCombinationId.Value() % 10)
						if err != nil {
							return err
						}
						if marker != sheetMarker {
							continue
						}

						inOrder2oddsRangeMap := s.createHitWinOddsRangeMap(ctx, data, 2)
						inOrder3oddsRangeMap := s.createHitWinOddsRangeMap(ctx, data, 3)
						oddsRangeRaceCountMap := raceCountMap[f][markerCombinationId]
						raceCount := 0
						for _, oddsRange := range oddsRanges {
							if n, ok := oddsRangeRaceCountMap[oddsRange]; ok {
								raceCount += n
							}
						}

						orderNo2MatchCount := 0
						orderNo3MatchCount := 0
						for _, calculable := range data.Calculables() {
							if calculable.OrderNo() <= 2 {
								orderNo2MatchCount++
							}
							if calculable.OrderNo() <= 3 {
								orderNo3MatchCount++
							}
						}

						valuesList[0][rowPosition] = append(valuesList[0][rowPosition], []interface{}{
							rateFormatFunc(orderNo2MatchCount, raceCount),
							rateFormatFunc(inOrder2oddsRangeMap[types.WinOddsRange1], oddsRangeRaceCountMap[types.WinOddsRange1]),
							rateFormatFunc(inOrder2oddsRangeMap[types.WinOddsRange2], oddsRangeRaceCountMap[types.WinOddsRange2]),
							rateFormatFunc(inOrder2oddsRangeMap[types.WinOddsRange3], oddsRangeRaceCountMap[types.WinOddsRange3]),
							rateFormatFunc(inOrder2oddsRangeMap[types.WinOddsRange4], oddsRangeRaceCountMap[types.WinOddsRange4]),
							rateFormatFunc(inOrder2oddsRangeMap[types.WinOddsRange5], oddsRangeRaceCountMap[types.WinOddsRange5]),
							rateFormatFunc(inOrder2oddsRangeMap[types.WinOddsRange6], oddsRangeRaceCountMap[types.WinOddsRange6]),
							rateFormatFunc(inOrder2oddsRangeMap[types.WinOddsRange7], oddsRangeRaceCountMap[types.WinOddsRange7]),
							rateFormatFunc(inOrder2oddsRangeMap[types.WinOddsRange8], oddsRangeRaceCountMap[types.WinOddsRange8]),
							rateFormatFunc(orderNo3MatchCount, raceCount),
							rateFormatFunc(inOrder3oddsRangeMap[types.WinOddsRange1], oddsRangeRaceCountMap[types.WinOddsRange1]),
							rateFormatFunc(inOrder3oddsRangeMap[types.WinOddsRange2], oddsRangeRaceCountMap[types.WinOddsRange2]),
							rateFormatFunc(inOrder3oddsRangeMap[types.WinOddsRange3], oddsRangeRaceCountMap[types.WinOddsRange3]),
							rateFormatFunc(inOrder3oddsRangeMap[types.WinOddsRange4], oddsRangeRaceCountMap[types.WinOddsRange4]),
							rateFormatFunc(inOrder3oddsRangeMap[types.WinOddsRange5], oddsRangeRaceCountMap[types.WinOddsRange5]),
							rateFormatFunc(inOrder3oddsRangeMap[types.WinOddsRange6], oddsRangeRaceCountMap[types.WinOddsRange6]),
							rateFormatFunc(inOrder3oddsRangeMap[types.WinOddsRange7], oddsRangeRaceCountMap[types.WinOddsRange7]),
							rateFormatFunc(inOrder3oddsRangeMap[types.WinOddsRange8], oddsRangeRaceCountMap[types.WinOddsRange8]),
						}...)
						valuesList[1][rowPosition] = append(valuesList[1][rowPosition], []interface{}{
							orderNo2MatchCount,
							inOrder2oddsRangeMap[types.WinOddsRange1],
							inOrder2oddsRangeMap[types.WinOddsRange2],
							inOrder2oddsRangeMap[types.WinOddsRange3],
							inOrder2oddsRangeMap[types.WinOddsRange4],
							inOrder2oddsRangeMap[types.WinOddsRange5],
							inOrder2oddsRangeMap[types.WinOddsRange6],
							inOrder2oddsRangeMap[types.WinOddsRange7],
							inOrder2oddsRangeMap[types.WinOddsRange8],
							orderNo3MatchCount,
							inOrder3oddsRangeMap[types.WinOddsRange1],
							inOrder3oddsRangeMap[types.WinOddsRange2],
							inOrder3oddsRangeMap[types.WinOddsRange3],
							inOrder3oddsRangeMap[types.WinOddsRange4],
							inOrder3oddsRangeMap[types.WinOddsRange5],
							inOrder3oddsRangeMap[types.WinOddsRange6],
							inOrder3oddsRangeMap[types.WinOddsRange7],
							inOrder3oddsRangeMap[types.WinOddsRange8],
						}...)
					}
				}
				data, ok = markerCombinationMap[f][markerCombinationId]
				if ok {
					switch markerCombinationId.TicketType() {
					case types.Win:
						marker, err := types.NewMarker(markerCombinationId.Value() % 10)
						if err != nil {
							return err
						}
						if marker != sheetMarker {
							continue
						}

						oddsRangeMap := s.createUnHitWinOddsRangeMap(ctx, data, 1)
						oddsRangeRaceCountMap := raceCountMap[f][markerCombinationId]
						raceCount := 0
						for _, oddsRange := range oddsRanges {
							if n, ok := oddsRangeRaceCountMap[oddsRange]; ok {
								raceCount += n
							}
						}
						matchCount := 0
						for _, calculable := range data.Calculables() {
							if calculable.OrderNo() > 1 {
								matchCount++
							}
						}

						valuesList[2] = append(valuesList[2], [][]interface{}{
							{
								f.String(),
								raceCount,
								matchCount,
								oddsRangeMap[types.WinOddsRange1],
								oddsRangeMap[types.WinOddsRange2],
								oddsRangeMap[types.WinOddsRange3],
								oddsRangeMap[types.WinOddsRange4],
								oddsRangeMap[types.WinOddsRange5],
								oddsRangeMap[types.WinOddsRange6],
								oddsRangeMap[types.WinOddsRange7],
								oddsRangeMap[types.WinOddsRange8],
							},
						}...)
					case types.Place:
						marker, err := types.NewMarker(markerCombinationId.Value() % 10)
						if err != nil {
							return err
						}
						if marker != sheetMarker {
							continue
						}

						inOrder2oddsRangeMap := s.createUnHitWinOddsRangeMap(ctx, data, 2)
						inOrder3oddsRangeMap := s.createUnHitWinOddsRangeMap(ctx, data, 3)
						oddsRangeRaceCountMap := raceCountMap[f][markerCombinationId]
						raceCount := 0
						for _, oddsRange := range oddsRanges {
							if n, ok := oddsRangeRaceCountMap[oddsRange]; ok {
								raceCount += n
							}
						}

						orderNo2UnMatchCount := 0
						orderNo3UnMatchCount := 0
						for _, calculable := range data.Calculables() {
							if calculable.OrderNo() > 2 {
								orderNo2UnMatchCount++
							}
							if calculable.OrderNo() > 3 {
								orderNo3UnMatchCount++
							}
						}

						valuesList[2][rowPosition] = append(valuesList[2][rowPosition], []interface{}{
							orderNo2UnMatchCount,
							inOrder2oddsRangeMap[types.WinOddsRange1],
							inOrder2oddsRangeMap[types.WinOddsRange2],
							inOrder2oddsRangeMap[types.WinOddsRange3],
							inOrder2oddsRangeMap[types.WinOddsRange4],
							inOrder2oddsRangeMap[types.WinOddsRange5],
							inOrder2oddsRangeMap[types.WinOddsRange6],
							inOrder2oddsRangeMap[types.WinOddsRange7],
							inOrder2oddsRangeMap[types.WinOddsRange8],
							orderNo3UnMatchCount,
							inOrder3oddsRangeMap[types.WinOddsRange1],
							inOrder3oddsRangeMap[types.WinOddsRange2],
							inOrder3oddsRangeMap[types.WinOddsRange3],
							inOrder3oddsRangeMap[types.WinOddsRange4],
							inOrder3oddsRangeMap[types.WinOddsRange5],
							inOrder3oddsRangeMap[types.WinOddsRange6],
							inOrder3oddsRangeMap[types.WinOddsRange7],
							inOrder3oddsRangeMap[types.WinOddsRange8],
						}...)
					}
				}
			}
		}

		for idx, values := range valuesList {
			writeRange := fmt.Sprintf("%s!%s", spreadSheetConfig.SheetName(), fmt.Sprintf("A%d", idx*(len(filters)+1)+1))
			_, err := s.client.Spreadsheets.Values.Update(spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
				Values: values,
			}).ValueInputOption("USER_ENTERED").Do()
			if err != nil {
				return err
			}
		}

		log.Println(ctx, fmt.Sprintf("write marker %s analysis end", sheetMarker.String()))
	}

	return nil
}

func (s *spreadSheetMarkerAnalysisRepository) createHitWinOddsRangeMap(
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

func (s *spreadSheetMarkerAnalysisRepository) createUnHitWinOddsRangeMap(
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
		if calculable.OrderNo() > inOrderNo {
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

func (s *spreadSheetMarkerAnalysisRepository) Style(
	ctx context.Context,
	analysisData *spreadsheet_entity.AnalysisData,
	filters []filter.Id,
) error {
	for _, spreadSheetConfig := range s.spreadSheetConfigs {
		var sheetMarker types.Marker
		switch spreadSheetConfig.SheetName() {
		case types.Favorite.String():
			sheetMarker = types.Favorite
		case types.Rival.String():
			sheetMarker = types.Rival
		case types.BrackTriangle.String():
			sheetMarker = types.BrackTriangle
		case types.WhiteTriangle.String():
			sheetMarker = types.WhiteTriangle
		case types.Star.String():
			sheetMarker = types.Star
		case types.Check.String():
			sheetMarker = types.Check
		default:
			return fmt.Errorf("invalid sheet name: %s", spreadSheetConfig.SheetName())
		}

		log.Println(ctx, fmt.Sprintf("write style marker %s analysis start", sheetMarker.String()))
		allMarkerCombinationIds := analysisData.AllMarkerCombinationIds()

		for _, markerCombinationId := range allMarkerCombinationIds {
			switch markerCombinationId.TicketType() {
			case types.Win:
				marker, err := types.NewMarker(markerCombinationId.Value() % 10)
				if err != nil {
					return err
				}
				if marker != sheetMarker {
					continue
				}
				for i := 0; i < 3; i++ {
					_, err = s.client.Spreadsheets.BatchUpdate(spreadSheetConfig.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
						Requests: []*sheets.Request{
							{
								RepeatCell: &sheets.RepeatCellRequest{
									Fields: "userEnteredFormat.textFormat.foregroundColor",
									Range: &sheets.GridRange{
										SheetId:          spreadSheetConfig.SheetId(),
										StartColumnIndex: 3,
										StartRowIndex:    int64(i * (1 + len(filters))),
										EndColumnIndex:   11,
										EndRowIndex:      int64(i*(1+len(filters)) + 1),
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
							{
								RepeatCell: &sheets.RepeatCellRequest{
									Fields: "userEnteredFormat.textFormat.foregroundColor",
									Range: &sheets.GridRange{
										SheetId:          spreadSheetConfig.SheetId(),
										StartColumnIndex: 12,
										StartRowIndex:    int64(i * (1 + len(filters))),
										EndColumnIndex:   20,
										EndRowIndex:      int64(i*(1+len(filters)) + 1),
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
							{
								RepeatCell: &sheets.RepeatCellRequest{
									Fields: "userEnteredFormat.textFormat.foregroundColor",
									Range: &sheets.GridRange{
										SheetId:          spreadSheetConfig.SheetId(),
										StartColumnIndex: 21,
										StartRowIndex:    int64(i * (1 + len(filters))),
										EndColumnIndex:   29,
										EndRowIndex:      int64(i*(1+len(filters)) + 1),
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
							{
								RepeatCell: &sheets.RepeatCellRequest{
									Fields: "userEnteredFormat.backgroundColor",
									Range: &sheets.GridRange{
										SheetId:          spreadSheetConfig.SheetId(),
										StartColumnIndex: 1,
										StartRowIndex:    int64(i * (1 + len(filters))),
										EndColumnIndex:   4,
										EndRowIndex:      int64(i*(1+len(filters)) + 1),
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
										SheetId:          spreadSheetConfig.SheetId(),
										StartColumnIndex: 11,
										StartRowIndex:    int64(i * (1 + len(filters))),
										EndColumnIndex:   12,
										EndRowIndex:      int64(i*(1+len(filters)) + 1),
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
										SheetId:          spreadSheetConfig.SheetId(),
										StartColumnIndex: 20,
										StartRowIndex:    int64(i * (1 + len(filters))),
										EndColumnIndex:   21,
										EndRowIndex:      int64(i*(1+len(filters)) + 1),
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
										SheetId:          spreadSheetConfig.SheetId(),
										StartColumnIndex: 3,
										StartRowIndex:    int64(i * (1 + len(filters))),
										EndColumnIndex:   11,
										EndRowIndex:      int64(i*(1+len(filters)) + 1),
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
									Fields: "userEnteredFormat.backgroundColor",
									Range: &sheets.GridRange{
										SheetId:          spreadSheetConfig.SheetId(),
										StartColumnIndex: 12,
										StartRowIndex:    int64(i * (1 + len(filters))),
										EndColumnIndex:   20,
										EndRowIndex:      int64(i*(1+len(filters)) + 1),
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
									Fields: "userEnteredFormat.backgroundColor",
									Range: &sheets.GridRange{
										SheetId:          spreadSheetConfig.SheetId(),
										StartColumnIndex: 21,
										StartRowIndex:    int64(i * (1 + len(filters))),
										EndColumnIndex:   29,
										EndRowIndex:      int64(i*(1+len(filters)) + 1),
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
									Fields: "userEnteredFormat.textFormat.bold",
									Range: &sheets.GridRange{
										SheetId:          spreadSheetConfig.SheetId(),
										StartColumnIndex: 1,
										StartRowIndex:    int64(i * (1 + len(filters))),
										EndColumnIndex:   29,
										EndRowIndex:      int64(i*(1+len(filters)) + 1),
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
										SheetId:          spreadSheetConfig.SheetId(),
										StartColumnIndex: 0,
										StartRowIndex:    int64(i*(1+len(filters)) + 1),
										EndColumnIndex:   1,
										EndRowIndex:      int64((i + 1) * (1 + len(filters))),
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
										SheetId:          spreadSheetConfig.SheetId(),
										StartColumnIndex: 0,
										StartRowIndex:    int64(i*(1+len(filters)) + 1),
										EndColumnIndex:   1,
										EndRowIndex:      int64((i + 1) * (1 + len(filters))),
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
						},
					}).Do()
					if err != nil {
						return err
					}
				}
			}
		}

		log.Println(ctx, fmt.Sprintf("write style marker %s analysis end", sheetMarker.String()))
	}

	return nil
}

func (s *spreadSheetMarkerAnalysisRepository) Clear(ctx context.Context) error {
	for _, spreadSheetConfig := range s.spreadSheetConfigs {
		requests := []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "*",
					Range: &sheets.GridRange{
						SheetId:          spreadSheetConfig.SheetId(),
						StartColumnIndex: 0,
						StartRowIndex:    0,
						EndColumnIndex:   40,
						EndRowIndex:      9999,
					},
					Cell: &sheets.CellData{},
				},
			},
		}
		_, err := s.client.Spreadsheets.BatchUpdate(spreadSheetConfig.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
			Requests: requests,
		}).Do()

		if err != nil {
			return err
		}
	}

	return nil
}
