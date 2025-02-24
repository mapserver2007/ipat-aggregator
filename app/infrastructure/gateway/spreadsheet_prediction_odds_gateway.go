package gateway

import (
	"context"
	"fmt"
	"sort"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"
)

const (
	spreadSheetPredictionOddsFileName = "spreadsheet_prediction_odds.json"
)

type SpreadSheetPredictionOddsGateway interface {
	Write(ctx context.Context, firstPlaceMap, secondPlaceMap, thirdPlaceMap map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace, raceCourseMap map[types.RaceCourse][]types.RaceId) error
	Style(ctx context.Context, firstPlaceMap, secondPlaceMap, thirdPlaceMap map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace, raceCourseMap map[types.RaceCourse][]types.RaceId) error
	Clear(ctx context.Context) error
}

type spreadSheetPredictionOddsGateway struct {
	logger *logrus.Logger
}

func NewSpreadSheetPredictionOddsGateway(
	logger *logrus.Logger,
) SpreadSheetPredictionOddsGateway {
	return &spreadSheetPredictionOddsGateway{
		logger: logger,
	}
}

func (s *spreadSheetPredictionOddsGateway) Write(
	ctx context.Context,
	firstPlaceMap,
	secondPlaceMap,
	thirdPlaceMap map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace,
	raceCourseMap map[types.RaceCourse][]types.RaceId,
) error {
	client, config, err := getSpreadSheetConfig(ctx, spreadSheetPredictionOddsFileName)
	if err != nil {
		return err
	}

	s.logger.Infof("write prediction odds start")

	raceCourseIds := make([]string, 0, len(raceCourseMap))
	for raceCourseId := range raceCourseMap {
		raceCourseIds = append(raceCourseIds, raceCourseId.Value())
	}
	sort.Strings(raceCourseIds)

	isHitIconFunc := func(isHit bool) string {
		if isHit {
			return "\U0001F3AF"
		}
		return ""
	}

	courseIdx := 0
	for _, raceCourseId := range raceCourseIds {
		raceIds := raceCourseMap[types.RaceCourse(raceCourseId)]
		var valuesList [][]interface{}
		for _, raceId := range raceIds {
			values := make([][][]interface{}, 4)
			values[0] = [][]interface{}{
				{
					"",
					"",
					"",
					"",
					"",
					"",
					"",
					"",
					"",
					"",
					"",
				},
			}
			for idx, placeMap := range []map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace{firstPlaceMap, secondPlaceMap, thirdPlaceMap} {
				switch idx {
				case 0:
					values[1] = [][]interface{}{
						{
							"",
							"1着率",
							types.WinOddsRange1.String(),
							types.WinOddsRange2.String(),
							types.WinOddsRange3.String(),
							types.WinOddsRange4.String(),
							types.WinOddsRange5.String(),
							types.WinOddsRange6.String(),
							types.WinOddsRange7.String(),
							types.WinOddsRange8.String(),
							types.WinOddsRange9.String(),
						},
					}
				case 1:
					values[2] = [][]interface{}{
						{
							"",
							"2着率",
							types.WinOddsRange1.String(),
							types.WinOddsRange2.String(),
							types.WinOddsRange3.String(),
							types.WinOddsRange4.String(),
							types.WinOddsRange5.String(),
							types.WinOddsRange6.String(),
							types.WinOddsRange7.String(),
							types.WinOddsRange8.String(),
							types.WinOddsRange9.String(),
						},
					}
				case 2:
					values[3] = [][]interface{}{
						{
							"",
							"3着率",
							types.WinOddsRange1.String(),
							types.WinOddsRange2.String(),
							types.WinOddsRange3.String(),
							types.WinOddsRange4.String(),
							types.WinOddsRange5.String(),
							types.WinOddsRange6.String(),
							types.WinOddsRange7.String(),
							types.WinOddsRange8.String(),
							types.WinOddsRange9.String(),
						},
					}
				}

				for predictionRace, markerPlaceMap := range placeMap {
					if predictionRace.RaceId() != raceId {
						continue
					}
					if values[0][0][1] == "" {
						title := fmt.Sprintf("%s%dR %s %s", predictionRace.RaceCourseId().Name(), predictionRace.RaceNumber(), predictionRace.RaceName(), predictionRace.FilterName())
						raceCount := markerPlaceMap[types.Favorite].RateData().RaceCount()
						values[0][0][1] = fmt.Sprintf("=HYPERLINK(\"%s\",\"%s(%d)\")", predictionRace.Url(), title, raceCount)
					}

					values[idx+1] = append(values[idx+1], [][]interface{}{
						{
							types.Favorite.String(),
							markerPlaceMap[types.Favorite].RateData().HitRateFormat(),
							isHitIconFunc(markerPlaceMap[types.Favorite].RateData().OddsRange1Hit()) + markerPlaceMap[types.Favorite].RateData().OddsRange1RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Favorite].RateData().OddsRange2Hit()) + markerPlaceMap[types.Favorite].RateData().OddsRange2RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Favorite].RateData().OddsRange3Hit()) + markerPlaceMap[types.Favorite].RateData().OddsRange3RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Favorite].RateData().OddsRange4Hit()) + markerPlaceMap[types.Favorite].RateData().OddsRange4RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Favorite].RateData().OddsRange5Hit()) + markerPlaceMap[types.Favorite].RateData().OddsRange5RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Favorite].RateData().OddsRange6Hit()) + markerPlaceMap[types.Favorite].RateData().OddsRange6RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Favorite].RateData().OddsRange7Hit()) + markerPlaceMap[types.Favorite].RateData().OddsRange7RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Favorite].RateData().OddsRange8Hit()) + markerPlaceMap[types.Favorite].RateData().OddsRange8RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Favorite].RateData().OddsRange9Hit()) + markerPlaceMap[types.Favorite].RateData().OddsRange9RateFormat(),
						},
					}...)
					values[idx+1] = append(values[idx+1], [][]interface{}{
						{
							types.Rival.String(),
							markerPlaceMap[types.Rival].RateData().HitRateFormat(),
							isHitIconFunc(markerPlaceMap[types.Rival].RateData().OddsRange1Hit()) + markerPlaceMap[types.Rival].RateData().OddsRange1RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Rival].RateData().OddsRange2Hit()) + markerPlaceMap[types.Rival].RateData().OddsRange2RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Rival].RateData().OddsRange3Hit()) + markerPlaceMap[types.Rival].RateData().OddsRange3RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Rival].RateData().OddsRange4Hit()) + markerPlaceMap[types.Rival].RateData().OddsRange4RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Rival].RateData().OddsRange5Hit()) + markerPlaceMap[types.Rival].RateData().OddsRange5RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Rival].RateData().OddsRange6Hit()) + markerPlaceMap[types.Rival].RateData().OddsRange6RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Rival].RateData().OddsRange7Hit()) + markerPlaceMap[types.Rival].RateData().OddsRange7RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Rival].RateData().OddsRange8Hit()) + markerPlaceMap[types.Rival].RateData().OddsRange8RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Rival].RateData().OddsRange9Hit()) + markerPlaceMap[types.Rival].RateData().OddsRange9RateFormat(),
						},
					}...)
					values[idx+1] = append(values[idx+1], [][]interface{}{
						{
							types.BrackTriangle.String(),
							markerPlaceMap[types.BrackTriangle].RateData().HitRateFormat(),
							isHitIconFunc(markerPlaceMap[types.BrackTriangle].RateData().OddsRange1Hit()) + markerPlaceMap[types.BrackTriangle].RateData().OddsRange1RateFormat(),
							isHitIconFunc(markerPlaceMap[types.BrackTriangle].RateData().OddsRange2Hit()) + markerPlaceMap[types.BrackTriangle].RateData().OddsRange2RateFormat(),
							isHitIconFunc(markerPlaceMap[types.BrackTriangle].RateData().OddsRange3Hit()) + markerPlaceMap[types.BrackTriangle].RateData().OddsRange3RateFormat(),
							isHitIconFunc(markerPlaceMap[types.BrackTriangle].RateData().OddsRange4Hit()) + markerPlaceMap[types.BrackTriangle].RateData().OddsRange4RateFormat(),
							isHitIconFunc(markerPlaceMap[types.BrackTriangle].RateData().OddsRange5Hit()) + markerPlaceMap[types.BrackTriangle].RateData().OddsRange5RateFormat(),
							isHitIconFunc(markerPlaceMap[types.BrackTriangle].RateData().OddsRange6Hit()) + markerPlaceMap[types.BrackTriangle].RateData().OddsRange6RateFormat(),
							isHitIconFunc(markerPlaceMap[types.BrackTriangle].RateData().OddsRange7Hit()) + markerPlaceMap[types.BrackTriangle].RateData().OddsRange7RateFormat(),
							isHitIconFunc(markerPlaceMap[types.BrackTriangle].RateData().OddsRange8Hit()) + markerPlaceMap[types.BrackTriangle].RateData().OddsRange8RateFormat(),
							isHitIconFunc(markerPlaceMap[types.BrackTriangle].RateData().OddsRange9Hit()) + markerPlaceMap[types.BrackTriangle].RateData().OddsRange9RateFormat(),
						},
					}...)
					values[idx+1] = append(values[idx+1], [][]interface{}{
						{
							types.WhiteTriangle.String(),
							markerPlaceMap[types.WhiteTriangle].RateData().HitRateFormat(),
							isHitIconFunc(markerPlaceMap[types.WhiteTriangle].RateData().OddsRange1Hit()) + markerPlaceMap[types.WhiteTriangle].RateData().OddsRange1RateFormat(),
							isHitIconFunc(markerPlaceMap[types.WhiteTriangle].RateData().OddsRange2Hit()) + markerPlaceMap[types.WhiteTriangle].RateData().OddsRange2RateFormat(),
							isHitIconFunc(markerPlaceMap[types.WhiteTriangle].RateData().OddsRange3Hit()) + markerPlaceMap[types.WhiteTriangle].RateData().OddsRange3RateFormat(),
							isHitIconFunc(markerPlaceMap[types.WhiteTriangle].RateData().OddsRange4Hit()) + markerPlaceMap[types.WhiteTriangle].RateData().OddsRange4RateFormat(),
							isHitIconFunc(markerPlaceMap[types.WhiteTriangle].RateData().OddsRange5Hit()) + markerPlaceMap[types.WhiteTriangle].RateData().OddsRange5RateFormat(),
							isHitIconFunc(markerPlaceMap[types.WhiteTriangle].RateData().OddsRange6Hit()) + markerPlaceMap[types.WhiteTriangle].RateData().OddsRange6RateFormat(),
							isHitIconFunc(markerPlaceMap[types.WhiteTriangle].RateData().OddsRange7Hit()) + markerPlaceMap[types.WhiteTriangle].RateData().OddsRange7RateFormat(),
							isHitIconFunc(markerPlaceMap[types.WhiteTriangle].RateData().OddsRange8Hit()) + markerPlaceMap[types.WhiteTriangle].RateData().OddsRange8RateFormat(),
							isHitIconFunc(markerPlaceMap[types.WhiteTriangle].RateData().OddsRange9Hit()) + markerPlaceMap[types.WhiteTriangle].RateData().OddsRange9RateFormat(),
						},
					}...)
					values[idx+1] = append(values[idx+1], [][]interface{}{
						{
							types.Star.String(),
							markerPlaceMap[types.Star].RateData().HitRateFormat(),
							isHitIconFunc(markerPlaceMap[types.Star].RateData().OddsRange1Hit()) + markerPlaceMap[types.Star].RateData().OddsRange1RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Star].RateData().OddsRange2Hit()) + markerPlaceMap[types.Star].RateData().OddsRange2RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Star].RateData().OddsRange3Hit()) + markerPlaceMap[types.Star].RateData().OddsRange3RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Star].RateData().OddsRange4Hit()) + markerPlaceMap[types.Star].RateData().OddsRange4RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Star].RateData().OddsRange5Hit()) + markerPlaceMap[types.Star].RateData().OddsRange5RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Star].RateData().OddsRange6Hit()) + markerPlaceMap[types.Star].RateData().OddsRange6RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Star].RateData().OddsRange7Hit()) + markerPlaceMap[types.Star].RateData().OddsRange7RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Star].RateData().OddsRange8Hit()) + markerPlaceMap[types.Star].RateData().OddsRange8RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Star].RateData().OddsRange9Hit()) + markerPlaceMap[types.Star].RateData().OddsRange9RateFormat(),
						},
					}...)
					values[idx+1] = append(values[idx+1], [][]interface{}{
						{
							types.Check.String(),
							markerPlaceMap[types.Check].RateData().HitRateFormat(),
							isHitIconFunc(markerPlaceMap[types.Check].RateData().OddsRange1Hit()) + markerPlaceMap[types.Check].RateData().OddsRange1RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Check].RateData().OddsRange2Hit()) + markerPlaceMap[types.Check].RateData().OddsRange2RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Check].RateData().OddsRange3Hit()) + markerPlaceMap[types.Check].RateData().OddsRange3RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Check].RateData().OddsRange4Hit()) + markerPlaceMap[types.Check].RateData().OddsRange4RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Check].RateData().OddsRange5Hit()) + markerPlaceMap[types.Check].RateData().OddsRange5RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Check].RateData().OddsRange6Hit()) + markerPlaceMap[types.Check].RateData().OddsRange6RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Check].RateData().OddsRange7Hit()) + markerPlaceMap[types.Check].RateData().OddsRange7RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Check].RateData().OddsRange8Hit()) + markerPlaceMap[types.Check].RateData().OddsRange8RateFormat(),
							isHitIconFunc(markerPlaceMap[types.Check].RateData().OddsRange9Hit()) + markerPlaceMap[types.Check].RateData().OddsRange9RateFormat(),
						},
					}...)
				}
			}
			for _, value := range values {
				valuesList = append(valuesList, value...)
			}
		}

		var cellId string
		switch courseIdx {
		case 0:
			cellId = "A1"
		case 1:
			cellId = "L1"
		case 2:
			cellId = "W1"
		}

		writeRange := fmt.Sprintf("%s!%s", config.SheetName(), cellId)
		_, err := client.Spreadsheets.Values.Update(config.SpreadSheetId(), writeRange, &sheets.ValueRange{
			Values: valuesList,
		}).ValueInputOption("USER_ENTERED").Do()
		if err != nil {
			return err
		}

		courseIdx++
	}

	s.logger.Infof("write prediction odds end")

	return nil
}

