package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	betting_ticket_entity "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/entity"
	race_entity "github.com/mapserver2007/tools/baken/app/domain/race/entity"
	race_vo "github.com/mapserver2007/tools/baken/app/domain/race/value_object"
	"github.com/mapserver2007/tools/baken/app/repository"
	"log"
	"os"
	"path/filepath"
	"time"
)

type RaceDB struct {
	raceClient repository.RaceClient
}

func NewRaceDB(
	raceClient repository.RaceClient,
) repository.RaceDB {
	return &RaceDB{
		raceClient: raceClient,
	}
}

func (r *RaceDB) ReadRaceResult(ctx context.Context, fileName string) (*race_entity.RaceInfo, error) {
	bytes, err := r.readFile(fileName)
	if err != nil {
		return nil, err
	}

	var raceInfo *race_entity.RaceInfo
	if err := json.Unmarshal(bytes, &raceInfo); err != nil {
		return nil, err
	}

	return raceInfo, nil
}

func (r *RaceDB) ReadRacingNumber(ctx context.Context, fileName string) (*race_entity.RacingNumberInfo, error) {
	bytes, err := r.readFile(fileName)
	if err != nil {
		return nil, err
	}

	var racingNumberInfo *race_entity.RacingNumberInfo
	if err := json.Unmarshal(bytes, &racingNumberInfo); err != nil {
		return nil, err
	}

	return racingNumberInfo, nil
}

// JRAの場合はDateからIDが特定できないので開催場所、日のデータをキャッシュしておいて
// 変換処理をする必要がある。NAR、海外はDateから特定可能
func (r *RaceDB) UpdateRacingNumber(ctx context.Context, fileName string, entities []*betting_ticket_entity.CsvEntity) error {
	// エラーだった場合はracing_number.jsonが空だった場合なので無視
	currentRacingNumber, _ := r.ReadRacingNumber(ctx, fileName)
	cacheKeyFunc := func(date, raceCourseName int) string {
		return fmt.Sprintf("%d_%d", date, raceCourseName)
	}

	cache := map[string]race_entity.RacingNumber{}
	if currentRacingNumber != nil {
		for _, racingNumber := range currentRacingNumber.RacingNumbers {
			key := cacheKeyFunc(racingNumber.Date, racingNumber.RaceCourseId)
			if _, ok := cache[key]; !ok {
				cache[key] = *racingNumber
			}
		}
	}

	for _, entity := range entities {
		if entity.RaceCourse.Organizer() != race_vo.JRA {
			continue
		}
		key := cacheKeyFunc(int(entity.RaceDate), entity.RaceCourse.Value())
		if _, ok := cache[key]; !ok {
			time.Sleep(time.Second * 1)
			newRacingNumbers, err := r.raceClient.GetRacingNumber(ctx, entity)
			if err != nil {
				return err
			}

			for i := 0; i < len(newRacingNumbers); i++ {
				newRacingNumber := newRacingNumbers[i]
				key = cacheKeyFunc(newRacingNumber.Date, newRacingNumber.RaceCourseId)
				cache[key] = *newRacingNumbers[i]
				log.Printf("updating key: %s ...", key)
			}
		}
	}

	var racingNumbers []*race_entity.RacingNumber
	for key := range cache {
		racingNumber, _ := cache[key]
		racingNumbers = append(racingNumbers, &racingNumber)
	}
	racingNumberInfo := race_entity.RacingNumberInfo{RacingNumbers: racingNumbers}

	bytes, err := json.Marshal(racingNumberInfo)
	if err != nil {
		return err
	}

	err = r.writeFile(fileName, bytes)
	if err != nil {
		return err
	}

	return nil
}

func (r *RaceDB) UpdateRaceResult(ctx context.Context, fileName string, racingNumbers []*race_entity.RacingNumber, entities []*betting_ticket_entity.CsvEntity) error {
	raceNumberMap := map[string]*race_entity.RacingNumber{}
	for _, racingNumber := range racingNumbers {
		key := fmt.Sprintf("%d_%d", racingNumber.Date, racingNumber.RaceCourseId)
		raceNumberMap[key] = racingNumber
	}

	// エラーだった場合はracing_number.jsonが空だった場合なので無視
	raceInfo, _ := r.ReadRaceResult(ctx, fileName)

	cache := map[race_vo.RaceId]*race_entity.Race{}
	if raceInfo != nil {
		for _, race := range raceInfo.Races {
			cache[race_vo.RaceId(race.RaceId)] = race
		}
	}
	for _, entity := range entities {
		var raceId race_vo.RaceId
		organizer := entity.RaceCourse.Organizer()

		switch organizer {
		case race_vo.JRA:
			key := fmt.Sprintf("%d_%d", entity.RaceDate, entity.RaceCourse.Value())
			racingNumber, ok := raceNumberMap[key]
			if !ok {
				return fmt.Errorf("undefined key: %s", key)
			}
			rawRaceId := fmt.Sprintf("%d%02d%02d%02d%02d", entity.RaceDate.Year(), racingNumber.RaceCourseId, racingNumber.Round, racingNumber.Day, entity.RaceNo)
			raceId = race_vo.RaceId(rawRaceId)
		case race_vo.NAR:
			rawRaceId := fmt.Sprintf("%d%02d%02d%02d%02d", entity.RaceDate.Year(), entity.RaceCourse.Value(), entity.RaceDate.Month(), entity.RaceDate.Day(), entity.RaceNo)
			raceId = race_vo.RaceId(rawRaceId)
		case race_vo.OverseaOrganizer:
			raceCourseIdForOversea := race_vo.ConvertToOverseaRaceCourseId(entity.RaceCourse)
			rawRaceId := fmt.Sprintf("%d%s%02d%02d%02d", entity.RaceDate.Year(), raceCourseIdForOversea, entity.RaceDate.Month(), entity.RaceDate.Day(), entity.RaceNo)
			raceId = race_vo.RaceId(rawRaceId)
		}

		if _, ok := cache[raceId]; ok {
			continue
		}

		time.Sleep(time.Second * 1)
		race := r.raceClient.GetRaceResult(ctx, raceId, organizer)
		race.RaceNumber = entity.RaceNo
		race.RaceDate = int(entity.RaceDate)
		race.RaceCourseId = entity.RaceCourse.Value()

		log.Printf("updating key: %s ...", raceId)

		cache[raceId] = race
	}

	var races []*race_entity.Race
	for _, race := range cache {
		races = append(races, race)
	}

	newRaceInfo := race_entity.RaceInfo{Races: races}

	bytes, err := json.Marshal(newRaceInfo)
	if err != nil {
		return err
	}

	err = r.writeFile(fileName, bytes)
	if err != nil {
		return err
	}

	return nil
}

func (r *RaceDB) readFile(fileName string) ([]byte, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	filePath, err := filepath.Abs(fmt.Sprintf("%s/cache/%s", rootPath, fileName))
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (r *RaceDB) writeFile(fileName string, data []byte) error {
	rootPath, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath, err := filepath.Abs(fmt.Sprintf("%s/cache/%s", rootPath, fileName))
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
