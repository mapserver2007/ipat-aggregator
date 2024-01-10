package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"os"
	"path/filepath"
	"sort"
)

const (
	spreadSheetSummaryFileName2 = "spreadsheet_summary.json"
)

type spreadsheetSummaryRepository struct {
	client            *sheets.Service
	spreadSheetConfig *spreadsheet_entity.SpreadSheetConfig
}

func NewSpreadSheetSummaryRepository() (repository.SpreadSheetSummaryRepository, error) {
	ctx := context.Background()
	client, spreadSheetConfig, err := getSpreadSheetConfig2(ctx, spreadSheetSummaryFileName2)
	if err != nil {
		return nil, err
	}

	return &spreadsheetSummaryRepository{
		client:            client,
		spreadSheetConfig: spreadSheetConfig,
	}, nil
}

func getSpreadSheetConfig2(
	ctx context.Context,
	spreadSheetConfigFileName string,
) (*sheets.Service, *spreadsheet_entity.SpreadSheetConfig, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}

	secretFilePath, err := filepath.Abs(fmt.Sprintf("%s/secret/%s", rootPath, secretFileName))
	if err != nil {
		return nil, nil, err
	}
	spreadSheetConfigFilePath, err := filepath.Abs(fmt.Sprintf("%s/secret/%s", rootPath, spreadSheetConfigFileName))
	if err != nil {
		return nil, nil, err
	}

	credential := option.WithCredentialsFile(secretFilePath)
	service, err := sheets.NewService(ctx, credential)
	if err != nil {
		return nil, nil, err
	}

	spreadSheetConfigBytes, err := os.ReadFile(spreadSheetConfigFilePath)
	if err != nil {
		return nil, nil, err
	}

	var rawSpreadSheetConfig raw_entity.SpreadSheetConfig
	if err = json.Unmarshal(spreadSheetConfigBytes, &rawSpreadSheetConfig); err != nil {
		return nil, nil, err
	}

	response, err := service.Spreadsheets.Get(rawSpreadSheetConfig.Id).Do()
	if err != nil {
		return nil, nil, err
	}

	var spreadSheetConfig *spreadsheet_entity.SpreadSheetConfig
	for _, sheet := range response.Sheets {
		if sheet.Properties.Title == rawSpreadSheetConfig.SheetName {
			spreadSheetConfig = spreadsheet_entity.NewSpreadSheetConfig(rawSpreadSheetConfig.Id, sheet.Properties.SheetId, sheet.Properties.Title)
		}
	}

	return service, spreadSheetConfig, nil
}

func (s *spreadsheetSummaryRepository) Write(
	ctx context.Context,
	summary *spreadsheet_entity.Summary,
) error {
	log.Println(ctx, "write spreadsheet start")

	err := s.writeAllResult(ctx, summary.AllTermResult())
	if err != nil {
		return err
	}
	err = s.writeYearResult(ctx, summary.YearTermResult())
	if err != nil {
		return err
	}
	err = s.writeMonthResult(ctx, summary.MonthTermResult())
	if err != nil {
		return err
	}
	err = s.writeTicketResult(ctx, summary.TicketResultMap())
	if err != nil {
		return err
	}
	err = s.writeGradeClassResult(ctx, summary.GradeClassResultMap())
	if err != nil {
		return err
	}
	err = s.writeCourseCategoryResult(ctx, summary.CourseCategoryResultMap())
	if err != nil {
		return err
	}
	err = s.writeDistanceCategoryResult(ctx, summary.DistanceCategoryResultMap())
	if err != nil {
		return err
	}
	err = s.writeRaceCourseResult(ctx, summary.RaceCourseResultMap())
	if err != nil {
		return err
	}
	err = s.writeMonthlyResult(ctx, summary.MonthlyResults())
	if err != nil {
		return err
	}

	log.Println(ctx, "write spreadsheet end")
	return nil
}

