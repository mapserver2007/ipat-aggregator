package master_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/config"
)

const horseFileName = "horse.json"

type Horse interface {
	Get(ctx context.Context) ([]*data_cache_entity.Horse, error)
	CreateOrUpdate(ctx context.Context, horses []*data_cache_entity.Horse) error
}

type horseService struct {
	horseRepository      repository.HorseRepository
	horseEntityConverter converter.HorseEntityConverter
}

func NewHorse(
	horseRepository repository.HorseRepository,
	horseEntityConverter converter.HorseEntityConverter,
) Horse {
	return &horseService{
		horseRepository:      horseRepository,
		horseEntityConverter: horseEntityConverter,
	}
}

func (h *horseService) Get(ctx context.Context) ([]*data_cache_entity.Horse, error) {
	rawHorseInfo, err := h.horseRepository.Read(ctx, fmt.Sprintf("%s/%s", config.CacheDir, horseFileName))
	if err != nil {
		return nil, err
	}

	var horses []*data_cache_entity.Horse
	if rawHorseInfo != nil {
		for _, rawHorse := range rawHorseInfo.Horses {
			horse, err := h.horseEntityConverter.RawToDataCache(rawHorse)
			if err != nil {
				return nil, err
			}
			horses = append(horses, horse)
		}
	}

	return horses, nil
}

func (h *horseService) CreateOrUpdate(
	ctx context.Context,
	horses []*data_cache_entity.Horse,
) error {
	caches, err := h.Get(ctx)
	if err != nil {
		return err
	}

	caches = append(caches, horses...)

	horseIdMap := converter.ConvertToMap(caches, func(horse *data_cache_entity.Horse) types.HorseId {
		return horse.HorseId()
	})

	newCaches := make([]*raw_entity.Horse, 0, len(caches)+len(horses))
	for _, horseId := range service.SortedHorseIdKeys(horseIdMap) {
		newCaches = append(newCaches, h.horseEntityConverter.DataCacheToRaw(horseIdMap[horseId]))
	}

	err = h.horseRepository.Write(ctx, fmt.Sprintf("%s/%s", config.CacheDir, horseFileName), &raw_entity.HorseInfo{
		Horses: newCaches,
	})
	if err != nil {
		return err
	}

	return nil
}
