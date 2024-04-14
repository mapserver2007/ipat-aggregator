package infrastructure

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"github.com/shopspring/decimal"
	"google.golang.org/api/sheets/v4"
	"log"
	"strconv"
	"strings"
)

const (
	spreadSheetTrioAnalysisFileName = "spreadsheet_trio_analysis.json"
)

type spreadSheetTrioAnalysisRepository struct {
	client             *sheets.Service
	spreadSheetConfigs []*spreadsheet_entity.SpreadSheetConfig
	spreadSheetService service.SpreadSheetService
}

func NewSpreadSheetTrioAnalysisRepository(
	spreadSheetService service.SpreadSheetService,
) (repository.SpreadSheetTrioAnalysisRepository, error) {
	ctx := context.Background()
	client, spreadSheetConfigs, err := getSpreadSheetConfigs(ctx, spreadSheetTrioAnalysisFileName)
	if err != nil {
		return nil, err
	}

	return &spreadSheetTrioAnalysisRepository{
		client:             client,
		spreadSheetConfigs: spreadSheetConfigs,
		spreadSheetService: spreadSheetService,
	}, nil
}

func (s *spreadSheetTrioAnalysisRepository) Write(
	ctx context.Context,
	analysisData *spreadsheet_entity.AnalysisData,
	races []*data_cache_entity.Race,
	odds []*data_cache_entity.Odds,
) error {
	for idx, spreadSheetConfig := range s.spreadSheetConfigs {
		pivotalMarker, _ := types.NewMarker(idx + 1)
		log.Println(ctx, fmt.Sprintf("write marker %s-印-印 analysis start", pivotalMarker.String()))

		markerCombinationMap := analysisData.MarkerCombinationFilterMap()

		// レース結果のオッズ
		winRaceResultOddsMap := map[types.RaceId][]*spreadsheet_entity.Odds{}  // 同着を考慮してslice
		trioRaceResultOddsMap := map[types.RaceId][]*spreadsheet_entity.Odds{} // 同着を考慮してslice

		for _, race := range races {
			for _, payoutResult := range race.PayoutResults() {
				switch payoutResult.TicketType() {
				case types.Win:
					winRaceResultOddsMap[race.RaceId()] = make([]*spreadsheet_entity.Odds, 0)
					// 軸を取得する目的なのでレース結果の1着の情報を保持(同着の場合は複数あり)
					for _, result := range race.RaceResults() {
						if result.OrderNo() == 1 {
							winRaceResultOddsMap[race.RaceId()] = append(winRaceResultOddsMap[race.RaceId()], spreadsheet_entity.NewOdds(
								payoutResult.TicketType(),
								result.Odds(),
								types.NewBetNumber(strconv.Itoa(result.HorseNumber())),
							))
						}
					}
				case types.Trio:
					trioRaceResultOddsMap[race.RaceId()] = make([]*spreadsheet_entity.Odds, 0)
					// オッズ計算のため3連複払戻結果の情報を保持(同着の場合は複数あり)
					for i := 0; i < len(payoutResult.Numbers()); i++ {
						trioRaceResultOddsMap[race.RaceId()] = append(trioRaceResultOddsMap[race.RaceId()], spreadsheet_entity.NewOdds(
							payoutResult.TicketType(),
							payoutResult.Odds()[i],
							payoutResult.Numbers()[i],
						))
					}
				}
			}
		}

		trioMarkerOddsMap := map[types.RaceId][]*spreadsheet_entity.Odds{}
		for _, o := range odds {
			if o.TicketType().OriginTicketType() == types.Trio {
				if _, ok := trioMarkerOddsMap[o.RaceId()]; !ok {
					trioMarkerOddsMap[o.RaceId()] = make([]*spreadsheet_entity.Odds, 0, 20)
				}
				trioMarkerOddsMap[o.RaceId()] = append(trioMarkerOddsMap[o.RaceId()], spreadsheet_entity.NewOdds(
					o.TicketType(),
					o.Odds(),
					o.Number(),
				))
			}
		}

		var valuesList [][][]interface{}
		for filterGroupIndex, f := range analysisData.Filters() {
			// 軸に対するレース単位のオッズ幅ごとの的中回数
			raceHitOddsRangeCountMap, err := s.getRaceHitOddsRangeCountMap(
				ctx,
				markerCombinationMap[f],
				winRaceResultOddsMap,
				trioRaceResultOddsMap,
			)
			if err != nil {
				return err
			}

			if f == filter.DirtSmallNumberOfHorses && idx == 4 {
				fmt.Println("wata")
			}

			// 軸に対するレース単位のオッズ幅ごとの出現回数
			raceOddsRangeCountMap, err := s.getRaceOddsRangeCountMap(
				ctx,
				pivotalMarker,
				markerCombinationMap[f],
				winRaceResultOddsMap,
				trioRaceResultOddsMap,
			)
			if err != nil {
				return err
			}

			// 軸に対する印単位のオッズ幅ごとの全回数
			markerAllOddsRangeCountMap, err := s.getMarkerAllOddsRangeCountMap(
				ctx,
				pivotalMarker,
				markerCombinationMap[f],
				winRaceResultOddsMap,
				trioMarkerOddsMap,
			)
			if err != nil {
				return err
			}

			// 軸に対するオッズ幅ごとの的中時オッズの合計
			pivotalMarkerHitTotalOddsMap, err := s.getPivotalMarkerHitTotalOddsMap(
				ctx,
				pivotalMarker,
				markerCombinationMap[f],
				winRaceResultOddsMap,
				trioRaceResultOddsMap,
			)
			if err != nil {
				return err
			}

			defaultValuesList := s.createDefaultValuesList()
			position := len(defaultValuesList) * filterGroupIndex
			valuesList = append(valuesList, defaultValuesList...)

			// 率計算
			for i := position; i < len(defaultValuesList)+position; i++ {
				var (
					hitCountMap  map[types.OddsRangeType]int
					allCountMap  map[types.OddsRangeType]int
					raceCountMap map[types.OddsRangeType]int
					hitOddsMap   map[types.OddsRangeType]decimal.Decimal
				)

				switch i - position {
				case 0:
					valuesList[i][0][4] = fmt.Sprintf("軸選択的中率 フィルタ条件: %s", f.String())
					valuesList[i][0][12] = fmt.Sprintf("軸選択回収率 フィルタ条件: %s", f.String())
					valuesList[i][0][20] = fmt.Sprintf("%s軸+印決着回数 フィルタ条件: %s", pivotalMarker.String(), f.String())
					valuesList[i][0][28] = fmt.Sprintf("%s軸+印/無含む決着回数) フィルタ条件: %s", pivotalMarker.String(), f.String())
					continue
				case 1:
					continue
				case 2:
					hitCountMap = map[types.OddsRangeType]int{}
					for _, oddsRange := range raceHitOddsRangeCountMap[pivotalMarker] {
						hitCountMap[types.TrioOddsRange1] += oddsRange[types.TrioOddsRange1]
						hitCountMap[types.TrioOddsRange2] += oddsRange[types.TrioOddsRange2]
						hitCountMap[types.TrioOddsRange3] += oddsRange[types.TrioOddsRange3]
						hitCountMap[types.TrioOddsRange4] += oddsRange[types.TrioOddsRange4]
						hitCountMap[types.TrioOddsRange5] += oddsRange[types.TrioOddsRange5]
						hitCountMap[types.TrioOddsRange6] += oddsRange[types.TrioOddsRange6]
						hitCountMap[types.TrioOddsRange7] += oddsRange[types.TrioOddsRange7]
						hitCountMap[types.TrioOddsRange8] += oddsRange[types.TrioOddsRange8]
					}
					allCountMap = map[types.OddsRangeType]int{}
					for _, oddsRange := range markerAllOddsRangeCountMap {
						allCountMap[types.TrioOddsRange1] += oddsRange[types.TrioOddsRange1]
						allCountMap[types.TrioOddsRange2] += oddsRange[types.TrioOddsRange2]
						allCountMap[types.TrioOddsRange3] += oddsRange[types.TrioOddsRange3]
						allCountMap[types.TrioOddsRange4] += oddsRange[types.TrioOddsRange4]
						allCountMap[types.TrioOddsRange5] += oddsRange[types.TrioOddsRange5]
						allCountMap[types.TrioOddsRange6] += oddsRange[types.TrioOddsRange6]
						allCountMap[types.TrioOddsRange7] += oddsRange[types.TrioOddsRange7]
						allCountMap[types.TrioOddsRange8] += oddsRange[types.TrioOddsRange8]
					}
					raceCountMap = map[types.OddsRangeType]int{}
					for _, oddsRange := range raceOddsRangeCountMap {
						raceCountMap[types.TrioOddsRange1] += oddsRange[types.TrioOddsRange1]
						raceCountMap[types.TrioOddsRange2] += oddsRange[types.TrioOddsRange2]
						raceCountMap[types.TrioOddsRange3] += oddsRange[types.TrioOddsRange3]
						raceCountMap[types.TrioOddsRange4] += oddsRange[types.TrioOddsRange4]
						raceCountMap[types.TrioOddsRange5] += oddsRange[types.TrioOddsRange5]
						raceCountMap[types.TrioOddsRange6] += oddsRange[types.TrioOddsRange6]
						raceCountMap[types.TrioOddsRange7] += oddsRange[types.TrioOddsRange7]
						raceCountMap[types.TrioOddsRange8] += oddsRange[types.TrioOddsRange8]
					}
					hitOddsMap = map[types.OddsRangeType]decimal.Decimal{}
					for _, oddsRange := range pivotalMarkerHitTotalOddsMap {
						hitOddsMap[types.TrioOddsRange1] = hitOddsMap[types.TrioOddsRange1].Add(oddsRange[types.TrioOddsRange1])
						hitOddsMap[types.TrioOddsRange2] = hitOddsMap[types.TrioOddsRange2].Add(oddsRange[types.TrioOddsRange2])
						hitOddsMap[types.TrioOddsRange3] = hitOddsMap[types.TrioOddsRange3].Add(oddsRange[types.TrioOddsRange3])
						hitOddsMap[types.TrioOddsRange4] = hitOddsMap[types.TrioOddsRange4].Add(oddsRange[types.TrioOddsRange4])
						hitOddsMap[types.TrioOddsRange5] = hitOddsMap[types.TrioOddsRange5].Add(oddsRange[types.TrioOddsRange5])
						hitOddsMap[types.TrioOddsRange6] = hitOddsMap[types.TrioOddsRange6].Add(oddsRange[types.TrioOddsRange6])
						hitOddsMap[types.TrioOddsRange7] = hitOddsMap[types.TrioOddsRange7].Add(oddsRange[types.TrioOddsRange7])
						hitOddsMap[types.TrioOddsRange8] = hitOddsMap[types.TrioOddsRange8].Add(oddsRange[types.TrioOddsRange8])
					}
					valuesList[i][0][0] = "単全部"
				case 3:
					hitCountMap = raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange1]
					allCountMap = markerAllOddsRangeCountMap[types.WinOddsRange1]
					raceCountMap = raceOddsRangeCountMap[types.WinOddsRange1]
					hitOddsMap = pivotalMarkerHitTotalOddsMap[types.WinOddsRange1]
					valuesList[i][0][0] = fmt.Sprintf("単%s", types.WinOddsRange1.String())
				case 4:
					hitCountMap = raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange2]
					allCountMap = markerAllOddsRangeCountMap[types.WinOddsRange2]
					raceCountMap = raceOddsRangeCountMap[types.WinOddsRange2]
					hitOddsMap = pivotalMarkerHitTotalOddsMap[types.WinOddsRange2]
					valuesList[i][0][0] = fmt.Sprintf("単%s", types.WinOddsRange2.String())
				case 5:
					hitCountMap = raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange3]
					allCountMap = markerAllOddsRangeCountMap[types.WinOddsRange3]
					raceCountMap = raceOddsRangeCountMap[types.WinOddsRange3]
					hitOddsMap = pivotalMarkerHitTotalOddsMap[types.WinOddsRange3]
					valuesList[i][0][0] = fmt.Sprintf("単%s", types.WinOddsRange3.String())
				case 6:
					hitCountMap = raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange4]
					allCountMap = markerAllOddsRangeCountMap[types.WinOddsRange4]
					raceCountMap = raceOddsRangeCountMap[types.WinOddsRange4]
					hitOddsMap = pivotalMarkerHitTotalOddsMap[types.WinOddsRange4]
					valuesList[i][0][0] = fmt.Sprintf("単%s", types.WinOddsRange4.String())
				case 7:
					hitCountMap = raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange5]
					allCountMap = markerAllOddsRangeCountMap[types.WinOddsRange5]
					raceCountMap = raceOddsRangeCountMap[types.WinOddsRange5]
					hitOddsMap = pivotalMarkerHitTotalOddsMap[types.WinOddsRange5]
					valuesList[i][0][0] = fmt.Sprintf("単%s", types.WinOddsRange5.String())
				case 8:
					hitCountMap = raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange6]
					allCountMap = markerAllOddsRangeCountMap[types.WinOddsRange6]
					raceCountMap = raceOddsRangeCountMap[types.WinOddsRange6]
					hitOddsMap = pivotalMarkerHitTotalOddsMap[types.WinOddsRange6]
					valuesList[i][0][0] = fmt.Sprintf("単%s", types.WinOddsRange6.String())
				case 9:
					hitCountMap = raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange7]
					allCountMap = markerAllOddsRangeCountMap[types.WinOddsRange7]
					raceCountMap = raceOddsRangeCountMap[types.WinOddsRange7]
					hitOddsMap = pivotalMarkerHitTotalOddsMap[types.WinOddsRange7]
					valuesList[i][0][0] = fmt.Sprintf("単%s", types.WinOddsRange7.String())
				case 10:
					hitCountMap = raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange8]
					allCountMap = markerAllOddsRangeCountMap[types.WinOddsRange8]
					raceCountMap = raceOddsRangeCountMap[types.WinOddsRange8]
					hitOddsMap = pivotalMarkerHitTotalOddsMap[types.WinOddsRange8]
					valuesList[i][0][0] = fmt.Sprintf("単%s", types.WinOddsRange8.String())
				}
				allCount := 0
				for _, count := range allCountMap {
					allCount += count
				}
				valuesList[i][0][1] = allCount
				hitCount := 0
				for _, count := range hitCountMap {
					hitCount += count
				}
				valuesList[i][0][2] = hitCount
				raceCount := 0
				for _, count := range raceCountMap {
					raceCount += count
				}
				valuesList[i][0][3] = raceCount
				valuesList[i][0][4] = HitRateFormat(hitCountMap[types.TrioOddsRange1], allCountMap[types.TrioOddsRange1])
				valuesList[i][0][5] = HitRateFormat(hitCountMap[types.TrioOddsRange2], allCountMap[types.TrioOddsRange2])
				valuesList[i][0][6] = HitRateFormat(hitCountMap[types.TrioOddsRange3], allCountMap[types.TrioOddsRange3])
				valuesList[i][0][7] = HitRateFormat(hitCountMap[types.TrioOddsRange4], allCountMap[types.TrioOddsRange4])
				valuesList[i][0][8] = HitRateFormat(hitCountMap[types.TrioOddsRange5], allCountMap[types.TrioOddsRange5])
				valuesList[i][0][9] = HitRateFormat(hitCountMap[types.TrioOddsRange6], allCountMap[types.TrioOddsRange6])
				valuesList[i][0][10] = HitRateFormat(hitCountMap[types.TrioOddsRange7], allCountMap[types.TrioOddsRange7])
				valuesList[i][0][11] = HitRateFormat(hitCountMap[types.TrioOddsRange8], allCountMap[types.TrioOddsRange8])
				valuesList[i][0][12] = PayoutRateFormat(hitOddsMap[types.TrioOddsRange1].InexactFloat64(), allCountMap[types.TrioOddsRange1])
				valuesList[i][0][13] = PayoutRateFormat(hitOddsMap[types.TrioOddsRange2].InexactFloat64(), allCountMap[types.TrioOddsRange2])
				valuesList[i][0][14] = PayoutRateFormat(hitOddsMap[types.TrioOddsRange3].InexactFloat64(), allCountMap[types.TrioOddsRange3])
				valuesList[i][0][15] = PayoutRateFormat(hitOddsMap[types.TrioOddsRange4].InexactFloat64(), allCountMap[types.TrioOddsRange4])
				valuesList[i][0][16] = PayoutRateFormat(hitOddsMap[types.TrioOddsRange5].InexactFloat64(), allCountMap[types.TrioOddsRange5])
				valuesList[i][0][17] = PayoutRateFormat(hitOddsMap[types.TrioOddsRange6].InexactFloat64(), allCountMap[types.TrioOddsRange6])
				valuesList[i][0][18] = PayoutRateFormat(hitOddsMap[types.TrioOddsRange7].InexactFloat64(), allCountMap[types.TrioOddsRange7])
				valuesList[i][0][19] = PayoutRateFormat(hitOddsMap[types.TrioOddsRange8].InexactFloat64(), allCountMap[types.TrioOddsRange8])
				valuesList[i][0][20] = hitCountMap[types.TrioOddsRange1]
				valuesList[i][0][21] = hitCountMap[types.TrioOddsRange2]
				valuesList[i][0][22] = hitCountMap[types.TrioOddsRange3]
				valuesList[i][0][23] = hitCountMap[types.TrioOddsRange4]
				valuesList[i][0][24] = hitCountMap[types.TrioOddsRange5]
				valuesList[i][0][25] = hitCountMap[types.TrioOddsRange6]
				valuesList[i][0][26] = hitCountMap[types.TrioOddsRange7]
				valuesList[i][0][27] = hitCountMap[types.TrioOddsRange8]
				valuesList[i][0][28] = allCountMap[types.TrioOddsRange1]
				valuesList[i][0][29] = allCountMap[types.TrioOddsRange2]
				valuesList[i][0][30] = allCountMap[types.TrioOddsRange3]
				valuesList[i][0][31] = allCountMap[types.TrioOddsRange4]
				valuesList[i][0][32] = allCountMap[types.TrioOddsRange5]
				valuesList[i][0][33] = allCountMap[types.TrioOddsRange6]
				valuesList[i][0][34] = allCountMap[types.TrioOddsRange7]
				valuesList[i][0][35] = allCountMap[types.TrioOddsRange8]
			}
		}

		var values [][]interface{}
		for _, v := range valuesList {
			values = append(values, v...)
		}
		writeRange := fmt.Sprintf("%s!%s", spreadSheetConfig.SheetName(), fmt.Sprintf("A1"))
		_, err := s.client.Spreadsheets.Values.Update(spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
			Values: values,
		}).ValueInputOption("USER_ENTERED").Do()
		if err != nil {
			return err
		}

		log.Println(ctx, fmt.Sprintf("write marker %s-印-印 analysis end", pivotalMarker.String()))
	}

	return nil
}

