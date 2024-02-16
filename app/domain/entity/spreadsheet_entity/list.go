package spreadsheet_entity

import (
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/list_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type List struct {
	rows   []*Row
	styles []*Style
}

type Row struct {
	raceDate                string
	class                   string
	courseCategory          string
	distance                string
	traceCondition          string
	raceName                string
	payment                 int
	payout                  int
	favoriteHorse           string
	favoriteJockey          string
	favoriteHorsePopular    int
	favoriteHorseOdds       string
	rivalHorse              string
	rivalJockey             string
	rivalHorsePopular       int
	rivalHorseOdds          string
	firstPlaceHorse         string
	firstPlaceJockey        string
	firstPlaceHorsePopular  int
	firstPlaceHorseOdds     string
	secondPlaceHorse        string
	secondPlaceJockey       string
	secondPlaceHorsePopular int
	secondPlaceHorseOdds    string
}

func NewRow(
	raceDate types.RaceDate,
	class types.GradeClass,
	courseCategory types.CourseCategory,
	distance int,
	traceCondition types.TrackCondition,
	raceName string,
	payment types.Payment,
	payout types.Payout,
	favoriteHorse *list_entity.Horse,
	favoriteJockey string,
	rivalHorse *list_entity.Horse,
	rivalJockey string,
	firstPlaceResult *list_entity.RaceResult,
	firstPlaceJockey string,
	secondPlaceResult *list_entity.RaceResult,
	secondPlaceJockey string,
) *Row {
	return &Row{
		raceDate:                raceDate.Format("2006/01/02"),
		class:                   class.String(),
		courseCategory:          courseCategory.String(),
		distance:                fmt.Sprintf("%dm", distance),
		traceCondition:          traceCondition.String(),
		raceName:                raceName,
		payment:                 payment.Value(),
		payout:                  payout.Value(),
		favoriteHorse:           favoriteHorse.HorseName(),
		favoriteHorsePopular:    favoriteHorse.PopularNumber(),
		favoriteHorseOdds:       favoriteHorse.Odds(),
		favoriteJockey:          favoriteJockey,
		rivalHorse:              rivalHorse.HorseName(),
		rivalJockey:             rivalJockey,
		rivalHorsePopular:       rivalHorse.PopularNumber(),
		rivalHorseOdds:          rivalHorse.Odds(),
		firstPlaceHorse:         firstPlaceResult.HorseName(),
		firstPlaceJockey:        firstPlaceJockey,
		firstPlaceHorsePopular:  firstPlaceResult.PopularNumber(),
		firstPlaceHorseOdds:     firstPlaceResult.Odds(),
		secondPlaceHorse:        secondPlaceResult.HorseName(),
		secondPlaceJockey:       secondPlaceJockey,
		secondPlaceHorsePopular: secondPlaceResult.PopularNumber(),
		secondPlaceHorseOdds:    secondPlaceResult.Odds(),
	}
}

func (r *Row) RaceDate() string {
	return r.raceDate
}

func (r *Row) Class() string {
	return r.class
}

func (r *Row) CourseCategory() string {
	return r.courseCategory
}

func (r *Row) Distance() string {
	return r.distance
}

func (r *Row) TraceCondition() string {
	return r.traceCondition
}

func (r *Row) RaceName() string {
	return r.raceName
}

func (r *Row) Payment() int {
	return r.payment
}

func (r *Row) Payout() int {
	return r.payout
}

func (r *Row) PayoutRate() string {
	return fmt.Sprintf("%.0f%s", float64(r.Payout())*float64(100)/float64(r.Payment()), "%")
}

func (r *Row) FavoriteHorse() string {
	return r.favoriteHorse
}

func (r *Row) FavoriteHorsePopular() int {
	return r.favoriteHorsePopular
}

func (r *Row) FavoriteHorseOdds() string {
	return r.favoriteHorseOdds
}

func (r *Row) FavoriteJockey() string {
	return r.favoriteJockey
}

func (r *Row) RivalHorse() string {
	return r.rivalHorse
}

func (r *Row) RivalJockey() string {
	return r.rivalJockey
}

func (r *Row) RivalHorsePopular() int {
	return r.rivalHorsePopular
}

func (r *Row) RivalHorseOdds() string {
	return r.rivalHorseOdds
}

type Style struct {
	classColor            string
	raceUrl               string
	payoutRateComment     string
	favoriteHorseColor    string
	rivalHorseColor       string
	firstPlaceHorseColor  string
	secondPlaceHorseColor string
}

func NewStyle(
	classColor string,
	raceUrl string,
	payoutRateComment string,
	favoriteHorseColor string,
	rivalHorseColor string,
	firstPlaceHorseColor string,
	secondPlaceHorseColor string,
) *Style {
	return &Style{
		classColor:            classColor,
		raceUrl:               raceUrl,
		payoutRateComment:     payoutRateComment,
		favoriteHorseColor:    favoriteHorseColor,
		rivalHorseColor:       rivalHorseColor,
		firstPlaceHorseColor:  firstPlaceHorseColor,
		secondPlaceHorseColor: secondPlaceHorseColor,
	}
}
