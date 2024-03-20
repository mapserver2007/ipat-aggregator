package infrastructure

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"google.golang.org/api/sheets/v4"
	"sort"
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
) error {
	rateFormatFunc := func(matchCount int, raceCount int) string {
		if raceCount == 0 {
			return "-"
		}
		return fmt.Sprintf("%.2f%%", float64(matchCount)*100/float64(raceCount))
	}

	allMarkerCombinationIds := analysisData.AllMarkerCombinationIds()
	markerCombinationMap := analysisData.MarkerCombinationFilterMap()

	var valuesList [][][]interface{}
	for _, f := range analysisData.Filters() {
		raceCount := analysisData.RaceCountFilterMap()[f][types.Trio]
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

		aggregationMarkerIndex := 0
		for _, rawId := range aggregationMarkerIds {
			defaultValuesList := s.createDefaultValuesList()
			position := len(defaultValuesList) * aggregationMarkerIndex
			pivotalMarker, err := types.NewMarker((rawId / 100) % 10)
			if err != nil {
				return err
			}

			valuesList = append(valuesList, defaultValuesList...)
			aggregationMarkerCombinationId := types.MarkerCombinationId(rawId)
			aggregationAnalysisList, ok := trioAggregationAnalysisListMap[aggregationMarkerCombinationId]
			if !ok {
				return fmt.Errorf("aggregationAnalysisList not found: %v", aggregationMarkerCombinationId)
			}
			// aggregationAnalysisListの中は、集約される前の印組合せが全部入っていて、calculablesはlistになっていて合算はされていない状態
			aggregationHitPivotalOddsRangeMap := s.createHitTrioOddsRangeMap(ctx, aggregationAnalysisList, pivotalMarker)
			total, max, min, average, median := s.aggregationOdds(ctx, aggregationAnalysisList)
			matchCount := 0
			for _, pivotalOddsRange := range aggregationHitPivotalOddsRangeMap {
				for _, count := range pivotalOddsRange {
					matchCount += count
				}
			}

			// 印組合せの概要集計
			for i := position; i < len(defaultValuesList)+position; i++ {
				switch i - position {
				case 0:
					valuesList[i][0][0] = fmt.Sprintf("%s / %s", aggregationMarkerCombinationId.String(), f.String())
				case 1:
					valuesList[i][0][1] = raceCount
				case 2:
					valuesList[i][0][1] = matchCount
				case 3:
					valuesList[i][0][1] = rateFormatFunc(matchCount, raceCount)
				case 4:
					valuesList[i][0][1] = fmt.Sprintf("%.2f%%", (total/float64(raceCount)*100)/10) // 10点買いなので10で割る
				case 5:
					valuesList[i][0][1] = fmt.Sprintf("%.2f", max)
				case 6:
					valuesList[i][0][1] = fmt.Sprintf("%.2f", min)
				case 7:
					valuesList[i][0][1] = fmt.Sprintf("%.2f", average)
				case 8:
					valuesList[i][0][1] = fmt.Sprintf("%.2f", median)
				}
			}

			// 印組合せのオッズ幅の集計
			for i := position; i < len(defaultValuesList)+position; i++ {
				switch i - position {
				case 1:
					valuesList[i][0][2] = "単全部"
					valuesList[i][0][3] = types.TrioOddsRange1.String()
					valuesList[i][0][4] = types.TrioOddsRange2.String()
					valuesList[i][0][5] = types.TrioOddsRange3.String()
					valuesList[i][0][6] = types.TrioOddsRange4.String()
					valuesList[i][0][7] = types.TrioOddsRange5.String()
					valuesList[i][0][8] = types.TrioOddsRange6.String()
					valuesList[i][0][9] = types.TrioOddsRange7.String()
					valuesList[i][0][10] = types.TrioOddsRange8.String()
				case 2:
					allWinOddsRangeMap := map[types.OddsRangeType]int{}
					for _, oddsRange := range aggregationHitPivotalOddsRangeMap {
						allWinOddsRangeMap[types.TrioOddsRange1] += oddsRange[types.TrioOddsRange1]
						allWinOddsRangeMap[types.TrioOddsRange2] += oddsRange[types.TrioOddsRange2]
						allWinOddsRangeMap[types.TrioOddsRange3] += oddsRange[types.TrioOddsRange3]
						allWinOddsRangeMap[types.TrioOddsRange4] += oddsRange[types.TrioOddsRange4]
						allWinOddsRangeMap[types.TrioOddsRange5] += oddsRange[types.TrioOddsRange5]
						allWinOddsRangeMap[types.TrioOddsRange6] += oddsRange[types.TrioOddsRange6]
						allWinOddsRangeMap[types.TrioOddsRange7] += oddsRange[types.TrioOddsRange7]
						allWinOddsRangeMap[types.TrioOddsRange8] += oddsRange[types.TrioOddsRange8]
					}
					valuesList[i][0][3] = allWinOddsRangeMap[types.TrioOddsRange1]
					valuesList[i][0][4] = allWinOddsRangeMap[types.TrioOddsRange2]
					valuesList[i][0][5] = allWinOddsRangeMap[types.TrioOddsRange3]
					valuesList[i][0][6] = allWinOddsRangeMap[types.TrioOddsRange4]
					valuesList[i][0][7] = allWinOddsRangeMap[types.TrioOddsRange5]
					valuesList[i][0][8] = allWinOddsRangeMap[types.TrioOddsRange6]
					valuesList[i][0][9] = allWinOddsRangeMap[types.TrioOddsRange7]
					valuesList[i][0][10] = allWinOddsRangeMap[types.TrioOddsRange8]
				case 3:
					allWinOddsRangeMap := map[types.OddsRangeType]int{}
					for _, oddsRange := range aggregationHitPivotalOddsRangeMap {
						allWinOddsRangeMap[types.TrioOddsRange1] += oddsRange[types.TrioOddsRange1]
						allWinOddsRangeMap[types.TrioOddsRange2] += oddsRange[types.TrioOddsRange2]
						allWinOddsRangeMap[types.TrioOddsRange3] += oddsRange[types.TrioOddsRange3]
						allWinOddsRangeMap[types.TrioOddsRange4] += oddsRange[types.TrioOddsRange4]
						allWinOddsRangeMap[types.TrioOddsRange5] += oddsRange[types.TrioOddsRange5]
						allWinOddsRangeMap[types.TrioOddsRange6] += oddsRange[types.TrioOddsRange6]
						allWinOddsRangeMap[types.TrioOddsRange7] += oddsRange[types.TrioOddsRange7]
						allWinOddsRangeMap[types.TrioOddsRange8] += oddsRange[types.TrioOddsRange8]
					}
					matchCount = 0
					for _, count := range allWinOddsRangeMap {
						matchCount += count
					}
					valuesList[i][0][3] = rateFormatFunc(allWinOddsRangeMap[types.TrioOddsRange1], matchCount)
					valuesList[i][0][4] = rateFormatFunc(allWinOddsRangeMap[types.TrioOddsRange2], matchCount)
					valuesList[i][0][5] = rateFormatFunc(allWinOddsRangeMap[types.TrioOddsRange3], matchCount)
					valuesList[i][0][6] = rateFormatFunc(allWinOddsRangeMap[types.TrioOddsRange4], matchCount)
					valuesList[i][0][7] = rateFormatFunc(allWinOddsRangeMap[types.TrioOddsRange5], matchCount)
					valuesList[i][0][8] = rateFormatFunc(allWinOddsRangeMap[types.TrioOddsRange6], matchCount)
					valuesList[i][0][9] = rateFormatFunc(allWinOddsRangeMap[types.TrioOddsRange7], matchCount)
					valuesList[i][0][10] = rateFormatFunc(allWinOddsRangeMap[types.TrioOddsRange8], matchCount)
				case 4:
					allWinOddsRangeMap := map[types.OddsRangeType]int{}
					for _, oddsRange := range aggregationHitPivotalOddsRangeMap {
						allWinOddsRangeMap[types.TrioOddsRange1] += oddsRange[types.TrioOddsRange1]
						allWinOddsRangeMap[types.TrioOddsRange2] += oddsRange[types.TrioOddsRange2]
						allWinOddsRangeMap[types.TrioOddsRange3] += oddsRange[types.TrioOddsRange3]
						allWinOddsRangeMap[types.TrioOddsRange4] += oddsRange[types.TrioOddsRange4]
						allWinOddsRangeMap[types.TrioOddsRange5] += oddsRange[types.TrioOddsRange5]
						allWinOddsRangeMap[types.TrioOddsRange6] += oddsRange[types.TrioOddsRange6]
						allWinOddsRangeMap[types.TrioOddsRange7] += oddsRange[types.TrioOddsRange7]
						allWinOddsRangeMap[types.TrioOddsRange8] += oddsRange[types.TrioOddsRange8]
					}
					valuesList[i][0][3] = rateFormatFunc(allWinOddsRangeMap[types.TrioOddsRange1], raceCount)
					valuesList[i][0][4] = rateFormatFunc(allWinOddsRangeMap[types.TrioOddsRange2], raceCount)
					valuesList[i][0][5] = rateFormatFunc(allWinOddsRangeMap[types.TrioOddsRange3], raceCount)
					valuesList[i][0][6] = rateFormatFunc(allWinOddsRangeMap[types.TrioOddsRange4], raceCount)
					valuesList[i][0][7] = rateFormatFunc(allWinOddsRangeMap[types.TrioOddsRange5], raceCount)
					valuesList[i][0][8] = rateFormatFunc(allWinOddsRangeMap[types.TrioOddsRange6], raceCount)
					valuesList[i][0][9] = rateFormatFunc(allWinOddsRangeMap[types.TrioOddsRange7], raceCount)
					valuesList[i][0][10] = rateFormatFunc(allWinOddsRangeMap[types.TrioOddsRange8], raceCount)
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
					valuesList[i][0][3] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange1]
					valuesList[i][0][4] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange2]
					valuesList[i][0][5] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange3]
					valuesList[i][0][6] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange4]
					valuesList[i][0][7] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange5]
					valuesList[i][0][8] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange6]
					valuesList[i][0][9] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange7]
					valuesList[i][0][10] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange8]
				case 7:
					matchCount = 0
					for _, count := range aggregationHitPivotalOddsRangeMap[types.WinOddsRange1] {
						matchCount += count
					}
					valuesList[i][0][3] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange1], matchCount)
					valuesList[i][0][4] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange2], matchCount)
					valuesList[i][0][5] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange3], matchCount)
					valuesList[i][0][6] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange4], matchCount)
					valuesList[i][0][7] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange5], matchCount)
					valuesList[i][0][8] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange6], matchCount)
					valuesList[i][0][9] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange7], matchCount)
					valuesList[i][0][10] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange8], matchCount)
				case 8:
					valuesList[i][0][3] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange1], raceCount)
					valuesList[i][0][4] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange2], raceCount)
					valuesList[i][0][5] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange3], raceCount)
					valuesList[i][0][6] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange4], raceCount)
					valuesList[i][0][7] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange5], raceCount)
					valuesList[i][0][8] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange6], raceCount)
					valuesList[i][0][9] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange7], raceCount)
					valuesList[i][0][10] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange1][types.TrioOddsRange8], raceCount)
				case 9:
					valuesList[i][0][2] = fmt.Sprintf("単%s", types.WinOddsRange2.String())
					valuesList[i][0][3] = types.TrioOddsRange1.String()
					valuesList[i][0][4] = types.TrioOddsRange2.String()
					valuesList[i][0][5] = types.TrioOddsRange3.String()
					valuesList[i][0][6] = types.TrioOddsRange4.String()
					valuesList[i][0][7] = types.TrioOddsRange5.String()
					valuesList[i][0][8] = types.TrioOddsRange6.String()
					valuesList[i][0][9] = types.TrioOddsRange7.String()
					valuesList[i][0][10] = types.TrioOddsRange8.String()
				case 10:
					valuesList[i][0][3] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange1]
					valuesList[i][0][4] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange2]
					valuesList[i][0][5] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange3]
					valuesList[i][0][6] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange4]
					valuesList[i][0][7] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange5]
					valuesList[i][0][8] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange6]
					valuesList[i][0][9] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange7]
					valuesList[i][0][10] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange8]
				case 11:
					matchCount = 0
					for _, count := range aggregationHitPivotalOddsRangeMap[types.WinOddsRange2] {
						matchCount += count
					}
					valuesList[i][0][3] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange1], matchCount)
					valuesList[i][0][4] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange2], matchCount)
					valuesList[i][0][5] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange3], matchCount)
					valuesList[i][0][6] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange4], matchCount)
					valuesList[i][0][7] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange5], matchCount)
					valuesList[i][0][8] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange6], matchCount)
					valuesList[i][0][9] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange7], matchCount)
					valuesList[i][0][10] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange8], matchCount)
				case 12:
					valuesList[i][0][3] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange1], raceCount)
					valuesList[i][0][4] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange2], raceCount)
					valuesList[i][0][5] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange3], raceCount)
					valuesList[i][0][6] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange4], raceCount)
					valuesList[i][0][7] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange5], raceCount)
					valuesList[i][0][8] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange6], raceCount)
					valuesList[i][0][9] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange7], raceCount)
					valuesList[i][0][10] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange2][types.TrioOddsRange8], raceCount)
				case 13:
					valuesList[i][0][2] = fmt.Sprintf("単%s", types.WinOddsRange3.String())
					valuesList[i][0][3] = types.TrioOddsRange1.String()
					valuesList[i][0][4] = types.TrioOddsRange2.String()
					valuesList[i][0][5] = types.TrioOddsRange3.String()
					valuesList[i][0][6] = types.TrioOddsRange4.String()
					valuesList[i][0][7] = types.TrioOddsRange5.String()
					valuesList[i][0][8] = types.TrioOddsRange6.String()
					valuesList[i][0][9] = types.TrioOddsRange7.String()
					valuesList[i][0][10] = types.TrioOddsRange8.String()
				case 14:
					valuesList[i][0][3] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange1]
					valuesList[i][0][4] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange2]
					valuesList[i][0][5] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange3]
					valuesList[i][0][6] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange4]
					valuesList[i][0][7] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange5]
					valuesList[i][0][8] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange6]
					valuesList[i][0][9] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange7]
					valuesList[i][0][10] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange8]
				case 15:
					matchCount = 0
					for _, count := range aggregationHitPivotalOddsRangeMap[types.WinOddsRange3] {
						matchCount += count
					}
					valuesList[i][0][3] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange1], matchCount)
					valuesList[i][0][4] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange2], matchCount)
					valuesList[i][0][5] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange3], matchCount)
					valuesList[i][0][6] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange4], matchCount)
					valuesList[i][0][7] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange5], matchCount)
					valuesList[i][0][8] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange6], matchCount)
					valuesList[i][0][9] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange7], matchCount)
					valuesList[i][0][10] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange8], matchCount)
				case 16:
					valuesList[i][0][3] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange1], raceCount)
					valuesList[i][0][4] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange2], raceCount)
					valuesList[i][0][5] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange3], raceCount)
					valuesList[i][0][6] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange4], raceCount)
					valuesList[i][0][7] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange5], raceCount)
					valuesList[i][0][8] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange6], raceCount)
					valuesList[i][0][9] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange7], raceCount)
					valuesList[i][0][10] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange3][types.TrioOddsRange8], raceCount)
				case 17:
					valuesList[i][0][2] = fmt.Sprintf("単%s", types.WinOddsRange4.String())
					valuesList[i][0][3] = types.TrioOddsRange1.String()
					valuesList[i][0][4] = types.TrioOddsRange2.String()
					valuesList[i][0][5] = types.TrioOddsRange3.String()
					valuesList[i][0][6] = types.TrioOddsRange4.String()
					valuesList[i][0][7] = types.TrioOddsRange5.String()
					valuesList[i][0][8] = types.TrioOddsRange6.String()
					valuesList[i][0][9] = types.TrioOddsRange7.String()
					valuesList[i][0][10] = types.TrioOddsRange8.String()
				case 18:
					valuesList[i][0][3] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange1]
					valuesList[i][0][4] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange2]
					valuesList[i][0][5] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange3]
					valuesList[i][0][6] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange4]
					valuesList[i][0][7] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange5]
					valuesList[i][0][8] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange6]
					valuesList[i][0][9] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange7]
					valuesList[i][0][10] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange8]
				case 19:
					matchCount = 0
					for _, count := range aggregationHitPivotalOddsRangeMap[types.WinOddsRange4] {
						matchCount += count
					}
					valuesList[i][0][3] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange1], matchCount)
					valuesList[i][0][4] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange2], matchCount)
					valuesList[i][0][5] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange3], matchCount)
					valuesList[i][0][6] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange4], matchCount)
					valuesList[i][0][7] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange5], matchCount)
					valuesList[i][0][8] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange6], matchCount)
					valuesList[i][0][9] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange7], matchCount)
					valuesList[i][0][10] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange8], matchCount)
				case 20:
					valuesList[i][0][3] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange1], raceCount)
					valuesList[i][0][4] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange2], raceCount)
					valuesList[i][0][5] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange3], raceCount)
					valuesList[i][0][6] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange4], raceCount)
					valuesList[i][0][7] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange5], raceCount)
					valuesList[i][0][8] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange6], raceCount)
					valuesList[i][0][9] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange7], raceCount)
					valuesList[i][0][10] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange4][types.TrioOddsRange8], raceCount)
				case 21:
					valuesList[i][0][2] = fmt.Sprintf("単%s", types.WinOddsRange5.String())
					valuesList[i][0][3] = types.TrioOddsRange1.String()
					valuesList[i][0][4] = types.TrioOddsRange2.String()
					valuesList[i][0][5] = types.TrioOddsRange3.String()
					valuesList[i][0][6] = types.TrioOddsRange4.String()
					valuesList[i][0][7] = types.TrioOddsRange5.String()
					valuesList[i][0][8] = types.TrioOddsRange6.String()
					valuesList[i][0][9] = types.TrioOddsRange7.String()
					valuesList[i][0][10] = types.TrioOddsRange8.String()
				case 22:
					valuesList[i][0][3] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange1]
					valuesList[i][0][4] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange2]
					valuesList[i][0][5] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange3]
					valuesList[i][0][6] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange4]
					valuesList[i][0][7] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange5]
					valuesList[i][0][8] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange6]
					valuesList[i][0][9] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange7]
					valuesList[i][0][10] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange8]
				case 23:
					matchCount = 0
					for _, count := range aggregationHitPivotalOddsRangeMap[types.WinOddsRange5] {
						matchCount += count
					}
					valuesList[i][0][3] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange1], matchCount)
					valuesList[i][0][4] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange2], matchCount)
					valuesList[i][0][5] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange3], matchCount)
					valuesList[i][0][6] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange4], matchCount)
					valuesList[i][0][7] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange5], matchCount)
					valuesList[i][0][8] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange6], matchCount)
					valuesList[i][0][9] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange7], matchCount)
					valuesList[i][0][10] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange8], matchCount)
				case 24:
					valuesList[i][0][3] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange1], raceCount)
					valuesList[i][0][4] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange2], raceCount)
					valuesList[i][0][5] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange3], raceCount)
					valuesList[i][0][6] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange4], raceCount)
					valuesList[i][0][7] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange5], raceCount)
					valuesList[i][0][8] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange6], raceCount)
					valuesList[i][0][9] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange7], raceCount)
					valuesList[i][0][10] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange5][types.TrioOddsRange8], raceCount)
				case 25:
					valuesList[i][0][2] = fmt.Sprintf("単%s", types.WinOddsRange6.String())
					valuesList[i][0][3] = types.TrioOddsRange1.String()
					valuesList[i][0][4] = types.TrioOddsRange2.String()
					valuesList[i][0][5] = types.TrioOddsRange3.String()
					valuesList[i][0][6] = types.TrioOddsRange4.String()
					valuesList[i][0][7] = types.TrioOddsRange5.String()
					valuesList[i][0][8] = types.TrioOddsRange6.String()
					valuesList[i][0][9] = types.TrioOddsRange7.String()
					valuesList[i][0][10] = types.TrioOddsRange8.String()
				case 26:
					valuesList[i][0][3] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange1]
					valuesList[i][0][4] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange2]
					valuesList[i][0][5] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange3]
					valuesList[i][0][6] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange4]
					valuesList[i][0][7] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange5]
					valuesList[i][0][8] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange6]
					valuesList[i][0][9] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange7]
					valuesList[i][0][10] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange8]
				case 27:
					matchCount = 0
					for _, count := range aggregationHitPivotalOddsRangeMap[types.WinOddsRange6] {
						matchCount += count
					}
					valuesList[i][0][3] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange1], matchCount)
					valuesList[i][0][4] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange2], matchCount)
					valuesList[i][0][5] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange3], matchCount)
					valuesList[i][0][6] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange4], matchCount)
					valuesList[i][0][7] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange5], matchCount)
					valuesList[i][0][8] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange6], matchCount)
					valuesList[i][0][9] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange7], matchCount)
					valuesList[i][0][10] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange8], matchCount)
				case 28:
					valuesList[i][0][3] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange1], raceCount)
					valuesList[i][0][4] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange2], raceCount)
					valuesList[i][0][5] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange3], raceCount)
					valuesList[i][0][6] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange4], raceCount)
					valuesList[i][0][7] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange5], raceCount)
					valuesList[i][0][8] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange6], raceCount)
					valuesList[i][0][9] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange7], raceCount)
					valuesList[i][0][10] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange6][types.TrioOddsRange8], raceCount)
				case 29:
					valuesList[i][0][2] = fmt.Sprintf("単%s", types.WinOddsRange7.String())
					valuesList[i][0][3] = types.TrioOddsRange1.String()
					valuesList[i][0][4] = types.TrioOddsRange2.String()
					valuesList[i][0][5] = types.TrioOddsRange3.String()
					valuesList[i][0][6] = types.TrioOddsRange4.String()
					valuesList[i][0][7] = types.TrioOddsRange5.String()
					valuesList[i][0][8] = types.TrioOddsRange6.String()
					valuesList[i][0][9] = types.TrioOddsRange7.String()
					valuesList[i][0][10] = types.TrioOddsRange8.String()
				case 30:
					valuesList[i][0][3] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange1]
					valuesList[i][0][4] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange2]
					valuesList[i][0][5] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange3]
					valuesList[i][0][6] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange4]
					valuesList[i][0][7] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange5]
					valuesList[i][0][8] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange6]
					valuesList[i][0][9] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange7]
					valuesList[i][0][10] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange8]
				case 31:
					matchCount = 0
					for _, count := range aggregationHitPivotalOddsRangeMap[types.WinOddsRange7] {
						matchCount += count
					}
					valuesList[i][0][3] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange1], matchCount)
					valuesList[i][0][4] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange2], matchCount)
					valuesList[i][0][5] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange3], matchCount)
					valuesList[i][0][6] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange4], matchCount)
					valuesList[i][0][7] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange5], matchCount)
					valuesList[i][0][8] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange6], matchCount)
					valuesList[i][0][9] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange7], matchCount)
					valuesList[i][0][10] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange8], matchCount)
				case 32:
					valuesList[i][0][3] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange1], raceCount)
					valuesList[i][0][4] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange2], raceCount)
					valuesList[i][0][5] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange3], raceCount)
					valuesList[i][0][6] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange4], raceCount)
					valuesList[i][0][7] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange5], raceCount)
					valuesList[i][0][8] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange6], raceCount)
					valuesList[i][0][9] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange7], raceCount)
					valuesList[i][0][10] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange7][types.TrioOddsRange8], raceCount)
				case 33:
					valuesList[i][0][2] = fmt.Sprintf("単%s", types.WinOddsRange8.String())
					valuesList[i][0][3] = types.TrioOddsRange1.String()
					valuesList[i][0][4] = types.TrioOddsRange2.String()
					valuesList[i][0][5] = types.TrioOddsRange3.String()
					valuesList[i][0][6] = types.TrioOddsRange4.String()
					valuesList[i][0][7] = types.TrioOddsRange5.String()
					valuesList[i][0][8] = types.TrioOddsRange6.String()
					valuesList[i][0][9] = types.TrioOddsRange7.String()
					valuesList[i][0][10] = types.TrioOddsRange8.String()
				case 34:
					valuesList[i][0][3] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange1]
					valuesList[i][0][4] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange2]
					valuesList[i][0][5] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange3]
					valuesList[i][0][6] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange4]
					valuesList[i][0][7] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange5]
					valuesList[i][0][8] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange6]
					valuesList[i][0][9] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange7]
					valuesList[i][0][10] = aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange8]
				case 35:
					matchCount = 0
					for _, count := range aggregationHitPivotalOddsRangeMap[types.WinOddsRange8] {
						matchCount += count
					}
					valuesList[i][0][3] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange1], matchCount)
					valuesList[i][0][4] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange2], matchCount)
					valuesList[i][0][5] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange3], matchCount)
					valuesList[i][0][6] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange4], matchCount)
					valuesList[i][0][7] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange5], matchCount)
					valuesList[i][0][8] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange6], matchCount)
					valuesList[i][0][9] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange7], matchCount)
					valuesList[i][0][10] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange8], matchCount)
				case 36:
					valuesList[i][0][3] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange1], raceCount)
					valuesList[i][0][4] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange2], raceCount)
					valuesList[i][0][5] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange3], raceCount)
					valuesList[i][0][6] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange4], raceCount)
					valuesList[i][0][7] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange5], raceCount)
					valuesList[i][0][8] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange6], raceCount)
					valuesList[i][0][9] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange7], raceCount)
					valuesList[i][0][10] = rateFormatFunc(aggregationHitPivotalOddsRangeMap[types.WinOddsRange8][types.TrioOddsRange8], raceCount)
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
			"的中割合",
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
			"的中全割合",
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
			"的中割合",
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
			"的中全割合",
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
			"",
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
			"的中割合",
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
			"的中全割合",
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
			"",
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
			"的中割合",
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
			"的中全割合",
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
			"",
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
			"的中割合",
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
			"的中全割合",
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
			"",
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
			"的中割合",
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
			"的中全割合",
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
			"的中割合",
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
			"的中全割合",
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
			"的中割合",
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
			"的中全割合",
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
			"的中割合",
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
			"的中全割合",
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

