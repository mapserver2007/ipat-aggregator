package gateway

import (
	"context"
	"fmt"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"
)

const (
	spreadSheetAnalysisPlaceUnhitFileName = "spreadsheet_analysis_place_unhit.json"
)

type SpreadSheetAnalysisPlaceUnhitGateway interface {
	Write(ctx context.Context, analysisPlaceUnhits []*spreadsheet_entity.AnalysisPlaceUnhit) error
	Style(ctx context.Context, analysisPlaceUnhits []*spreadsheet_entity.AnalysisPlaceUnhit) error
}

type spreadSheetAnalysisPlaceUnhitGateway struct {
	spreadSheetConfigGateway SpreadSheetConfigGateway
	logger                   *logrus.Logger
}

func NewSpreadSheetAnalysisPlaceUnhitGateway(
	spreadSheetConfigGateway SpreadSheetConfigGateway,
	logger *logrus.Logger,
) SpreadSheetAnalysisPlaceUnhitGateway {
	return &spreadSheetAnalysisPlaceUnhitGateway{
		spreadSheetConfigGateway: spreadSheetConfigGateway,
		logger:                   logger,
	}
}

func (g *spreadSheetAnalysisPlaceUnhitGateway) Write(
	ctx context.Context,
	analysisPlaceUnhits []*spreadsheet_entity.AnalysisPlaceUnhit,
) error {
	client, config, err := g.spreadSheetConfigGateway.GetConfig(ctx, spreadSheetAnalysisPlaceUnhitFileName)
	if err != nil {
		return err
	}

	g.logger.Infof("write analysis place unhit start")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "A1")
	values := [][]any{
		{
			"レース条件",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"頭数",
			"馬番",
			"馬名",
			"騎手",
			"馬体重",
			"増減",
			"印",
			"人気",
			"オッズ",
			"着",
			"100人気",
			"赤単数",
			"断層1",
			"断層2",
			"馬連続",
			"連平均",
			"軸出現",
		},
	}

	for _, analysisPlaceUnhit := range analysisPlaceUnhits {
		values = append(values, []any{
			analysisPlaceUnhit.RaceDate(),
			analysisPlaceUnhit.Class().String(),
			fmt.Sprintf("=HYPERLINK(\"%s\",\"%s\")", analysisPlaceUnhit.RaceUrl(), analysisPlaceUnhit.RaceName()),
			analysisPlaceUnhit.CourseCategory().String(),
			analysisPlaceUnhit.Distance(),
			analysisPlaceUnhit.RaceWeightCondition().String(),
			analysisPlaceUnhit.JockeyWeight(),
			analysisPlaceUnhit.TrackCondition().String(),
			analysisPlaceUnhit.Entries(),
			analysisPlaceUnhit.HorseNumber(),
			fmt.Sprintf("=HYPERLINK(\"%s\",\"%s\")", analysisPlaceUnhit.HorseUrl(), analysisPlaceUnhit.HorseName()),
			fmt.Sprintf("=HYPERLINK(\"%s\",\"%s\")", analysisPlaceUnhit.JockeyUrl(), analysisPlaceUnhit.JockeyName()),
			analysisPlaceUnhit.HorseWeight(),
			analysisPlaceUnhit.HorseWeightAdd(),
			analysisPlaceUnhit.Marker().String(),
			analysisPlaceUnhit.PopularNumber(),
			analysisPlaceUnhit.Odds().String(),
			func() string {
				if analysisPlaceUnhit.OrderNo() == 99 {
					return "中止"
				}
				return fmt.Sprintf("%d", analysisPlaceUnhit.OrderNo())
			}(),
			analysisPlaceUnhit.TrioOdds100().String(),
			analysisPlaceUnhit.WinRedOddsNum(),
			analysisPlaceUnhit.OddsFault1().String(),
			analysisPlaceUnhit.OddsFault2().String(),
			analysisPlaceUnhit.QuinellaConsecutiveNumber(),
			analysisPlaceUnhit.QuinellaWheelAverageOdds().Round(1).String(),
			analysisPlaceUnhit.TrioFavoriteCount(),
		})
	}

	_, err = client.Spreadsheets.Values.Update(config.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	g.logger.Infof("write analysis place unhit end")

	return nil
}

func (g *spreadSheetAnalysisPlaceUnhitGateway) Style(
	ctx context.Context,
	analysisPlaceUnhits []*spreadsheet_entity.AnalysisPlaceUnhit,
) error {
	client, config, err := g.spreadSheetConfigGateway.GetConfig(ctx, spreadSheetAnalysisPlaceUnhitFileName)
	if err != nil {
		return err
	}

	g.logger.Infof("write analysis place unhit style start")
	err = g.writeStyleTrioOdds100(client, config, analysisPlaceUnhits)
	if err != nil {
		return err
	}

	err = g.writeStyleWinRedOddsNum(client, config, analysisPlaceUnhits)
	if err != nil {
		return err
	}

	err = g.writeStyleOddsFault(client, config, analysisPlaceUnhits)
	if err != nil {
		return err
	}

	g.logger.Infof("write analysis place unhit style end")

	return nil
}

