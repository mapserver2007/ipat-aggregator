package infrastructure

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"google.golang.org/api/sheets/v4"
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
	filters []filter.Id,
) error {
	valuesList := make([][][]interface{}, 0)
	valuesList = append(valuesList, [][]interface{}{
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
		},
	})
	valuesList = append(valuesList, [][]interface{}{
		{
			"",
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

	//rateFormatFunc := func(matchCount int, raceCount int) string {
	//	if raceCount == 0 {
	//		return "-"
	//	}
	//	return fmt.Sprintf("%.2f%%", float64(matchCount)*100/float64(raceCount))
	//}

	allMarkerCombinationIds := analysisData.AllMarkerCombinationIds()
	markerCombinationMap := analysisData.MarkerCombinationFilterMap()
	raceCountMap := analysisData.OddsRangeRaceCountFilterMap()

	oddsRanges := []types.OddsRangeType{
		types.TrioOddsRange1,
		types.TrioOddsRange2,
		types.TrioOddsRange3,
		types.TrioOddsRange4,
		types.TrioOddsRange5,
		types.TrioOddsRange6,
		types.TrioOddsRange7,
		types.TrioOddsRange8,
	}
	_ = oddsRanges

	for _, f := range filters {
		trioAggregationAnalysisListMap, trioAggregationRaceCountMap, err := s.spreadSheetService.CreateTrioMarkerCombinationAggregationData(ctx, allMarkerCombinationIds, markerCombinationMap[f], raceCountMap[f])
		if err != nil {
			return err
		}

		raceCount := 0
		for _, oddsRangeRaceCountMap := range trioAggregationRaceCountMap {
			for _, count := range oddsRangeRaceCountMap {
				raceCount += count
			}
		}

		fmt.Println(raceCount)

		_ = trioAggregationAnalysisListMap
		_ = trioAggregationRaceCountMap
	}

	//for _, markerCombinationId := range allMarkerCombinationIds {
	//	f := filter.All // TODO 各フィルタを取得する
	//
	//	// TODO あとでフィルタ
	//	oddsRangeMap := raceCountMap[f]
	//	fmt.Println(oddsRangeMap)
	//
	//	data, ok := markerCombinationMap[f][markerCombinationId]
	//	if ok && markerCombinationId.TicketType().OriginTicketType() == types.Trio {
	//		s.spreadSheetService.CreateTrioMarkerCombinationAggregationData(ctx, data)
	//	}
	//}

	return nil
}

func (s *spreadSheetTrioAnalysisRepository) Style(ctx context.Context, analysisData *spreadsheet_entity.AnalysisData, filters []filter.Id) error {
	//TODO implement me
	panic("implement me")
}

func (s *spreadSheetTrioAnalysisRepository) Clear(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
