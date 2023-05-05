package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	analyse_entity "github.com/mapserver2007/ipat-aggregator/app/domain/analyse/entity"
	betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"
	predict_entity "github.com/mapserver2007/ipat-aggregator/app/domain/predict/entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	spreadsheet_entity "github.com/mapserver2007/ipat-aggregator/app/domain/spreadsheet/entity"
	spreadsheet_vo "github.com/mapserver2007/ipat-aggregator/app/domain/spreadsheet/value_object"
	"github.com/mapserver2007/ipat-aggregator/app/repository"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	secretFileName             = "secret.json"
	spreadSheetCalcFileName    = "spreadsheet_calc.json"
	spreadSheetListFileName    = "spreadsheet_list.json"
	spreadSheetAnalyseFileName = "spreadsheet_analyse.json"
)

type SpreadSheetClient struct {
	client            *sheets.Service
	spreadSheetConfig spreadsheet_entity.SpreadSheetConfig
	sheetId           int64
}

func NewSpreadSheetClient(
	ctx context.Context,
) repository.SpreadSheetClient {
	service, spreadSheetConfig, sheetId := getSpreadSheetConfig(ctx, spreadSheetCalcFileName)
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
) repository.SpreadSheetListClient {
	service, spreadSheetConfig, sheetId := getSpreadSheetConfig(ctx, spreadSheetListFileName)
	return &SpreadSheetListClient{
		client:            service,
		spreadSheetConfig: spreadSheetConfig,
		sheetId:           sheetId,
	}
}

type SpreadSheetAnalyseClient struct {
	client            *sheets.Service
	spreadSheetConfig spreadsheet_entity.SpreadSheetAnalyseConfig
	sheetMap          map[spreadsheet_vo.AnalyseType]*sheets.SheetProperties
}

func NewSpreadSheetAnalyseClient(
	ctx context.Context,
) repository.SpreadSheetAnalyseClient {
	service, spreadSheetConfig, sheetMap := getSpreadSheetAnalyseConfig(ctx, spreadSheetAnalyseFileName)
	return &SpreadSheetAnalyseClient{
		client:            service,
		spreadSheetConfig: spreadSheetConfig,
		sheetMap:          sheetMap,
	}
}

