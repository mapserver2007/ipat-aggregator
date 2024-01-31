package infrastructure

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"google.golang.org/api/sheets/v4"
	"log"
)

const (
	spreadSheetMarkerAnalysisFileName = "spreadsheet_marker_analysis.json"
)

type spreadSheetMarkerAnalysisRepository struct {
	client            *sheets.Service
	spreadSheetConfig *spreadsheet_entity.SpreadSheetConfig
}

func NewSpreadSheetMarkerAnalysisRepository() (repository.SpreadSheetMarkerAnalysisRepository, error) {
	ctx := context.Background()
	client, spreadSheetConfig, err := getSpreadSheetConfig(ctx, spreadSheetMarkerAnalysisFileName)
	if err != nil {
		return nil, err
	}

	return &spreadSheetMarkerAnalysisRepository{
		client:            client,
		spreadSheetConfig: spreadSheetConfig,
	}, nil
}

func (s *spreadSheetMarkerAnalysisRepository) Write(
	ctx context.Context,
	analysisData *spreadsheet_entity.AnalysisData,
) error {
	log.Println(ctx, "write marker analysis start")
	var valuesList [4][][]interface{}
	valuesList[0] = [][]interface{}{
		{
			"",
			"対象レース数",
			"印的中率",
			types.WinOddsRange1.String(),
			types.WinOddsRange2.String(),
			types.WinOddsRange3.String(),
			types.WinOddsRange4.String(),
			types.WinOddsRange5.String(),
			types.WinOddsRange6.String(),
			types.WinOddsRange7.String(),
			types.WinOddsRange8.String(),
		},
	}
	valuesList[1] = [][]interface{}{
		{
			"",
			"対象レース数",
			"印的中数",
			types.WinOddsRange1.String(),
			types.WinOddsRange2.String(),
			types.WinOddsRange3.String(),
			types.WinOddsRange4.String(),
			types.WinOddsRange5.String(),
			types.WinOddsRange6.String(),
			types.WinOddsRange7.String(),
			types.WinOddsRange8.String(),
		},
	}
	valuesList[2] = [][]interface{}{
		{
			"",
			"対象レース数",
			"印不的中率",
			types.WinOddsRange1.String(),
			types.WinOddsRange2.String(),
			types.WinOddsRange3.String(),
			types.WinOddsRange4.String(),
			types.WinOddsRange5.String(),
			types.WinOddsRange6.String(),
			types.WinOddsRange7.String(),
			types.WinOddsRange8.String(),
		},
	}
	valuesList[3] = [][]interface{}{
		{
			"",
			"対象レース数",
			"印不的中数",
			types.WinOddsRange1.String(),
			types.WinOddsRange2.String(),
			types.WinOddsRange3.String(),
			types.WinOddsRange4.String(),
			types.WinOddsRange5.String(),
			types.WinOddsRange6.String(),
			types.WinOddsRange7.String(),
			types.WinOddsRange8.String(),
		},
	}

	rateFormatFunc := func(matchCount int, raceCount int) string {
		return fmt.Sprintf("%.2f%%", float64(matchCount)*100/float64(raceCount))
	}

	allMarkerCombinationIds := analysisData.AllMarkerCombinationIds()
	hitDataMap := analysisData.HitDataMapByFilter()
	unHitDataMap := analysisData.UnHitDataMapByFilter()
	raceCountMap := analysisData.RaceCountByFilter()

	// TODO 関数化してループでフィルタ条件を回す
	filters := []filter.Id{filter.All, filter.Turf}
	for _, f := range filters {
		for _, markerCombinationId := range allMarkerCombinationIds {
			data, ok := hitDataMap[f][markerCombinationId]
			if ok {
				switch markerCombinationId.TicketType() {
				case types.Win:
					marker, err := types.NewMarker(markerCombinationId.Value() % 10)
					if err != nil {
						return err
					}
					if marker != types.Favorite {
						continue
					}

					oddsRangeMap := s.createWinOddsRangeMap(ctx, data)
					// TODO フィルタ条件でvaluesの行を作っていく。フィルタの条件の個数だけforループまわす
					valuesList[0] = append(valuesList[0], [][]interface{}{
						{
							f.String(),
							raceCountMap[f],
							data.MatchRateFormat(),
							rateFormatFunc(oddsRangeMap[types.WinOddsRange1], data.MatchCount()),
							rateFormatFunc(oddsRangeMap[types.WinOddsRange2], data.MatchCount()),
							rateFormatFunc(oddsRangeMap[types.WinOddsRange3], data.MatchCount()),
							rateFormatFunc(oddsRangeMap[types.WinOddsRange4], data.MatchCount()),
							rateFormatFunc(oddsRangeMap[types.WinOddsRange5], data.MatchCount()),
							rateFormatFunc(oddsRangeMap[types.WinOddsRange6], data.MatchCount()),
							rateFormatFunc(oddsRangeMap[types.WinOddsRange7], data.MatchCount()),
							rateFormatFunc(oddsRangeMap[types.WinOddsRange8], data.MatchCount()),
						},
					}...)

					valuesList[1] = append(valuesList[1], [][]interface{}{
						{
							f.String(),
							raceCountMap[f],
							data.MatchCount(),
							oddsRangeMap[types.WinOddsRange1],
							oddsRangeMap[types.WinOddsRange2],
							oddsRangeMap[types.WinOddsRange3],
							oddsRangeMap[types.WinOddsRange4],
							oddsRangeMap[types.WinOddsRange5],
							oddsRangeMap[types.WinOddsRange6],
							oddsRangeMap[types.WinOddsRange7],
							oddsRangeMap[types.WinOddsRange8],
						},
					}...)
				}
			}
			data, ok = unHitDataMap[f][markerCombinationId]
			if ok {
				switch markerCombinationId.TicketType() {
				case types.Win:
					marker, err := types.NewMarker(markerCombinationId.Value() % 10)
					if err != nil {
						return err
					}
					if marker != types.Favorite {
						continue
					}

					oddsRangeMap := s.createWinOddsRangeMap(ctx, data)
					valuesList[2] = append(valuesList[2], [][]interface{}{
						{
							f.String(),
							raceCountMap[f],
							data.MatchRateFormat(),
							rateFormatFunc(oddsRangeMap[types.WinOddsRange1], data.MatchCount()),
							rateFormatFunc(oddsRangeMap[types.WinOddsRange2], data.MatchCount()),
							rateFormatFunc(oddsRangeMap[types.WinOddsRange3], data.MatchCount()),
							rateFormatFunc(oddsRangeMap[types.WinOddsRange4], data.MatchCount()),
							rateFormatFunc(oddsRangeMap[types.WinOddsRange5], data.MatchCount()),
							rateFormatFunc(oddsRangeMap[types.WinOddsRange6], data.MatchCount()),
							rateFormatFunc(oddsRangeMap[types.WinOddsRange7], data.MatchCount()),
							rateFormatFunc(oddsRangeMap[types.WinOddsRange8], data.MatchCount()),
						},
					}...)

					valuesList[3] = append(valuesList[3], [][]interface{}{
						{
							f.String(),
							raceCountMap[f],
							data.MatchCount(),
							oddsRangeMap[types.WinOddsRange1],
							oddsRangeMap[types.WinOddsRange2],
							oddsRangeMap[types.WinOddsRange3],
							oddsRangeMap[types.WinOddsRange4],
							oddsRangeMap[types.WinOddsRange5],
							oddsRangeMap[types.WinOddsRange6],
							oddsRangeMap[types.WinOddsRange7],
							oddsRangeMap[types.WinOddsRange8],
						},
					}...)
				}
			}
		}
	}

	for idx, values := range valuesList {
		writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName(), fmt.Sprintf("A%d", idx*(len(filters)+1)+1))
		_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
			Values: values,
		}).ValueInputOption("USER_ENTERED").Do()
		if err != nil {
			return err
		}
	}

	log.Println(ctx, "write marker analysis end")

	return nil
}

