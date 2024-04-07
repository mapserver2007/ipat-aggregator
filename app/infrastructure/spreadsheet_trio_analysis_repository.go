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
	for _, f := range analysisData.Filters() {
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

		// 軸に対するレース単位のオッズ幅ごとの出現回数
		raceOddsRangeCountMap, err := s.getRaceOddsRangeCountMap(
			ctx,
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
			markerCombinationMap[f],
			winRaceResultOddsMap,
			trioRaceResultOddsMap,
		)
		if err != nil {
			return err
		}

		aggregationMarkerIndex := 0
		for _, rawMarkerId := range []int{1, 2, 3, 4, 5, 6} {
			pivotalMarker, _ := types.NewMarker(rawMarkerId)
			defaultValuesList := s.createDefaultValuesList()
			position := len(defaultValuesList) * aggregationMarkerIndex
			valuesList = append(valuesList, defaultValuesList...)

			// 印組合せのオッズ幅の集計
			for i := position; i < len(defaultValuesList)+position; i++ {
				switch i - position {
				case 0:
					valuesList[i][0][0] = fmt.Sprintf("%s-印-印, フィルタ条件: %s", pivotalMarker.String(), f.String())
				case 1:
					valuesList[i][0][0] = "単全部"
					valuesList[i][0][1] = types.TrioOddsRange1.String()
					valuesList[i][0][2] = types.TrioOddsRange2.String()
					valuesList[i][0][3] = types.TrioOddsRange3.String()
					valuesList[i][0][4] = types.TrioOddsRange4.String()
					valuesList[i][0][5] = types.TrioOddsRange5.String()
					valuesList[i][0][6] = types.TrioOddsRange6.String()
					valuesList[i][0][7] = types.TrioOddsRange7.String()
					valuesList[i][0][8] = types.TrioOddsRange8.String()
				case 2:
					oddsMap := map[types.OddsRangeType]int{}
					for _, oddsRangeMap := range raceHitOddsRangeCountMap[pivotalMarker] {
						oddsMap[types.TrioOddsRange1] += oddsRangeMap[types.TrioOddsRange1]
						oddsMap[types.TrioOddsRange2] += oddsRangeMap[types.TrioOddsRange2]
						oddsMap[types.TrioOddsRange3] += oddsRangeMap[types.TrioOddsRange3]
						oddsMap[types.TrioOddsRange4] += oddsRangeMap[types.TrioOddsRange4]
						oddsMap[types.TrioOddsRange5] += oddsRangeMap[types.TrioOddsRange5]
						oddsMap[types.TrioOddsRange6] += oddsRangeMap[types.TrioOddsRange6]
						oddsMap[types.TrioOddsRange7] += oddsRangeMap[types.TrioOddsRange7]
						oddsMap[types.TrioOddsRange8] += oddsRangeMap[types.TrioOddsRange8]
					}
					valuesList[i][0][1] = oddsMap[types.TrioOddsRange1]
					valuesList[i][0][2] = oddsMap[types.TrioOddsRange2]
					valuesList[i][0][3] = oddsMap[types.TrioOddsRange3]
					valuesList[i][0][4] = oddsMap[types.TrioOddsRange4]
					valuesList[i][0][5] = oddsMap[types.TrioOddsRange5]
					valuesList[i][0][6] = oddsMap[types.TrioOddsRange6]
					valuesList[i][0][7] = oddsMap[types.TrioOddsRange7]
					valuesList[i][0][8] = oddsMap[types.TrioOddsRange8]
				case 3:
					oddsMap := map[types.OddsRangeType]int{}
					for _, oddsRange := range raceOddsRangeCountMap[pivotalMarker] {
						oddsMap[types.TrioOddsRange1] += oddsRange[types.TrioOddsRange1]
						oddsMap[types.TrioOddsRange2] += oddsRange[types.TrioOddsRange2]
						oddsMap[types.TrioOddsRange3] += oddsRange[types.TrioOddsRange3]
						oddsMap[types.TrioOddsRange4] += oddsRange[types.TrioOddsRange4]
						oddsMap[types.TrioOddsRange5] += oddsRange[types.TrioOddsRange5]
						oddsMap[types.TrioOddsRange6] += oddsRange[types.TrioOddsRange6]
						oddsMap[types.TrioOddsRange7] += oddsRange[types.TrioOddsRange7]
						oddsMap[types.TrioOddsRange8] += oddsRange[types.TrioOddsRange8]
					}
					valuesList[i][0][1] = oddsMap[types.TrioOddsRange1]
					valuesList[i][0][2] = oddsMap[types.TrioOddsRange2]
					valuesList[i][0][3] = oddsMap[types.TrioOddsRange3]
					valuesList[i][0][4] = oddsMap[types.TrioOddsRange4]
					valuesList[i][0][5] = oddsMap[types.TrioOddsRange5]
					valuesList[i][0][6] = oddsMap[types.TrioOddsRange6]
					valuesList[i][0][7] = oddsMap[types.TrioOddsRange7]
					valuesList[i][0][8] = oddsMap[types.TrioOddsRange8]
				case 4:
					hitCountMap := map[types.OddsRangeType]int{}
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
					allCountMap := map[types.OddsRangeType]int{}
					for _, oddsRange := range markerAllOddsRangeCountMap[pivotalMarker] {
						allCountMap[types.TrioOddsRange1] += oddsRange[types.TrioOddsRange1]
						allCountMap[types.TrioOddsRange2] += oddsRange[types.TrioOddsRange2]
						allCountMap[types.TrioOddsRange3] += oddsRange[types.TrioOddsRange3]
						allCountMap[types.TrioOddsRange4] += oddsRange[types.TrioOddsRange4]
						allCountMap[types.TrioOddsRange5] += oddsRange[types.TrioOddsRange5]
						allCountMap[types.TrioOddsRange6] += oddsRange[types.TrioOddsRange6]
						allCountMap[types.TrioOddsRange7] += oddsRange[types.TrioOddsRange7]
						allCountMap[types.TrioOddsRange8] += oddsRange[types.TrioOddsRange8]
					}
					valuesList[i][0][1] = hitRateFormat(hitCountMap[types.TrioOddsRange1], allCountMap[types.TrioOddsRange1])
					valuesList[i][0][2] = hitRateFormat(hitCountMap[types.TrioOddsRange2], allCountMap[types.TrioOddsRange2])
					valuesList[i][0][3] = hitRateFormat(hitCountMap[types.TrioOddsRange3], allCountMap[types.TrioOddsRange3])
					valuesList[i][0][4] = hitRateFormat(hitCountMap[types.TrioOddsRange4], allCountMap[types.TrioOddsRange4])
					valuesList[i][0][5] = hitRateFormat(hitCountMap[types.TrioOddsRange5], allCountMap[types.TrioOddsRange5])
					valuesList[i][0][6] = hitRateFormat(hitCountMap[types.TrioOddsRange6], allCountMap[types.TrioOddsRange6])
					valuesList[i][0][7] = hitRateFormat(hitCountMap[types.TrioOddsRange7], allCountMap[types.TrioOddsRange7])
					valuesList[i][0][8] = hitRateFormat(hitCountMap[types.TrioOddsRange8], allCountMap[types.TrioOddsRange8])
				case 5:
					hitOddsMap := map[types.OddsRangeType]decimal.Decimal{}
					for _, oddsRange := range pivotalMarkerHitTotalOddsMap[pivotalMarker] {
						hitOddsMap[types.TrioOddsRange1] = hitOddsMap[types.TrioOddsRange1].Add(oddsRange[types.TrioOddsRange1])
						hitOddsMap[types.TrioOddsRange2] = hitOddsMap[types.TrioOddsRange2].Add(oddsRange[types.TrioOddsRange2])
						hitOddsMap[types.TrioOddsRange3] = hitOddsMap[types.TrioOddsRange3].Add(oddsRange[types.TrioOddsRange3])
						hitOddsMap[types.TrioOddsRange4] = hitOddsMap[types.TrioOddsRange4].Add(oddsRange[types.TrioOddsRange4])
						hitOddsMap[types.TrioOddsRange5] = hitOddsMap[types.TrioOddsRange5].Add(oddsRange[types.TrioOddsRange5])
						hitOddsMap[types.TrioOddsRange6] = hitOddsMap[types.TrioOddsRange6].Add(oddsRange[types.TrioOddsRange6])
						hitOddsMap[types.TrioOddsRange7] = hitOddsMap[types.TrioOddsRange7].Add(oddsRange[types.TrioOddsRange7])
						hitOddsMap[types.TrioOddsRange8] = hitOddsMap[types.TrioOddsRange8].Add(oddsRange[types.TrioOddsRange8])
					}
					allCountMap := map[types.OddsRangeType]int{}
					for _, oddsRange := range markerAllOddsRangeCountMap[pivotalMarker] {
						allCountMap[types.TrioOddsRange1] += oddsRange[types.TrioOddsRange1]
						allCountMap[types.TrioOddsRange2] += oddsRange[types.TrioOddsRange2]
						allCountMap[types.TrioOddsRange3] += oddsRange[types.TrioOddsRange3]
						allCountMap[types.TrioOddsRange4] += oddsRange[types.TrioOddsRange4]
						allCountMap[types.TrioOddsRange5] += oddsRange[types.TrioOddsRange5]
						allCountMap[types.TrioOddsRange6] += oddsRange[types.TrioOddsRange6]
						allCountMap[types.TrioOddsRange7] += oddsRange[types.TrioOddsRange7]
						allCountMap[types.TrioOddsRange8] += oddsRange[types.TrioOddsRange8]
					}
					valuesList[i][0][1] = payoutRateFormat(hitOddsMap[types.TrioOddsRange1].InexactFloat64(), allCountMap[types.TrioOddsRange1])
					valuesList[i][0][2] = payoutRateFormat(hitOddsMap[types.TrioOddsRange2].InexactFloat64(), allCountMap[types.TrioOddsRange2])
					valuesList[i][0][3] = payoutRateFormat(hitOddsMap[types.TrioOddsRange3].InexactFloat64(), allCountMap[types.TrioOddsRange3])
					valuesList[i][0][4] = payoutRateFormat(hitOddsMap[types.TrioOddsRange4].InexactFloat64(), allCountMap[types.TrioOddsRange4])
					valuesList[i][0][5] = payoutRateFormat(hitOddsMap[types.TrioOddsRange5].InexactFloat64(), allCountMap[types.TrioOddsRange5])
					valuesList[i][0][6] = payoutRateFormat(hitOddsMap[types.TrioOddsRange6].InexactFloat64(), allCountMap[types.TrioOddsRange6])
					valuesList[i][0][7] = payoutRateFormat(hitOddsMap[types.TrioOddsRange7].InexactFloat64(), allCountMap[types.TrioOddsRange7])
					valuesList[i][0][8] = payoutRateFormat(hitOddsMap[types.TrioOddsRange8].InexactFloat64(), allCountMap[types.TrioOddsRange8])
				case 6:
					valuesList[i][0][0] = fmt.Sprintf("単%s", types.WinOddsRange1.String())
					valuesList[i][0][1] = types.TrioOddsRange1.String()
					valuesList[i][0][2] = types.TrioOddsRange2.String()
					valuesList[i][0][3] = types.TrioOddsRange3.String()
					valuesList[i][0][4] = types.TrioOddsRange4.String()
					valuesList[i][0][5] = types.TrioOddsRange5.String()
					valuesList[i][0][6] = types.TrioOddsRange6.String()
					valuesList[i][0][7] = types.TrioOddsRange7.String()
					valuesList[i][0][8] = types.TrioOddsRange8.String()
				case 7:
					oddsMap := raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange1]
					valuesList[i][0][1] = oddsMap[types.TrioOddsRange1]
					valuesList[i][0][2] = oddsMap[types.TrioOddsRange2]
					valuesList[i][0][3] = oddsMap[types.TrioOddsRange3]
					valuesList[i][0][4] = oddsMap[types.TrioOddsRange4]
					valuesList[i][0][5] = oddsMap[types.TrioOddsRange5]
					valuesList[i][0][6] = oddsMap[types.TrioOddsRange6]
					valuesList[i][0][7] = oddsMap[types.TrioOddsRange7]
					valuesList[i][0][8] = oddsMap[types.TrioOddsRange8]
				case 8:
					oddsMap := raceOddsRangeCountMap[pivotalMarker][types.WinOddsRange1]
					valuesList[i][0][1] = oddsMap[types.TrioOddsRange1]
					valuesList[i][0][2] = oddsMap[types.TrioOddsRange2]
					valuesList[i][0][3] = oddsMap[types.TrioOddsRange3]
					valuesList[i][0][4] = oddsMap[types.TrioOddsRange4]
					valuesList[i][0][5] = oddsMap[types.TrioOddsRange5]
					valuesList[i][0][6] = oddsMap[types.TrioOddsRange6]
					valuesList[i][0][7] = oddsMap[types.TrioOddsRange7]
					valuesList[i][0][8] = oddsMap[types.TrioOddsRange8]
				case 9:
					hitCountMap := raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange1]
					allCountMap := markerAllOddsRangeCountMap[pivotalMarker][types.WinOddsRange1]
					valuesList[i][0][1] = hitRateFormat(hitCountMap[types.TrioOddsRange1], allCountMap[types.TrioOddsRange1])
					valuesList[i][0][2] = hitRateFormat(hitCountMap[types.TrioOddsRange2], allCountMap[types.TrioOddsRange2])
					valuesList[i][0][3] = hitRateFormat(hitCountMap[types.TrioOddsRange3], allCountMap[types.TrioOddsRange3])
					valuesList[i][0][4] = hitRateFormat(hitCountMap[types.TrioOddsRange4], allCountMap[types.TrioOddsRange4])
					valuesList[i][0][5] = hitRateFormat(hitCountMap[types.TrioOddsRange5], allCountMap[types.TrioOddsRange5])
					valuesList[i][0][6] = hitRateFormat(hitCountMap[types.TrioOddsRange6], allCountMap[types.TrioOddsRange6])
					valuesList[i][0][7] = hitRateFormat(hitCountMap[types.TrioOddsRange7], allCountMap[types.TrioOddsRange7])
					valuesList[i][0][8] = hitRateFormat(hitCountMap[types.TrioOddsRange8], allCountMap[types.TrioOddsRange8])
				case 10:
					hitOddsMap := pivotalMarkerHitTotalOddsMap[pivotalMarker][types.WinOddsRange1]
					allCountMap := markerAllOddsRangeCountMap[pivotalMarker][types.WinOddsRange1]
					valuesList[i][0][1] = payoutRateFormat(hitOddsMap[types.TrioOddsRange1].InexactFloat64(), allCountMap[types.TrioOddsRange1])
					valuesList[i][0][2] = payoutRateFormat(hitOddsMap[types.TrioOddsRange2].InexactFloat64(), allCountMap[types.TrioOddsRange2])
					valuesList[i][0][3] = payoutRateFormat(hitOddsMap[types.TrioOddsRange3].InexactFloat64(), allCountMap[types.TrioOddsRange3])
					valuesList[i][0][4] = payoutRateFormat(hitOddsMap[types.TrioOddsRange4].InexactFloat64(), allCountMap[types.TrioOddsRange4])
					valuesList[i][0][5] = payoutRateFormat(hitOddsMap[types.TrioOddsRange5].InexactFloat64(), allCountMap[types.TrioOddsRange5])
					valuesList[i][0][6] = payoutRateFormat(hitOddsMap[types.TrioOddsRange6].InexactFloat64(), allCountMap[types.TrioOddsRange6])
					valuesList[i][0][7] = payoutRateFormat(hitOddsMap[types.TrioOddsRange7].InexactFloat64(), allCountMap[types.TrioOddsRange7])
					valuesList[i][0][8] = payoutRateFormat(hitOddsMap[types.TrioOddsRange8].InexactFloat64(), allCountMap[types.TrioOddsRange8])
				case 11:
					valuesList[i][0][0] = fmt.Sprintf("単%s", types.WinOddsRange2.String())
					valuesList[i][0][1] = types.TrioOddsRange1.String()
					valuesList[i][0][2] = types.TrioOddsRange2.String()
					valuesList[i][0][3] = types.TrioOddsRange3.String()
					valuesList[i][0][4] = types.TrioOddsRange4.String()
					valuesList[i][0][5] = types.TrioOddsRange5.String()
					valuesList[i][0][6] = types.TrioOddsRange6.String()
					valuesList[i][0][7] = types.TrioOddsRange7.String()
					valuesList[i][0][8] = types.TrioOddsRange8.String()
				case 12:
					oddsMap := raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange2]
					valuesList[i][0][1] = oddsMap[types.TrioOddsRange1]
					valuesList[i][0][2] = oddsMap[types.TrioOddsRange2]
					valuesList[i][0][3] = oddsMap[types.TrioOddsRange3]
					valuesList[i][0][4] = oddsMap[types.TrioOddsRange4]
					valuesList[i][0][5] = oddsMap[types.TrioOddsRange5]
					valuesList[i][0][6] = oddsMap[types.TrioOddsRange6]
					valuesList[i][0][7] = oddsMap[types.TrioOddsRange7]
					valuesList[i][0][8] = oddsMap[types.TrioOddsRange8]
				case 13:
					oddsMap := raceOddsRangeCountMap[pivotalMarker][types.WinOddsRange2]
					valuesList[i][0][1] = oddsMap[types.TrioOddsRange1]
					valuesList[i][0][2] = oddsMap[types.TrioOddsRange2]
					valuesList[i][0][3] = oddsMap[types.TrioOddsRange3]
					valuesList[i][0][4] = oddsMap[types.TrioOddsRange4]
					valuesList[i][0][5] = oddsMap[types.TrioOddsRange5]
					valuesList[i][0][6] = oddsMap[types.TrioOddsRange6]
					valuesList[i][0][7] = oddsMap[types.TrioOddsRange7]
					valuesList[i][0][8] = oddsMap[types.TrioOddsRange8]
				case 14:
					hitCountMap := raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange2]
					allCountMap := markerAllOddsRangeCountMap[pivotalMarker][types.WinOddsRange2]
					valuesList[i][0][1] = hitRateFormat(hitCountMap[types.TrioOddsRange1], allCountMap[types.TrioOddsRange1])
					valuesList[i][0][2] = hitRateFormat(hitCountMap[types.TrioOddsRange2], allCountMap[types.TrioOddsRange2])
					valuesList[i][0][3] = hitRateFormat(hitCountMap[types.TrioOddsRange3], allCountMap[types.TrioOddsRange3])
					valuesList[i][0][4] = hitRateFormat(hitCountMap[types.TrioOddsRange4], allCountMap[types.TrioOddsRange4])
					valuesList[i][0][5] = hitRateFormat(hitCountMap[types.TrioOddsRange5], allCountMap[types.TrioOddsRange5])
					valuesList[i][0][6] = hitRateFormat(hitCountMap[types.TrioOddsRange6], allCountMap[types.TrioOddsRange6])
					valuesList[i][0][7] = hitRateFormat(hitCountMap[types.TrioOddsRange7], allCountMap[types.TrioOddsRange7])
					valuesList[i][0][8] = hitRateFormat(hitCountMap[types.TrioOddsRange8], allCountMap[types.TrioOddsRange8])
				case 15:
					hitOddsMap := pivotalMarkerHitTotalOddsMap[pivotalMarker][types.WinOddsRange2]
					allCountMap := markerAllOddsRangeCountMap[pivotalMarker][types.WinOddsRange2]
					valuesList[i][0][1] = payoutRateFormat(hitOddsMap[types.TrioOddsRange1].InexactFloat64(), allCountMap[types.TrioOddsRange1])
					valuesList[i][0][2] = payoutRateFormat(hitOddsMap[types.TrioOddsRange2].InexactFloat64(), allCountMap[types.TrioOddsRange2])
					valuesList[i][0][3] = payoutRateFormat(hitOddsMap[types.TrioOddsRange3].InexactFloat64(), allCountMap[types.TrioOddsRange3])
					valuesList[i][0][4] = payoutRateFormat(hitOddsMap[types.TrioOddsRange4].InexactFloat64(), allCountMap[types.TrioOddsRange4])
					valuesList[i][0][5] = payoutRateFormat(hitOddsMap[types.TrioOddsRange5].InexactFloat64(), allCountMap[types.TrioOddsRange5])
					valuesList[i][0][6] = payoutRateFormat(hitOddsMap[types.TrioOddsRange6].InexactFloat64(), allCountMap[types.TrioOddsRange6])
					valuesList[i][0][7] = payoutRateFormat(hitOddsMap[types.TrioOddsRange7].InexactFloat64(), allCountMap[types.TrioOddsRange7])
					valuesList[i][0][8] = payoutRateFormat(hitOddsMap[types.TrioOddsRange8].InexactFloat64(), allCountMap[types.TrioOddsRange8])
				case 16:
					valuesList[i][0][0] = fmt.Sprintf("単%s", types.WinOddsRange3.String())
					valuesList[i][0][1] = types.TrioOddsRange1.String()
					valuesList[i][0][2] = types.TrioOddsRange2.String()
					valuesList[i][0][3] = types.TrioOddsRange3.String()
					valuesList[i][0][4] = types.TrioOddsRange4.String()
					valuesList[i][0][5] = types.TrioOddsRange5.String()
					valuesList[i][0][6] = types.TrioOddsRange6.String()
					valuesList[i][0][7] = types.TrioOddsRange7.String()
					valuesList[i][0][8] = types.TrioOddsRange8.String()
				case 17:
					oddsMap := raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange3]
					valuesList[i][0][1] = oddsMap[types.TrioOddsRange1]
					valuesList[i][0][2] = oddsMap[types.TrioOddsRange2]
					valuesList[i][0][3] = oddsMap[types.TrioOddsRange3]
					valuesList[i][0][4] = oddsMap[types.TrioOddsRange4]
					valuesList[i][0][5] = oddsMap[types.TrioOddsRange5]
					valuesList[i][0][6] = oddsMap[types.TrioOddsRange6]
					valuesList[i][0][7] = oddsMap[types.TrioOddsRange7]
					valuesList[i][0][8] = oddsMap[types.TrioOddsRange8]
				case 18:
					oddsMap := raceOddsRangeCountMap[pivotalMarker][types.WinOddsRange3]
					valuesList[i][0][1] = oddsMap[types.TrioOddsRange1]
					valuesList[i][0][2] = oddsMap[types.TrioOddsRange2]
					valuesList[i][0][3] = oddsMap[types.TrioOddsRange3]
					valuesList[i][0][4] = oddsMap[types.TrioOddsRange4]
					valuesList[i][0][5] = oddsMap[types.TrioOddsRange5]
					valuesList[i][0][6] = oddsMap[types.TrioOddsRange6]
					valuesList[i][0][7] = oddsMap[types.TrioOddsRange7]
					valuesList[i][0][8] = oddsMap[types.TrioOddsRange8]
				case 19:
					hitCountMap := raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange3]
					allCountMap := markerAllOddsRangeCountMap[pivotalMarker][types.WinOddsRange3]
					valuesList[i][0][1] = hitRateFormat(hitCountMap[types.TrioOddsRange1], allCountMap[types.TrioOddsRange1])
					valuesList[i][0][2] = hitRateFormat(hitCountMap[types.TrioOddsRange2], allCountMap[types.TrioOddsRange2])
					valuesList[i][0][3] = hitRateFormat(hitCountMap[types.TrioOddsRange3], allCountMap[types.TrioOddsRange3])
					valuesList[i][0][4] = hitRateFormat(hitCountMap[types.TrioOddsRange4], allCountMap[types.TrioOddsRange4])
					valuesList[i][0][5] = hitRateFormat(hitCountMap[types.TrioOddsRange5], allCountMap[types.TrioOddsRange5])
					valuesList[i][0][6] = hitRateFormat(hitCountMap[types.TrioOddsRange6], allCountMap[types.TrioOddsRange6])
					valuesList[i][0][7] = hitRateFormat(hitCountMap[types.TrioOddsRange7], allCountMap[types.TrioOddsRange7])
					valuesList[i][0][8] = hitRateFormat(hitCountMap[types.TrioOddsRange8], allCountMap[types.TrioOddsRange8])
				case 20:
					hitOddsMap := pivotalMarkerHitTotalOddsMap[pivotalMarker][types.WinOddsRange3]
					allCountMap := markerAllOddsRangeCountMap[pivotalMarker][types.WinOddsRange3]
					valuesList[i][0][1] = payoutRateFormat(hitOddsMap[types.TrioOddsRange1].InexactFloat64(), allCountMap[types.TrioOddsRange1])
					valuesList[i][0][2] = payoutRateFormat(hitOddsMap[types.TrioOddsRange2].InexactFloat64(), allCountMap[types.TrioOddsRange2])
					valuesList[i][0][3] = payoutRateFormat(hitOddsMap[types.TrioOddsRange3].InexactFloat64(), allCountMap[types.TrioOddsRange3])
					valuesList[i][0][4] = payoutRateFormat(hitOddsMap[types.TrioOddsRange4].InexactFloat64(), allCountMap[types.TrioOddsRange4])
					valuesList[i][0][5] = payoutRateFormat(hitOddsMap[types.TrioOddsRange5].InexactFloat64(), allCountMap[types.TrioOddsRange5])
					valuesList[i][0][6] = payoutRateFormat(hitOddsMap[types.TrioOddsRange6].InexactFloat64(), allCountMap[types.TrioOddsRange6])
					valuesList[i][0][7] = payoutRateFormat(hitOddsMap[types.TrioOddsRange7].InexactFloat64(), allCountMap[types.TrioOddsRange7])
					valuesList[i][0][8] = payoutRateFormat(hitOddsMap[types.TrioOddsRange8].InexactFloat64(), allCountMap[types.TrioOddsRange8])
				case 21:
					valuesList[i][0][0] = fmt.Sprintf("単%s", types.WinOddsRange4.String())
					valuesList[i][0][1] = types.TrioOddsRange1.String()
					valuesList[i][0][2] = types.TrioOddsRange2.String()
					valuesList[i][0][3] = types.TrioOddsRange3.String()
					valuesList[i][0][4] = types.TrioOddsRange4.String()
					valuesList[i][0][5] = types.TrioOddsRange5.String()
					valuesList[i][0][6] = types.TrioOddsRange6.String()
					valuesList[i][0][7] = types.TrioOddsRange7.String()
					valuesList[i][0][8] = types.TrioOddsRange8.String()
				case 22:
					oddsMap := raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange4]
					valuesList[i][0][1] = oddsMap[types.TrioOddsRange1]
					valuesList[i][0][2] = oddsMap[types.TrioOddsRange2]
					valuesList[i][0][3] = oddsMap[types.TrioOddsRange3]
					valuesList[i][0][4] = oddsMap[types.TrioOddsRange4]
					valuesList[i][0][5] = oddsMap[types.TrioOddsRange5]
					valuesList[i][0][6] = oddsMap[types.TrioOddsRange6]
					valuesList[i][0][7] = oddsMap[types.TrioOddsRange7]
					valuesList[i][0][8] = oddsMap[types.TrioOddsRange8]
				case 23:
					oddsMap := raceOddsRangeCountMap[pivotalMarker][types.WinOddsRange4]
					valuesList[i][0][1] = oddsMap[types.TrioOddsRange1]
					valuesList[i][0][2] = oddsMap[types.TrioOddsRange2]
					valuesList[i][0][3] = oddsMap[types.TrioOddsRange3]
					valuesList[i][0][4] = oddsMap[types.TrioOddsRange4]
					valuesList[i][0][5] = oddsMap[types.TrioOddsRange5]
					valuesList[i][0][6] = oddsMap[types.TrioOddsRange6]
					valuesList[i][0][7] = oddsMap[types.TrioOddsRange7]
					valuesList[i][0][8] = oddsMap[types.TrioOddsRange8]
				case 24:
					hitCountMap := raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange4]
					allCountMap := markerAllOddsRangeCountMap[pivotalMarker][types.WinOddsRange4]
					valuesList[i][0][1] = hitRateFormat(hitCountMap[types.TrioOddsRange1], allCountMap[types.TrioOddsRange1])
					valuesList[i][0][2] = hitRateFormat(hitCountMap[types.TrioOddsRange2], allCountMap[types.TrioOddsRange2])
					valuesList[i][0][3] = hitRateFormat(hitCountMap[types.TrioOddsRange3], allCountMap[types.TrioOddsRange3])
					valuesList[i][0][4] = hitRateFormat(hitCountMap[types.TrioOddsRange4], allCountMap[types.TrioOddsRange4])
					valuesList[i][0][5] = hitRateFormat(hitCountMap[types.TrioOddsRange5], allCountMap[types.TrioOddsRange5])
					valuesList[i][0][6] = hitRateFormat(hitCountMap[types.TrioOddsRange6], allCountMap[types.TrioOddsRange6])
					valuesList[i][0][7] = hitRateFormat(hitCountMap[types.TrioOddsRange7], allCountMap[types.TrioOddsRange7])
					valuesList[i][0][8] = hitRateFormat(hitCountMap[types.TrioOddsRange8], allCountMap[types.TrioOddsRange8])
				case 25:
					hitOddsMap := pivotalMarkerHitTotalOddsMap[pivotalMarker][types.WinOddsRange4]
					allCountMap := markerAllOddsRangeCountMap[pivotalMarker][types.WinOddsRange4]
					valuesList[i][0][1] = payoutRateFormat(hitOddsMap[types.TrioOddsRange1].InexactFloat64(), allCountMap[types.TrioOddsRange1])
					valuesList[i][0][2] = payoutRateFormat(hitOddsMap[types.TrioOddsRange2].InexactFloat64(), allCountMap[types.TrioOddsRange2])
					valuesList[i][0][3] = payoutRateFormat(hitOddsMap[types.TrioOddsRange3].InexactFloat64(), allCountMap[types.TrioOddsRange3])
					valuesList[i][0][4] = payoutRateFormat(hitOddsMap[types.TrioOddsRange4].InexactFloat64(), allCountMap[types.TrioOddsRange4])
					valuesList[i][0][5] = payoutRateFormat(hitOddsMap[types.TrioOddsRange5].InexactFloat64(), allCountMap[types.TrioOddsRange5])
					valuesList[i][0][6] = payoutRateFormat(hitOddsMap[types.TrioOddsRange6].InexactFloat64(), allCountMap[types.TrioOddsRange6])
					valuesList[i][0][7] = payoutRateFormat(hitOddsMap[types.TrioOddsRange7].InexactFloat64(), allCountMap[types.TrioOddsRange7])
					valuesList[i][0][8] = payoutRateFormat(hitOddsMap[types.TrioOddsRange8].InexactFloat64(), allCountMap[types.TrioOddsRange8])
				case 26:
					valuesList[i][0][0] = fmt.Sprintf("単%s", types.WinOddsRange5.String())
					valuesList[i][0][1] = types.TrioOddsRange1.String()
					valuesList[i][0][2] = types.TrioOddsRange2.String()
					valuesList[i][0][3] = types.TrioOddsRange3.String()
					valuesList[i][0][4] = types.TrioOddsRange4.String()
					valuesList[i][0][5] = types.TrioOddsRange5.String()
					valuesList[i][0][6] = types.TrioOddsRange6.String()
					valuesList[i][0][7] = types.TrioOddsRange7.String()
					valuesList[i][0][8] = types.TrioOddsRange8.String()
				case 27:
					oddsMap := raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange5]
					valuesList[i][0][1] = oddsMap[types.TrioOddsRange1]
					valuesList[i][0][2] = oddsMap[types.TrioOddsRange2]
					valuesList[i][0][3] = oddsMap[types.TrioOddsRange3]
					valuesList[i][0][4] = oddsMap[types.TrioOddsRange4]
					valuesList[i][0][5] = oddsMap[types.TrioOddsRange5]
					valuesList[i][0][6] = oddsMap[types.TrioOddsRange6]
					valuesList[i][0][7] = oddsMap[types.TrioOddsRange7]
					valuesList[i][0][8] = oddsMap[types.TrioOddsRange8]
				case 28:
					oddsMap := raceOddsRangeCountMap[pivotalMarker][types.WinOddsRange5]
					valuesList[i][0][1] = oddsMap[types.TrioOddsRange1]
					valuesList[i][0][2] = oddsMap[types.TrioOddsRange2]
					valuesList[i][0][3] = oddsMap[types.TrioOddsRange3]
					valuesList[i][0][4] = oddsMap[types.TrioOddsRange4]
					valuesList[i][0][5] = oddsMap[types.TrioOddsRange5]
					valuesList[i][0][6] = oddsMap[types.TrioOddsRange6]
					valuesList[i][0][7] = oddsMap[types.TrioOddsRange7]
					valuesList[i][0][8] = oddsMap[types.TrioOddsRange8]
				case 29:
					hitCountMap := raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange5]
					allCountMap := markerAllOddsRangeCountMap[pivotalMarker][types.WinOddsRange5]
					valuesList[i][0][1] = hitRateFormat(hitCountMap[types.TrioOddsRange1], allCountMap[types.TrioOddsRange1])
					valuesList[i][0][2] = hitRateFormat(hitCountMap[types.TrioOddsRange2], allCountMap[types.TrioOddsRange2])
					valuesList[i][0][3] = hitRateFormat(hitCountMap[types.TrioOddsRange3], allCountMap[types.TrioOddsRange3])
					valuesList[i][0][4] = hitRateFormat(hitCountMap[types.TrioOddsRange4], allCountMap[types.TrioOddsRange4])
					valuesList[i][0][5] = hitRateFormat(hitCountMap[types.TrioOddsRange5], allCountMap[types.TrioOddsRange5])
					valuesList[i][0][6] = hitRateFormat(hitCountMap[types.TrioOddsRange6], allCountMap[types.TrioOddsRange6])
					valuesList[i][0][7] = hitRateFormat(hitCountMap[types.TrioOddsRange7], allCountMap[types.TrioOddsRange7])
					valuesList[i][0][8] = hitRateFormat(hitCountMap[types.TrioOddsRange8], allCountMap[types.TrioOddsRange8])
				case 30:
					hitOddsMap := pivotalMarkerHitTotalOddsMap[pivotalMarker][types.WinOddsRange5]
					allCountMap := markerAllOddsRangeCountMap[pivotalMarker][types.WinOddsRange5]
					valuesList[i][0][1] = payoutRateFormat(hitOddsMap[types.TrioOddsRange1].InexactFloat64(), allCountMap[types.TrioOddsRange1])
					valuesList[i][0][2] = payoutRateFormat(hitOddsMap[types.TrioOddsRange2].InexactFloat64(), allCountMap[types.TrioOddsRange2])
					valuesList[i][0][3] = payoutRateFormat(hitOddsMap[types.TrioOddsRange3].InexactFloat64(), allCountMap[types.TrioOddsRange3])
					valuesList[i][0][4] = payoutRateFormat(hitOddsMap[types.TrioOddsRange4].InexactFloat64(), allCountMap[types.TrioOddsRange4])
					valuesList[i][0][5] = payoutRateFormat(hitOddsMap[types.TrioOddsRange5].InexactFloat64(), allCountMap[types.TrioOddsRange5])
					valuesList[i][0][6] = payoutRateFormat(hitOddsMap[types.TrioOddsRange6].InexactFloat64(), allCountMap[types.TrioOddsRange6])
					valuesList[i][0][7] = payoutRateFormat(hitOddsMap[types.TrioOddsRange7].InexactFloat64(), allCountMap[types.TrioOddsRange7])
					valuesList[i][0][8] = payoutRateFormat(hitOddsMap[types.TrioOddsRange8].InexactFloat64(), allCountMap[types.TrioOddsRange8])
				case 31:
					valuesList[i][0][0] = fmt.Sprintf("単%s", types.WinOddsRange6.String())
					valuesList[i][0][1] = types.TrioOddsRange1.String()
					valuesList[i][0][2] = types.TrioOddsRange2.String()
					valuesList[i][0][3] = types.TrioOddsRange3.String()
					valuesList[i][0][4] = types.TrioOddsRange4.String()
					valuesList[i][0][5] = types.TrioOddsRange5.String()
					valuesList[i][0][6] = types.TrioOddsRange6.String()
					valuesList[i][0][7] = types.TrioOddsRange7.String()
					valuesList[i][0][8] = types.TrioOddsRange8.String()
				case 32:
					oddsMap := raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange6]
					valuesList[i][0][1] = oddsMap[types.TrioOddsRange1]
					valuesList[i][0][2] = oddsMap[types.TrioOddsRange2]
					valuesList[i][0][3] = oddsMap[types.TrioOddsRange3]
					valuesList[i][0][4] = oddsMap[types.TrioOddsRange4]
					valuesList[i][0][5] = oddsMap[types.TrioOddsRange5]
					valuesList[i][0][6] = oddsMap[types.TrioOddsRange6]
					valuesList[i][0][7] = oddsMap[types.TrioOddsRange7]
					valuesList[i][0][8] = oddsMap[types.TrioOddsRange8]
				case 33:
					oddsMap := raceOddsRangeCountMap[pivotalMarker][types.WinOddsRange6]
					valuesList[i][0][1] = oddsMap[types.TrioOddsRange1]
					valuesList[i][0][2] = oddsMap[types.TrioOddsRange2]
					valuesList[i][0][3] = oddsMap[types.TrioOddsRange3]
					valuesList[i][0][4] = oddsMap[types.TrioOddsRange4]
					valuesList[i][0][5] = oddsMap[types.TrioOddsRange5]
					valuesList[i][0][6] = oddsMap[types.TrioOddsRange6]
					valuesList[i][0][7] = oddsMap[types.TrioOddsRange7]
					valuesList[i][0][8] = oddsMap[types.TrioOddsRange8]
				case 34:
					hitCountMap := raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange6]
					allCountMap := markerAllOddsRangeCountMap[pivotalMarker][types.WinOddsRange6]
					valuesList[i][0][1] = hitRateFormat(hitCountMap[types.TrioOddsRange1], allCountMap[types.TrioOddsRange1])
					valuesList[i][0][2] = hitRateFormat(hitCountMap[types.TrioOddsRange2], allCountMap[types.TrioOddsRange2])
					valuesList[i][0][3] = hitRateFormat(hitCountMap[types.TrioOddsRange3], allCountMap[types.TrioOddsRange3])
					valuesList[i][0][4] = hitRateFormat(hitCountMap[types.TrioOddsRange4], allCountMap[types.TrioOddsRange4])
					valuesList[i][0][5] = hitRateFormat(hitCountMap[types.TrioOddsRange5], allCountMap[types.TrioOddsRange5])
					valuesList[i][0][6] = hitRateFormat(hitCountMap[types.TrioOddsRange6], allCountMap[types.TrioOddsRange6])
					valuesList[i][0][7] = hitRateFormat(hitCountMap[types.TrioOddsRange7], allCountMap[types.TrioOddsRange7])
					valuesList[i][0][8] = hitRateFormat(hitCountMap[types.TrioOddsRange8], allCountMap[types.TrioOddsRange8])
				case 35:
					hitOddsMap := pivotalMarkerHitTotalOddsMap[pivotalMarker][types.WinOddsRange6]
					allCountMap := markerAllOddsRangeCountMap[pivotalMarker][types.WinOddsRange6]
					valuesList[i][0][1] = payoutRateFormat(hitOddsMap[types.TrioOddsRange1].InexactFloat64(), allCountMap[types.TrioOddsRange1])
					valuesList[i][0][2] = payoutRateFormat(hitOddsMap[types.TrioOddsRange2].InexactFloat64(), allCountMap[types.TrioOddsRange2])
					valuesList[i][0][3] = payoutRateFormat(hitOddsMap[types.TrioOddsRange3].InexactFloat64(), allCountMap[types.TrioOddsRange3])
					valuesList[i][0][4] = payoutRateFormat(hitOddsMap[types.TrioOddsRange4].InexactFloat64(), allCountMap[types.TrioOddsRange4])
					valuesList[i][0][5] = payoutRateFormat(hitOddsMap[types.TrioOddsRange5].InexactFloat64(), allCountMap[types.TrioOddsRange5])
					valuesList[i][0][6] = payoutRateFormat(hitOddsMap[types.TrioOddsRange6].InexactFloat64(), allCountMap[types.TrioOddsRange6])
					valuesList[i][0][7] = payoutRateFormat(hitOddsMap[types.TrioOddsRange7].InexactFloat64(), allCountMap[types.TrioOddsRange7])
					valuesList[i][0][8] = payoutRateFormat(hitOddsMap[types.TrioOddsRange8].InexactFloat64(), allCountMap[types.TrioOddsRange8])
				case 36:
					valuesList[i][0][0] = fmt.Sprintf("単%s", types.WinOddsRange7.String())
					valuesList[i][0][1] = types.TrioOddsRange1.String()
					valuesList[i][0][2] = types.TrioOddsRange2.String()
					valuesList[i][0][3] = types.TrioOddsRange3.String()
					valuesList[i][0][4] = types.TrioOddsRange4.String()
					valuesList[i][0][5] = types.TrioOddsRange5.String()
					valuesList[i][0][6] = types.TrioOddsRange6.String()
					valuesList[i][0][7] = types.TrioOddsRange7.String()
					valuesList[i][0][8] = types.TrioOddsRange8.String()
				case 37:
					oddsMap := raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange7]
					valuesList[i][0][1] = oddsMap[types.TrioOddsRange1]
					valuesList[i][0][2] = oddsMap[types.TrioOddsRange2]
					valuesList[i][0][3] = oddsMap[types.TrioOddsRange3]
					valuesList[i][0][4] = oddsMap[types.TrioOddsRange4]
					valuesList[i][0][5] = oddsMap[types.TrioOddsRange5]
					valuesList[i][0][6] = oddsMap[types.TrioOddsRange6]
					valuesList[i][0][7] = oddsMap[types.TrioOddsRange7]
					valuesList[i][0][8] = oddsMap[types.TrioOddsRange8]
				case 38:
					oddsMap := raceOddsRangeCountMap[pivotalMarker][types.WinOddsRange7]
					valuesList[i][0][1] = oddsMap[types.TrioOddsRange1]
					valuesList[i][0][2] = oddsMap[types.TrioOddsRange2]
					valuesList[i][0][3] = oddsMap[types.TrioOddsRange3]
					valuesList[i][0][4] = oddsMap[types.TrioOddsRange4]
					valuesList[i][0][5] = oddsMap[types.TrioOddsRange5]
					valuesList[i][0][6] = oddsMap[types.TrioOddsRange6]
					valuesList[i][0][7] = oddsMap[types.TrioOddsRange7]
					valuesList[i][0][8] = oddsMap[types.TrioOddsRange8]
				case 39:
					hitCountMap := raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange7]
					allCountMap := markerAllOddsRangeCountMap[pivotalMarker][types.WinOddsRange7]
					valuesList[i][0][1] = hitRateFormat(hitCountMap[types.TrioOddsRange1], allCountMap[types.TrioOddsRange1])
					valuesList[i][0][2] = hitRateFormat(hitCountMap[types.TrioOddsRange2], allCountMap[types.TrioOddsRange2])
					valuesList[i][0][3] = hitRateFormat(hitCountMap[types.TrioOddsRange3], allCountMap[types.TrioOddsRange3])
					valuesList[i][0][4] = hitRateFormat(hitCountMap[types.TrioOddsRange4], allCountMap[types.TrioOddsRange4])
					valuesList[i][0][5] = hitRateFormat(hitCountMap[types.TrioOddsRange5], allCountMap[types.TrioOddsRange5])
					valuesList[i][0][6] = hitRateFormat(hitCountMap[types.TrioOddsRange6], allCountMap[types.TrioOddsRange6])
					valuesList[i][0][7] = hitRateFormat(hitCountMap[types.TrioOddsRange7], allCountMap[types.TrioOddsRange7])
					valuesList[i][0][8] = hitRateFormat(hitCountMap[types.TrioOddsRange8], allCountMap[types.TrioOddsRange8])
				case 40:
					hitOddsMap := pivotalMarkerHitTotalOddsMap[pivotalMarker][types.WinOddsRange7]
					allCountMap := markerAllOddsRangeCountMap[pivotalMarker][types.WinOddsRange7]
					valuesList[i][0][1] = payoutRateFormat(hitOddsMap[types.TrioOddsRange1].InexactFloat64(), allCountMap[types.TrioOddsRange1])
					valuesList[i][0][2] = payoutRateFormat(hitOddsMap[types.TrioOddsRange2].InexactFloat64(), allCountMap[types.TrioOddsRange2])
					valuesList[i][0][3] = payoutRateFormat(hitOddsMap[types.TrioOddsRange3].InexactFloat64(), allCountMap[types.TrioOddsRange3])
					valuesList[i][0][4] = payoutRateFormat(hitOddsMap[types.TrioOddsRange4].InexactFloat64(), allCountMap[types.TrioOddsRange4])
					valuesList[i][0][5] = payoutRateFormat(hitOddsMap[types.TrioOddsRange5].InexactFloat64(), allCountMap[types.TrioOddsRange5])
					valuesList[i][0][6] = payoutRateFormat(hitOddsMap[types.TrioOddsRange6].InexactFloat64(), allCountMap[types.TrioOddsRange6])
					valuesList[i][0][7] = payoutRateFormat(hitOddsMap[types.TrioOddsRange7].InexactFloat64(), allCountMap[types.TrioOddsRange7])
					valuesList[i][0][8] = payoutRateFormat(hitOddsMap[types.TrioOddsRange8].InexactFloat64(), allCountMap[types.TrioOddsRange8])
				case 41:
					valuesList[i][0][0] = fmt.Sprintf("単%s", types.WinOddsRange8.String())
					valuesList[i][0][1] = types.TrioOddsRange1.String()
					valuesList[i][0][2] = types.TrioOddsRange2.String()
					valuesList[i][0][3] = types.TrioOddsRange3.String()
					valuesList[i][0][4] = types.TrioOddsRange4.String()
					valuesList[i][0][5] = types.TrioOddsRange5.String()
					valuesList[i][0][6] = types.TrioOddsRange6.String()
					valuesList[i][0][7] = types.TrioOddsRange7.String()
					valuesList[i][0][8] = types.TrioOddsRange8.String()
				case 42:
					oddsMap := raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange8]
					valuesList[i][0][1] = oddsMap[types.TrioOddsRange1]
					valuesList[i][0][2] = oddsMap[types.TrioOddsRange2]
					valuesList[i][0][3] = oddsMap[types.TrioOddsRange3]
					valuesList[i][0][4] = oddsMap[types.TrioOddsRange4]
					valuesList[i][0][5] = oddsMap[types.TrioOddsRange5]
					valuesList[i][0][6] = oddsMap[types.TrioOddsRange6]
					valuesList[i][0][7] = oddsMap[types.TrioOddsRange7]
					valuesList[i][0][8] = oddsMap[types.TrioOddsRange8]
				case 43:
					oddsMap := raceOddsRangeCountMap[pivotalMarker][types.WinOddsRange8]
					valuesList[i][0][1] = oddsMap[types.TrioOddsRange1]
					valuesList[i][0][2] = oddsMap[types.TrioOddsRange2]
					valuesList[i][0][3] = oddsMap[types.TrioOddsRange3]
					valuesList[i][0][4] = oddsMap[types.TrioOddsRange4]
					valuesList[i][0][5] = oddsMap[types.TrioOddsRange5]
					valuesList[i][0][6] = oddsMap[types.TrioOddsRange6]
					valuesList[i][0][7] = oddsMap[types.TrioOddsRange7]
					valuesList[i][0][8] = oddsMap[types.TrioOddsRange8]
				case 44:
					hitCountMap := raceHitOddsRangeCountMap[pivotalMarker][types.WinOddsRange8]
					allCountMap := markerAllOddsRangeCountMap[pivotalMarker][types.WinOddsRange8]
					valuesList[i][0][1] = hitRateFormat(hitCountMap[types.TrioOddsRange1], allCountMap[types.TrioOddsRange1])
					valuesList[i][0][2] = hitRateFormat(hitCountMap[types.TrioOddsRange2], allCountMap[types.TrioOddsRange2])
					valuesList[i][0][3] = hitRateFormat(hitCountMap[types.TrioOddsRange3], allCountMap[types.TrioOddsRange3])
					valuesList[i][0][4] = hitRateFormat(hitCountMap[types.TrioOddsRange4], allCountMap[types.TrioOddsRange4])
					valuesList[i][0][5] = hitRateFormat(hitCountMap[types.TrioOddsRange5], allCountMap[types.TrioOddsRange5])
					valuesList[i][0][6] = hitRateFormat(hitCountMap[types.TrioOddsRange6], allCountMap[types.TrioOddsRange6])
					valuesList[i][0][7] = hitRateFormat(hitCountMap[types.TrioOddsRange7], allCountMap[types.TrioOddsRange7])
					valuesList[i][0][8] = hitRateFormat(hitCountMap[types.TrioOddsRange8], allCountMap[types.TrioOddsRange8])
				case 45:
					hitOddsMap := pivotalMarkerHitTotalOddsMap[pivotalMarker][types.WinOddsRange8]
					allCountMap := markerAllOddsRangeCountMap[pivotalMarker][types.WinOddsRange8]
					valuesList[i][0][1] = payoutRateFormat(hitOddsMap[types.TrioOddsRange1].InexactFloat64(), allCountMap[types.TrioOddsRange1])
					valuesList[i][0][2] = payoutRateFormat(hitOddsMap[types.TrioOddsRange2].InexactFloat64(), allCountMap[types.TrioOddsRange2])
					valuesList[i][0][3] = payoutRateFormat(hitOddsMap[types.TrioOddsRange3].InexactFloat64(), allCountMap[types.TrioOddsRange3])
					valuesList[i][0][4] = payoutRateFormat(hitOddsMap[types.TrioOddsRange4].InexactFloat64(), allCountMap[types.TrioOddsRange4])
					valuesList[i][0][5] = payoutRateFormat(hitOddsMap[types.TrioOddsRange5].InexactFloat64(), allCountMap[types.TrioOddsRange5])
					valuesList[i][0][6] = payoutRateFormat(hitOddsMap[types.TrioOddsRange6].InexactFloat64(), allCountMap[types.TrioOddsRange6])
					valuesList[i][0][7] = payoutRateFormat(hitOddsMap[types.TrioOddsRange7].InexactFloat64(), allCountMap[types.TrioOddsRange7])
					valuesList[i][0][8] = payoutRateFormat(hitOddsMap[types.TrioOddsRange8].InexactFloat64(), allCountMap[types.TrioOddsRange8])
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
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	})
	for i := 0; i < 9; i++ {
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
		valuesList = append(valuesList, [][]interface{}{
			{
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
				"出現回数",
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
	markerCombinationAnalysisMap map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis,
	winRaceOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
	trioRaceOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
) (map[types.Marker]map[types.OddsRangeType]map[types.OddsRangeType]int, error) {
	raceOddsRangeCountMap := map[types.Marker]map[types.OddsRangeType]map[types.OddsRangeType]int{}
	pivotalMarkers := []types.Marker{
		types.Favorite,
		types.Rival,
		types.BrackTriangle,
		types.WhiteTriangle,
		types.Star,
		types.Check,
	}

	markerCalculablesMap := map[types.Marker][]*analysis_entity.Calculable{}
	for _, marker := range pivotalMarkers {
		for markerCombinationId, markerCombinationAnalysis := range markerCombinationAnalysisMap {
			if markerCombinationId.TicketType().OriginTicketType() != types.Trio {
				continue
			}
			for _, calculable := range markerCombinationAnalysis.Calculables() {
				markerCalculablesMap[marker] = append(markerCalculablesMap[marker], calculable)
			}
		}
	}

	for marker, calculables := range markerCalculablesMap {
		raceOddsRangeCountMap[marker] = map[types.OddsRangeType]map[types.OddsRangeType]int{}
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

				if _, ok := raceOddsRangeCountMap[marker][pivotalMarkerOddsRange]; !ok {
					raceOddsRangeCountMap[marker][pivotalMarkerOddsRange] = map[types.OddsRangeType]int{}
				}

				for _, trioOdds := range trioOddsList {
					odds = trioOdds.Odds().InexactFloat64()
					if odds >= 1.0 && odds <= 9.9 {
						raceOddsRangeCountMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange1]++
					} else if odds >= 10.0 && odds <= 19.9 {
						raceOddsRangeCountMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange2]++
					} else if odds >= 20.0 && odds <= 29.9 {
						raceOddsRangeCountMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange3]++
					} else if odds >= 30.0 && odds <= 49.9 {
						raceOddsRangeCountMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange4]++
					} else if odds >= 50.0 && odds <= 99.9 {
						raceOddsRangeCountMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange5]++
					} else if odds >= 100.0 && odds <= 299.9 {
						raceOddsRangeCountMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange6]++
					} else if odds >= 300.0 && odds <= 499.9 {
						raceOddsRangeCountMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange7]++
					} else if odds >= 500.0 {
						raceOddsRangeCountMap[marker][pivotalMarkerOddsRange][types.TrioOddsRange8]++
					}
				}
			}
		}
	}

	return raceOddsRangeCountMap, nil
}

