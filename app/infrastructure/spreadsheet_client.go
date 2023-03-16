package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	betting_ticket_vo "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/value_object"
	predict_entity "github.com/mapserver2007/tools/baken/app/domain/predict/entity"
	race_vo "github.com/mapserver2007/tools/baken/app/domain/race/value_object"
	spreadsheet_entity "github.com/mapserver2007/tools/baken/app/domain/spreadsheet/entity"
	spreadsheet_vo "github.com/mapserver2007/tools/baken/app/domain/spreadsheet/value_object"
	"github.com/mapserver2007/tools/baken/app/repository"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type SpreadSheetClient struct {
	client            *sheets.Service
	spreadSheetConfig spreadsheet_entity.SpreadSheetConfig
	sheetId           int64
}

func NewSpreadSheetClient(
	ctx context.Context,
	secretFileName,
	spreadSheetConfigFileName string,
) repository.SpreadSheetClient {
	service, spreadSheetConfig, sheetId := getSpreadSheetConfig(ctx, secretFileName, spreadSheetConfigFileName)

	return &SpreadSheetClient{
		client:            service,
		spreadSheetConfig: spreadSheetConfig,
		sheetId:           sheetId,
	}
}

type SpreadSheetListClient struct {
	client            *sheets.Service
	spreadSheetConfig spreadsheet_entity.SpreadSheetConfig
	sheetId           int64
}

func NewSpreadSheetListClient(
	ctx context.Context,
	secretFileName,
	spreadSheetConfigFileName string,
) repository.SpreadSheetListClient {
	service, spreadSheetConfig, sheetId := getSpreadSheetConfig(ctx, secretFileName, spreadSheetConfigFileName)

	return &SpreadSheetListClient{
		client:            service,
		spreadSheetConfig: spreadSheetConfig,
		sheetId:           sheetId,
	}
}

func getSpreadSheetConfig(
	ctx context.Context,
	secretFileName,
	spreadSheetConfigFileName string,
) (*sheets.Service, spreadsheet_entity.SpreadSheetConfig, int64) {
	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	secretFilePath, err := filepath.Abs(fmt.Sprintf("%s/secret/%s", rootPath, secretFileName))
	if err != nil {
		panic(err)
	}
	spreadSheetConfigFilePath, err := filepath.Abs(fmt.Sprintf("%s/secret/%s", rootPath, spreadSheetConfigFileName))
	if err != nil {
		panic(err)
	}

	credential := option.WithCredentialsFile(secretFilePath)
	service, err := sheets.NewService(ctx, credential)
	if err != nil {
		panic(err)
	}

	spreadSheetConfigBytes, err := os.ReadFile(spreadSheetConfigFilePath)
	if err != nil {
		panic(err)
	}

	var spreadSheetConfig spreadsheet_entity.SpreadSheetConfig
	if err = json.Unmarshal(spreadSheetConfigBytes, &spreadSheetConfig); err != nil {
		panic(err)
	}

	response, err := service.Spreadsheets.Get(spreadSheetConfig.Id).Do()
	if err != nil {
		panic(err)
	}

	var sheetId int64
	for _, sheet := range response.Sheets {
		if sheet.Properties.Title == spreadSheetConfig.SheetName {
			sheetId = sheet.Properties.SheetId
		}
	}

	return service, spreadSheetConfig, sheetId
}