func (s *spreadSheetMarkerAnalysisRepository) createWinOddsRangeMap(
	ctx context.Context,
	markerCombinationAnalysis *spreadsheet_entity.MarkerCombinationAnalysis,
) map[types.OddsRangeType]int {
	oddsRangeMap := map[types.OddsRangeType]int{
		types.WinOddsRange1: 0,
		types.WinOddsRange2: 0,
		types.WinOddsRange3: 0,
		types.WinOddsRange4: 0,
		types.WinOddsRange5: 0,
		types.WinOddsRange6: 0,
		types.WinOddsRange7: 0,
		types.WinOddsRange8: 0,
	}

	for _, decimalOdds := range markerCombinationAnalysis.Odds() {
		odds := decimalOdds.InexactFloat64()
		if odds >= 1.0 && odds <= 1.5 {
			oddsRangeMap[types.WinOddsRange1]++
		} else if odds >= 1.6 && odds <= 2.0 {
			oddsRangeMap[types.WinOddsRange2]++
		} else if odds >= 2.1 && odds <= 2.9 {
			oddsRangeMap[types.WinOddsRange3]++
		} else if odds >= 3.0 && odds <= 4.9 {
			oddsRangeMap[types.WinOddsRange4]++
		} else if odds >= 5.0 && odds <= 9.9 {
			oddsRangeMap[types.WinOddsRange5]++
		} else if odds >= 10.0 && odds <= 19.9 {
			oddsRangeMap[types.WinOddsRange6]++
		} else if odds >= 20.0 && odds <= 49.9 {
			oddsRangeMap[types.WinOddsRange7]++
		} else if odds >= 50.0 {
			oddsRangeMap[types.WinOddsRange8]++
		}
	}

	return oddsRangeMap
}

