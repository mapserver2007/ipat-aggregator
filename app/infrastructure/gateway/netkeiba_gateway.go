package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/sirupsen/logrus"
)

type NetKeibaGateway interface {
	FetchRaceId(ctx context.Context, url string) ([]string, error)
	FetchRace(ctx context.Context, url string) (*netkeiba_entity.Race, error)
	FetchRaceCard(ctx context.Context, url string) (*netkeiba_entity.Race, error)
	FetchJockey(ctx context.Context, url string) (*netkeiba_entity.Jockey, error)
	FetchHorse(ctx context.Context, url string) (*netkeiba_entity.Horse, error)
	FetchTrainer(ctx context.Context, url string) (*netkeiba_entity.Trainer, error)
	FetchMarker(ctx context.Context, url string) ([]*netkeiba_entity.Marker, error)
	FetchWinOdds(ctx context.Context, url string) ([]*netkeiba_entity.Odds, error)
	FetchPlaceOdds(ctx context.Context, url string) ([]*netkeiba_entity.Odds, error)
	FetchTrioOdds(ctx context.Context, url string) ([]*netkeiba_entity.Odds, error)
}

type netKeibaGateway struct {
	collector NetKeibaCollector
	logger    *logrus.Logger
	mu        sync.Mutex
}

func NewNetKeibaGateway(
	collector NetKeibaCollector,
	logger *logrus.Logger,
) NetKeibaGateway {
	return &netKeibaGateway{
		collector: collector,
		logger:    logger,
	}
}

func (n *netKeibaGateway) FetchRaceId(
	ctx context.Context,
	url string,
) ([]string, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	var rawRaceIds []string
	n.collector.Client().OnHTML(".RaceList_DataItem > a:first-child", func(e *colly.HTMLElement) {
		regex := regexp.MustCompile(`race_id=(\d+)`)
		matches := regex.FindAllStringSubmatch(e.Attr("href"), -1)
		raceId := matches[0][1]
		rawRaceIds = append(rawRaceIds, raceId)
	})

	n.logger.Infof("fetching race id from %s", url)
	err := n.collector.Client().Visit(url)
	if err != nil {
		if err.Error() == "EOF" { // unreachable url
			return nil, nil
		}
		return nil, err
	}

	return rawRaceIds, nil
}

