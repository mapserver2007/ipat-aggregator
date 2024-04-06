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
	"github.com/shopspring/decimal"
	"google.golang.org/api/sheets/v4"
	"sort"
	"strconv"
)

const (
	spreadSheetTrioAnalysisFileName = "spreadsheet_trio_analysis.json"
)

type spreadSheetTrioAnalysisRepository struct {
	client             *sheets.Service
	spreadSheetConfig  *spreadsheet_entity.SpreadSheetConfig
	spreadSheetService service.SpreadSheetService
}

func NewSpreadSheetTrioAnalysisRepository(
	spreadSheetService service.SpreadSheetService,
) (repository.SpreadSheetTrioAnalysisRepository, error) {
	ctx := context.Background()
	client, spreadSheetConfig, err := getSpreadSheetConfig(ctx, spreadSheetTrioAnalysisFileName)
	if err != nil {
		return nil, err
	}

	return &spreadSheetTrioAnalysisRepository{
		client:             client,
		spreadSheetConfig:  spreadSheetConfig,
		spreadSheetService: spreadSheetService,
	}, nil
}

func (s *spreadSheetTrioAnalysisRepository) Write(
	ctx context.Context,
	analysisData *spreadsheet_entity.AnalysisData,
	races []*data_cache_entity.Race,
	odds []*data_cache_entity.Odds,
) error {
	//markerCombinationTotalCount := 10 // 軸1頭相手5頭に対する点数は10点
	hitRateFormat := func(matchCount, raceCount int) string {
		if raceCount == 0 {
			return "-"
		}
		return fmt.Sprintf("%.2f%%", float64(matchCount)*100/float64(raceCount))
	}
	payoutRateFormat := func(payout float64, raceCount int) string {
		if raceCount == 0 {
			return "-"
		}
		return fmt.Sprintf("%.2f%%", payout*100/float64(raceCount))
	}

	//allMarkerCombinationIds := analysisData.AllMarkerCombinationIds()
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

			//payoutResultMap[race.RaceId()][payoutResult.TicketType()] = payoutResult
		}
	}

	var valuesList [][][]interface{}
	for _, f := range analysisData.Filters() {
		raceCount := analysisData.RaceCountFilterMap()[f][types.Trio]

		// 軸に対するオッズ幅ごとの的中回数
		pivotalMarkerHitOddsRangeCountMap, err := s.getPivotalMarkerHitOddsRangeCountMap(ctx, markerCombinationMap[f], winRaceResultOddsMap, trioRaceResultOddsMap)
		if err != nil {
			return err
		}

		// 軸に対するオッズ幅ごとの全回数
		pivotalMarkerAllOddsRangeCountMap, err := s.getPivotalMarkerAllOddsRangeCountMap(ctx, markerCombinationMap[f], winRaceResultOddsMap, trioRaceResultOddsMap)
		if err != nil {
			return err
		}

		// 軸に対するオッズ幅ごとの的中時オッズの合計
		pivotalMarkerHitTotalOddsMap, err := s.getPivotalMarkerHitTotalOddsMap(ctx, markerCombinationMap[f], winRaceResultOddsMap, trioRaceResultOddsMap)
		if err != nil {
			return err
		}

		// 軸に対するオッズ幅ごとの全オッズの合計
		pivotalMarkerAllTotalOddsMap, err := s.getPivotalMarkerAllTotalOddsMap(ctx, markerCombinationMap[f], winRaceResultOddsMap, trioRaceResultOddsMap)
		if err != nil {
			return err
		}

		// 軸に対する的中+不的中の全組み合わせのオッズ幅ごとの出現回数

		_ = raceCount
		_ = pivotalMarkerHitOddsRangeCountMap
		_ = pivotalMarkerAllOddsRangeCountMap
		_ = pivotalMarkerHitTotalOddsMap
		_ = pivotalMarkerAllTotalOddsMap

		//var rawPivotalMarkers []int
		//for pivotalMarker := range pivotalMarkerOddsRangeCountMap {
		//	rawPivotalMarkers = append(rawPivotalMarkers, pivotalMarker.Value())
		//}
		//sort.Ints(rawPivotalMarkers)

		aggregationMarkerIndex := 0
		for _, rawMarkerId := range []int{1, 2, 3, 4, 5, 6} {
			pivotalMarker, _ := types.NewMarker(rawMarkerId)
			defaultValuesList := s.createDefaultValuesList()
			position := len(defaultValuesList) * aggregationMarkerIndex

			//oddsRangeCountMap, ok := pivotalMarkerOddsRangeCountMap[pivotalMarker]
			//if !ok {
			//	return fmt.Errorf("oddsRangeCountMap not found: %v", pivotalMarker)
			//}
			//
			//totalOddsMap, ok := pivotalMarkerTotalOddsMap[pivotalMarker]
			//if !ok {
			//	return fmt.Errorf("oddsRangeTotalMap not found: %v", pivotalMarker)
			//}

			valuesList = append(valuesList, defaultValuesList...)
			//aggregationMarkerCombinationId := types.MarkerCombinationId(rawId)
			//aggregationAnalysisList, ok := trioAggregationAnalysisListMap[aggregationMarkerCombinationId]
			//if !ok {
			//	return fmt.Errorf("aggregationAnalysisList not found: %v", aggregationMarkerCombinationId)
			//}
			//
			//// aggregationAnalysisListの中は、集約される前の印組合せが全部入っていて、calculablesはlistになっていて合算はされていない状態
			//aggregationHitPivotalOddsRangeMap := s.createHitTrioOddsRangeMap(ctx, aggregationAnalysisList, pivotalMarker)
			////total, max, min, average, median := s.aggregationOdds(ctx, aggregationAnalysisList)
			//matchCount := 0
			//for _, pivotalOddsRange := range aggregationHitPivotalOddsRangeMap {
			//	for _, count := range pivotalOddsRange {
			//		matchCount += count
			//	}
			//}

			// 印組合せの概要集計
			//for i := position; i < len(defaultValuesList)+position; i++ {
			//	switch i - position {
			//	case 0:
			//		valuesList[i][0][0] = fmt.Sprintf("%s / %s", aggregationMarkerCombinationId.String(), f.String())
			//	case 1:
			//		valuesList[i][0][1] = raceCount
			//	case 2:
			//		valuesList[i][0][1] = matchCount
			//	case 3:
			//		valuesList[i][0][1] = rateFormatFunc(matchCount, raceCount)
			//	case 4:
			//		valuesList[i][0][1] = fmt.Sprintf("%.2f%%", (total/float64(raceCount)*100)/10) // 10点買いなので10で割る
			//	case 5:
			//		valuesList[i][0][1] = fmt.Sprintf("%.2f", max)
			//	case 6:
			//		valuesList[i][0][1] = fmt.Sprintf("%.2f", min)
			//	case 7:
			//		valuesList[i][0][1] = fmt.Sprintf("%.2f", average)
			//	case 8:
			//		valuesList[i][0][1] = fmt.Sprintf("%.2f", median)
			//	}
			//}

			// TODO 不的中回数計算のためには、やはり全組み合わせのオッズがわからないと無理

			// 印組合せのオッズ幅の集計
			for i := position; i < len(defaultValuesList)+position; i++ {
				switch i - position {
				//case 1:
				//	valuesList[i][0][2] = "単全部"
				//	valuesList[i][0][3] = types.TrioOddsRange1.String()
				//	valuesList[i][0][4] = types.TrioOddsRange2.String()
				//	valuesList[i][0][5] = types.TrioOddsRange3.String()
				//	valuesList[i][0][6] = types.TrioOddsRange4.String()
				//	valuesList[i][0][7] = types.TrioOddsRange5.String()
				//	valuesList[i][0][8] = types.TrioOddsRange6.String()
				//	valuesList[i][0][9] = types.TrioOddsRange7.String()
				//	valuesList[i][0][10] = types.TrioOddsRange8.String()
				//case 2:
				//	allWinOddsRangeMap := map[types.OddsRangeType]int{}
				//	for _, oddsRange := range aggregationHitPivotalOddsRangeMap {
				//		allWinOddsRangeMap[types.TrioOddsRange1] += oddsRange[types.TrioOddsRange1]
				//		allWinOddsRangeMap[types.TrioOddsRange2] += oddsRange[types.TrioOddsRange2]
				//		allWinOddsRangeMap[types.TrioOddsRange3] += oddsRange[types.TrioOddsRange3]
				//		allWinOddsRangeMap[types.TrioOddsRange4] += oddsRange[types.TrioOddsRange4]
				//		allWinOddsRangeMap[types.TrioOddsRange5] += oddsRange[types.TrioOddsRange5]
				//		allWinOddsRangeMap[types.TrioOddsRange6] += oddsRange[types.TrioOddsRange6]
				//		allWinOddsRangeMap[types.TrioOddsRange7] += oddsRange[types.TrioOddsRange7]
				//		allWinOddsRangeMap[types.TrioOddsRange8] += oddsRange[types.TrioOddsRange8]
				//	}
				//	valuesList[i][0][3] = allWinOddsRangeMap[types.TrioOddsRange1]
				//	valuesList[i][0][4] = allWinOddsRangeMap[types.TrioOddsRange2]
				//	valuesList[i][0][5] = allWinOddsRangeMap[types.TrioOddsRange3]
				//	valuesList[i][0][6] = allWinOddsRangeMap[types.TrioOddsRange4]
				//	valuesList[i][0][7] = allWinOddsRangeMap[types.TrioOddsRange5]
				//	valuesList[i][0][8] = allWinOddsRangeMap[types.TrioOddsRange6]
				//	valuesList[i][0][9] = allWinOddsRangeMap[types.TrioOddsRange7]
				//	valuesList[i][0][10] = allWinOddsRangeMap[types.TrioOddsRange8]
				//case 3:
				//	allWinOddsRangeMap := map[types.OddsRangeType]int{}
				//	for _, oddsRange := range aggregationHitPivotalOddsRangeMap {
				//		allWinOddsRangeMap[types.TrioOddsRange1] += oddsRange[types.TrioOddsRange1]
				//		allWinOddsRangeMap[types.TrioOddsRange2] += oddsRange[types.TrioOddsRange2]
				//		allWinOddsRangeMap[types.TrioOddsRange3] += oddsRange[types.TrioOddsRange3]
				//		allWinOddsRangeMap[types.TrioOddsRange4] += oddsRange[types.TrioOddsRange4]
				//		allWinOddsRangeMap[types.TrioOddsRange5] += oddsRange[types.TrioOddsRange5]
				//		allWinOddsRangeMap[types.TrioOddsRange6] += oddsRange[types.TrioOddsRange6]
				//		allWinOddsRangeMap[types.TrioOddsRange7] += oddsRange[types.TrioOddsRange7]
				//		allWinOddsRangeMap[types.TrioOddsRange8] += oddsRange[types.TrioOddsRange8]
				//	}
				//	matchCount = 0
				//	for _, count := range allWinOddsRangeMap {
				//		matchCount += count
				//	}
				//	valuesList[i][0][3] = hitRateFormat(allWinOddsRangeMap[types.TrioOddsRange1], matchCount)
				//	valuesList[i][0][4] = hitRateFormat(allWinOddsRangeMap[types.TrioOddsRange2], matchCount)
				//	valuesList[i][0][5] = hitRateFormat(allWinOddsRangeMap[types.TrioOddsRange3], matchCount)
				//	valuesList[i][0][6] = hitRateFormat(allWinOddsRangeMap[types.TrioOddsRange4], matchCount)
				//	valuesList[i][0][7] = hitRateFormat(allWinOddsRangeMap[types.TrioOddsRange5], matchCount)
				//	valuesList[i][0][8] = hitRateFormat(allWinOddsRangeMap[types.TrioOddsRange6], matchCount)
				//	valuesList[i][0][9] = hitRateFormat(allWinOddsRangeMap[types.TrioOddsRange7], matchCount)
				//	valuesList[i][0][10] = hitRateFormat(allWinOddsRangeMap[types.TrioOddsRange8], matchCount)
				//case 4:
				//	allWinOddsRangeMap := map[types.OddsRangeType]int{}
				//	for _, oddsRange := range aggregationHitPivotalOddsRangeMap {
				//		allWinOddsRangeMap[types.TrioOddsRange1] += oddsRange[types.TrioOddsRange1]
				//		allWinOddsRangeMap[types.TrioOddsRange2] += oddsRange[types.TrioOddsRange2]
				//		allWinOddsRangeMap[types.TrioOddsRange3] += oddsRange[types.TrioOddsRange3]
				//		allWinOddsRangeMap[types.TrioOddsRange4] += oddsRange[types.TrioOddsRange4]
				//		allWinOddsRangeMap[types.TrioOddsRange5] += oddsRange[types.TrioOddsRange5]
				//		allWinOddsRangeMap[types.TrioOddsRange6] += oddsRange[types.TrioOddsRange6]
				//		allWinOddsRangeMap[types.TrioOddsRange7] += oddsRange[types.TrioOddsRange7]
				//		allWinOddsRangeMap[types.TrioOddsRange8] += oddsRange[types.TrioOddsRange8]
				//	}
				//	valuesList[i][0][3] = payoutRateFormat(allWinOddsRangeMap[types.TrioOddsRange1], raceCount)
				//	valuesList[i][0][4] = payoutRateFormat(allWinOddsRangeMap[types.TrioOddsRange2], raceCount)
				//	valuesList[i][0][5] = payoutRateFormat(allWinOddsRangeMap[types.TrioOddsRange3], raceCount)
				//	valuesList[i][0][6] = payoutRateFormat(allWinOddsRangeMap[types.TrioOddsRange4], raceCount)
				//	valuesList[i][0][7] = payoutRateFormat(allWinOddsRangeMap[types.TrioOddsRange5], raceCount)
				//	valuesList[i][0][8] = payoutRateFormat(allWinOddsRangeMap[types.TrioOddsRange6], raceCount)
				//	valuesList[i][0][9] = payoutRateFormat(allWinOddsRangeMap[types.TrioOddsRange7], raceCount)
				//	valuesList[i][0][10] = payoutRateFormat(allWinOddsRangeMap[types.TrioOddsRange8], raceCount)
				case 5:
					valuesList[i][0][2] = fmt.Sprintf("単%s", types.WinOddsRange1.String())
					valuesList[i][0][3] = types.TrioOddsRange1.String()
					valuesList[i][0][4] = types.TrioOddsRange2.String()
					valuesList[i][0][5] = types.TrioOddsRange3.String()
					valuesList[i][0][6] = types.TrioOddsRange4.String()
					valuesList[i][0][7] = types.TrioOddsRange5.String()
					valuesList[i][0][8] = types.TrioOddsRange6.String()
					valuesList[i][0][9] = types.TrioOddsRange7.String()
					valuesList[i][0][10] = types.TrioOddsRange8.String()
				case 6:
					oddsMap := pivotalMarkerHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange1]
					valuesList[i][0][3] = oddsMap[types.TrioOddsRange1]
					valuesList[i][0][4] = oddsMap[types.TrioOddsRange2]
					valuesList[i][0][5] = oddsMap[types.TrioOddsRange3]
					valuesList[i][0][6] = oddsMap[types.TrioOddsRange4]
					valuesList[i][0][7] = oddsMap[types.TrioOddsRange5]
					valuesList[i][0][8] = oddsMap[types.TrioOddsRange6]
					valuesList[i][0][9] = oddsMap[types.TrioOddsRange7]
					valuesList[i][0][10] = oddsMap[types.TrioOddsRange8]
				case 7:
					hitCountMap := pivotalMarkerHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange1]
					allCountMap := pivotalMarkerAllOddsRangeCountMap[pivotalMarker][types.WinOddsRange1]
					valuesList[i][0][3] = hitRateFormat(hitCountMap[types.TrioOddsRange1], allCountMap[types.TrioOddsRange1])
					valuesList[i][0][4] = hitRateFormat(hitCountMap[types.TrioOddsRange2], allCountMap[types.TrioOddsRange2])
					valuesList[i][0][5] = hitRateFormat(hitCountMap[types.TrioOddsRange3], allCountMap[types.TrioOddsRange3])
					valuesList[i][0][6] = hitRateFormat(hitCountMap[types.TrioOddsRange4], allCountMap[types.TrioOddsRange4])
					valuesList[i][0][7] = hitRateFormat(hitCountMap[types.TrioOddsRange5], allCountMap[types.TrioOddsRange5])
					valuesList[i][0][8] = hitRateFormat(hitCountMap[types.TrioOddsRange6], allCountMap[types.TrioOddsRange6])
					valuesList[i][0][9] = hitRateFormat(hitCountMap[types.TrioOddsRange7], allCountMap[types.TrioOddsRange7])
					valuesList[i][0][10] = hitRateFormat(hitCountMap[types.TrioOddsRange8], allCountMap[types.TrioOddsRange8])
				case 8:
					hitOddsMap := pivotalMarkerHitTotalOddsMap[pivotalMarker][types.WinOddsRange1]
					allCountMap := pivotalMarkerAllOddsRangeCountMap[pivotalMarker][types.WinOddsRange1]
					valuesList[i][0][3] = payoutRateFormat(hitOddsMap[types.TrioOddsRange1].InexactFloat64(), allCountMap[types.TrioOddsRange1])
					valuesList[i][0][4] = payoutRateFormat(hitOddsMap[types.TrioOddsRange2].InexactFloat64(), allCountMap[types.TrioOddsRange2])
					valuesList[i][0][5] = payoutRateFormat(hitOddsMap[types.TrioOddsRange3].InexactFloat64(), allCountMap[types.TrioOddsRange3])
					valuesList[i][0][6] = payoutRateFormat(hitOddsMap[types.TrioOddsRange4].InexactFloat64(), allCountMap[types.TrioOddsRange4])
					valuesList[i][0][7] = payoutRateFormat(hitOddsMap[types.TrioOddsRange5].InexactFloat64(), allCountMap[types.TrioOddsRange5])
					valuesList[i][0][8] = payoutRateFormat(hitOddsMap[types.TrioOddsRange6].InexactFloat64(), allCountMap[types.TrioOddsRange6])
					valuesList[i][0][9] = payoutRateFormat(hitOddsMap[types.TrioOddsRange7].InexactFloat64(), allCountMap[types.TrioOddsRange7])
					valuesList[i][0][10] = payoutRateFormat(hitOddsMap[types.TrioOddsRange8].InexactFloat64(), allCountMap[types.TrioOddsRange8])
					//case 9:
					//	valuesList[i][0][2] = fmt.Sprintf("単%s", types.WinOddsRange2.String())
					//	valuesList[i][0][3] = types.TrioOddsRange1.String()
					//	valuesList[i][0][4] = types.TrioOddsRange2.String()
					//	valuesList[i][0][5] = types.TrioOddsRange3.String()
					//	valuesList[i][0][6] = types.TrioOddsRange4.String()
					//	valuesList[i][0][7] = types.TrioOddsRange5.String()
					//	valuesList[i][0][8] = types.TrioOddsRange6.String()
					//	valuesList[i][0][9] = types.TrioOddsRange7.String()
					//	valuesList[i][0][10] = types.TrioOddsRange8.String()
					//case 10:
					//	valuesList[i][0][3] = len(totalOddsMap[types.WinOddsRange2][types.TrioOddsRange1])
					//	valuesList[i][0][4] = len(totalOddsMap[types.WinOddsRange2][types.TrioOddsRange2])
					//	valuesList[i][0][5] = len(totalOddsMap[types.WinOddsRange2][types.TrioOddsRange3])
					//	valuesList[i][0][6] = len(totalOddsMap[types.WinOddsRange2][types.TrioOddsRange4])
					//	valuesList[i][0][7] = len(totalOddsMap[types.WinOddsRange2][types.TrioOddsRange5])
					//	valuesList[i][0][8] = len(totalOddsMap[types.WinOddsRange2][types.TrioOddsRange6])
					//	valuesList[i][0][9] = len(totalOddsMap[types.WinOddsRange2][types.TrioOddsRange7])
					//	valuesList[i][0][10] = len(totalOddsMap[types.WinOddsRange2][types.TrioOddsRange8])
					//case 11:
					//	valuesList[i][0][3] = hitRateFormat(len(totalOddsMap[types.WinOddsRange2][types.TrioOddsRange1]), oddsRangeCountMap[types.WinOddsRange2][types.TrioOddsRange1])
					//	valuesList[i][0][4] = hitRateFormat(len(totalOddsMap[types.WinOddsRange2][types.TrioOddsRange2]), oddsRangeCountMap[types.WinOddsRange2][types.TrioOddsRange2])
					//	valuesList[i][0][5] = hitRateFormat(len(totalOddsMap[types.WinOddsRange2][types.TrioOddsRange3]), oddsRangeCountMap[types.WinOddsRange2][types.TrioOddsRange3])
					//	valuesList[i][0][6] = hitRateFormat(len(totalOddsMap[types.WinOddsRange2][types.TrioOddsRange4]), oddsRangeCountMap[types.WinOddsRange2][types.TrioOddsRange4])
					//	valuesList[i][0][7] = hitRateFormat(len(totalOddsMap[types.WinOddsRange2][types.TrioOddsRange5]), oddsRangeCountMap[types.WinOddsRange2][types.TrioOddsRange5])
					//	valuesList[i][0][8] = hitRateFormat(len(totalOddsMap[types.WinOddsRange2][types.TrioOddsRange6]), oddsRangeCountMap[types.WinOddsRange2][types.TrioOddsRange6])
					//	valuesList[i][0][9] = hitRateFormat(len(totalOddsMap[types.WinOddsRange2][types.TrioOddsRange7]), oddsRangeCountMap[types.WinOddsRange2][types.TrioOddsRange7])
					//	valuesList[i][0][10] = hitRateFormat(len(totalOddsMap[types.WinOddsRange2][types.TrioOddsRange8]), oddsRangeCountMap[types.WinOddsRange2][types.TrioOddsRange8])
					//case 12:
					//	valuesList[i][0][3] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange2][types.TrioOddsRange1]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange2][types.TrioOddsRange1])
					//	valuesList[i][0][4] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange2][types.TrioOddsRange2]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange2][types.TrioOddsRange2])
					//	valuesList[i][0][5] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange2][types.TrioOddsRange3]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange2][types.TrioOddsRange3])
					//	valuesList[i][0][6] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange2][types.TrioOddsRange4]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange2][types.TrioOddsRange4])
					//	valuesList[i][0][7] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange2][types.TrioOddsRange5]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange2][types.TrioOddsRange5])
					//	valuesList[i][0][8] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange2][types.TrioOddsRange6]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange2][types.TrioOddsRange6])
					//	valuesList[i][0][9] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange2][types.TrioOddsRange7]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange2][types.TrioOddsRange7])
					//	valuesList[i][0][10] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange2][types.TrioOddsRange8]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange2][types.TrioOddsRange8])
					//case 13:
					//	valuesList[i][0][2] = fmt.Sprintf("単%s", types.WinOddsRange3.String())
					//	valuesList[i][0][3] = types.TrioOddsRange1.String()
					//	valuesList[i][0][4] = types.TrioOddsRange2.String()
					//	valuesList[i][0][5] = types.TrioOddsRange3.String()
					//	valuesList[i][0][6] = types.TrioOddsRange4.String()
					//	valuesList[i][0][7] = types.TrioOddsRange5.String()
					//	valuesList[i][0][8] = types.TrioOddsRange6.String()
					//	valuesList[i][0][9] = types.TrioOddsRange7.String()
					//	valuesList[i][0][10] = types.TrioOddsRange8.String()
					//case 14:
					//	valuesList[i][0][3] = len(totalOddsMap[types.WinOddsRange3][types.TrioOddsRange1])
					//	valuesList[i][0][4] = len(totalOddsMap[types.WinOddsRange3][types.TrioOddsRange2])
					//	valuesList[i][0][5] = len(totalOddsMap[types.WinOddsRange3][types.TrioOddsRange3])
					//	valuesList[i][0][6] = len(totalOddsMap[types.WinOddsRange3][types.TrioOddsRange4])
					//	valuesList[i][0][7] = len(totalOddsMap[types.WinOddsRange3][types.TrioOddsRange5])
					//	valuesList[i][0][8] = len(totalOddsMap[types.WinOddsRange3][types.TrioOddsRange6])
					//	valuesList[i][0][9] = len(totalOddsMap[types.WinOddsRange3][types.TrioOddsRange7])
					//	valuesList[i][0][10] = len(totalOddsMap[types.WinOddsRange3][types.TrioOddsRange8])
					//case 15:
					//	valuesList[i][0][3] = hitRateFormat(len(totalOddsMap[types.WinOddsRange3][types.TrioOddsRange1]), oddsRangeCountMap[types.WinOddsRange3][types.TrioOddsRange1])
					//	valuesList[i][0][4] = hitRateFormat(len(totalOddsMap[types.WinOddsRange3][types.TrioOddsRange2]), oddsRangeCountMap[types.WinOddsRange3][types.TrioOddsRange2])
					//	valuesList[i][0][5] = hitRateFormat(len(totalOddsMap[types.WinOddsRange3][types.TrioOddsRange3]), oddsRangeCountMap[types.WinOddsRange3][types.TrioOddsRange3])
					//	valuesList[i][0][6] = hitRateFormat(len(totalOddsMap[types.WinOddsRange3][types.TrioOddsRange4]), oddsRangeCountMap[types.WinOddsRange3][types.TrioOddsRange4])
					//	valuesList[i][0][7] = hitRateFormat(len(totalOddsMap[types.WinOddsRange3][types.TrioOddsRange5]), oddsRangeCountMap[types.WinOddsRange3][types.TrioOddsRange5])
					//	valuesList[i][0][8] = hitRateFormat(len(totalOddsMap[types.WinOddsRange3][types.TrioOddsRange6]), oddsRangeCountMap[types.WinOddsRange3][types.TrioOddsRange6])
					//	valuesList[i][0][9] = hitRateFormat(len(totalOddsMap[types.WinOddsRange3][types.TrioOddsRange7]), oddsRangeCountMap[types.WinOddsRange3][types.TrioOddsRange7])
					//	valuesList[i][0][10] = hitRateFormat(len(totalOddsMap[types.WinOddsRange3][types.TrioOddsRange8]), oddsRangeCountMap[types.WinOddsRange3][types.TrioOddsRange8])
					//case 16:
					//	valuesList[i][0][3] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange3][types.TrioOddsRange1]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange3][types.TrioOddsRange1])
					//	valuesList[i][0][4] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange3][types.TrioOddsRange2]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange3][types.TrioOddsRange2])
					//	valuesList[i][0][5] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange3][types.TrioOddsRange3]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange3][types.TrioOddsRange3])
					//	valuesList[i][0][6] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange3][types.TrioOddsRange4]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange3][types.TrioOddsRange4])
					//	valuesList[i][0][7] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange3][types.TrioOddsRange5]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange3][types.TrioOddsRange5])
					//	valuesList[i][0][8] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange3][types.TrioOddsRange6]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange3][types.TrioOddsRange6])
					//	valuesList[i][0][9] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange3][types.TrioOddsRange7]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange3][types.TrioOddsRange7])
					//	valuesList[i][0][10] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange3][types.TrioOddsRange8]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange3][types.TrioOddsRange8])
					//case 17:
					//	valuesList[i][0][2] = fmt.Sprintf("単%s", types.WinOddsRange4.String())
					//	valuesList[i][0][3] = types.TrioOddsRange1.String()
					//	valuesList[i][0][4] = types.TrioOddsRange2.String()
					//	valuesList[i][0][5] = types.TrioOddsRange3.String()
					//	valuesList[i][0][6] = types.TrioOddsRange4.String()
					//	valuesList[i][0][7] = types.TrioOddsRange5.String()
					//	valuesList[i][0][8] = types.TrioOddsRange6.String()
					//	valuesList[i][0][9] = types.TrioOddsRange7.String()
					//	valuesList[i][0][10] = types.TrioOddsRange8.String()
					//case 18:
					//	valuesList[i][0][3] = len(totalOddsMap[types.WinOddsRange4][types.TrioOddsRange1])
					//	valuesList[i][0][4] = len(totalOddsMap[types.WinOddsRange4][types.TrioOddsRange2])
					//	valuesList[i][0][5] = len(totalOddsMap[types.WinOddsRange4][types.TrioOddsRange3])
					//	valuesList[i][0][6] = len(totalOddsMap[types.WinOddsRange4][types.TrioOddsRange4])
					//	valuesList[i][0][7] = len(totalOddsMap[types.WinOddsRange4][types.TrioOddsRange5])
					//	valuesList[i][0][8] = len(totalOddsMap[types.WinOddsRange4][types.TrioOddsRange6])
					//	valuesList[i][0][9] = len(totalOddsMap[types.WinOddsRange4][types.TrioOddsRange7])
					//	valuesList[i][0][10] = len(totalOddsMap[types.WinOddsRange4][types.TrioOddsRange8])
					//case 19:
					//	valuesList[i][0][3] = hitRateFormat(len(totalOddsMap[types.WinOddsRange4][types.TrioOddsRange1]), oddsRangeCountMap[types.WinOddsRange4][types.TrioOddsRange1])
					//	valuesList[i][0][4] = hitRateFormat(len(totalOddsMap[types.WinOddsRange4][types.TrioOddsRange2]), oddsRangeCountMap[types.WinOddsRange4][types.TrioOddsRange2])
					//	valuesList[i][0][5] = hitRateFormat(len(totalOddsMap[types.WinOddsRange4][types.TrioOddsRange3]), oddsRangeCountMap[types.WinOddsRange4][types.TrioOddsRange3])
					//	valuesList[i][0][6] = hitRateFormat(len(totalOddsMap[types.WinOddsRange4][types.TrioOddsRange4]), oddsRangeCountMap[types.WinOddsRange4][types.TrioOddsRange4])
					//	valuesList[i][0][7] = hitRateFormat(len(totalOddsMap[types.WinOddsRange4][types.TrioOddsRange5]), oddsRangeCountMap[types.WinOddsRange4][types.TrioOddsRange5])
					//	valuesList[i][0][8] = hitRateFormat(len(totalOddsMap[types.WinOddsRange4][types.TrioOddsRange6]), oddsRangeCountMap[types.WinOddsRange4][types.TrioOddsRange6])
					//	valuesList[i][0][9] = hitRateFormat(len(totalOddsMap[types.WinOddsRange4][types.TrioOddsRange7]), oddsRangeCountMap[types.WinOddsRange4][types.TrioOddsRange7])
					//	valuesList[i][0][10] = hitRateFormat(len(totalOddsMap[types.WinOddsRange4][types.TrioOddsRange8]), oddsRangeCountMap[types.WinOddsRange4][types.TrioOddsRange8])
					//case 20:
					//	valuesList[i][0][3] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange4][types.TrioOddsRange1]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange4][types.TrioOddsRange1])
					//	valuesList[i][0][4] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange4][types.TrioOddsRange2]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange4][types.TrioOddsRange2])
					//	valuesList[i][0][5] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange4][types.TrioOddsRange3]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange4][types.TrioOddsRange3])
					//	valuesList[i][0][6] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange4][types.TrioOddsRange4]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange4][types.TrioOddsRange4])
					//	valuesList[i][0][7] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange4][types.TrioOddsRange5]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange4][types.TrioOddsRange5])
					//	valuesList[i][0][8] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange4][types.TrioOddsRange6]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange4][types.TrioOddsRange6])
					//	valuesList[i][0][9] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange4][types.TrioOddsRange7]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange4][types.TrioOddsRange7])
					//	valuesList[i][0][10] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange4][types.TrioOddsRange8]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange4][types.TrioOddsRange8])
					//case 21:
					//	valuesList[i][0][2] = fmt.Sprintf("単%s", types.WinOddsRange5.String())
					//	valuesList[i][0][3] = types.TrioOddsRange1.String()
					//	valuesList[i][0][4] = types.TrioOddsRange2.String()
					//	valuesList[i][0][5] = types.TrioOddsRange3.String()
					//	valuesList[i][0][6] = types.TrioOddsRange4.String()
					//	valuesList[i][0][7] = types.TrioOddsRange5.String()
					//	valuesList[i][0][8] = types.TrioOddsRange6.String()
					//	valuesList[i][0][9] = types.TrioOddsRange7.String()
					//	valuesList[i][0][10] = types.TrioOddsRange8.String()
					//case 22:
					//	valuesList[i][0][3] = len(totalOddsMap[types.WinOddsRange5][types.TrioOddsRange1])
					//	valuesList[i][0][4] = len(totalOddsMap[types.WinOddsRange5][types.TrioOddsRange2])
					//	valuesList[i][0][5] = len(totalOddsMap[types.WinOddsRange5][types.TrioOddsRange3])
					//	valuesList[i][0][6] = len(totalOddsMap[types.WinOddsRange5][types.TrioOddsRange4])
					//	valuesList[i][0][7] = len(totalOddsMap[types.WinOddsRange5][types.TrioOddsRange5])
					//	valuesList[i][0][8] = len(totalOddsMap[types.WinOddsRange5][types.TrioOddsRange6])
					//	valuesList[i][0][9] = len(totalOddsMap[types.WinOddsRange5][types.TrioOddsRange7])
					//	valuesList[i][0][10] = len(totalOddsMap[types.WinOddsRange5][types.TrioOddsRange8])
					//case 23:
					//	valuesList[i][0][3] = hitRateFormat(len(totalOddsMap[types.WinOddsRange5][types.TrioOddsRange1]), oddsRangeCountMap[types.WinOddsRange5][types.TrioOddsRange1])
					//	valuesList[i][0][4] = hitRateFormat(len(totalOddsMap[types.WinOddsRange5][types.TrioOddsRange2]), oddsRangeCountMap[types.WinOddsRange5][types.TrioOddsRange2])
					//	valuesList[i][0][5] = hitRateFormat(len(totalOddsMap[types.WinOddsRange5][types.TrioOddsRange3]), oddsRangeCountMap[types.WinOddsRange5][types.TrioOddsRange3])
					//	valuesList[i][0][6] = hitRateFormat(len(totalOddsMap[types.WinOddsRange5][types.TrioOddsRange4]), oddsRangeCountMap[types.WinOddsRange5][types.TrioOddsRange4])
					//	valuesList[i][0][7] = hitRateFormat(len(totalOddsMap[types.WinOddsRange5][types.TrioOddsRange5]), oddsRangeCountMap[types.WinOddsRange5][types.TrioOddsRange5])
					//	valuesList[i][0][8] = hitRateFormat(len(totalOddsMap[types.WinOddsRange5][types.TrioOddsRange6]), oddsRangeCountMap[types.WinOddsRange5][types.TrioOddsRange6])
					//	valuesList[i][0][9] = hitRateFormat(len(totalOddsMap[types.WinOddsRange5][types.TrioOddsRange7]), oddsRangeCountMap[types.WinOddsRange5][types.TrioOddsRange7])
					//	valuesList[i][0][10] = hitRateFormat(len(totalOddsMap[types.WinOddsRange5][types.TrioOddsRange8]), oddsRangeCountMap[types.WinOddsRange5][types.TrioOddsRange8])
					//case 24:
					//	valuesList[i][0][3] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange5][types.TrioOddsRange1]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange5][types.TrioOddsRange1])
					//	valuesList[i][0][4] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange5][types.TrioOddsRange2]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange5][types.TrioOddsRange2])
					//	valuesList[i][0][5] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange5][types.TrioOddsRange3]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange5][types.TrioOddsRange3])
					//	valuesList[i][0][6] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange5][types.TrioOddsRange4]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange5][types.TrioOddsRange4])
					//	valuesList[i][0][7] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange5][types.TrioOddsRange5]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange5][types.TrioOddsRange5])
					//	valuesList[i][0][8] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange5][types.TrioOddsRange6]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange5][types.TrioOddsRange6])
					//	valuesList[i][0][9] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange5][types.TrioOddsRange7]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange5][types.TrioOddsRange7])
					//	valuesList[i][0][10] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange5][types.TrioOddsRange8]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange5][types.TrioOddsRange8])
					//case 25:
					//	valuesList[i][0][2] = fmt.Sprintf("単%s", types.WinOddsRange6.String())
					//	valuesList[i][0][3] = types.TrioOddsRange1.String()
					//	valuesList[i][0][4] = types.TrioOddsRange2.String()
					//	valuesList[i][0][5] = types.TrioOddsRange3.String()
					//	valuesList[i][0][6] = types.TrioOddsRange4.String()
					//	valuesList[i][0][7] = types.TrioOddsRange5.String()
					//	valuesList[i][0][8] = types.TrioOddsRange6.String()
					//	valuesList[i][0][9] = types.TrioOddsRange7.String()
					//	valuesList[i][0][10] = types.TrioOddsRange8.String()
					//case 26:
					//	valuesList[i][0][3] = len(totalOddsMap[types.WinOddsRange6][types.TrioOddsRange1])
					//	valuesList[i][0][4] = len(totalOddsMap[types.WinOddsRange6][types.TrioOddsRange2])
					//	valuesList[i][0][5] = len(totalOddsMap[types.WinOddsRange6][types.TrioOddsRange3])
					//	valuesList[i][0][6] = len(totalOddsMap[types.WinOddsRange6][types.TrioOddsRange4])
					//	valuesList[i][0][7] = len(totalOddsMap[types.WinOddsRange6][types.TrioOddsRange5])
					//	valuesList[i][0][8] = len(totalOddsMap[types.WinOddsRange6][types.TrioOddsRange6])
					//	valuesList[i][0][9] = len(totalOddsMap[types.WinOddsRange6][types.TrioOddsRange7])
					//	valuesList[i][0][10] = len(totalOddsMap[types.WinOddsRange6][types.TrioOddsRange8])
					//case 27:
					//	valuesList[i][0][3] = hitRateFormat(len(totalOddsMap[types.WinOddsRange6][types.TrioOddsRange1]), oddsRangeCountMap[types.WinOddsRange6][types.TrioOddsRange1])
					//	valuesList[i][0][4] = hitRateFormat(len(totalOddsMap[types.WinOddsRange6][types.TrioOddsRange2]), oddsRangeCountMap[types.WinOddsRange6][types.TrioOddsRange2])
					//	valuesList[i][0][5] = hitRateFormat(len(totalOddsMap[types.WinOddsRange6][types.TrioOddsRange3]), oddsRangeCountMap[types.WinOddsRange6][types.TrioOddsRange3])
					//	valuesList[i][0][6] = hitRateFormat(len(totalOddsMap[types.WinOddsRange6][types.TrioOddsRange4]), oddsRangeCountMap[types.WinOddsRange6][types.TrioOddsRange4])
					//	valuesList[i][0][7] = hitRateFormat(len(totalOddsMap[types.WinOddsRange6][types.TrioOddsRange5]), oddsRangeCountMap[types.WinOddsRange6][types.TrioOddsRange5])
					//	valuesList[i][0][8] = hitRateFormat(len(totalOddsMap[types.WinOddsRange6][types.TrioOddsRange6]), oddsRangeCountMap[types.WinOddsRange6][types.TrioOddsRange6])
					//	valuesList[i][0][9] = hitRateFormat(len(totalOddsMap[types.WinOddsRange6][types.TrioOddsRange7]), oddsRangeCountMap[types.WinOddsRange6][types.TrioOddsRange7])
					//	valuesList[i][0][10] = hitRateFormat(len(totalOddsMap[types.WinOddsRange6][types.TrioOddsRange8]), oddsRangeCountMap[types.WinOddsRange6][types.TrioOddsRange8])
					//case 28:
					//	valuesList[i][0][3] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange6][types.TrioOddsRange1]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange6][types.TrioOddsRange1])
					//	valuesList[i][0][4] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange6][types.TrioOddsRange2]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange6][types.TrioOddsRange2])
					//	valuesList[i][0][5] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange6][types.TrioOddsRange3]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange6][types.TrioOddsRange3])
					//	valuesList[i][0][6] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange6][types.TrioOddsRange4]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange6][types.TrioOddsRange4])
					//	valuesList[i][0][7] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange6][types.TrioOddsRange5]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange6][types.TrioOddsRange5])
					//	valuesList[i][0][8] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange6][types.TrioOddsRange6]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange6][types.TrioOddsRange6])
					//	valuesList[i][0][9] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange6][types.TrioOddsRange7]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange6][types.TrioOddsRange7])
					//	valuesList[i][0][10] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange6][types.TrioOddsRange8]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange6][types.TrioOddsRange8])
					//case 29:
					//	valuesList[i][0][2] = fmt.Sprintf("単%s", types.WinOddsRange7.String())
					//	valuesList[i][0][3] = types.TrioOddsRange1.String()
					//	valuesList[i][0][4] = types.TrioOddsRange2.String()
					//	valuesList[i][0][5] = types.TrioOddsRange3.String()
					//	valuesList[i][0][6] = types.TrioOddsRange4.String()
					//	valuesList[i][0][7] = types.TrioOddsRange5.String()
					//	valuesList[i][0][8] = types.TrioOddsRange6.String()
					//	valuesList[i][0][9] = types.TrioOddsRange7.String()
					//	valuesList[i][0][10] = types.TrioOddsRange8.String()
					//case 30:
					//	valuesList[i][0][3] = len(totalOddsMap[types.WinOddsRange7][types.TrioOddsRange1])
					//	valuesList[i][0][4] = len(totalOddsMap[types.WinOddsRange7][types.TrioOddsRange2])
					//	valuesList[i][0][5] = len(totalOddsMap[types.WinOddsRange7][types.TrioOddsRange3])
					//	valuesList[i][0][6] = len(totalOddsMap[types.WinOddsRange7][types.TrioOddsRange4])
					//	valuesList[i][0][7] = len(totalOddsMap[types.WinOddsRange7][types.TrioOddsRange5])
					//	valuesList[i][0][8] = len(totalOddsMap[types.WinOddsRange7][types.TrioOddsRange6])
					//	valuesList[i][0][9] = len(totalOddsMap[types.WinOddsRange7][types.TrioOddsRange7])
					//	valuesList[i][0][10] = len(totalOddsMap[types.WinOddsRange7][types.TrioOddsRange8])
					//case 31:
					//	valuesList[i][0][3] = hitRateFormat(len(totalOddsMap[types.WinOddsRange7][types.TrioOddsRange1]), oddsRangeCountMap[types.WinOddsRange7][types.TrioOddsRange1])
					//	valuesList[i][0][4] = hitRateFormat(len(totalOddsMap[types.WinOddsRange7][types.TrioOddsRange2]), oddsRangeCountMap[types.WinOddsRange7][types.TrioOddsRange2])
					//	valuesList[i][0][5] = hitRateFormat(len(totalOddsMap[types.WinOddsRange7][types.TrioOddsRange3]), oddsRangeCountMap[types.WinOddsRange7][types.TrioOddsRange3])
					//	valuesList[i][0][6] = hitRateFormat(len(totalOddsMap[types.WinOddsRange7][types.TrioOddsRange4]), oddsRangeCountMap[types.WinOddsRange7][types.TrioOddsRange4])
					//	valuesList[i][0][7] = hitRateFormat(len(totalOddsMap[types.WinOddsRange7][types.TrioOddsRange5]), oddsRangeCountMap[types.WinOddsRange7][types.TrioOddsRange5])
					//	valuesList[i][0][8] = hitRateFormat(len(totalOddsMap[types.WinOddsRange7][types.TrioOddsRange6]), oddsRangeCountMap[types.WinOddsRange7][types.TrioOddsRange6])
					//	valuesList[i][0][9] = hitRateFormat(len(totalOddsMap[types.WinOddsRange7][types.TrioOddsRange7]), oddsRangeCountMap[types.WinOddsRange7][types.TrioOddsRange7])
					//	valuesList[i][0][10] = hitRateFormat(len(totalOddsMap[types.WinOddsRange7][types.TrioOddsRange8]), oddsRangeCountMap[types.WinOddsRange7][types.TrioOddsRange8])
					//case 32:
					//	valuesList[i][0][3] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange7][types.TrioOddsRange1]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange7][types.TrioOddsRange1])
					//	valuesList[i][0][4] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange7][types.TrioOddsRange2]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange7][types.TrioOddsRange2])
					//	valuesList[i][0][5] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange7][types.TrioOddsRange3]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange7][types.TrioOddsRange3])
					//	valuesList[i][0][6] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange7][types.TrioOddsRange4]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange7][types.TrioOddsRange4])
					//	valuesList[i][0][7] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange7][types.TrioOddsRange5]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange7][types.TrioOddsRange5])
					//	valuesList[i][0][8] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange7][types.TrioOddsRange6]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange7][types.TrioOddsRange6])
					//	valuesList[i][0][9] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange7][types.TrioOddsRange7]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange7][types.TrioOddsRange7])
					//	valuesList[i][0][10] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange7][types.TrioOddsRange8]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange7][types.TrioOddsRange8])
					//case 33:
					//	valuesList[i][0][2] = fmt.Sprintf("単%s", types.WinOddsRange8.String())
					//	valuesList[i][0][3] = types.TrioOddsRange1.String()
					//	valuesList[i][0][4] = types.TrioOddsRange2.String()
					//	valuesList[i][0][5] = types.TrioOddsRange3.String()
					//	valuesList[i][0][6] = types.TrioOddsRange4.String()
					//	valuesList[i][0][7] = types.TrioOddsRange5.String()
					//	valuesList[i][0][8] = types.TrioOddsRange6.String()
					//	valuesList[i][0][9] = types.TrioOddsRange7.String()
					//	valuesList[i][0][10] = types.TrioOddsRange8.String()
					//case 34:
					//	valuesList[i][0][3] = len(totalOddsMap[types.WinOddsRange8][types.TrioOddsRange1])
					//	valuesList[i][0][4] = len(totalOddsMap[types.WinOddsRange8][types.TrioOddsRange2])
					//	valuesList[i][0][5] = len(totalOddsMap[types.WinOddsRange8][types.TrioOddsRange3])
					//	valuesList[i][0][6] = len(totalOddsMap[types.WinOddsRange8][types.TrioOddsRange4])
					//	valuesList[i][0][7] = len(totalOddsMap[types.WinOddsRange8][types.TrioOddsRange5])
					//	valuesList[i][0][8] = len(totalOddsMap[types.WinOddsRange8][types.TrioOddsRange6])
					//	valuesList[i][0][9] = len(totalOddsMap[types.WinOddsRange8][types.TrioOddsRange7])
					//	valuesList[i][0][10] = len(totalOddsMap[types.WinOddsRange8][types.TrioOddsRange8])
					//case 35:
					//	valuesList[i][0][3] = hitRateFormat(len(totalOddsMap[types.WinOddsRange8][types.TrioOddsRange1]), oddsRangeCountMap[types.WinOddsRange8][types.TrioOddsRange1])
					//	valuesList[i][0][4] = hitRateFormat(len(totalOddsMap[types.WinOddsRange8][types.TrioOddsRange2]), oddsRangeCountMap[types.WinOddsRange8][types.TrioOddsRange2])
					//	valuesList[i][0][5] = hitRateFormat(len(totalOddsMap[types.WinOddsRange8][types.TrioOddsRange3]), oddsRangeCountMap[types.WinOddsRange8][types.TrioOddsRange3])
					//	valuesList[i][0][6] = hitRateFormat(len(totalOddsMap[types.WinOddsRange8][types.TrioOddsRange4]), oddsRangeCountMap[types.WinOddsRange8][types.TrioOddsRange4])
					//	valuesList[i][0][7] = hitRateFormat(len(totalOddsMap[types.WinOddsRange8][types.TrioOddsRange5]), oddsRangeCountMap[types.WinOddsRange8][types.TrioOddsRange5])
					//	valuesList[i][0][8] = hitRateFormat(len(totalOddsMap[types.WinOddsRange8][types.TrioOddsRange6]), oddsRangeCountMap[types.WinOddsRange8][types.TrioOddsRange6])
					//	valuesList[i][0][9] = hitRateFormat(len(totalOddsMap[types.WinOddsRange8][types.TrioOddsRange7]), oddsRangeCountMap[types.WinOddsRange8][types.TrioOddsRange7])
					//	valuesList[i][0][10] = hitRateFormat(len(totalOddsMap[types.WinOddsRange8][types.TrioOddsRange8]), oddsRangeCountMap[types.WinOddsRange8][types.TrioOddsRange8])
					//case 36:
					//	valuesList[i][0][3] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange8][types.TrioOddsRange1]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange8][types.TrioOddsRange1])
					//	valuesList[i][0][4] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange8][types.TrioOddsRange2]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange8][types.TrioOddsRange2])
					//	valuesList[i][0][5] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange8][types.TrioOddsRange3]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange8][types.TrioOddsRange3])
					//	valuesList[i][0][6] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange8][types.TrioOddsRange4]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange8][types.TrioOddsRange4])
					//	valuesList[i][0][7] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange8][types.TrioOddsRange5]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange8][types.TrioOddsRange5])
					//	valuesList[i][0][8] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange8][types.TrioOddsRange6]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange8][types.TrioOddsRange6])
					//	valuesList[i][0][9] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange8][types.TrioOddsRange7]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange8][types.TrioOddsRange7])
					//	valuesList[i][0][10] = payoutRateFormat(s.getTotalOdds(ctx, totalOddsMap[types.WinOddsRange8][types.TrioOddsRange8]).InexactFloat64(), oddsRangeCountMap[types.WinOddsRange8][types.TrioOddsRange8])
					//}
				}
			}

			aggregationMarkerIndex++
		}
	}

	var values [][]interface{}
	for _, v := range valuesList {
		values = append(values, v...)
	}
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName(), fmt.Sprintf("A1"))
	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetTrioAnalysisRepository) createDefaultValuesList() [][][]interface{} {
	valuesList := make([][][]interface{}, 0)
	valuesList = append(valuesList, [][]interface{}{
		{
			"",
			"",
			"軸の単勝オッズに対する三連複のオッズ幅の集計",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	// OddsRange単全部
	valuesList = append(valuesList, [][]interface{}{
		{
			"レース数",
			"",
			"単全部",
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
	valuesList = append(valuesList, [][]interface{}{
		{
			"的中回数",
			"",
			"的中回数",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			"的中率",
			"",
			"軸選択的中率",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			"回収率",
			"",
			"軸選択的中率",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	// OddsRange単1
	valuesList = append(valuesList, [][]interface{}{
		{
			"最大オッズ",
			"",
			fmt.Sprintf("単%s", types.WinOddsRange1.String()),
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
	valuesList = append(valuesList, [][]interface{}{
		{
			"最小オッズ",
			"",
			"的中回数",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			"平均オッズ",
			"",
			"軸選択的中率",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			"中央オッズ",
			"",
			"軸選択回収率",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	// OddsRange単2
	valuesList = append(valuesList, [][]interface{}{
		{
			fmt.Sprintf("%s率", types.TrioOddsRange1.String()),
			"",
			fmt.Sprintf("単%s", types.WinOddsRange2.String()),
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
	valuesList = append(valuesList, [][]interface{}{
		{
			fmt.Sprintf("%s率", types.TrioOddsRange2.String()),
			"",
			"的中回数",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			fmt.Sprintf("%s率", types.TrioOddsRange3.String()),
			"",
			"軸選択的中率",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			fmt.Sprintf("%s率", types.TrioOddsRange4.String()),
			"",
			"軸選択回収率",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	// OddsRange単3
	valuesList = append(valuesList, [][]interface{}{
		{
			fmt.Sprintf("%s率", types.TrioOddsRange5.String()),
			"",
			fmt.Sprintf("単%s", types.WinOddsRange3.String()),
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
	valuesList = append(valuesList, [][]interface{}{
		{
			fmt.Sprintf("%s率", types.TrioOddsRange6.String()),
			"",
			"的中回数",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			fmt.Sprintf("%s率", types.TrioOddsRange7.String()),
			"",
			"軸選択的中率",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			fmt.Sprintf("%s率", types.TrioOddsRange8.String()),
			"",
			"軸選択回収率",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	// OddsRange単4
	valuesList = append(valuesList, [][]interface{}{
		{
			fmt.Sprintf("%s率", types.WinOddsRange1.String()),
			"",
			fmt.Sprintf("単%s", types.WinOddsRange4.String()),
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
	valuesList = append(valuesList, [][]interface{}{
		{
			fmt.Sprintf("%s率", types.WinOddsRange2.String()),
			"",
			"的中回数",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			fmt.Sprintf("%s率", types.WinOddsRange3.String()),
			"",
			"軸選択的中率",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			fmt.Sprintf("%s率", types.WinOddsRange4.String()),
			"",
			"軸選択回収率",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	// OddsRange単5
	valuesList = append(valuesList, [][]interface{}{
		{
			fmt.Sprintf("%s率", types.WinOddsRange5.String()),
			"",
			fmt.Sprintf("単%s", types.WinOddsRange5.String()),
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
	valuesList = append(valuesList, [][]interface{}{
		{
			fmt.Sprintf("%s率", types.WinOddsRange6.String()),
			"",
			"的中回数",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			fmt.Sprintf("%s率", types.WinOddsRange7.String()),
			"",
			"軸選択的中率",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			fmt.Sprintf("%s率", types.WinOddsRange8.String()),
			"",
			"軸選択回収率",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	// OddsRange単6
	valuesList = append(valuesList, [][]interface{}{
		{
			"",
			"",
			fmt.Sprintf("単%s", types.WinOddsRange6.String()),
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
	valuesList = append(valuesList, [][]interface{}{
		{
			"",
			"",
			"的中回数",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			"",
			"",
			"軸選択的中率",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			"",
			"",
			"軸選択回収率",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	// OddsRange単7
	valuesList = append(valuesList, [][]interface{}{
		{
			"",
			"",
			fmt.Sprintf("単%s", types.WinOddsRange7.String()),
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
	valuesList = append(valuesList, [][]interface{}{
		{
			"",
			"",
			"的中回数",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			"",
			"",
			"軸選択的中率",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			"",
			"",
			"軸選択回収率",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	// OddsRange単8
	valuesList = append(valuesList, [][]interface{}{
		{
			"",
			"",
			fmt.Sprintf("単%s", types.WinOddsRange8.String()),
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
	valuesList = append(valuesList, [][]interface{}{
		{
			"",
			"",
			"的中回数",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			"",
			"",
			"軸選択的中率",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			"",
			"",
			"軸選択回収率",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})

	return valuesList
}

//func (s *spreadSheetTrioAnalysisRepository) createHitTrioOddsRangeMap(
//	ctx context.Context,
//
//) {
//
//}

//func (s *spreadSheetTrioAnalysisRepository) createHitTrioOddsRangeMap(
//	ctx context.Context,
//	markerCombinationAnalysisList []*spreadsheet_entity.MarkerCombinationAnalysis,
//	pivotalMarker types.Marker,
//) map[types.OddsRangeType]map[types.OddsRangeType]int {
//	pivotalOddsRangeMap := map[types.OddsRangeType]map[types.OddsRangeType]int{}
//
//	for _, markerCombinationAnalysis := range markerCombinationAnalysisList {
//		for _, calculable := range markerCombinationAnalysis.Calculables() {
//			var (
//				isContainPivotal bool
//				pivotalOddsRange types.OddsRangeType
//			)
//			for _, pivotal := range calculable.Pivotals() {
//				if pivotalMarker == pivotal.Marker() {
//					isContainPivotal = true
//					odds := pivotal.Odds().InexactFloat64()
//					if odds >= 1.0 && odds <= 1.5 {
//						pivotalOddsRange = types.WinOddsRange1
//					} else if odds >= 1.6 && odds <= 2.0 {
//						pivotalOddsRange = types.WinOddsRange2
//					} else if odds >= 2.1 && odds <= 2.9 {
//						pivotalOddsRange = types.WinOddsRange3
//					} else if odds >= 3.0 && odds <= 4.9 {
//						pivotalOddsRange = types.WinOddsRange4
//					} else if odds >= 5.0 && odds <= 9.9 {
//						pivotalOddsRange = types.WinOddsRange5
//					} else if odds >= 10.0 && odds <= 19.9 {
//						pivotalOddsRange = types.WinOddsRange6
//					} else if odds >= 20.0 && odds <= 49.9 {
//						pivotalOddsRange = types.WinOddsRange7
//					} else if odds >= 50.0 {
//						pivotalOddsRange = types.WinOddsRange8
//					}
//				}
//			}
//
//			if !isContainPivotal {
//				continue
//			}
//
//			if _, ok := pivotalOddsRangeMap[pivotalOddsRange]; !ok {
//				pivotalOddsRangeMap[pivotalOddsRange] = map[types.OddsRangeType]int{}
//			}
//
//			odds := calculable.Odds().InexactFloat64()
//			if odds >= 1.0 && odds <= 9.9 {
//				if _, ok := pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange1]; !ok {
//					pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange1] = 0
//				}
//				pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange1]++
//			} else if odds >= 10.0 && odds <= 19.9 {
//				if _, ok := pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange2]; !ok {
//					pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange2] = 0
//				}
//				pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange2]++
//			} else if odds >= 20.0 && odds <= 29.9 {
//				if _, ok := pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange3]; !ok {
//					pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange3] = 0
//				}
//				pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange3]++
//			} else if odds >= 30.0 && odds <= 49.9 {
//				if _, ok := pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange4]; !ok {
//					pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange4] = 0
//				}
//				pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange4]++
//			} else if odds >= 50.0 && odds <= 99.9 {
//				if _, ok := pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange5]; !ok {
//					pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange5] = 0
//				}
//				pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange5]++
//			} else if odds >= 100.0 && odds <= 299.9 {
//				if _, ok := pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange6]; !ok {
//					pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange6] = 0
//				}
//				pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange6]++
//			} else if odds >= 300.0 && odds <= 499.9 {
//				if _, ok := pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange7]; !ok {
//					pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange7] = 0
//				}
//				pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange7]++
//			} else if odds >= 500.0 {
//				if _, ok := pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange8]; !ok {
//					pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange8] = 0
//				}
//				pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange8]++
//			}
//		}
//	}
//
//	return pivotalOddsRangeMap
//}

func (s *spreadSheetTrioAnalysisRepository) getPivotalMarkerHitOddsRangeCountMap(
	ctx context.Context,
	markerCombinationAnalysisMap map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis,
	winRaceOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
	trioRaceOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
) (map[types.Marker]map[types.OddsRangeType]map[types.OddsRangeType]int, error) {
	pivotalMarkerHitCountOddsRangeMap := map[types.Marker]map[types.OddsRangeType]map[types.OddsRangeType]int{}
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
			for _, calculable := range markerCombinationAnalysis.Calculables() {
				if calculable.IsHit() {
					hitMarkerCalculablesMap[marker] = append(hitMarkerCalculablesMap[marker], calculable)
				}
			}
		}
	}

	for marker, calculables := range hitMarkerCalculablesMap {
		pivotalMarkerHitCountOddsRangeMap[marker] = map[types.OddsRangeType]map[types.OddsRangeType]int{}
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

				if _, ok := pivotalMarkerHitCountOddsRangeMap[marker][pivotalMarkerOddsRange]; !ok {
					pivotalMarkerHitCountOddsRangeMap[marker][pivotalMarkerOddsRange] = map[types.OddsRangeType]int{}
				}

				for _, trioOdds := range trioOddsList {
					odds = trioOdds.Odds().InexactFloat64()
					if odds >= 1.0 && odds <= 9.9 {
						pivotalMarkerHitCountOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange1]++
					} else if odds >= 10.0 && odds <= 19.9 {
						pivotalMarkerHitCountOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange2]++
					} else if odds >= 20.0 && odds <= 29.9 {
						pivotalMarkerHitCountOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange3]++
					} else if odds >= 30.0 && odds <= 49.9 {
						pivotalMarkerHitCountOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange4]++
					} else if odds >= 50.0 && odds <= 99.9 {
						pivotalMarkerHitCountOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange5]++
					} else if odds >= 100.0 && odds <= 299.9 {
						pivotalMarkerHitCountOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange6]++
					} else if odds >= 300.0 && odds <= 499.9 {
						pivotalMarkerHitCountOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange7]++
					} else if odds >= 500.0 {
						pivotalMarkerHitCountOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange8]++
					}
				}
			}
		}
	}

	return pivotalMarkerHitCountOddsRangeMap, nil
}

func (s *spreadSheetTrioAnalysisRepository) getPivotalMarkerAllOddsRangeCountMap(
	ctx context.Context,
	markerCombinationAnalysisMap map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis,
	winRaceOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
	trioRaceOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
) (map[types.Marker]map[types.OddsRangeType]map[types.OddsRangeType]int, error) {
	pivotalMarkerAllCountOddsRangeMap := map[types.Marker]map[types.OddsRangeType]map[types.OddsRangeType]int{}
	pivotalMarkers := []types.Marker{
		types.Favorite,
		types.Rival,
		types.BrackTriangle,
		types.WhiteTriangle,
		types.Star,
		types.Check,
	}

	allMarkerCalculablesMap := map[types.Marker][]*analysis_entity.Calculable{}
	for _, marker := range pivotalMarkers {
		for markerCombinationId, markerCombinationAnalysis := range markerCombinationAnalysisMap {
			if markerCombinationId.TicketType().OriginTicketType() != types.Trio {
				continue
			}
			for _, calculable := range markerCombinationAnalysis.Calculables() {
				allMarkerCalculablesMap[marker] = append(allMarkerCalculablesMap[marker], calculable)
			}
		}
	}

	for marker, calculables := range allMarkerCalculablesMap {
		pivotalMarkerAllCountOddsRangeMap[marker] = map[types.OddsRangeType]map[types.OddsRangeType]int{}
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

				if _, ok := pivotalMarkerAllCountOddsRangeMap[marker][pivotalMarkerOddsRange]; !ok {
					pivotalMarkerAllCountOddsRangeMap[marker][pivotalMarkerOddsRange] = map[types.OddsRangeType]int{}
				}

				for _, trioOdds := range trioOddsList {
					odds = trioOdds.Odds().InexactFloat64()
					if odds >= 1.0 && odds <= 9.9 {
						pivotalMarkerAllCountOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange1]++
					} else if odds >= 10.0 && odds <= 19.9 {
						pivotalMarkerAllCountOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange2]++
					} else if odds >= 20.0 && odds <= 29.9 {
						pivotalMarkerAllCountOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange3]++
					} else if odds >= 30.0 && odds <= 49.9 {
						pivotalMarkerAllCountOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange4]++
					} else if odds >= 50.0 && odds <= 99.9 {
						pivotalMarkerAllCountOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange5]++
					} else if odds >= 100.0 && odds <= 299.9 {
						pivotalMarkerAllCountOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange6]++
					} else if odds >= 300.0 && odds <= 499.9 {
						pivotalMarkerAllCountOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange7]++
					} else if odds >= 500.0 {
						pivotalMarkerAllCountOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange8]++
					}
				}
			}
		}
	}

	return pivotalMarkerAllCountOddsRangeMap, nil
}

func (s *spreadSheetTrioAnalysisRepository) getPivotalMarkerHitTotalOddsMap(
	ctx context.Context,
	markerCombinationAnalysisMap map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis,
	winRaceOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
	trioRaceOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
) (map[types.Marker]map[types.OddsRangeType]map[types.OddsRangeType]decimal.Decimal, error) {
	pivotalMarkerHitTotalOddsRangeMap := map[types.Marker]map[types.OddsRangeType]map[types.OddsRangeType]decimal.Decimal{}
	pivotalMarkers := []types.Marker{
		types.Favorite,
		types.Rival,
		types.BrackTriangle,
		types.WhiteTriangle,
		types.Star,
		types.Check,
	}

	allMarkerCalculablesMap := map[types.Marker][]*analysis_entity.Calculable{}
	for _, marker := range pivotalMarkers {
		for markerCombinationId, markerCombinationAnalysis := range markerCombinationAnalysisMap {
			if markerCombinationId.TicketType().OriginTicketType() != types.Trio {
				continue
			}
			for _, calculable := range markerCombinationAnalysis.Calculables() {
				allMarkerCalculablesMap[marker] = append(allMarkerCalculablesMap[marker], calculable)
			}
		}
	}

	for marker, calculables := range allMarkerCalculablesMap {
		pivotalMarkerHitTotalOddsRangeMap[marker] = map[types.OddsRangeType]map[types.OddsRangeType]decimal.Decimal{}
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

				if _, ok := pivotalMarkerHitTotalOddsRangeMap[marker][pivotalMarkerOddsRange]; !ok {
					pivotalMarkerHitTotalOddsRangeMap[marker][pivotalMarkerOddsRange] = map[types.OddsRangeType]decimal.Decimal{}
				}

				for _, trioOdds := range trioOddsList {
					odds = trioOdds.Odds().InexactFloat64()
					if odds >= 1.0 && odds <= 9.9 {
						pivotalMarkerHitTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange1] =
							pivotalMarkerHitTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange1].Add(trioOdds.Odds())
					} else if odds >= 10.0 && odds <= 19.9 {
						pivotalMarkerHitTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange2] =
							pivotalMarkerHitTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange2].Add(trioOdds.Odds())
					} else if odds >= 20.0 && odds <= 29.9 {
						pivotalMarkerHitTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange3] =
							pivotalMarkerHitTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange3].Add(trioOdds.Odds())
					} else if odds >= 30.0 && odds <= 49.9 {
						pivotalMarkerHitTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange4] =
							pivotalMarkerHitTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange4].Add(trioOdds.Odds())
					} else if odds >= 50.0 && odds <= 99.9 {
						pivotalMarkerHitTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange5] =
							pivotalMarkerHitTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange5].Add(trioOdds.Odds())
					} else if odds >= 100.0 && odds <= 299.9 {
						pivotalMarkerHitTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange6] =
							pivotalMarkerHitTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange6].Add(trioOdds.Odds())
					} else if odds >= 300.0 && odds <= 499.9 {
						pivotalMarkerHitTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange7] =
							pivotalMarkerHitTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange7].Add(trioOdds.Odds())
					} else if odds >= 500.0 {
						pivotalMarkerHitTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange8] =
							pivotalMarkerHitTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange8].Add(trioOdds.Odds())
					}
				}
			}
		}
	}

	return pivotalMarkerHitTotalOddsRangeMap, nil
}

func (s *spreadSheetTrioAnalysisRepository) getPivotalMarkerAllTotalOddsMap(
	ctx context.Context,
	markerCombinationAnalysisMap map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis,
	winRaceOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
	trioRaceOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
) (map[types.Marker]map[types.OddsRangeType]map[types.OddsRangeType]decimal.Decimal, error) {
	pivotalMarkerAllTotalOddsRangeMap := map[types.Marker]map[types.OddsRangeType]map[types.OddsRangeType]decimal.Decimal{}
	pivotalMarkers := []types.Marker{
		types.Favorite,
		types.Rival,
		types.BrackTriangle,
		types.WhiteTriangle,
		types.Star,
		types.Check,
	}

	allMarkerCalculablesMap := map[types.Marker][]*analysis_entity.Calculable{}
	for _, marker := range pivotalMarkers {
		for markerCombinationId, markerCombinationAnalysis := range markerCombinationAnalysisMap {
			if markerCombinationId.TicketType().OriginTicketType() != types.Trio {
				continue
			}
			for _, calculable := range markerCombinationAnalysis.Calculables() {
				allMarkerCalculablesMap[marker] = append(allMarkerCalculablesMap[marker], calculable)
			}
		}
	}

	for marker, calculables := range allMarkerCalculablesMap {
		pivotalMarkerAllTotalOddsRangeMap[marker] = map[types.OddsRangeType]map[types.OddsRangeType]decimal.Decimal{}
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

				if _, ok := pivotalMarkerAllTotalOddsRangeMap[marker][pivotalMarkerOddsRange]; !ok {
					pivotalMarkerAllTotalOddsRangeMap[marker][pivotalMarkerOddsRange] = map[types.OddsRangeType]decimal.Decimal{}
				}

				for _, trioOdds := range trioOddsList {
					odds = trioOdds.Odds().InexactFloat64()
					if odds >= 1.0 && odds <= 9.9 {
						pivotalMarkerAllTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange1] =
							pivotalMarkerAllTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange1].Add(trioOdds.Odds())
					} else if odds >= 10.0 && odds <= 19.9 {
						pivotalMarkerAllTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange2] =
							pivotalMarkerAllTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange2].Add(trioOdds.Odds())
					} else if odds >= 20.0 && odds <= 29.9 {
						pivotalMarkerAllTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange3] =
							pivotalMarkerAllTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange3].Add(trioOdds.Odds())
					} else if odds >= 30.0 && odds <= 49.9 {
						pivotalMarkerAllTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange4] =
							pivotalMarkerAllTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange4].Add(trioOdds.Odds())
					} else if odds >= 50.0 && odds <= 99.9 {
						pivotalMarkerAllTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange5] =
							pivotalMarkerAllTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange5].Add(trioOdds.Odds())
					} else if odds >= 100.0 && odds <= 299.9 {
						pivotalMarkerAllTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange6] =
							pivotalMarkerAllTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange6].Add(trioOdds.Odds())
					} else if odds >= 300.0 && odds <= 499.9 {
						pivotalMarkerAllTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange7] =
							pivotalMarkerAllTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange7].Add(trioOdds.Odds())
					} else if odds >= 500.0 {
						pivotalMarkerAllTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange8] =
							pivotalMarkerAllTotalOddsRangeMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange8].Add(trioOdds.Odds())
					}
				}
			}
		}
	}

	return pivotalMarkerAllTotalOddsRangeMap, nil
}

