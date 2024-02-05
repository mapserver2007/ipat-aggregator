package service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"sort"
)

type AnalysisService interface {
	AddAnalysisData(ctx context.Context, markerCombinationId types.MarkerCombinationId, race *data_cache_entity.Race, numerical *analysis_entity.Calculable) error
	GetAnalysisData() *analysis_entity.Layer1
	GetSearchFilters() []filter.Id
	GetHitMarkerCombinationIds(ctx context.Context, result *data_cache_entity.PayoutResult, marker *marker_csv_entity.Yamato) []types.MarkerCombinationId
	GetUnHitMarkerCombinationIds(ctx context.Context, result *data_cache_entity.PayoutResult, marker *marker_csv_entity.Yamato) []types.MarkerCombinationId
	CreateSpreadSheetAnalysisData(ctx context.Context, analysisData *analysis_entity.Layer1) *spreadsheet_entity.AnalysisData
	CreateAnalysisFilters(ctx context.Context, race *data_cache_entity.Race, raceResultByMarker *data_cache_entity.RaceResult) []filter.Id
}

type analysisService struct {
	analysisData  *analysis_entity.Layer1
	searchFilters []filter.Id
}

func NewAnalysisService() AnalysisService {
	analysisData := analysis_entity.Layer1{
		MarkerCombination: make(map[types.MarkerCombinationId]*analysis_entity.Layer2),
	}
	searchFilters := []filter.Id{
		filter.All,
		filter.TurfShortDistance,
		filter.TurfMiddleDistance,
		filter.TurfLongDistance,
		filter.DirtShortDistance,
		filter.DirtLongDistance,
		filter.TurfShortDistanceJockeyTop1,
		filter.TurfMiddleDistanceJockeyTop1,
		filter.TurfLongDistanceJockeyTop1,
		filter.DirtShortDistanceJockeyTop1,
		filter.DirtLongDistanceJockeyTop1,
		filter.TurfShortDistanceJockeyTop2,
		filter.TurfMiddleDistanceJockeyTop2,
		filter.TurfLongDistanceJockeyTop2,
		filter.DirtShortDistanceJockeyTop2,
		filter.DirtLongDistanceJockeyTop2,
		filter.TurfShortDistanceJockeyOther,
		filter.TurfMiddleDistanceJockeyOther,
		filter.TurfLongDistanceJockeyOther,
		filter.DirtShortDistanceJockeyOther,
		filter.DirtLongDistanceJockeyOther,
	}

	return &analysisService{
		analysisData:  &analysisData,
		searchFilters: searchFilters,
	}
}

func (p *analysisService) AddAnalysisData(
	ctx context.Context,
	markerCombinationId types.MarkerCombinationId,
	race *data_cache_entity.Race,
	calculable *analysis_entity.Calculable,
) error {
	layer1 := p.analysisData.MarkerCombination
	if _, ok := layer1[markerCombinationId]; !ok {
		layer1[markerCombinationId] = &analysis_entity.Layer2{
			RaceDate: make(map[types.RaceDate]*analysis_entity.Layer3),
		}
	}
	layer2 := layer1[markerCombinationId].RaceDate
	if _, ok := layer2[race.RaceDate()]; !ok {
		layer2[race.RaceDate()] = &analysis_entity.Layer3{
			RaceId: make(map[types.RaceId][]*analysis_entity.Calculable),
		}
	}
	layer2[race.RaceDate()].RaceId[race.RaceId()] = append(layer2[race.RaceDate()].RaceId[race.RaceId()], calculable)

	return nil
}

func (p *analysisService) GetAnalysisData() *analysis_entity.Layer1 {
	return p.analysisData
}

func (p *analysisService) GetSearchFilters() []filter.Id {
	return p.searchFilters
}