func (n *netKeibaGateway) FetchRace(
	ctx context.Context,
	url string,
) (*netkeiba_entity.Race, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	var (
		raceResults       []*netkeiba_entity.RaceResult
		payoutResults     []*netkeiba_entity.PayoutResult
		raceName          string
		trackCondition    types.TrackCondition
		startTime         string
		raceTime          string
		courseCategory    types.CourseCategory
		distance, entries int
		gradeClass        types.GradeClass
	)
	raceSexCondition := types.NoRaceSexCondition
	raceWeightCondition := types.FixedWeight

	n.collector.Client().OnHTML("#All_Result_Table", func(e *colly.HTMLElement) {
		raceTime = e.DOM.Find(".Time > .RaceTime").Eq(0).Text()
		e.ForEach("tr.HorseList", func(i int, ce *colly.HTMLElement) {
			var numbers []int
			var oddsList []string
			query := ce.Request.URL.Query()
			rawCurrentOrganizer, _ := strconv.Atoi(query.Get("organizer"))
			currentOrganizer := types.NewOrganizer(rawCurrentOrganizer)

			if currentOrganizer == types.JRA {
				ce.ForEach(".Num > div", func(j int, ce2 *colly.HTMLElement) {
					num, _ := strconv.Atoi(ce2.DOM.Text())
					numbers = append(numbers, num)
				})
				ce.ForEach(".Odds span", func(j int, ce2 *colly.HTMLElement) {
					oddsList = append(oddsList, ce2.DOM.Text())
				})

				popularNumber, _ := strconv.Atoi(oddsList[0])
				linkUrl, _ := ce.DOM.Find(".Jockey > a").Attr("href")
				regex := regexp.MustCompile(`(\d{5})`)
				result := regex.FindStringSubmatch(linkUrl)
				// 一部の騎手で引っかからないjockeyIdの場合があるが、ダミーIDで不明扱いしておく
				jockeyId := "00000"
				if result != nil {
					jockeyId = result[1]
				}
				horseName := Trim(ce.DOM.Find(".Horse_Name > a").Text())
				linkUrl, _ = ce.DOM.Find(".Horse_Name > a").Attr("href")
				segments := strings.Split(linkUrl, "/")
				horseId := segments[4]
				orderNo, _ := strconv.Atoi(ce.DOM.Find(".Rank").Text())

				raceResults = append(raceResults, netkeiba_entity.NewRaceResult(
					orderNo,
					horseId,
					horseName,
					numbers[0],
					numbers[1],
					jockeyId,
					oddsList[1],
					popularNumber,
				))
			} else if currentOrganizer == types.OverseaOrganizer {
				ce.ForEach(".Num > div", func(j int, ce2 *colly.HTMLElement) {
					num, _ := strconv.Atoi(ce2.DOM.Text())
					numbers = append(numbers, num)
				})
				ce.ForEach(".Odds span", func(j int, ce2 *colly.HTMLElement) {
					oddsList = append(oddsList, ce2.DOM.Text())
				})
				popularNumber, _ := strconv.Atoi(oddsList[0])
				linkUrl, _ := ce.DOM.Find(".Jockey > a").Attr("href")
				regex := regexp.MustCompile(`(\d{5})`)
				result := regex.FindStringSubmatch(linkUrl)
				// 一部の騎手で引っかからないjockeyIdの場合があるが、ダミーIDで不明扱いしておく
				jockeyId := "00000"
				if result != nil {
					jockeyId = result[1]
				}
				horseName := Trim(ce.DOM.Find(".Horse_Name > a").Text())
				linkUrl, _ = ce.DOM.Find(".Horse_Name > a").Attr("href")
				segments := strings.Split(linkUrl, "/")
				horseId := segments[4]
				orderNo, _ := strconv.Atoi(ce.DOM.Find(".Rank").Text())

				raceResults = append(raceResults, netkeiba_entity.NewRaceResult(
					orderNo,
					horseId,
					horseName,
					numbers[0],
					numbers[1],
					jockeyId,
					oddsList[1],
					popularNumber,
				))
			}
		})
		e.ForEach("#All_Result_Table > tbody > tr", func(i int, ce *colly.HTMLElement) {
			var numbers []int
			var oddsList []string
			query := ce.Request.URL.Query()
			rawCurrentOrganizer, _ := strconv.Atoi(query.Get("organizer"))
			currentOrganizer := types.NewOrganizer(rawCurrentOrganizer)

			if currentOrganizer == types.NAR {
				ce.ForEach(".Num > div", func(j int, ce2 *colly.HTMLElement) {
					num, _ := strconv.Atoi(ce2.DOM.Text())
					numbers = append(numbers, num)
				})
				ce.ForEach(".Odds span", func(j int, ce2 *colly.HTMLElement) {
					oddsList = append(oddsList, ce2.DOM.Text())
				})
				popularNumber, _ := strconv.Atoi(oddsList[0])
				linkUrl, _ := ce.DOM.Find(".Jockey > a").Attr("href")
				regex := regexp.MustCompile(`(\d{5})`)
				result := regex.FindStringSubmatch(linkUrl)
				// 一部の騎手で引っかからないjockeyIdの場合があるが、ダミーIDで不明扱いしておく
				jockeyId := "00000"
				if result != nil {
					jockeyId = result[1]
				}
				horseName := Trim(ce.DOM.Find(".Horse_Name > a").Text())
				linkUrl, _ = ce.DOM.Find(".Horse_Name > a").Attr("href")
				segments := strings.Split(linkUrl, "/")
				horseId := segments[4]
				orderNo, _ := strconv.Atoi(ce.DOM.Find(".Rank").Text())

				raceResults = append(raceResults, netkeiba_entity.NewRaceResult(
					orderNo,
					horseId,
					horseName,
					numbers[0],
					numbers[1],
					jockeyId,
					oddsList[1],
					popularNumber,
				))
			}
		})
	})

	n.collector.Client().OnHTML("div.RaceList_Item02", func(e *colly.HTMLElement) {
		e.ForEach("h1", func(_ int, ce *colly.HTMLElement) {
			regex := regexp.MustCompile(`(.+)\s`)
			matches := regex.FindAllStringSubmatch(ce.DOM.Text(), -1)
			raceName = Trim(matches[0][1])
			gradeClass = types.AllowanceClass
			if len(ce.DOM.Find(".Icon_GradeType1").Nodes) > 0 {
				gradeClass = types.Grade1
			} else if len(ce.DOM.Find(".Icon_GradeType2").Nodes) > 0 {
				gradeClass = types.Grade2
			} else if len(ce.DOM.Find(".Icon_GradeType3").Nodes) > 0 {
				gradeClass = types.Grade3
			} else if len(ce.DOM.Find(".Icon_GradeType5").Nodes) > 0 {
				if strings.Contains(raceName, "障害") {
					gradeClass = types.JumpOpenClass
				} else {
					gradeClass = types.OpenClass
				}
			} else if len(ce.DOM.Find(".Icon_GradeType10").Nodes) > 0 {
				gradeClass = types.JumpGrade1
			} else if len(ce.DOM.Find(".Icon_GradeType11").Nodes) > 0 {
				gradeClass = types.JumpGrade2
			} else if len(ce.DOM.Find(".Icon_GradeType12").Nodes) > 0 {
				gradeClass = types.JumpGrade3
			} else if len(ce.DOM.Find(".Icon_GradeType15").Nodes) > 0 {
				gradeClass = types.ListedClass
			} else if len(ce.DOM.Find(".Icon_GradeType16").Nodes) > 0 { // 3勝クラス
				gradeClass = types.ThreeWinClass
			} else if len(ce.DOM.Find(".Icon_GradeType17").Nodes) > 0 { // 2勝クラス
				gradeClass = types.TwoWinClass
			} else if len(ce.DOM.Find(".Icon_GradeType18").Nodes) > 0 { // 1勝クラス
				gradeClass = types.OneWinClass
			} else if len(ce.DOM.Find(".Icon_GradeType19").Nodes) > 0 {
				gradeClass = types.Jpn1
			} else if len(ce.DOM.Find(".Icon_GradeType20").Nodes) > 0 {
				gradeClass = types.Jpn2
			} else if len(ce.DOM.Find(".Icon_GradeType21").Nodes) > 0 {
				gradeClass = types.Jpn3
			} else if len(ce.DOM.Find(".Icon_GradeType4").Nodes) > 0 {
				gradeClass = types.LocalGrade
			} else {
				// 条件戦の特別戦、OP、L以外の平場はアイコンが無いのでレース名からクラスを判定する
				if strings.Contains(raceName, "新馬") {
					gradeClass = types.MakeDebut
				} else if strings.Contains(raceName, "未勝利") {
					if strings.Contains(raceName, "障害") {
						gradeClass = types.JumpMaiden
					} else {
						gradeClass = types.Maiden
					}
				} else if strings.Contains(raceName, "1勝クラス") {
					gradeClass = types.OneWinClass
				} else if strings.Contains(raceName, "2勝クラス") {
					gradeClass = types.TwoWinClass
				} else if strings.Contains(raceName, "3勝クラス") {
					gradeClass = types.ThreeWinClass
				}
			}
		})
		e.ForEach("div", func(i int, ce *colly.HTMLElement) {
			query := ce.Request.URL.Query()
			rawCurrentOrganizer, _ := strconv.Atoi(query.Get("organizer"))
			currentOrganizer := types.NewOrganizer(rawCurrentOrganizer)
			if currentOrganizer == types.JRA {
				switch i {
				case 0:
					text := Trim(ce.DOM.Text())
					regex := regexp.MustCompile(`(\d+\:\d+).+(ダ|芝|障)(\d+)[\s\S]+馬場:(.+)`)
					trackConditionText := "良" // 前日だと馬場が表示されない場合があるのでデフォルトで良を設定
					matches := regex.FindAllStringSubmatch(text, -1)
					if matches == nil {
						regex := regexp.MustCompile(`(\d+\:\d+).+(ダ|芝|障)(\d+)`)
						matches = regex.FindAllStringSubmatch(text, -1)
					} else {
						startTime = matches[0][1]
						courseCategory = types.NewCourseCategory(matches[0][2])
						distance, _ = strconv.Atoi(matches[0][3])
						trackConditionText = matches[0][4]
					}

					if strings.Contains(trackConditionText, "良") {
						trackCondition = types.GoodToFirm
					} else if strings.Contains(trackConditionText, "稍") {
						trackCondition = types.Good
					} else if strings.Contains(trackConditionText, "重") {
						trackCondition = types.Yielding
					} else if strings.Contains(trackConditionText, "不") {
						trackCondition = types.Soft
					}
				case 1:
					ce.ForEach("span", func(j int, ce2 *colly.HTMLElement) {
						text := Trim(ce.DOM.Text())
						switch j {
						case 5:
							if strings.Contains(text, "牝") {
								raceSexCondition = types.FillyAndMareLimited
							}
						case 6:
							if text == types.AgeWeight.String() {
								raceWeightCondition = types.AgeWeight
							} else if text == types.FixedWeight.String() {
								raceWeightCondition = types.FixedWeight
							} else if text == types.SpecialWeight.String() {
								raceWeightCondition = types.SpecialWeight
							} else if text == types.HandicapWeight.String() {
								raceWeightCondition = types.HandicapWeight
							}
						case 7:
							regex := regexp.MustCompile(`(\d+)頭`)
							matches := regex.FindAllStringSubmatch(text, -1)
							entries, _ = strconv.Atoi(matches[0][1])
						}
					})
				}
			} else if currentOrganizer == types.NAR {
				switch i {
				case 0:
					regex := regexp.MustCompile(`(.+)\s`)
					matches := regex.FindAllStringSubmatch(ce.DOM.Text(), -1)
					raceName = Trim(matches[0][1])
					gradeClass = types.AllowanceClass
					if len(ce.DOM.Find(".Icon_GradeType1").Nodes) > 0 {
						gradeClass = types.Grade1
					} else if len(ce.DOM.Find(".Icon_GradeType2").Nodes) > 0 {
						gradeClass = types.Grade2
					} else if len(ce.DOM.Find(".Icon_GradeType3").Nodes) > 0 {
						gradeClass = types.Grade3
					} else if len(ce.DOM.Find(".Icon_GradeType19").Nodes) > 0 {
						gradeClass = types.Jpn1
					} else if len(ce.DOM.Find(".Icon_GradeType20").Nodes) > 0 {
						gradeClass = types.Jpn2
					} else if len(ce.DOM.Find(".Icon_GradeType21").Nodes) > 0 {
						gradeClass = types.Jpn3
					} else if len(ce.DOM.Find(".Icon_GradeType4").Nodes) > 0 {
						gradeClass = types.LocalGrade
					}
				case 1:
					text := Trim(ce.DOM.Text())
					regex := regexp.MustCompile(`(\d+\:\d+).+(ダ|芝|障)(\d+)[\s\S]+馬場:(.+)`)
					matches := regex.FindAllStringSubmatch(text, -1)
					startTime = matches[0][1]
					courseCategory = types.NewCourseCategory(matches[0][2])
					distance, _ = strconv.Atoi(matches[0][3])

					trackConditionText := matches[0][4]
					if strings.Contains(trackConditionText, "良") {
						trackCondition = types.GoodToFirm
					} else if strings.Contains(trackConditionText, "稍") {
						trackCondition = types.Good
					} else if strings.Contains(trackConditionText, "重") {
						trackCondition = types.Yielding
					} else if strings.Contains(trackConditionText, "不") {
						trackCondition = types.Soft
					}
				case 2:
					ce.ForEach("span", func(j int, ce2 *colly.HTMLElement) {
						switch j {
						case 5:
							text := Trim(ce.DOM.Text())
							regex := regexp.MustCompile(`(\d+)頭`)
							matches := regex.FindAllStringSubmatch(text, -1)
							entries, _ = strconv.Atoi(matches[0][1])
						}
					})
				}
			} else if currentOrganizer == types.OverseaOrganizer {
				switch i {
				case 0:
					raceName = Trim(ce.DOM.Text())
					gradeClass = types.NonGrade
					if len(ce.DOM.Find(".Icon_GradeType1").Nodes) > 0 {
						gradeClass = types.Grade1
					} else if len(ce.DOM.Find(".Icon_GradeType2").Nodes) > 0 {
						gradeClass = types.Grade2
					} else if len(ce.DOM.Find(".Icon_GradeType3").Nodes) > 0 {
						gradeClass = types.Grade3
					}
				case 1:
					ce.ForEach("span", func(j int, ce2 *colly.HTMLElement) {
						text := Trim(ce2.DOM.Text())
						switch j {
						case 0:
							regex := regexp.MustCompile(`(ダ|芝)(\d+)`)
							matches := regex.FindAllStringSubmatch(text, -1)
							courseCategory = types.NewCourseCategory(matches[0][1])
							distance, _ = strconv.Atoi(matches[0][2])
						case 1:
							regex := regexp.MustCompile(`(\d+)頭`)
							matches := regex.FindAllStringSubmatch(text, -1)
							entries, _ = strconv.Atoi(matches[0][1])
						case 2:
							if text != "" { // 天気アイコンある場合はtextが空
								regex := regexp.MustCompile(`：(.+)`)
								matches := regex.FindAllStringSubmatch(text, -1)
								trackConditionText := matches[0][1]
								if strings.Contains(trackConditionText, "良") {
									trackCondition = types.GoodToFirm
								} else if strings.Contains(trackConditionText, "稍") {
									trackCondition = types.Good
								} else if strings.Contains(trackConditionText, "重") {
									trackCondition = types.Yielding
								} else if strings.Contains(trackConditionText, "不") {
									trackCondition = types.Soft
								}
							}
						case 3:
							regex := regexp.MustCompile(`：(.+)`)
							matches := regex.FindAllStringSubmatch(text, -1)
							trackConditionText := matches[0][1]
							if strings.Contains(trackConditionText, "良") {
								trackCondition = types.GoodToFirm
							} else if strings.Contains(trackConditionText, "稍") {
								trackCondition = types.Good
							} else if strings.Contains(trackConditionText, "重") {
								trackCondition = types.Yielding
							} else if strings.Contains(trackConditionText, "不") {
								trackCondition = types.Soft
							}
						}
					})
				}
			}
		})
	})

	n.collector.Client().OnHTML("div.Result_Pay_Back table tbody", func(e *colly.HTMLElement) {
		ticketTypeMap := map[string]types.TicketType{
			"Tansho":  types.Win,
			"Fukusho": types.Place,
			"Wakuren": types.BracketQuinella,
			"Umaren":  types.Quinella,
			"Wide":    types.QuinellaPlace,
			"Umatan":  types.Exacta,
			"Fuku3":   types.Trio,
			"Tan3":    types.Trifecta,
		}
		e.ForEach("tr", func(i int, ce *colly.HTMLElement) {
			var (
				numbers, odds                                         []string
				populars                                              []int
				resultSelector, payoutSelector, popularNumberSelector string
			)

			ticketClassName, _ := ce.DOM.Attr("class")
			ticketType, _ := ticketTypeMap[ticketClassName]

			readNumber := func(ce2 *colly.HTMLElement) string {
				str := Trim(ce2.DOM.Text())
				if len(str) == 1 {
					str = fmt.Sprintf("0%s", str)
				}
				return str
			}
			readOdds := func(ce2 *colly.HTMLElement) []string {
				values := strings.Split(Trim(ce2.DOM.Text()), "円")
				values = values[0 : len(values)-1]
				var fixValues []string
				for _, value := range values {
					v := strings.Replace(value, ",", "", -1)
					fixValue, _ := strconv.Atoi(v)
					fixValues = append(fixValues, fmt.Sprintf("%.1f", float64(fixValue)/float64(100)))
				}
				return fixValues
			}
			readPopulars := func(ce2 *colly.HTMLElement) []int {
				values := strings.Split(Trim(ce2.DOM.Text()), "人気")
				values = values[0 : len(values)-1]
				var fixValues []int
				for _, value := range values {
					fixValue, _ := strconv.Atoi(value)
					fixValues = append(fixValues, fixValue)
				}
				return fixValues
			}

			switch ticketType {
			case types.Win, types.Place:
				resultSelector = fmt.Sprintf(".%s > .Result > div", ticketClassName)
				payoutSelector = fmt.Sprintf(".%s > .Payout", ticketClassName)
				popularNumberSelector = fmt.Sprintf(".%s > .Ninki", ticketClassName)
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
				ce.ForEach(popularNumberSelector, func(j int, ce2 *colly.HTMLElement) {
					populars = readPopulars(ce2)
				})
			case types.BracketQuinella, types.Quinella, types.QuinellaPlace, types.Trio:
				resultSelector = fmt.Sprintf(".%s > .Result > ul > li", ticketClassName)
				payoutSelector = fmt.Sprintf(".%s > .Payout", ticketClassName)
				popularNumberSelector = fmt.Sprintf(".%s > .Ninki", ticketClassName)
				size := 2
				if ticketType == types.Trio {
					size = 3
				}
				numberElems := make([]string, 0, size)
				ce.ForEach(resultSelector, func(j int, ce2 *colly.HTMLElement) {
					numberElem := readNumber(ce2)
					if numberElem != "" {
						numberElems = append(numberElems, numberElem)
						if len(numberElems) == size {
							numbers = append(numbers, strings.Join(numberElems, types.QuinellaSeparator))
							numberElems = make([]string, 0, size)
						}
					}
				})
				ce.ForEach(payoutSelector, func(j int, ce2 *colly.HTMLElement) {
					odds = readOdds(ce2)
				})
				ce.ForEach(popularNumberSelector, func(j int, ce2 *colly.HTMLElement) {
					populars = readPopulars(ce2)
				})
			case types.Exacta, types.Trifecta:
				resultSelector = fmt.Sprintf(".%s > .Result > ul > li", ticketClassName)
				payoutSelector = fmt.Sprintf(".%s > .Payout", ticketClassName)
				popularNumberSelector = fmt.Sprintf(".%s > .Ninki", ticketClassName)
				size := 2
				if ticketType == types.Trifecta {
					size = 3
				}
				numberElems := make([]string, 0, size)
				ce.ForEach(resultSelector, func(j int, ce2 *colly.HTMLElement) {
					numberElem := readNumber(ce2)
					if numberElem != "" {
						numberElems = append(numberElems, numberElem)
						if len(numberElems) == size {
							numbers = append(numbers, strings.Join(numberElems, types.ExactaSeparator))
							numberElems = make([]string, 0, size)
						}
					}
				})
				ce.ForEach(payoutSelector, func(j int, ce2 *colly.HTMLElement) {
					odds = readOdds(ce2)
				})
				ce.ForEach(popularNumberSelector, func(j int, ce2 *colly.HTMLElement) {
					populars = readPopulars(ce2)
				})
			default:
				// NARの場合、枠単があるが今の所集計するつもりがない
				return
			}

			payoutResults = append(payoutResults, netkeiba_entity.NewPayoutResult(
				ticketType.Value(),
				numbers,
				odds,
				populars,
			))
		})
	})

	parsedUrl, err := neturl.Parse(url)
	if err != nil {
		return nil, err
	}
	queryParams, err := neturl.ParseQuery(parsedUrl.RawQuery)
	if err != nil {
		return nil, err
	}
	raceId := queryParams.Get("race_id")
	organizer, err := strconv.Atoi(queryParams.Get("organizer"))
	if err != nil {
		return nil, err
	}
	raceDate, err := strconv.Atoi(queryParams.Get("race_date"))
	if err != nil {
		return nil, err
	}

	n.logger.Infof("fetching race from %s", url)
	err = n.collector.Client().Visit(url)
	if err != nil {
		return nil, fmt.Errorf("failed to visit url: %s, %v", url, err)
	}

	return netkeiba_entity.NewRace(
		raceId,
		raceDate,
		raceName,
		organizer,
		url,
		raceTime,
		startTime,
		entries,
		distance,
		gradeClass.Value(),
		courseCategory.Value(),
		trackCondition.Value(),
		raceSexCondition.Value(),
		raceWeightCondition.Value(),
		nil,
		raceResults,
		payoutResults,
	), nil
}

