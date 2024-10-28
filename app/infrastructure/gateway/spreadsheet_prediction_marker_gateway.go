package gateway

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"google.golang.org/api/sheets/v4"
	"log"
)

const (
	spreadSheetPredictionMarkerFileName = "spreadsheet_prediction_marker.json"
)

type SpreadSheetPredictionMarkerGateway interface {
	Write(ctx context.Context, rows []*spreadsheet_entity.PredictionMarker) error
	Clear(ctx context.Context) error
}

type spreadSheetPredictionMarkerGateway struct{}

func NewSpreadSheetPredictionMarkerGateway() SpreadSheetPredictionMarkerGateway {
	return &spreadSheetPredictionMarkerGateway{}
}

func (s *spreadSheetPredictionMarkerGateway) Write(
	ctx context.Context,
	rows []*spreadsheet_entity.PredictionMarker,
) error {
	client, config, err := getSpreadSheetConfig(ctx, spreadSheetPredictionMarkerFileName)
	if err != nil {
		return err
	}

	log.Println(ctx, "write prediction marker start")

	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "A1")
	values := [][]interface{}{
		{
			"レースID",
			"◎",
			"◯",
			"▲",
			"△",
			"☆",
			"✓",
		},
	}

	for _, row := range rows {
		values = append(values, []interface{}{
			row.RaceId(),
			row.FavoriteHorseNumber(),
			row.RivalHorseNumber(),
			row.BrackTriangleHorseNumber(),
			row.WhiteTriangleHorseNumber(),
			row.StarHorseNumber(),
			row.CheckHorseNumber(),
		})
	}

	_, err = client.Spreadsheets.Values.Update(config.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	log.Println(ctx, "write prediction marker end")

	return nil
}

func (s *spreadSheetPredictionMarkerGateway) Clear(ctx context.Context) error {
	client, config, err := getSpreadSheetConfig(ctx, spreadSheetPredictionMarkerFileName)
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
					EndColumnIndex:   7,
					EndRowIndex:      100,
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
