package infrastructure

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"github.com/gocolly/colly"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"io"
	"net/http"
	neturl "net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type predictionDataRepository struct {
	client *colly.Collector
}

func NewPredictionDataRepository() repository.PredictionDataRepository {
	return &predictionDataRepository{
		client: colly.NewCollector(),
	}
}

func (p *predictionDataRepository) Read(ctx context.Context, filePath string) ([]*marker_csv_entity.PredictionMarker, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var markers []*marker_csv_entity.PredictionMarker
	reader := csv.NewReader(f)
	rowNum := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if rowNum == 0 {
			rowNum++
			continue
		}

		marker := marker_csv_entity.NewPredictionMarker(
			record[0],
			record[1],
			record[2],
			record[3],
			record[4],
			record[5],
			record[6],
		)
		if err != nil {
			return nil, err
		}

		markers = append(markers, marker)
		rowNum++
	}

	return markers, nil
}

func (p *predictionDataRepository) Fetch(
	ctx context.Context,
	raceUrl string,
	oddsUrl string,
) (*netkeiba_entity.Race, []*netkeiba_entity.Odds, error) {
	odds, err := p.fetchOdds(ctx, oddsUrl)
	if err != nil {
		return nil, nil, err
	}

	race, err := p.fetchRace(ctx, raceUrl)
	if err != nil {
		return nil, nil, err
	}

	return race, odds, nil
}

func (p *predictionDataRepository) fetchRace(ctx context.Context, url string) (*netkeiba_entity.Race, error) {
	var (
		raceName          string
		trackCondition    types.TrackCondition
		startTime         string
		courseCategory    types.CourseCategory
		distance, entries int
		gradeClass        types.GradeClass
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
	raceId := queryParams.Get("race_id")

	p.client.OnHTML("div.RaceList_Item02", func(e *colly.HTMLElement) {
		e.ForEach("div", func(i int, ce *colly.HTMLElement) {
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
						text := ConvertFromEucJPToUtf8(ce.DOM.Text())
						if strings.Contains(text, "牝") {
							raceSexCondition = types.FillyAndMareLimited
						}
					case 6:
						text := ConvertFromEucJPToUtf8(ce.DOM.Text())
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
						text := ConvertFromEucJPToUtf8(ce.DOM.Text())
						regex := regexp.MustCompile(`(\d+)頭`)
						matches := regex.FindAllStringSubmatch(text, -1)
						entries, _ = strconv.Atoi(matches[0][1])
					}
				})
			}
		})

	})

	err = p.client.Visit(url)
	if err != nil {
		return nil, err
	}

	return netkeiba_entity.NewRace(
		raceId,
		0,
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
		nil,
		nil,
	), nil
}

func (p *predictionDataRepository) fetchOdds(ctx context.Context, url string) ([]*netkeiba_entity.Odds, error) {
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

	var odds []*netkeiba_entity.Odds
	for rawNumber, list := range oddsInfo.Data.Odds.List {
		popularNumber, _ := strconv.Atoi(list[2])
		horseNumber, _ := strconv.Atoi(rawNumber)
		odds = append(odds, netkeiba_entity.NewOdds(
			list[0], popularNumber, horseNumber,
		))
	}

	return odds, nil
}
