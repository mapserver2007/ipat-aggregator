package service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/predict_analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/predict_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"sort"
)

type PredictAnalysisService interface {
	AddAnalysisData(ctx context.Context, markerCombinationId types.MarkerCombinationId, race *data_cache_entity.Race, numerical *predict_analysis_entity.Calculable) error
	GetAnalysisData() *predict_analysis_entity.Layer1
	GetMarkerCombinationIds(ctx context.Context, result *data_cache_entity.PayoutResult, predict *predict_csv_entity.Yamato) []types.MarkerCombinationId
	CreateSpreadSheetAnalysisData(ctx context.Context, analysisData *predict_analysis_entity.Layer1) map[types.MarkerCombinationId]*spreadsheet_entity.MarkerAnalysis
}

type predictAnalysisService struct {
	analysisData *predict_analysis_entity.Layer1
}

func NewPredictAnalysisService() PredictAnalysisService {
	analysisData := predict_analysis_entity.Layer1{
		MarkerCombination: make(map[types.MarkerCombinationId]*predict_analysis_entity.Layer2),
	}

	return &predictAnalysisService{
		analysisData: &analysisData,
	}
}

func (p *predictAnalysisService) AddAnalysisData(
	ctx context.Context,
	markerCombinationId types.MarkerCombinationId,
	race *data_cache_entity.Race,
	numerical *predict_analysis_entity.Calculable,
) error {
	layer1 := p.analysisData.MarkerCombination
	if _, ok := layer1[markerCombinationId]; !ok {
		layer1[markerCombinationId] = &predict_analysis_entity.Layer2{
			RaceDate: make(map[types.RaceDate]*predict_analysis_entity.Layer3),
		}
	}
	layer2 := layer1[markerCombinationId].RaceDate
	if _, ok := layer2[race.RaceDate()]; !ok {
		layer2[race.RaceDate()] = &predict_analysis_entity.Layer3{
			RaceId: make(map[types.RaceId][]*predict_analysis_entity.Result),
		}
	}
	result := predict_analysis_entity.NewResult(numerical, p.createFilters(race))
	layer2[race.RaceDate()].RaceId[race.RaceId()] = append(layer2[race.RaceDate()].RaceId[race.RaceId()], result)

	return nil
}

func (p *predictAnalysisService) GetAnalysisData() *predict_analysis_entity.Layer1 {
	return p.analysisData
}

