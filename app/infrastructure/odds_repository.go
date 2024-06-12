package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure/gateway"
	neturl "net/url"
	"os"
	"path/filepath"
)

type oddsRepository struct {
	netKeibaGateway gateway.NetKeibaGateway
}

func NewOddsRepository(
	netKeibaGateway gateway.NetKeibaGateway,
) repository.OddsRepository {
	return &oddsRepository{
		netKeibaGateway: netKeibaGateway,
	}
}

func (o *oddsRepository) List(
	ctx context.Context,
	path string,
) ([]string, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(fmt.Sprintf("%s/%s", rootPath, path))
	if err != nil {
		return nil, err
	}

	pattern := filepath.Join(absPath, "*.json")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		fileNames = append(fileNames, filepath.Base(file))
	}

	return fileNames, nil
}

func (o *oddsRepository) Read(
	ctx context.Context,
	path string,
) ([]*raw_entity.RaceOdds, error) {
	raceOdds := make([]*raw_entity.RaceOdds, 0)
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(fmt.Sprintf("%s/%s", rootPath, path))
	if err != nil {
		return nil, err
	}

	// ファイルが存在しない場合はエラーは返さず処理を継続する
	bytes, err := os.ReadFile(absPath)
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

func (o *oddsRepository) Write(
	ctx context.Context,
	path string,
	data *raw_entity.RaceOddsInfo,
) error {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(data)
	if err != nil {
		return err
	}

	rootPath, err := os.Getwd()
	if err != nil {
		return err
	}

	absPath, err := filepath.Abs(fmt.Sprintf("%s/%s", rootPath, path))
	if err != nil {
		return err
	}

	err = os.WriteFile(absPath, buffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (o *oddsRepository) Fetch(
	ctx context.Context,
	url string,
) ([]*netkeiba_entity.Odds, error) {
	queryParams, err := neturl.ParseQuery(url)
	if err != nil {
		return nil, err
	}
	ticketType := queryParams.Get("type")

	var odds []*netkeiba_entity.Odds
	switch ticketType {
	case "1": // win
		odds, err = o.netKeibaGateway.FetchWinOdds(ctx, url)
	case "2": // place
		odds, err = o.netKeibaGateway.FetchPlaceOdds(ctx, url)
	case "7": // trio
		odds, err = o.netKeibaGateway.FetchTrioOdds(ctx, url)
	}
	if err != nil {
		return nil, err
	}
	return odds, nil
}