func (n *netKeibaGateway) FetchRaceCard(
	ctx context.Context,
	url string,
) (*netkeiba_entity.Race, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	var (
		raceName          string
		raceDate          types.RaceDate
		trackCondition    types.TrackCondition
		startTime         string
		courseCategory    types.CourseCategory
		distance, entries int
		gradeClass        types.GradeClass
		raceEntryHorses   []*netkeiba_entity.RaceEntryHorse
	)
	raceSexCondition := types.NoRaceSexCondition
	raceWeightCondition := types.FixedWeight

	parsedUrl, err := neturl.Parse(url)
	if err != nil {
		return nil, err
	}
	queryParams, err := neturl.ParseQuery(parsedUrl.RawQuery)
	if err != nil {
		return nil, err
	}

	cache := true
	if queryParams.Get("cache") == "false" {
		cache = false
	}
	n.collector.Cache(cache)

	raceId := queryParams.Get("race_id")

	n.collector.Client().OnHTML("#RaceList_DateList dd.Active a", func(e *colly.HTMLElement) {
		path := e.Attr("href")
		u, err := neturl.Parse(path)
		if err != nil {
			fmt.Errorf("failed to parse url: %s, %v", path, err)
		}
		rawRaceDate := u.Query().Get("kaisai_date")
		raceDate, err = types.NewRaceDate(rawRaceDate)
		if err != nil {
			fmt.Errorf("failed to convert to raceDate: %s, %v", rawRaceDate, err)
		}
	})

	n.collector.Client().OnHTML("div.RaceList_Item02", func(e *colly.HTMLElement) {
		e.ForEach("h1", func(_ int, ce *colly.HTMLElement) {
			regex := regexp.MustCompile(`(.+)\s`)
			matches := regex.FindAllStringSubmatch(ce.DOM.Text(), -1)
			raceName = matches[0][1]
			gradeClass = types.AllowanceClass
			if len(ce.DOM.Find(".Icon_GradeType1").Nodes) > 0 {
				gradeClass = types.Grade1
			} else if len(ce.DOM.Find(".Icon_GradeType2").Nodes) > 0 {
				gradeClass = types.Grade2
			} else if len(ce.DOM.Find(".Icon_GradeType3").Nodes) > 0 {
				gradeClass = types.Grade3
			} else if len(ce.DOM.Find(".Icon_GradeType5").Nodes) > 0 {
				if strings.Contains(raceName, "障害") {
					gradeClass = types.JumpOpenClass
				} else {
					gradeClass = types.OpenClass
				}
			} else if len(ce.DOM.Find(".Icon_GradeType10").Nodes) > 0 {
				gradeClass = types.JumpGrade1
			} else if len(ce.DOM.Find(".Icon_GradeType11").Nodes) > 0 {
				gradeClass = types.JumpGrade2
			} else if len(ce.DOM.Find(".Icon_GradeType12").Nodes) > 0 {
				gradeClass = types.JumpGrade3
			} else if len(ce.DOM.Find(".Icon_GradeType15").Nodes) > 0 {
				gradeClass = types.ListedClass
			} else if len(ce.DOM.Find(".Icon_GradeType16").Nodes) > 0 { // 3勝クラス
				gradeClass = types.ThreeWinClass
			} else if len(ce.DOM.Find(".Icon_GradeType17").Nodes) > 0 { // 2勝クラス
				gradeClass = types.TwoWinClass
			} else if len(ce.DOM.Find(".Icon_GradeType18").Nodes) > 0 { // 1勝クラス
				gradeClass = types.OneWinClass
			} else if len(ce.DOM.Find(".Icon_GradeType19").Nodes) > 0 {
				gradeClass = types.Jpn1
			} else if len(ce.DOM.Find(".Icon_GradeType20").Nodes) > 0 {
				gradeClass = types.Jpn2
			} else if len(ce.DOM.Find(".Icon_GradeType21").Nodes) > 0 {
				gradeClass = types.Jpn3
			} else if len(ce.DOM.Find(".Icon_GradeType4").Nodes) > 0 {
				gradeClass = types.LocalGrade
			} else {
				// 条件戦の特別戦、OP、L以外の平場はアイコンが無いのでレース名からクラスを判定する
				if strings.Contains(raceName, "新馬") {
					gradeClass = types.MakeDebut
				} else if strings.Contains(raceName, "未勝利") {
					if strings.Contains(raceName, "障害") {
						gradeClass = types.JumpMaiden
					} else {
						gradeClass = types.Maiden
					}
				} else if strings.Contains(raceName, "1勝クラス") {
					gradeClass = types.OneWinClass
				} else if strings.Contains(raceName, "2勝クラス") {
					gradeClass = types.TwoWinClass
				} else if strings.Contains(raceName, "3勝クラス") {
					gradeClass = types.ThreeWinClass
				}
			}
		})
		e.ForEach("div", func(i int, ce *colly.HTMLElement) {
			switch i {
			case 0:
				text := ce.DOM.Text()
				regex := regexp.MustCompile(`(\d+:\d+).+(ダ|芝|障)(\d+)(?:[\s\S]*馬場:(.+))?`)
				matches := regex.FindAllStringSubmatch(text, -1)
				startTime = matches[0][1]
				courseCategory = types.NewCourseCategory(matches[0][2])
				distance, _ = strconv.Atoi(matches[0][3])

				trackConditionText := matches[0][4]
				// 前日の早い段階では馬場が確定していないため、matches[0][4]は空になる場合があるので暫定で初期値は良にしておく
				trackCondition = types.GoodToFirm
				if strings.Contains(trackConditionText, "良") {
					trackCondition = types.GoodToFirm
				} else if strings.Contains(trackConditionText, "稍") {
					trackCondition = types.Good
				} else if strings.Contains(trackConditionText, "重") {
					trackCondition = types.Yielding
				} else if strings.Contains(trackConditionText, "不") {
					trackCondition = types.Soft
				}
			case 1:
				ce.ForEach("span", func(j int, ce2 *colly.HTMLElement) {
					text := ce.DOM.Text()
					switch j {
					case 5:
						if strings.Contains(text, "牝") {
							raceSexCondition = types.FillyAndMareLimited
						}
					case 6:
						if text == types.AgeWeight.String() {
							raceWeightCondition = types.AgeWeight
						} else if text == types.FixedWeight.String() {
							raceWeightCondition = types.FixedWeight
						} else if text == types.SpecialWeight.String() {
							raceWeightCondition = types.SpecialWeight
						} else if text == types.HandicapWeight.String() {
							raceWeightCondition = types.HandicapWeight
						}
					case 7:
						regex := regexp.MustCompile(`(\d+)頭`)
						matches := regex.FindAllStringSubmatch(text, -1)
						entries, _ = strconv.Atoi(matches[0][1])
					}
				})
			}
		})
	})

	n.collector.Client().OnHTML("div.RaceTableArea table tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, ce *colly.HTMLElement) {
			rawBracketNumber, _ := strconv.Atoi(ce.DOM.Find("td:nth-child(1)").Text())
			rawHorseNumber, _ := strconv.Atoi(ce.DOM.Find("td:nth-child(2)").Text())
			rawHorseName := Trim(ce.DOM.Find("td:nth-child(4)").Text())

			rawHorseId := func() string {
				rawHorseUrl, _ := ce.DOM.Find("td:nth-child(4) .HorseName a").Attr("href")
				parsedUrl, _ = neturl.Parse(rawHorseUrl)
				pathSegments := strings.Split(parsedUrl.Path, "/")
				rawHorseId := pathSegments[2]

				return rawHorseId
			}()

			rawJockeyId := func() string {
				rawJockeyUrl, _ := ce.DOM.Find("td:nth-child(7) a").Attr("href")
				if rawJockeyUrl == "" { // 騎手未定の場合
					return ""
				}
				parsedUrl, _ = neturl.Parse(rawJockeyUrl)
				pathSegments := strings.Split(parsedUrl.Path, "/")
				rawJockeyId := pathSegments[4]

				return rawJockeyId
			}()

			rawTrainerId := func() string {
				rawTrainerUrl, _ := ce.DOM.Find("td:nth-child(8) a").Attr("href")
				parsedUrl, _ = neturl.Parse(rawTrainerUrl)
				pathSegments := strings.Split(parsedUrl.Path, "/")
				rawTrainerId := pathSegments[4]

				return rawTrainerId
			}()

			rawHorseWeight := func() float64 {
				rawHorseWeightText := ce.DOM.Find("td:nth-child(6)").Text()
				rawHorseWeight, _ := strconv.ParseFloat(rawHorseWeightText, 64)

				return rawHorseWeight
			}()

			raceEntryHorses = append(raceEntryHorses, netkeiba_entity.NewRaceEntryHorse(
				rawHorseId,
				rawHorseName,
				rawBracketNumber,
				rawHorseNumber,
				rawJockeyId,
				rawTrainerId,
				rawHorseWeight,
			))
		})
	})

	n.logger.Infof("fetching race card from %s", url)
	err = n.collector.Client().Visit(url)
	if err != nil {
		return nil, err
	}

	return netkeiba_entity.NewRace(
		raceId,
		raceDate.Value(),
		raceName,
		int(types.JRA),
		url,
		"",
		startTime,
		entries,
		distance,
		gradeClass.Value(),
		courseCategory.Value(),
		trackCondition.Value(),
		raceSexCondition.Value(),
		raceWeightCondition.Value(),
		raceEntryHorses,
		nil,
		nil,
	), nil
}

