package service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"github.com/shopspring/decimal"
	"google.golang.org/api/sheets/v4"
)

type SpreadSheetService interface {
	CreateMarkerCombinationAnalysisData(ctx context.Context, analysisData *analysis_entity.Layer1, filter filter.Id) map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis
	CreateOddsRangeRaceCountMap(ctx context.Context, analysisData *analysis_entity.Layer1, filter filter.Id) map[types.MarkerCombinationId]map[types.OddsRangeType]int
	CreatePredictionOdds(ctx context.Context, marker *marker_csv_entity.PredictionMarker, race *prediction_entity.Race) map[types.Marker]types.OddsRangeType
	GetCellColor(ctx context.Context, colorType types.CellColorType) *sheets.Color
}

type spreadSheetService struct{}

func NewSpreadSheetService() SpreadSheetService {
	return &spreadSheetService{}
}

func (s *spreadSheetService) CreateMarkerCombinationAnalysisData(
	ctx context.Context,
	analysisData *analysis_entity.Layer1,
	filter filter.Id,
) map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis {
	markerCombinationDataMap := map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis{}
	raceCountMap := s.CreateOddsRangeRaceCountMap(ctx, analysisData, filter)
	for markerCombinationId, data := range analysisData.MarkerCombination {
		for _, data2 := range data.RaceDate {
			for _, data3 := range data2.RaceId {
				for _, calculable := range data3 {
					switch markerCombinationId.TicketType() {
					case types.Win, types.Place:
						if _, ok := markerCombinationDataMap[markerCombinationId]; !ok {
							markerCombinationDataMap[markerCombinationId] = spreadsheet_entity.NewMarkerCombinationAnalysis(raceCountMap[markerCombinationId])
						}
						match := true
						for _, f := range calculable.Filters() {
							if f&filter == 0 {
								match = false
								break
							}
						}
						if match {
							markerCombinationDataMap[markerCombinationId].AddCalculable(calculable)
						}
					}
				}
			}
		}
	}

	return markerCombinationDataMap
}

func (s *spreadSheetService) CreateOddsRangeRaceCountMap(
	ctx context.Context,
	analysisData *analysis_entity.Layer1,
	filter filter.Id,
) map[types.MarkerCombinationId]map[types.OddsRangeType]int {
	markerCombinationOddsRangeCountMap := map[types.MarkerCombinationId]map[types.OddsRangeType]int{}
	for markerCombinationId, data := range analysisData.MarkerCombination {
		if _, ok := markerCombinationOddsRangeCountMap[markerCombinationId]; !ok {
			markerCombinationOddsRangeCountMap[markerCombinationId] = map[types.OddsRangeType]int{}
		}

		for _, data2 := range data.RaceDate {
			for _, data3 := range data2.RaceId {
				match := true
				for _, calculable := range data3 {
					// レースIDに対して複数の結果があるケースは、複勝ワイド、同着のケース
					for _, f := range calculable.Filters() {
						// フィルタマッチ条件は同一レースになるため、ループを回さなくても1件目のチェックとおなじになるはず
						// だが一応全部チェックして1つでもマッチしなければフィルタマッチしないとする
						if f&filter == 0 {
							match = false
							break
						}
					}
					if match {
						odds := calculable.Odds().InexactFloat64()
						if odds >= 1.0 && odds <= 1.5 {
							markerCombinationOddsRangeCountMap[markerCombinationId][types.WinOddsRange1]++
						} else if odds >= 1.6 && odds <= 2.0 {
							markerCombinationOddsRangeCountMap[markerCombinationId][types.WinOddsRange2]++
						} else if odds >= 2.1 && odds <= 2.9 {
							markerCombinationOddsRangeCountMap[markerCombinationId][types.WinOddsRange3]++
						} else if odds >= 3.0 && odds <= 4.9 {
							markerCombinationOddsRangeCountMap[markerCombinationId][types.WinOddsRange4]++
						} else if odds >= 5.0 && odds <= 9.9 {
							markerCombinationOddsRangeCountMap[markerCombinationId][types.WinOddsRange5]++
						} else if odds >= 10.0 && odds <= 19.9 {
							markerCombinationOddsRangeCountMap[markerCombinationId][types.WinOddsRange6]++
						} else if odds >= 20.0 && odds <= 49.9 {
							markerCombinationOddsRangeCountMap[markerCombinationId][types.WinOddsRange7]++
						} else if odds >= 50.0 {
							markerCombinationOddsRangeCountMap[markerCombinationId][types.WinOddsRange8]++
						}
					}
				}
			}
		}
	}

	return markerCombinationOddsRangeCountMap
}

func (s *spreadSheetService) CreatePredictionOdds(
	ctx context.Context,
	marker *marker_csv_entity.PredictionMarker,
	race *prediction_entity.Race,
) map[types.Marker]types.OddsRangeType {
	markerOddsRangeMap := map[types.Marker]types.OddsRangeType{}
	getOddsRange := func(decimalOdds decimal.Decimal) types.OddsRangeType {
		odds := decimalOdds.InexactFloat64()
		if odds >= 1.0 && odds <= 1.5 {
			return types.WinOddsRange1
		} else if odds >= 1.6 && odds <= 2.0 {
			return types.WinOddsRange2
		} else if odds >= 2.1 && odds <= 2.9 {
			return types.WinOddsRange3
		} else if odds >= 3.0 && odds <= 4.9 {
			return types.WinOddsRange4
		} else if odds >= 5.0 && odds <= 9.9 {
			return types.WinOddsRange5
		} else if odds >= 10.0 && odds <= 19.9 {
			return types.WinOddsRange6
		} else if odds >= 20.0 && odds <= 49.9 {
			return types.WinOddsRange7
		} else if odds >= 50.0 {
			return types.WinOddsRange8
		}
		return types.UnknownOddsRangeType
	}

	for _, odds := range race.Odds() {
		switch odds.HorseNumber() {
		case marker.Favorite():
			markerOddsRangeMap[types.Favorite] = getOddsRange(odds.Odds())
		case marker.Rival():
			markerOddsRangeMap[types.Rival] = getOddsRange(odds.Odds())
		case marker.BrackTriangle():
			markerOddsRangeMap[types.BrackTriangle] = getOddsRange(odds.Odds())
		case marker.WhiteTriangle():
			markerOddsRangeMap[types.WhiteTriangle] = getOddsRange(odds.Odds())
		case marker.Star():
			markerOddsRangeMap[types.Star] = getOddsRange(odds.Odds())
		case marker.Check():
			markerOddsRangeMap[types.Check] = getOddsRange(odds.Odds())
		}
	}

	return markerOddsRangeMap
}

func (s *spreadSheetService) GetCellColor(
	ctx context.Context,
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
