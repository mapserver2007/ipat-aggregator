package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"os"
	"path/filepath"
)

type racingNumberDataRepository struct {
	raceConverter service.RaceConverter
}

func NewRacingNumberDataRepository(
	raceConverter service.RaceConverter,
) repository.RacingNumberDataRepository {
	return &racingNumberDataRepository{
		raceConverter: raceConverter,
	}
}

func (r *racingNumberDataRepository) Read(ctx context.Context, fileName string) ([]*raw_entity.RacingNumber, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path, err := filepath.Abs(fmt.Sprintf("%s/cache/%s", rootPath, fileName))
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var racingNumberInfo *raw_entity.RacingNumberInfo
	if err := json.Unmarshal(bytes, &racingNumberInfo); err != nil {
		return nil, err
	}

	return racingNumberInfo.RacingNumbers, nil
}

func (r *racingNumberDataRepository) Write(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (r *racingNumberDataRepository) Fetch(
	ctx context.Context,
	racingNumbers []*raw_entity.RacingNumber,
	tickets []*ticket_csv_entity.Ticket,
) error {

	racingNumberMap := r.raceConverter.RawRacingNumberToRawRacingNumberMap(ctx, racingNumbers)

	fmt.Println(racingNumberMap)

	return nil
}
