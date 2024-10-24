package prediction_service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type PlaceCheckList interface {
	OkEntries(ctx context.Context, input *PlaceCheckListInput) bool
	OkWinOdds(ctx context.Context, input *PlaceCheckListInput) bool
	OkInThirdPlaceRatio(ctx context.Context, input *PlaceCheckListInput) bool
	OkNotChangeCourseCategory(ctx context.Context, input *PlaceCheckListInput) bool
	OkSameDistance(ctx context.Context, input *PlaceCheckListInput) bool
	OkSameCourseCategory(ctx context.Context, input *PlaceCheckListInput) bool
	OkInThirdPlaceRecent(ctx context.Context, input *PlaceCheckListInput) bool
	OkTrackConditionExperience(ctx context.Context, input *PlaceCheckListInput) bool
	OkNotHorseWeightUp(ctx context.Context, input *PlaceCheckListInput) bool
	OkNotClassUp(ctx context.Context, input *PlaceCheckListInput) bool
	OkContinueOrEnhancementJockey(ctx context.Context, input *PlaceCheckListInput) bool
	OkNotSlowStart(ctx context.Context, input *PlaceCheckListInput) bool
	OkFavoriteRatio(ctx context.Context, input *PlaceCheckListInput) bool
	OkOnlyFavoriteAndRival(ctx context.Context, input *PlaceCheckListInput) bool
	OkIsHighlyRecommended(ctx context.Context, input *PlaceCheckListInput) bool
}

type placeCheckListService struct{}

type PlaceCheckListInput struct {
	Race     *prediction_entity.Race
	Horse    *prediction_entity.Horse
	Forecast *prediction_entity.RaceForecast
}

func NewPlaceCheckList() PlaceCheckList {
	return &placeCheckListService{}
}

// OkEntries 13頭立て以下であること
func (p *placeCheckListService) OkEntries(ctx context.Context, input *PlaceCheckListInput) bool {
	return input.Race.Entries() <= 13
}

// OkWinOdds 単勝1倍台であること
func (p *placeCheckListService) OkWinOdds(ctx context.Context, input *PlaceCheckListInput) bool {
	for _, entryHorse := range input.Race.RaceEntryHorses() {
		entryHorse.HorseNumber()
		if entryHorse.HorseId() == input.Horse.HorseId() {
			for _, o := range input.Race.Odds() {
				if o.HorseNumber() == entryHorse.HorseNumber() {
					if o.Odds().InexactFloat64() < 2.0 {
						return true
					}
				}
			}
		}
	}

	return false
}

// OkInThirdPlaceRatio 3着以内率80%であること
func (p *placeCheckListService) OkInThirdPlaceRatio(ctx context.Context, input *PlaceCheckListInput) bool {
	isOk := false
	for _, entryHorse := range input.Race.RaceEntryHorses() {
		if entryHorse.HorseId() == input.Horse.HorseId() {
			var placed float64
			raceNum := float64(len(input.Horse.HorseResults()))
			for _, horseResult := range input.Horse.HorseResults() {
				// 履歴を取るタイミングで今走が履歴に含まれる場合はスキップする
				if horseResult.RaceId() == input.Race.RaceId() {
					continue
				}
				if horseResult.OrderNo() <= 3 {
					placed++
				}
			}
			isOk = placed/raceNum >= 0.8
		}
	}

	return isOk
}

// OkNotChangeCourseCategory 芝ダート替わりでないこと
func (p *placeCheckListService) OkNotChangeCourseCategory(ctx context.Context, input *PlaceCheckListInput) bool {
	for _, entryHorse := range input.Race.RaceEntryHorses() {
		if entryHorse.HorseId() == input.Horse.HorseId() {
			for _, horseResult := range input.Horse.HorseResults() {
				// 履歴を取るタイミングで今走が履歴に含まれる場合はスキップする
				if horseResult.RaceId() == input.Race.RaceId() {
					continue
				}
				return input.Race.CourseCategory() == horseResult.CourseCategory()
			}
		}
	}

	return false
}