func getSpreadSheetConfig(
	ctx context.Context,
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

func getSpreadSheetAnalyseConfig(
	ctx context.Context,
	spreadSheetConfigFileName string,
) (*sheets.Service, spreadsheet_entity.SpreadSheetAnalyseConfig, map[spreadsheet_vo.AnalyseType]*sheets.SheetProperties) {
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

	var spreadSheetConfig spreadsheet_entity.SpreadSheetAnalyseConfig
	if err = json.Unmarshal(spreadSheetConfigBytes, &spreadSheetConfig); err != nil {
		panic(err)
	}

	response, err := service.Spreadsheets.Get(spreadSheetConfig.Id).Do()
	if err != nil {
		panic(err)
	}

	var sheetReverseMap = map[string]spreadsheet_vo.AnalyseType{}
	for _, sheetName := range spreadSheetConfig.SheetNames {
		sheetReverseMap[sheetName.Name] = spreadsheet_vo.AnalyseType(sheetName.Type)
	}

	sheetMap := map[spreadsheet_vo.AnalyseType]*sheets.SheetProperties{}
	for _, sheet := range response.Sheets {
		if analyseType, ok := sheetReverseMap[sheet.Properties.Title]; ok {
			sheetMap[analyseType] = sheet.Properties
		}
	}

	return service, spreadSheetConfig, sheetMap
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

func (s *SpreadSheetClient) WriteForCurrentYearSummary(ctx context.Context, summary spreadsheet_entity.ResultSummary) error {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "E1")
	values := [][]interface{}{
		{
			"年間累計",
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

func (s *SpreadSheetClient) WriteStyleForCurrentYearSummary(ctx context.Context) error {
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			// 1行目のセルをマージ
			{
				MergeCells: &sheets.MergeCellsRequest{
					MergeType: "MERGE_ROWS",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
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
						SheetId:          s.sheetId,
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
						EndRowIndex:      40,
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
						EndRowIndex:      40,
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

func (s *SpreadSheetListClient) WriteList(ctx context.Context, records []*predict_entity.PredictEntity) (map[race_vo.RaceId]*spreadsheet_entity.ResultStyle, error) {
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
			"騎手",
			"人気",
			"オッズ",
			"対抗",
			"騎手",
			"人気",
			"オッズ",
			"1着",
			"騎手",
			"人気",
			"オッズ",
			"2着",
			"騎手",
			"人気",
			"オッズ",
		},
	}

	var rivalHorseName, rivalPopularNumber, rivalOdds, rivalJockeyName string

	sort.SliceStable(records, func(i, j int) bool {
		return records[i].Race().StartTime() > records[j].Race().StartTime()
	})
	sort.SliceStable(records, func(i, j int) bool {
		return records[i].Race().RaceDate() > records[j].Race().RaceDate()
	})

	styleMap := map[race_vo.RaceId]*spreadsheet_entity.ResultStyle{}
	for idx, record := range records {
		var (
			favoriteColor, rivalColor spreadsheet_vo.PlaceColor
			gradeClassColor           spreadsheet_vo.GradeClassColor
			repaymentComment          string
		)
		raceResults := record.Race().RaceResults()
		raceResultOfFirst := raceResults[0]
		raceResultOfSecond := raceResults[1]

		if record.FavoriteHorse().HorseName() == raceResultOfFirst.HorseName() {
			favoriteColor = spreadsheet_vo.FirstPlace
		} else if record.FavoriteHorse().HorseName() == raceResultOfSecond.HorseName() {
			favoriteColor = spreadsheet_vo.SecondPlace
		}

		if record.RivalHorse() != nil {
			if record.RivalHorse().HorseName() == raceResultOfFirst.HorseName() {
				rivalColor = spreadsheet_vo.FirstPlace
			} else if record.RivalHorse().HorseName() == raceResultOfSecond.HorseName() {
				rivalColor = spreadsheet_vo.SecondPlace
			}
		}

		switch record.Race().Class() {
		case race_vo.Grade1, race_vo.Jpn1, race_vo.JumpGrade1:
			gradeClassColor = spreadsheet_vo.Grade1
		case race_vo.Grade2, race_vo.Jpn2, race_vo.JumpGrade2:
			gradeClassColor = spreadsheet_vo.Grade2
		case race_vo.Grade3, race_vo.Jpn3, race_vo.JumpGrade3:
			gradeClassColor = spreadsheet_vo.Grade3
		}

		if record.WinningTickets() != nil {
			var comments []string
			for _, winningTicket := range record.WinningTickets() {
				comment := fmt.Sprintf("%s %s %s倍 %d円", winningTicket.BettingTicket.Name(), winningTicket.BetNumber.String(), winningTicket.Odds, winningTicket.Repayment)
				comments = append(comments, comment)
			}
			repaymentComment = strings.Join(comments, "\n")
		}

		styleMap[record.Race().RaceId()] = spreadsheet_entity.NewResultStyle(
			idx+1, favoriteColor, rivalColor, gradeClassColor, repaymentComment,
		)

		if record.RivalHorse() != nil {
			rivalHorseName = record.RivalHorse().HorseName()
			rivalPopularNumber = strconv.Itoa(record.RivalHorse().PopularNumber())
			rivalOdds = record.RivalHorse().Odds()
			rivalJockeyName = record.RivalJockey().JockeyName()
		} else {
			rivalHorseName = "-"
			rivalPopularNumber = "-"
			rivalOdds = "-"
			rivalJockeyName = "-"
		}

		values = append(values, []interface{}{
			record.Race().RaceDate().DateFormat(),
			record.Race().Class().String(),
			record.Race().CourseCategory().String(),
			fmt.Sprintf("%d%s", record.Race().Distance(), "m"),
			record.Race().TrackCondition(),
			fmt.Sprintf("=HYPERLINK(\"%s\",\"%s\")", record.Race().Url(), record.Race().RaceName()),
			record.Payment(),
			record.Repayment(),
			fmt.Sprintf("%.0f%s", float64(record.Repayment())*float64(100)/float64(record.Payment()), "%"),
			record.FavoriteHorse().HorseName(),
			record.FavoriteJockey().JockeyName(),
			record.FavoriteHorse().PopularNumber(),
			record.FavoriteHorse().Odds(),
			rivalHorseName,
			rivalJockeyName,
			rivalPopularNumber,
			rivalOdds,
			raceResultOfFirst.HorseName(),
			raceResultOfFirst.JockeyName(),
			raceResultOfFirst.PopularNumber(),
			raceResultOfFirst.Odds(),
			raceResultOfSecond.HorseName(),
			raceResultOfSecond.JockeyName(),
			raceResultOfSecond.PopularNumber(),
			raceResultOfSecond.Odds(),
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

func (s *SpreadSheetListClient) WriteStyleList(ctx context.Context, records []*predict_entity.PredictEntity, styleMap map[race_vo.RaceId]*spreadsheet_entity.ResultStyle) error {
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
						EndColumnIndex:   25,
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
						EndColumnIndex:   25,
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
						EndIndex:   11,
						SheetId:    s.sheetId,
						StartIndex: 10,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 75,
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
						StartIndex: 11,
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
						EndIndex:   14,
						SheetId:    s.sheetId,
						StartIndex: 13,
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
						StartIndex: 14,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 75,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   17,
						SheetId:    s.sheetId,
						StartIndex: 15,
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
						EndIndex:   18,
						SheetId:    s.sheetId,
						StartIndex: 17,
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
						EndIndex:   19,
						SheetId:    s.sheetId,
						StartIndex: 18,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 75,
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
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   22,
						SheetId:    s.sheetId,
						StartIndex: 21,
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
						EndIndex:   23,
						SheetId:    s.sheetId,
						StartIndex: 22,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 75,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   26,
						SheetId:    s.sheetId,
						StartIndex: 23,
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
		if style, ok := styleMap[record.Race().RaceId()]; ok {
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
						StartColumnIndex: 13,
						StartRowIndex:    int64(style.RowIndex),
						EndColumnIndex:   14,
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

	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetAnalyseClient) WriteWin(ctx context.Context, summary *analyse_entity.WinAnalyseSummary) error {
	sheetProperties, ok := s.sheetMap[spreadsheet_vo.Win]
	if !ok {
		return fmt.Errorf("sheet not found")
	}

	writeRange := fmt.Sprintf("%s!%s", sheetProperties.Title, "A1")
	var values [][]interface{}
	addValues := func(values [][]interface{}, summaries []*analyse_entity.PopularAnalyseSummary, title string) [][]interface{} {
		values = append(values, [][]interface{}{
			{
				title,
				"1人気",
				"2人気",
				"3人気",
				"4人気",
				"5人気",
				"6人気",
				"7人気",
				"8人気",
				"9人気",
				"10人気",
				"11人気",
				"12人気",
				"13人気",
				"14人気",
				"15人気",
				"16人気",
				"17人気",
				"18人気",
			},
		}...)
		capacity := 19
		betCounts := make([]interface{}, 0, capacity)
		hitCounts := make([]interface{}, 0, capacity)
		hitRates := make([]interface{}, 0, capacity)
		payoutRate := make([]interface{}, 0, capacity)
		payoutUpside := make([]interface{}, 0, capacity)
		averageOddsAtVotes := make([]interface{}, 0, capacity)
		averageOddsAtHits := make([]interface{}, 0, capacity)
		averageOddsAtUnHits := make([]interface{}, 0, capacity)
		maxOddsAtHits := make([]interface{}, 0, capacity)
		minOddsAtHits := make([]interface{}, 0, capacity)
		totalPayments := make([]interface{}, 0, capacity)
		totalPayouts := make([]interface{}, 0, capacity)
		averagePayments := make([]interface{}, 0, capacity)
		averagePayouts := make([]interface{}, 0, capacity)
		medianPayments := make([]interface{}, 0, capacity)
		medianPayouts := make([]interface{}, 0, capacity)
		maxPayouts := make([]interface{}, 0, capacity)
		minPayouts := make([]interface{}, 0, capacity)

		betCounts = append(betCounts, "投票回数")
		hitCounts = append(hitCounts, "的中回数")
		hitRates = append(hitRates, "的中率")
		payoutRate = append(payoutRate, "回収率")
		payoutUpside = append(payoutUpside, "回収上振れ率")
		averageOddsAtVotes = append(averageOddsAtVotes, "投票時平均オッズ")
		averageOddsAtHits = append(averageOddsAtHits, "的中時平均オッズ")
		averageOddsAtUnHits = append(averageOddsAtUnHits, "不的中時平均オッズ")
		maxOddsAtHits = append(maxOddsAtHits, "的中時最大オッズ")
		minOddsAtHits = append(minOddsAtHits, "的中時最小オッズ")
		totalPayments = append(totalPayments, "投票金額合計")
		totalPayouts = append(totalPayouts, "払戻金額合計")
		averagePayments = append(averagePayments, "平均投票金額")
		medianPayments = append(medianPayments, "中央値投票金額")
		averagePayouts = append(averagePayouts, "平均払戻金額")
		medianPayouts = append(medianPayouts, "中央値払戻金額")
		maxPayouts = append(maxPayouts, "最大払戻金額")
		minPayouts = append(minPayouts, "最小払戻金額")

		for _, record := range summaries {
			betCounts = append(betCounts, record.BetCount())
			hitCounts = append(hitCounts, record.HitCount())
			hitRates = append(hitRates, record.FormattedHitRate())
			payoutRate = append(payoutRate, record.FormattedPayoutRate())
			payoutUpside = append(payoutUpside, record.FormattedPayoutUpsideRate())
			averageOddsAtVotes = append(averageOddsAtVotes, record.AverageOddsAtVote())
			averageOddsAtHits = append(averageOddsAtHits, record.AverageOddsAtHit())
			averageOddsAtUnHits = append(averageOddsAtUnHits, record.AverageOddsAtUnHit())
			maxOddsAtHits = append(maxOddsAtHits, record.MaxOddsAtHit())
			minOddsAtHits = append(minOddsAtHits, record.MinOddsAtHit())
			totalPayments = append(totalPayments, record.TotalPayment())
			totalPayouts = append(totalPayouts, record.TotalPayout())
			averagePayments = append(averagePayments, record.AveragePayment())
			medianPayments = append(medianPayments, record.MedianPayment())
			averagePayouts = append(averagePayouts, record.AveragePayout())
			medianPayouts = append(medianPayouts, record.MedianPayout())
			maxPayouts = append(maxPayouts, record.MaxPayout())
			minPayouts = append(minPayouts, record.MinPayout())
		}

		values = append(values, betCounts, hitCounts, hitRates, payoutRate, payoutUpside, averageOddsAtVotes,
			averageOddsAtHits, averageOddsAtUnHits, maxOddsAtHits, minOddsAtHits, totalPayments, totalPayouts, averagePayments, medianPayments, averagePayouts, medianPayouts, maxPayouts, minPayouts)

		return values
	}

	values = addValues(values, summary.AllSummaries(), "全レース集計")
	values = addValues(values, summary.Grade1Summaries(), "JRA G1集計")
	values = addValues(values, summary.Grade2Summaries(), "JRA G2集計")
	values = addValues(values, summary.Grade3Summaries(), "JRA G3集計")
	values = addValues(values, summary.AllowanceClassSummaries(), "JRA 平場集計")

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetAnalyseClient) WriteStyleWin(ctx context.Context, summary *analyse_entity.WinAnalyseSummary) error {
	sheetProperties, _ := s.sheetMap[spreadsheet_vo.Win]

	// 全レース
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    0,
						EndColumnIndex:   1,
						EndRowIndex:      99,
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
						SheetId:    sheetProperties.SheetId,
						StartIndex: 0,
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
						EndIndex:   19,
						SheetId:    sheetProperties.SheetId,
						StartIndex: 1,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 50,
					},
					Fields: "pixelSize",
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    0,
						EndColumnIndex:   19,
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
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    0,
						EndColumnIndex:   19,
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
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    1,
						EndColumnIndex:   1,
						EndRowIndex:      6,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    6,
						EndColumnIndex:   1,
						EndRowIndex:      11,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    11,
						EndColumnIndex:   1,
						EndRowIndex:      19,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  1.0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.numberFormat",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    3,
						EndColumnIndex:   19,
						EndRowIndex:      6,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							NumberFormat: &sheets.NumberFormat{
								Type:    "PERCENT",
								Pattern: "0%",
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    19,
						EndColumnIndex:   19,
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
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    19,
						EndColumnIndex:   19,
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
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    20,
						EndColumnIndex:   1,
						EndRowIndex:      25,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    25,
						EndColumnIndex:   1,
						EndRowIndex:      30,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    30,
						EndColumnIndex:   1,
						EndRowIndex:      38,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  1.0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.numberFormat",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    22,
						EndColumnIndex:   19,
						EndRowIndex:      25,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							NumberFormat: &sheets.NumberFormat{
								Type:    "PERCENT",
								Pattern: "0%",
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    38,
						EndColumnIndex:   19,
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
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    38,
						EndColumnIndex:   19,
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
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    39,
						EndColumnIndex:   1,
						EndRowIndex:      44,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    44,
						EndColumnIndex:   1,
						EndRowIndex:      49,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    49,
						EndColumnIndex:   1,
						EndRowIndex:      57,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  1.0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.numberFormat",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    41,
						EndColumnIndex:   19,
						EndRowIndex:      44,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							NumberFormat: &sheets.NumberFormat{
								Type:    "PERCENT",
								Pattern: "0%",
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    57,
						EndColumnIndex:   19,
						EndRowIndex:      58,
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
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    57,
						EndColumnIndex:   19,
						EndRowIndex:      58,
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
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    58,
						EndColumnIndex:   1,
						EndRowIndex:      63,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    63,
						EndColumnIndex:   1,
						EndRowIndex:      68,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    68,
						EndColumnIndex:   1,
						EndRowIndex:      76,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  1.0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.numberFormat",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    60,
						EndColumnIndex:   19,
						EndRowIndex:      63,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							NumberFormat: &sheets.NumberFormat{
								Type:    "PERCENT",
								Pattern: "0%",
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    76,
						EndColumnIndex:   19,
						EndRowIndex:      77,
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
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    76,
						EndColumnIndex:   19,
						EndRowIndex:      77,
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
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    77,
						EndColumnIndex:   1,
						EndRowIndex:      82,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    82,
						EndColumnIndex:   1,
						EndRowIndex:      87,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    87,
						EndColumnIndex:   1,
						EndRowIndex:      95,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  1.0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.numberFormat",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    79,
						EndColumnIndex:   19,
						EndRowIndex:      82,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							NumberFormat: &sheets.NumberFormat{
								Type:    "PERCENT",
								Pattern: "0%",
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

	payoutUpsideIndexMap := map[int64][]int64{}
	noPaymentIndexMap := map[int64][]int64{}
	payoutUpsideRowIndexForAll := int64(5)
	payoutUpsideRowIndexForG1 := int64(24)
	payoutUpsideRowIndexForG2 := int64(43)
	payoutUpsideRowIndexForG3 := int64(62)
	payoutUpsideRowIndexForAllowance := int64(81)
	noPaymentRowIndexForAll := int64(1)
	noPaymentRowIndexForG1 := int64(20)
	noPaymentRowIndexForG2 := int64(39)
	noPaymentRowIndexForG3 := int64(58)
	noPaymentRowIndexForAllowance := int64(77)
	payoutUpsideIndexMap[payoutUpsideRowIndexForAll] = []int64{}
	payoutUpsideIndexMap[payoutUpsideRowIndexForG1] = []int64{}
	payoutUpsideIndexMap[payoutUpsideRowIndexForG2] = []int64{}
	payoutUpsideIndexMap[payoutUpsideRowIndexForG3] = []int64{}
	payoutUpsideIndexMap[payoutUpsideRowIndexForAllowance] = []int64{}
	noPaymentIndexMap[noPaymentRowIndexForAll] = []int64{}
	noPaymentIndexMap[noPaymentRowIndexForG1] = []int64{}
	noPaymentIndexMap[noPaymentRowIndexForG2] = []int64{}
	noPaymentIndexMap[noPaymentRowIndexForG3] = []int64{}
	noPaymentIndexMap[noPaymentRowIndexForAllowance] = []int64{}

	for idx, record := range summary.AllSummaries() {
		if record.PayoutUpsideRate() < 0 {
			payoutUpsideIndexMap[payoutUpsideRowIndexForAll] = append(payoutUpsideIndexMap[payoutUpsideRowIndexForAll], int64(idx+1))
		}
		if record.BetCount() == 0 {
			noPaymentIndexMap[noPaymentRowIndexForAll] = append(noPaymentIndexMap[noPaymentRowIndexForAll], int64(idx+1))
		}
	}
	for idx, record := range summary.Grade1Summaries() {
		if record.PayoutUpsideRate() < 0 {
			payoutUpsideIndexMap[payoutUpsideRowIndexForG1] = append(payoutUpsideIndexMap[payoutUpsideRowIndexForG1], int64(idx+1))
		}
		if record.BetCount() == 0 {
			noPaymentIndexMap[noPaymentRowIndexForG1] = append(noPaymentIndexMap[noPaymentRowIndexForG1], int64(idx+1))
		}
	}
	for idx, record := range summary.Grade2Summaries() {
		if record.PayoutUpsideRate() < 0 {
			payoutUpsideIndexMap[payoutUpsideRowIndexForG2] = append(payoutUpsideIndexMap[payoutUpsideRowIndexForG2], int64(idx+1))
		}
		if record.BetCount() == 0 {
			noPaymentIndexMap[noPaymentRowIndexForG2] = append(noPaymentIndexMap[noPaymentRowIndexForG2], int64(idx+1))
		}
	}
	for idx, record := range summary.Grade3Summaries() {
		if record.PayoutUpsideRate() < 0 {
			payoutUpsideIndexMap[payoutUpsideRowIndexForG3] = append(payoutUpsideIndexMap[payoutUpsideRowIndexForG3], int64(idx+1))
		}
		if record.BetCount() == 0 {
			noPaymentIndexMap[noPaymentRowIndexForG3] = append(noPaymentIndexMap[noPaymentRowIndexForG3], int64(idx+1))
		}
	}
	for idx, record := range summary.AllowanceClassSummaries() {
		if record.PayoutUpsideRate() < 0 {
			payoutUpsideIndexMap[payoutUpsideRowIndexForAllowance] = append(payoutUpsideIndexMap[payoutUpsideRowIndexForAllowance], int64(idx+1))
		}
		if record.BetCount() == 0 {
			noPaymentIndexMap[noPaymentRowIndexForAllowance] = append(noPaymentIndexMap[noPaymentRowIndexForAllowance], int64(idx+1))
		}
	}

	var requests []*sheets.Request
	// 回収上振れ率の専用style
	for rowIndex, columnIndexes := range payoutUpsideIndexMap {
		for _, columnIndex := range columnIndexes {
			requests = append(requests, &sheets.Request{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.foregroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: columnIndex,
						StartRowIndex:    rowIndex,
						EndColumnIndex:   columnIndex + 1,
						EndRowIndex:      rowIndex + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								ForegroundColor: &sheets.Color{
									Red:   1.0,
									Blue:  0,
									Green: 0,
								},
							},
						},
					},
				},
			})
		}

	}

	// 購入実績なしのセルはグレーアウト
	for rowIndex, columnIndexes := range noPaymentIndexMap {
		for _, columnIndex := range columnIndexes {
			requests = append(requests, &sheets.Request{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: columnIndex,
						StartRowIndex:    rowIndex,
						EndColumnIndex:   columnIndex + 1,
						EndRowIndex:      rowIndex + 18,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.7,
								Blue:  0.7,
								Green: 0.7,
							},
						},
					},
				},
			})
		}
	}

	_, err = s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetListClient) Clear(ctx context.Context) error {
	requests := []*sheets.Request{
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Fields: "*",
				Range: &sheets.GridRange{
					SheetId:          s.sheetId,
					StartColumnIndex: 0,
					StartRowIndex:    0,
					EndColumnIndex:   25,
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

func (s *SpreadSheetAnalyseClient) Clear(ctx context.Context) error {
	for _, sheetProperties := range s.sheetMap {
		requests := []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "*",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    0,
						EndColumnIndex:   25,
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
	}

	return nil
}