func (n *netKeibaGateway) FetchJockey(
	ctx context.Context,
	url string,
) (*netkeiba_entity.Jockey, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	var name string
	n.collector.Client().OnHTML("div.Name h1", func(e *colly.HTMLElement) {
		list := strings.Split(e.DOM.Text(), "\n")
		name = Trim(list[1][:len(list[1])-2])
	})
	n.collector.Client().OnError(func(r *colly.Response, err error) {
		n.logger.Errorf("GetJockey error: %v", err)
	})

	regex := regexp.MustCompile(`\/jockey\/([0-9a-z]+)\/`)
	result := regex.FindStringSubmatch(url)

	n.logger.Infof("fetching jockey from %s", url)

	err := n.collector.Client().Visit(url)
	if err != nil {
		if err.Error() == "EOF" { // unreachable url
			return netkeiba_entity.NewJockey(result[1], ""), nil
		}
		return nil, err
	}

	return netkeiba_entity.NewJockey(result[1], name), nil
}

func (n *netKeibaGateway) FetchHorse(
	ctx context.Context,
	url string,
) (*netkeiba_entity.Horse, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	var (
		horseId, horseName            string
		trainerId, ownerId, breederId string
		sireId, broodmareSireId       string
		birthDay                      int
		horseBlood                    *netkeiba_entity.HorseBlood
		horseResults                  []*netkeiba_entity.HorseResult
	)

	err := n.collector.Login(ctx)
	if err != nil {
		return nil, err
	}

	parsedUrl, err := neturl.Parse(url)
	if err != nil {
		return nil, err
	}
	queryParams, err := neturl.ParseQuery(parsedUrl.RawQuery)
	if err != nil {
		return nil, err
	}

	cache := true
	if queryParams.Get("cache") == "false" {
		cache = false
	}
	n.collector.Cache(cache)

	segments := strings.Split(parsedUrl.Path, "/")
	horseId = segments[2]

	n.collector.Client().OnHTML("div.horse_title h1", func(e *colly.HTMLElement) {
		horseName = Trim(e.DOM.Text())
	})

	n.collector.Client().OnHTML("table.db_prof_table tbody", func(e *colly.HTMLElement) {
		rowCount := e.DOM.Find("tr").Length()
		e.ForEach("tr", func(i int, ce *colly.HTMLElement) {
			switch i {
			case 0:
				birthDayStr := ce.DOM.Find("td:nth-child(2)").Text()
				layout := "2006年1月2日"
				date, _ := time.Parse(layout, birthDayStr)
				rawBirthDay, _ := strconv.Atoi(date.Format("20060102"))
				birthDay = rawBirthDay
			case 1:
				path, _ := ce.DOM.Find("td:nth-child(2) a").Attr("href")
				segments = strings.Split(path, "/")
				trainerId = segments[2]
			case 2:
				path, _ := ce.DOM.Find("td:nth-child(2) a").Attr("href")
				segments = strings.Split(path, "/")
				ownerId = segments[2]
			}

			if rowCount == 10 && i == 3 || rowCount == 11 && i == 4 { // 個人馬主 or 一口会員
				path, _ := ce.DOM.Find("td:nth-child(2) a").Attr("href")
				segments = strings.Split(path, "/")
				breederId = segments[2]
			}
		})
	})

	n.collector.Client().OnHTML("table.blood_table tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, ce *colly.HTMLElement) {
			switch i {
			case 0:
				path, _ := ce.DOM.Find("td:nth-child(1) a").Attr("href")
				segments = strings.Split(path, "/")
				sireId = segments[3]
			case 2:
				path, _ := ce.DOM.Find("td:nth-child(2) a").Attr("href")
				segments = strings.Split(path, "/")
				broodmareSireId = segments[3]
			}
		})
		horseBlood = netkeiba_entity.NewHorseBlood(sireId, broodmareSireId)
	})

	n.collector.Client().OnHTML("table.db_h_race_results tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, ce *colly.HTMLElement) {
			raceDateStr := strings.Replace(ce.DOM.Find("td:nth-child(1)").Text(), "/", "", -1)
			raceDate, _ := strconv.Atoi(raceDateStr)

			path, _ := ce.DOM.Find("td:nth-child(2) a").Attr("href")
			segments = strings.Split(path, "/")
			typedRaceCourse := types.RaceCourse(segments[3])
			raceCourseId := typedRaceCourse.Value() // primitiveで扱いたいのであえて戻す

			// 他のクラスについては判別が困難なので一律NonGradeとして返す
			// 障害はそもそも予想しないので対応しない
			typedGradeClass := types.NonGrade
			{
				rawRaceName := ce.DOM.Find("td:nth-child(5)").Text()
				if strings.Contains(rawRaceName, "(GIII)") {
					typedGradeClass = types.Grade3
				} else if strings.Contains(rawRaceName, "(GII)") {
					typedGradeClass = types.Grade2
				} else if strings.Contains(rawRaceName, "(GI)") {
					typedGradeClass = types.Grade1
				} else if strings.Contains(rawRaceName, "(JpnIII)") {
					typedGradeClass = types.Jpn3
				} else if strings.Contains(rawRaceName, "(JpnII)") {
					typedGradeClass = types.Jpn2
				} else if strings.Contains(rawRaceName, "(JpnI)") {
					typedGradeClass = types.Jpn1
				} else if strings.Contains(rawRaceName, "(OP)") {
					typedGradeClass = types.OpenClass
				} else if strings.Contains(rawRaceName, "(L)") {
					typedGradeClass = types.ListedClass
				} else if strings.Contains(rawRaceName, "(3勝クラス)") {
					typedGradeClass = types.ThreeWinClass
				} else if strings.Contains(rawRaceName, "(2勝クラス)") {
					typedGradeClass = types.TwoWinClass
				} else if strings.Contains(rawRaceName, "(1勝クラス)") {
					typedGradeClass = types.OneWinClass
				} else if strings.Contains(rawRaceName, "未勝利") {
					typedGradeClass = types.Maiden
				} else if strings.Contains(rawRaceName, "新馬") {
					typedGradeClass = types.MakeDebut
				}
			}
			gradeClass := typedGradeClass.Value()

			path, _ = ce.DOM.Find("td:nth-child(5) a").Attr("href")
			segments = strings.Split(path, "/")
			raceId := segments[2]

			raceName := ce.DOM.Find("td:nth-child(5) a").Text()
			entries, _ := strconv.Atoi(ce.DOM.Find("td:nth-child(7)").Text())
			horseNumber, _ := strconv.Atoi(ce.DOM.Find("td:nth-child(9)").Text())
			odds := Trim(ce.DOM.Find("td:nth-child(10)").Text()) // 海外レースなどでは空になる場合あり

			popularNumber, _ := strconv.Atoi(ce.DOM.Find("td:nth-child(11)").Text())
			rawOrderNo := ce.DOM.Find("td:nth-child(12)").Text()
			orderNo := 0
			if rawOrderNo != "中" { // 競走中止
				orderNo, _ = strconv.Atoi(ce.DOM.Find("td:nth-child(12)").Text())
			}
			raceWeight, _ := strconv.ParseFloat(ce.DOM.Find("td:nth-child(14)").Text(), 64)

			path, _ = ce.DOM.Find("td:nth-child(13) a").Attr("href")
			segments = strings.Split(path, "/")
			jockeyId := segments[4]

			rawHorseWeight := ce.DOM.Find("td:nth-child(24)").Text()
			horseWeight := 0
			if rawHorseWeight != "計不" { // 海外レースなどでは計量不可
				regex := regexp.MustCompile(`^\d+`)
				horseWeight, _ = strconv.Atoi(regex.FindString(rawHorseWeight))
			}

			trackAndDistanceStr := ce.DOM.Find("td:nth-child(15)").Text()
			regex := regexp.MustCompile(`(ダ|芝|障)(\d+)`)
			matches := regex.FindAllStringSubmatch(trackAndDistanceStr, -1)

			typedCourseCategory := types.NewCourseCategory(matches[0][1])
			courseCategoryId := typedCourseCategory.Value() // primitiveで扱いたいのであえて戻す
			distance, _ := strconv.Atoi(matches[0][2])

			trackConditionStr := ce.DOM.Find("td:nth-child(16)").Text()
			trackCondition := types.NewTrackCondition(trackConditionStr).Value() // primitiveで扱いたいのであえて戻す
			comment := Trim(ce.DOM.Find("td:nth-child(26)").Text())

			horseResults = append(horseResults, netkeiba_entity.NewHorseResult(
				raceId,
				raceDate,
				raceName,
				jockeyId,
				orderNo,
				popularNumber,
				horseNumber,
				odds,
				gradeClass,
				entries,
				distance,
				raceCourseId,
				courseCategoryId,
				trackCondition,
				horseWeight,
				raceWeight,
				comment,
			))
		})
	})

	n.collector.Client().OnError(func(r *colly.Response, err error) {
		n.logger.Errorf("GetHorse error: %v", err)
	})

	n.logger.Infof("fetching horse from %s", url)

	err = n.collector.Client().Visit(url)
	if err != nil {
		return nil, err
	}

	return netkeiba_entity.NewHorse(
		horseId,
		horseName,
		birthDay,
		trainerId,
		ownerId,
		breederId,
		horseBlood,
		horseResults,
	), nil
}

