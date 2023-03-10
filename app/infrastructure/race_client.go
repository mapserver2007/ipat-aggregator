package infrastructure

import (
	"context"
	"fmt"
	"github.com/gocolly/colly"
	betting_ticket_entity "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/entity"
	betting_ticket_vo "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/value_object"
	race_entity "github.com/mapserver2007/tools/baken/app/domain/race/entity"
	race_vo "github.com/mapserver2007/tools/baken/app/domain/race/value_object"
	"github.com/mapserver2007/tools/baken/app/repository"
	"regexp"
	"strconv"
	"strings"
)

const (
	raceListUrlForJRA       = "https://race.netkeiba.com/top/race_list_sub.html?kaisai_date=%d"
	raceResultUrlForJRA     = "https://race.netkeiba.com/race/result.html?race_id=%s&organizer=%d"
	raceResultUrlForNAR     = "https://nar.netkeiba.com/race/result.html?race_id=%s&organizer=%d"
	raceResultUrlForOversea = "https://race.netkeiba.com/race/result.html?race_id=%s&organizer=%d"
)

type RaceClient struct {
	client *colly.Collector
}

func NewRaceClient(client *colly.Collector) repository.RaceClient {
	return &RaceClient{
		client: client,
	}
}

func (r *RaceClient) GetRacingNumber(ctx context.Context, entity *betting_ticket_entity.CsvEntity) ([]*race_entity.RacingNumber, error) {
	var racingNumbers []*race_entity.RacingNumber
	url := fmt.Sprintf(raceListUrlForJRA, int(entity.RaceDate))
	r.client.OnHTML(".RaceList_DataList", func(e *colly.HTMLElement) {
		e.ForEach(".RaceList_DataTitle", func(i int, ce *colly.HTMLElement) {
			regex := regexp.MustCompile(`(\d+)回\s+(.+)\s+(\d+)日目`)
			matches := regex.FindAllStringSubmatch(ce.Text, -1)
			round, _ := strconv.Atoi(matches[0][1])
			day, _ := strconv.Atoi(matches[0][3])
			raceCourse := race_vo.ConvertToRaceCourse(matches[0][2])

			racingNumbers = append(racingNumbers, &race_entity.RacingNumber{
				Date:         int(entity.RaceDate),
				Round:        round,
				Day:          day,
				RaceCourseId: raceCourse.Value(),
			})
		})
	})
	err := r.client.Visit(url)
	if err != nil {
		return nil, err
	}

	return racingNumbers, nil
}

