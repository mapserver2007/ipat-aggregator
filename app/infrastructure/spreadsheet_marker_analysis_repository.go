package infrastructure

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"google.golang.org/api/sheets/v4"
	"log"
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
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName(), "A3")

	var values [][]interface{}
	allMarkerCombinationIds := analysisData.AllMarkerCombinationIds()
	markerCombinationAnalysisMap := analysisData.MarkerCombinationAnalysisMap()
	currentTicketType := types.Win
	rowNo := 2
	headerRowNos := make([]int, 0, 7)
	headerRowNos = append(headerRowNos, rowNo)
	for _, markerCombinationId := range allMarkerCombinationIds {
		if currentTicketType != markerCombinationId.TicketType() {
			rowNo++
			values = append(values, []interface{}{})
			currentTicketType = markerCombinationId.TicketType()
			headerRowNos = append(headerRowNos, rowNo)
		}
		log.Println(ctx, fmt.Sprintf("write marker %s", markerCombinationId.String()))
		data, ok := markerCombinationAnalysisMap[markerCombinationId]
		if ok {
			values = append(values, []interface{}{
				fmt.Sprintf("%s(%d)", markerCombinationId.String(), markerCombinationId.Value()),
				data.HitRateFormat(),
				data.HitCount(),
			})
		} else {
			values = append(values, []interface{}{
				fmt.Sprintf("%s(%d)", markerCombinationId.String(), markerCombinationId.Value()),
				0,
				0,
			})
		}

		rowNo++
	}

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	s.writeHeader(ctx, analysisData.RaceCount(), headerRowNos)

	log.Println(ctx, "write marker analysis end")

	return nil
}

func (s *spreadSheetMarkerAnalysisRepository) Style(
	ctx context.Context,
	analysisData *spreadsheet_entity.AnalysisData,
) error {
	log.Println(ctx, "write style marker analysis start")
	currentTicketType := types.Win
	rowNo := 2
	headerRowNos := make([]int, 0, 7)
	headerRowNos = append(headerRowNos, rowNo)
	for _, markerCombinationId := range analysisData.AllMarkerCombinationIds() {
		if currentTicketType != markerCombinationId.TicketType() {
			rowNo++
			currentTicketType = markerCombinationId.TicketType()
			headerRowNos = append(headerRowNos, rowNo)
		}
		rowNo++
	}

	err := s.writeStyleHeader(ctx, headerRowNos)
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetMarkerAnalysisRepository) writeFilter(
	ctx context.Context,
) error {
	log.Println(ctx, "writing spreadsheet writeFilter in marker analysis")

	return nil
}

func (s *spreadSheetMarkerAnalysisRepository) writeHeader(
	ctx context.Context,
	raceCount int,
	headerRowNos []int,
) error {
	log.Println(ctx, "writing spreadsheet writeHeader in marker analysis")

	values := [][]interface{}{
		{
			fmt.Sprintf("レース数: %d, フィルタ条件: %s", raceCount, "なし"),
		},
	}
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName(), "A1")
	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	values = [][]interface{}{
		{
			"印組合せ",
			"印的中率",
			"印的中回数",
			"投票回数",
			"回収率",
			"払戻平均値",
			"払戻中央値",
			"払戻最大値",
			"払戻最小値",
			"平均人気",
			"最大人気",
			"最小人気",
			"平均オッズ",
			"最大オッズ",
			"最小オッズ",
		},
	}

	for _, rowNo := range headerRowNos {
		writeRange = fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName(), fmt.Sprintf("A%d", rowNo))
		_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
			Values: values,
		}).ValueInputOption("USER_ENTERED").Do()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *spreadSheetMarkerAnalysisRepository) writeStyleHeader(
	ctx context.Context,
	headerRowNos []int,
) error {
	log.Println(ctx, "writing spreadsheet writeStyleHeader in marker analysis")
	for _, rowNo := range headerRowNos {
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
							EndColumnIndex:   15,
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
							EndColumnIndex:   15,
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
							EndIndex:   15,
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
	}

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
