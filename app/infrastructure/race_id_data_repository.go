package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"os"
	"path/filepath"
	"regexp"
)

type raceIdDataRepository struct {
	client *colly.Collector
}

func NewRaceIdDataRepository() repository.RaceIdDataRepository {
	return &raceIdDataRepository{
		client: colly.NewCollector(),
	}
}

func (r *raceIdDataRepository) Read(ctx context.Context, fileName string) ([]*raw_entity.RaceDate, []int, error) {
	raceDates := make([]*raw_entity.RaceDate, 0)
	excludeDates := make([]int, 0)
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}

	path, err := filepath.Abs(fmt.Sprintf("%s/cache/%s", rootPath, fileName))
	if err != nil {
		return nil, nil, err
	}

	// ファイルが存在しない場合はエラーは返さず処理を継続する
	bytes, err := os.ReadFile(path)
	if err != nil {
		return raceDates, excludeDates, nil
	}

	var raceIdInfo *raw_entity.RaceIdInfo
	if err := json.Unmarshal(bytes, &raceIdInfo); err != nil {
		return nil, nil, err
	}
	raceDates = raceIdInfo.RaceDates
	excludeDates = raceIdInfo.ExcludeDates

	return raceDates, excludeDates, nil
}

func (r *raceIdDataRepository) Write(ctx context.Context, fileName string, raceIdInfo *raw_entity.RaceIdInfo) error {
	bytes, err := json.Marshal(raceIdInfo)
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

func (r *raceIdDataRepository) Fetch(ctx context.Context, url string) ([]string, error) {
	var rawRaceIds []string
	r.client.OnHTML(".RaceList_DataItem > a:first-child", func(e *colly.HTMLElement) {
		regex := regexp.MustCompile(`race_id=(\d+)`)
		matches := regex.FindAllStringSubmatch(e.Attr("href"), -1)
		raceId := matches[0][1]
		rawRaceIds = append(rawRaceIds, raceId)
	})
	err := r.client.Visit(url)
	if err != nil {
		return nil, err
	}

	return rawRaceIds, nil
}