func (s *spreadSheetTrioAnalysisRepository) getTotalOdds(
	ctx context.Context,
	oddsList []decimal.Decimal,
) decimal.Decimal {
	totalOdds := decimal.NewFromFloat(0.0)
	for _, odds := range oddsList {
		totalOdds = totalOdds.Add(odds)
	}
	return totalOdds
}

func (s *spreadSheetTrioAnalysisRepository) aggregationOdds(
	ctx context.Context,
	markerCombinationAnalysisList []*spreadsheet_entity.MarkerCombinationAnalysis,
) (float64, float64, float64, float64, float64) {
	var rawOddsList []float64
	for _, markerCombinationAnalysis := range markerCombinationAnalysisList {
		for _, calculable := range markerCombinationAnalysis.Calculables() {
			rawOddsList = append(rawOddsList, calculable.Odds().InexactFloat64())
		}
	}

	if len(rawOddsList) == 0 {
		return 0, 0, 0, 0, 0
	}

	// 初期値を設定
	min, max := rawOddsList[0], rawOddsList[0]
	total := 0.0

	// 合計と最小値、最大値を計算
	for _, rawOdds := range rawOddsList {
		if rawOdds < min {
			min = rawOdds
		}
		if rawOdds > max {
			max = rawOdds
		}
		total += rawOdds
	}

	// 平均値を計算
	average := total / float64(len(rawOddsList))

	// 中央値の計算のためにスライスをソート
	sortedRawOddsList := make([]float64, len(rawOddsList))
	copy(sortedRawOddsList, rawOddsList)
	sort.Float64s(sortedRawOddsList)

	// 中央値を計算
	middle := len(sortedRawOddsList) / 2
	median := 0.0
	if len(sortedRawOddsList)%2 == 0 {
		median = (sortedRawOddsList[middle-1] + sortedRawOddsList[middle]) / 2
	} else {
		median = sortedRawOddsList[middle]
	}

	return total, max, min, average, median
}