func (p *predictAnalysisService) GetMarkerCombinationIds(
	ctx context.Context,
	result *data_cache_entity.PayoutResult,
	predict *predict_csv_entity.Yamato,
) []types.MarkerCombinationId {
	var markerCombinationIds []types.MarkerCombinationId
	switch result.TicketType() {
	case types.Win:
		rawHorseNumber := result.Numbers()[0].List()[0]
		markerCombinationId, _ := types.NewMarkerCombinationId(19)
		switch rawHorseNumber {
		case predict.Favorite():
			markerCombinationId, _ = types.NewMarkerCombinationId(11)
		case predict.Rival():
			markerCombinationId, _ = types.NewMarkerCombinationId(12)
		case predict.BrackTriangle():
			markerCombinationId, _ = types.NewMarkerCombinationId(13)
		case predict.WhiteTriangle():
			markerCombinationId, _ = types.NewMarkerCombinationId(14)
		case predict.Star():
			markerCombinationId, _ = types.NewMarkerCombinationId(15)
		case predict.Check():
			markerCombinationId, _ = types.NewMarkerCombinationId(16)
		}
		markerCombinationIds = append(markerCombinationIds, markerCombinationId)
	case types.Place:
		rawHorseNumber := result.Numbers()[0].List()[0]
		markerCombinationId, _ := types.NewMarkerCombinationId(29)
		switch rawHorseNumber {
		case predict.Favorite():
			markerCombinationId, _ = types.NewMarkerCombinationId(21)
		case predict.Rival():
			markerCombinationId, _ = types.NewMarkerCombinationId(22)
		case predict.BrackTriangle():
			markerCombinationId, _ = types.NewMarkerCombinationId(23)
		case predict.WhiteTriangle():
			markerCombinationId, _ = types.NewMarkerCombinationId(24)
		case predict.Star():
			markerCombinationId, _ = types.NewMarkerCombinationId(25)
		case predict.Check():
			markerCombinationId, _ = types.NewMarkerCombinationId(26)
		}
		markerCombinationIds = append(markerCombinationIds, markerCombinationId)
	case types.QuinellaPlace:
		for _, horseNumber := range result.Numbers() {
			// 馬番に対する印の昇順ソート
			rawHorseNumbers := []int{99, 99} // 初期値は無扱いの99
			// mapは順序保証効かないのでキーだけを取り出してスライスに保存
			markerMap := predict.MarkerMap()
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
			case predict.Favorite():
				switch rawHorseNumbers[1] {
				case predict.Rival():
					markerCombinationId, _ = types.NewMarkerCombinationId(312)
				case predict.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(313)
				case predict.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(314)
				case predict.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(315)
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(316)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(319)
				}
			case predict.Rival():
				switch rawHorseNumbers[1] {
				case predict.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(323)
				case predict.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(324)
				case predict.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(325)
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(326)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(329)
				}
			case predict.BrackTriangle():
				switch rawHorseNumbers[1] {
				case predict.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(334)
				case predict.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(335)
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(336)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(339)
				}
			case predict.WhiteTriangle():
				switch rawHorseNumbers[1] {
				case predict.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(345)
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(346)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(349)
				}
			case predict.Star():
				switch rawHorseNumbers[1] {
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(356)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(359)
				}
			case predict.Check():
				markerCombinationId, _ = types.NewMarkerCombinationId(369)
			}
			markerCombinationIds = append(markerCombinationIds, markerCombinationId)
		}
	case types.Quinella:
		for _, horseNumber := range result.Numbers() {
			// 馬番に対する印の昇順ソート
			rawHorseNumbers := []int{99, 99} // 初期値は無扱いの99
			// mapは順序保証効かないのでキーだけを取り出してスライスに保存
			markerMap := predict.MarkerMap()
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
			case predict.Favorite():
				switch rawHorseNumbers[1] {
				case predict.Rival():
					markerCombinationId, _ = types.NewMarkerCombinationId(412)
				case predict.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(413)
				case predict.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(414)
				case predict.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(415)
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(416)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(419)
				}
			case predict.Rival():
				switch rawHorseNumbers[1] {
				case predict.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(423)
				case predict.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(424)
				case predict.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(425)
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(426)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(429)
				}
			case predict.BrackTriangle():
				switch rawHorseNumbers[1] {
				case predict.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(434)
				case predict.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(435)
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(436)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(439)
				}
			case predict.WhiteTriangle():
				switch rawHorseNumbers[1] {
				case predict.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(445)
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(446)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(449)
				}
			case predict.Star():
				switch rawHorseNumbers[1] {
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(456)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(459)
				}
			case predict.Check():
				markerCombinationId, _ = types.NewMarkerCombinationId(469)
			}
			markerCombinationIds = append(markerCombinationIds, markerCombinationId)
		}
	case types.Exacta:
		for _, horseNumber := range result.Numbers() {
			rawHorseNumbers := []int{99, 99} // 初期値は無扱いの99
			markerMap := predict.MarkerMap()
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
			case predict.Favorite():
				switch rawHorseNumbers[1] {
				case predict.Rival():
					markerCombinationId, _ = types.NewMarkerCombinationId(512)
				case predict.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(513)
				case predict.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(514)
				case predict.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(515)
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(516)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(519)
				}
			case predict.Rival():
				switch rawHorseNumbers[1] {
				case predict.Favorite():
					markerCombinationId, _ = types.NewMarkerCombinationId(521)
				case predict.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(523)
				case predict.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(524)
				case predict.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(525)
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(526)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(529)
				}
			case predict.BrackTriangle():
				switch rawHorseNumbers[1] {
				case predict.Favorite():
					markerCombinationId, _ = types.NewMarkerCombinationId(531)
				case predict.Rival():
					markerCombinationId, _ = types.NewMarkerCombinationId(532)
				case predict.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(534)
				case predict.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(535)
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(536)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(539)
				}
			case predict.WhiteTriangle():
				switch rawHorseNumbers[1] {
				case predict.Favorite():
					markerCombinationId, _ = types.NewMarkerCombinationId(541)
				case predict.Rival():
					markerCombinationId, _ = types.NewMarkerCombinationId(542)
				case predict.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(543)
				case predict.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(545)
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(546)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(549)
				}
			case predict.Star():
				switch rawHorseNumbers[1] {
				case predict.Favorite():
					markerCombinationId, _ = types.NewMarkerCombinationId(551)
				case predict.Rival():
					markerCombinationId, _ = types.NewMarkerCombinationId(552)
				case predict.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(553)
				case predict.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(554)
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(556)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(559)
				}
			case predict.Check():
				switch rawHorseNumbers[1] {
				case predict.Favorite():
					markerCombinationId, _ = types.NewMarkerCombinationId(561)
				case predict.Rival():
					markerCombinationId, _ = types.NewMarkerCombinationId(562)
				case predict.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(563)
				case predict.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(564)
				case predict.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(565)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(569)
				}
			default:
				switch rawHorseNumbers[1] {
				case predict.Favorite():
					markerCombinationId, _ = types.NewMarkerCombinationId(591)
				case predict.Rival():
					markerCombinationId, _ = types.NewMarkerCombinationId(592)
				case predict.BrackTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(593)
				case predict.WhiteTriangle():
					markerCombinationId, _ = types.NewMarkerCombinationId(594)
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(595)
				case predict.Star():
					markerCombinationId, _ = types.NewMarkerCombinationId(596)
				}
			}
			markerCombinationIds = append(markerCombinationIds, markerCombinationId)
		}
	case types.Trio:
		for _, horseNumber := range result.Numbers() {
			rawHorseNumbers := []int{99, 99, 99} // 初期値は無扱いの99
			markerMap := predict.MarkerMap()
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
			case predict.Favorite():
				switch rawHorseNumbers[1] {
				case predict.Rival():
					switch rawHorseNumbers[2] {
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(6123)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(6124)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(6125)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6126)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6129)
					}
				case predict.BrackTriangle():
					switch rawHorseNumbers[2] {
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(6134)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(6135)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6136)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6139)
					}
				case predict.WhiteTriangle():
					switch rawHorseNumbers[2] {
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(6145)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6146)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6149)
					}
				case predict.Star():
					switch rawHorseNumbers[2] {
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6156)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6159)
					}
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(6169)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(6199)
				}
			case predict.Rival():
				switch rawHorseNumbers[1] {
				case predict.BrackTriangle():
					switch rawHorseNumbers[2] {
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(6234)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(6235)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6236)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6239)
					}
				case predict.WhiteTriangle():
					switch rawHorseNumbers[2] {
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(6245)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6246)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6249)
					}
				case predict.Star():
					switch rawHorseNumbers[2] {
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6256)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6259)
					}
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(6269)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(6299)
				}
			case predict.BrackTriangle():
				switch rawHorseNumbers[1] {
				case predict.WhiteTriangle():
					switch rawHorseNumbers[2] {
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(6345)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6346)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6349)
					}
				case predict.Star():
					switch rawHorseNumbers[2] {
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6356)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6359)
					}
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(6369)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(6399)
				}
			case predict.WhiteTriangle():
				switch rawHorseNumbers[1] {
				case predict.Star():
					switch rawHorseNumbers[2] {
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(6456)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(6459)
					}
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(6469)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(6499)
				}
			case predict.Star():
				switch rawHorseNumbers[1] {
				case predict.Check():
					markerCombinationId, _ = types.NewMarkerCombinationId(6569)
				default:
					markerCombinationId, _ = types.NewMarkerCombinationId(6599)
				}
			case predict.Check():
				markerCombinationId, _ = types.NewMarkerCombinationId(6699)
			}
			markerCombinationIds = append(markerCombinationIds, markerCombinationId)
		}
	case types.Trifecta:
		for _, horseNumber := range result.Numbers() {
			rawHorseNumbers := []int{99, 99, 99} // 初期値は無扱いの99
			markerMap := predict.MarkerMap()
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
			case predict.Favorite():
				switch rawHorseNumbers[1] {
				case predict.Rival():
					switch rawHorseNumbers[2] {
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7123)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7124)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7125)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7126)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7129)
					}
				case predict.BrackTriangle():
					switch rawHorseNumbers[2] {
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7132)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7134)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7135)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7136)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7139)
					}
				case predict.WhiteTriangle():
					switch rawHorseNumbers[2] {
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7142)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7143)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7145)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7146)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7149)
					}
				case predict.Star():
					switch rawHorseNumbers[2] {
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7152)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7153)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7154)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7156)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7159)
					}
				case predict.Check():
					switch rawHorseNumbers[2] {
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7162)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7163)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7164)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7165)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7169)
					}
				default:
					switch rawHorseNumbers[2] {
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7192)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7193)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7194)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7195)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7196)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7199)
					}
				}
			case predict.Rival():
				switch rawHorseNumbers[1] {
				case predict.Favorite():
					switch rawHorseNumbers[2] {
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7213)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7214)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7215)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7216)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7219)
					}
				case predict.BrackTriangle():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7231)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7234)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7235)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7236)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7239)
					}
				case predict.WhiteTriangle():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7241)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7243)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7245)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7246)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7249)
					}
				case predict.Star():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7251)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7253)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7254)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7256)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7259)
					}
				case predict.Check():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7261)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7263)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7264)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7265)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7269)
					}
				default:
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7291)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7293)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7294)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7295)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7296)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7299)
					}
				}
			case predict.BrackTriangle():
				switch rawHorseNumbers[1] {
				case predict.Favorite():
					switch rawHorseNumbers[2] {
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7312)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7314)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7315)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7316)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7319)
					}
				case predict.Rival():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7321)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7324)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7325)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7326)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7329)
					}
				case predict.WhiteTriangle():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7341)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7342)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7345)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7346)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7349)
					}
				case predict.Star():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7351)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7352)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7354)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7356)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7359)
					}
				case predict.Check():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7361)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7362)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7364)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7365)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7369)
					}
				default:
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7391)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7392)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7394)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7395)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7396)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7399)
					}
				}
			case predict.WhiteTriangle():
				switch rawHorseNumbers[1] {
				case predict.Favorite():
					switch rawHorseNumbers[2] {
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7412)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7413)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7415)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7416)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7419)
					}
				case predict.Rival():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7421)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7423)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7425)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7426)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7429)
					}
				case predict.BrackTriangle():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7431)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7432)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7435)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7436)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7439)
					}
				case predict.Star():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7451)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7452)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7453)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7456)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7459)
					}
				case predict.Check():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7461)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7462)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7463)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7465)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7469)
					}
				default:
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7491)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7492)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7493)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7495)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7496)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7499)
					}
				}
			case predict.Star():
				switch rawHorseNumbers[1] {
				case predict.Favorite():
					switch rawHorseNumbers[2] {
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7512)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7513)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7514)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7516)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7519)
					}
				case predict.Rival():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7521)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7523)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7524)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7526)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7529)
					}
				case predict.BrackTriangle():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7531)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7532)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7534)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7536)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7539)
					}
				case predict.WhiteTriangle():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7541)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7542)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7543)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7546)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7549)
					}
				case predict.Check():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7561)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7562)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7563)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7564)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7569)
					}
				default:
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7591)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7592)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7593)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7594)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7596)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7599)
					}
				}
			case predict.Check():
				switch rawHorseNumbers[1] {
				case predict.Favorite():
					switch rawHorseNumbers[2] {
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7612)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7613)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7614)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7615)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7619)
					}
				case predict.Rival():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7621)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7623)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7624)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7625)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7629)
					}
				case predict.BrackTriangle():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7631)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7632)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7634)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7635)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7639)
					}
				case predict.WhiteTriangle():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7641)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7642)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7643)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7645)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7649)
					}
				case predict.Star():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7651)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7652)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7653)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7654)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7659)
					}
				default:
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7691)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7692)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7693)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7694)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7695)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7699)
					}
				}
			default:
				switch rawHorseNumbers[1] {
				case predict.Favorite():
					switch rawHorseNumbers[2] {
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7912)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7913)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7914)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7915)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7916)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7919)
					}
				case predict.Rival():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7921)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7923)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7924)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7925)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7926)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7929)
					}
				case predict.BrackTriangle():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7931)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7932)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7934)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7935)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7936)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7939)
					}
				case predict.WhiteTriangle():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7941)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7942)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7943)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7945)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7946)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7949)
					}
				case predict.Star():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7951)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7952)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7953)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7954)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7956)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7959)
					}
				case predict.Check():
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7961)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7962)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7963)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7964)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7965)
					default:
						markerCombinationId, _ = types.NewMarkerCombinationId(7969)
					}
				default:
					switch rawHorseNumbers[2] {
					case predict.Favorite():
						markerCombinationId, _ = types.NewMarkerCombinationId(7991)
					case predict.Rival():
						markerCombinationId, _ = types.NewMarkerCombinationId(7992)
					case predict.BrackTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7993)
					case predict.WhiteTriangle():
						markerCombinationId, _ = types.NewMarkerCombinationId(7994)
					case predict.Star():
						markerCombinationId, _ = types.NewMarkerCombinationId(7995)
					case predict.Check():
						markerCombinationId, _ = types.NewMarkerCombinationId(7996)
					}
				}
			}
			markerCombinationIds = append(markerCombinationIds, markerCombinationId)
		}
	}

	return markerCombinationIds
}