func (s *spreadSheetTrioAnalysisRepository) createDefaultValuesList() [][][]interface{} {
	valuesList := make([][][]interface{}, 0)
	for i := 0; i < 11; i++ {
		if i == 1 {
			valuesList = append(valuesList, [][]interface{}{
				{
					"",
					"投票回数",
					"的中回数",
					"レース数",
					types.TrioOddsRange1.String(),
					types.TrioOddsRange2.String(),
					types.TrioOddsRange3.String(),
					types.TrioOddsRange4.String(),
					types.TrioOddsRange5.String(),
					types.TrioOddsRange6.String(),
					types.TrioOddsRange7.String(),
					types.TrioOddsRange8.String(),
					types.TrioOddsRange1.String(),
					types.TrioOddsRange2.String(),
					types.TrioOddsRange3.String(),
					types.TrioOddsRange4.String(),
					types.TrioOddsRange5.String(),
					types.TrioOddsRange6.String(),
					types.TrioOddsRange7.String(),
					types.TrioOddsRange8.String(),
					types.TrioOddsRange1.String(),
					types.TrioOddsRange2.String(),
					types.TrioOddsRange3.String(),
					types.TrioOddsRange4.String(),
					types.TrioOddsRange5.String(),
					types.TrioOddsRange6.String(),
					types.TrioOddsRange7.String(),
					types.TrioOddsRange8.String(),
					types.TrioOddsRange1.String(),
					types.TrioOddsRange2.String(),
					types.TrioOddsRange3.String(),
					types.TrioOddsRange4.String(),
					types.TrioOddsRange5.String(),
					types.TrioOddsRange6.String(),
					types.TrioOddsRange7.String(),
					types.TrioOddsRange8.String(),
				},
			})
		} else {
			valuesList = append(valuesList, [][]interface{}{
				{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
			})
		}

	}

	return valuesList
}

func (s *spreadSheetTrioAnalysisRepository) getRaceHitOddsRangeCountMap(
	ctx context.Context,
	markerCombinationAnalysisMap map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis,
	winRaceOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
	trioRaceOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
) (map[types.Marker]map[types.OddsRangeType]map[types.OddsRangeType]int, error) {
	hitRaceOddsRangeCountMap := map[types.Marker]map[types.OddsRangeType]map[types.OddsRangeType]int{}
	pivotalMarkers := []types.Marker{
		types.Favorite,
		types.Rival,
		types.BrackTriangle,
		types.WhiteTriangle,
		types.Star,
		types.Check,
	}

	hitMarkerCalculablesMap := map[types.Marker][]*analysis_entity.Calculable{}
	for _, marker := range pivotalMarkers {
		for markerCombinationId, markerCombinationAnalysis := range markerCombinationAnalysisMap {
			if markerCombinationId.TicketType().OriginTicketType() != types.Trio {
				continue
			}
			if !strings.Contains(strconv.Itoa(markerCombinationId.Value()%1000), strconv.Itoa(marker.Value())) {
				continue
			}
			for _, calculable := range markerCombinationAnalysis.Calculables() {
				if calculable.IsHit() {
					hitMarkerCalculablesMap[marker] = append(hitMarkerCalculablesMap[marker], calculable)
				}
			}
		}
	}

	for marker, calculables := range hitMarkerCalculablesMap {
		hitRaceOddsRangeCountMap[marker] = map[types.OddsRangeType]map[types.OddsRangeType]int{}
		for _, calculable := range calculables {
			pivotalOddsList, ok := winRaceOddsMap[calculable.RaceId()]
			if !ok {
				return nil, fmt.Errorf("winRaceOdds not found. raceId: %s", calculable.RaceId())
			}
			trioOddsList, ok := trioRaceOddsMap[calculable.RaceId()]
			if !ok {
				return nil, fmt.Errorf("trioRaceOdds not found. raceId: %s", calculable.RaceId())
			}

			for _, pivotalOdds := range pivotalOddsList {
				var pivotalMarkerOddsRange types.OddsRangeType
				odds := pivotalOdds.Odds().InexactFloat64()
				if odds >= 1.0 && odds <= 1.5 {
					pivotalMarkerOddsRange = types.WinOddsRange1
				} else if odds >= 1.6 && odds <= 2.0 {
					pivotalMarkerOddsRange = types.WinOddsRange2
				} else if odds >= 2.1 && odds <= 2.9 {
					pivotalMarkerOddsRange = types.WinOddsRange3
				} else if odds >= 3.0 && odds <= 4.9 {
					pivotalMarkerOddsRange = types.WinOddsRange4
				} else if odds >= 5.0 && odds <= 9.9 {
					pivotalMarkerOddsRange = types.WinOddsRange5
				} else if odds >= 10.0 && odds <= 19.9 {
					pivotalMarkerOddsRange = types.WinOddsRange6
				} else if odds >= 20.0 && odds <= 49.9 {
					pivotalMarkerOddsRange = types.WinOddsRange7
				} else if odds >= 50.0 {
					pivotalMarkerOddsRange = types.WinOddsRange8
				}

				if _, ok := hitRaceOddsRangeCountMap[marker][pivotalMarkerOddsRange]; !ok {
					hitRaceOddsRangeCountMap[marker][pivotalMarkerOddsRange] = map[types.OddsRangeType]int{}
				}

				for _, trioOdds := range trioOddsList {
					odds = trioOdds.Odds().InexactFloat64()
					if odds >= 1.0 && odds <= 9.9 {
						hitRaceOddsRangeCountMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange1]++
					} else if odds >= 10.0 && odds <= 19.9 {
						hitRaceOddsRangeCountMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange2]++
					} else if odds >= 20.0 && odds <= 29.9 {
						hitRaceOddsRangeCountMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange3]++
					} else if odds >= 30.0 && odds <= 49.9 {
						hitRaceOddsRangeCountMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange4]++
					} else if odds >= 50.0 && odds <= 99.9 {
						hitRaceOddsRangeCountMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange5]++
					} else if odds >= 100.0 && odds <= 299.9 {
						hitRaceOddsRangeCountMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange6]++
					} else if odds >= 300.0 && odds <= 499.9 {
						hitRaceOddsRangeCountMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange7]++
					} else if odds >= 500.0 {
						hitRaceOddsRangeCountMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange8]++
					}
				}
			}
		}
	}

	return hitRaceOddsRangeCountMap, nil
}

