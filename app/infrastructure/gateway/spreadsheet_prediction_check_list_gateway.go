package gateway

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"google.golang.org/api/sheets/v4"
	"log"
)

const (
	spreadSheetPredictionCheckListFileName = "spreadsheet_prediction_check_list.json"
)

var checkListItems = []string{
	"13頭立て以下であること",
	"単勝1倍台であること",
	"3着以内率80%であること",
	"芝ダート替わりでないこと",
	"前走または2走前と今走の距離が同じなこと",
	"前走または2走前と今走のコースが同じなこと",
	"前走または2走前に馬券内なこと",
	"今走の馬場状態と同じ馬場状態で馬券内経験があること",
	"斤量増でないこと",
	"昇級初戦でないこと",
	"継続騎乗もしくは鞍上強化であること",
	"近2走出遅れがないこと",
	"東スポ印◎が50%以上であること",
	"東スポ印が◎◯のみで構成されていること",
	"調教イチ押しであること",
}

type SpreadSheetPredictionCheckListGateway interface {
	Write(ctx context.Context, checkLists []*spreadsheet_entity.PredictionCheckList) error
	Style(ctx context.Context) error
	Clear(ctx context.Context) error
}

type spreadSheetPredictionCheckListGateway struct{}

func NewSpreadSheetPredictionCheckListGateway() SpreadSheetPredictionCheckListGateway {
	return &spreadSheetPredictionCheckListGateway{}
}

func (s *spreadSheetPredictionCheckListGateway) Write(
	ctx context.Context,
	rows []*spreadsheet_entity.PredictionCheckList,
) error {
	client, config, err := getSpreadSheetConfig(ctx, spreadSheetPredictionCheckListFileName)
	if err != nil {
		return err
	}

	log.Println(ctx, "write prediction check list start")

	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "A1")
	values := [][]interface{}{
		{
			"日付",
			"レース名",
			"馬名",
			"単勝",
			"印",
			"1着率",
			"2着率",
			"3着率",
			"1",
			"2",
			"3",
			"4",
			"5",
			"6",
			"7",
			"8",
			"9",
			"10",
			"11",
			"12",
			"13",
			"14",
			"15",
			"計",
			"◎",
			"◯",
			"印数",
			"推",
			"厩舎コメント",
		},
	}

	for _, row := range rows {
		values = append(values, []interface{}{
			row.RaceDate(),
			fmt.Sprintf("=HYPERLINK(\"%s\",\"%s\")", row.RaceUrl(), row.RaceName()),
			fmt.Sprintf("=HYPERLINK(\"%s\",\"%s\")", row.HorseUrl(), row.HorseName()),
			row.WinOdds(),
			row.Marker(),
			row.FirstPlaceRate(),
			row.SecondPlaceRate(),
			row.ThirdPlaceRate(),
			row.CheckList()[0],
			row.CheckList()[1],
			row.CheckList()[2],
			row.CheckList()[3],
			row.CheckList()[4],
			row.CheckList()[5],
			row.CheckList()[6],
			row.CheckList()[7],
			row.CheckList()[8],
			row.CheckList()[9],
			row.CheckList()[10],
			row.CheckList()[11],
			row.CheckList()[12],
			row.CheckList()[13],
			row.CheckList()[14],
			func() int {
				count := 0
				for _, check := range row.CheckList() {
					if check == "◯" {
						count++
					}
				}
				return count
			}(),
			row.FavoriteNum(),
			row.RivalNum(),
			row.MarkerNum(),
			row.HighlyRecommended(),
			row.TrainingComment(),
		})
	}

	_, err = client.Spreadsheets.Values.Update(config.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	log.Println(ctx, "write prediction check list end")

	return nil
}

func (s *spreadSheetPredictionCheckListGateway) Style(ctx context.Context) error {
	client, config, err := getSpreadSheetConfig(ctx, spreadSheetPredictionCheckListFileName)
	if err != nil {
		return err
	}

	log.Println(ctx, "write prediction check list style start")

	var requests []*sheets.Request
	requests = append(requests, []*sheets.Request{
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Fields: "userEnteredFormat.backgroundColor",
				Range: &sheets.GridRange{
					SheetId:          config.SheetId(),
					StartColumnIndex: 0,
					StartRowIndex:    0,
					EndColumnIndex:   29,
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
		},
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Fields: "userEnteredFormat.backgroundColor",
				Range: &sheets.GridRange{
					SheetId:          config.SheetId(),
					StartColumnIndex: 24,
					StartRowIndex:    0,
					EndColumnIndex:   29,
					EndRowIndex:      1,
				},
				Cell: &sheets.CellData{
					UserEnteredFormat: &sheets.CellFormat{
						BackgroundColor: &sheets.Color{
							Red:   0.0,
							Green: 0.0,
							Blue:  1.0,
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
					EndColumnIndex:   29,
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
				Fields: "userEnteredFormat.textFormat.foregroundColor",
				Range: &sheets.GridRange{
					SheetId:          config.SheetId(),
					StartColumnIndex: 24,
					StartRowIndex:    0,
					EndColumnIndex:   29,
					EndRowIndex:      1,
				},
				Cell: &sheets.CellData{
					UserEnteredFormat: &sheets.CellFormat{
						TextFormat: &sheets.TextFormat{
							ForegroundColor: &sheets.Color{
								Red:   1.0,
								Green: 1.0,
								Blue:  1.0,
							},
						},
					},
				},
			},
		},
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Fields: "userEnteredFormat(verticalAlignment)",
				Range: &sheets.GridRange{
					SheetId:          config.SheetId(),
					StartColumnIndex: 0,
					StartRowIndex:    1,
					EndColumnIndex:   29,
					EndRowIndex:      999,
				},
				Cell: &sheets.CellData{
					UserEnteredFormat: &sheets.CellFormat{
						VerticalAlignment: "TOP",
					},
				},
			},
		},
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Fields: "userEnteredFormat(horizontalAlignment,wrapStrategy)",
				Range: &sheets.GridRange{
					SheetId:          config.SheetId(),
					StartColumnIndex: 28,
					StartRowIndex:    1,
					EndColumnIndex:   29,
					EndRowIndex:      999,
				},
				Cell: &sheets.CellData{
					UserEnteredFormat: &sheets.CellFormat{
						WrapStrategy: "WRAP",
					},
				},
			},
		},
	}...)

	for i := int64(0); i < int64(15); i++ {
		requests = append(requests, &sheets.Request{
			RepeatCell: &sheets.RepeatCellRequest{
				Fields: "note",
				Range: &sheets.GridRange{
					SheetId:          config.SheetId(),
					StartColumnIndex: i + 8,
					StartRowIndex:    0,
					EndColumnIndex:   i + 9,
					EndRowIndex:      1,
				},
				Cell: &sheets.CellData{
					Note: checkListItems[i],
				},
			},
		})
	}

	_, err = client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		return err
	}

	log.Println(ctx, "write prediction check list style end")

	return nil
}

func (s *spreadSheetPredictionCheckListGateway) Clear(ctx context.Context) error {
	client, config, err := getSpreadSheetConfig(ctx, spreadSheetPredictionCheckListFileName)
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