func (p *analysisService) GetHitMarkerCombinationIds(
	ctx context.Context,
	result *data_cache_entity.PayoutResult,
	marker *marker_csv_entity.Yamato,
) []types.MarkerCombinationId {
	var markerCombinationIds []types.MarkerCombinationId
	switch result.TicketType() {
	case types.Win:
		rawHorseNumber := result.Numbers()[0].List()[0]
		markerCombinationId, _ := types.NewMarkerCombinationId(19)
		switch rawHorseNumber {
		case marker.Favorite():
			markerCombinationId, _ = types.NewMarkerCombinationId(11)
		case marker.Rival():
			markerCombinationId, _ = types.NewMarkerCombinationId(12)
		case marker.BrackTriangle():
			markerCombinationId, _ = types.NewMarkerCombinationId(13)
		case marker.WhiteTriangle():
			markerCombinationId, _ = types.NewMarkerCombinationId(14)
		case marker.Star():
			markerCombinationId, _ = types.NewMarkerCombinationId(15)
		case marker.Check():
			markerCombinationId, _ = types.NewMarkerCombinationId(16)
		}
		markerCombinationIds = append(markerCombinationIds, markerCombinationId)
	case types.Place:
		for _, number := range result.Numbers() {
			rawHorseNumber := number.List()[0]
			markerCombinationId, _ := types.NewMarkerCombinationId(29)
			switch rawHorseNumber {
			case marker.Favorite():
				markerCombinationId, _ = types.NewMarkerCombinationId(21)
			case marker.Rival():
				markerCombinationId, _ = types.NewMarkerCombinationId(22)
			case marker.BrackTriangle():
				markerCombinationId, _ = types.NewMarkerCombinationId(23)
			case marker.WhiteTriangle():
				markerCombinationId, _ = types.NewMarkerCombinationId(24)
			case marker.Star():
				markerCombinationId, _ = types.NewMarkerCombinationId(25)
			case marker.Check():
				markerCombinationId, _ = types.NewMarkerCombinationId(26)
			}
			markerCombinationIds = append(markerCombinationIds, markerCombinationId)
		}
	case types.QuinellaPlace:
		for _, horseNumber := range result.Numbers() {
			// 馬番に対する印の昇順ソート
			rawHorseNumbers := []int{99, 99} // 初期値は無扱いの99
			// mapは順序保証効かないのでキーだけを取り出してスライスに保存
			markerMap := marker.MarkerMap()
			// 印の重い順で処理
			for _, k := range []int{1, 2, 3, 4, 5, 6} {
				horseNumberByMarker := markerMap[types.Marker(k)]
				for idx, rawHorseNumber := range horseNumber.List() {
					if horseNumberByMarker == rawHorseNumber {
						rawHorseNumbers[idx] = rawHorseNumber
					}
				}
			}
			// 無扱いの99必ず末尾にするためにソート
			sort.Ints(rawHorseNumbers)

			markerCombinationId, _ := types.NewMarkerCombinationId(399)
			switch rawHorseNumbers[0] {
			case marker.Favorite():
				switch rawHorseNumbers[1] {
				case marker.Rival():
					markerCombinationId, _ = types.NewMarkerCombinationId(312)
				case marker.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(313)
				case marker.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(314)
				case marker.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(315)
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(316)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(319)
				}
			case marker.Rival():
				switch rawHorseNumbers[1] {
				case marker.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(323)
				case marker.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(324)
				case marker.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(325)
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(326)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(329)
				}
			case marker.BrackTriangle():
				switch rawHorseNumbers[1] {
				case marker.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(334)
				case marker.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(335)
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(336)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(339)
				}
			case marker.WhiteTriangle():
				switch rawHorseNumbers[1] {
				case marker.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(345)
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(346)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(349)
				}
			case marker.Star():
				switch rawHorseNumbers[1] {
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(356)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(359)
				}
			case marker.Check():
				markerCombinationId, _ = types.NewMarkerCombinationId(369)
			}
			markerCombinationIds = append(markerCombinationIds, markerCombinationId)
		}
	case types.Quinella:
		for _, horseNumber := range result.Numbers() {
			// 馬番に対する印の昇順ソート
			rawHorseNumbers := []int{99, 99} // 初期値は無扱いの99
			// mapは順序保証効かないのでキーだけを取り出してスライスに保存
			markerMap := marker.MarkerMap()
			// 印の重い順で処理
			for _, k := range []int{1, 2, 3, 4, 5, 6} {
				horseNumberByMarker := markerMap[types.Marker(k)]
				for idx, rawHorseNumber := range horseNumber.List() {
					if horseNumberByMarker == rawHorseNumber {
						rawHorseNumbers[idx] = rawHorseNumber
					}
				}
			}
			// 無扱いの99必ず末尾にするためにソート
			sort.Ints(rawHorseNumbers)

			markerCombinationId, _ := types.NewMarkerCombinationId(499)
			switch rawHorseNumbers[0] {
			case marker.Favorite():
				switch rawHorseNumbers[1] {
				case marker.Rival():
					markerCombinationId, _ = types.NewMarkerCombinationId(412)
				case marker.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(413)
				case marker.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(414)
				case marker.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(415)
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(416)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(419)
				}
			case marker.Rival():
				switch rawHorseNumbers[1] {
				case marker.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(423)
				case marker.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(424)
				case marker.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(425)
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(426)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(429)
				}
			case marker.BrackTriangle():
				switch rawHorseNumbers[1] {
				case marker.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(434)
				case marker.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(435)
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(436)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(439)
				}
			case marker.WhiteTriangle():
				switch rawHorseNumbers[1] {
				case marker.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(445)
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(446)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(449)
				}
			case marker.Star():
				switch rawHorseNumbers[1] {
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(456)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(459)
				}
			case marker.Check():
				markerCombinationId, _ = types.NewMarkerCombinationId(469)
			}
			markerCombinationIds = append(markerCombinationIds, markerCombinationId)
		}
	case types.Exacta:
		for _, horseNumber := range result.Numbers() {
			rawHorseNumbers := []int{99, 99} // 初期値は無扱いの99
			markerMap := marker.MarkerMap()
			// 印の重い順で処理
			for _, k := range []int{1, 2, 3, 4, 5, 6} {
				horseNumberByMarker := markerMap[types.Marker(k)]
				for idx, rawHorseNumber := range horseNumber.List() {
					if horseNumberByMarker == rawHorseNumber {
						rawHorseNumbers[idx] = rawHorseNumber
					}
				}
			}

			markerCombinationId, _ := types.NewMarkerCombinationId(599)
			switch rawHorseNumbers[0] {
			case marker.Favorite():
				switch rawHorseNumbers[1] {
				case marker.Rival():
					markerCombinationId, _ = types.NewMarkerCombinationId(512)
				case marker.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(513)
				case marker.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(514)
				case marker.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(515)
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(516)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(519)
				}
			case marker.Rival():
				switch rawHorseNumbers[1] {
				case marker.Favorite():
					markerCombinationId, _ = types.NewMarkerCombinationId(521)
				case marker.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(523)
				case marker.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(524)
				case marker.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(525)
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(526)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(529)
				}
			case marker.BrackTriangle():
				switch rawHorseNumbers[1] {
				case marker.Favorite():
					markerCombinationId, _ = types.NewMarkerCombinationId(531)
				case marker.Rival():
					markerCombinationId, _ = types.NewMarkerCombinationId(532)
				case marker.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(534)
				case marker.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(535)
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(536)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(539)
				}
			case marker.WhiteTriangle():
				switch rawHorseNumbers[1] {
				case marker.Favorite():
					markerCombinationId, _ = types.NewMarkerCombinationId(541)
				case marker.Rival():
					markerCombinationId, _ = types.NewMarkerCombinationId(542)
				case marker.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(543)
				case marker.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(545)
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(546)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(549)
				}
			case marker.Star():
				switch rawHorseNumbers[1] {
				case marker.Favorite():
					markerCombinationId, _ = types.NewMarkerCombinationId(551)
				case marker.Rival():
					markerCombinationId, _ = types.NewMarkerCombinationId(552)
				case marker.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(553)
				case marker.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(554)
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(556)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(559)
				}
			case marker.Check():
				switch rawHorseNumbers[1] {
				case marker.Favorite():
					markerCombinationId, _ = types.NewMarkerCombinationId(561)
				case marker.Rival():
					markerCombinationId, _ = types.NewMarkerCombinationId(562)
				case marker.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(563)
				case marker.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(564)
				case marker.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(565)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(569)
				}
			default:
				switch rawHorseNumbers[1] {
				case marker.Favorite():
					markerCombinationId, _ = types.NewMarkerCombinationId(591)
				case marker.Rival():
					markerCombinationId, _ = types.NewMarkerCombinationId(592)
				case marker.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(593)
				case marker.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(594)
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(595)
				case marker.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(596)
				}
			}
			markerCombinationIds = append(markerCombinationIds, markerCombinationId)
		}
	case types.Trio:
		for _, horseNumber := range result.Numbers() {
			rawHorseNumbers := []int{99, 99, 99} // 初期値は無扱いの99
			markerMap := marker.MarkerMap()
			// 印の重い順で処理
			for _, k := range []int{1, 2, 3, 4, 5, 6} {
				horseNumberByMarker := markerMap[types.Marker(k)]
				for idx, rawHorseNumber := range horseNumber.List() {
					if horseNumberByMarker == rawHorseNumber {
						rawHorseNumbers[idx] = rawHorseNumber
					}
				}
			}

			// 無扱いの99必ず末尾にするためにソート
			sort.Ints(rawHorseNumbers)

			markerCombinationId, _ := types.NewMarkerCombinationId(6999)
			switch rawHorseNumbers[0] {
			case marker.Favorite():
				switch rawHorseNumbers[1] {
				case marker.Rival():
					switch rawHorseNumbers[2] {
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(6123)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(6124)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(6125)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6126)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6129)
					}
				case marker.BrackTriangle():
					switch rawHorseNumbers[2] {
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(6134)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(6135)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6136)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6139)
					}
				case marker.WhiteTriangle():
					switch rawHorseNumbers[2] {
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(6145)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6146)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6149)
					}
				case marker.Star():
					switch rawHorseNumbers[2] {
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6156)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6159)
					}
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(6169)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(6199)
				}
			case marker.Rival():
				switch rawHorseNumbers[1] {
				case marker.BrackTriangle():
					switch rawHorseNumbers[2] {
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(6234)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(6235)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6236)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6239)
					}
				case marker.WhiteTriangle():
					switch rawHorseNumbers[2] {
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(6245)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6246)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6249)
					}
				case marker.Star():
					switch rawHorseNumbers[2] {
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6256)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6259)
					}
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(6269)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(6299)
				}
			case marker.BrackTriangle():
				switch rawHorseNumbers[1] {
				case marker.WhiteTriangle():
					switch rawHorseNumbers[2] {
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(6345)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6346)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6349)
					}
				case marker.Star():
					switch rawHorseNumbers[2] {
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6356)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6359)
					}
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(6369)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(6399)
				}
			case marker.WhiteTriangle():
				switch rawHorseNumbers[1] {
				case marker.Star():
					switch rawHorseNumbers[2] {
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6456)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6459)
					}
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(6469)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(6499)
				}
			case marker.Star():
				switch rawHorseNumbers[1] {
				case marker.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(6569)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(6599)
				}
			case marker.Check():
				markerCombinationId, _ = types.NewMarkerCombinationId(6699)
			}
			markerCombinationIds = append(markerCombinationIds, markerCombinationId)
		}
	case types.Trifecta:
		for _, horseNumber := range result.Numbers() {
			rawHorseNumbers := []int{99, 99, 99} // 初期値は無扱いの99
			markerMap := marker.MarkerMap()
			// 印の重い順で処理
			for _, k := range []int{1, 2, 3, 4, 5, 6} {
				horseNumberByMarker := markerMap[types.Marker(k)]
				for idx, rawHorseNumber := range horseNumber.List() {
					if horseNumberByMarker == rawHorseNumber {
						rawHorseNumbers[idx] = rawHorseNumber
					}
				}
			}

			markerCombinationId, _ := types.NewMarkerCombinationId(7999)
			switch rawHorseNumbers[0] {
			case marker.Favorite():
				switch rawHorseNumbers[1] {
				case marker.Rival():
					switch rawHorseNumbers[2] {
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7123)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7124)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7125)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7126)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7129)
					}
				case marker.BrackTriangle():
					switch rawHorseNumbers[2] {
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7132)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7134)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7135)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7136)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7139)
					}
				case marker.WhiteTriangle():
					switch rawHorseNumbers[2] {
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7142)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7143)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7145)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7146)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7149)
					}
				case marker.Star():
					switch rawHorseNumbers[2] {
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7152)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7153)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7154)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7156)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7159)
					}
				case marker.Check():
					switch rawHorseNumbers[2] {
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7162)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7163)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7164)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7165)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7169)
					}
				default:
					switch rawHorseNumbers[2] {
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7192)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7193)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7194)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7195)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7196)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7199)
					}
				}
			case marker.Rival():
				switch rawHorseNumbers[1] {
				case marker.Favorite():
					switch rawHorseNumbers[2] {
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7213)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7214)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7215)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7216)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7219)
					}
				case marker.BrackTriangle():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7231)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7234)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7235)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7236)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7239)
					}
				case marker.WhiteTriangle():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7241)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7243)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7245)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7246)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7249)
					}
				case marker.Star():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7251)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7253)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7254)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7256)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7259)
					}
				case marker.Check():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7261)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7263)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7264)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7265)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7269)
					}
				default:
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7291)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7293)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7294)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7295)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7296)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7299)
					}
				}
			case marker.BrackTriangle():
				switch rawHorseNumbers[1] {
				case marker.Favorite():
					switch rawHorseNumbers[2] {
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7312)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7314)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7315)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7316)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7319)
					}
				case marker.Rival():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7321)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7324)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7325)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7326)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7329)
					}
				case marker.WhiteTriangle():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7341)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7342)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7345)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7346)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7349)
					}
				case marker.Star():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7351)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7352)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7354)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7356)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7359)
					}
				case marker.Check():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7361)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7362)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7364)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7365)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7369)
					}
				default:
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7391)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7392)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7394)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7395)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7396)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7399)
					}
				}
			case marker.WhiteTriangle():
				switch rawHorseNumbers[1] {
				case marker.Favorite():
					switch rawHorseNumbers[2] {
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7412)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7413)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7415)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7416)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7419)
					}
				case marker.Rival():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7421)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7423)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7425)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7426)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7429)
					}
				case marker.BrackTriangle():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7431)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7432)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7435)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7436)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7439)
					}
				case marker.Star():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7451)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7452)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7453)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7456)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7459)
					}
				case marker.Check():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7461)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7462)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7463)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7465)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7469)
					}
				default:
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7491)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7492)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7493)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7495)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7496)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7499)
					}
				}
			case marker.Star():
				switch rawHorseNumbers[1] {
				case marker.Favorite():
					switch rawHorseNumbers[2] {
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7512)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7513)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7514)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7516)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7519)
					}
				case marker.Rival():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7521)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7523)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7524)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7526)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7529)
					}
				case marker.BrackTriangle():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7531)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7532)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7534)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7536)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7539)
					}
				case marker.WhiteTriangle():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7541)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7542)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7543)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7546)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7549)
					}
				case marker.Check():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7561)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7562)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7563)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7564)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7569)
					}
				default:
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7591)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7592)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7593)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7594)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7596)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7599)
					}
				}
			case marker.Check():
				switch rawHorseNumbers[1] {
				case marker.Favorite():
					switch rawHorseNumbers[2] {
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7612)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7613)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7614)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7615)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7619)
					}
				case marker.Rival():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7621)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7623)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7624)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7625)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7629)
					}
				case marker.BrackTriangle():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7631)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7632)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7634)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7635)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7639)
					}
				case marker.WhiteTriangle():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7641)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7642)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7643)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7645)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7649)
					}
				case marker.Star():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7651)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7652)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7653)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7654)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7659)
					}
				default:
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7691)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7692)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7693)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7694)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7695)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7699)
					}
				}
			default:
				switch rawHorseNumbers[1] {
				case marker.Favorite():
					switch rawHorseNumbers[2] {
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7912)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7913)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7914)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7915)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7916)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7919)
					}
				case marker.Rival():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7921)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7923)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7924)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7925)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7926)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7929)
					}
				case marker.BrackTriangle():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7931)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7932)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7934)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7935)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7936)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7939)
					}
				case marker.WhiteTriangle():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7941)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7942)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7943)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7945)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7946)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7949)
					}
				case marker.Star():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7951)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7952)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7953)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7954)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7956)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7959)
					}
				case marker.Check():
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7961)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7962)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7963)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7964)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7965)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7969)
					}
				default:
					switch rawHorseNumbers[2] {
					case marker.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7991)
					case marker.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7992)
					case marker.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7993)
					case marker.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7994)
					case marker.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7995)
					case marker.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7996)
					}
				}
			}
			markerCombinationIds = append(markerCombinationIds, markerCombinationId)
		}
	}

	return markerCombinationIds
}