func (s *spreadSheetTrioAnalysisRepository) getRaceOddsRangeCountMap(
	ctx context.Context,
	marker types.Marker,
	markerCombinationAnalysisMap map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis,
	winRaceOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
	trioRaceOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
) (map[types.OddsRangeType]map[types.OddsRangeType]int, error) {
	raceOddsRangeCountMap := map[types.OddsRangeType]map[types.OddsRangeType]int{}
	var calculables []*analysis_entity.Calculable
	for markerCombinationId, markerCombinationAnalysis := range markerCombinationAnalysisMap {
		if markerCombinationId.TicketType().OriginTicketType() != types.Trio {
			continue
		}
		if !strings.Contains(strconv.Itoa(markerCombinationId.Value()%1000), strconv.Itoa(marker.Value())) {
			continue
		}
		if len(markerCombinationAnalysis.Calculables()) == 0 {
			continue
		}
		calculables = append(calculables, markerCombinationAnalysis.Calculables()...)
	}

	for _, calculable := range calculables {
		pivotalOddsList, ok := winRaceOddsMap[calculable.RaceId()]
		if !ok {
			return nil, fmt.Errorf("winRaceOdds not found. raceId: %s", calculable.RaceId())
		}
		trioOddsList, ok := trioRaceOddsMap[calculable.RaceId()]
		if !ok {
			return nil, fmt.Errorf("trioRaceOdds not found. raceId: %s", calculable.RaceId())
		}

		for _, pivotalOdds := range pivotalOddsList {
			var pivotalMarkerOddsRange types.OddsRangeType
			odds := pivotalOdds.Odds().InexactFloat64()
			if odds >= 1.0 && odds <= 1.5 {
				pivotalMarkerOddsRange = types.WinOddsRange1
			} else if odds >= 1.6 && odds <= 2.0 {
				pivotalMarkerOddsRange = types.WinOddsRange2
			} else if odds >= 2.1 && odds <= 2.9 {
				pivotalMarkerOddsRange = types.WinOddsRange3
			} else if odds >= 3.0 && odds <= 4.9 {
				pivotalMarkerOddsRange = types.WinOddsRange4
			} else if odds >= 5.0 && odds <= 9.9 {
				pivotalMarkerOddsRange = types.WinOddsRange5
			} else if odds >= 10.0 && odds <= 19.9 {
				pivotalMarkerOddsRange = types.WinOddsRange6
			} else if odds >= 20.0 && odds <= 49.9 {
				pivotalMarkerOddsRange = types.WinOddsRange7
			} else if odds >= 50.0 {
				pivotalMarkerOddsRange = types.WinOddsRange8
			}

			if _, ok := raceOddsRangeCountMap[pivotalMarkerOddsRange]; !ok {
				raceOddsRangeCountMap[pivotalMarkerOddsRange] = map[types.OddsRangeType]int{}
			}

			for _, trioOdds := range trioOddsList {
				odds = trioOdds.Odds().InexactFloat64()
				if odds >= 1.0 && odds <= 9.9 {
					raceOddsRangeCountMap[pivotalMarkerOddsRange][types.TrioOddsRange1]++
				} else if odds >= 10.0 && odds <= 19.9 {
					raceOddsRangeCountMap[pivotalMarkerOddsRange][types.TrioOddsRange2]++
				} else if odds >= 20.0 && odds <= 29.9 {
					raceOddsRangeCountMap[pivotalMarkerOddsRange][types.TrioOddsRange3]++
				} else if odds >= 30.0 && odds <= 49.9 {
					raceOddsRangeCountMap[pivotalMarkerOddsRange][types.TrioOddsRange4]++
				} else if odds >= 50.0 && odds <= 99.9 {
					raceOddsRangeCountMap[pivotalMarkerOddsRange][types.TrioOddsRange5]++
				} else if odds >= 100.0 && odds <= 299.9 {
					raceOddsRangeCountMap[pivotalMarkerOddsRange][types.TrioOddsRange6]++
				} else if odds >= 300.0 && odds <= 499.9 {
					raceOddsRangeCountMap[pivotalMarkerOddsRange][types.TrioOddsRange7]++
				} else if odds >= 500.0 {
					raceOddsRangeCountMap[pivotalMarkerOddsRange][types.TrioOddsRange8]++
				}
			}
		}
	}

	return raceOddsRangeCountMap, nil
}

