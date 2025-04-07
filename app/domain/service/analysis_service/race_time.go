package analysis_service

import (
	"context"
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/filter_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

const (
	timeFormat    = "0:00.0"
	rapTimeFormat = "0.0"
)

type RaceTime interface {
	Create(ctx context.Context,
		races []*data_cache_entity.Race,
		raceTimes []*data_cache_entity.RaceTime,
	) ([]*analysis_entity.RaceTimeCalculable, error)
	Convert(ctx context.Context,
		calculables []*analysis_entity.RaceTimeCalculable,
	) (map[filter.AttributeId]*spreadsheet_entity.AnalysisRaceTime, []filter.AttributeId, []filter.AttributeId)
	Write(ctx context.Context,
		analysisRaceTimeMap map[filter.AttributeId]*spreadsheet_entity.AnalysisRaceTime,
		attributeFilters []filter.AttributeId,
		conditionFilters []filter.AttributeId,
	) error
}

type raceTimeService struct {
	filterService         filter_service.AnalysisFilter
	spreadSheetRepository repository.SpreadSheetRepository
}

func NewRaceTime(
	filterService filter_service.AnalysisFilter,
	spreadSheetRepository repository.SpreadSheetRepository,
) RaceTime {
	return &raceTimeService{
		filterService:         filterService,
		spreadSheetRepository: spreadSheetRepository,
	}
}

func (r *raceTimeService) Create(
	ctx context.Context,
	races []*data_cache_entity.Race,
	raceTimes []*data_cache_entity.RaceTime,
) ([]*analysis_entity.RaceTimeCalculable, error) {
	raceMap := make(map[types.RaceId]*data_cache_entity.Race)

	for _, race := range races {
		switch race.Class() {
		case types.JumpMaiden, types.JumpGrade1, types.JumpGrade2, types.JumpGrade3, types.JumpOpenClass:
			continue
		default:
			raceMap[race.RaceId()] = race
		}
	}

	raceCalculables := make([]*analysis_entity.RaceTimeCalculable, 0, len(raceTimes))
	for _, raceTime := range raceTimes {
		if race, ok := raceMap[raceTime.RaceId()]; ok {
			attributeFilterIds := r.filterService.CreateRaceTimeFilters(ctx, race)
			raceCalculable, err := analysis_entity.NewRaceTimeCalculable(
				race.RaceId(),
				race.RaceDate(),
				raceTime.Time(),
				raceTime.TimeIndex(),
				raceTime.TrackIndex(),
				raceTime.RapTimes(),
				raceTime.First3f(),
				raceTime.First4f(),
				raceTime.Last3f(),
				raceTime.Last4f(),
				raceTime.Rap5f(),
				attributeFilterIds,
			)
			if err != nil {
				return nil, err
			}
			raceCalculables = append(raceCalculables, raceCalculable)
		}
	}

	return raceCalculables, nil
}

func (r *raceTimeService) Convert(
	ctx context.Context,
	calculables []*analysis_entity.RaceTimeCalculable,
) (map[filter.AttributeId]*spreadsheet_entity.AnalysisRaceTime, []filter.AttributeId, []filter.AttributeId) {
	filterRaceTimeMap := make(map[filter.AttributeId][]*analysis_entity.RaceTimeCalculable)
	attributeFilters := r.getAttributeFilters()
	conditionFilters := r.getConditionFilters()
	for _, attributeFilter := range attributeFilters {
		for _, calculable := range calculables {
			var calcFilter filter.AttributeId
			for _, f := range calculable.AttributeFilterIds() {
				calcFilter |= f
			}
			if attributeFilter&calcFilter == attributeFilter {
				filterRaceTimeMap[attributeFilter] = append(filterRaceTimeMap[attributeFilter], calculable)
			}
		}
	}

	analysisRaceTimeMap := make(map[filter.AttributeId]*spreadsheet_entity.AnalysisRaceTime)
	for attributeFilter, calculables := range filterRaceTimeMap {
		metrics := []struct {
			time func(*analysis_entity.RaceTimeCalculable) time.Duration
		}{
			{func(c *analysis_entity.RaceTimeCalculable) time.Duration { return c.Time() }},
			{func(c *analysis_entity.RaceTimeCalculable) time.Duration { return c.First3f() }},
			{func(c *analysis_entity.RaceTimeCalculable) time.Duration { return c.First4f() }},
			{func(c *analysis_entity.RaceTimeCalculable) time.Duration { return c.Last3f() }},
			{func(c *analysis_entity.RaceTimeCalculable) time.Duration { return c.Last4f() }},
			{func(c *analysis_entity.RaceTimeCalculable) time.Duration { return c.Rap5f() }},
		}

		times := make([]string, 0, len(metrics)*2)
		for i, metric := range metrics {
			format := rapTimeFormat
			if i == 0 {
				format = timeFormat
			}
			times = append(times, r.calcAverageTime(calculables, metric.time, format))
			times = append(times, r.calcMedianTime(calculables, metric.time, format))
		}

		trackIndices := make([]int, 0, len(calculables))
		timeIndices := make([]int, 0, len(calculables))
		for _, calculable := range calculables {
			trackIndices = append(trackIndices, calculable.TrackIndex())
			timeIndices = append(timeIndices, calculable.TimeIndex())
		}

		averageTrackIndex := r.calcAverage(trackIndices)
		maxTrackIndex := r.calcMax(trackIndices)
		minTrackIndex := r.calcMin(trackIndices)
		averageTimeIndex := r.calcAverage(timeIndices)

		analysisRaceTimeMap[attributeFilter] = spreadsheet_entity.NewAnalysisRaceTime(
			times[0],
			times[1],
			times[2],
			times[3],
			times[4],
			times[5],
			times[6],
			times[7],
			times[8],
			times[9],
			times[10],
			times[11],
			averageTrackIndex,
			maxTrackIndex,
			minTrackIndex,
			averageTimeIndex,
			len(calculables),
		)
	}

	return analysisRaceTimeMap, attributeFilters, conditionFilters
}

func (r *raceTimeService) Write(
	ctx context.Context,
	analysisRaceTimeMap map[filter.AttributeId]*spreadsheet_entity.AnalysisRaceTime,
	attributeFilters []filter.AttributeId,
	conditionFilters []filter.AttributeId,
) error {
	return r.spreadSheetRepository.WriteAnalysisRaceTime(ctx, analysisRaceTimeMap, attributeFilters, conditionFilters)
}

func (r *raceTimeService) calcAverageTime(
	calculables []*analysis_entity.RaceTimeCalculable,
	getter func(*analysis_entity.RaceTimeCalculable) time.Duration,
	format string,
) string {
	if len(calculables) == 0 {
		return format
	}

	var totalTime time.Duration
	for _, calculable := range calculables {
		totalTime += getter(calculable)
	}
	averageTime := totalTime / time.Duration(len(calculables))

	totalSeconds := averageTime.Seconds()
	minutes := int(totalSeconds) / 60
	seconds := totalSeconds - float64(minutes*60)

	if format == timeFormat {
		return fmt.Sprintf("%d:%04.1f", minutes, seconds)
	}

	return fmt.Sprintf("%.1f", totalSeconds)
}

func (r *raceTimeService) calcMedianTime(
	calculables []*analysis_entity.RaceTimeCalculable,
	getter func(*analysis_entity.RaceTimeCalculable) time.Duration,
	format string,
) string {
	if len(calculables) == 0 {
		return format
	}

	durations := make([]time.Duration, 0, len(calculables))
	for _, c := range calculables {
		durations = append(durations, getter(c))
	}

	slices.Sort(durations)

	var median time.Duration
	n := len(durations)
	if n%2 == 1 {
		median = durations[n/2]
	} else {
		median = (durations[n/2-1] + durations[n/2]) / 2
	}

	totalSeconds := median.Seconds()
	minutes := int(totalSeconds) / 60
	seconds := totalSeconds - float64(minutes*60)

	if format == timeFormat {
		return fmt.Sprintf("%d:%04.1f", minutes, seconds)
	}

	return fmt.Sprintf("%.1f", totalSeconds)
}

func (r *raceTimeService) calcAverage(values []int) int {
	if len(values) == 0 {
		return 0
	}

	sum := 0
	for _, v := range values {
		sum += v
	}
	return int(math.Round(float64(sum) / float64(len(values))))
}

func (r *raceTimeService) calcMax(values []int) int {
	if len(values) == 0 {
		return 0
	}

	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func (r *raceTimeService) calcMin(values []int) int {
	if len(values) == 0 {
		return 0
	}

	min := values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

func (r *raceTimeService) getConditionFilters() []filter.AttributeId {
	return []filter.AttributeId{
		filter.Tokyo | filter.Nakayama | filter.Kyoto | filter.Hanshin | filter.Niigata | filter.Chukyo | filter.Sapporo | filter.Hakodate | filter.Fukushima | filter.Kokura,
		filter.Turf | filter.Dirt,
		filter.Distance1000m | filter.Distance1150m | filter.Distance1200m | filter.Distance1300m | filter.Distance1400m | filter.Distance1500m | filter.Distance1600m | filter.Distance1700m | filter.Distance1800m | filter.Distance1900m | filter.Distance2000m | filter.Distance2100m | filter.Distance2200m | filter.Distance2300m | filter.Distance2400m | filter.Distance2500m | filter.Distance2600m | filter.Distance3000m | filter.Distance3200m | filter.Distance3400m | filter.Distance3600m,
		filter.GoodToFirm | filter.Good | filter.Yielding | filter.Soft,
		filter.Maiden | filter.OneWinClass | filter.TwoWinClass | filter.ThreeWinClass | filter.OpenListedClass | filter.Grade3 | filter.Grade2 | filter.Grade1,
		filter.TwoYearsOld | filter.ThreeYearsOld | filter.ThreeYearsAndOlder | filter.FourYearsAndOlder,
	}
}

func (r *raceTimeService) getAttributeFilters() []filter.AttributeId {
	var filters []filter.AttributeId

	// 競馬場
	courses := []filter.AttributeId{
		filter.Tokyo,
		filter.Nakayama,
		filter.Kyoto,
		filter.Hanshin,
		filter.Niigata,
		filter.Chukyo,
		filter.Sapporo,
		filter.Hakodate,
		filter.Fukushima,
		filter.Kokura,
	}

	// 馬場
	surfaces := []filter.AttributeId{
		filter.Turf,
		filter.Dirt,
	}

	// 距離（AttributeIdで定義されているすべての距離）
	distances := []filter.AttributeId{
		filter.Distance1000m,
		filter.Distance1150m,
		filter.Distance1200m,
		filter.Distance1300m,
		filter.Distance1400m,
		filter.Distance1500m,
		filter.Distance1600m,
		filter.Distance1700m,
		filter.Distance1800m,
		filter.Distance1900m,
		filter.Distance2000m,
		filter.Distance2100m,
		filter.Distance2200m,
		filter.Distance2300m,
		filter.Distance2400m,
		filter.Distance2500m,
		filter.Distance2600m,
		filter.Distance3000m,
		filter.Distance3200m,
		filter.Distance3400m,
		filter.Distance3600m,
	}

	// クラス
	classes := []filter.AttributeId{
		filter.Maiden,
		filter.OneWinClass,
		filter.TwoWinClass,
		filter.ThreeWinClass,
		filter.OpenListedClass,
		filter.Grade3,
		filter.Grade2,
		filter.Grade1,
	}

	// 馬場状態
	conditions := []filter.AttributeId{
		filter.GoodToFirm,
		filter.Good,
		filter.Yielding,
		filter.Soft,
	}

	// 年齢
	ages := []filter.AttributeId{
		filter.TwoYearsOld,
		filter.ThreeYearsOld,
		filter.ThreeYearsAndOlder,
		filter.FourYearsAndOlder,
	}

	// フィルターの組み合わせを生成
	for _, course := range courses {
		for _, surface := range surfaces {
			for _, distance := range distances {
				for _, class := range classes {
					for _, condition := range conditions {
						for _, age := range ages {
							filters = append(filters, course|surface|distance|class|condition|age)
						}
					}
				}
			}
		}
	}

	return filters
}
