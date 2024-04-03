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
	CreateOddsRangeCountMap(ctx context.Context, analysisData *analysis_entity.Layer1, filter filter.Id) map[types.MarkerCombinationId]map[types.OddsRangeType]int
	CreateTicketTypeRaceCountMap(ctx context.Context, analysisData *analysis_entity.Layer1, filter filter.Id) map[types.TicketType]int
	CreatePredictionOdds(ctx context.Context, marker *marker_csv_entity.PredictionMarker, race *prediction_entity.Race) map[types.Marker]*prediction_entity.OddsRange
	CreateTrioMarkerCombinationAggregationData(ctx context.Context, markerCombinationIds []types.MarkerCombinationId, markerCombinationMap map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis) (map[types.MarkerCombinationId][]*spreadsheet_entity.MarkerCombinationAnalysis, error)
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
	//oddRangeCountMap := s.CreateOddsRangeCountMap(ctx, analysisData, filter)
	for markerCombinationId, data := range analysisData.MarkerCombination {
		for _, data2 := range data.RaceDate {
			for _, data3 := range data2.RaceId {
				for _, calculable := range data3 {
					switch markerCombinationId.TicketType() {
					case types.Win, types.Place:
						if _, ok := markerCombinationDataMap[markerCombinationId]; !ok {
							markerCombinationDataMap[markerCombinationId] = spreadsheet_entity.NewMarkerCombinationAnalysis()
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
					case types.Trio, types.TrioFormation, types.TrioWheelOfFirst, types.TrioWheelOfSecond, types.TrioBox:
						if _, ok := markerCombinationDataMap[markerCombinationId]; !ok {
							markerCombinationDataMap[markerCombinationId] = spreadsheet_entity.NewMarkerCombinationAnalysis()
						}
						// TODO フィルタ
						markerCombinationDataMap[markerCombinationId].AddCalculable(calculable)
					}
				}
			}
		}
	}

	return markerCombinationDataMap
}

func (s *spreadSheetService) CreateOddsRangeCountMap(
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
						switch calculable.TicketType() {
						case types.Win, types.Place:
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
						case types.Trio:
							if odds >= 1.0 && odds <= 9.9 {
								markerCombinationOddsRangeCountMap[markerCombinationId][types.TrioOddsRange1]++
							} else if odds >= 10.0 && odds <= 19.9 {
								markerCombinationOddsRangeCountMap[markerCombinationId][types.TrioOddsRange2]++
							} else if odds >= 20.0 && odds <= 29.9 {
								markerCombinationOddsRangeCountMap[markerCombinationId][types.TrioOddsRange3]++
							} else if odds >= 30.0 && odds <= 49.9 {
								markerCombinationOddsRangeCountMap[markerCombinationId][types.TrioOddsRange4]++
							} else if odds >= 50.0 && odds <= 99.9 {
								markerCombinationOddsRangeCountMap[markerCombinationId][types.TrioOddsRange5]++
							} else if odds >= 100.0 && odds <= 299.9 {
								markerCombinationOddsRangeCountMap[markerCombinationId][types.TrioOddsRange6]++
							} else if odds >= 300.0 && odds <= 499.9 {
								markerCombinationOddsRangeCountMap[markerCombinationId][types.TrioOddsRange7]++
							} else if odds >= 500.0 {
								markerCombinationOddsRangeCountMap[markerCombinationId][types.TrioOddsRange8]++
							}
						}
					}
				}
			}
		}
	}

	return markerCombinationOddsRangeCountMap
}

func (s *spreadSheetService) CreateTicketTypeRaceCountMap(
	ctx context.Context,
	analysisData *analysis_entity.Layer1,
	filter filter.Id,
) map[types.TicketType]int {
	ticketTypeRaceIdCountMap := map[types.TicketType]map[types.RaceId]bool{}
	for markerCombinationId, data := range analysisData.MarkerCombination {
		ticketType := markerCombinationId.TicketType().OriginTicketType()
		if _, ok := ticketTypeRaceIdCountMap[ticketType]; !ok {
			ticketTypeRaceIdCountMap[ticketType] = map[types.RaceId]bool{}
		}

		for _, data2 := range data.RaceDate {
			for raceId, data3 := range data2.RaceId {
				if _, ok := ticketTypeRaceIdCountMap[ticketType][raceId]; ok {
					continue
				}
				ticketTypeRaceIdCountMap[ticketType][raceId] = false
				for _, calculable := range data3 {
					match := true
					for _, f := range calculable.Filters() {
						if f&filter == 0 {
							match = false
							break
						}
					}
					if match {
						ticketTypeRaceIdCountMap[ticketType][raceId] = true
					}
				}
			}
		}
	}

	ticketTypeRaceCountMap := map[types.TicketType]int{}
	for ticketType, raceIdCountMap := range ticketTypeRaceIdCountMap {
		for _, b := range raceIdCountMap {
			if b {
				ticketTypeRaceCountMap[ticketType] += 1
			}
		}
	}

	return ticketTypeRaceCountMap
}

