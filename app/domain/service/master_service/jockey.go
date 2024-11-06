package master_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/config"
	"sort"
	"time"
)

const (
	jockeyUrl               = "https://db.netkeiba.com/jockey/%s/"
	jockeyFileName          = "jockey.json"
	horseFileName           = "horse.json"
	beginIdForJRA           = 422
	endIdForJRA             = 2000
	beginIdForNARandOversea = 5000
	endIdForNARandOversea   = 5999
)

type Jockey interface {
	Get(ctx context.Context) ([]*data_cache_entity.Jockey, []types.JockeyId, error)
	CreateOrUpdate(ctx context.Context, jockeys []*data_cache_entity.Jockey, excludeJockeyIds []types.JockeyId) error
}

type jockeyService struct {
	jockeyRepository      repository.JockeyRepository
	jockeyEntityConverter converter.JockeyEntityConverter
}

func NewJockey(
	jockeyRepository repository.JockeyRepository,
	jockeyEntityConverter converter.JockeyEntityConverter,
) Jockey {
	return &jockeyService{
		jockeyRepository:      jockeyRepository,
		jockeyEntityConverter: jockeyEntityConverter,
	}
}

func (j *jockeyService) Get(ctx context.Context) ([]*data_cache_entity.Jockey, []types.JockeyId, error) {
	rawJockeyInfo, err := j.jockeyRepository.Read(ctx, fmt.Sprintf("%s/%s", config.CacheDir, jockeyFileName))
	if err != nil {
		return nil, nil, err
	}

	var (
		jockeys          []*data_cache_entity.Jockey
		excludeJockeyIds []types.JockeyId
	)
	if rawJockeyInfo != nil {
		for _, rawJockey := range rawJockeyInfo.Jockeys {
			jockeys = append(jockeys, j.jockeyEntityConverter.RawToDataCache(rawJockey))
		}
		for _, rawExcludeJockeyId := range rawJockeyInfo.ExcludeJockeyIds {
			excludeJockeyIds = append(excludeJockeyIds, types.JockeyId(rawExcludeJockeyId))
		}
	}

	return jockeys, excludeJockeyIds, nil
}

func (j *jockeyService) CreateOrUpdate(
	ctx context.Context,
	jockeys []*data_cache_entity.Jockey,
	excludeJockeyIds []types.JockeyId,
) error {
	urls := j.createJockeyUrls(jockeys, excludeJockeyIds)
	if len(urls) == 0 {
		return nil
	}

	var (
		rawJockeys          []*raw_entity.Jockey
		rawExcludeJockeyIds []string
	)

	for _, url := range urls {
		time.Sleep(time.Millisecond)
		jockey, err := j.jockeyRepository.Fetch(ctx, url)
		if err != nil {
			return err
		}
		if jockey.Name() == "" {
			rawExcludeJockeyIds = append(rawExcludeJockeyIds, jockey.Id())
		} else {
			rawJockeys = append(rawJockeys, j.jockeyEntityConverter.NetKeibaToRaw(jockey))
		}
	}

	for _, jockey := range jockeys {
		rawJockeys = append(rawJockeys, j.jockeyEntityConverter.DataCacheToRaw(jockey))
	}

	sort.Slice(rawJockeys, func(i, j int) bool {
		return rawJockeys[i].JockeyId < rawJockeys[j].JockeyId
	})

	for _, excludeJockeyId := range excludeJockeyIds {
		rawExcludeJockeyIds = append(rawExcludeJockeyIds, excludeJockeyId.Value())
	}

	sort.Strings(rawExcludeJockeyIds)

	err := j.jockeyRepository.Write(ctx, fmt.Sprintf("%s/%s", config.CacheDir, jockeyFileName), &raw_entity.JockeyInfo{
		Jockeys:          rawJockeys,
		ExcludeJockeyIds: rawExcludeJockeyIds,
	})
	if err != nil {
		return err
	}

	return nil
}

func (j *jockeyService) createJockeyUrls(
	jockeys []*data_cache_entity.Jockey,
	excludeJockeyIds []types.JockeyId,
) []string {
	jockeysMap := map[string]bool{}
	for _, jockeyData := range jockeys {
		jockeysMap[jockeyData.JockeyId().Value()] = true
	}

	excludeJockeyIdsMap := map[string]bool{}
	for _, jockeyId := range excludeJockeyIds {
		excludeJockeyIdsMap[jockeyId.Value()] = true
	}

	var urls []string
	for i := beginIdForJRA; i <= endIdForJRA; i++ {
		id := fmt.Sprintf("%05d", i)
		// 除外リストに含まれてたら何もしない
		if _, ok := excludeJockeyIdsMap[id]; ok {
			continue
		}
		// 既に取得済みの場合は何もしない
		if _, ok := jockeysMap[id]; ok {
			continue
		}
		jockeyId := types.JockeyId(id)
		urls = append(urls, fmt.Sprintf(jockeyUrl, jockeyId.Value()))
	}
	for i := beginIdForNARandOversea; i <= endIdForNARandOversea; i++ {
		id := fmt.Sprintf("%05d", i)
		// 除外リストに含まれてたら何もしない
		if _, ok := excludeJockeyIdsMap[id]; ok {
			continue
		}
		// 既に取得済みの場合は何もしない
		if _, ok := jockeysMap[id]; ok {
			continue
		}
		jockeyId := types.JockeyId(id)
		urls = append(urls, fmt.Sprintf(jockeyUrl, jockeyId.Value()))
	}
	// 特殊IDの騎手を追加
	otherJockeyIds := []string{
		"a02d7", // 西啓太
		"a050d", // 宮内勇樹
	}
	for _, jockeyId := range otherJockeyIds {
		// 除外リストに含まれてたら何もしない
		if _, ok := excludeJockeyIdsMap[jockeyId]; ok {
			continue
		}
		// 既に取得済みの場合は何もしない
		if _, ok := jockeysMap[jockeyId]; ok {
			continue
		}
		urls = append(urls, fmt.Sprintf(jockeyUrl, jockeyId))
	}

	return urls
}
