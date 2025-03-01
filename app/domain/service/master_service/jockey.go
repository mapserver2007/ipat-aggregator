package master_service

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/config"
	"github.com/sirupsen/logrus"
)

const (
	jockeyUrl               = "https://db.netkeiba.com/jockey/%s/"
	jockeyFileName          = "jockey.json"
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
	logger                *logrus.Logger
}

func NewJockey(
	jockeyRepository repository.JockeyRepository,
	jockeyEntityConverter converter.JockeyEntityConverter,
	logger *logrus.Logger,
) Jockey {
	return &jockeyService{
		jockeyRepository:      jockeyRepository,
		jockeyEntityConverter: jockeyEntityConverter,
		logger:                logger,
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
	taskCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	urls := j.createJockeyUrls(jockeys, excludeJockeyIds)
	if len(urls) == 0 {
		return nil
	}

	var (
		rawJockeys          []*raw_entity.Jockey
		rawExcludeJockeyIds []string
	)

	var wg sync.WaitGroup
	const raceIdParallel = 10
	errorCh := make(chan error, 1)
	resultCh := make(chan []*netkeiba_entity.Jockey, raceIdParallel)
	chunkSize := (len(urls) + raceIdParallel - 1) / raceIdParallel

	for i := 0; i < len(urls); i += chunkSize {
		end := i + chunkSize
		if end > len(urls) {
			end = len(urls)
		}

		wg.Add(1)
		go func(splitUrls []string) {
			defer wg.Done()
			localJockeys := make([]*netkeiba_entity.Jockey, 0, len(splitUrls))
			j.logger.Infof("jockey fetch processing: %v/%v", end, len(urls))
			for _, url := range splitUrls {
				time.Sleep(time.Millisecond)
				select {
				case <-taskCtx.Done():
					return
				default:
					jockey, err := j.jockeyRepository.Fetch(taskCtx, url)
					if err != nil {
						select {
						case errorCh <- err:
							cancel()
						}
						return
					}
					localJockeys = append(localJockeys, jockey)
				}
			}

			resultCh <- localJockeys
		}(urls[i:end])
	}

	wg.Wait()
	close(errorCh)
	close(resultCh)

	if err := <-errorCh; err != nil {
		return err
	}

	for results := range resultCh {
		for _, jockey := range results {
			if jockey.Name() == "" {
				rawExcludeJockeyIds = append(rawExcludeJockeyIds, jockey.Id())
			} else {
				rawJockeys = append(rawJockeys, j.jockeyEntityConverter.NetKeibaToRaw(jockey))
			}
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