func (p *predictAnalysisService) CreateSpreadSheetAnalysisData(
	ctx context.Context,
	analysisData *predict_analysis_entity.Layer1,
) map[types.MarkerCombinationId]*spreadsheet_entity.MarkerAnalysis {
	markerCombinationMap := make(map[types.MarkerCombinationId]*spreadsheet_entity.MarkerAnalysis)

	// 印単位での集計をかける
	for markerCombinationId, data := range analysisData.MarkerCombination {
		if _, ok := markerCombinationMap[markerCombinationId]; !ok {
			markerCombinationMap[markerCombinationId] = spreadsheet_entity.NewMarkerAnalysis()
		}
		for _, data2 := range data.RaceDate {
			for _, data3 := range data2.RaceId {
				for _, data4 := range data3 {
					markerCombinationMap[markerCombinationId].AddRaceCount()
					markerCombinationMap[markerCombinationId].AddPopular(data4.Calculable().Popular())
					markerCombinationMap[markerCombinationId].AddOdds(data4.Calculable().Odds())

					if data4.Calculable().Payment() > 0 {
						markerCombinationMap[markerCombinationId].AddVoteCount()
						markerCombinationMap[markerCombinationId].AddPayment(data4.Calculable().Payment().Value())
					}
					if data4.Calculable().Payout() > 0 {
						markerCombinationMap[markerCombinationId].AddHitCount()
						markerCombinationMap[markerCombinationId].AddPayout(data4.Calculable().Payout().Value())
					}
				}
			}
		}
	}

	return markerCombinationMap
}

func (p *predictAnalysisService) createFilters(
	race *data_cache_entity.Race,
) []predict_analysis_entity.IFilter {
	var filters []predict_analysis_entity.IFilter
	filters = append(filters, predict_analysis_entity.NewFilter(
		types.GradeClassFilter,
		race.Class(),
	))
	filters = append(filters, predict_analysis_entity.NewFilter(
		types.CourseCategoryFilter,
		race.CourseCategory(),
	))

	return filters
}