func (n *netKeibaGateway) FetchTrainer(
	ctx context.Context,
	url string,
) (*netkeiba_entity.Trainer, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	var trainerId, trainerName, locationName string

	u, _ := neturl.Parse(url)
	segments := strings.Split(u.Path, "/")
	trainerId = segments[2]

	n.collector.Client().OnHTML("div.db_head_name .Name", func(e *colly.HTMLElement) {
		str1 := strings.ReplaceAll(e.DOM.Find("h1").Text(), " ", "")
		segments = strings.Split(str1, "\n")
		trainerName = Trim(segments[1])
		str2 := strings.ReplaceAll(e.DOM.Find("p").Text(), " ", "")
		segments = strings.Split(str2, "\n")
		locationName = Trim(segments[3])
	})

	n.logger.Infof("fetching trainer from %s", url)

	err := n.collector.Client().Visit(url)
	if err != nil {
		return nil, err
	}

	return netkeiba_entity.NewTrainer(
		trainerId,
		trainerName,
		locationName,
	), nil
}

func (n *netKeibaGateway) FetchMarker(
	ctx context.Context,
	url string,
) ([]*netkeiba_entity.Marker, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	cookies, err := n.collector.Cookies(ctx)
	if err != nil {
		return nil, err
	}

	parsedUrl, err := neturl.Parse(url)
	if err != nil {
		return nil, err
	}
	queryParams, err := neturl.ParseQuery(parsedUrl.RawQuery)
	if err != nil {
		return nil, err
	}
	raceId := queryParams.Get("race_id")

	data := neturl.Values{}
	data.Set("action", "get")
	data.Set("pid", "api_post_social_cart")
	data.Set("group", fmt.Sprintf("horse_%s", raceId))

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var markerInfo *raw_entity.MarkerInfo
	json.Unmarshal(body, &markerInfo)

	if markerInfo == nil {
		return nil, nil
	}

	markers := make([]*netkeiba_entity.Marker, 0, len(markerInfo.Data))
	for _, d := range markerInfo.Data {
		segments := strings.Split(d.Code, "_")
		marker, err := netkeiba_entity.NewMarker(
			segments[0],
			segments[1],
		)
		if err != nil {
			return nil, err
		}
		if marker == nil {
			continue
		}
		markers = append(markers, marker)
	}

	n.logger.Infof("fetching marker from %s", url)

	return markers, nil
}