func (s *spreadSheetTrioAnalysisRepository) getMarkerAllOddsRangeCountMap(
	ctx context.Context,
	markerCombinationAnalysisMap map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis,
	winRaceOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
	trioMarkerOddsMap map[types.RaceId][]*spreadsheet_entity.Odds,
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

				if _, ok := pivotalMarkerAllCountOddsRangeMap[marker][pivotalMarkerOddsRange]; !ok {
					pivotalMarkerAllCountOddsRangeMap[marker][pivotalMarkerOddsRange] = map[types.OddsRangeType]int{}
				}

				for _, trioOdds := range markerOdds {
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
							EndColumnIndex:   9,
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
							StartColumnIndex: 0,
							StartRowIndex:    int64(idx*rowGroupSize + 1),
							EndColumnIndex:   1,
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
								StartColumnIndex: 0,
								StartRowIndex:    int64(1 + (i * 5) + idx*rowGroupSize),
								EndColumnIndex:   1,
								EndRowIndex:      int64(2 + (i * 5) + idx*rowGroupSize),
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
								StartColumnIndex: 1,
								StartRowIndex:    int64(1 + (i * 5) + idx*rowGroupSize),
								EndColumnIndex:   9,
								EndRowIndex:      int64(2 + (i * 5) + idx*rowGroupSize),
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
								StartColumnIndex: 0,
								StartRowIndex:    int64(1 + (i * 5) + idx*rowGroupSize),
								EndColumnIndex:   9,
								EndRowIndex:      int64(2 + (i * 5) + idx*rowGroupSize),
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
								StartColumnIndex: 0,
								StartRowIndex:    int64(1 + (i * 5) + idx*rowGroupSize),
								EndColumnIndex:   9,
								EndRowIndex:      int64(2 + (i * 5) + idx*rowGroupSize),
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
								StartColumnIndex: 0,
								StartRowIndex:    int64(2 + (i * 5) + idx*rowGroupSize),
								EndColumnIndex:   1,
								EndRowIndex:      int64(6 + (i * 5) + idx*rowGroupSize),
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
