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
	Write(ctx context.Context, summary *spreadsheet_entity.Summary) error
	WriteV2(ctx context.Context, summary *spreadsheet_entity.Summary) error
	Style(ctx context.Context, summary *spreadsheet_entity.Summary) error
	StyleV2(ctx context.Context, summary *spreadsheet_entity.Summary) error
	Clear(ctx context.Context) error
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

func (s *spreadSheetSummaryGateway) Write(
	ctx context.Context,
	summary *spreadsheet_entity.Summary,
) error {
	client, config, err := s.spreadSheetConfigGateway.GetConfig(ctx, spreadSheetSummaryFileName)
	if err != nil {
		return err
	}

	s.logger.Infof("write summary start")
	err = s.writeAllResult(summary.AllTermResult(), client, config)
	if err != nil {
		return err
	}
	err = s.writeYearResult(summary.YearTermResult(), client, config)
	if err != nil {
		return err
	}
	err = s.writeMonthResult(summary.MonthTermResult(), client, config)
	if err != nil {
		return err
	}
	err = s.writeTicketResult(summary.TicketResultMap(), client, config)
	if err != nil {
		return err
	}
	err = s.writeGradeClassResult(summary.GradeClassResultMap(), client, config)
	if err != nil {
		return err
	}
	err = s.writeMonthlyResult(summary.MonthlyResults(), client, config)
	if err != nil {
		return err
	}
	err = s.writeCourseCategoryResult(summary.CourseCategoryResultMap(), client, config)
	if err != nil {
		return err
	}
	err = s.writeDistanceCategoryResult(summary.DistanceCategoryResultMap(), client, config)
	if err != nil {
		return err
	}
	err = s.writeRaceCourseResult(summary.RaceCourseResultMap(), client, config)
	if err != nil {
		return err
	}

	s.logger.Infof("write summary end")
	return nil
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

func (s *spreadSheetSummaryGateway) writeAllResult(
	result *spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeAllResult")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "A1")
	values := [][]any{
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
	result *spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeYearResult")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "E1")
	values := [][]any{
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
	result *spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeMonthResult")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "C1")
	values := [][]any{
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
	results map[types.TicketType]*spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeTicketResult")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "A6")
	values := [][]any{
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
		values = append(values, []any{
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
	results map[types.GradeClass]*spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeGradeClassResult")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "A15")
	values := [][]any{
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
		values = append(values, []any{
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
	results map[time.Time]*spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeMonthlyResult")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "A28")
	values := [][]any{
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

	keys := make([]time.Time, 0, len(results))
	for month := range results {
		keys = append(keys, month)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	for _, month := range keys {
		result := results[month]
		values = append(values, []any{
			month.Format("2006年01月"),
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
	results map[types.CourseCategory]*spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeCourseCategoryResult")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "I6")
	values := [][]any{
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
		values = append(values, []any{
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
	results map[types.DistanceCategory]*spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeDistanceCategoryResult")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "I10")
	values := [][]any{
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
		values = append(values, []any{
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
	results map[types.RaceCourse]*spreadsheet_entity.TicketResult,
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeRaceCourseResult")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "I21")
	values := [][]any{
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
		values = append(values, []any{
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
	client, config, err := s.spreadSheetConfigGateway.GetConfig(ctx, spreadSheetSummaryFileName)
	if err != nil {
		return err
	}

	s.logger.Infof("write spreadsheet style start")
	err = s.writeStyleAllResult(client, config)
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

func (s *spreadSheetSummaryGateway) writeStyleAllResult(
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
						EndRowIndex:      3,
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

func (s *spreadSheetSummaryGateway) writeStyleAllResultV2(
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
) error {
	s.logger.Infof("writing spreadsheet writeStyleAllResultV2")
	requests := []*sheets.Request{
		s.createBackgroundColorRequest(config.SheetId(), 0, 1, 1, 4, 1.0, 0.937, 0.498),
		s.createBackgroundColorRequest(config.SheetId(), 0, 5, 1, 8, 1.0, 0.937, 0.498),
		s.createBackgroundColorRequest(config.SheetId(), 0, 9, 1, 12, 1.0, 0.937, 0.498),
		s.createBackgroundColorRequest(config.SheetId(), 0, 13, 1, 16, 1.0, 0.937, 0.498),
		s.createBoldTextRequest(config.SheetId(), 0, 0, 2, 16),
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
	client, config, err := s.spreadSheetConfigGateway.GetConfig(ctx, spreadSheetSummaryFileName)
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
