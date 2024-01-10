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
	neturl "net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type racingNumberDataRepository struct {
	client *colly.Collector
}

func NewRacingNumberDataRepository() repository.RacingNumberDataRepository {
	return &racingNumberDataRepository{
		client: colly.NewCollector(),
	}
}

func (r *racingNumberDataRepository) Read(ctx context.Context, fileName string) ([]*raw_entity.RacingNumber, error) {
	racingNumbers := make([]*raw_entity.RacingNumber, 0)
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
		return racingNumbers, nil
	}

	var racingNumberInfo *raw_entity.RacingNumberInfo
	if err := json.Unmarshal(bytes, &racingNumberInfo); err != nil {
		return nil, err
	}
	racingNumbers = racingNumberInfo.RacingNumbers

	return racingNumbers, nil
}

// Write JRAの場合はDateからIDが特定できないので開催場所、日のデータをキャッシュしておいて
// 変換処理をする必要がある。NAR、海外はDateから特定可能
func (r *racingNumberDataRepository) Write(
	ctx context.Context,
	fileName string,
	racingNumberInfo *raw_entity.RacingNumberInfo,
) error {
	bytes, err := json.Marshal(racingNumberInfo)
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

func (r *racingNumberDataRepository) Fetch(
	ctx context.Context,
	url string,
) ([]*netkeiba_entity.RacingNumber, error) {
	var racingNumbers []*netkeiba_entity.RacingNumber
	r.client.OnHTML(".RaceList_DataList", func(e *colly.HTMLElement) {
		e.ForEach(".RaceList_DataTitle", func(i int, ce *colly.HTMLElement) {
			regex := regexp.MustCompile(`(\d+)回\s+(.+)\s+(\d+)日目`)
			matches := regex.FindAllStringSubmatch(ce.Text, -1)
			round, _ := strconv.Atoi(matches[0][1])
			day, _ := strconv.Atoi(matches[0][3])
			raceCourse := types.NewRaceCourse(matches[0][2])
			u, _ := neturl.Parse(url)
			query := u.Query()
			raceDate, _ := strconv.Atoi(query.Get("kaisai_date"))

			racingNumbers = append(racingNumbers, netkeiba_entity.NewRacingNumber(
				raceDate,
				round,
				day,
				raceCourse.Value(),
			))
		})
	})
	err := r.client.Visit(url)
	if err != nil {
		return nil, err
	}

	return racingNumbers, nil
}
