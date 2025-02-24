package gateway

import (
	"context"
	"fmt"
	"strings"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"
)

const (
	spreadSheetListFileName = "spreadsheet_list.json"
)

type SpreadSheetListGateway interface {
	Write(ctx context.Context, rows []*spreadsheet_entity.ListRow) error
	Style(ctx context.Context, rows []*spreadsheet_entity.ListRow) error
	Clear(ctx context.Context) error
}

type spreadSheetListGateway struct {
	spreadSheetConfigGateway SpreadSheetConfigGateway
	logger                   *logrus.Logger
}

func NewSpreadSheetListGateway(
	logger *logrus.Logger,
	spreadSheetConfigGateway SpreadSheetConfigGateway,
) SpreadSheetListGateway {
	return &spreadSheetListGateway{
		spreadSheetConfigGateway: spreadSheetConfigGateway,
		logger:                   logger,
	}
}

func (s *spreadSheetListGateway) Write(
	ctx context.Context,
	rows []*spreadsheet_entity.ListRow,
) error {
	client, config, err := s.spreadSheetConfigGateway.GetConfig(ctx, spreadSheetListFileName)
	if err != nil {
		return err
	}

	s.logger.Infof("write list start")

	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "A1")
	values := [][]any{
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

	for _, row := range rows {
		values = append(values, []any{
			row.Data().RaceDate(),
			row.Data().Class(),
			row.Data().CourseCategory(),
			row.Data().Distance(),
			row.Data().TraceCondition(),
			fmt.Sprintf("=HYPERLINK(\"%s\",\"%s\")", row.Data().Url(), row.Data().RaceName()),
			row.Data().Payment(),
			row.Data().Payout(),
			row.Data().PayoutRate(),
			row.Data().FavoriteHorse(),
			row.Data().FavoriteJockey(),
			row.Data().FavoriteHorsePopular(),
			row.Data().FavoriteHorseOdds(),
			row.Data().RivalHorse(),
			row.Data().RivalJockey(),
			func(s string) string {
				if s == "0" {
					return "-"
				}
				return s
			}(row.Data().RivalHorsePopular()),
			row.Data().RivalHorseOdds(),
			row.Data().FirstPlaceHorse(),
			row.Data().FirstPlaceJockey(),
			row.Data().FirstPlaceHorsePopular(),
			row.Data().FirstPlaceHorseOdds(),
			row.Data().SecondPlaceHorse(),
			row.Data().SecondPlaceJockey(),
			row.Data().SecondPlaceHorsePopular(),
			row.Data().SecondPlaceHorseOdds(),
		})
	}

	_, err = client.Spreadsheets.Values.Update(config.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	s.logger.Infof("write list end")

	return nil
}

func (s *spreadSheetListGateway) Style(
	ctx context.Context,
	rows []*spreadsheet_entity.ListRow,
) error {
	client, config, err := s.spreadSheetConfigGateway.GetConfig(ctx, spreadSheetListFileName)
	if err != nil {
		return err
	}

	s.logger.Infof("write list style start")

	var requests []*sheets.Request
	requests = append(requests, &sheets.Request{
		RepeatCell: &sheets.RepeatCellRequest{
			Fields: "userEnteredFormat.backgroundColor",
			Range: &sheets.GridRange{
				SheetId:          config.SheetId(),
				StartColumnIndex: 0,
				StartRowIndex:    0,
				EndColumnIndex:   25,
				EndRowIndex:      1,
			},
			Cell: &sheets.CellData{
				UserEnteredFormat: &sheets.CellFormat{
					BackgroundColor: &sheets.Color{
						Red:   1.0,
						Green: 1.0,
						Blue:  0.0,
					},
				},
			},
		},
	})

	for idx, row := range rows {
		rowNo := int64(idx + 1)
		requests = append(requests, []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 1,
						StartRowIndex:    rowNo,
						EndColumnIndex:   2,
						EndRowIndex:      rowNo + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: s.getCellColor(row.Style().ClassColor()),
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 9,
						StartRowIndex:    rowNo,
						EndColumnIndex:   10,
						EndRowIndex:      rowNo + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: s.getCellColor(row.Style().FavoriteHorseColor()),
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
						StartRowIndex:    rowNo,
						EndColumnIndex:   14,
						EndRowIndex:      rowNo + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: s.getCellColor(row.Style().RivalHorseColor()),
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 17,
						StartRowIndex:    rowNo,
						EndColumnIndex:   18,
						EndRowIndex:      rowNo + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: s.getCellColor(row.Style().FirstPlaceHorseColor()),
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 21,
						StartRowIndex:    rowNo,
						EndColumnIndex:   22,
						EndRowIndex:      rowNo + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: s.getCellColor(row.Style().SecondPlaceHorseColor()),
						},
					},
				},
			},
		}...)

		if len(row.Style().PayoutComments()) > 0 {
			requests = append(requests, &sheets.Request{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "note",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 7,
						StartRowIndex:    rowNo,
						EndColumnIndex:   8,
						EndRowIndex:      rowNo + 1,
					},
					Cell: &sheets.CellData{
						Note: strings.Join(row.Style().PayoutComments(), "\n"),
					},
				},
			})
		}
	}

	_, err = client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		return err
	}

	s.logger.Infof("write list style end")

	return nil
}

func (s *spreadSheetListGateway) Clear(ctx context.Context) error {
	client, config, err := s.spreadSheetConfigGateway.GetConfig(ctx, spreadSheetListFileName)
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
					EndColumnIndex:   40,
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

func (s *spreadSheetListGateway) getCellColor(
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