func (r *RaceClient) GetRaceResult(ctx context.Context, raceId race_vo.RaceId, organizer race_vo.Organizer) *race_entity.Race {
	var url string
	switch organizer {
	case race_vo.JRA:
		url = fmt.Sprintf(raceResultUrlForJRA, raceId, organizer)
	case race_vo.NAR:
		url = fmt.Sprintf(raceResultUrlForNAR, raceId, organizer)
	case race_vo.OverseaOrganizer:
		url = fmt.Sprintf(raceResultUrlForOversea, raceId, organizer)
	}

	var (
		raceResults              []*race_entity.RaceResult
		payoutResult             []*race_entity.PayoutResult
		raceName, trackCondition string
		raceTimes                []string
		courseCategory           race_vo.CourseCategory
		distance, entries        int
		gradeClass               race_vo.GradeClass
	)
	r.client.OnHTML("#All_Result_Table", func(e *colly.HTMLElement) {
		e.ForEach("tr.HorseList", func(i int, ce *colly.HTMLElement) {
			var numbers []int
			var oddsList []string
			query := ce.Request.URL.Query()
			rawCurrentOrganizer, _ := strconv.Atoi(query.Get("organizer"))
			currentOrganizer := race_vo.Organizer(rawCurrentOrganizer)

			if currentOrganizer == race_vo.JRA {
				ce.ForEach(".Num > div", func(j int, ce2 *colly.HTMLElement) {
					num, _ := strconv.Atoi(ce2.DOM.Text())
					numbers = append(numbers, num)
				})
				ce.ForEach(".Odds span", func(j int, ce2 *colly.HTMLElement) {
					oddsList = append(oddsList, ce2.DOM.Text())
				})
				raceTimes = append(raceTimes, ce.DOM.Find(".Time > .RaceTime").Text())
				popularNumber, _ := strconv.Atoi(oddsList[0])

				raceResults = append(raceResults, &race_entity.RaceResult{
					OrderNo:       i + 1,
					HorseName:     ConvertFromEucJPToUtf8(ce.DOM.Find(".Horse_Name > a").Text()),
					BracketNumber: numbers[0],
					HorseNumber:   numbers[1],
					Odds:          oddsList[1],
					PopularNumber: popularNumber,
				})
			} else if currentOrganizer == race_vo.OverseaOrganizer {
				ce.ForEach(".Num > div", func(j int, ce2 *colly.HTMLElement) {
					num, _ := strconv.Atoi(ce2.DOM.Text())
					numbers = append(numbers, num)
				})
				ce.ForEach(".Odds span", func(j int, ce2 *colly.HTMLElement) {
					oddsList = append(oddsList, ce2.DOM.Text())
				})
				raceTimes = append(raceTimes, ce.DOM.Find(".Time > .RaceTime").Text())
				popularNumber, _ := strconv.Atoi(oddsList[0])

				raceResults = append(raceResults, &race_entity.RaceResult{
					OrderNo:       i + 1,
					HorseName:     ConvertFromEucJPToUtf8(ce.DOM.Find(".Horse_Name > a").Text()),
					BracketNumber: numbers[0],
					HorseNumber:   numbers[1],
					Odds:          oddsList[1],
					PopularNumber: popularNumber,
				})
			}
		})
		e.ForEach("#All_Result_Table > tbody > tr", func(i int, ce *colly.HTMLElement) {
			var numbers []int
			var oddsList []string
			query := ce.Request.URL.Query()
			rawCurrentOrganizer, _ := strconv.Atoi(query.Get("organizer"))
			currentOrganizer := race_vo.Organizer(rawCurrentOrganizer)

			if currentOrganizer == race_vo.NAR {
				ce.ForEach(".Num > div", func(j int, ce2 *colly.HTMLElement) {
					num, _ := strconv.Atoi(ce2.DOM.Text())
					numbers = append(numbers, num)
				})
				ce.ForEach(".Odds span", func(j int, ce2 *colly.HTMLElement) {
					oddsList = append(oddsList, ce2.DOM.Text())
				})
				raceTimes = append(raceTimes, ce.DOM.Find(".Time > .RaceTime").Text())
				popularNumber, _ := strconv.Atoi(oddsList[0])

				raceResults = append(raceResults, &race_entity.RaceResult{
					OrderNo:       i + 1,
					HorseName:     ConvertFromEucJPToUtf8(ce.DOM.Find(".Horse_Name > a").Text()),
					BracketNumber: numbers[0],
					HorseNumber:   numbers[1],
					Odds:          oddsList[1],
					PopularNumber: popularNumber,
				})
			}
		})
	})

	r.client.OnHTML("div.RaceList_Item02", func(e *colly.HTMLElement) {
		e.ForEach("div", func(i int, ce *colly.HTMLElement) {
			query := ce.Request.URL.Query()
			rawCurrentOrganizer, _ := strconv.Atoi(query.Get("organizer"))
			currentOrganizer := race_vo.Organizer(rawCurrentOrganizer)
			if currentOrganizer == race_vo.JRA || currentOrganizer == race_vo.NAR {
				switch i {
				case 0:
					regex := regexp.MustCompile(`(.+)\s`)
					matches := regex.FindAllStringSubmatch(ce.DOM.Text(), -1)
					raceName = ConvertFromEucJPToUtf8(matches[0][1])
					gradeClass = race_vo.AllowanceClass
					if len(ce.DOM.Find(".Icon_GradeType1").Nodes) > 0 {
						gradeClass = race_vo.Grade1
					} else if len(ce.DOM.Find(".Icon_GradeType2").Nodes) > 0 {
						gradeClass = race_vo.Grade2
					} else if len(ce.DOM.Find(".Icon_GradeType3").Nodes) > 0 {
						gradeClass = race_vo.Grade3
					} else if len(ce.DOM.Find(".Icon_GradeType5").Nodes) > 0 {
						gradeClass = race_vo.OpenClass
					} else if len(ce.DOM.Find(".Icon_GradeType10").Nodes) > 0 {
						gradeClass = race_vo.JumpGrade1
					} else if len(ce.DOM.Find(".Icon_GradeType11").Nodes) > 0 {
						gradeClass = race_vo.JumpGrade2
					} else if len(ce.DOM.Find(".Icon_GradeType12").Nodes) > 0 {
						gradeClass = race_vo.JumpGrade3
					} else if len(ce.DOM.Find(".Icon_GradeType15").Nodes) > 0 {
						gradeClass = race_vo.ListedClass
					} else if len(ce.DOM.Find(".Icon_GradeType16").Nodes) > 0 {
						gradeClass = race_vo.AllowanceClass
					} else if len(ce.DOM.Find(".Icon_GradeType17").Nodes) > 0 {
						gradeClass = race_vo.AllowanceClass
					} else if len(ce.DOM.Find(".Icon_GradeType18").Nodes) > 0 {
						gradeClass = race_vo.AllowanceClass
					} else if len(ce.DOM.Find(".Icon_GradeType19").Nodes) > 0 {
						gradeClass = race_vo.Jpn1
					} else if len(ce.DOM.Find(".Icon_GradeType20").Nodes) > 0 {
						gradeClass = race_vo.Jpn2
					} else if len(ce.DOM.Find(".Icon_GradeType21").Nodes) > 0 {
						gradeClass = race_vo.Jpn3
					}
				case 1:
					text := ConvertFromEucJPToUtf8(ce.DOM.Text())
					regex := regexp.MustCompile(`(ダ|芝|障)(\d+)[\s\S]+馬場:(.+)`)
					matches := regex.FindAllStringSubmatch(text, -1)
					courseCategory = race_vo.NewCourseCategory(matches[0][1])
					distance, _ = strconv.Atoi(matches[0][2])
					trackCondition = matches[0][3]
				case 2:
					text := ConvertFromEucJPToUtf8(ce.DOM.Text())
					regex := regexp.MustCompile(`(\d+)頭`)
					matches := regex.FindAllStringSubmatch(text, -1)
					entries, _ = strconv.Atoi(matches[0][1])
				}
			} else if currentOrganizer == race_vo.OverseaOrganizer {
				switch i {
				case 0:
					raceName = ConvertFromEucJPToUtf8(ce.DOM.Text())
					gradeClass = race_vo.NonGrade
					if len(ce.DOM.Find(".Icon_GradeType1").Nodes) > 0 {
						gradeClass = race_vo.Grade1
					} else if len(ce.DOM.Find(".Icon_GradeType2").Nodes) > 0 {
						gradeClass = race_vo.Grade2
					} else if len(ce.DOM.Find(".Icon_GradeType3").Nodes) > 0 {
						gradeClass = race_vo.Grade3
					}
				case 1:
					ce.ForEach("span", func(j int, ce2 *colly.HTMLElement) {
						text := ConvertFromEucJPToUtf8(ce2.DOM.Text())
						switch j {
						case 0:
							regex := regexp.MustCompile(`(ダ|芝)(\d+)`)
							matches := regex.FindAllStringSubmatch(text, -1)
							courseCategory = race_vo.NewCourseCategory(matches[0][1])
							distance, _ = strconv.Atoi(matches[0][2])
						case 1:
							regex := regexp.MustCompile(`(\d+)頭`)
							matches := regex.FindAllStringSubmatch(text, -1)
							entries, _ = strconv.Atoi(matches[0][1])
						case 2:
							regex := regexp.MustCompile(`：(.+)`)
							matches := regex.FindAllStringSubmatch(text, -1)
							trackCondition = matches[0][1]
						}
					})
				}
			}
		})
	})

	r.client.OnHTML("div.Result_Pay_Back table tbody", func(e *colly.HTMLElement) {
		ticketTypeMap := map[string]betting_ticket_vo.BettingTicket{
			"Tansho":  betting_ticket_vo.Win,
			"Fukusho": betting_ticket_vo.Place,
			"Wakuren": betting_ticket_vo.BracketQuinella,
			"Umaren":  betting_ticket_vo.Quinella,
			"Wide":    betting_ticket_vo.QuinellaPlace,
			"Umatan":  betting_ticket_vo.Exacta,
			"Fuku3":   betting_ticket_vo.Trio,
			"Tan3":    betting_ticket_vo.Trifecta,
		}
		e.ForEach("tr", func(i int, ce *colly.HTMLElement) {
			var (
				numbers                        []string
				odds                           []string
				resultSelector, payoutSelector string
			)

			ticketClassName, _ := ce.DOM.Attr("class")
			ticketType, _ := ticketTypeMap[ticketClassName]

			readNumber := func(ce2 *colly.HTMLElement) string {
				str := ConvertFromEucJPToUtf8(ce2.DOM.Text())
				if len(str) == 1 {
					str = fmt.Sprintf("0%s", str)
				}
				return str
			}
			readOdds := func(ce2 *colly.HTMLElement) []string {
				values := strings.Split(ConvertFromEucJPToUtf8(ce2.DOM.Text()), "円")
				values = values[0 : len(values)-1]
				var fixValues []string
				for _, value := range values {
					v := strings.Replace(value, ",", "", -1)
					fixValue, _ := strconv.Atoi(v)
					fixValues = append(fixValues, fmt.Sprintf("%.1f", float64(fixValue)/float64(100)))
				}
				return fixValues
			}

			switch ticketType {
			case betting_ticket_vo.Win, betting_ticket_vo.Place:
				resultSelector = fmt.Sprintf(".%s > .Result > div", ticketClassName)
				payoutSelector = fmt.Sprintf(".%s > .Payout", ticketClassName)
				ce.ForEach(resultSelector, func(j int, ce2 *colly.HTMLElement) {
					switch j {
					case 0, 3, 6:
						number := readNumber(ce2)
						if number != "" {
							numbers = append(numbers, readNumber(ce2))
						}
					}
				})
				ce.ForEach(payoutSelector, func(j int, ce2 *colly.HTMLElement) {
					odds = readOdds(ce2)
				})
			case betting_ticket_vo.BracketQuinella, betting_ticket_vo.Quinella, betting_ticket_vo.QuinellaPlace, betting_ticket_vo.Trio:
				resultSelector = fmt.Sprintf(".%s > .Result > ul > li", ticketClassName)
				payoutSelector = fmt.Sprintf(".%s > .Payout", ticketClassName)
				size := 2
				if ticketType == betting_ticket_vo.Trio {
					size = 3
				}
				numberElems := make([]string, 0, size)
				ce.ForEach(resultSelector, func(j int, ce2 *colly.HTMLElement) {
					numberElem := readNumber(ce2)
					if numberElem != "" {
						numberElems = append(numberElems, numberElem)
						if len(numberElems) == size {
							numbers = append(numbers, strings.Join(numberElems, betting_ticket_vo.QuinellaSeparator))
							numberElems = make([]string, 0, size)
						}
					}
				})
				ce.ForEach(payoutSelector, func(j int, ce2 *colly.HTMLElement) {
					odds = readOdds(ce2)
				})
			case betting_ticket_vo.Exacta, betting_ticket_vo.Trifecta:
				resultSelector = fmt.Sprintf(".%s > .Result > ul > li", ticketClassName)
				payoutSelector = fmt.Sprintf(".%s > .Payout", ticketClassName)
				size := 2
				if ticketType == betting_ticket_vo.Trifecta {
					size = 3
				}
				numberElems := make([]string, 0, size)
				ce.ForEach(resultSelector, func(j int, ce2 *colly.HTMLElement) {
					numberElem := readNumber(ce2)
					if numberElem != "" {
						numberElems = append(numberElems, numberElem)
						if len(numberElems) == size {
							numbers = append(numbers, strings.Join(numberElems, betting_ticket_vo.ExactaSeparator))
							numberElems = make([]string, 0, size)
						}
					}
				})
				ce.ForEach(payoutSelector, func(j int, ce2 *colly.HTMLElement) {
					odds = readOdds(ce2)
				})
			default:
				// NARの場合、枠単があるが今の所集計するつもりがない
				return
			}

			payoutResult = append(payoutResult, &race_entity.PayoutResult{
				TicketType: ticketType.Value(),
				Numbers:    numbers,
				Odds:       odds,
			})
		})
	})

	r.client.Visit(url)

	return &race_entity.Race{
		RaceId:         string(raceId),
		RaceName:       raceName,
		Url:            url,
		Time:           raceTimes[0],
		CourseCategory: int(courseCategory),
		Distance:       distance,
		Entries:        entries,
		Class:          int(gradeClass),
		TrackCondition: trackCondition,
		RaceResults:    raceResults,
		PayoutResults:  payoutResult,
	}
}