func (s *SpreadSheetClient) WriteForTotalSummary(ctx context.Context, summary spreadsheet_entity.ResultSummary) error {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "A1")
	values := [][]interface{}{
		{
			"累計",
			"",
		},
		{
			"投資",
			summary.Payments,
		},
		{
			"回収",
			summary.Repayments,
		},
		{
			"回収率",
			summary.CalcReturnOnInvestment(),
		},
	}

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteForCurrentMonthSummary(ctx context.Context, summary spreadsheet_entity.ResultSummary) error {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "C1")
	values := [][]interface{}{
		{
			"月間累計",
			"",
		},
		{
			"投資",
			summary.Payments,
		},
		{
			"回収",
			summary.Repayments,
		},
		{
			"回収率",
			summary.CalcReturnOnInvestment(),
		},
	}

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteForTotalBettingTicketRateSummary(ctx context.Context, summary spreadsheet_entity.BettingTicketSummary) error {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "A6")
	values := [][]interface{}{
		{
			"券種別",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	var bettingTicketKindList []betting_ticket_vo.BettingTicket
	for bettingTicketKind := range summary.BettingTicketRates {
		bettingTicketKindList = append(bettingTicketKindList, bettingTicketKind)
	}

	sort.Slice(bettingTicketKindList, func(i, j int) bool {
		return bettingTicketKindList[i].Value() < bettingTicketKindList[j].Value()
	})

	var totalResultRate spreadsheet_entity.ResultRate
	for _, bettingTicketKind := range bettingTicketKindList {
		bettingTicketRate, _ := summary.BettingTicketRates[bettingTicketKind]
		values = append(values, []interface{}{
			bettingTicketKind.Name(),
			bettingTicketRate.VoteCount,
			bettingTicketRate.HitCount,
			bettingTicketRate.HitRateFormat(),
			bettingTicketRate.Payments,
			bettingTicketRate.Repayments,
			bettingTicketRate.ReturnOnInvestmentFormat(),
		})

		totalResultRate.VoteCount += bettingTicketRate.VoteCount
		totalResultRate.HitCount += bettingTicketRate.HitCount
		totalResultRate.Payments += bettingTicketRate.Payments
		totalResultRate.Repayments += bettingTicketRate.Repayments
	}

	values = append(values, []interface{}{
		"累計",
		totalResultRate.VoteCount,
		totalResultRate.HitCount,
		totalResultRate.HitRateFormat(),
		totalResultRate.Payments,
		totalResultRate.Repayments,
		totalResultRate.ReturnOnInvestmentFormat(),
	})

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteForRaceClassRateSummary(ctx context.Context, summary spreadsheet_entity.RaceClassSummary) error {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "A15")

	values := [][]interface{}{
		{
			"クラス別",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	var dateList []race_vo.GradeClass
	for key := range summary.RaceClassRates {
		dateList = append(dateList, key)
	}
	sort.Slice(dateList, func(i, j int) bool {
		return dateList[i] < dateList[j]
	})

	for _, raceClass := range dateList {
		monthlyRate := summary.RaceClassRates[raceClass]
		values = append(values, []interface{}{
			raceClass.String(),
			monthlyRate.VoteCount,
			monthlyRate.HitCount,
			monthlyRate.HitRateFormat(),
			monthlyRate.Payments,
			monthlyRate.Repayments,
			monthlyRate.ReturnOnInvestmentFormat(),
		})
	}

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteForCourseCategoryRateSummary(ctx context.Context, summary spreadsheet_entity.CourseCategorySummary) error {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "H6")

	values := [][]interface{}{
		{
			"路面別",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	var dateList []race_vo.CourseCategory
	for key := range summary.CourseCategoryRates {
		dateList = append(dateList, key)
	}
	sort.Slice(dateList, func(i, j int) bool {
		return dateList[i] < dateList[j]
	})

	for _, courseCategory := range dateList {
		courseCategoryRate := summary.CourseCategoryRates[courseCategory]
		values = append(values, []interface{}{
			courseCategory.String(),
			courseCategoryRate.VoteCount,
			courseCategoryRate.HitCount,
			courseCategoryRate.HitRateFormat(),
			courseCategoryRate.Payments,
			courseCategoryRate.Repayments,
			courseCategoryRate.ReturnOnInvestmentFormat(),
		})
	}

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteForDistanceCategoryRateSummary(ctx context.Context, summary spreadsheet_entity.DistanceCategorySummary) error {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "H10")

	values := [][]interface{}{
		{
			"距離別",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	var dateList []race_vo.DistanceCategory
	for key := range summary.DistanceCategoryRates {
		dateList = append(dateList, key)
	}
	sort.Slice(dateList, func(i, j int) bool {
		return dateList[i] < dateList[j]
	})

	for _, distanceCategory := range dateList {
		distanceCategoryRate := summary.DistanceCategoryRates[distanceCategory]
		values = append(values, []interface{}{
			distanceCategory.String(),
			distanceCategoryRate.VoteCount,
			distanceCategoryRate.HitCount,
			distanceCategoryRate.HitRateFormat(),
			distanceCategoryRate.Payments,
			distanceCategoryRate.Repayments,
			distanceCategoryRate.ReturnOnInvestmentFormat(),
		})
	}

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteForRaceCourseRateSummary(ctx context.Context, summary spreadsheet_entity.RaceCourseSummary) error {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "H21")

	values := [][]interface{}{
		{
			"開催別",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	var dateList []race_vo.RaceCourse
	for key := range summary.RaceCourseRates {
		dateList = append(dateList, key)
	}
	sort.Slice(dateList, func(i, j int) bool {
		return dateList[i].Value() < dateList[j].Value()
	})

	for _, raceCourse := range dateList {
		raceCourseRate := summary.RaceCourseRates[raceCourse]
		values = append(values, []interface{}{
			raceCourse.Name(),
			raceCourseRate.VoteCount,
			raceCourseRate.HitCount,
			raceCourseRate.HitRateFormat(),
			raceCourseRate.Payments,
			raceCourseRate.Repayments,
			raceCourseRate.ReturnOnInvestmentFormat(),
		})
	}

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteForMonthlyRateSummary(ctx context.Context, summary spreadsheet_entity.MonthlySummary) error {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "A20")

	values := [][]interface{}{
		{
			"月別",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	var dateList []int
	for key := range summary.MonthlyRates {
		dateList = append(dateList, key)
	}
	sort.Slice(dateList, func(i, j int) bool {
		return dateList[i] > dateList[j]
	})

	for _, date := range dateList {
		monthlyRate := summary.MonthlyRates[date]
		values = append(values, []interface{}{
			strconv.Itoa(date),
			monthlyRate.VoteCount,
			monthlyRate.HitCount,
			monthlyRate.HitRateFormat(),
			monthlyRate.Payments,
			monthlyRate.Repayments,
			monthlyRate.ReturnOnInvestmentFormat(),
		})
	}

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteStyleForTotalSummary(ctx context.Context) error {
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			// 1行目のセルをマージ
			{
				MergeCells: &sheets.MergeCellsRequest{
					MergeType: "MERGE_ROWS",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
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
						SheetId:    s.sheetId,
						StartIndex: 0,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 90,
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
						SheetId:    s.sheetId,
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
						SheetId:          s.sheetId,
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
						SheetId:          s.sheetId,
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

func (s *SpreadSheetClient) WriteStyleForCurrentMonthlySummary(ctx context.Context) error {
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			// 1行目のセルをマージ
			{
				MergeCells: &sheets.MergeCellsRequest{
					MergeType: "MERGE_ROWS",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
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
						SheetId:    s.sheetId,
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
						SheetId:          s.sheetId,
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
						SheetId:          s.sheetId,
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

func (s *SpreadSheetClient) WriteStyleForTotalBettingTicketRateSummary(ctx context.Context) error {
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    5,
						EndColumnIndex:   7,
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
						SheetId:          s.sheetId,
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
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    5,
						EndColumnIndex:   7,
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
						SheetId:          s.sheetId,
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

func (s *SpreadSheetClient) WriteStyleForRaceClassRateSummary(ctx context.Context) error {
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    14,
						EndColumnIndex:   7,
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
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    15,
						EndColumnIndex:   1,
						EndRowIndex:      19,
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
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    14,
						EndColumnIndex:   7,
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
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    15,
						EndColumnIndex:   1,
						EndRowIndex:      19,
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

func (s *SpreadSheetClient) WriteStyleForCourseCategoryRateSummary(ctx context.Context) error {
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 7,
						StartRowIndex:    5,
						EndColumnIndex:   14,
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
						SheetId:          s.sheetId,
						StartColumnIndex: 7,
						StartRowIndex:    6,
						EndColumnIndex:   8,
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
						SheetId:          s.sheetId,
						StartColumnIndex: 7,
						StartRowIndex:    5,
						EndColumnIndex:   14,
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
						SheetId:          s.sheetId,
						StartColumnIndex: 7,
						StartRowIndex:    5,
						EndColumnIndex:   8,
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

func (s *SpreadSheetClient) WriteStyleForDistanceCategoryRateSummary(ctx context.Context) error {
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 7,
						StartRowIndex:    9,
						EndColumnIndex:   14,
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
						SheetId:          s.sheetId,
						StartColumnIndex: 7,
						StartRowIndex:    10,
						EndColumnIndex:   8,
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
						SheetId:          s.sheetId,
						StartColumnIndex: 7,
						StartRowIndex:    9,
						EndColumnIndex:   14,
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
						SheetId:          s.sheetId,
						StartColumnIndex: 7,
						StartRowIndex:    10,
						EndColumnIndex:   8,
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

func (s *SpreadSheetClient) WriteStyleForRaceCourseRateSummary(ctx context.Context) error {
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 7,
						StartRowIndex:    20,
						EndColumnIndex:   14,
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
						SheetId:          s.sheetId,
						StartColumnIndex: 7,
						StartRowIndex:    21,
						EndColumnIndex:   8,
						EndRowIndex:      39,
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
						SheetId:          s.sheetId,
						StartColumnIndex: 7,
						StartRowIndex:    20,
						EndColumnIndex:   14,
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
						SheetId:          s.sheetId,
						StartColumnIndex: 7,
						StartRowIndex:    21,
						EndColumnIndex:   8,
						EndRowIndex:      39,
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

func (s *SpreadSheetClient) WriteStyleForMonthlyRateSummary(ctx context.Context, summary spreadsheet_entity.MonthlySummary) error {
	rowCount := len(summary.MonthlyRates)
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    19,
						EndColumnIndex:   7,
						EndRowIndex:      20,
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
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    20,
						EndColumnIndex:   1,
						EndRowIndex:      20 + int64(rowCount),
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
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    19,
						EndColumnIndex:   7,
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
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    20,
						EndColumnIndex:   1,
						EndRowIndex:      20 + int64(rowCount),
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

func (s SpreadSheetListClient) WriteList(ctx context.Context, records []*predict_entity.PredictEntity) (map[race_vo.RaceId]*spreadsheet_entity.ResultStyle, error) {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "A1")
	values := [][]interface{}{
		{
			"レース条件",
			"",
			"",
			"",
			"",
			"",
			"投資額",
			"回収額",
			"回収率",
			"本命",
			"人気",
			"オッズ",
			"対抗",
			"人気",
			"オッズ",
			"1着",
			"人気",
			"オッズ",
			"2着",
			"人気",
			"オッズ",
		},
	}

	var rivalHorseName, rivalPopularNumber, rivalOdds string

	sort.SliceStable(records, func(i, j int) bool {
		return records[i].Race.RaceNumber > records[j].Race.RaceNumber
	})
	sort.SliceStable(records, func(i, j int) bool {
		return records[i].Race.RaceDate > records[j].Race.RaceDate
	})

	styleMap := map[race_vo.RaceId]*spreadsheet_entity.ResultStyle{}
	for idx, record := range records {
		var (
			favoriteColor, rivalColor spreadsheet_vo.PlaceColor
			gradeClassColor           spreadsheet_vo.GradeClassColor
			repaymentComment          string
		)
		if record.FavoriteHorse.HorseName == record.Race.RaceResults[0].HorseName {
			favoriteColor = spreadsheet_vo.FirstPlace
		} else if record.FavoriteHorse.HorseName == record.Race.RaceResults[1].HorseName {
			favoriteColor = spreadsheet_vo.SecondPlace
		}

		if record.RivalHorse != nil {
			if record.RivalHorse.HorseName == record.Race.RaceResults[0].HorseName {
				rivalColor = spreadsheet_vo.FirstPlace
			} else if record.RivalHorse.HorseName == record.Race.RaceResults[1].HorseName {
				rivalColor = spreadsheet_vo.SecondPlace
			}
		}

		switch record.Race.Class {
		case race_vo.Grade1, race_vo.Jpn1, race_vo.JumpGrade1:
			gradeClassColor = spreadsheet_vo.Grade1
		case race_vo.Grade2, race_vo.Jpn2, race_vo.JumpGrade2:
			gradeClassColor = spreadsheet_vo.Grade2
		case race_vo.Grade3, race_vo.Jpn3, race_vo.JumpGrade3:
			gradeClassColor = spreadsheet_vo.Grade3
		}

		if record.WinningTickets != nil {
			var comments []string
			for _, winningTicket := range record.WinningTickets {
				comment := fmt.Sprintf("%s %s %s倍 %d円", winningTicket.BettingTicket.Name(), winningTicket.BetNumber.String(), winningTicket.Odds, winningTicket.Repayment)
				comments = append(comments, comment)
			}
			repaymentComment = strings.Join(comments, "\n")
		}

		styleMap[record.Race.RaceId] = spreadsheet_entity.NewResultStyle(
			idx+1, favoriteColor, rivalColor, gradeClassColor, repaymentComment,
		)

		if record.RivalHorse != nil {
			rivalHorseName = record.RivalHorse.HorseName
			rivalPopularNumber = strconv.Itoa(record.RivalHorse.PopularNumber)
			rivalOdds = record.RivalHorse.Odds
		} else {
			rivalHorseName = "-"
			rivalPopularNumber = "-"
			rivalOdds = "-"
		}

		values = append(values, []interface{}{
			record.Race.RaceDate.DateFormat(),
			record.Race.Class.String(),
			record.Race.CourseCategory.String(),
			fmt.Sprintf("%d%s", record.Race.Distance, "m"),
			record.Race.TrackCondition,
			fmt.Sprintf("=HYPERLINK(\"%s\",\"%s\")", record.Race.Url, record.Race.RaceName),
			record.Payment,
			record.Repayment,
			fmt.Sprintf("%.0f%s", float64(record.Repayment)*float64(100)/float64(record.Payment), "%"),
			record.FavoriteHorse.HorseName,
			record.FavoriteHorse.PopularNumber,
			record.FavoriteHorse.Odds,
			rivalHorseName,
			rivalPopularNumber,
			rivalOdds,
			record.Race.RaceResults[0].HorseName,
			record.Race.RaceResults[0].PopularNumber,
			record.Race.RaceResults[0].Odds,
			record.Race.RaceResults[1].HorseName,
			record.Race.RaceResults[1].PopularNumber,
			record.Race.RaceResults[1].Odds,
		})
	}

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return nil, err
	}

	return styleMap, nil
}

func (s SpreadSheetListClient) WriteStyleList(ctx context.Context, records []*predict_entity.PredictEntity, styleMap map[race_vo.RaceId]*spreadsheet_entity.ResultStyle) error {
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				MergeCells: &sheets.MergeCellsRequest{
					MergeType: "MERGE_ROWS",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    0,
						EndColumnIndex:   6,
						EndRowIndex:      1,
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    0,
						EndColumnIndex:   21,
						EndRowIndex:      1,
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
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    0,
						EndColumnIndex:   21,
						EndRowIndex:      1,
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
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   1,
						SheetId:    s.sheetId,
						StartIndex: 0,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 80,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   2,
						SheetId:    s.sheetId,
						StartIndex: 1,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 30,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   3,
						SheetId:    s.sheetId,
						StartIndex: 2,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 45,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   4,
						SheetId:    s.sheetId,
						StartIndex: 3,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 50,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   5,
						SheetId:    s.sheetId,
						StartIndex: 4,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 25,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   6,
						SheetId:    s.sheetId,
						StartIndex: 5,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 130,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   9,
						SheetId:    s.sheetId,
						StartIndex: 6,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 60,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   10,
						SheetId:    s.sheetId,
						StartIndex: 9,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 135,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   12,
						SheetId:    s.sheetId,
						StartIndex: 10,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 50,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   13,
						SheetId:    s.sheetId,
						StartIndex: 12,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 135,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   15,
						SheetId:    s.sheetId,
						StartIndex: 13,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 50,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   16,
						SheetId:    s.sheetId,
						StartIndex: 15,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 135,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   18,
						SheetId:    s.sheetId,
						StartIndex: 16,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 50,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   19,
						SheetId:    s.sheetId,
						StartIndex: 18,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 135,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   21,
						SheetId:    s.sheetId,
						StartIndex: 19,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 50,
					},
					Fields: "pixelSize",
				},
			},
		},
	}).Do()

	if err != nil {
		return err
	}

	var requests []*sheets.Request
	for _, record := range records {
		if style, ok := styleMap[record.Race.RaceId]; ok {
			if style.GradeClassColor != spreadsheet_vo.NonGrade {
				color := &sheets.Color{
					Red:   1.0,
					Blue:  1.0,
					Green: 1.0,
				}
				if style.GradeClassColor == spreadsheet_vo.Grade1 {
					color = &sheets.Color{
						Red:   1.0,
						Green: 0.937,
						Blue:  0.498,
					}
				} else if style.GradeClassColor == spreadsheet_vo.Grade2 {
					color = &sheets.Color{
						Red:   0.796,
						Green: 0.871,
						Blue:  1.0,
					}
				} else if style.GradeClassColor == spreadsheet_vo.Grade3 {
					color = &sheets.Color{
						Red:   0.937,
						Green: 0.78,
						Blue:  0.624,
					}
				}

				cellRequest := &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 1,
						StartRowIndex:    int64(style.RowIndex),
						EndColumnIndex:   2,
						EndRowIndex:      int64(style.RowIndex) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: color,
						},
					},
				}
				requests = append(requests, &sheets.Request{
					RepeatCell: cellRequest,
				})
			}
			if style.FavoriteColor != spreadsheet_vo.OtherPlace {
				color := &sheets.Color{
					Red:   1.0,
					Blue:  1.0,
					Green: 1.0,
				}
				if style.FavoriteColor == spreadsheet_vo.FirstPlace {
					color = &sheets.Color{
						Red:   1.0,
						Green: 0.937,
						Blue:  0.498,
					}
				} else if style.FavoriteColor == spreadsheet_vo.SecondPlace {
					color = &sheets.Color{
						Red:   0.796,
						Green: 0.871,
						Blue:  1.0,
					}
				}

				cellRequest := &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 9,
						StartRowIndex:    int64(style.RowIndex),
						EndColumnIndex:   10,
						EndRowIndex:      int64(style.RowIndex) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: color,
						},
					},
				}
				requests = append(requests, &sheets.Request{
					RepeatCell: cellRequest,
				})
			}
			if style.RivalColor != spreadsheet_vo.OtherPlace {
				color := &sheets.Color{
					Red:   1.0,
					Blue:  1.0,
					Green: 1.0,
				}
				if style.RivalColor == spreadsheet_vo.FirstPlace {
					color = &sheets.Color{
						Red:   1.0,
						Green: 0.937,
						Blue:  0.498,
					}
				} else if style.RivalColor == spreadsheet_vo.SecondPlace {
					color = &sheets.Color{
						Red:   0.796,
						Green: 0.871,
						Blue:  1.0,
					}
				}

				cellRequest := &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 12,
						StartRowIndex:    int64(style.RowIndex),
						EndColumnIndex:   13,
						EndRowIndex:      int64(style.RowIndex) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: color,
						},
					},
				}
				requests = append(requests, &sheets.Request{
					RepeatCell: cellRequest,
				})
			}
			if len(style.RepaymentComment) > 0 {
				cellRequest := &sheets.RepeatCellRequest{
					Fields: "note",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 7,
						StartRowIndex:    int64(style.RowIndex),
						EndColumnIndex:   8,
						EndRowIndex:      int64(style.RowIndex) + 1,
					},
					Cell: &sheets.CellData{
						Note: style.RepaymentComment,
					},
				}
				requests = append(requests, &sheets.Request{
					RepeatCell: cellRequest,
				})
			}
		}
	}

	_, err = s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	return nil
}

func (s SpreadSheetListClient) Clear(ctx context.Context) error {
	requests := []*sheets.Request{
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Fields: "*",
				Range: &sheets.GridRange{
					SheetId:          s.sheetId,
					StartColumnIndex: 0,
					StartRowIndex:    0,
					EndColumnIndex:   21,
					EndRowIndex:      9999,
				},
				Cell: &sheets.CellData{},
			},
		},
	}
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
}