func (s *spreadSheetTrioAnalysisRepository) Style(
	ctx context.Context,
	analysisData *spreadsheet_entity.AnalysisData,
) error {
	var requests []*sheets.Request
	allMarkerCombinationIds := analysisData.AllMarkerCombinationIds()
	markerCombinationMap := analysisData.MarkerCombinationFilterMap()
	rowGroupSize := len(s.createDefaultValuesList())

	for _, f := range analysisData.Filters() {
		trioAggregationAnalysisListMap, err := s.spreadSheetService.CreateTrioMarkerCombinationAggregationData(ctx, allMarkerCombinationIds, markerCombinationMap[f])
		if err != nil {
			return err
		}
		aggregationMarkerIds := make([]int, 0, len(trioAggregationAnalysisListMap))
		for id := range trioAggregationAnalysisListMap {
			if id.Value()%10 == types.NoMarker.Value() {
				// TODO 一旦無が含まれる場合をスルーする
				continue
			}
			aggregationMarkerIds = append(aggregationMarkerIds, id.Value())
		}
		sort.Ints(aggregationMarkerIds)

		for idx := range aggregationMarkerIds {
			requests = append(requests, []*sheets.Request{
				{
					RepeatCell: &sheets.RepeatCellRequest{
						Fields: "userEnteredFormat.backgroundColor",
						Range: &sheets.GridRange{
							SheetId:          s.spreadSheetConfig.SheetId(),
							StartColumnIndex: 0,
							StartRowIndex:    int64(idx * rowGroupSize),
							EndColumnIndex:   11,
							EndRowIndex:      int64(idx*rowGroupSize + 1),
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
							SheetId:          s.spreadSheetConfig.SheetId(),
							StartColumnIndex: 0,
							StartRowIndex:    int64(idx * rowGroupSize),
							EndColumnIndex:   11,
							EndRowIndex:      int64(idx*rowGroupSize + 1),
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
							SheetId:          s.spreadSheetConfig.SheetId(),
							StartColumnIndex: 0,
							StartRowIndex:    int64(idx * rowGroupSize),
							EndColumnIndex:   11,
							EndRowIndex:      int64(idx*rowGroupSize + 1),
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
							SheetId:          s.spreadSheetConfig.SheetId(),
							StartColumnIndex: 0,
							StartRowIndex:    int64(idx*rowGroupSize + 1),
							EndColumnIndex:   1,
							EndRowIndex:      int64(idx*rowGroupSize + 9),
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
							SheetId:          s.spreadSheetConfig.SheetId(),
							StartColumnIndex: 2,
							StartRowIndex:    int64(idx*rowGroupSize + 1),
							EndColumnIndex:   3,
							EndRowIndex:      int64(idx*rowGroupSize + rowGroupSize),
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
							StartRowIndex:    int64(idx*rowGroupSize + 1),
							EndColumnIndex:   1,
							EndRowIndex:      int64(idx*rowGroupSize + 9),
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
			}...)
			for i := 0; i < 9; i++ {
				requests = append(requests, []*sheets.Request{
					{
						RepeatCell: &sheets.RepeatCellRequest{
							Fields: "userEnteredFormat.backgroundColor",
							Range: &sheets.GridRange{
								SheetId:          s.spreadSheetConfig.SheetId(),
								StartColumnIndex: 2,
								StartRowIndex:    int64(1 + (i * 4) + idx*rowGroupSize),
								EndColumnIndex:   3,
								EndRowIndex:      int64(2 + (i * 4) + idx*rowGroupSize),
							},
							Cell: &sheets.CellData{
								UserEnteredFormat: &sheets.CellFormat{
									BackgroundColor: &sheets.Color{
										Red:   0.0,
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
								SheetId:          s.spreadSheetConfig.SheetId(),
								StartColumnIndex: 3,
								StartRowIndex:    int64(1 + (i * 4) + idx*rowGroupSize),
								EndColumnIndex:   11,
								EndRowIndex:      int64(2 + (i * 4) + idx*rowGroupSize),
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
								StartColumnIndex: 2,
								StartRowIndex:    int64(1 + (i * 4) + idx*rowGroupSize),
								EndColumnIndex:   11,
								EndRowIndex:      int64(2 + (i * 4) + idx*rowGroupSize),
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
								SheetId:          s.spreadSheetConfig.SheetId(),
								StartColumnIndex: 2,
								StartRowIndex:    int64(1 + (i * 4) + idx*rowGroupSize),
								EndColumnIndex:   11,
								EndRowIndex:      int64(2 + (i * 4) + idx*rowGroupSize),
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
								StartColumnIndex: 2,
								StartRowIndex:    int64(2 + (i * 4) + idx*rowGroupSize),
								EndColumnIndex:   3,
								EndRowIndex:      int64(5 + (i * 4) + idx*rowGroupSize),
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
				}...)
			}
		}

	}

	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.SpreadSheetId(), &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetTrioAnalysisRepository) Clear(ctx context.Context) error {
	requests := []*sheets.Request{
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Fields: "*",
				Range: &sheets.GridRange{
					SheetId:          s.spreadSheetConfig.SheetId(),
					StartColumnIndex: 0,
					StartRowIndex:    0,
					EndColumnIndex:   12,
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