func (s *spreadSheetTrioAnalysisRepository) getMarkerAllOddsRangeCountMap(
	ctx context.Context,
	marker types.Marker,
	markerCombinationAnalysisMap map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis,
	winRaceOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
	trioMarkerOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
) (map[types.OddsRangeType]map[types.OddsRangeType]int, error) {
	pivotalMarkerAllCountOddsRangeMap := map[types.OddsRangeType]map[types.OddsRangeType]int{}
	var calculables []*analysis_entity.Calculable

	for markerCombinationId, markerCombinationAnalysis := range markerCombinationAnalysisMap {
		if markerCombinationId.TicketType().OriginTicketType() != types.Trio {
			continue
		}
		if !strings.Contains(strconv.Itoa(markerCombinationId.Value()%1000), strconv.Itoa(marker.Value())) {
			continue
		}
		if len(markerCombinationAnalysis.Calculables()) == 0 {
			continue
		}
		calculables = append(calculables, markerCombinationAnalysis.Calculables()...)
	}

	for _, calculable := range calculables {
		markerOdds, ok := trioMarkerOddsMap[calculable.RaceId()]
		if !ok {
			return nil, fmt.Errorf("trioMarkerOdds not found. raceId: %s", calculable.RaceId())
		}

		pivotalOddsList, ok := winRaceOddsMap[calculable.RaceId()]
		if !ok {
			return nil, fmt.Errorf("winRaceOdds not found. raceId: %s", calculable.RaceId())
		}

		for _, pivotalOdds := range pivotalOddsList {
			var pivotalMarkerOddsRange types.OddsRangeType
			odds := pivotalOdds.Odds().InexactFloat64()
			if odds >= 1.0 && odds <= 1.5 {
				pivotalMarkerOddsRange = types.WinOddsRange1
			} else if odds >= 1.6 && odds <= 2.0 {
				pivotalMarkerOddsRange = types.WinOddsRange2
			} else if odds >= 2.1 && odds <= 2.9 {
				pivotalMarkerOddsRange = types.WinOddsRange3
			} else if odds >= 3.0 && odds <= 4.9 {
				pivotalMarkerOddsRange = types.WinOddsRange4
			} else if odds >= 5.0 && odds <= 9.9 {
				pivotalMarkerOddsRange = types.WinOddsRange5
			} else if odds >= 10.0 && odds <= 19.9 {
				pivotalMarkerOddsRange = types.WinOddsRange6
			} else if odds >= 20.0 && odds <= 49.9 {
				pivotalMarkerOddsRange = types.WinOddsRange7
			} else if odds >= 50.0 {
				pivotalMarkerOddsRange = types.WinOddsRange8
			}

			if _, ok := pivotalMarkerAllCountOddsRangeMap[pivotalMarkerOddsRange]; !ok {
				pivotalMarkerAllCountOddsRangeMap[pivotalMarkerOddsRange] = map[types.OddsRangeType]int{}
			}

			for _, trioOdds := range markerOdds {
				odds = trioOdds.Odds().InexactFloat64()
				if odds >= 1.0 && odds <= 9.9 {
					pivotalMarkerAllCountOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange1]++
				} else if odds >= 10.0 && odds <= 19.9 {
					pivotalMarkerAllCountOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange2]++
				} else if odds >= 20.0 && odds <= 29.9 {
					pivotalMarkerAllCountOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange3]++
				} else if odds >= 30.0 && odds <= 49.9 {
					pivotalMarkerAllCountOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange4]++
				} else if odds >= 50.0 && odds <= 99.9 {
					pivotalMarkerAllCountOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange5]++
				} else if odds >= 100.0 && odds <= 299.9 {
					pivotalMarkerAllCountOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange6]++
				} else if odds >= 300.0 && odds <= 499.9 {
					pivotalMarkerAllCountOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange7]++
				} else if odds >= 500.0 {
					pivotalMarkerAllCountOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange8]++
				}
			}
		}
	}

	return pivotalMarkerAllCountOddsRangeMap, nil
}

