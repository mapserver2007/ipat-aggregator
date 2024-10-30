package spreadsheet_entity

import (
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
)

type PredictionCheckList struct {
	raceDate          string
	raceName          string
	raceUrl           string
	horseName         string
	horseUrl          string
	jockeyName        string
	jockeyUrl         string
	trainerName       string
	trainerUrl        string
	locationName      string
	winOdds           string
	marker            string
	firstPlaceRate    string
	secondPlaceRate   string
	thirdPlaceRate    string
	checkList         []string
	favoriteNum       int
	rivalNum          int
	markerNum         int
	highlyRecommended string
	trainingComment   string
	paperUrl          string
}

func NewPredictionCheckList(
	raceId types.RaceId,
	raceDate types.RaceDate,
	raceName string,
	raceNumber int,
	raceCourse types.RaceCourse,
	horseId types.HorseId,
	horseName string,
	jockeyId types.JockeyId,
	jockeyName string,
	trainerId types.TrainerId,
	trainerName string,
	locationId types.LocationId,
	winOdds decimal.Decimal,
	marker types.Marker,
	firstPlaceRate string,
	secondPlaceRate string,
	thirdPlaceRate string,
	checkList []bool,
	favoriteNum int,
	rivalNum int,
	markerNum int,
	highlyRecommended bool,
	trainingComment string,
) *PredictionCheckList {
	checkListFormat := make([]string, len(checkList))
	for idx, isCheck := range checkList {
		if isCheck {
			checkListFormat[idx] = "◯"
		} else {
			checkListFormat[idx] = "-"
		}
	}

	var highlyRecommendedFormat string
	if highlyRecommended {
		highlyRecommendedFormat = "◯"
	} else {
		highlyRecommendedFormat = "-"
	}

	return &PredictionCheckList{
		raceDate:          raceDate.Format("2006/01/02"),
		raceName:          fmt.Sprintf("%s %dR %s", raceCourse.Name(), raceNumber, raceName),
		raceUrl:           fmt.Sprintf("https://race.netkeiba.com/race/shutuba.html?race_id=%s", raceId.String()),
		horseName:         horseName,
		horseUrl:          fmt.Sprintf("https://db.netkeiba.com/horse/%s", horseId),
		jockeyName:        jockeyName,
		jockeyUrl:         fmt.Sprintf("https://db.netkeiba.com/jockey/%s", jockeyId),
		trainerName:       trainerName,
		trainerUrl:        fmt.Sprintf("https://db.netkeiba.com/trainer/%s", trainerId),
		locationName:      locationId.Name(),
		winOdds:           winOdds.String(),
		marker:            marker.String(),
		firstPlaceRate:    firstPlaceRate,
		secondPlaceRate:   secondPlaceRate,
		thirdPlaceRate:    thirdPlaceRate,
		checkList:         checkListFormat,
		favoriteNum:       favoriteNum,
		rivalNum:          rivalNum,
		markerNum:         markerNum,
		highlyRecommended: highlyRecommendedFormat,
		trainingComment:   trainingComment,
		paperUrl:          "https://tospo-keiba.jp/newspaper-list",
	}
}

func (p *PredictionCheckList) RaceDate() string {
	return p.raceDate
}

func (p *PredictionCheckList) RaceName() string {
	return p.raceName
}

func (p *PredictionCheckList) RaceUrl() string {
	return p.raceUrl
}

func (p *PredictionCheckList) HorseName() string {
	return p.horseName
}

func (p *PredictionCheckList) HorseUrl() string {
	return p.horseUrl
}

func (p *PredictionCheckList) JockeyName() string {
	return p.jockeyName
}

func (p *PredictionCheckList) JockeyUrl() string {
	return p.jockeyUrl
}

func (p *PredictionCheckList) TrainerName() string {
	return p.trainerName
}

func (p *PredictionCheckList) TrainerUrl() string {
	return p.trainerUrl
}

func (p *PredictionCheckList) LocationName() string {
	return p.locationName
}

func (p *PredictionCheckList) WinOdds() string {
	return p.winOdds
}

func (p *PredictionCheckList) Marker() string {
	return p.marker
}

func (p *PredictionCheckList) FirstPlaceRate() string {
	return p.firstPlaceRate
}

func (p *PredictionCheckList) SecondPlaceRate() string {
	return p.secondPlaceRate
}

func (p *PredictionCheckList) ThirdPlaceRate() string {
	return p.thirdPlaceRate
}

func (p *PredictionCheckList) CheckList() []string {
	return p.checkList
}

func (p *PredictionCheckList) FavoriteNum() int {
	return p.favoriteNum
}

func (p *PredictionCheckList) RivalNum() int {
	return p.rivalNum
}

func (p *PredictionCheckList) MarkerNum() int {
	return p.markerNum
}

func (p *PredictionCheckList) HighlyRecommended() string {
	return p.highlyRecommended
}

func (p *PredictionCheckList) TrainingComment() string {
	return p.trainingComment
}

func (p *PredictionCheckList) PaperUrl() string {
	return p.paperUrl
}
