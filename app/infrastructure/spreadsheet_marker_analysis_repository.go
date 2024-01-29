package infrastructure

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"google.golang.org/api/sheets/v4"
	"log"
	"strconv"
)

const (
	spreadSheetMarkerAnalysisFileName = "spreadsheet_marker_analysis.json"
)

type spreadSheetMarkerAnalysisRepository struct {
	client            *sheets.Service
	spreadSheetConfig *spreadsheet_entity.SpreadSheetConfig
}

func NewSpreadSheetMarkerAnalysisRepository() (repository.SpreadSheetMarkerAnalysisRepository, error) {
	ctx := context.Background()
	client, spreadSheetConfig, err := getSpreadSheetConfig(ctx, spreadSheetMarkerAnalysisFileName)
	if err != nil {
		return nil, err
	}

	return &spreadSheetMarkerAnalysisRepository{
		client:            client,
		spreadSheetConfig: spreadSheetConfig,
	}, nil
}

func (s *spreadSheetMarkerAnalysisRepository) Write(
	ctx context.Context,
	analysisData *spreadsheet_entity.AnalysisData,
) error {
	log.Println(ctx, "write marker analysis start")
	values := [][]interface{}{
		{
			fmt.Sprintf("レース数: %d, フィルタ条件: %s", analysisData.RaceCount(), "なし"),
		},
	}

	allMarkerCombinationIds := analysisData.AllMarkerCombinationIds() // TODO 多分いらなくなる
	markerCombinationAnalysisMap := analysisData.MarkerCombinationAnalysisMap()
	currentTicketType := types.UnknownTicketType
	for _, markerCombinationId := range allMarkerCombinationIds {
		if currentTicketType != markerCombinationId.TicketType() {
			currentTicketType = markerCombinationId.TicketType()
			switch currentTicketType {
			case types.Win:
				values = append(values, [][]interface{}{
					{
						"印組合せ",
						"印的中率",
						"印的中回数",
						"投票回数",
						"回収率",
						types.WinOddsRange1.String(),
						types.WinOddsRange2.String(),
						types.WinOddsRange3.String(),
						types.WinOddsRange4.String(),
						types.WinOddsRange5.String(),
						types.WinOddsRange6.String(),
						types.WinOddsRange7.String(),
						types.WinOddsRange8.String(),
					},
				}...)
			default:
				// TODO 単勝だけにとりあえず注力するので塞いでおく
				//values = append(values, [][]interface{}{
				//	{
				//		"印組合せ",
				//		"印的中率",
				//		"印的中回数",
				//		"投票回数",
				//		"回収率",
				//		"払戻平均値",
				//		"払戻中央値",
				//		"払戻最大値",
				//		"払戻最小値",
				//		"平均人気",
				//		"最大人気",
				//		"最小人気",
				//		"平均オッズ",
				//		"最大オッズ",
				//		"最小オッズ",
				//	},
				//}...)
			}
		}
		log.Println(ctx, fmt.Sprintf("write marker %s", markerCombinationId.String()))
		data, ok := markerCombinationAnalysisMap[markerCombinationId]

		switch markerCombinationId.TicketType() {
		case types.Win:
			oddsRangeMap := s.createWinOddsRangeMap(ctx, data)
			if ok {
				values = append(values, []interface{}{
					fmt.Sprintf("%s(%d)", markerCombinationId.String(), markerCombinationId.Value()),
					data.HitRateFormat(),
					data.HitCount(),
					"",
					"",
					fmt.Sprintf("%s%s", strconv.FormatFloat(float64(oddsRangeMap[types.WinOddsRange1])*float64(100)/float64(analysisData.RaceCount()), 'f', 2, 64), "%"),
					fmt.Sprintf("%s%s", strconv.FormatFloat(float64(oddsRangeMap[types.WinOddsRange2])*float64(100)/float64(analysisData.RaceCount()), 'f', 2, 64), "%"),
					fmt.Sprintf("%s%s", strconv.FormatFloat(float64(oddsRangeMap[types.WinOddsRange3])*float64(100)/float64(analysisData.RaceCount()), 'f', 2, 64), "%"),
					fmt.Sprintf("%s%s", strconv.FormatFloat(float64(oddsRangeMap[types.WinOddsRange4])*float64(100)/float64(analysisData.RaceCount()), 'f', 2, 64), "%"),
					fmt.Sprintf("%s%s", strconv.FormatFloat(float64(oddsRangeMap[types.WinOddsRange5])*float64(100)/float64(analysisData.RaceCount()), 'f', 2, 64), "%"),
					fmt.Sprintf("%s%s", strconv.FormatFloat(float64(oddsRangeMap[types.WinOddsRange6])*float64(100)/float64(analysisData.RaceCount()), 'f', 2, 64), "%"),
					fmt.Sprintf("%s%s", strconv.FormatFloat(float64(oddsRangeMap[types.WinOddsRange7])*float64(100)/float64(analysisData.RaceCount()), 'f', 2, 64), "%"),
					fmt.Sprintf("%s%s", strconv.FormatFloat(float64(oddsRangeMap[types.WinOddsRange8])*float64(100)/float64(analysisData.RaceCount()), 'f', 2, 64), "%"),
				})
			} else {
				values = append(values, []interface{}{
					fmt.Sprintf("%s(%d)", markerCombinationId.String(), markerCombinationId.Value()),
					0,
					0,
					"",
					"",
				})
			}
		default:
			// TODO 単勝だけにとりあえず注力するので塞いでおく
			//if ok {
			//	values = append(values, []interface{}{
			//		fmt.Sprintf("%s(%d)", markerCombinationId.String(), markerCombinationId.Value()),
			//		data.HitRateFormat(),
			//		data.HitCount(),
			//		"",
			//		"",
			//	})
			//} else {
			//	values = append(values, []interface{}{
			//		fmt.Sprintf("%s(%d)", markerCombinationId.String(), markerCombinationId.Value()),
			//		0,
			//		0,
			//		"",
			//		"",
			//	})
			//}
		}
	}

	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName(), fmt.Sprintf("A1"))
	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	log.Println(ctx, "write marker analysis end")

	return nil
}