func (p *analysisService) GetUnHitMarkerCombinationIds(
	ctx context.Context,
	result *data_cache_entity.PayoutResult,
	marker *marker_csv_entity.Yamato,
) []types.MarkerCombinationId {
	var (
		unHitMarkerCombinationIds   []types.MarkerCombinationId
		unHitMarkerCombinationIdMap map[types.MarkerCombinationId]bool
	)
	switch result.TicketType() {
	case types.Win:
		rawHorseNumber := result.Numbers()[0].List()[0]
		unHitMarkerCombinationIdMap = map[types.MarkerCombinationId]bool{
			types.MarkerCombinationId(11): true,
			types.MarkerCombinationId(12): true,
			types.MarkerCombinationId(13): true,
			types.MarkerCombinationId(14): true,
			types.MarkerCombinationId(15): true,
			types.MarkerCombinationId(16): true,
			types.MarkerCombinationId(19): true,
		}
		switch rawHorseNumber {
		case marker.Favorite():
			unHitMarkerCombinationIdMap[types.MarkerCombinationId(11)] = false
		case marker.Rival():
			unHitMarkerCombinationIdMap[types.MarkerCombinationId(12)] = false
		case marker.BrackTriangle():
			unHitMarkerCombinationIdMap[types.MarkerCombinationId(13)] = false
		case marker.WhiteTriangle():
			unHitMarkerCombinationIdMap[types.MarkerCombinationId(14)] = false
		case marker.Star():
			unHitMarkerCombinationIdMap[types.MarkerCombinationId(15)] = false
		case marker.Check():
			unHitMarkerCombinationIdMap[types.MarkerCombinationId(16)] = false
		}
	case types.Place:
		var rawHorseNumbers []int
		for _, numbers := range result.Numbers() {
			rawHorseNumbers = append(rawHorseNumbers, numbers.List()[0])
		}
		unHitMarkerCombinationIdMap = map[types.MarkerCombinationId]bool{
			types.MarkerCombinationId(21): true,
			types.MarkerCombinationId(22): true,
			types.MarkerCombinationId(23): true,
			types.MarkerCombinationId(24): true,
			types.MarkerCombinationId(25): true,
			types.MarkerCombinationId(26): true,
			types.MarkerCombinationId(29): true,
		}
		for _, rawHorseNumber := range rawHorseNumbers {
			switch rawHorseNumber {
			case marker.Favorite():
				unHitMarkerCombinationIdMap[types.MarkerCombinationId(21)] = false
			case marker.Rival():
				unHitMarkerCombinationIdMap[types.MarkerCombinationId(22)] = false
			case marker.BrackTriangle():
				unHitMarkerCombinationIdMap[types.MarkerCombinationId(23)] = false
			case marker.WhiteTriangle():
				unHitMarkerCombinationIdMap[types.MarkerCombinationId(24)] = false
			case marker.Star():
				unHitMarkerCombinationIdMap[types.MarkerCombinationId(25)] = false
			case marker.Check():
				unHitMarkerCombinationIdMap[types.MarkerCombinationId(26)] = false
			}
		}
	}

	for markerCombinationId, unHit := range unHitMarkerCombinationIdMap {
		if unHit {
			unHitMarkerCombinationIds = append(unHitMarkerCombinationIds, markerCombinationId)
		}
	}

	return unHitMarkerCombinationIds
}

