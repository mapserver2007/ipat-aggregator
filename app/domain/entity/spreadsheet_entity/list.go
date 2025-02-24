package spreadsheet_entity

import (
	"fmt"
	"strconv"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/list_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
)

type ListRow struct {
	data  *ListData
	style *ListStyle
}

func (l *ListRow) Data() *ListData {
	return l.data
}

func (l *ListRow) Style() *ListStyle {
	return l.style
}

type ListData struct {
	raceDate                string
	raceStartTime           string
	class                   string
	courseCategory          string
	distance                string
	trackCondition          string
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

func (l *ListData) RaceDate() string {
	return l.raceDate
}

func (l *ListData) RaceStartTime() string {
	return l.raceStartTime
}

func (l *ListData) Class() string {
	return l.class
}

func (l *ListData) CourseCategory() string {
	return l.courseCategory
}

func (l *ListData) Distance() string {
	return l.distance
}

func (l *ListData) TraceCondition() string {
	return l.trackCondition
}

func (l *ListData) RaceName() string {
	return l.raceName
}

func (l *ListData) Url() string {
	return l.url
}

func (l *ListData) Payment() int {
	return l.payment
}

func (l *ListData) Payout() int {
	return l.payout
}

func (l *ListData) PayoutRate() string {
	return fmt.Sprintf("%.0f%s", float64(l.Payout())*float64(100)/float64(l.Payment()), "%")
}

func (l *ListData) FavoriteHorse() string {
	return l.favoriteHorse
}

func (l *ListData) FavoriteJockey() string {
	return l.favoriteJockey
}

func (l *ListData) FavoriteHorsePopular() string {
	return l.favoriteHorsePopular
}

func (l *ListData) FavoriteHorseOdds() string {
	return l.favoriteHorseOdds
}

func (l *ListData) RivalHorse() string {
	return l.rivalHorse
}

func (l *ListData) RivalJockey() string {
	return l.rivalJockey
}

func (l *ListData) RivalHorsePopular() string {
	return l.rivalHorsePopular
}

func (l *ListData) RivalHorseOdds() string {
	return l.rivalHorseOdds
}

func (l *ListData) FirstPlaceHorse() string {
	return l.firstPlaceHorse
}

func (l *ListData) FirstPlaceJockey() string {
	return l.firstPlaceJockey
}

func (l *ListData) FirstPlaceHorsePopular() string {
	return l.firstPlaceHorsePopular
}

func (l *ListData) FirstPlaceHorseOdds() string {
	return l.firstPlaceHorseOdds
}

func (l *ListData) SecondPlaceHorse() string {
	return l.secondPlaceHorse
}

func (l *ListData) SecondPlaceJockey() string {
	return l.secondPlaceJockey
}

func (l *ListData) SecondPlaceHorsePopular() string {
	return l.secondPlaceHorsePopular
}

func (l *ListData) SecondPlaceHorseOdds() string {
	return l.secondPlaceHorseOdds
}

type ListStyle struct {
	classColor            types.CellColorType
	payoutComments        []string
	favoriteHorseColor    types.CellColorType
	rivalHorseColor       types.CellColorType
	firstPlaceHorseColor  types.CellColorType
	secondPlaceHorseColor types.CellColorType
}

func (l *ListStyle) ClassColor() types.CellColorType {
	return l.classColor
}

func (l *ListStyle) PayoutComments() []string {
	return l.payoutComments
}

func (l *ListStyle) FavoriteHorseColor() types.CellColorType {
	return l.favoriteHorseColor
}

func (l *ListStyle) RivalHorseColor() types.CellColorType {
	return l.rivalHorseColor
}

func (l *ListStyle) FirstPlaceHorseColor() types.CellColorType {
	return l.firstPlaceHorseColor
}

func (l *ListStyle) SecondPlaceHorseColor() types.CellColorType {
	return l.secondPlaceHorseColor
}

func NewListRow(
	race *list_entity.Race,
	favoriteHorse *list_entity.Horse,
	rivalHorse *list_entity.Horse,
	favoriteJockey *list_entity.Jockey,
	rivalJockey *list_entity.Jockey,
	firstPlaceResult *list_entity.RaceResult,
	firstPlaceJockey *list_entity.Jockey,
	secondPlaceResult *list_entity.RaceResult,
	secondPlaceJockey *list_entity.Jockey,
	tickets []*list_entity.Ticket,
	payment types.Payment,
	payout types.Payout,
) *ListRow {
	var payoutComments []string
	for _, ticket := range tickets {
		payoutComments = append(payoutComments, fmt.Sprintf("%s %s %s倍 %d円 %d人気",
			ticket.TicketType().OriginTicketType().Name(), ticket.BetNumber().String(), ticket.Odds(), ticket.Payout(), ticket.Popular()))
	}

	oddsFormatFunc := func(odds decimal.Decimal) string {
		if odds.IsZero() {
			return "-"
		}
		return odds.String()
	}

	listData := &ListData{
		raceName:                race.RaceName(),
		raceDate:                race.RaceDate().Format("2006/01/02"),
		raceStartTime:           race.StartTime(),
		class:                   race.Class().String(),
		courseCategory:          race.CourseCategory().String(),
		distance:                fmt.Sprintf("%dm", race.Distance()),
		trackCondition:          race.TrackCondition().String(),
		url:                     race.Url(),
		payment:                 payment.Value(),
		payout:                  payout.Value(),
		favoriteHorse:           favoriteHorse.HorseName(),
		favoriteJockey:          favoriteJockey.JockeyName(),
		favoriteHorsePopular:    strconv.Itoa(favoriteHorse.PopularNumber()),
		favoriteHorseOdds:       oddsFormatFunc(favoriteHorse.Odds()),
		rivalHorse:              rivalHorse.HorseName(),
		rivalJockey:             rivalJockey.JockeyName(),
		rivalHorsePopular:       strconv.Itoa(rivalHorse.PopularNumber()),
		rivalHorseOdds:          oddsFormatFunc(rivalHorse.Odds()),
		firstPlaceHorse:         firstPlaceResult.HorseName(),
		firstPlaceJockey:        firstPlaceJockey.JockeyName(),
		firstPlaceHorsePopular:  strconv.Itoa(firstPlaceResult.PopularNumber()),
		firstPlaceHorseOdds:     oddsFormatFunc(firstPlaceResult.Odds()),
		secondPlaceHorse:        secondPlaceResult.HorseName(),
		secondPlaceJockey:       secondPlaceJockey.JockeyName(),
		secondPlaceHorsePopular: strconv.Itoa(secondPlaceResult.PopularNumber()),
		secondPlaceHorseOdds:    oddsFormatFunc(secondPlaceResult.Odds()),
	}

	classColor := types.NoneColor
	favoriteHorseColor := types.NoneColor
	rivalHorseColor := types.NoneColor
	firstPlaceHorseColor := types.NoneColor
	secondPlaceHorseColor := types.NoneColor

	switch race.Class() {
	case types.Grade1, types.Jpn1:
		classColor = types.FirstColor
	case types.Grade2, types.Jpn2:
		classColor = types.SecondColor
	case types.Grade3, types.Jpn3:
		classColor = types.ThirdColor
	}

	for _, raceResult := range race.RaceResults() {
		if raceResult.HorseName() == favoriteHorse.HorseName() {
			switch raceResult.OrderNo() {
			case 1:
				favoriteHorseColor = types.FirstColor
			case 2:
				favoriteHorseColor = types.SecondColor
			case 3:
				favoriteHorseColor = types.ThirdColor
			}
			switch raceResult.PopularNumber() {
			case 1:
				firstPlaceHorseColor = types.FirstColor
			case 2:
				firstPlaceHorseColor = types.SecondColor
			case 3:
				firstPlaceHorseColor = types.ThirdColor
			}
		}
		if raceResult.HorseName() == rivalHorse.HorseName() {
			switch raceResult.OrderNo() {
			case 1:
				rivalHorseColor = types.FirstColor
			case 2:
				rivalHorseColor = types.SecondColor
			case 3:
				rivalHorseColor = types.ThirdColor
			}
			switch raceResult.PopularNumber() {
			case 1:
				secondPlaceHorseColor = types.FirstColor
			case 2:
				secondPlaceHorseColor = types.SecondColor
			case 3:
				secondPlaceHorseColor = types.ThirdColor
			}
		}
	}

	listStyle := &ListStyle{
		classColor:            classColor,
		payoutComments:        payoutComments,
		favoriteHorseColor:    favoriteHorseColor,
		rivalHorseColor:       rivalHorseColor,
		firstPlaceHorseColor:  firstPlaceHorseColor,
		secondPlaceHorseColor: secondPlaceHorseColor,
	}

	return &ListRow{
		data:  listData,
		style: listStyle,
	}
}
