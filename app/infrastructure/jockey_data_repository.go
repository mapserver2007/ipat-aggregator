package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type jockeyDataRepository struct {
	client *colly.Collector
}

func NewJockeyDataRepository() repository.JockeyDataRepository {
	return &jockeyDataRepository{
		client: colly.NewCollector(),
	}
}

func (j *jockeyDataRepository) Read(ctx context.Context, fileName string) ([]*raw_entity.Jockey, []int, error) {
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
		return make([]*raw_entity.Jockey, 0), make([]int, 0), nil
	}

	var jockeyInfo *raw_entity.JockeyInfo
	if err := json.Unmarshal(bytes, &jockeyInfo); err != nil {
		return nil, nil, err
	}

	return jockeyInfo.Jockeys, jockeyInfo.ExcludeJockeyIds, nil
}

func (j *jockeyDataRepository) Write(
	ctx context.Context,
	fileName string,
	jockeyInfo *raw_entity.JockeyInfo,
) error {
	bytes, err := json.Marshal(jockeyInfo)
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

func (j *jockeyDataRepository) Fetch(
	ctx context.Context,
	url string,
) (*netkeiba_entity.Jockey, error) {
	var name string
	j.client.OnHTML("div.Name h1", func(e *colly.HTMLElement) {
		list := strings.Split(e.DOM.Text(), "\n")
		name = ConvertFromEucJPToUtf8(list[1][:len(list[1])-2])
	})
	j.client.OnError(func(r *colly.Response, err error) {
		log.Printf("GetJockey error: %v", err)
	})

	regex := regexp.MustCompile(`\/jockey\/(\d+)\/`)
	result := regex.FindStringSubmatch(url)
	id, _ := strconv.Atoi(result[1])

	err := j.client.Visit(url)
	if err != nil {
		return nil, err
	}

	return netkeiba_entity.NewJockey(id, name), nil
}
