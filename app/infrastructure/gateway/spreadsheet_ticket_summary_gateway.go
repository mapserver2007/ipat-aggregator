package gateway

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"
	"strconv"
)

const (
	spreadSheetTicketSummaryFileName = "spreadsheet_ticket_summary.json"
)

type SpreadSheetTicketSummaryGateway interface {
	Write(ctx context.Context, ticketSummaryMap map[int]*spreadsheet_entity.TicketSummary) error
	Style(ctx context.Context, ticketSummaryMap map[int]*spreadsheet_entity.TicketSummary) error
	Clear(ctx context.Context) error
}

type spreadSheetTicketSummaryGateway struct {
	logger *logrus.Logger
}

func NewSpreadSheetTicketSummaryGateway(
	logger *logrus.Logger,
) SpreadSheetTicketSummaryGateway {
	return &spreadSheetTicketSummaryGateway{
		logger: logger,
	}
}

func (s *spreadSheetTicketSummaryGateway) Write(
	ctx context.Context,
	ticketSummaryMap map[int]*spreadsheet_entity.TicketSummary,
) error {
	client, config, err := getSpreadSheetConfig(ctx, spreadSheetTicketSummaryFileName)
	if err != nil {
		return err
	}

	s.logger.Infof("write ticket summary start")

	defaultValuesFunc := func(ticketType types.TicketType) [][]interface{} {
		return [][]interface{}{
			{
				ticketType.Name(),
				"月別",
				"投票回数",
				"的中回数",
				"的中率",
				"投資額",
				"回収額",
				"回収率",
				"平均払戻金額",
				"最大払戻金額",
				"最小払戻金額",
			},
		}
	}
	winSummaryValues := defaultValuesFunc(types.Win)
	placeSummaryValues := defaultValuesFunc(types.Place)
	quinellaSummaryValues := defaultValuesFunc(types.Quinella)
	exactaSummaryValues := defaultValuesFunc(types.Exacta)
	quinellaPlaceSummaryValues := defaultValuesFunc(types.QuinellaPlace)
	trioSummaryValues := defaultValuesFunc(types.Trio)
	trifectaSummaryValues := defaultValuesFunc(types.Trifecta)

	for _, date := range SortedIntKeys(ticketSummaryMap) {
		ticketSummary := ticketSummaryMap[date]
		winSummaryValues = s.append(winSummaryValues, date, ticketSummary.WinTermResult())
		placeSummaryValues = s.append(placeSummaryValues, date, ticketSummary.PlaceTermResult())
		quinellaSummaryValues = s.append(quinellaSummaryValues, date, ticketSummary.QuinellaTermResult())
		exactaSummaryValues = s.append(exactaSummaryValues, date, ticketSummary.ExactaTermResult())
		quinellaPlaceSummaryValues = s.append(quinellaPlaceSummaryValues, date, ticketSummary.QuinellaPlaceTermResult())
		trioSummaryValues = s.append(trioSummaryValues, date, ticketSummary.TrioTermResult())
		trifectaSummaryValues = s.append(trifectaSummaryValues, date, ticketSummary.TrifectaTermResult())
	}

	values := s.concatSlices(winSummaryValues, placeSummaryValues, quinellaSummaryValues, exactaSummaryValues, quinellaPlaceSummaryValues, trioSummaryValues, trifectaSummaryValues)

	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "A1")
	_, err = client.Spreadsheets.Values.Update(config.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	s.logger.Infof("write ticket summary end")

	return nil
}

func (s *spreadSheetTicketSummaryGateway) Style(
	ctx context.Context,
	ticketSummaryMap map[int]*spreadsheet_entity.TicketSummary,
) error {
	client, config, err := getSpreadSheetConfig(ctx, spreadSheetTicketSummaryFileName)
	if err != nil {
		return err
	}

	s.logger.Infof("write ticket summary style start")
	var requests []*sheets.Request
	alignment := len(ticketSummaryMap) + 1
	for idx := 0; idx < 7; idx++ {
		requests = append(requests, []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          config.SheetId(),
						StartColumnIndex: 0,
						StartRowIndex:    int64(idx * alignment),
						EndColumnIndex:   1,
						EndRowIndex:      int64(idx*alignment) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0.0,
								Green: 0.0,
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
						StartColumnIndex: 1,
						StartRowIndex:    int64(idx * alignment),
						EndColumnIndex:   11,
						EndRowIndex:      int64(idx*alignment) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0.0,
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
						StartRowIndex:    int64(idx * alignment),
						EndColumnIndex:   11,
						EndRowIndex:      int64(idx*alignment) + 1,
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
						StartColumnIndex: 0,
						StartRowIndex:    int64(idx * alignment),
						EndColumnIndex:   1,
						EndRowIndex:      int64(idx*alignment) + 1,
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
		}...)
	}

	_, err = client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	s.logger.Infof("write ticket summary style end")

	return nil
}

func (s *spreadSheetTicketSummaryGateway) Clear(ctx context.Context) error {
	client, config, err := getSpreadSheetConfig(ctx, spreadSheetTicketSummaryFileName)
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
					EndColumnIndex:   11,
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

func (s *spreadSheetTicketSummaryGateway) append(
	values [][]interface{},
	date int,
	result *spreadsheet_entity.TicketResult,
) [][]interface{} {
	values = append(values, []interface{}{
		"",
		strconv.Itoa(date),
		result.BetCount(),
		result.HitCount(),
		result.HitRate(),
		result.Payment(),
		result.Payout(),
		result.PayoutRate(),
		result.AveragePayout(),
		result.MaxPayout(),
		result.MinPayout(),
	})

	return values
}

func (s *spreadSheetTicketSummaryGateway) concatSlices(slices ...[][]interface{}) [][]interface{} {
	totalLength := 0
	for _, s := range slices {
		totalLength += len(s)
	}

	result := make([][]interface{}, totalLength)
	pos := 0
	for _, s := range slices {
		pos += copy(result[pos:], s)
	}

	return result
}