func (s *spreadsheetSummaryRepository) writeAllResult(
	ctx context.Context,
	result *spreadsheet_entity.TicketResult,
) error {
	log.Println(ctx, "writing spreadsheet writeAllResult")
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName(), "A1")
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

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadsheetSummaryRepository) writeYearResult(
	ctx context.Context,
	result *spreadsheet_entity.TicketResult,
) error {
	log.Println(ctx, "writing spreadsheet writeYearResult")
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName(), "E1")
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

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadsheetSummaryRepository) writeMonthResult(
	ctx context.Context,
	result *spreadsheet_entity.TicketResult,
) error {
	log.Println(ctx, "writing spreadsheet writeMonthResult")
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName(), "C1")
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

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadsheetSummaryRepository) writeMonthlyResult(
	ctx context.Context,
	results map[int]*spreadsheet_entity.TicketResult,
) error {
	log.Println(ctx, "writing spreadsheet writeMonthlyResult")
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName(), "A28")
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

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadsheetSummaryRepository) writeTicketResult(
	ctx context.Context,
	results map[types.TicketType]*spreadsheet_entity.TicketResult,
) error {
	log.Println(ctx, "writing spreadsheet writeTicketResult")
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName(), "A6")
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

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadsheetSummaryRepository) writeGradeClassResult(
	ctx context.Context,
	results map[types.GradeClass]*spreadsheet_entity.TicketResult,
) error {
	log.Println(ctx, "writing spreadsheet writeGradeClassResult")
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName(), "A15")
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

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadsheetSummaryRepository) writeCourseCategoryResult(
	ctx context.Context,
	results map[types.CourseCategory]*spreadsheet_entity.TicketResult,
) error {
	log.Println(ctx, "writing spreadsheet writeCourseCategoryResult")
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName(), "I6")
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

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadsheetSummaryRepository) writeDistanceCategoryResult(
	ctx context.Context,
	results map[types.DistanceCategory]*spreadsheet_entity.TicketResult,
) error {
	log.Println(ctx, "writing spreadsheet writeDistanceCategoryResult")
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName(), "I10")
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

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadsheetSummaryRepository) writeRaceCourseResult(
	ctx context.Context,
	results map[types.RaceCourse]*spreadsheet_entity.TicketResult,
) error {
	log.Println(ctx, "writing spreadsheet writeRaceCourseResult")
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName(), "I21")
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

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadsheetSummaryRepository) Style(
	ctx context.Context,
	summary *spreadsheet_entity.Summary,
) error {
	log.Println(ctx, "write spreadsheet style start")

	err := s.writeStyleAllResult(ctx)
	if err != nil {
		return err
	}
	err = s.writeStyleYearResult(ctx)
	if err != nil {
		return err
	}
	err = s.writeStyleMonthResult(ctx)
	if err != nil {
		return err
	}
	err = s.writeStyleTicketResult(ctx)
	if err != nil {
		return err
	}
	err = s.writeStyleGradeClassResult(ctx)
	if err != nil {
		return err
	}
	err = s.writeStyleMonthlyResult(ctx, len(summary.MonthlyResults()))
	if err != nil {
		return err
	}
	err = s.writeStyleCourseCategoryResult(ctx)
	if err != nil {
		return err
	}
	err = s.writeStyleDistanceCategoryResult(ctx)
	if err != nil {
		return err
	}
	err = s.writeStyleRaceCourseResult(ctx)
	if err != nil {
		return err
	}

	log.Println(ctx, "write spreadsheet style end")

	return nil
}

func (s *spreadsheetSummaryRepository) writeStyleAllResult(ctx context.Context) error {
	log.Println(ctx, "writing spreadsheet writeStyleAllResult")
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			// 1行目のセルをマージ
			{
				MergeCells: &sheets.MergeCellsRequest{
					MergeType: "MERGE_ROWS",
					Range: &sheets.GridRange{
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:    s.spreadSheetConfig.SheetId(),
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
						SheetId:    s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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

func (s *spreadsheetSummaryRepository) writeStyleYearResult(ctx context.Context) error {
	log.Println(ctx, "writing spreadsheet writeStyleYearResult")
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			// 1行目のセルをマージ
			{
				MergeCells: &sheets.MergeCellsRequest{
					MergeType: "MERGE_ROWS",
					Range: &sheets.GridRange{
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:    s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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

func (s *spreadsheetSummaryRepository) writeStyleMonthResult(ctx context.Context) error {
	log.Println(ctx, "writing spreadsheet writeStyleMonthResult")
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			// 1行目のセルをマージ
			{
				MergeCells: &sheets.MergeCellsRequest{
					MergeType: "MERGE_ROWS",
					Range: &sheets.GridRange{
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:    s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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

func (s *spreadsheetSummaryRepository) writeStyleTicketResult(ctx context.Context) error {
	log.Println(ctx, "writing spreadsheet writeStyleTicketResult")
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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

func (s *spreadsheetSummaryRepository) writeStyleGradeClassResult(ctx context.Context) error {
	log.Println(ctx, "writing spreadsheet writeStyleGradeClassResult")
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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

func (s *spreadsheetSummaryRepository) writeStyleMonthlyResult(ctx context.Context, rowCount int) error {
	log.Println(ctx, "writing spreadsheet writeStyleMonthlyResult")
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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

func (s *spreadsheetSummaryRepository) writeStyleCourseCategoryResult(ctx context.Context) error {
	log.Println(ctx, "writing spreadsheet writeStyleCourseCategoryResult")
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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

func (s *spreadsheetSummaryRepository) writeStyleDistanceCategoryResult(ctx context.Context) error {
	log.Println(ctx, "writing spreadsheet writeStyleDistanceCategoryResult")
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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

func (s *spreadsheetSummaryRepository) writeStyleRaceCourseResult(ctx context.Context) error {
	log.Println(ctx, "writing spreadsheet writeStyleRaceCourseResult")
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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
						SheetId:          s.spreadSheetConfig.SheetId(),
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

func (s *spreadsheetSummaryRepository) Clear(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
