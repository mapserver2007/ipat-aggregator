package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type oddsDataRepository struct{}

func NewOddsDataRepository() repository.OddsDataRepository {
	return &oddsDataRepository{}
}

func (o *oddsDataRepository) Read(ctx context.Context, filePath string) ([]*raw_entity.RaceOdds, error) {
	raceOdds := make([]*raw_entity.RaceOdds, 0)
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path, err := filepath.Abs(fmt.Sprintf("%s/cache/%s", rootPath, filePath))
	if err != nil {
		return nil, err
	}

	// ファイルが存在しない場合はエラーは返さず処理を継続する
	bytes, err := os.ReadFile(path)
	if err != nil {
		return raceOdds, nil
	}

	var raceOddsInfo *raw_entity.RaceOddsInfo
	if err := json.Unmarshal(bytes, &raceOddsInfo); err != nil {
		return nil, err
	}
	raceOdds = raceOddsInfo.RaceOdds

	return raceOdds, nil
}

func (o *oddsDataRepository) Write(ctx context.Context, filePath string, oddsInfo *raw_entity.RaceOddsInfo) error {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(oddsInfo)
	if err != nil {
		return err
	}

	rootPath, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath, err = filepath.Abs(fmt.Sprintf("%s/cache/%s", rootPath, filePath))
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, buffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (o *oddsDataRepository) Fetch(ctx context.Context, url string) ([]*netkeiba_entity.Odds, error) {
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
		horseNumber1, _ := strconv.Atoi(rawNumber[0:2])
		horseNumber2, _ := strconv.Atoi(rawNumber[2:4])
		horseNumber3, _ := strconv.Atoi(rawNumber[4:6])
		popularNumber, _ := strconv.Atoi(list[2])
		odds = append(odds, netkeiba_entity.NewOdds(
			types.Trio, list[0], popularNumber, []int{horseNumber1, horseNumber2, horseNumber3}, raceDate,
		))
	}

	return odds, nil
}