func (s *spreadSheetTrioAnalysisRepository) getPivotalMarkerHitTotalOddsMap(
	ctx context.Context,
	marker types.Marker,
	markerCombinationAnalysisMap map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis,
	winRaceOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
	trioRaceOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
) (map[types.OddsRangeType]map[types.OddsRangeType]decimal.Decimal, error) {
	var calculables []*analysis_entity.Calculable
	pivotalMarkerHitTotalOddsRangeMap := map[types.OddsRangeType]map[types.OddsRangeType]decimal.Decimal{}

	for markerCombinationId, markerCombinationAnalysis := range markerCombinationAnalysisMap {
		if markerCombinationId.TicketType().OriginTicketType() != types.Trio {
			continue
		}
		if !strings.Contains(strconv.Itoa(markerCombinationId.Value()%1000), strconv.Itoa(marker.Value())) {
			continue
		}
		if len(markerCombinationAnalysis.Calculables()) == 0 {
			continue
		}
		calculables = append(calculables, markerCombinationAnalysis.Calculables()...)
	}

	for _, calculable := range calculables {
		if !calculable.IsHit() {
			continue
		}

		pivotalOddsList, ok := winRaceOddsMap[calculable.RaceId()]
		if !ok {
			return nil, fmt.Errorf("winRaceOdds not found. raceId: %s", calculable.RaceId())
		}
		trioOddsList, ok := trioRaceOddsMap[calculable.RaceId()]
		if !ok {
			return nil, fmt.Errorf("trioRaceOdds not found. raceId: %s", calculable.RaceId())
		}

		for _, pivotalOdds := range pivotalOddsList {
			var pivotalMarkerOddsRange types.OddsRangeType
			odds := pivotalOdds.Odds().InexactFloat64()
			if odds >= 1.0 && odds <= 1.5 {
				pivotalMarkerOddsRange = types.WinOddsRange1
			} else if odds >= 1.6 && odds <= 2.0 {
				pivotalMarkerOddsRange = types.WinOddsRange2
			} else if odds >= 2.1 && odds <= 2.9 {
				pivotalMarkerOddsRange = types.WinOddsRange3
			} else if odds >= 3.0 && odds <= 4.9 {
				pivotalMarkerOddsRange = types.WinOddsRange4
			} else if odds >= 5.0 && odds <= 9.9 {
				pivotalMarkerOddsRange = types.WinOddsRange5
			} else if odds >= 10.0 && odds <= 19.9 {
				pivotalMarkerOddsRange = types.WinOddsRange6
			} else if odds >= 20.0 && odds <= 49.9 {
				pivotalMarkerOddsRange = types.WinOddsRange7
			} else if odds >= 50.0 {
				pivotalMarkerOddsRange = types.WinOddsRange8
			}

			if _, ok := pivotalMarkerHitTotalOddsRangeMap[pivotalMarkerOddsRange]; !ok {
				pivotalMarkerHitTotalOddsRangeMap[pivotalMarkerOddsRange] = map[types.OddsRangeType]decimal.Decimal{}
			}

			for _, trioOdds := range trioOddsList {
				odds = trioOdds.Odds().InexactFloat64()
				if odds >= 1.0 && odds <= 9.9 {
					pivotalMarkerHitTotalOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange1] =
						pivotalMarkerHitTotalOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange1].Add(trioOdds.Odds())
				} else if odds >= 10.0 && odds <= 19.9 {
					pivotalMarkerHitTotalOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange2] =
						pivotalMarkerHitTotalOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange2].Add(trioOdds.Odds())
				} else if odds >= 20.0 && odds <= 29.9 {
					pivotalMarkerHitTotalOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange3] =
						pivotalMarkerHitTotalOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange3].Add(trioOdds.Odds())
				} else if odds >= 30.0 && odds <= 49.9 {
					pivotalMarkerHitTotalOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange4] =
						pivotalMarkerHitTotalOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange4].Add(trioOdds.Odds())
				} else if odds >= 50.0 && odds <= 99.9 {
					pivotalMarkerHitTotalOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange5] =
						pivotalMarkerHitTotalOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange5].Add(trioOdds.Odds())
				} else if odds >= 100.0 && odds <= 299.9 {
					pivotalMarkerHitTotalOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange6] =
						pivotalMarkerHitTotalOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange6].Add(trioOdds.Odds())
				} else if odds >= 300.0 && odds <= 499.9 {
					pivotalMarkerHitTotalOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange7] =
						pivotalMarkerHitTotalOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange7].Add(trioOdds.Odds())
				} else if odds >= 500.0 {
					pivotalMarkerHitTotalOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange8] =
						pivotalMarkerHitTotalOddsRangeMap[pivotalMarkerOddsRange][types.TrioOddsRange8].Add(trioOdds.Odds())
				}
			}
		}
	}

	return pivotalMarkerHitTotalOddsRangeMap, nil
}