func (n *netKeibaGateway) FetchWinOdds(
	ctx context.Context,
	url string,
) ([]*netkeiba_entity.Odds, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.logger.Infof("fetching win odds from %s", url)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var oddsInfo *raw_entity.OddsInfo
	if err := json.Unmarshal(body, &oddsInfo); err != nil {
		n.logger.Errorf("Odds is not published: %s", url)
		return nil, err
	}

	dateTime, err := time.Parse("2006-01-02 15:04:05", oddsInfo.Data.OfficialDatetime)
	if err != nil {
		return nil, err
	}
	raceDate, err := types.NewRaceDate(dateTime.Format("20060102"))
	if err != nil {
		return nil, err
	}

	var odds []*netkeiba_entity.Odds
	for rawNumber, list := range oddsInfo.Data.Odds.Wins {
		popularNumber, _ := strconv.Atoi(list[2])
		// 9999人気の値は取り消しという仕様なので除外する
		if popularNumber == 9999 {
			continue
		}
		rawHorseNumber, _ := strconv.Atoi(rawNumber)
		horseNumber := types.HorseNumber(rawHorseNumber)
		odds = append(odds, netkeiba_entity.NewOdds(
			types.Win, []string{list[0], list[1]}, popularNumber, []types.HorseNumber{horseNumber}, raceDate,
		))
	}

	return odds, nil
}

