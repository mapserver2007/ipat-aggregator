package gateway

import (
	"context"
	"fmt"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"
)

const (
	spreadSheetAnalysisRaceTimeFileName = "spreadsheet_analysis_race_time.json"
)

type SpreadSheetAnalysisRaceTimeGateway interface {
	Write(ctx context.Context,
		analysisRaceTimeMap map[filter.AttributeId]*spreadsheet_entity.AnalysisRaceTime,
		attributeFilters []filter.AttributeId,
		conditionFilters []filter.AttributeId,
	) error
	Style(ctx context.Context,
		analysisRaceTimeMap map[filter.AttributeId]*spreadsheet_entity.AnalysisRaceTime,
		attributeFilters []filter.AttributeId,
		conditionFilters []filter.AttributeId,
	) error
	Clear(ctx context.Context) error
}

type spreadSheetAnalysisRaceTimeGateway struct {
	spreadSheetConfigGateway SpreadSheetConfigGateway
	logger                   *logrus.Logger
}

func NewSpreadSheetAnalysisRaceTimeGateway(
	spreadSheetConfigGateway SpreadSheetConfigGateway,
	logger *logrus.Logger,
) SpreadSheetAnalysisRaceTimeGateway {
	return &spreadSheetAnalysisRaceTimeGateway{
		spreadSheetConfigGateway: spreadSheetConfigGateway,
		logger:                   logger,
	}
}

func (s *spreadSheetAnalysisRaceTimeGateway) Write(
	ctx context.Context,
	analysisRaceTimeMap map[filter.AttributeId]*spreadsheet_entity.AnalysisRaceTime,
	attributeFilters []filter.AttributeId,
	conditionFilters []filter.AttributeId,
) error {
	client, config, err := s.spreadSheetConfigGateway.GetConfig(ctx, spreadSheetAnalysisRaceTimeFileName)
	if err != nil {
		return err
	}

	s.logger.Infof("write analysis race time start")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "A1")
	values := [][]any{
		{
			"レース",
			"場所",
			"距離",
			"クラス",
			"馬場",
			"馬齢",
			"タイム",
			"前3f",
			"前4f",
			"5f通過",
			"後3f",
			"後4f",
			"馬場(平均)",
			"馬場(最遅)",
			"馬場(最速)",
			"タイム指数",
		},
	}

	for _, analysisFilter := range attributeFilters {
		analysisRaceTime, ok := analysisRaceTimeMap[analysisFilter]
		if !ok {
			continue
		}
		var (
			raceCourseFilter       filter.AttributeId
			courseCategoryFilter   filter.AttributeId
			distanceFilter         filter.AttributeId
			trackConditionFilter   filter.AttributeId
			classFilter            filter.AttributeId
			raceAgeConditionFilter filter.AttributeId
		)

		for _, originFilter := range analysisFilter.OriginFilters() {
			if originFilter&conditionFilters[0] != 0 {
				raceCourseFilter = originFilter
			}
			if originFilter&conditionFilters[1] != 0 {
				courseCategoryFilter = originFilter
			}
			if originFilter&conditionFilters[2] != 0 {
				distanceFilter = originFilter
			}
			if originFilter&conditionFilters[3] != 0 {
				trackConditionFilter = originFilter
			}
			if originFilter&conditionFilters[4] != 0 {
				classFilter = originFilter
			}
			if originFilter&conditionFilters[5] != 0 {
				raceAgeConditionFilter = originFilter
			}
		}

		values = append(values, []any{
			analysisRaceTime.RaceCount(),
			raceCourseFilter.String(),
			courseCategoryFilter.String() + distanceFilter.String(),
			classFilter.String(),
			trackConditionFilter.String(),
			raceAgeConditionFilter.String(),
			analysisRaceTime.AverageRaceTime(),
			analysisRaceTime.AverageFirst3f(),
			analysisRaceTime.AverageFirst4f(),
			analysisRaceTime.AverageRap5f(),
			analysisRaceTime.AverageLast3f(),
			analysisRaceTime.AverageLast4f(),
			analysisRaceTime.AverageTrackIndex(),
			analysisRaceTime.MaxTrackIndex(),
			analysisRaceTime.MinTrackIndex(),
			analysisRaceTime.AverageTimeIndex(),
		})
	}

	_, err = client.Spreadsheets.Values.Update(config.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	s.logger.Infof("write analysis race time end")

	return nil
}

func (s *spreadSheetAnalysisRaceTimeGateway) Style(
	ctx context.Context,
	analysisRaceTimeMap map[filter.AttributeId]*spreadsheet_entity.AnalysisRaceTime,
	attributeFilters []filter.AttributeId,
	conditionFilters []filter.AttributeId,
) error {
	client, config, err := s.spreadSheetConfigGateway.GetConfig(ctx, spreadSheetAnalysisRaceTimeFileName)
	if err != nil {
		return err
	}

	s.logger.Infof("write analysis race time style start")
	requests := make([]*sheets.Request, 0)
	requests = append(requests, s.createBackgroundColorRequest(
		config.SheetId(),
		0, 0, 6, 1,
		1.0, 1.0, 0.0,
	))
	requests = append(requests, s.createBackgroundColorRequest(
		config.SheetId(),
		6, 0, 16, 1,
		1.0, 0.0, 0.0,
	))
	requests = append(requests, s.createTextFormatRequest(
		config.SheetId(),
		0, 0, 6, 1,
		0.0, 0.0, 0.0,
	))
	requests = append(requests, s.createTextFormatRequest(
		config.SheetId(),
		6, 0, 16, 1,
		1.0, 1.0, 1.0,
	))
	requests = append(requests, s.createTextBoldRequest(
		config.SheetId(),
		0, 0, 6, 1,
		true,
	))
	requests = append(requests, s.createTextBoldRequest(
		config.SheetId(),
		6, 0, 16, 1,
		true,
	))

	_, err = client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	s.logger.Infof("write analysis race time style end")

	return nil
}

func (s *spreadSheetAnalysisRaceTimeGateway) Clear(ctx context.Context) error {
	client, config, err := s.spreadSheetConfigGateway.GetConfig(ctx, spreadSheetAnalysisRaceTimeFileName)
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

func (s *spreadSheetAnalysisRaceTimeGateway) createTextFormatRequest(
	sheetId int64,
	startCol, startRow, endCol, endRow int,
	red, green, blue float64,
) *sheets.Request {
	return &sheets.Request{
		RepeatCell: &sheets.RepeatCellRequest{
			Fields: "userEnteredFormat.textFormat.foregroundColor",
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
						ForegroundColor: &sheets.Color{
							Red:   red,
							Green: green,
							Blue:  blue,
						},
					},
				},
			},
		},
	}
}

func (s *spreadSheetAnalysisRaceTimeGateway) createTextBoldRequest(
	sheetId int64,
	startCol, startRow, endCol, endRow int,
	bold bool,
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
						Bold: bold,
					},
				},
			},
		},
	}
}

func (s *spreadSheetAnalysisRaceTimeGateway) createBackgroundColorRequest(
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