// OkSameDistance 前走または2走前と今走の距離が同じなこと
func (p *placeCheckListService) OkSameDistance(ctx context.Context, input *PlaceCheckListInput) bool {
	for _, entryHorse := range input.Race.RaceEntryHorses() {
		if entryHorse.HorseId() == input.Horse.HorseId() {
			historyCount := 0
			for _, horseResult := range input.Horse.HorseResults() {
				// 履歴を取るタイミングで今走が履歴に含まれる場合はスキップする
				if horseResult.RaceId() == input.Race.RaceId() {
					continue
				}
				if historyCount > 2 { // 2走前まで
					break
				}
				if input.Race.Distance() == horseResult.Distance() {
					return true
				}
				historyCount++
			}
		}
	}

	return false
}

// OkSameCourseCategory 前走または2走前と今走のコースが同じなこと
func (p *placeCheckListService) OkSameCourseCategory(ctx context.Context, input *PlaceCheckListInput) bool {
	for _, entryHorse := range input.Race.RaceEntryHorses() {
		if entryHorse.HorseId() == input.Horse.HorseId() {
			historyCount := 0
			for _, horseResult := range input.Horse.HorseResults() {
				// 履歴を取るタイミングで今走が履歴に含まれる場合はスキップする
				if horseResult.RaceId() == input.Race.RaceId() {
					continue
				}
				if historyCount > 2 { // 2走前まで
					break
				}
				if input.Race.RaceCourse() == horseResult.RaceCourse() {
					return true
				}
				historyCount++
			}
		}
	}

	return false
}

// OkInThirdPlaceRecent 前走または2走前に馬券内なこと
func (p *placeCheckListService) OkInThirdPlaceRecent(ctx context.Context, input *PlaceCheckListInput) bool {
	for _, entryHorse := range input.Race.RaceEntryHorses() {
		if entryHorse.HorseId() == input.Horse.HorseId() {
			historyCount := 0
			for _, horseResult := range input.Horse.HorseResults() {
				// 履歴を取るタイミングで今走が履歴に含まれる場合はスキップする
				if horseResult.RaceId() == input.Race.RaceId() {
					continue
				}
				if historyCount > 2 { // 2走前まで
					break
				}
				if horseResult.OrderNo() <= 3 {
					return true
				}
				historyCount++
			}
		}
	}

	return false
}

// OkTrackConditionExperience 今走の馬場状態と同じ馬場状態で馬券内経験があること
func (p *placeCheckListService) OkTrackConditionExperience(ctx context.Context, input *PlaceCheckListInput) bool {
	for _, entryHorse := range input.Race.RaceEntryHorses() {
		if entryHorse.HorseId() == input.Horse.HorseId() {
			for _, horseResult := range input.Horse.HorseResults() {
				// 履歴を取るタイミングで今走が履歴に含まれる場合はスキップする
				if horseResult.RaceId() == input.Race.RaceId() {
					continue
				}
				if horseResult.TrackCondition() == input.Race.TrackCondition() && horseResult.OrderNo() <= 3 {
					return true
				}
			}
		}
	}

	return false
}

// OkNotHorseWeightUp 斤量増でないこと
func (p *placeCheckListService) OkNotHorseWeightUp(ctx context.Context, input *PlaceCheckListInput) bool {
	for _, entryHorse := range input.Race.RaceEntryHorses() {
		if entryHorse.HorseId() == input.Horse.HorseId() {
			if len(input.Horse.HorseResults()) == 0 {
				return true
			}
			for _, horseResult := range input.Horse.HorseResults() {
				// 履歴を取るタイミングで今走が履歴に含まれる場合はスキップする
				if horseResult.RaceId() == input.Race.RaceId() {
					continue
				}
				return horseResult.RaceWeight() >= entryHorse.RaceWeight()
			}
		}
	}

	return false
}