func (s *spreadSheetTrioAnalysisRepository) createHitTrioOddsRangeMap(
	ctx context.Context,
	markerCombinationAnalysisList []*spreadsheet_entity.MarkerCombinationAnalysis,
	pivotalMarker types.Marker,
) map[types.OddsRangeType]map[types.OddsRangeType]int {
	pivotalOddsRangeMap := map[types.OddsRangeType]map[types.OddsRangeType]int{}

	for _, markerCombinationAnalysis := range markerCombinationAnalysisList {
		for _, calculable := range markerCombinationAnalysis.Calculables() {
			var (
				isContainPivotal bool
				pivotalOddsRange types.OddsRangeType
			)
			for _, pivotal := range calculable.Pivotals() {
				if pivotalMarker == pivotal.Marker() {
					isContainPivotal = true
					odds := pivotal.Odds().InexactFloat64()
					if odds >= 1.0 && odds <= 1.5 {
						pivotalOddsRange = types.WinOddsRange1
					} else if odds >= 1.6 && odds <= 2.0 {
						pivotalOddsRange = types.WinOddsRange2
					} else if odds >= 2.1 && odds <= 2.9 {
						pivotalOddsRange = types.WinOddsRange3
					} else if odds >= 3.0 && odds <= 4.9 {
						pivotalOddsRange = types.WinOddsRange4
					} else if odds >= 5.0 && odds <= 9.9 {
						pivotalOddsRange = types.WinOddsRange5
					} else if odds >= 10.0 && odds <= 19.9 {
						pivotalOddsRange = types.WinOddsRange6
					} else if odds >= 20.0 && odds <= 49.9 {
						pivotalOddsRange = types.WinOddsRange7
					} else if odds >= 50.0 {
						pivotalOddsRange = types.WinOddsRange8
					}
				}
			}

			if !isContainPivotal {
				continue
			}

			if _, ok := pivotalOddsRangeMap[pivotalOddsRange]; !ok {
				pivotalOddsRangeMap[pivotalOddsRange] = map[types.OddsRangeType]int{}
			}

			odds := calculable.Odds().InexactFloat64()
			if odds >= 1.0 && odds <= 9.9 {
				if _, ok := pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange1]; !ok {
					pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange1] = 0
				}
				pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange1]++
			} else if odds >= 10.0 && odds <= 19.9 {
				if _, ok := pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange2]; !ok {
					pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange2] = 0
				}
				pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange2]++
			} else if odds >= 20.0 && odds <= 29.9 {
				if _, ok := pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange3]; !ok {
					pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange3] = 0
				}
				pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange3]++
			} else if odds >= 30.0 && odds <= 49.9 {
				if _, ok := pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange4]; !ok {
					pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange4] = 0
				}
				pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange4]++
			} else if odds >= 50.0 && odds <= 99.9 {
				if _, ok := pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange5]; !ok {
					pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange5] = 0
				}
				pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange5]++
			} else if odds >= 100.0 && odds <= 299.9 {
				if _, ok := pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange6]; !ok {
					pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange6] = 0
				}
				pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange6]++
			} else if odds >= 300.0 && odds <= 499.9 {
				if _, ok := pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange7]; !ok {
					pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange7] = 0
				}
				pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange7]++
			} else if odds >= 500.0 {
				if _, ok := pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange8]; !ok {
					pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange8] = 0
				}
				pivotalOddsRangeMap[pivotalOddsRange][types.TrioOddsRange8]++
			}
		}
	}

	return pivotalOddsRangeMap
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
