package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type raceDataRepository struct {
	client *colly.Collector
}

func NewRaceDataRepository() repository.RaceDataRepository {
	return &raceDataRepository{
		client: colly.NewCollector(),
	}
}

func (r *raceDataRepository) Read(ctx context.Context, fileName string) ([]*raw_entity.Race, error) {
	races := make([]*raw_entity.Race, 0)
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path, err := filepath.Abs(fmt.Sprintf("%s/cache/%s", rootPath, fileName))
	if err != nil {
		return nil, err
	}

	// ファイルが存在しない場合はエラーは返さず処理を継続する
	bytes, err := os.ReadFile(path)
	if err != nil {
		return races, nil
	}

	var raceInfo *raw_entity.RaceInfo
	if err := json.Unmarshal(bytes, &raceInfo); err != nil {
		return nil, err
	}
	races = raceInfo.Races

	return races, nil
}

func (r *raceDataRepository) Write(
	ctx context.Context,
	fileName string,
	raceInfo *raw_entity.RaceInfo,
) error {
	bytes, err := json.Marshal(raceInfo)
	if err != nil {
		return err
	}

	rootPath, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath, err := filepath.Abs(fmt.Sprintf("%s/cache/%s", rootPath, fileName))
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *raceDataRepository) Fetch(
	ctx context.Context,
	url string,
) (*netkeiba_entity.Race, error) {
	var (
		raceResults              []*netkeiba_entity.RaceResult
		payoutResults            []*netkeiba_entity.PayoutResult
		raceName, trackCondition string
		startTime                string
		raceTimes                []string
		courseCategory           types.CourseCategory
		distance, entries        int
		gradeClass               types.GradeClass
	)
	r.client.OnHTML("#All_Result_Table", func(e *colly.HTMLElement) {
		e.ForEach("tr.HorseList", func(i int, ce *colly.HTMLElement) {
			var numbers []int
			var oddsList []string
			query := ce.Request.URL.Query()
			rawCurrentOrganizer, _ := strconv.Atoi(query.Get("organizer"))
			currentOrganizer := types.Organizer(rawCurrentOrganizer)

			if currentOrganizer == types.JRA {
				ce.ForEach(".Num > div", func(j int, ce2 *colly.HTMLElement) {
					num, _ := strconv.Atoi(ce2.DOM.Text())
					numbers = append(numbers, num)
				})
				ce.ForEach(".Odds span", func(j int, ce2 *colly.HTMLElement) {
					oddsList = append(oddsList, ce2.DOM.Text())
				})
				raceTimes = append(raceTimes, ce.DOM.Find(".Time > .RaceTime").Text())
				popularNumber, _ := strconv.Atoi(oddsList[0])
				linkUrl, _ := ce.DOM.Find(".Jockey > a").Attr("href")
				regex := regexp.MustCompile(`(\d{5})`)
				result := regex.FindStringSubmatch(linkUrl)
				// 一部の騎手で\d{5}で引っかからないjockeyIdの場合があるが、マイナーな騎手なので無視する
				jockeyId := 0
				if result != nil {
					jockeyId, _ = strconv.Atoi(result[1])
				}

				raceResults = append(raceResults, netkeiba_entity.NewRaceResult(
					i+1,
					ConvertFromEucJPToUtf8(ce.DOM.Find(".Horse_Name > a").Text()),
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
				raceTimes = append(raceTimes, ce.DOM.Find(".Time > .RaceTime").Text())
				popularNumber, _ := strconv.Atoi(oddsList[0])
				linkUrl, _ := ce.DOM.Find(".Jockey > a").Attr("href")
				regex := regexp.MustCompile(`(\d{5})`)
				result := regex.FindStringSubmatch(linkUrl)
				// 一部の騎手で\d{5}で引っかからないjockeyIdの場合があるが、マイナーな騎手なので無視する
				jockeyId := 0
				if result != nil {
					jockeyId, _ = strconv.Atoi(result[1])
				}

				raceResults = append(raceResults, netkeiba_entity.NewRaceResult(
					i+1,
					ConvertFromEucJPToUtf8(ce.DOM.Find(".Horse_Name > a").Text()),
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
			currentOrganizer := types.Organizer(rawCurrentOrganizer)

			if currentOrganizer == types.NAR {
				ce.ForEach(".Num > div", func(j int, ce2 *colly.HTMLElement) {
					num, _ := strconv.Atoi(ce2.DOM.Text())
					numbers = append(numbers, num)
				})
				ce.ForEach(".Odds span", func(j int, ce2 *colly.HTMLElement) {
					oddsList = append(oddsList, ce2.DOM.Text())
				})
				raceTimes = append(raceTimes, ce.DOM.Find(".Time > .RaceTime").Text())
				popularNumber, _ := strconv.Atoi(oddsList[0])
				linkUrl, _ := ce.DOM.Find(".Jockey > a").Attr("href")
				regex := regexp.MustCompile(`(\d{5})`)
				result := regex.FindStringSubmatch(linkUrl)
				// 一部の騎手で\d{5}で引っかからないjockeyIdの場合があるが、マイナーな騎手なので無視する
				jockeyId := 0
				if result != nil {
					jockeyId, _ = strconv.Atoi(result[1])
				}

				raceResults = append(raceResults, netkeiba_entity.NewRaceResult(
					i+1,
					ConvertFromEucJPToUtf8(ce.DOM.Find(".Horse_Name > a").Text()),
					numbers[0],
					numbers[1],
					jockeyId,
					oddsList[1],
					popularNumber,
				))
			}
		})
	})

	r.client.OnHTML("div.RaceList_Item02", func(e *colly.HTMLElement) {
		e.ForEach("div", func(i int, ce *colly.HTMLElement) {
			query := ce.Request.URL.Query()
			rawCurrentOrganizer, _ := strconv.Atoi(query.Get("organizer"))
			currentOrganizer := types.Organizer(rawCurrentOrganizer)
			if currentOrganizer == types.JRA || currentOrganizer == types.NAR {
				switch i {
				case 0:
					regex := regexp.MustCompile(`(.+)\s`)
					matches := regex.FindAllStringSubmatch(ce.DOM.Text(), -1)
					raceName = ConvertFromEucJPToUtf8(matches[0][1])
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
				case 1:
					text := ConvertFromEucJPToUtf8(ce.DOM.Text())
					regex := regexp.MustCompile(`(\d+\:\d+).+(ダ|芝|障)(\d+)[\s\S]+馬場:(.+)`)
					matches := regex.FindAllStringSubmatch(text, -1)
					startTime = matches[0][1]
					courseCategory = types.NewCourseCategory(matches[0][2])
					distance, _ = strconv.Atoi(matches[0][3])
					trackCondition = matches[0][4]
				case 2:
					text := ConvertFromEucJPToUtf8(ce.DOM.Text())
					regex := regexp.MustCompile(`(\d+)頭`)
					matches := regex.FindAllStringSubmatch(text, -1)
					entries, _ = strconv.Atoi(matches[0][1])
				}
			} else if currentOrganizer == types.OverseaOrganizer {
				switch i {
				case 0:
					raceName = ConvertFromEucJPToUtf8(ce.DOM.Text())
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
						text := ConvertFromEucJPToUtf8(ce2.DOM.Text())
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
								trackCondition = matches[0][1]
							}
						case 3:
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
			readPopulars := func(ce2 *colly.HTMLElement) []int {
				values := strings.Split(ConvertFromEucJPToUtf8(ce2.DOM.Text()), "人気")
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

	err := r.client.Visit(url)
	if err != nil {
		return nil, err
	}

	return netkeiba_entity.NewRace(
		raceName,
		url,
		raceTimes[0],
		startTime,
		entries,
		distance,
		int(gradeClass),
		int(courseCategory),
		trackCondition,
		raceResults,
		payoutResults,
	), nil
}
