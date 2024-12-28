package gateway

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"google.golang.org/api/sheets/v4"
	"log"
)

const (
	spreadSheetAnalysisPlaceFileName = "spreadsheet_analysis_place.json"
)

type SpreadSheetAnalysisPlaceGateway interface {
	Write(ctx context.Context, firstPlaceMap, secondPlaceMap, thirdPlaceMap map[types.Marker]map[filter.Id]*spreadsheet_entity.AnalysisPlace, analysisFilters []filter.Id) error
	Style(ctx context.Context, firstPlaceMap, secondPlaceMap, thirdPlaceMap map[types.Marker]map[filter.Id]*spreadsheet_entity.AnalysisPlace, analysisFilters []filter.Id) error
	Clear(ctx context.Context) error
}

type spreadSheetAnalysisPlaceGateway struct{}

func NewSpreadSheetAnalysisPlaceGateway() SpreadSheetAnalysisPlaceGateway {
	return &spreadSheetAnalysisPlaceGateway{}
}

func (s *spreadSheetAnalysisPlaceGateway) Write(
	ctx context.Context,
	firstPlaceMap,
	secondPlaceMap,
	thirdPlaceMap map[types.Marker]map[filter.Id]*spreadsheet_entity.AnalysisPlace,
	analysisFilters []filter.Id,
) error {
	client, configs, err := getSpreadSheetConfigs(ctx, spreadSheetAnalysisPlaceFileName)
	if err != nil {
		return err
	}

	for _, config := range configs {
		var sheetMarker types.Marker
		switch config.SheetName() {
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
			return fmt.Errorf("invalid sheet name: %s", config.SheetName())
		}

		log.Println(ctx, fmt.Sprintf("write analysis place %s start", sheetMarker.String()))
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
				types.WinOddsRange9.String(),
				"2着以内率",
				types.WinOddsRange1.String(),
				types.WinOddsRange2.String(),
				types.WinOddsRange3.String(),
				types.WinOddsRange4.String(),
				types.WinOddsRange5.String(),
				types.WinOddsRange6.String(),
				types.WinOddsRange7.String(),
				types.WinOddsRange8.String(),
				types.WinOddsRange9.String(),
				"3着以内率",
				types.WinOddsRange1.String(),
				types.WinOddsRange2.String(),
				types.WinOddsRange3.String(),
				types.WinOddsRange4.String(),
				types.WinOddsRange5.String(),
				types.WinOddsRange6.String(),
				types.WinOddsRange7.String(),
				types.WinOddsRange8.String(),
				types.WinOddsRange9.String(),
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
				types.WinOddsRange9.String(),
				"2着以内率",
				types.WinOddsRange1.String(),
				types.WinOddsRange2.String(),
				types.WinOddsRange3.String(),
				types.WinOddsRange4.String(),
				types.WinOddsRange5.String(),
				types.WinOddsRange6.String(),
				types.WinOddsRange7.String(),
				types.WinOddsRange8.String(),
				types.WinOddsRange9.String(),
				"3着以内率",
				types.WinOddsRange1.String(),
				types.WinOddsRange2.String(),
				types.WinOddsRange3.String(),
				types.WinOddsRange4.String(),
				types.WinOddsRange5.String(),
				types.WinOddsRange6.String(),
				types.WinOddsRange7.String(),
				types.WinOddsRange8.String(),
				types.WinOddsRange9.String(),
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
				types.WinOddsRange9.String(),
				"3着以下数",
				types.WinOddsRange1.String(),
				types.WinOddsRange2.String(),
				types.WinOddsRange3.String(),
				types.WinOddsRange4.String(),
				types.WinOddsRange5.String(),
				types.WinOddsRange6.String(),
				types.WinOddsRange7.String(),
				types.WinOddsRange8.String(),
				types.WinOddsRange9.String(),
				"4着以下数",
				types.WinOddsRange1.String(),
				types.WinOddsRange2.String(),
				types.WinOddsRange3.String(),
				types.WinOddsRange4.String(),
				types.WinOddsRange5.String(),
				types.WinOddsRange6.String(),
				types.WinOddsRange7.String(),
				types.WinOddsRange8.String(),
				types.WinOddsRange9.String(),
			},
		}

		firstFilterMap := firstPlaceMap[sheetMarker]
		secondFilterMap := secondPlaceMap[sheetMarker]
		thirdFilterMap := thirdPlaceMap[sheetMarker]

		for _, analysisFilter := range analysisFilters {
			var filterName string
			for _, f := range analysisFilter.OriginFilters() {
				filterName += f.String()
			}

			firstPlace := firstFilterMap[analysisFilter]
			secondPlace := secondFilterMap[analysisFilter]
			thirdPlace := thirdFilterMap[analysisFilter]

			valuesList[0] = append(valuesList[0], [][]interface{}{
				{
					filterName,
					firstPlace.RateData().RaceCount(),
					firstPlace.RateData().HitRateFormat(),
					firstPlace.RateData().OddsRange1RateFormat(),
					firstPlace.RateData().OddsRange2RateFormat(),
					firstPlace.RateData().OddsRange3RateFormat(),
					firstPlace.RateData().OddsRange4RateFormat(),
					firstPlace.RateData().OddsRange5RateFormat(),
					firstPlace.RateData().OddsRange6RateFormat(),
					firstPlace.RateData().OddsRange7RateFormat(),
					firstPlace.RateData().OddsRange8RateFormat(),
					firstPlace.RateData().OddsRange9RateFormat(),
					secondPlace.RateData().HitRateFormat(),
					secondPlace.RateData().OddsRange1RateFormat(),
					secondPlace.RateData().OddsRange2RateFormat(),
					secondPlace.RateData().OddsRange3RateFormat(),
					secondPlace.RateData().OddsRange4RateFormat(),
					secondPlace.RateData().OddsRange5RateFormat(),
					secondPlace.RateData().OddsRange6RateFormat(),
					secondPlace.RateData().OddsRange7RateFormat(),
					secondPlace.RateData().OddsRange8RateFormat(),
					secondPlace.RateData().OddsRange9RateFormat(),
					thirdPlace.RateData().HitRateFormat(),
					thirdPlace.RateData().OddsRange1RateFormat(),
					thirdPlace.RateData().OddsRange2RateFormat(),
					thirdPlace.RateData().OddsRange3RateFormat(),
					thirdPlace.RateData().OddsRange4RateFormat(),
					thirdPlace.RateData().OddsRange5RateFormat(),
					thirdPlace.RateData().OddsRange6RateFormat(),
					thirdPlace.RateData().OddsRange7RateFormat(),
					thirdPlace.RateData().OddsRange8RateFormat(),
					thirdPlace.RateData().OddsRange9RateFormat(),
				},
			}...)
			valuesList[1] = append(valuesList[1], [][]interface{}{
				{
					filterName,
					firstPlace.HitCountData().RaceCount(),
					firstPlace.HitCountData().HitCount(),
					firstPlace.HitCountData().OddsRange1Count(),
					firstPlace.HitCountData().OddsRange2Count(),
					firstPlace.HitCountData().OddsRange3Count(),
					firstPlace.HitCountData().OddsRange4Count(),
					firstPlace.HitCountData().OddsRange5Count(),
					firstPlace.HitCountData().OddsRange6Count(),
					firstPlace.HitCountData().OddsRange7Count(),
					firstPlace.HitCountData().OddsRange8Count(),
					firstPlace.HitCountData().OddsRange9Count(),
					secondPlace.HitCountData().HitCount(),
					secondPlace.HitCountData().OddsRange1Count(),
					secondPlace.HitCountData().OddsRange2Count(),
					secondPlace.HitCountData().OddsRange3Count(),
					secondPlace.HitCountData().OddsRange4Count(),
					secondPlace.HitCountData().OddsRange5Count(),
					secondPlace.HitCountData().OddsRange6Count(),
					secondPlace.HitCountData().OddsRange7Count(),
					secondPlace.HitCountData().OddsRange8Count(),
					secondPlace.HitCountData().OddsRange9Count(),
					thirdPlace.HitCountData().HitCount(),
					thirdPlace.HitCountData().OddsRange1Count(),
					thirdPlace.HitCountData().OddsRange2Count(),
					thirdPlace.HitCountData().OddsRange3Count(),
					thirdPlace.HitCountData().OddsRange4Count(),
					thirdPlace.HitCountData().OddsRange5Count(),
					thirdPlace.HitCountData().OddsRange6Count(),
					thirdPlace.HitCountData().OddsRange7Count(),
					thirdPlace.HitCountData().OddsRange8Count(),
					thirdPlace.HitCountData().OddsRange9Count(),
				},
			}...)
			valuesList[2] = append(valuesList[2], [][]interface{}{
				{
					filterName,
					firstPlace.UnHitCountData().RaceCount(),
					firstPlace.UnHitCountData().UnHitCount(),
					firstPlace.UnHitCountData().OddsRange1Count(),
					firstPlace.UnHitCountData().OddsRange2Count(),
					firstPlace.UnHitCountData().OddsRange3Count(),
					firstPlace.UnHitCountData().OddsRange4Count(),
					firstPlace.UnHitCountData().OddsRange5Count(),
					firstPlace.UnHitCountData().OddsRange6Count(),
					firstPlace.UnHitCountData().OddsRange7Count(),
					firstPlace.UnHitCountData().OddsRange8Count(),
					firstPlace.UnHitCountData().OddsRange9Count(),
					secondPlace.UnHitCountData().UnHitCount(),
					secondPlace.UnHitCountData().OddsRange1Count(),
					secondPlace.UnHitCountData().OddsRange2Count(),
					secondPlace.UnHitCountData().OddsRange3Count(),
					secondPlace.UnHitCountData().OddsRange4Count(),
					secondPlace.UnHitCountData().OddsRange5Count(),
					secondPlace.UnHitCountData().OddsRange6Count(),
					secondPlace.UnHitCountData().OddsRange7Count(),
					secondPlace.UnHitCountData().OddsRange8Count(),
					secondPlace.UnHitCountData().OddsRange9Count(),
					thirdPlace.UnHitCountData().UnHitCount(),
					thirdPlace.UnHitCountData().OddsRange1Count(),
					thirdPlace.UnHitCountData().OddsRange2Count(),
					thirdPlace.UnHitCountData().OddsRange3Count(),
					thirdPlace.UnHitCountData().OddsRange4Count(),
					thirdPlace.UnHitCountData().OddsRange5Count(),
					thirdPlace.UnHitCountData().OddsRange6Count(),
					thirdPlace.UnHitCountData().OddsRange7Count(),
					thirdPlace.UnHitCountData().OddsRange8Count(),
					thirdPlace.UnHitCountData().OddsRange9Count(),
				},
			}...)
		}

		for idx, values := range valuesList {
			writeRange := fmt.Sprintf("%s!%s", config.SheetName(), fmt.Sprintf("A%d", idx*(len(analysisFilters)+1)+1))
			_, err := client.Spreadsheets.Values.Update(config.SpreadSheetId(), writeRange, &sheets.ValueRange{
				Values: values,
			}).ValueInputOption("USER_ENTERED").Do()
			if err != nil {
				return err
			}
		}

		log.Println(ctx, fmt.Sprintf("write analysis place %s end", sheetMarker.String()))
	}

	return nil
}

