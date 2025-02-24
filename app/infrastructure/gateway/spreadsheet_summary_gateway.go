package gateway

import (
	"context"
	"fmt"
	"sort"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"
)

const (
	spreadSheetSummaryFileName = "spreadsheet_summary.json"
)

type SpreadSheetSummaryGateway interface {
	Write(ctx context.Context, summary *spreadsheet_entity.Summary) error
	Style(ctx context.Context, summary *spreadsheet_entity.Summary) error
	Clear(ctx context.Context) error
}

type spreadSheetSummaryGateway struct {
	logger *logrus.Logger
}

func NewSpreadSheetSummaryGateway(
	logger *logrus.Logger,
) SpreadSheetSummaryGateway {
	return &spreadSheetSummaryGateway{
		logger: logger,
	}
}

func (s *spreadSheetSummaryGateway) Write(
	ctx context.Context,
	summary *spreadsheet_entity.Summary,
) error {
	client, config, err := getSpreadSheetConfig(ctx, spreadSheetSummaryFileName)
	if err != nil {
		return err
	}

	s.logger.Infof("write summary start")
	err = s.writeAllResult(ctx, summary.AllTermResult(), client, config)
	if err != nil {
		return err
	}
	err = s.writeYearResult(ctx, summary.YearTermResult(), client, config)
	if err != nil {
		return err
	}
	err = s.writeMonthResult(ctx, summary.MonthTermResult(), client, config)
	if err != nil {
		return err
	}
	err = s.writeTicketResult(ctx, summary.TicketResultMap(), client, config)
	if err != nil {
		return err
	}
	err = s.writeGradeClassResult(ctx, summary.GradeClassResultMap(), client, config)
	if err != nil {
		return err
	}
	err = s.writeMonthlyResult(ctx, summary.MonthlyResults(), client, config)
	if err != nil {
		return err
	}
	err = s.writeCourseCategoryResult(ctx, summary.CourseCategoryResultMap(), client, config)
	if err != nil {
		return err
	}
	err = s.writeDistanceCategoryResult(ctx, summary.DistanceCategoryResultMap(), client, config)
	if err != nil {
		return err
	}
	err = s.writeRaceCourseResult(ctx, summary.RaceCourseResultMap(), client, config)
	if err != nil {
		return err
	}

	s.logger.Infof("write summary end")
	return nil
}

