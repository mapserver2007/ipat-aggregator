package infrastructure

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"google.golang.org/api/sheets/v4"
	"log"
	"strings"
)

const (
	spreadSheetListFileName2 = "spreadsheet_list2.json"
)

type spreadSheetListRepository struct {
	client             *sheets.Service
	spreadSheetConfig  *spreadsheet_entity.SpreadSheetConfig
	spreadSheetService service.SpreadSheetService
}

func NewSpreadSheetListRepository(
	spreadSheetService service.SpreadSheetService,
) (repository.SpreadsheetListRepository, error) {
	ctx := context.Background()
	client, spreadSheetConfig, err := getSpreadSheetConfig(ctx, spreadSheetListFileName2)
	if err != nil {
		return nil, err
	}

	return &spreadSheetListRepository{
		client:             client,
		spreadSheetConfig:  spreadSheetConfig,
		spreadSheetService: spreadSheetService,
	}, nil
}

func (s *spreadSheetListRepository) Write(ctx context.Context, rows []*spreadsheet_entity.Row) error {
	log.Println(ctx, fmt.Sprintf("write list start"))
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName(), "A1")
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

	for _, row := range rows {
		values = append(values, []interface{}{
			row.RaceDate(),
			row.Class(),
			row.CourseCategory(),
			row.Distance(),
			row.TraceCondition(),
			fmt.Sprintf("=HYPERLINK(\"%s\",\"%s\")", row.Url(), row.RaceName()),
			row.Payment(),
			row.Payout(),
			row.PayoutRate(),
			row.FavoriteHorse(),
			row.FavoriteJockey(),
			row.FavoriteHorsePopular(),
			row.FavoriteHorseOdds(),
			row.RivalHorse(),
			row.RivalJockey(),
			row.RivalHorsePopular(),
			row.RivalHorseOdds(),
			row.FirstPlaceHorse(),
			row.FirstPlaceJockey(),
			row.FirstPlaceHorsePopular(),
			row.FirstPlaceHorseOdds(),
			row.SecondPlaceHorse(),
			row.SecondPlaceJockey(),
			row.SecondPlaceHorsePopular(),
			row.SecondPlaceHorseOdds(),
		})
	}

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	log.Println(ctx, fmt.Sprintf("write list end"))

	return nil
}

func (s *spreadSheetListRepository) Style(ctx context.Context, styles []*spreadsheet_entity.Style) error {
	log.Println(ctx, fmt.Sprintf("write list style start"))

	var requests []*sheets.Request
	requests = append(requests, &sheets.Request{
		RepeatCell: &sheets.RepeatCellRequest{
			Fields: "userEnteredFormat.backgroundColor",
			Range: &sheets.GridRange{
				SheetId:          s.spreadSheetConfig.SheetId(),
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

	for idx, style := range styles {
		rowNo := int64(idx + 1)
		requests = append(requests, []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.spreadSheetConfig.SheetId(),
						StartColumnIndex: 1,
						StartRowIndex:    rowNo,
						EndColumnIndex:   2,
						EndRowIndex:      rowNo + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: s.spreadSheetService.GetCellColor(ctx, style.ClassColor()),
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.spreadSheetConfig.SheetId(),
						StartColumnIndex: 9,
						StartRowIndex:    rowNo,
						EndColumnIndex:   10,
						EndRowIndex:      rowNo + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: s.spreadSheetService.GetCellColor(ctx, style.FavoriteHorseColor()),
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.spreadSheetConfig.SheetId(),
						StartColumnIndex: 13,
						StartRowIndex:    rowNo,
						EndColumnIndex:   14,
						EndRowIndex:      rowNo + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: s.spreadSheetService.GetCellColor(ctx, style.RivalHorseColor()),
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.spreadSheetConfig.SheetId(),
						StartColumnIndex: 17,
						StartRowIndex:    rowNo,
						EndColumnIndex:   18,
						EndRowIndex:      rowNo + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: s.spreadSheetService.GetCellColor(ctx, style.FirstPlaceHorseColor()),
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.spreadSheetConfig.SheetId(),
						StartColumnIndex: 21,
						StartRowIndex:    rowNo,
						EndColumnIndex:   22,
						EndRowIndex:      rowNo + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: s.spreadSheetService.GetCellColor(ctx, style.SecondPlaceHorseColor()),
						},
					},
				},
			},
		}...)

		if len(style.PayoutComments()) > 0 {
			requests = append(requests, &sheets.Request{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "note",
					Range: &sheets.GridRange{
						SheetId:          s.spreadSheetConfig.SheetId(),
						StartColumnIndex: 7,
						StartRowIndex:    rowNo,
						EndColumnIndex:   8,
						EndRowIndex:      rowNo + 1,
					},
					Cell: &sheets.CellData{
						Note: strings.Join(style.PayoutComments(), "\n"),
					},
				},
			})
		}
	}

	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		return err
	}

	log.Println(ctx, fmt.Sprintf("write list style end"))

	return nil
}

func (s *spreadSheetListRepository) Clear(ctx context.Context) error {
	requests := []*sheets.Request{
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Fields: "*",
				Range: &sheets.GridRange{
					SheetId:          s.spreadSheetConfig.SheetId(),
					StartColumnIndex: 0,
					StartRowIndex:    0,
					EndColumnIndex:   40,
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