func (s *spreadSheetMarkerAnalysisRepository) createWinOddsRangeMap(
	ctx context.Context,
	markerCombinationAnalysis *spreadsheet_entity.MarkerCombinationAnalysis,
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

	for _, decimalOdds := range markerCombinationAnalysis.Odds() {
		odds := decimalOdds.InexactFloat64()
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

	return oddsRangeMap
}

func (s *spreadSheetMarkerAnalysisRepository) Style(
	ctx context.Context,
	analysisData *spreadsheet_entity.AnalysisData,
) error {
	log.Println(ctx, "write style marker analysis start")
	rowNo := 2
	currentTicketType := types.UnknownTicketType
	allMarkerCombinationIds := analysisData.AllMarkerCombinationIds()
	for _, markerCombinationId := range allMarkerCombinationIds {
		if currentTicketType != markerCombinationId.TicketType() {
			currentTicketType = markerCombinationId.TicketType()
			colNo := 0
			switch currentTicketType {
			case types.Win:
				colNo = 13
			default:
				continue // TODO 単勝だけにとりあえず注力するので塞いでおく
				colNo = 15
			}
			_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
				Requests: []*sheets.Request{
					{
						RepeatCell: &sheets.RepeatCellRequest{
							Fields: "userEnteredFormat.textFormat.foregroundColor",
							Range: &sheets.GridRange{
								SheetId:          s.spreadSheetConfig.SheetId(),
								StartColumnIndex: 0,
								StartRowIndex:    int64(rowNo - 1),
								EndColumnIndex:   1,
								EndRowIndex:      int64(rowNo),
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
								SheetId:          s.spreadSheetConfig.SheetId(),
								StartColumnIndex: 0,
								StartRowIndex:    int64(rowNo - 1),
								EndColumnIndex:   1,
								EndRowIndex:      int64(rowNo),
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
							Fields: "userEnteredFormat.backgroundColor",
							Range: &sheets.GridRange{
								SheetId:          s.spreadSheetConfig.SheetId(),
								StartColumnIndex: 1,
								StartRowIndex:    int64(rowNo - 1),
								EndColumnIndex:   int64(colNo),
								EndRowIndex:      int64(rowNo),
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
								StartColumnIndex: 0,
								StartRowIndex:    int64(rowNo - 1),
								EndColumnIndex:   int64(colNo),
								EndRowIndex:      int64(rowNo),
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
						UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
							Range: &sheets.DimensionRange{
								Dimension:  "COLUMNS",
								EndIndex:   1,
								SheetId:    s.spreadSheetConfig.SheetId(),
								StartIndex: 0,
							},
							Properties: &sheets.DimensionProperties{
								PixelSize: 120,
							},
							Fields: "pixelSize",
						},
					},
					{
						UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
							Range: &sheets.DimensionRange{
								Dimension:  "COLUMNS",
								EndIndex:   int64(colNo),
								SheetId:    s.spreadSheetConfig.SheetId(),
								StartIndex: 1,
							},
							Properties: &sheets.DimensionProperties{
								PixelSize: 80,
							},
							Fields: "pixelSize",
						},
					},
				},
			}).Do()

			if err != nil {
				return err
			}

			rowNo++
		}
		rowNo++
	}

	return nil
}

func (s *spreadSheetMarkerAnalysisRepository) writeFilter(
	ctx context.Context,
) error {
	log.Println(ctx, "writing spreadsheet writeFilter in marker analysis")

	return nil
}

func (s *spreadSheetMarkerAnalysisRepository) Clear(ctx context.Context) error {
	requests := []*sheets.Request{
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Fields: "*",
				Range: &sheets.GridRange{
					SheetId:          s.spreadSheetConfig.SheetId(),
					StartColumnIndex: 0,
					StartRowIndex:    0,
					EndColumnIndex:   16,
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