func (p *analysisService) CreateSpreadSheetAnalysisData(
	ctx context.Context,
	analysisData *analysis_entity.Layer1,
) *spreadsheet_entity.AnalysisData {
	// TODO entityにしたほうがいいかもしれない。データ構造が複雑になってきた
	hitDataMapByFilter := map[filter.Id]map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis{}
	unHitDataMapByFilter := map[filter.Id]map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis{}
	raceCountMapByFilter := map[filter.Id]map[types.MarkerCombinationId]map[types.OddsRangeType]int{}

	for _, f := range p.searchFilters {
		raceCountMapByFilter[f] = p.calcMarkerCombinationRaceCountByFilter(analysisData, f)
		hitData, unHitData := p.createMarkerCombinationDataByFilter(analysisData, f)
		hitDataMapByFilter[f] = hitData
		unHitDataMapByFilter[f] = unHitData
	}

	return spreadsheet_entity.NewAnalysisData(
		hitDataMapByFilter,
		unHitDataMapByFilter,
		raceCountMapByFilter,
		p.createAllMarkerCombinations(),
	)
}

func (p *analysisService) createMarkerCombinationDataByFilter(
	analysisData *analysis_entity.Layer1,
	searchFilter filter.Id,
) (map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis, map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis) {
	hitMarkerCombinationDataMap := map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis{}
	unHitMarkerCombinationDataMap := map[types.MarkerCombinationId]*spreadsheet_entity.MarkerCombinationAnalysis{}

	raceCountMap := p.calcMarkerCombinationRaceCountByFilter(analysisData, searchFilter)
	for markerCombinationId, data := range analysisData.MarkerCombination {
		for _, data2 := range data.RaceDate {
			for _, data3 := range data2.RaceId {
				for _, calculable := range data3 {
					orderNo := calculable.OrderNo()
					switch markerCombinationId.TicketType() {
					case types.Win:
						if orderNo == 1 {
							if _, ok := hitMarkerCombinationDataMap[markerCombinationId]; !ok {
								hitMarkerCombinationDataMap[markerCombinationId] = spreadsheet_entity.NewMarkerCombinationAnalysis(raceCountMap[markerCombinationId])
							}
							match := true
							for _, f := range calculable.Filters() {
								if f&searchFilter == 0 {
									match = false
									break
								}
							}
							if match {
								hitMarkerCombinationDataMap[markerCombinationId].AddCalculable(calculable)
							}
						} else {
							if _, ok := unHitMarkerCombinationDataMap[markerCombinationId]; !ok {
								unHitMarkerCombinationDataMap[markerCombinationId] = spreadsheet_entity.NewMarkerCombinationAnalysis(raceCountMap[markerCombinationId])
							}
							match := true
							for _, f := range calculable.Filters() {
								if f&searchFilter == 0 {
									match = false
									break
								}
							}
							if match {
								unHitMarkerCombinationDataMap[markerCombinationId].AddCalculable(calculable)
							}
						}
					case types.Place:
						if orderNo <= 3 {
							if _, ok := hitMarkerCombinationDataMap[markerCombinationId]; !ok {
								hitMarkerCombinationDataMap[markerCombinationId] = spreadsheet_entity.NewMarkerCombinationAnalysis(raceCountMap[markerCombinationId])
							}
							match := true
							for _, f := range calculable.Filters() {
								if f&searchFilter == 0 {
									match = false
									break
								}
							}
							if match {
								hitMarkerCombinationDataMap[markerCombinationId].AddCalculable(calculable)
							}
						} else {
							if _, ok := unHitMarkerCombinationDataMap[markerCombinationId]; !ok {
								unHitMarkerCombinationDataMap[markerCombinationId] = spreadsheet_entity.NewMarkerCombinationAnalysis(raceCountMap[markerCombinationId])
							}
							match := true
							for _, f := range calculable.Filters() {
								if f&searchFilter == 0 {
									match = false
									break
								}
							}
							if match {
								unHitMarkerCombinationDataMap[markerCombinationId].AddCalculable(calculable)
							}
						}
					}
				}
			}
		}
	}

	return hitMarkerCombinationDataMap, unHitMarkerCombinationDataMap
}