func (s *spreadSheetMarkerAnalysisRepository) Style(
	ctx context.Context,
	analysisData *spreadsheet_entity.AnalysisData,
) error {
	log.Println(ctx, "write style marker analysis start")
	rowNum := 3
	currentTicketType := types.UnknownTicketType
	allMarkerCombinationIds := analysisData.AllMarkerCombinationIds()
	for _, markerCombinationId := range allMarkerCombinationIds {
		if currentTicketType != markerCombinationId.TicketType() {
			currentTicketType = markerCombinationId.TicketType()
			switch currentTicketType {
			case types.Win:
				marker, err := types.NewMarker(markerCombinationId.Value() % 10)
				if err != nil {
					return err
				}
				if marker != types.Favorite {
					continue
				}

				for i := 0; i < 4; i++ {
					_, err = s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
						Requests: []*sheets.Request{
							{
								RepeatCell: &sheets.RepeatCellRequest{
									Fields: "userEnteredFormat.textFormat.foregroundColor",
									Range: &sheets.GridRange{
										SheetId:          s.spreadSheetConfig.SheetId(),
										StartColumnIndex: 3,
										StartRowIndex:    int64(i * rowNum),
										EndColumnIndex:   11,
										EndRowIndex:      int64(i*rowNum) + 1,
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
									Fields: "userEnteredFormat.backgroundColor",
									Range: &sheets.GridRange{
										SheetId:          s.spreadSheetConfig.SheetId(),
										StartColumnIndex: 1,
										StartRowIndex:    int64(i * rowNum),
										EndColumnIndex:   4,
										EndRowIndex:      int64(i*rowNum) + 1,
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
									Fields: "userEnteredFormat.backgroundColor",
									Range: &sheets.GridRange{
										SheetId:          s.spreadSheetConfig.SheetId(),
										StartColumnIndex: 3,
										StartRowIndex:    int64(i * rowNum),
										EndColumnIndex:   11,
										EndRowIndex:      int64(i*rowNum) + 1,
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
									Fields: "userEnteredFormat.textFormat.bold",
									Range: &sheets.GridRange{
										SheetId:          s.spreadSheetConfig.SheetId(),
										StartColumnIndex: 1,
										StartRowIndex:    int64(i * rowNum),
										EndColumnIndex:   11,
										EndRowIndex:      int64(i*rowNum) + 1,
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
									Fields: "userEnteredFormat.backgroundColor",
									Range: &sheets.GridRange{
										SheetId:          s.spreadSheetConfig.SheetId(),
										StartColumnIndex: 0,
										StartRowIndex:    int64(i*rowNum) + 1,
										EndColumnIndex:   1,
										EndRowIndex:      int64((i + 1) * rowNum),
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
										SheetId:          s.spreadSheetConfig.SheetId(),
										StartColumnIndex: 0,
										StartRowIndex:    int64(i*rowNum) + 1,
										EndColumnIndex:   1,
										EndRowIndex:      int64((i + 1) * rowNum),
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
				}

			default:
				continue // TODO 単勝だけにとりあえず注力するので塞いでおく
			}
		}
	}

	return nil
}

func (s *spreadSheetMarkerAnalysisRepository) writeFilter(
	ctx context.Context,
) error {
	log.Println(ctx, "writing spreadsheet writeFilter in marker analysis")

	return nil
}

func (s *spreadSheetMarkerAnalysisRepository) Clear(ctx context.Context) error {
	requests := []*sheets.Request{
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Fields: "*",
				Range: &sheets.GridRange{
					SheetId:          s.spreadSheetConfig.SheetId(),
					StartColumnIndex: 0,
					StartRowIndex:    0,
					EndColumnIndex:   16,
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