func (n *netKeibaGateway) FetchPlaceOdds(
	ctx context.Context,
	url string,
) ([]*netkeiba_entity.Odds, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.logger.Infof("fetching place odds from %s", url)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var oddsInfo *raw_entity.OddsInfo
	if err := json.Unmarshal(body, &oddsInfo); err != nil {
		n.logger.Errorf("Odds is not published: %s", url)
		return nil, err
	}

	dateTime, err := time.Parse("2006-01-02 15:04:05", oddsInfo.Data.OfficialDatetime)
	if err != nil {
		return nil, err
	}
	raceDate, err := types.NewRaceDate(dateTime.Format("20060102"))
	if err != nil {
		return nil, err
	}

	var odds []*netkeiba_entity.Odds
	for rawNumber, list := range oddsInfo.Data.Odds.Places {
		popularNumber, _ := strconv.Atoi(list[2])
		// 9999人気の値は取り消しという仕様なので除外する
		if popularNumber == 9999 {
			continue
		}
		rawHorseNumber, _ := strconv.Atoi(rawNumber)
		horseNumber := types.HorseNumber(rawHorseNumber)
		odds = append(odds, netkeiba_entity.NewOdds(
			types.Place, []string{list[0], list[1]}, popularNumber, []types.HorseNumber{horseNumber}, raceDate,
		))
	}

	return odds, nil
}