// OkNotClassUp 昇級初戦でないこと
func (p *placeCheckListService) OkNotClassUp(ctx context.Context, input *PlaceCheckListInput) bool {
	classMap := map[types.GradeClass]int{
		types.NonGrade:      1,
		types.OneWinClass:   2,
		types.TwoWinClass:   3,
		types.ThreeWinClass: 4,
		types.Grade1:        5,
		types.Grade2:        5,
		types.Grade3:        5,
		types.OpenClass:     5,
		types.ListedClass:   5,
	}

	for _, entryHorse := range input.Race.RaceEntryHorses() {
		if entryHorse.HorseId() == input.Horse.HorseId() {
			if len(input.Horse.HorseResults()) == 0 {
				return true
			}
			for _, horseResult := range input.Horse.HorseResults() {
				// 履歴を取るタイミングで今走が履歴に含まれる場合はスキップする
				if horseResult.RaceId() == input.Race.RaceId() {
					continue
				}
				if classMap[input.Race.Class()] > classMap[horseResult.Class()] {
					return false
				}
				return true
			}
		}
	}
	return false
}

// OkContinueOrEnhancementJockey 継続騎乗もしくは鞍上強化であること
func (p *placeCheckListService) OkContinueOrEnhancementJockey(ctx context.Context, input *PlaceCheckListInput) bool {
	containsInSlice := func(slice []types.JockeyId, jockeyId types.JockeyId) bool {
		for _, c := range slice {
			if c == jockeyId {
				return true
			}
		}
		return false
	}

	for _, entryHorse := range input.Race.RaceEntryHorses() {
		if entryHorse.HorseId() == input.Horse.HorseId() {
			if len(input.Horse.HorseResults()) == 0 {
				return true
			}
			for _, horseResult := range input.Horse.HorseResults() {
				// 履歴を取るタイミングで今走が履歴に含まれる場合はスキップする
				if horseResult.RaceId() == input.Race.RaceId() {
					continue
				}
				// 継続騎乗
				if horseResult.JockeyId() == entryHorse.JockeyId() {
					return true
				}
				// 鞍上強化
				return containsInSlice([]types.JockeyId{
					5339, // ルメール
					5509, // モレイラ
					5585, // レーン
					5473, // C.デムーロ
					5366, // ムーア
					1088, // 川田将雅
					5299, // 吉原寛人
				}, entryHorse.JockeyId())
			}
		}
	}

	return false
}

// OkNotSlowStart 近2走出遅れがないこと
func (p *placeCheckListService) OkNotSlowStart(ctx context.Context, input *PlaceCheckListInput) bool {
	for _, entryHorse := range input.Race.RaceEntryHorses() {
		if entryHorse.HorseId() == input.Horse.HorseId() {
			if len(input.Horse.HorseResults()) == 0 {
				return true
			}
			historyCount := 0
			for _, horseResult := range input.Horse.HorseResults() {
				// 履歴を取るタイミングで今走が履歴に含まれる場合はスキップする
				if horseResult.RaceId() == input.Race.RaceId() {
					continue
				}
				if historyCount > 2 { // 2走前まで
					break
				}
				if horseResult.Comment() == "出遅れ" {
					return false
				}
				historyCount++
			}
		}
	}

	return true
}

// OkFavoriteRatio 東スポ印◎が50%以上であること
func (p *placeCheckListService) OkFavoriteRatio(ctx context.Context, input *PlaceCheckListInput) bool {
	return float64(input.Forecast.FavoriteNum())/float64(input.Forecast.MarkerNum()) >= 0.5
}

// OkOnlyFavoriteAndRival 東スポ印◎と○東スポ印が◎◯のみで構成されていること
func (p *placeCheckListService) OkOnlyFavoriteAndRival(ctx context.Context, input *PlaceCheckListInput) bool {
	return (input.Forecast.FavoriteNum() + input.Forecast.RivalNum()) == input.Forecast.MarkerNum()
}

// OkIsHighlyRecommended 調教イチ押しであること
func (p *placeCheckListService) OkIsHighlyRecommended(ctx context.Context, input *PlaceCheckListInput) bool {
	return input.Forecast.IsHighlyRecommended()
}