func (s *spreadSheetSummaryGateway) writeAllResult(
	ctx context.Context,
	result *spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeAllResult")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "A1")
	values := [][]interface{}{
		{
			"累計",
			"",
		},
		{
			"投資",
			result.Payment(),
		},
		{
			"回収",
			result.Payout(),
		},
		{
			"回収率",
			result.PayoutRate(),
		},
	}

	_, err := client.Spreadsheets.Values.Update(config.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetSummaryGateway) writeYearResult(
	ctx context.Context,
	result *spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeYearResult")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "E1")
	values := [][]interface{}{
		{
			"年間累計",
			"",
		},
		{
			"投資",
			result.Payment(),
		},
		{
			"回収",
			result.Payout(),
		},
		{
			"回収率",
			result.PayoutRate(),
		},
	}

	_, err := client.Spreadsheets.Values.Update(config.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetSummaryGateway) writeMonthResult(
	ctx context.Context,
	result *spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeMonthResult")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "C1")
	values := [][]interface{}{
		{
			"月間累計",
			"",
		},
		{
			"投資",
			result.Payment(),
		},
		{
			"回収",
			result.Payout(),
		},
		{
			"回収率",
			result.PayoutRate(),
		},
	}

	_, err := client.Spreadsheets.Values.Update(config.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetSummaryGateway) writeTicketResult(
	ctx context.Context,
	results map[types.TicketType]*spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeTicketResult")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "A6")
	values := [][]interface{}{
		{
			"券種別",
			"投票レース数",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	keys := make([]int, 0, len(results))
	for k := range results {
		keys = append(keys, k.Value())
	}
	sort.Ints(keys)

	for _, k := range keys {
		ticketType := types.TicketType(k)
		result := results[ticketType]
		values = append(values, []interface{}{
			ticketType.Name(),
			result.RaceCount(),
			result.BetCount(),
			result.HitCount(),
			result.HitRate(),
			result.Payment(),
			result.Payout(),
			result.PayoutRate(),
		})
	}

	_, err := client.Spreadsheets.Values.Update(config.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetSummaryGateway) writeGradeClassResult(
	ctx context.Context,
	results map[types.GradeClass]*spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeGradeClassResult")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "A15")
	values := [][]interface{}{
		{
			"券種別",
			"投票レース数",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	keys := make([]int, 0, len(results))
	for k := range results {
		keys = append(keys, k.Value())
	}
	sort.Ints(keys)

	for _, k := range keys {
		gradeClass := types.GradeClass(k)
		result := results[gradeClass]
		values = append(values, []interface{}{
			gradeClass.String(),
			result.RaceCount(),
			result.BetCount(),
			result.HitCount(),
			result.HitRate(),
			result.Payment(),
			result.Payout(),
			result.PayoutRate(),
		})
	}

	_, err := client.Spreadsheets.Values.Update(config.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetSummaryGateway) writeMonthlyResult(
	ctx context.Context,
	results map[int]*spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeMonthlyResult")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "A28")
	values := [][]interface{}{
		{
			"券種別",
			"投票レース数",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	keys := make([]int, 0, len(results))
	for month := range results {
		keys = append(keys, month)
	}
	sort.Ints(keys)

	for _, month := range keys {
		result := results[month]
		values = append(values, []interface{}{
			month,
			result.RaceCount(),
			result.BetCount(),
			result.HitCount(),
			result.HitRate(),
			result.Payment(),
			result.Payout(),
			result.PayoutRate(),
		})
	}

	_, err := client.Spreadsheets.Values.Update(config.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetSummaryGateway) writeCourseCategoryResult(
	ctx context.Context,
	results map[types.CourseCategory]*spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeCourseCategoryResult")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "I6")
	values := [][]interface{}{
		{
			"券種別",
			"投票レース数",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	keys := make([]int, 0, len(results))
	for k := range results {
		keys = append(keys, k.Value())
	}
	sort.Ints(keys)

	for _, k := range keys {
		courseCategory := types.CourseCategory(k)
		result := results[courseCategory]
		values = append(values, []interface{}{
			courseCategory.String(),
			result.RaceCount(),
			result.BetCount(),
			result.HitCount(),
			result.HitRate(),
			result.Payment(),
			result.Payout(),
			result.PayoutRate(),
		})
	}

	_, err := client.Spreadsheets.Values.Update(config.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetSummaryGateway) writeDistanceCategoryResult(
	ctx context.Context,
	results map[types.DistanceCategory]*spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeDistanceCategoryResult")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "I10")
	values := [][]interface{}{
		{
			"券種別",
			"投票レース数",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	keys := make([]int, 0, len(results))
	for k := range results {
		keys = append(keys, k.Value())
	}
	sort.Ints(keys)

	for _, k := range keys {
		distanceCategory := types.DistanceCategory(k)
		result := results[distanceCategory]
		values = append(values, []interface{}{
			distanceCategory.String(),
			result.RaceCount(),
			result.BetCount(),
			result.HitCount(),
			result.HitRate(),
			result.Payment(),
			result.Payout(),
			result.PayoutRate(),
		})
	}

	_, err := client.Spreadsheets.Values.Update(config.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetSummaryGateway) writeRaceCourseResult(
	ctx context.Context,
	results map[types.RaceCourse]*spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeRaceCourseResult")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "I21")
	values := [][]interface{}{
		{
			"券種別",
			"投票レース数",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	keys := make([]string, 0, len(results))
	for k := range results {
		keys = append(keys, k.Value())
	}
	sort.Strings(keys)

	for _, k := range keys {
		raceCourse := types.RaceCourse(k)
		result := results[raceCourse]
		values = append(values, []interface{}{
			raceCourse.Name(),
			result.RaceCount(),
			result.BetCount(),
			result.HitCount(),
			result.HitRate(),
			result.Payment(),
			result.Payout(),
			result.PayoutRate(),
		})
	}

	_, err := client.Spreadsheets.Values.Update(config.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetSummaryGateway) Style(
	ctx context.Context,
	summary *spreadsheet_entity.Summary,
) error {
	client, config, err := getSpreadSheetConfig(ctx, spreadSheetSummaryFileName)
	if err != nil {
		return err
	}

	s.logger.Infof("write spreadsheet style start")
	err = s.writeStyleAllResult(ctx, client, config)
	if err != nil {
		return err
	}
	err = s.writeStyleYearResult(ctx, client, config)
	if err != nil {
		return err
	}
	err = s.writeStyleMonthResult(ctx, client, config)
	if err != nil {
		return err
	}
	err = s.writeStyleTicketResult(ctx, client, config)
	if err != nil {
		return err
	}
	err = s.writeStyleGradeClassResult(ctx, client, config)
	if err != nil {
		return err
	}
	err = s.writeStyleMonthlyResult(ctx, len(summary.MonthlyResults()), client, config)
	if err != nil {
		return err
	}
	err = s.writeStyleCourseCategoryResult(ctx, client, config)
	if err != nil {
		return err
	}
	err = s.writeStyleDistanceCategoryResult(ctx, client, config)
	if err != nil {
		return err
	}
	err = s.writeStyleRaceCourseResult(ctx, client, config)
	if err != nil {
		return err
	}

	s.logger.Infof("write spreadsheet style end")
	return nil
}

func (s *spreadSheetSummaryGateway) writeStyleAllResult(
	ctx context.Context,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeStyleAllResult")
	_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			// 1行目のセルをマージ
			{
				MergeCells: &sheets.MergeCellsRequest{
					MergeType: "MERGE_ROWS",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 0,
						StartRowIndex:    0,
						EndColumnIndex:   2,
						EndRowIndex:      1,
					},
				},
			},
			// 1列目のセル幅調整
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   1,
						SheetId:    config.SheetId(),
						StartIndex: 0,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 95,
					},
					Fields: "pixelSize",
				},
			},
			// 2列目のセル幅調整
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   2,
						SheetId:    config.SheetId(),
						StartIndex: 1,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 90,
					},
					Fields: "pixelSize",
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 0,
						StartRowIndex:    1,
						EndColumnIndex:   1,
						EndRowIndex:      4,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
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
						StartRowIndex:    0,
						EndColumnIndex:   2,
						EndRowIndex:      4,
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

	return nil
}

func (s *spreadSheetSummaryGateway) writeStyleYearResult(
	ctx context.Context,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeStyleYearResult")
	_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			// 1行目のセルをマージ
			{
				MergeCells: &sheets.MergeCellsRequest{
					MergeType: "MERGE_ROWS",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 4,
						StartRowIndex:    0,
						EndColumnIndex:   6,
						EndRowIndex:      1,
					},
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   3,
						SheetId:    config.SheetId(),
						StartIndex: 2,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 90,
					},
					Fields: "pixelSize",
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 4,
						StartRowIndex:    1,
						EndColumnIndex:   5,
						EndRowIndex:      4,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
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
						StartColumnIndex: 4,
						StartRowIndex:    0,
						EndColumnIndex:   6,
						EndRowIndex:      4,
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

	return nil
}

func (s *spreadSheetSummaryGateway) writeStyleMonthResult(
	ctx context.Context,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeStyleMonthResult")
	_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			// 1行目のセルをマージ
			{
				MergeCells: &sheets.MergeCellsRequest{
					MergeType: "MERGE_ROWS",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 2,
						StartRowIndex:    0,
						EndColumnIndex:   4,
						EndRowIndex:      1,
					},
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   3,
						SheetId:    config.SheetId(),
						StartIndex: 2,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 90,
					},
					Fields: "pixelSize",
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 2,
						StartRowIndex:    1,
						EndColumnIndex:   3,
						EndRowIndex:      4,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
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
						StartColumnIndex: 2,
						StartRowIndex:    0,
						EndColumnIndex:   4,
						EndRowIndex:      4,
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

	return nil
}

func (s *spreadSheetSummaryGateway) writeStyleTicketResult(
	ctx context.Context,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeStyleTicketResult")
	_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 0,
						StartRowIndex:    5,
						EndColumnIndex:   8,
						EndRowIndex:      6,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.6,
								Blue:  0,
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
						StartColumnIndex: 0,
						StartRowIndex:    6,
						EndColumnIndex:   1,
						EndRowIndex:      15,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
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
						StartRowIndex:    5,
						EndColumnIndex:   8,
						EndRowIndex:      6,
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
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 0,
						StartRowIndex:    6,
						EndColumnIndex:   1,
						EndRowIndex:      15,
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

	return nil
}

func (s *spreadSheetSummaryGateway) writeStyleGradeClassResult(
	ctx context.Context,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeStyleGradeClassResult")
	_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 0,
						StartRowIndex:    14,
						EndColumnIndex:   8,
						EndRowIndex:      15,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.6,
								Blue:  0,
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
						StartColumnIndex: 0,
						StartRowIndex:    15,
						EndColumnIndex:   1,
						EndRowIndex:      27,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
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
						StartRowIndex:    14,
						EndColumnIndex:   8,
						EndRowIndex:      15,
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
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 0,
						StartRowIndex:    15,
						EndColumnIndex:   1,
						EndRowIndex:      27,
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

	return nil
}

func (s *spreadSheetSummaryGateway) writeStyleMonthlyResult(
	ctx context.Context,
	rowCount int,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeStyleMonthlyResult")
	_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 0,
						StartRowIndex:    27,
						EndColumnIndex:   8,
						EndRowIndex:      28,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.6,
								Blue:  0,
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
						StartColumnIndex: 0,
						StartRowIndex:    28,
						EndColumnIndex:   1,
						EndRowIndex:      28 + int64(rowCount),
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
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
						StartRowIndex:    27,
						EndColumnIndex:   8,
						EndRowIndex:      28,
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
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 0,
						StartRowIndex:    28,
						EndColumnIndex:   1,
						EndRowIndex:      28 + int64(rowCount),
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

	return nil
}

func (s *spreadSheetSummaryGateway) writeStyleCourseCategoryResult(
	ctx context.Context,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeStyleCourseCategoryResult")
	_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 8,
						StartRowIndex:    5,
						EndColumnIndex:   16,
						EndRowIndex:      6,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.6,
								Blue:  0,
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
						StartColumnIndex: 8,
						StartRowIndex:    6,
						EndColumnIndex:   9,
						EndRowIndex:      9,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
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
						StartColumnIndex: 8,
						StartRowIndex:    5,
						EndColumnIndex:   16,
						EndRowIndex:      6,
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
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 8,
						StartRowIndex:    5,
						EndColumnIndex:   9,
						EndRowIndex:      9,
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

	return nil
}

func (s *spreadSheetSummaryGateway) writeStyleDistanceCategoryResult(
	ctx context.Context,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeStyleDistanceCategoryResult")
	_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 8,
						StartRowIndex:    9,
						EndColumnIndex:   16,
						EndRowIndex:      10,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.6,
								Blue:  0,
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
						StartColumnIndex: 8,
						StartRowIndex:    10,
						EndColumnIndex:   9,
						EndRowIndex:      20,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
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
						StartColumnIndex: 8,
						StartRowIndex:    9,
						EndColumnIndex:   16,
						EndRowIndex:      10,
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
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 8,
						StartRowIndex:    10,
						EndColumnIndex:   9,
						EndRowIndex:      20,
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

	return nil
}

func (s *spreadSheetSummaryGateway) writeStyleRaceCourseResult(
	ctx context.Context,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeStyleRaceCourseResult")
	_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 8,
						StartRowIndex:    20,
						EndColumnIndex:   16,
						EndRowIndex:      21,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.6,
								Blue:  0,
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
						StartColumnIndex: 8,
						StartRowIndex:    21,
						EndColumnIndex:   9,
						EndRowIndex:      43,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
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
						StartColumnIndex: 8,
						StartRowIndex:    20,
						EndColumnIndex:   16,
						EndRowIndex:      21,
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
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 8,
						StartRowIndex:    21,
						EndColumnIndex:   9,
						EndRowIndex:      43,
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

	return nil
}

func (s *spreadSheetSummaryGateway) Clear(
	ctx context.Context,
) error {
	client, config, err := getSpreadSheetConfig(ctx, spreadSheetSummaryFileName)
	if err != nil {
		return err
	}

	requests := []*sheets.Request{
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Fields: "*",
				Range: &sheets.GridRange{
					SheetId:          config.SheetId(),
					StartColumnIndex: 0,
					StartRowIndex:    0,
					EndColumnIndex:   16,
					EndRowIndex:      9999,
				},
				Cell: &sheets.CellData{},
			},
		},
	}
	_, err = client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
}
