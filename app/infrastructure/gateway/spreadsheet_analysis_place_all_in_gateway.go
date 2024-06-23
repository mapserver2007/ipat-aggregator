package gateway

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"google.golang.org/api/sheets/v4"
	"log"
	"time"
)

const (
	spreadSheetAnalysisPlaceAllInFileName = "spreadsheet_analysis_place_all_in.json"
)

type SpreadSheetAnalysisPlaceAllInGateway interface {
	Write(ctx context.Context, placeAllInMap map[filter.Id]*spreadsheet_entity.AnalysisPlaceAllIn, filters []filter.Id) error
	Style(ctx context.Context, placeAllInMap map[filter.Id]*spreadsheet_entity.AnalysisPlaceAllIn, filters []filter.Id) error
	Clear(ctx context.Context) error
}

type spreadSheetAnalysisPlaceAllInGateway struct{}

func NewSpreadSheetAnalysisPlaceAllInGateway() SpreadSheetAnalysisPlaceAllInGateway {
	return &spreadSheetAnalysisPlaceAllInGateway{}
}

func (s *spreadSheetAnalysisPlaceAllInGateway) Write(
	ctx context.Context,
	placeAllInMap map[filter.Id]*spreadsheet_entity.AnalysisPlaceAllIn,
	analysisFilters []filter.Id,
) error {
	client, config, err := getSpreadSheetConfig(ctx, spreadSheetAnalysisPlaceAllInFileName)
	if err != nil {
		return err
	}

	valuesList := make([][][]interface{}, len(analysisFilters)+1)
	valuesList[0] = [][]interface{}{
		{
			"",
			"1.1", "1.2", "1.3", "1.4", "1.5", "1.6", "1.7", "1.8", "1.9",
			"2.0", "2.1", "2.2", "2.3", "2.4", "2.5", "2.6", "2.7", "2.8", "2.9",
			"3.0", "3.1", "3.2", "3.3", "3.4", "3.5", "3.6", "3.7", "3.8", "3.9",
		},
	}

	log.Println(ctx, "write analysis place all in start")
	for idx, analysisFilter := range analysisFilters {
		placeAllIn, ok := placeAllInMap[analysisFilter]
		if !ok {
			return fmt.Errorf("placeAllInMap filter is not found: %d", analysisFilter.OriginFilters())
		}

		var filterName string
		for _, f := range analysisFilter.OriginFilters() {
			filterName += f.String()
		}

		valuesList[idx+1] = [][]interface{}{
			{
				filterName,
				placeAllIn.RateData().WinOdds11HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds12HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds13HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds14HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds15HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds16HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds17HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds18HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds19HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds20HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds21HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds22HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds23HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds24HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds25HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds26HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds27HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds28HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds29HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds30HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds31HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds32HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds33HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds34HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds35HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds36HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds37HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds38HitData().HitRateFormat(),
				placeAllIn.RateData().WinOdds39HitData().HitRateFormat(),
			},
		}
	}

	for idx, values := range valuesList {
		time.Sleep(time.Second)
		if idx > 0 {
			var filterName string
			for _, f := range analysisFilters[idx-1].OriginFilters() {
				filterName += f.String()
			}
			log.Println(ctx, fmt.Sprintf("write analysis place all in filter %s start", filterName))
		}
		writeRange := fmt.Sprintf("%s!%s", config.SheetName(), fmt.Sprintf("A%d", idx+1))
		_, err := client.Spreadsheets.Values.Update(config.SpreadSheetId(), writeRange, &sheets.ValueRange{
			Values: values,
		}).ValueInputOption("USER_ENTERED").Do()
		if err != nil {
			return err
		}
	}

	log.Println(ctx, "write analysis place all in end")

	return nil
}

