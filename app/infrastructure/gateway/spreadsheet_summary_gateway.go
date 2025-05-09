package gateway

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"
)

const (
	spreadSheetSummaryFileName   = "spreadsheet_summary.json"
	spreadSheetSummaryV2FileName = "spreadsheet_summary_v2.json"
)

type SpreadSheetSummaryGateway interface {
	WriteV2(ctx context.Context, summary *spreadsheet_entity.Summary) error
	StyleV2(ctx context.Context, summary *spreadsheet_entity.Summary) error
	ClearV2(ctx context.Context) error
}

type spreadSheetSummaryGateway struct {
	spreadSheetConfigGateway SpreadSheetConfigGateway
	logger                   *logrus.Logger
}

func NewSpreadSheetSummaryGateway(
	logger *logrus.Logger,
	spreadSheetConfigGateway SpreadSheetConfigGateway,
) SpreadSheetSummaryGateway {
	return &spreadSheetSummaryGateway{
		spreadSheetConfigGateway: spreadSheetConfigGateway,
		logger:                   logger,
	}
}

func (s *spreadSheetSummaryGateway) WriteV2(
	ctx context.Context,
	summary *spreadsheet_entity.Summary,
) error {
	client, config, err := s.spreadSheetConfigGateway.GetConfig(ctx, spreadSheetSummaryV2FileName)
	if err != nil {
		return err
	}

	s.logger.Infof("write summary v2 start")

	if err = s.writeTotalResultV2(
		summary.AllTermResult(),
		summary.YearTermResult(),
		summary.MonthTermResult(),
		summary.WeekTermResult(),
		client,
		config,
		1,
	); err != nil {
		return err
	}

	if err = s.writeWeeklyResultV2(
		summary.WeeklyResults(),
		client,
		config,
		1,
	); err != nil {
		return err
	}

	if err = s.writeMonthlyResultV2(
		summary.MonthlyResults(),
		client,
		config,
		2+len(summary.WeeklyResults()),
	); err != nil {
		return err
	}

	if err = s.writeTicketResultV2(
		summary.TicketResultMap(),
		client,
		config,
		"H",
		1,
		"券種(全)",
	); err != nil {
		return err
	}

	if err = s.writeTicketResultV2(
		summary.TicketYearlyResultMap(),
		client,
		config,
		"H",
		2+len(summary.TicketResultMap()),
		"券種(年)",
	); err != nil {
		return err
	}

	if err = s.writeTicketResultV2(
		summary.TicketMonthlyResultMap(),
		client,
		config,
		"H",
		3+len(summary.TicketResultMap())+len(summary.TicketYearlyResultMap()),
		"券種(月)",
	); err != nil {
		return err
	}

	if err = s.writeGradeClassResultV2(
		summary.GradeClassResultMap(),
		client,
		config,
		"M",
		1,
		"クラス(全)",
	); err != nil {
		return err
	}

	if err = s.writeGradeClassResultV2(
		summary.GradeClassYearlyResultMap(),
		client,
		config,
		"M",
		2+len(summary.GradeClassResultMap()),
		"クラス(年)",
	); err != nil {
		return err
	}

	if err = s.writeGradeClassResultV2(
		summary.GradeClassMonthlyResultMap(),
		client,
		config,
		"M",
		3+len(summary.GradeClassResultMap())+len(summary.GradeClassYearlyResultMap()),
		"クラス(月)",
	); err != nil {
		return err
	}

	if err = s.writeDistanceCategoryResultV2(
		summary.DistanceCategoryResultMap(),
		client,
		config,
		"R",
		1,
		"距離(全)",
	); err != nil {
		return err
	}

	if err = s.writeDistanceCategoryResultV2(
		summary.DistanceCategoryYearlyResultMap(),
		client,
		config,
		"R",
		2+len(summary.DistanceCategoryResultMap()),
		"距離(年)",
	); err != nil {
		return err
	}

	if err = s.writeDistanceCategoryResultV2(
		summary.DistanceCategoryMonthlyResultMap(),
		client,
		config,
		"R",
		3+len(summary.DistanceCategoryResultMap())+len(summary.DistanceCategoryYearlyResultMap()),
		"距離(月)",
	); err != nil {
		return err
	}

	if err = s.writeRaceCourseResultV2(
		summary.RaceCourseResultMap(),
		client,
		config,
		"W",
		1,
		"開催(全)",
	); err != nil {
		return err
	}

	if err = s.writeRaceCourseResultV2(
		summary.RaceCourseYearlyResultMap(),
		client,
		config,
		"W",
		2+len(summary.RaceCourseResultMap()),
		"開催(年)",
	); err != nil {
		return err
	}

	if err = s.writeRaceCourseResultV2(
		summary.RaceCourseMonthlyResultMap(),
		client,
		config,
		"W",
		3+len(summary.RaceCourseResultMap())+len(summary.RaceCourseYearlyResultMap()),
		"開催(月)",
	); err != nil {
		return err
	}

	s.logger.Infof("write summary v2 end")
	return nil
}