func (s *spreadSheetTrioAnalysisRepository) Style(
	ctx context.Context,
	analysisData *spreadsheet_entity.AnalysisData,
) error {
	var requests []*sheets.Request
	rowGroupSize := len(s.createDefaultValuesList())

	for idx, spreadSheetConfig := range s.spreadSheetConfigs {
		pivotalMarker, _ := types.NewMarker(idx + 1)
		log.Println(ctx, fmt.Sprintf("write style marker %s-印-印 analysis start", pivotalMarker.String()))

		for filterGroupIndex := range analysisData.Filters() {
			position := rowGroupSize * filterGroupIndex
			requests = append(requests, []*sheets.Request{
				{
					RepeatCell: &sheets.RepeatCellRequest{
						Fields: "userEnteredFormat.backgroundColor",
						Range: &sheets.GridRange{
							SheetId:          spreadSheetConfig.SheetId(),
							StartColumnIndex: 4,
							StartRowIndex:    int64(position),
							EndColumnIndex:   36,
							EndRowIndex:      int64(position + 1),
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
							SheetId:          spreadSheetConfig.SheetId(),
							StartColumnIndex: 4,
							StartRowIndex:    int64(position),
							EndColumnIndex:   36,
							EndRowIndex:      int64(position + 2),
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
							SheetId:          spreadSheetConfig.SheetId(),
							StartColumnIndex: 0,
							StartRowIndex:    int64(position),
							EndColumnIndex:   36,
							EndRowIndex:      int64(position + 2),
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
						Fields: "userEnteredFormat.textFormat.bold",
						Range: &sheets.GridRange{
							SheetId:          spreadSheetConfig.SheetId(),
							StartColumnIndex: 0,
							StartRowIndex:    int64(position + 1),
							EndColumnIndex:   1,
							EndRowIndex:      int64(position + 11),
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
							SheetId:          spreadSheetConfig.SheetId(),
							StartColumnIndex: 0,
							StartRowIndex:    int64(position + 1),
							EndColumnIndex:   1,
							EndRowIndex:      int64(position + 11),
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
							SheetId:          spreadSheetConfig.SheetId(),
							StartColumnIndex: 0,
							StartRowIndex:    int64(position + 1),
							EndColumnIndex:   4,
							EndRowIndex:      int64(position + 2),
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
							SheetId:          spreadSheetConfig.SheetId(),
							StartColumnIndex: 4,
							StartRowIndex:    int64(position + 1),
							EndColumnIndex:   36,
							EndRowIndex:      int64(position + 2),
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
			}...)
		}

		_, err := s.client.Spreadsheets.BatchUpdate(spreadSheetConfig.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
			Requests: requests,
		}).Do()
		if err != nil {
			return err
		}

		log.Println(ctx, fmt.Sprintf("write style marker %s-印-印 analysis end", pivotalMarker.String()))
	}

	return nil
}

func (s *spreadSheetTrioAnalysisRepository) Clear(ctx context.Context) error {
	for _, spreadSheetConfig := range s.spreadSheetConfigs {
		requests := []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "*",
					Range: &sheets.GridRange{
						SheetId:          spreadSheetConfig.SheetId(),
						StartColumnIndex: 0,
						StartRowIndex:    0,
						EndColumnIndex:   40,
						EndRowIndex:      9999,
					},
					Cell: &sheets.CellData{},
				},
			},
		}
		_, err := s.client.Spreadsheets.BatchUpdate(spreadSheetConfig.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
			Requests: requests,
		}).Do()

		if err != nil {
			return err
		}
	}

	return nil
}