func (s *spreadSheetPredictionOddsGateway) Style(
	ctx context.Context,
	firstPlaceMap,
	secondPlaceMap,
	thirdPlaceMap map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace,
	raceCourseMap map[types.RaceCourse][]types.RaceId,
) error {
	client, config, err := getSpreadSheetConfig(ctx, spreadSheetPredictionOddsFileName)
	if err != nil {
		return err
	}

	s.logger.Infof("write prediction odds style start")

	raceCourseIds := make([]string, 0, len(raceCourseMap))
	for raceCourseId := range raceCourseMap {
		raceCourseIds = append(raceCourseIds, raceCourseId.Value())
	}
	sort.Strings(raceCourseIds)

	var requests []*sheets.Request
	raceCourseCount := 0
	for _, raceCourseId := range raceCourseIds {
		raceIds := raceCourseMap[types.RaceCourse(raceCourseId)]
		for raceIndex, raceId := range raceIds {
			values := make([][][]interface{}, 4)
			values[0] = [][]interface{}{
				{
					"",
					"",
					"",
					"",
					"",
					"",
					"",
					"",
					"",
					"",
					"",
				},
			}
			for placeIndex, placeMap := range []map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace{firstPlaceMap, secondPlaceMap, thirdPlaceMap} {
				for race, markerPlaceMap := range placeMap {
					if race.RaceId() != raceId {
						continue
					}
					for markerIndex, marker := range []types.Marker{types.Favorite, types.Rival, types.BrackTriangle, types.WhiteTriangle, types.Star, types.Check} {
						oddsRangeIndex := markerPlaceMap[marker].RateStyle().MatchOddsRangeIndex()
						requests = append(requests, []*sheets.Request{
							{
								RepeatCell: &sheets.RepeatCellRequest{
									Fields: "userEnteredFormat.backgroundColor",
									Range: &sheets.GridRange{
										SheetId:          config.SheetId(),
										StartColumnIndex: 2 + int64(oddsRangeIndex) + int64(raceCourseCount*11),
										StartRowIndex:    2 + int64(raceIndex*22+markerIndex) + int64(placeIndex*7),
										EndColumnIndex:   3 + int64(oddsRangeIndex) + int64(raceCourseCount*11),
										EndRowIndex:      3 + int64(raceIndex*22+markerIndex) + int64(placeIndex*7),
									},
									Cell: &sheets.CellData{
										UserEnteredFormat: &sheets.CellFormat{
											BackgroundColor: &sheets.Color{
												Red:   1.0,
												Green: 0.937,
												Blue:  0.498,
											},
										},
									},
								},
							},
							{
								RepeatCell: &sheets.RepeatCellRequest{
									Fields: "userEnteredFormat.horizontalAlignment",
									Range: &sheets.GridRange{
										SheetId:          config.SheetId(),
										StartColumnIndex: 2 + int64(oddsRangeIndex) + int64(raceCourseCount*11),
										StartRowIndex:    2 + int64(raceIndex*22+markerIndex) + int64(placeIndex*7),
										EndColumnIndex:   3 + int64(oddsRangeIndex) + int64(raceCourseCount*11),
										EndRowIndex:      3 + int64(raceIndex*22+markerIndex) + int64(placeIndex*7),
									},
									Cell: &sheets.CellData{
										UserEnteredFormat: &sheets.CellFormat{
											HorizontalAlignment: "RIGHT",
										},
									},
								},
							},
							{
								RepeatCell: &sheets.RepeatCellRequest{
									Fields: "userEnteredFormat.textFormat.bold",
									Range: &sheets.GridRange{
										SheetId:          config.SheetId(),
										StartColumnIndex: 2 + int64(oddsRangeIndex) + int64(raceCourseCount*11),
										StartRowIndex:    2 + int64(raceIndex*22+markerIndex) + int64(placeIndex*7),
										EndColumnIndex:   3 + int64(oddsRangeIndex) + int64(raceCourseCount*11),
										EndRowIndex:      3 + int64(raceIndex*22+markerIndex) + int64(placeIndex*7),
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
					}
					requests = append(requests, []*sheets.Request{
						{
							RepeatCell: &sheets.RepeatCellRequest{
								Fields: "userEnteredFormat.backgroundColor",
								Range: &sheets.GridRange{
									SheetId:          config.SheetId(),
									StartColumnIndex: 1 + int64(raceCourseCount*11),
									StartRowIndex:    1 + int64(placeIndex*7) + int64(raceIndex*22),
									EndColumnIndex:   2 + int64(raceCourseCount*11),
									EndRowIndex:      2 + int64(placeIndex*7) + int64(raceIndex*22),
								},
								Cell: &sheets.CellData{
									UserEnteredFormat: &sheets.CellFormat{
										BackgroundColor: &sheets.Color{
											Red:   1.0,
											Green: 0.937,
											Blue:  0.498,
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
									StartColumnIndex: 1 + int64(raceCourseCount*11),
									StartRowIndex:    1 + int64(placeIndex*7) + int64(raceIndex*22),
									EndColumnIndex:   11 + int64(raceCourseCount*11),
									EndRowIndex:      2 + int64(placeIndex*7) + int64(raceIndex*22),
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
									SheetId:          config.SheetId(),
									StartColumnIndex: 2 + int64(raceCourseCount*11),
									StartRowIndex:    1 + int64(placeIndex*7) + int64(raceIndex*22),
									EndColumnIndex:   11 + int64(raceCourseCount*11),
									EndRowIndex:      2 + int64(placeIndex*7) + int64(raceIndex*22),
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
								Fields: "userEnteredFormat.textFormat.foregroundColor",
								Range: &sheets.GridRange{
									SheetId:          config.SheetId(),
									StartColumnIndex: 2 + int64(raceCourseCount*11),
									StartRowIndex:    1 + int64(placeIndex*7) + int64(raceIndex*22),
									EndColumnIndex:   11 + int64(raceCourseCount*11),
									EndRowIndex:      2 + int64(placeIndex*7) + int64(raceIndex*22),
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
			}
			requests = append(requests, []*sheets.Request{
				{
					RepeatCell: &sheets.RepeatCellRequest{
						Fields: "userEnteredFormat.backgroundColor",
						Range: &sheets.GridRange{
							SheetId:          config.SheetId(),
							StartColumnIndex: 1 + int64(raceCourseCount*11),
							StartRowIndex:    int64(raceIndex * 22),
							EndColumnIndex:   11 + int64(raceCourseCount*11),
							EndRowIndex:      1 + int64(raceIndex*22),
						},
						Cell: &sheets.CellData{
							UserEnteredFormat: &sheets.CellFormat{
								BackgroundColor: &sheets.Color{
									Red:   0.0,
									Blue:  1.0,
									Green: 0.0,
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
							StartColumnIndex: 1 + int64(raceCourseCount*11),
							StartRowIndex:    int64(raceIndex * 22),
							EndColumnIndex:   11 + int64(raceCourseCount*11),
							EndRowIndex:      1 + int64(raceIndex*22),
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
							StartColumnIndex: 1 + int64(raceCourseCount*11),
							StartRowIndex:    int64(raceIndex * 22),
							EndColumnIndex:   11 + int64(raceCourseCount*11),
							EndRowIndex:      1 + int64(raceIndex*22),
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
		}
		raceCourseCount++
	}

	_, err = client.Spreadsheets.BatchUpdate(config.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		return err
	}

	s.logger.Infof("write prediction odds style end")

	return nil
}

func (s *spreadSheetPredictionOddsGateway) Clear(ctx context.Context) error {
	client, config, err := getSpreadSheetConfig(ctx, spreadSheetPredictionOddsFileName)
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

func (s *spreadSheetPredictionOddsGateway) getCellColor(
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