func (s *spreadSheetAnalysisPlaceAllInGateway) Style(
	ctx context.Context,
	placeAllInMap map[filter.Id]*spreadsheet_entity.AnalysisPlaceAllIn,
	analysisFilters []filter.Id,
) error {
	client, config, err := getSpreadSheetConfig(ctx, spreadSheetAnalysisPlaceAllInFileName)
	if err != nil {
		return err
	}

	var requests []*sheets.Request

	log.Println(ctx, "write style analysis place all in start")
	for rowIdx, analysisFilter := range analysisFilters {
		placeAllIn, ok := placeAllInMap[analysisFilter]
		if !ok {
			return fmt.Errorf("placeAllInMap filter is not found: %d", analysisFilter.OriginFilters())
		}

		cellComments := []string{
			s.getCellComments(placeAllIn.RateData().WinOdds11HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds12HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds13HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds14HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds15HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds16HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds17HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds18HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds19HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds20HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds21HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds22HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds23HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds24HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds25HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds26HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds27HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds28HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds29HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds30HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds31HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds32HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds33HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds34HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds35HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds36HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds37HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds38HitData()),
			s.getCellComments(placeAllIn.RateData().WinOdds39HitData()),
		}
		cellColorTypes := []types.CellColorType{
			placeAllIn.RateStyle().WinOdds11CellColorType(),
			placeAllIn.RateStyle().WinOdds12CellColorType(),
			placeAllIn.RateStyle().WinOdds13CellColorType(),
			placeAllIn.RateStyle().WinOdds14CellColorType(),
			placeAllIn.RateStyle().WinOdds15CellColorType(),
			placeAllIn.RateStyle().WinOdds16CellColorType(),
			placeAllIn.RateStyle().WinOdds17CellColorType(),
			placeAllIn.RateStyle().WinOdds18CellColorType(),
			placeAllIn.RateStyle().WinOdds19CellColorType(),
			placeAllIn.RateStyle().WinOdds20CellColorType(),
			placeAllIn.RateStyle().WinOdds21CellColorType(),
			placeAllIn.RateStyle().WinOdds22CellColorType(),
			placeAllIn.RateStyle().WinOdds23CellColorType(),
			placeAllIn.RateStyle().WinOdds24CellColorType(),
			placeAllIn.RateStyle().WinOdds25CellColorType(),
			placeAllIn.RateStyle().WinOdds26CellColorType(),
			placeAllIn.RateStyle().WinOdds27CellColorType(),
			placeAllIn.RateStyle().WinOdds28CellColorType(),
			placeAllIn.RateStyle().WinOdds29CellColorType(),
			placeAllIn.RateStyle().WinOdds30CellColorType(),
			placeAllIn.RateStyle().WinOdds31CellColorType(),
			placeAllIn.RateStyle().WinOdds32CellColorType(),
			placeAllIn.RateStyle().WinOdds33CellColorType(),
			placeAllIn.RateStyle().WinOdds34CellColorType(),
			placeAllIn.RateStyle().WinOdds35CellColorType(),
			placeAllIn.RateStyle().WinOdds36CellColorType(),
			placeAllIn.RateStyle().WinOdds37CellColorType(),
			placeAllIn.RateStyle().WinOdds38CellColorType(),
			placeAllIn.RateStyle().WinOdds39CellColorType(),
		}
		for colIdx, cellColorType := range cellColorTypes {
			requests = append(requests, []*sheets.Request{
				{
					RepeatCell: &sheets.RepeatCellRequest{
						Fields: "userEnteredFormat.backgroundColor",
						Range: &sheets.GridRange{
							SheetId:          config.SheetId(),
							StartColumnIndex: 1 + int64(colIdx),
							StartRowIndex:    1 + int64(rowIdx),
							EndColumnIndex:   2 + int64(colIdx),
							EndRowIndex:      2 + int64(rowIdx),
						},
						Cell: &sheets.CellData{
							UserEnteredFormat: &sheets.CellFormat{
								BackgroundColor: s.getCellColor(cellColorType),
							},
						},
					},
				},
			}...)
			if len(cellComments[colIdx]) > 0 {
				requests = append(requests, []*sheets.Request{
					{
						RepeatCell: &sheets.RepeatCellRequest{
							Fields: "note",
							Range: &sheets.GridRange{
								SheetId:          config.SheetId(),
								StartColumnIndex: 1 + int64(colIdx),
								StartRowIndex:    1 + int64(rowIdx),
								EndColumnIndex:   2 + int64(colIdx),
								EndRowIndex:      2 + int64(rowIdx),
							},
							Cell: &sheets.CellData{
								Note: cellComments[colIdx],
							},
						},
					},
				}...)
			}
		}
	}

	requests = append(requests, []*sheets.Request{
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Fields: "userEnteredFormat.backgroundColor",
				Range: &sheets.GridRange{
					SheetId:          config.SheetId(),
					StartColumnIndex: 0,
					StartRowIndex:    1,
					EndColumnIndex:   1,
					EndRowIndex:      int64(len(analysisFilters) + 1),
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
					StartColumnIndex: 1,
					StartRowIndex:    0,
					EndColumnIndex:   30,
					EndRowIndex:      1,
				},
				Cell: &sheets.CellData{
					UserEnteredFormat: &sheets.CellFormat{
						BackgroundColor: &sheets.Color{
							Red:   1.0,
							Green: 0.0,
							Blue:  0.0,
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
					StartColumnIndex: 1,
					StartRowIndex:    0,
					EndColumnIndex:   30,
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
					StartColumnIndex: 1,
					StartRowIndex:    0,
					EndColumnIndex:   30,
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
				Fields: "userEnteredFormat.textFormat.bold",
				Range: &sheets.GridRange{
					SheetId:          config.SheetId(),
					StartColumnIndex: 0,
					StartRowIndex:    1,
					EndColumnIndex:   1,
					EndRowIndex:      int64(len(analysisFilters) + 1),
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
	}...)

	_, err = client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		return err
	}

	log.Println(ctx, "write style analysis place all in end")

	return nil
}

func (s *spreadSheetAnalysisPlaceAllInGateway) Clear(ctx context.Context) error {
	client, config, err := getSpreadSheetConfig(ctx, spreadSheetAnalysisPlaceAllInFileName)
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

func (s *spreadSheetAnalysisPlaceAllInGateway) getCellComments(
	data *spreadsheet_entity.PlaceAllInHitData,
) string {
	if data.HitCount()+data.UnHitCount() == 0 {
		return ""
	}
	return fmt.Sprintf("的中%d, 不的中%d", data.HitCount(), data.UnHitCount())
}

func (s *spreadSheetAnalysisPlaceAllInGateway) getCellColor(
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