func (n *netKeibaGateway) FetchTrioOdds(
	ctx context.Context,
	url string,
) ([]*netkeiba_entity.Odds, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.logger.Infof("fetching trio odds from %s", url)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var oddsInfo *raw_entity.OddsInfo
	if err := json.Unmarshal(body, &oddsInfo); err != nil {
		return nil, err
	}

	dateTime, err := time.Parse("2006-01-02 15:04:05", oddsInfo.Data.OfficialDatetime)
	if err != nil {
		return nil, err
	}
	raceDate, err := types.NewRaceDate(dateTime.Format("20060102"))
	if err != nil {
		return nil, err
	}

	var odds []*netkeiba_entity.Odds
	for rawNumber, list := range oddsInfo.Data.Odds.Trios {
		rawHorseNumber1, _ := strconv.Atoi(rawNumber[0:2])
		rawHorseNumber2, _ := strconv.Atoi(rawNumber[2:4])
		rawHorseNumber3, _ := strconv.Atoi(rawNumber[4:6])
		horseNumber1 := types.HorseNumber(rawHorseNumber1)
		horseNumber2 := types.HorseNumber(rawHorseNumber2)
		horseNumber3 := types.HorseNumber(rawHorseNumber3)
		popularNumber, _ := strconv.Atoi(list[2])
		odds = append(odds, netkeiba_entity.NewOdds(
			types.Trio, []string{list[0], list[1]}, popularNumber, []types.HorseNumber{horseNumber1, horseNumber2, horseNumber3}, raceDate,
		))
	}

	return odds, nil
}