func (s *spreadSheetSummaryGateway) writeTotalResultV2(
	allResult *spreadsheet_entity.TicketResult,
	yearResult *spreadsheet_entity.TicketResult,
	monthResult *spreadsheet_entity.TicketResult,
	weekResult *spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
	rowCellNum int,
) error {
	s.logger.Infof("writing spreadsheet writeTotalResultV2")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), fmt.Sprintf("A%d", rowCellNum))
	values := [][]any{
		{
			"累計",
			"",
		},
		{
			"投資",
			allResult.Payment(),
		},
		{
			"回収",
			allResult.Payout(),
		},
		{
			"回収率",
			allResult.PayoutRate(),
		},
		{
			"年間累計",
			"",
		},
		{
			"投資",
			yearResult.Payment(),
		},
		{
			"回収",
			yearResult.Payout(),
		},
		{
			"回収率",
			yearResult.PayoutRate(),
		},
		{
			"収支",
			yearResult.Profit(),
		},
		{
			"月間累計",
			"",
		},
		{
			"投資",
			monthResult.Payment(),
		},
		{
			"回収",
			monthResult.Payout(),
		},
		{
			"回収率",
			monthResult.PayoutRate(),
		},
		{
			"収支",
			monthResult.Profit(),
		},
		{
			"週間累計",
			"",
		},
		{
			"投資",
			weekResult.Payment(),
		},
		{
			"回収",
			weekResult.Payout(),
		},
		{
			"回収率",
			weekResult.PayoutRate(),
		},
		{
			"収支",
			weekResult.Profit(),
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

func (s *spreadSheetSummaryGateway) writeWeeklyResultV2(
	results map[time.Time]*spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
	rowCellNum int,
) error {
	s.logger.Infof("writing spreadsheet writeWeeklyResultV2")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), fmt.Sprintf("C%d", rowCellNum))
	values := [][]any{
		{
			"週別",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	keys := make([]time.Time, 0, len(results))
	for k := range results {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	for _, key := range keys {
		result := results[key]
		values = append(values, []any{
			key.Format("2006/01/02"),
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

func (s *spreadSheetSummaryGateway) writeMonthlyResultV2(
	results map[time.Time]*spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
	rowCellNum int,
) error {
	s.logger.Infof("writing spreadsheet writeMonthlyResultV2")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), fmt.Sprintf("C%d", rowCellNum))
	values := [][]any{
		{
			"月別",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	keys := make([]time.Time, 0, len(results))
	for k := range results {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	for _, key := range keys {
		result := results[key]
		values = append(values, []any{
			"'" + key.Format("2006年01月"),
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

func (s *spreadSheetSummaryGateway) writeTicketResultV2(
	results map[types.TicketType]*spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
	colCell string,
	rowCell int,
	title string,
) error {
	cell := fmt.Sprintf("%s%d", colCell, rowCell)
	s.logger.Infof("writing spreadsheet writeTicketResultV2 %s", cell)
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), cell)
	values := [][]any{
		{
			title,
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
		values = append(values, []any{
			ticketType.Name(),
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

func (s *spreadSheetSummaryGateway) writeGradeClassResultV2(
	results map[types.GradeClass]*spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
	colCell string,
	rowCell int,
	title string,
) error {
	cell := fmt.Sprintf("%s%d", colCell, rowCell)
	s.logger.Infof("writing spreadsheet writeGradeClassResultV2 %s", cell)
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), cell)
	values := [][]any{
		{
			title,
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
		values = append(values, []any{
			gradeClass.String(),
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

func (s *spreadSheetSummaryGateway) writeDistanceCategoryResultV2(
	results map[types.DistanceCategory]*spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
	colCell string,
	rowCell int,
	title string,
) error {
	cell := fmt.Sprintf("%s%d", colCell, rowCell)
	s.logger.Infof("writing spreadsheet writeDistanceCategoryResultV2 %s", cell)
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), cell)
	values := [][]any{
		{
			title,
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
		values = append(values, []any{
			distanceCategory.String(),
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

func (s *spreadSheetSummaryGateway) writeRaceCourseResultV2(
	results map[types.RaceCourse]*spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
	colCell string,
	rowCell int,
	title string,
) error {
	cell := fmt.Sprintf("%s%d", colCell, rowCell)
	s.logger.Infof("writing spreadsheet writeRaceCourseResultV2 %s", cell)
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), cell)
	values := [][]any{
		{
			title,
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
		values = append(values, []any{
			raceCourse.Name(),
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

func (s *spreadSheetSummaryGateway) StyleV2(
	ctx context.Context,
	summary *spreadsheet_entity.Summary,
) error {
	client, config, err := s.spreadSheetConfigGateway.GetConfig(ctx, spreadSheetSummaryV2FileName)
	if err != nil {
		return err
	}

	s.logger.Infof("write spreadsheet style v2 start")
	if err = s.writeStyleAllResultV2(client, config); err != nil {
		return err
	}
	if err = s.writeStyleWeeklyResultV2(client, config, summary.WeeklyResults()); err != nil {
		return err
	}
	rowPosition := len(summary.WeeklyResults())
	if err = s.writeStyleMonthlyResultV2(client, config, summary.MonthlyResults(), rowPosition); err != nil {
		return err
	}
	if err = s.writeStyleTicketResultV2(client, config, summary.TicketResultMap(), 0); err != nil {
		return err
	}
	rowPosition = len(summary.TicketResultMap()) + 1
	if err = s.writeStyleTicketResultV2(client, config, summary.TicketYearlyResultMap(), rowPosition); err != nil {
		return err
	}
	rowPosition = len(summary.TicketResultMap()) + len(summary.TicketYearlyResultMap()) + 2
	if err = s.writeStyleTicketResultV2(client, config, summary.TicketYearlyResultMap(), rowPosition); err != nil {
		return err
	}
	if err = s.writeStyleGradeClassResultV2(client, config, summary.GradeClassResultMap(), 0); err != nil {
		return err
	}
	rowPosition = len(summary.GradeClassResultMap()) + 1
	if err = s.writeStyleGradeClassResultV2(client, config, summary.GradeClassYearlyResultMap(), rowPosition); err != nil {
		return err
	}
	rowPosition = len(summary.GradeClassResultMap()) + len(summary.GradeClassYearlyResultMap()) + 2
	if err = s.writeStyleGradeClassResultV2(client, config, summary.GradeClassMonthlyResultMap(), rowPosition); err != nil {
		return err
	}
	if err = s.writeStyleDistanceCategoryResultV2(client, config, summary.DistanceCategoryResultMap(), 0); err != nil {
		return err
	}
	rowPosition = len(summary.DistanceCategoryResultMap()) + 1
	if err = s.writeStyleDistanceCategoryResultV2(client, config, summary.DistanceCategoryYearlyResultMap(), rowPosition); err != nil {
		return err
	}
	rowPosition = len(summary.DistanceCategoryResultMap()) + len(summary.DistanceCategoryYearlyResultMap()) + 2
	if err = s.writeStyleDistanceCategoryResultV2(client, config, summary.DistanceCategoryMonthlyResultMap(), rowPosition); err != nil {
		return err
	}
	if err = s.writeStyleRaceCourseResultV2(client, config, summary.RaceCourseResultMap(), 0); err != nil {
		return err
	}
	rowPosition = len(summary.RaceCourseResultMap()) + 1
	if err = s.writeStyleRaceCourseResultV2(client, config, summary.RaceCourseYearlyResultMap(), rowPosition); err != nil {
		return err
	}
	rowPosition = len(summary.RaceCourseResultMap()) + len(summary.RaceCourseYearlyResultMap()) + 2
	if err = s.writeStyleRaceCourseResultV2(client, config, summary.RaceCourseMonthlyResultMap(), rowPosition); err != nil {
		return err
	}

	s.logger.Infof("write spreadsheet style v2 end")
	return nil
}

func (s *spreadSheetSummaryGateway) writeStyleAllResultV2(
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeStyleAllResultV2")
	requests := []*sheets.Request{
		s.createBackgroundColorRequest(config.SheetId(), 0, 1, 1, 4, 1.0, 0.937, 0.498),
		s.createBackgroundColorRequest(config.SheetId(), 0, 5, 1, 9, 1.0, 0.937, 0.498),
		s.createBackgroundColorRequest(config.SheetId(), 0, 10, 1, 14, 1.0, 0.937, 0.498),
		s.createBackgroundColorRequest(config.SheetId(), 0, 15, 1, 19, 1.0, 0.937, 0.498),
		s.createBoldTextRequest(config.SheetId(), 0, 0, 2, 19),
	}
	_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetSummaryGateway) writeStyleWeeklyResultV2(
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
	weeklyResults map[time.Time]*spreadsheet_entity.TicketResult,
) error {
	s.logger.Infof("writing spreadsheet writeStyleWeeklyResultV2")
	requests := []*sheets.Request{
		s.createBackgroundColorRequest(config.SheetId(), 2, 0, 7, 1, 1.0, 1.0, 0),
		s.createBackgroundColorRequest(config.SheetId(), 2, 1, 3, 1+len(weeklyResults), 1.0, 0.937, 0.498),
		s.createTextFormatRequest(config.SheetId(), 2, 0, 7, 1, "TEXT", true),
		s.createTextFormatRequest(config.SheetId(), 2, 1, 3, 1+len(weeklyResults), "DATE", true),
		s.createTextFormatRequest(config.SheetId(), 3, 1, 4, 1+len(weeklyResults), "PERCENT", false),
		s.createTextFormatRequest(config.SheetId(), 4, 1, 6, 1+len(weeklyResults), "TEXT", false),
		s.createTextFormatRequest(config.SheetId(), 6, 1, 7, 1+len(weeklyResults), "PERCENT", false),
	}
	_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetSummaryGateway) writeStyleMonthlyResultV2(
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
	monthlyResults map[time.Time]*spreadsheet_entity.TicketResult,
	rowPosition int,
) error {
	s.logger.Infof("writing spreadsheet writeStyleMonthlyResultV2")
	requests := []*sheets.Request{
		s.createBackgroundColorRequest(config.SheetId(), 2, 1+rowPosition, 7, 2+rowPosition, 1.0, 1.0, 0),
		s.createBackgroundColorRequest(config.SheetId(), 2, 2+rowPosition, 3, 2+rowPosition+len(monthlyResults), 1.0, 0.937, 0.498),
		s.createTextFormatRequest(config.SheetId(), 2, 1+rowPosition, 7, 2+rowPosition, "TEXT", true),
		s.createTextFormatRequest(config.SheetId(), 2, 2+rowPosition, 3, 2+rowPosition+len(monthlyResults), "TEXT", true),
		s.createTextFormatRequest(config.SheetId(), 3, 2+rowPosition, 4, 2+rowPosition+len(monthlyResults), "PERCENT", false),
		s.createTextFormatRequest(config.SheetId(), 4, 2+rowPosition, 6, 2+rowPosition+len(monthlyResults), "TEXT", false),
		s.createTextFormatRequest(config.SheetId(), 6, 2+rowPosition, 7, 2+rowPosition+len(monthlyResults), "PERCENT", false),
	}
	_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetSummaryGateway) writeStyleTicketResultV2(
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
	ticketResults map[types.TicketType]*spreadsheet_entity.TicketResult,
	rowPosition int,
) error {
	s.logger.Infof("writing spreadsheet writeStyleTicketResultV2 %d", rowPosition)
	requests := []*sheets.Request{
		s.createBackgroundColorRequest(config.SheetId(), 7, rowPosition, 12, 1+rowPosition, 1.0, 1.0, 0),
		s.createBackgroundColorRequest(config.SheetId(), 7, 1+rowPosition, 8, 1+rowPosition+len(ticketResults), 1.0, 0.937, 0.498),
		s.createTextFormatRequest(config.SheetId(), 7, rowPosition, 12, 1+rowPosition, "TEXT", true),
		s.createTextFormatRequest(config.SheetId(), 7, 1+rowPosition, 8, 1+rowPosition+len(ticketResults), "TEXT", true),
		s.createTextFormatRequest(config.SheetId(), 8, 1+rowPosition, 9, 1+rowPosition+len(ticketResults), "PERCENT", false),
		s.createTextFormatRequest(config.SheetId(), 9, 1+rowPosition, 11, 1+rowPosition+len(ticketResults), "TEXT", false),
		s.createTextFormatRequest(config.SheetId(), 11, 1+rowPosition, 12, 1+rowPosition+len(ticketResults), "PERCENT", false),
	}
	_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetSummaryGateway) writeStyleGradeClassResultV2(
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
	gradeClassResults map[types.GradeClass]*spreadsheet_entity.TicketResult,
	rowPosition int,
) error {
	s.logger.Infof("writing spreadsheet writeStyleGradeClassResultV2 %d", rowPosition)
	requests := []*sheets.Request{
		s.createBackgroundColorRequest(config.SheetId(), 12, rowPosition, 17, 1+rowPosition, 1.0, 1.0, 0),
		s.createBackgroundColorRequest(config.SheetId(), 12, 1+rowPosition, 13, 1+rowPosition+len(gradeClassResults), 1.0, 0.937, 0.498),
		s.createTextFormatRequest(config.SheetId(), 12, rowPosition, 17, 1+rowPosition, "TEXT", true),
		s.createTextFormatRequest(config.SheetId(), 12, 1+rowPosition, 13, 1+rowPosition+len(gradeClassResults), "TEXT", true),
		s.createTextFormatRequest(config.SheetId(), 13, 1+rowPosition, 14, 1+rowPosition+len(gradeClassResults), "PERCENT", false),
		s.createTextFormatRequest(config.SheetId(), 14, 1+rowPosition, 16, 1+rowPosition+len(gradeClassResults), "TEXT", false),
		s.createTextFormatRequest(config.SheetId(), 16, 1+rowPosition, 17, 1+rowPosition+len(gradeClassResults), "PERCENT", false),
	}
	_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetSummaryGateway) writeStyleDistanceCategoryResultV2(
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
	distanceCategoryResults map[types.DistanceCategory]*spreadsheet_entity.TicketResult,
	rowPosition int,
) error {
	s.logger.Infof("writing spreadsheet writeStyleDistanceCategoryResultV2 %d", rowPosition)
	requests := []*sheets.Request{
		s.createBackgroundColorRequest(config.SheetId(), 17, rowPosition, 22, 1+rowPosition, 1.0, 1.0, 0),
		s.createBackgroundColorRequest(config.SheetId(), 17, 1+rowPosition, 18, 1+rowPosition+len(distanceCategoryResults), 1.0, 0.937, 0.498),
		s.createTextFormatRequest(config.SheetId(), 17, rowPosition, 22, 1+rowPosition, "TEXT", true),
		s.createTextFormatRequest(config.SheetId(), 17, 1+rowPosition, 18, 1+rowPosition+len(distanceCategoryResults), "TEXT", true),
		s.createTextFormatRequest(config.SheetId(), 18, 1+rowPosition, 19, 1+rowPosition+len(distanceCategoryResults), "PERCENT", false),
		s.createTextFormatRequest(config.SheetId(), 19, 1+rowPosition, 21, 1+rowPosition+len(distanceCategoryResults), "TEXT", false),
		s.createTextFormatRequest(config.SheetId(), 21, 1+rowPosition, 22, 1+rowPosition+len(distanceCategoryResults), "PERCENT", false),
	}
	_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetSummaryGateway) writeStyleRaceCourseResultV2(
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
	raceCourseResults map[types.RaceCourse]*spreadsheet_entity.TicketResult,
	rowPosition int,
) error {
	s.logger.Infof("writing spreadsheet writeStyleRaceCourseResultV2 %d", rowPosition)
	requests := []*sheets.Request{
		s.createBackgroundColorRequest(config.SheetId(), 22, rowPosition, 27, 1+rowPosition, 1.0, 1.0, 0),
		s.createBackgroundColorRequest(config.SheetId(), 22, 1+rowPosition, 23, 1+rowPosition+len(raceCourseResults), 1.0, 0.937, 0.498),
		s.createTextFormatRequest(config.SheetId(), 22, rowPosition, 27, 1+rowPosition, "TEXT", true),
		s.createTextFormatRequest(config.SheetId(), 22, 1+rowPosition, 23, 1+rowPosition+len(raceCourseResults), "TEXT", true),
		s.createTextFormatRequest(config.SheetId(), 23, 1+rowPosition, 24, 1+rowPosition+len(raceCourseResults), "PERCENT", false),
		s.createTextFormatRequest(config.SheetId(), 24, 1+rowPosition, 26, 1+rowPosition+len(raceCourseResults), "TEXT", false),
		s.createTextFormatRequest(config.SheetId(), 26, 1+rowPosition, 27, 1+rowPosition+len(raceCourseResults), "PERCENT", false),
	}
	_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetSummaryGateway) createTextFormatRequest(
	sheetId int64,
	startCol, startRow, endCol, endRow int,
	formatType string,
	bold bool,
) *sheets.Request {
	cellFormat := &sheets.CellFormat{}
	switch formatType {
	case "DATE":
		cellFormat.NumberFormat = &sheets.NumberFormat{
			Type:    "DATE",
			Pattern: "yyyy/MM/dd",
		}
		cellFormat.TextFormat = &sheets.TextFormat{
			Bold: bold,
		}
	case "PERCENT":
		cellFormat.NumberFormat = &sheets.NumberFormat{
			Type:    "PERCENT",
			Pattern: "0.00%",
		}
		cellFormat.TextFormat = &sheets.TextFormat{
			Bold: bold,
		}
	case "TEXT":
		cellFormat.TextFormat = &sheets.TextFormat{
			FontSize: 10,
			Bold:     bold,
		}
	}

	return &sheets.Request{
		RepeatCell: &sheets.RepeatCellRequest{
			Fields: "userEnteredFormat.textFormat,userEnteredFormat.numberFormat",
			Range: &sheets.GridRange{
				SheetId:          sheetId,
				StartColumnIndex: int64(startCol),
				StartRowIndex:    int64(startRow),
				EndColumnIndex:   int64(endCol),
				EndRowIndex:      int64(endRow),
			},
			Cell: &sheets.CellData{
				UserEnteredFormat: cellFormat,
			},
		},
	}
}

func (s *spreadSheetSummaryGateway) createBackgroundColorRequest(
	sheetId int64,
	startCol, startRow, endCol, endRow int,
	red, green, blue float64,
) *sheets.Request {
	cellFormat := &sheets.CellFormat{
		BackgroundColor: &sheets.Color{
			Red:   red,
			Green: green,
			Blue:  blue,
		},
	}

	return &sheets.Request{
		RepeatCell: &sheets.RepeatCellRequest{
			Fields: "userEnteredFormat.backgroundColor,userEnteredFormat.numberFormat,userEnteredFormat.textFormat",
			Range: &sheets.GridRange{
				SheetId:          sheetId,
				StartColumnIndex: int64(startCol),
				StartRowIndex:    int64(startRow),
				EndColumnIndex:   int64(endCol),
				EndRowIndex:      int64(endRow),
			},
			Cell: &sheets.CellData{
				UserEnteredFormat: cellFormat,
			},
		},
	}
}

func (s *spreadSheetSummaryGateway) createBoldTextRequest(
	sheetId int64,
	startCol, startRow, endCol, endRow int,
) *sheets.Request {
	return &sheets.Request{
		RepeatCell: &sheets.RepeatCellRequest{
			Fields: "userEnteredFormat.textFormat.bold",
			Range: &sheets.GridRange{
				SheetId:          sheetId,
				StartColumnIndex: int64(startCol),
				StartRowIndex:    int64(startRow),
				EndColumnIndex:   int64(endCol),
				EndRowIndex:      int64(endRow),
			},
			Cell: &sheets.CellData{
				UserEnteredFormat: &sheets.CellFormat{
					TextFormat: &sheets.TextFormat{
						Bold: true,
					},
				},
			},
		},
	}
}

func (s *spreadSheetSummaryGateway) ClearV2(ctx context.Context) error {
	client, config, err := s.spreadSheetConfigGateway.GetConfig(ctx, spreadSheetSummaryV2FileName)
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
					EndColumnIndex:   50,
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
