package spreadsheet_entity

import (
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/list_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"strconv"
)

type List struct {
	rows   []*Row
	styles []*Style
}

func NewList(
	rows []*Row,
	styles []*Style,
) *List {
	return &List{
		rows:   rows,
		styles: styles,
	}
}

func (l *List) Rows() []*Row {
	return l.rows
}

func (l *List) Styles() []*Style {
	return l.styles
}

type Row struct {
	raceDate                string
	class                   string
	courseCategory          string
	distance                string
	traceCondition          string
	raceName                string
	url                     string
	payment                 int
	payout                  int
	favoriteHorse           string
	favoriteJockey          string
	favoriteHorsePopular    string
	favoriteHorseOdds       string
	rivalHorse              string
	rivalJockey             string
	rivalHorsePopular       string
	rivalHorseOdds          string
	firstPlaceHorse         string
	firstPlaceJockey        string
	firstPlaceHorsePopular  string
	firstPlaceHorseOdds     string
	secondPlaceHorse        string
	secondPlaceJockey       string
	secondPlaceHorsePopular string
	secondPlaceHorseOdds    string
}

func NewRow(
	raceDate types.RaceDate,
	class types.GradeClass,
	courseCategory types.CourseCategory,
	distance int,
	traceCondition types.TrackCondition,
	raceName string,
	url string,
	payment types.Payment,
	payout types.Payout,
	favoriteHorse *list_entity.Horse,
	favoriteJockey *list_entity.Jockey,
	rivalHorse *list_entity.Horse,
	rivalJockey *list_entity.Jockey,
	firstPlaceResult *list_entity.RaceResult,
	firstPlaceJockey *list_entity.Jockey,
	secondPlaceResult *list_entity.RaceResult,
	secondPlaceJockey *list_entity.Jockey,
) *Row {
	return &Row{
		raceDate:                raceDate.Format("2006/01/02"),
		class:                   class.String(),
		courseCategory:          courseCategory.String(),
		distance:                fmt.Sprintf("%dm", distance),
		traceCondition:          traceCondition.String(),
		raceName:                raceName,
		url:                     url,
		payment:                 payment.Value(),
		payout:                  payout.Value(),
		favoriteHorse:           favoriteHorse.HorseName(),
		favoriteHorsePopular:    strconv.Itoa(favoriteHorse.PopularNumber()),
		favoriteHorseOdds:       favoriteHorse.Odds(),
		favoriteJockey:          favoriteJockey.JockeyName(),
		rivalHorse:              rivalHorse.HorseName(),
		rivalJockey:             rivalJockey.JockeyName(),
		rivalHorsePopular:       strconv.Itoa(rivalHorse.PopularNumber()),
		rivalHorseOdds:          rivalHorse.Odds(),
		firstPlaceHorse:         firstPlaceResult.HorseName(),
		firstPlaceJockey:        firstPlaceJockey.JockeyName(),
		firstPlaceHorsePopular:  firstPlaceResult.HorseName(),
		firstPlaceHorseOdds:     firstPlaceResult.Odds(),
		secondPlaceHorse:        secondPlaceResult.HorseName(),
		secondPlaceJockey:       secondPlaceJockey.JockeyName(),
		secondPlaceHorsePopular: secondPlaceResult.HorseName(),
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

func (r *Row) Url() string {
	return r.url
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

func (r *Row) FavoriteHorsePopular() string {
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

func (r *Row) RivalHorsePopular() string {
	return r.rivalHorsePopular
}

func (r *Row) RivalHorseOdds() string {
	return r.rivalHorseOdds
}

func (r *Row) FirstPlaceHorse() string {
	return r.firstPlaceHorse
}

func (r *Row) FirstPlaceJockey() string {
	return r.firstPlaceJockey
}

func (r *Row) FirstPlaceHorsePopular() string {
	return r.firstPlaceHorsePopular
}

func (r *Row) FirstPlaceHorseOdds() string {
	return r.firstPlaceHorseOdds
}

func (r *Row) SecondPlaceHorse() string {
	return r.secondPlaceHorse
}

func (r *Row) SecondPlaceJockey() string {
	return r.secondPlaceJockey

}

func (r *Row) SecondPlaceHorsePopular() string {
	return r.secondPlaceHorsePopular
}

func (r *Row) SecondPlaceHorseOdds() string {
	return r.secondPlaceHorseOdds

}

type Style struct {
	classColor            types.CellColorType
	payoutComments        []string
	favoriteHorseColor    types.CellColorType
	rivalHorseColor       types.CellColorType
	firstPlaceHorseColor  types.CellColorType
	secondPlaceHorseColor types.CellColorType
}

func NewStyle(
	classColor types.CellColorType,
	payoutComments []string,
	favoriteHorseColor types.CellColorType,
	rivalHorseColor types.CellColorType,
	firstPlaceHorseColor types.CellColorType,
	secondPlaceHorseColor types.CellColorType,
) *Style {
	return &Style{
		classColor:            classColor,
		payoutComments:        payoutComments,
		favoriteHorseColor:    favoriteHorseColor,
		rivalHorseColor:       rivalHorseColor,
		firstPlaceHorseColor:  firstPlaceHorseColor,
		secondPlaceHorseColor: secondPlaceHorseColor,
	}
}

func (s *Style) ClassColor() types.CellColorType {
	return s.classColor
}

func (s *Style) PayoutComments() []string {
	return s.payoutComments
}

func (s *Style) FavoriteHorseColor() types.CellColorType {
	return s.favoriteHorseColor
}

func (s *Style) RivalHorseColor() types.CellColorType {
	return s.rivalHorseColor
}

func (s *Style) FirstPlaceHorseColor() types.CellColorType {
	return s.firstPlaceHorseColor
}

func (s *Style) SecondPlaceHorseColor() types.CellColorType {
	return s.secondPlaceHorseColor
}