func (s *spreadSheetAnalysisPlaceGateway) Style(
	ctx context.Context,
	firstPlaceMap,
	secondPlaceMap,
	thirdPlaceMap map[types.Marker]map[filter.Id]*spreadsheet_entity.AnalysisPlace,
	analysisFilters []filter.Id,
) error {
	client, configs, err := getSpreadSheetConfigs(ctx, spreadSheetAnalysisPlaceFileName)
	if err != nil {
		return err
	}

	var requests []*sheets.Request
	for _, config := range configs {
		var sheetMarker types.Marker
		switch config.SheetName() {
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
			return fmt.Errorf("invalid sheet name: %s", config.SheetName())
		}

		log.Println(ctx, fmt.Sprintf("write style analysis place %s start", sheetMarker.String()))

		firstFilterMap := firstPlaceMap[sheetMarker]
		secondFilterMap := secondPlaceMap[sheetMarker]
		thirdFilterMap := thirdPlaceMap[sheetMarker]

		for rowIdx, analysisFilter := range analysisFilters {
			firstPlace := firstFilterMap[analysisFilter]
			secondPlace := secondFilterMap[analysisFilter]
			thirdPlace := thirdFilterMap[analysisFilter]

			for colIdx := 0; colIdx < 9; colIdx++ {
				var cellColorType types.CellColorType
				switch colIdx {
				case 0:
					cellColorType = firstPlace.RateStyle().OddsRange1CellColorType()
				case 1:
					cellColorType = firstPlace.RateStyle().OddsRange2CellColorType()
				case 2:
					cellColorType = firstPlace.RateStyle().OddsRange3CellColorType()
				case 3:
					cellColorType = firstPlace.RateStyle().OddsRange4CellColorType()
				case 4:
					cellColorType = firstPlace.RateStyle().OddsRange5CellColorType()
				case 5:
					cellColorType = firstPlace.RateStyle().OddsRange6CellColorType()
				case 6:
					cellColorType = firstPlace.RateStyle().OddsRange7CellColorType()
				case 7:
					cellColorType = firstPlace.RateStyle().OddsRange8CellColorType()
				case 8:
					cellColorType = firstPlace.RateStyle().OddsRange9CellColorType()
				}

				rowSpace := int64(1)
				colSpace := int64(3)
				requests = append(requests, []*sheets.Request{
					{
						RepeatCell: &sheets.RepeatCellRequest{
							Fields: "userEnteredFormat.backgroundColor",
							Range: &sheets.GridRange{
								SheetId:          config.SheetId(),
								StartColumnIndex: colSpace + int64(colIdx),
								StartRowIndex:    rowSpace + int64(rowIdx),
								EndColumnIndex:   colSpace + int64(colIdx) + 1,
								EndRowIndex:      rowSpace + int64(rowIdx) + 1,
							},
							Cell: &sheets.CellData{
								UserEnteredFormat: &sheets.CellFormat{
									BackgroundColor: s.getCellColor(cellColorType),
								},
							},
						},
					},
				}...)
			}

			for colIdx := 0; colIdx < 9; colIdx++ {
				var cellColorType types.CellColorType
				switch colIdx {
				case 0:
					cellColorType = secondPlace.RateStyle().OddsRange1CellColorType()
				case 1:
					cellColorType = secondPlace.RateStyle().OddsRange2CellColorType()
				case 2:
					cellColorType = secondPlace.RateStyle().OddsRange3CellColorType()
				case 3:
					cellColorType = secondPlace.RateStyle().OddsRange4CellColorType()
				case 4:
					cellColorType = secondPlace.RateStyle().OddsRange5CellColorType()
				case 5:
					cellColorType = secondPlace.RateStyle().OddsRange6CellColorType()
				case 6:
					cellColorType = secondPlace.RateStyle().OddsRange7CellColorType()
				case 7:
					cellColorType = secondPlace.RateStyle().OddsRange8CellColorType()
				case 8:
					cellColorType = secondPlace.RateStyle().OddsRange9CellColorType()
				}

				rowSpace := int64(1)
				colSpace := int64(13)
				requests = append(requests, []*sheets.Request{
					{
						RepeatCell: &sheets.RepeatCellRequest{
							Fields: "userEnteredFormat.backgroundColor",
							Range: &sheets.GridRange{
								SheetId:          config.SheetId(),
								StartColumnIndex: colSpace + int64(colIdx),
								StartRowIndex:    rowSpace + int64(rowIdx),
								EndColumnIndex:   colSpace + int64(colIdx) + 1,
								EndRowIndex:      rowSpace + int64(rowIdx) + 1,
							},
							Cell: &sheets.CellData{
								UserEnteredFormat: &sheets.CellFormat{
									BackgroundColor: s.getCellColor(cellColorType),
								},
							},
						},
					},
				}...)
			}

			for colIdx := 0; colIdx < 9; colIdx++ {
				var cellColorType types.CellColorType
				switch colIdx {
				case 0:
					cellColorType = thirdPlace.RateStyle().OddsRange1CellColorType()
				case 1:
					cellColorType = thirdPlace.RateStyle().OddsRange2CellColorType()
				case 2:
					cellColorType = thirdPlace.RateStyle().OddsRange3CellColorType()
				case 3:
					cellColorType = thirdPlace.RateStyle().OddsRange4CellColorType()
				case 4:
					cellColorType = thirdPlace.RateStyle().OddsRange5CellColorType()
				case 5:
					cellColorType = thirdPlace.RateStyle().OddsRange6CellColorType()
				case 6:
					cellColorType = thirdPlace.RateStyle().OddsRange7CellColorType()
				case 7:
					cellColorType = thirdPlace.RateStyle().OddsRange8CellColorType()
				case 8:
					cellColorType = thirdPlace.RateStyle().OddsRange9CellColorType()
				}

				rowSpace := int64(1)
				colSpace := int64(23)
				requests = append(requests, []*sheets.Request{
					{
						RepeatCell: &sheets.RepeatCellRequest{
							Fields: "userEnteredFormat.backgroundColor",
							Range: &sheets.GridRange{
								SheetId:          config.SheetId(),
								StartColumnIndex: colSpace + int64(colIdx),
								StartRowIndex:    rowSpace + int64(rowIdx),
								EndColumnIndex:   colSpace + int64(colIdx) + 1,
								EndRowIndex:      rowSpace + int64(rowIdx) + 1,
							},
							Cell: &sheets.CellData{
								UserEnteredFormat: &sheets.CellFormat{
									BackgroundColor: s.getCellColor(cellColorType),
								},
							},
						},
					},
				}...)
			}
		}

		for i := 0; i < 3; i++ {
			requests = append(requests, []*sheets.Request{
				{
					RepeatCell: &sheets.RepeatCellRequest{
						Fields: "userEnteredFormat.textFormat.foregroundColor",
						Range: &sheets.GridRange{
							SheetId:          config.SheetId(),
							StartColumnIndex: 3,
							StartRowIndex:    int64(i * (1 + len(analysisFilters))),
							EndColumnIndex:   12,
							EndRowIndex:      int64(i*(1+len(analysisFilters)) + 1),
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
							SheetId:          config.SheetId(),
							StartColumnIndex: 13,
							StartRowIndex:    int64(i * (1 + len(analysisFilters))),
							EndColumnIndex:   22,
							EndRowIndex:      int64(i*(1+len(analysisFilters)) + 1),
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
							SheetId:          config.SheetId(),
							StartColumnIndex: 23,
							StartRowIndex:    int64(i * (1 + len(analysisFilters))),
							EndColumnIndex:   32,
							EndRowIndex:      int64(i*(1+len(analysisFilters)) + 1),
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
							SheetId:          config.SheetId(),
							StartColumnIndex: 1,
							StartRowIndex:    int64(i * (1 + len(analysisFilters))),
							EndColumnIndex:   4,
							EndRowIndex:      int64(i*(1+len(analysisFilters)) + 1),
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
							SheetId:          config.SheetId(),
							StartColumnIndex: 12,
							StartRowIndex:    int64(i * (1 + len(analysisFilters))),
							EndColumnIndex:   13,
							EndRowIndex:      int64(i*(1+len(analysisFilters)) + 1),
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
							SheetId:          config.SheetId(),
							StartColumnIndex: 22,
							StartRowIndex:    int64(i * (1 + len(analysisFilters))),
							EndColumnIndex:   23,
							EndRowIndex:      int64(i*(1+len(analysisFilters)) + 1),
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
							SheetId:          config.SheetId(),
							StartColumnIndex: 3,
							StartRowIndex:    int64(i * (1 + len(analysisFilters))),
							EndColumnIndex:   12,
							EndRowIndex:      int64(i*(1+len(analysisFilters)) + 1),
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
							SheetId:          config.SheetId(),
							StartColumnIndex: 13,
							StartRowIndex:    int64(i * (1 + len(analysisFilters))),
							EndColumnIndex:   22,
							EndRowIndex:      int64(i*(1+len(analysisFilters)) + 1),
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
							SheetId:          config.SheetId(),
							StartColumnIndex: 23,
							StartRowIndex:    int64(i * (1 + len(analysisFilters))),
							EndColumnIndex:   32,
							EndRowIndex:      int64(i*(1+len(analysisFilters)) + 1),
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
							SheetId:          config.SheetId(),
							StartColumnIndex: 1,
							StartRowIndex:    int64(i * (1 + len(analysisFilters))),
							EndColumnIndex:   32,
							EndRowIndex:      int64(i*(1+len(analysisFilters)) + 1),
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
							SheetId:          config.SheetId(),
							StartColumnIndex: 0,
							StartRowIndex:    int64(i*(1+len(analysisFilters)) + 1),
							EndColumnIndex:   1,
							EndRowIndex:      int64((i + 1) * (1 + len(analysisFilters))),
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
							SheetId:          config.SheetId(),
							StartColumnIndex: 0,
							StartRowIndex:    int64(i*(1+len(analysisFilters)) + 1),
							EndColumnIndex:   1,
							EndRowIndex:      int64((i + 1) * (1 + len(analysisFilters))),
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

		_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
			Requests: requests,
		}).Do()
		if err != nil {
			return err
		}

		log.Println(ctx, fmt.Sprintf("write style analysis place %s end", sheetMarker.String()))
	}

	return nil
}

func (s *spreadSheetAnalysisPlaceGateway) Clear(ctx context.Context) error {
	client, configs, err := getSpreadSheetConfigs(ctx, spreadSheetAnalysisPlaceFileName)
	if err != nil {
		return err
	}

	for _, config := range configs {
		requests := []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "*",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 0,
						StartRowIndex:    0,
						EndColumnIndex:   40,
						EndRowIndex:      9999,
					},
					Cell: &sheets.CellData{},
				},
			},
		}
		_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
			Requests: requests,
		}).Do()

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *spreadSheetAnalysisPlaceGateway) getCellColor(
	colorType types.CellColorType,
) *sheets.Color {
	switch colorType {
	case types.FirstColor:
		return &sheets.Color{
			Red:   1.0,
			Green: 0.937,
			Blue:  0.498,
		}
	case types.SecondColor:
		return &sheets.Color{
			Red:   0.796,
			Green: 0.871,
			Blue:  1.0,
		}
	case types.ThirdColor:
		return &sheets.Color{
			Red:   0.937,
			Green: 0.78,
			Blue:  0.624,
		}
	}
	return &sheets.Color{
		Red:   1.0,
		Blue:  1.0,
		Green: 1.0,
	}
}