func (g *spreadSheetAnalysisPlaceUnhitGateway) writeStyleTrioOdds100(
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
	analysisPlaceUnhits []*spreadsheet_entity.AnalysisPlaceUnhit,
) error {
	g.logger.Infof("writing spreadsheet writeStyleTrioOdds100")

	requests := make([]*sheets.Request, 0, len(analysisPlaceUnhits))
	for idx, analysisPlaceUnhit := range analysisPlaceUnhits {
		rowNum := 1 + idx
		if analysisPlaceUnhit.TrioOdds100().Equal(decimal.Zero) {
			continue
		}
		if analysisPlaceUnhit.TrioOdds100().LessThanOrEqual(decimal.NewFromInt(1350)) && analysisPlaceUnhit.TrioOdds100().GreaterThanOrEqual(decimal.NewFromInt(650)) {
			requests = append(requests, g.createBackgroundColorRequest(
				config.SheetId(),
				18, rowNum, 19, rowNum+1,
				1.0, 1.0, 0.0,
			))
		} else {
			requests = append(requests, g.createBackgroundColorRequest(
				config.SheetId(),
				18, rowNum, 19, rowNum+1,
				1.0, 0.0, 0.0,
			))
			requests = append(requests, g.createTextFormatRequest(
				config.SheetId(),
				18, rowNum, 19, rowNum+1,
				1.0, 1.0, 1.0,
			))
		}
	}

	_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (g *spreadSheetAnalysisPlaceUnhitGateway) writeStyleWinRedOddsNum(
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
	analysisPlaceUnhits []*spreadsheet_entity.AnalysisPlaceUnhit,
) error {
	g.logger.Infof("writing spreadsheet writeStyleWinRedOddsNum")

	requests := make([]*sheets.Request, 0, len(analysisPlaceUnhits))
	for idx, analysisPlaceUnhit := range analysisPlaceUnhits {
		rowNum := 1 + idx
		if analysisPlaceUnhit.WinRedOddsNum() > 3 {
			requests = append(requests, g.createBackgroundColorRequest(
				config.SheetId(),
				19, rowNum, 20, rowNum+1,
				1.0, 0.0, 0.0,
			))
			requests = append(requests, g.createTextFormatRequest(
				config.SheetId(),
				19, rowNum, 20, rowNum+1,
				1.0, 1.0, 1.0,
			))
		} else {
			requests = append(requests, g.createBackgroundColorRequest(
				config.SheetId(),
				19, rowNum, 20, rowNum+1,
				1.0, 1.0, 0.0,
			))
		}
	}

	_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (g *spreadSheetAnalysisPlaceUnhitGateway) writeStyleOddsFault(
	client *sheets.Service,
	config *spreadsheet_entity.SpreadSheetConfig,
	analysisPlaceUnhits []*spreadsheet_entity.AnalysisPlaceUnhit,
) error {
	g.logger.Infof("writing spreadsheet writeStyleOddsFault")

	requests := make([]*sheets.Request, 0, len(analysisPlaceUnhits))
	for idx, analysisPlaceUnhit := range analysisPlaceUnhits {
		rowNum := 1 + idx
		if analysisPlaceUnhit.OddsFault1().LessThan(decimal.NewFromFloat(3.0)) {
			requests = append(requests, g.createBackgroundColorRequest(
				config.SheetId(),
				19, rowNum, 20, rowNum+1,
				1.0, 0.0, 0.0,
			))
			requests = append(requests, g.createTextFormatRequest(
				config.SheetId(),
				19, rowNum, 20, rowNum+1,
				1.0, 1.0, 1.0,
			))
		} else {
			requests = append(requests, g.createBackgroundColorRequest(
				config.SheetId(),
				19, rowNum, 20, rowNum+1,
				1.0, 1.0, 0.0,
			))
			requests = append(requests, g.createTextFormatRequest(
				config.SheetId(),
				19, rowNum, 20, rowNum+1,
				0.0, 0.0, 0.0,
			))
		}
		if analysisPlaceUnhit.OddsFault2().LessThan(decimal.NewFromFloat(3.0)) {
			requests = append(requests, g.createBackgroundColorRequest(
				config.SheetId(),
				20, rowNum, 21, rowNum+1,
				1.0, 0.0, 0.0,
			))
			requests = append(requests, g.createTextFormatRequest(
				config.SheetId(),
				20, rowNum, 21, rowNum+1,
				1.0, 1.0, 1.0,
			))
		} else {
			requests = append(requests, g.createBackgroundColorRequest(
				config.SheetId(),
				20, rowNum, 21, rowNum+1,
				1.0, 1.0, 0.0,
			))
			requests = append(requests, g.createTextFormatRequest(
				config.SheetId(),
				20, rowNum, 21, rowNum+1,
				0.0, 0.0, 0.0,
			))
		}
	}

	_, err := client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (g *spreadSheetAnalysisPlaceUnhitGateway) createTextFormatRequest(
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

func (g *spreadSheetAnalysisPlaceUnhitGateway) createBackgroundColorRequest(
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