func (s *spreadSheetService) CreatePredictionOdds(
	ctx context.Context,
	marker *marker_csv_entity.PredictionMarker,
	race *prediction_entity.Race,
) map[types.Marker]*prediction_entity.OddsRange {
	markerOddsRangeMap := map[types.Marker]*prediction_entity.OddsRange{}
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

	containsHorseNumber := func(slice []int, value int) types.InOrder {
		for idx, v := range slice {
			if v == value {
				return types.InOrder(idx + 1)
			}
		}
		return types.OutOfPlace
	}

	for _, odds := range race.Odds() {
		inOrder := containsHorseNumber(race.RaceResultHorseNumbers(), odds.HorseNumber())
		switch odds.HorseNumber() {
		case marker.Favorite():
			markerOddsRangeMap[types.Favorite] = prediction_entity.NewOddsRange(getOddsRange(odds.Odds()), inOrder)
		case marker.Rival():
			markerOddsRangeMap[types.Rival] = prediction_entity.NewOddsRange(getOddsRange(odds.Odds()), inOrder)
		case marker.BrackTriangle():
			markerOddsRangeMap[types.BrackTriangle] = prediction_entity.NewOddsRange(getOddsRange(odds.Odds()), inOrder)
		case marker.WhiteTriangle():
			markerOddsRangeMap[types.WhiteTriangle] = prediction_entity.NewOddsRange(getOddsRange(odds.Odds()), inOrder)
		case marker.Star():
			markerOddsRangeMap[types.Star] = prediction_entity.NewOddsRange(getOddsRange(odds.Odds()), inOrder)
		case marker.Check():
			markerOddsRangeMap[types.Check] = prediction_entity.NewOddsRange(getOddsRange(odds.Odds()), inOrder)
		}
	}

	return markerOddsRangeMap
}

// Deprecated
// CreateHitTrioMarkerCombinationAggregationData 3連複の各印の組合せ(的中)を表示用の印に再集計する
func (s *spreadSheetService) CreateTrioMarkerCombinationAggregationData(
	ctx context.Context,
	markerCombinationIds []types.MarkerCombinationId,
	markerCombinationMap map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis,
) (map[types.MarkerCombinationId][]*spreadsheet_entity.MarkerCombinationAnalysis, error) {
	aggregationMarkerCombinationIds := []types.MarkerCombinationId{
		types.MarkerCombinationId(6100), // ◎-印-印
		types.MarkerCombinationId(6200), // ◯-印-印
		types.MarkerCombinationId(6300), // ▲-印-印
		types.MarkerCombinationId(6400), // △-印-印
		types.MarkerCombinationId(6500), // ☆-印-印
		types.MarkerCombinationId(6600), // ✓-印-印
	}

	aggregationAnalysisListMap := map[types.MarkerCombinationId][]*spreadsheet_entity.MarkerCombinationAnalysis{}
	aggregationRaceCountOddsRangeMap := map[types.MarkerCombinationId][]map[types.OddsRangeType]int{}
	for _, markerCombinationId := range markerCombinationIds {
		if markerCombinationId.TicketType() == types.Trio {
			for _, aggregationMarkerCombinationId := range aggregationMarkerCombinationIds {
				if _, ok := aggregationAnalysisListMap[aggregationMarkerCombinationId]; !ok {
					aggregationAnalysisListMap[aggregationMarkerCombinationId] = make([]*spreadsheet_entity.MarkerCombinationAnalysis, 0)
				}
				if _, ok := aggregationRaceCountOddsRangeMap[aggregationMarkerCombinationId]; !ok {
					aggregationRaceCountOddsRangeMap[aggregationMarkerCombinationId] = make([]map[types.OddsRangeType]int, 0)
				}

				var rawMarkerCombinationIds []int
				switch aggregationMarkerCombinationId.Value() {
				case 6100:
					rawMarkerCombinationIds = []int{6123, 6124, 6125, 6126, 6134, 6135, 6136, 6145, 6146, 6156}
				case 6200:
					rawMarkerCombinationIds = []int{6123, 6124, 6125, 6126, 6234, 6235, 6236, 6245, 6246, 6256}
				case 6300:
					rawMarkerCombinationIds = []int{6123, 6134, 6135, 6136, 6234, 6235, 6236, 6345, 6346, 6356}
				case 6400:
					rawMarkerCombinationIds = []int{6124, 6134, 6145, 6146, 6234, 6245, 6246, 6345, 6346, 6456}
				case 6500:
					rawMarkerCombinationIds = []int{6125, 6135, 6145, 6156, 6235, 6245, 6256, 6345, 6356, 6456}
				case 6600:
					rawMarkerCombinationIds = []int{6126, 6136, 6146, 6156, 6236, 6246, 6256, 6346, 6356, 6456}
				case 6109:
					rawMarkerCombinationIds = []int{6129, 6139, 6149, 6159, 6169}
				case 6209:
					rawMarkerCombinationIds = []int{6129, 6239, 6249, 6259, 6269}
				case 6309:
					rawMarkerCombinationIds = []int{6139, 6239, 6349, 6359, 6369}
				case 6409:
					rawMarkerCombinationIds = []int{6149, 6249, 6349, 6459, 6469}
				case 6509:
					rawMarkerCombinationIds = []int{6159, 6259, 6359, 6459, 6569}
				case 6609:
					rawMarkerCombinationIds = []int{6169, 6269, 6369, 6469, 6569}
				}

				aggregationAnalysis, err := s.getTrioAggregationTrioAnalysis(
					ctx,
					rawMarkerCombinationIds,
					markerCombinationId,
					markerCombinationMap,
				)
				if err != nil {
					return nil, err
				}
				if aggregationAnalysis != nil {
					aggregationAnalysisListMap[aggregationMarkerCombinationId] = append(aggregationAnalysisListMap[aggregationMarkerCombinationId], aggregationAnalysis)
				}
			}
		}
	}

	return aggregationAnalysisListMap, nil
}

func (s *spreadSheetService) getTrioAggregationTrioAnalysis(
	ctx context.Context,
	rawMarkerCombinationIds []int,
	markerCombinationId types.MarkerCombinationId,
	markerCombinationMap map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis,
) (*spreadsheet_entity.MarkerCombinationAnalysis, error) {
	if contains(rawMarkerCombinationIds, markerCombinationId.Value()) {
		if analysisData, ok := markerCombinationMap[markerCombinationId]; ok {
			return analysisData, nil
		}
	}
	return nil, nil
}

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
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
