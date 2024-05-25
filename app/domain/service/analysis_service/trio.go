package analysis_service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/filter_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type Trio interface {
	// Create15 3連複 軸1頭相手5頭流し 10点買い分析
	Create15(ctx context.Context, markers []*marker_csv_entity.AnalysisMarker, races []*data_cache_entity.Race) error
	// Create24 3連複 軸2頭相手4頭流し 4点買い分析
	Create24(ctx context.Context, markers []*marker_csv_entity.AnalysisMarker, races []*data_cache_entity.Race) error
}

type trioService struct {
	filterService filter_service.AnalysisFilter
}

func (t *trioService) Create24(ctx context.Context, markers []*marker_csv_entity.AnalysisMarker, races []*data_cache_entity.Race) error {
	//TODO implement me
	panic("implement me")
}

func NewTrio(
	filterService filter_service.AnalysisFilter,
) Trio {
	return &trioService{
		filterService: filterService,
	}
}

func (t *trioService) Create15(
	ctx context.Context,
	markers []*marker_csv_entity.AnalysisMarker,
	races []*data_cache_entity.Race,
) error {
	markerMap := converter.ConvertToMap(markers, func(marker *marker_csv_entity.AnalysisMarker) types.RaceId {
		return marker.RaceId()
	})

	for _, race := range races {
		raceResultMap := converter.ConvertToMap(race.RaceResults(), func(raceResult *data_cache_entity.RaceResult) int {
			return raceResult.HorseNumber()
		})

		marker, ok := markerMap[race.RaceId()]
		if !ok {
			switch race.Class() {
			case types.MakeDebut, types.JumpMaiden, types.JumpGrade1, types.JumpGrade2, types.JumpGrade3:
				// 新馬・障害は分析印なしなのでスキップ
			default:
				// 印が不完全な場合がたまにあり(同じ印がついていたり、取り消しによる印6個未満の場合)、その場合はスキップ
				// log.Println(fmt.Sprintf("raceId not found in place markers: %s", race.RaceId()))
			}
			continue
		}

		filters := t.filterService.Create(ctx, race)

		_ = raceResultMap
		_ = marker
		_ = filters

	}

	return nil
}
