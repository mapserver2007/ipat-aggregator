package analysis_service

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

// // 危険レースチェックリスト：合計点数が高いほど危険
// // https://www.youtube.com/watch?v=zfbJv2arw28&t=425s
// var dangerCheckListItems = []string{
// 	"15頭立て以上の多頭数戦であること:2",
// 	"ダート1000~1200m戦であること:3", // かなり紛れる
// 	"今走の馬場が不良であること:4",       // 相当危険
// 	"今走の馬場が重であること:2",
// 	"前走から騎手乗り替わりであること:1", // 意外と乗り替わりでも成績は悪くない
// 	"前走3着以下であること:2",
// 	"前走4着以下であること:2", // 4着以下は更に成績が悪い
// 	"3ヶ月以上の休み明けであること:1",
// 	"前走から2F以上距離が変更されていること:2", // 1Fの場合は危険度が低い
// 	"昇級初戦であること:1",            // そこまで危険度は高くない
// }

type PlaceNegativeCheck interface {
	// 15頭立て以上の多頭数戦であること(2点)
	LargeFieldRace(ctx context.Context, input *PlaceNegativeCheckInput) *spreadsheet_entity.AnalysisPlaceCheckPoint
	// ダート1000~1200m戦であること(3点)
	DirtShortDistanceRace(ctx context.Context, input *PlaceNegativeCheckInput) *spreadsheet_entity.AnalysisPlaceCheckPoint
	// 今走の馬場が不良または重であること(2点or4点)
	SoftTrackCondition(ctx context.Context, input *PlaceNegativeCheckInput) *spreadsheet_entity.AnalysisPlaceCheckPoint
	// 前走から騎手乗り替わりであること(1点)
	ChangeJockey(ctx context.Context, input *PlaceNegativeCheckInput) *spreadsheet_entity.AnalysisPlaceCheckPoint
}

type placeNegativeCheckService struct{}

type PlaceNegativeCheckInput struct {
	Race *analysis_entity.Race
	// Horse *analysis_entity.Horse
	// Forecast *analysis_entity.RaceForecast
}

func NewPlaceNegativeCheck() PlaceNegativeCheck {
	return &placeNegativeCheckService{}
}

// 15頭立て以上の多頭数戦であること(2点)
func (p *placeNegativeCheckService) LargeFieldRace(
	ctx context.Context,
	input *PlaceNegativeCheckInput,
) *spreadsheet_entity.AnalysisPlaceCheckPoint {
	point := 0
	if input.Race.Entries() >= 15 {
		point = 2
	}

	return spreadsheet_entity.NewAnalysisPlaceCheckPoint(
		"15頭立て以上の多頭数戦であること",
		point,
	)
}

// ダート1000~1200m戦であること(3点)
func (p *placeNegativeCheckService) DirtShortDistanceRace(
	ctx context.Context,
	input *PlaceNegativeCheckInput,
) *spreadsheet_entity.AnalysisPlaceCheckPoint {
	point := 0
	if input.Race.CourseCategory() == types.Dirt && input.Race.Distance() >= 1000 && input.Race.Distance() <= 1200 {
		point = 3
	}

	return spreadsheet_entity.NewAnalysisPlaceCheckPoint(
		"ダート1000~1200m戦であること",
		point,
	)
}

// 今走の馬場が不良であること(4点)
func (p *placeNegativeCheckService) SoftTrackCondition(
	ctx context.Context,
	input *PlaceNegativeCheckInput,
) *spreadsheet_entity.AnalysisPlaceCheckPoint {
	point := 0
	if input.Race.TrackCondition() == types.Soft {
		point = 4
	} else if input.Race.TrackCondition() == types.Yielding {
		point = 2
	}

	return spreadsheet_entity.NewAnalysisPlaceCheckPoint(
		"今走の馬場が不良または重であること",
		point,
	)
}

// 前走から騎手乗り替わりであること(1点)
func (p *placeNegativeCheckService) ChangeJockey(
	ctx context.Context,
	input *PlaceNegativeCheckInput,
) *spreadsheet_entity.AnalysisPlaceCheckPoint {
	point := 0
	// 騎手の乗り替わりのロジックを追加
	// 例: if input.Horse.JockeyId() != input.Race.PreviousJockeyId() {
	//     point = 1
	// }

	return spreadsheet_entity.NewAnalysisPlaceCheckPoint(
		"前走から騎手乗り替わりであること",
		point,
	)
}