func (p *analysisService) calcMarkerCombinationRaceCountByFilter(
	analysisData *analysis_entity.Layer1,
	searchFilter filter.Id,
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
						if f&searchFilter == 0 {
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

func (p *analysisService) createAllMarkerCombinations() []types.MarkerCombinationId {
	var markerCombinationIds []types.MarkerCombinationId
	for _, rawTicketType := range []int{1, 2, 3, 4, 5, 6, 7} {
		switch rawTicketType {
		case 1, 2:
			for _, rawMakerId := range []int{1, 2, 3, 4, 5, 6, 9} {
				markerCombinationId, _ := types.NewMarkerCombinationId(rawTicketType*10 + rawMakerId)
				markerCombinationIds = append(markerCombinationIds, markerCombinationId)
			}
		case 3, 4:
			for _, rawMakerId := range []int{1, 2, 3, 4, 5, 6, 9} {
				for _, rawMakerId2 := range []int{1, 2, 3, 4, 5, 6, 9} {
					if rawMakerId >= rawMakerId2 {
						continue
					}
					markerCombinationId, _ := types.NewMarkerCombinationId(rawTicketType*100 + rawMakerId*10 + rawMakerId2)
					markerCombinationIds = append(markerCombinationIds, markerCombinationId)
				}
			}
			markerCombinationId, _ := types.NewMarkerCombinationId(rawTicketType*100 + 99)
			markerCombinationIds = append(markerCombinationIds, markerCombinationId)
		case 5:
			for _, rawMakerId := range []int{1, 2, 3, 4, 5, 6, 9} {
				for _, rawMakerId2 := range []int{1, 2, 3, 4, 5, 6, 9} {
					if rawMakerId == rawMakerId2 {
						continue
					}
					markerCombinationId, _ := types.NewMarkerCombinationId(rawTicketType*100 + rawMakerId*10 + rawMakerId2)
					markerCombinationIds = append(markerCombinationIds, markerCombinationId)
				}
			}
		case 6:
			for _, rawMakerId := range []int{1, 2, 3, 4, 5, 6, 9} {
				for _, rawMakerId2 := range []int{1, 2, 3, 4, 5, 6, 9} {
					if rawMakerId >= rawMakerId2 {
						continue
					}
					for _, rawMakerId3 := range []int{1, 2, 3, 4, 5, 6, 9} {
						if rawMakerId2 >= rawMakerId3 {
							continue
						}
						markerCombinationId, _ := types.NewMarkerCombinationId(rawTicketType*1000 + rawMakerId*100 + rawMakerId2*10 + rawMakerId3)
						markerCombinationIds = append(markerCombinationIds, markerCombinationId)
					}
				}
			}
			markerCombinationId, _ := types.NewMarkerCombinationId(6999)
			markerCombinationIds = append(markerCombinationIds, markerCombinationId)
		case 7:
			for _, rawMakerId := range []int{1, 2, 3, 4, 5, 6, 9} {
				for _, rawMakerId2 := range []int{1, 2, 3, 4, 5, 6, 9} {
					if rawMakerId == rawMakerId2 {
						continue
					}
					for _, rawMakerId3 := range []int{1, 2, 3, 4, 5, 6, 9} {
						if rawMakerId == rawMakerId3 || rawMakerId2 == rawMakerId3 {
							continue
						}
						markerCombinationId, _ := types.NewMarkerCombinationId(rawTicketType*1000 + rawMakerId*100 + rawMakerId2*10 + rawMakerId3)
						markerCombinationIds = append(markerCombinationIds, markerCombinationId)
					}
				}
			}
			markerCombinationId, _ := types.NewMarkerCombinationId(7999)
			markerCombinationIds = append(markerCombinationIds, markerCombinationId)
		}
	}

	return markerCombinationIds
}

func (p *analysisService) CreateAnalysisFilters(
	ctx context.Context,
	race *data_cache_entity.Race,
	raceResultByMarker *data_cache_entity.RaceResult,
) []filter.Id {
	var filterIds []filter.Id
	switch race.CourseCategory() {
	case types.Turf:
		filterIds = append(filterIds, filter.Turf)
	case types.Dirt:
		filterIds = append(filterIds, filter.Dirt)
	}
	if race.Distance() >= 1000 && race.Distance() <= 1600 {
		filterIds = append(filterIds, filter.ShortDistance)
	} else if race.Distance() >= 1601 && race.Distance() <= 2000 {
		filterIds = append(filterIds, filter.MiddleDistance)
	} else if race.Distance() >= 2001 {
		filterIds = append(filterIds, filter.LongDistance)
	}
	switch raceResultByMarker.JockeyId() {
	case 5339: // C.ルメール
		filterIds = append(filterIds, filter.JokeyTop1)
	case 1088: // 川田将雅
		filterIds = append(filterIds, filter.JokeyTop2)
	default:
		filterIds = append(filterIds, filter.JokeyOther)
	}

	return filterIds
}
