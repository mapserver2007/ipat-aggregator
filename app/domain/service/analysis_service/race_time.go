package analysis_service

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/filter_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type RaceTime interface {
	Create(ctx context.Context,
		races []*data_cache_entity.Race,
		raceTimes []*data_cache_entity.RaceTime,
	) ([]*analysis_entity.RaceTimeCalculable, error)
	Convert(ctx context.Context,
		calculables []*analysis_entity.RaceTimeCalculable,
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
) error {
	filterRaceTimeMap := make(map[filter.AttributeId][]*analysis_entity.RaceTimeCalculable)

	for _, attributeFilter := range r.getAttributeFilters() {
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

	for attributeFilter, calculables := range filterRaceTimeMap {
		_ = attributeFilter
		averageRaceTime := r.calcAverageRaceTime(calculables)
		medianRaceTime := r.calcMedianRaceTime(calculables)
		averageFirst3f := r.calcAverageFirst3f(calculables)
		medianFirst3f := r.calcMedianFirst3f(calculables)
		averageFirst4f := r.calcAverageFirst4f(calculables)
		medianFirst4f := r.calcMedianFirst4f(calculables)
		averageLast3f := r.calcAverageLast3f(calculables)
		medianLast3f := r.calcMedianLast3f(calculables)
		averageLast4f := r.calcAverageLast4f(calculables)
		medianLast4f := r.calcMedianLast4f(calculables)
		averageRap5f := r.calcAverageRap5f(calculables)
		medianRap5f := r.calcMedianRap5f(calculables)

		fmt.Println(averageRaceTime)
		fmt.Println(medianRaceTime)
		fmt.Println(averageFirst3f)
		fmt.Println(medianFirst3f)
		fmt.Println(averageFirst4f)
		fmt.Println(medianFirst4f)
		fmt.Println(averageLast3f)
		fmt.Println(medianLast3f)
		fmt.Println(averageLast4f)
		fmt.Println(medianLast4f)
		fmt.Println(averageRap5f)
		fmt.Println(medianRap5f)
	}

	return nil
}

func (r *raceTimeService) calcAverageRaceTime(
	calculables []*analysis_entity.RaceTimeCalculable,
) string {
	if len(calculables) == 0 {
		return "0:00.0"
	}

	var totalRaceTime time.Duration
	for _, calculable := range calculables {
		totalRaceTime += calculable.Time()
	}
	averageRaceTime := totalRaceTime / time.Duration(len(calculables))

	totalSeconds := averageRaceTime.Seconds()
	minutes := int(totalSeconds) / 60
	seconds := totalSeconds - float64(minutes*60)

	if minutes > 0 {
		return fmt.Sprintf("%d:%04.1f", minutes, seconds)
	}
	return fmt.Sprintf("%.1f", seconds)
}

func (r *raceTimeService) calcMedianRaceTime(
	calculables []*analysis_entity.RaceTimeCalculable,
) string {
	if len(calculables) == 0 {
		return "0:00.0"
	}

	durations := make([]time.Duration, 0, len(calculables))
	for _, c := range calculables {
		durations = append(durations, c.Time())
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

	if minutes > 0 {
		return fmt.Sprintf("%d:%04.1f", minutes, seconds)
	}
	return fmt.Sprintf("%.1f", seconds)
}

func (r *raceTimeService) calcAverageFirst3f(
	calculables []*analysis_entity.RaceTimeCalculable,
) string {
	if len(calculables) == 0 {
		return "0.0"
	}

	var totalFirst3f time.Duration
	for _, calculable := range calculables {
		totalFirst3f += calculable.First3f()
	}
	averageRaceTime := totalFirst3f / time.Duration(len(calculables))

	totalSeconds := averageRaceTime.Seconds()
	minutes := int(totalSeconds) / 60
	seconds := totalSeconds - float64(minutes*60)

	return fmt.Sprintf("%.1f", seconds)
}

func (r *raceTimeService) calcMedianFirst3f(
	calculables []*analysis_entity.RaceTimeCalculable,
) string {
	if len(calculables) == 0 {
		return "0.0"
	}

	durations := make([]time.Duration, 0, len(calculables))
	for _, c := range calculables {
		durations = append(durations, c.First3f())
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

	return fmt.Sprintf("%.1f", seconds)
}

func (r *raceTimeService) calcAverageFirst4f(
	calculables []*analysis_entity.RaceTimeCalculable,
) string {
	if len(calculables) == 0 {
		return "0.0"
	}

	var totalFirst4f time.Duration
	for _, calculable := range calculables {
		totalFirst4f += calculable.First4f()
	}
	averageFirst4f := totalFirst4f / time.Duration(len(calculables))

	totalSeconds := averageFirst4f.Seconds()
	minutes := int(totalSeconds) / 60
	seconds := totalSeconds - float64(minutes*60)

	return fmt.Sprintf("%.1f", seconds)
}

func (r *raceTimeService) calcMedianFirst4f(
	calculables []*analysis_entity.RaceTimeCalculable,
) string {
	if len(calculables) == 0 {
		return "0.0"
	}

	durations := make([]time.Duration, 0, len(calculables))
	for _, c := range calculables {
		durations = append(durations, c.First4f())
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

	return fmt.Sprintf("%.1f", seconds)
}

func (r *raceTimeService) calcAverageLast3f(
	calculables []*analysis_entity.RaceTimeCalculable,
) string {
	if len(calculables) == 0 {
		return "0.0"
	}

	var totalLast3f time.Duration
	for _, calculable := range calculables {
		totalLast3f += calculable.Last3f()
	}
	averageLast3f := totalLast3f / time.Duration(len(calculables))

	totalSeconds := averageLast3f.Seconds()
	minutes := int(totalSeconds) / 60
	seconds := totalSeconds - float64(minutes*60)

	return fmt.Sprintf("%.1f", seconds)
}

func (r *raceTimeService) calcMedianLast3f(
	calculables []*analysis_entity.RaceTimeCalculable,
) string {
	if len(calculables) == 0 {
		return "0.0"
	}

	durations := make([]time.Duration, 0, len(calculables))
	for _, c := range calculables {
		durations = append(durations, c.Last3f())
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

	return fmt.Sprintf("%.1f", seconds)
}

func (r *raceTimeService) calcAverageLast4f(
	calculables []*analysis_entity.RaceTimeCalculable,
) string {
	if len(calculables) == 0 {
		return "0.0"
	}

	var totalLast4f time.Duration
	for _, calculable := range calculables {
		totalLast4f += calculable.Last4f()
	}
	averageLast4f := totalLast4f / time.Duration(len(calculables))

	totalSeconds := averageLast4f.Seconds()
	minutes := int(totalSeconds) / 60
	seconds := totalSeconds - float64(minutes*60)

	return fmt.Sprintf("%.1f", seconds)
}

func (r *raceTimeService) calcMedianLast4f(
	calculables []*analysis_entity.RaceTimeCalculable,
) string {
	if len(calculables) == 0 {
		return "0.0"
	}

	durations := make([]time.Duration, 0, len(calculables))
	for _, c := range calculables {
		durations = append(durations, c.Last4f())
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

	return fmt.Sprintf("%.1f", seconds)
}

func (r *raceTimeService) calcAverageRap5f(
	calculables []*analysis_entity.RaceTimeCalculable,
) string {
	if len(calculables) == 0 {
		return "0.0"
	}

	var totalRap5f time.Duration
	for _, calculable := range calculables {
		totalRap5f += calculable.Rap5f()
	}
	averageRap5f := totalRap5f / time.Duration(len(calculables))

	totalSeconds := averageRap5f.Seconds()
	minutes := int(totalSeconds) / 60
	seconds := totalSeconds - float64(minutes*60)

	return fmt.Sprintf("%.1f", seconds)
}

func (r *raceTimeService) calcMedianRap5f(
	calculables []*analysis_entity.RaceTimeCalculable,
) string {
	if len(calculables) == 0 {
		return "0.0"
	}

	durations := make([]time.Duration, 0, len(calculables))
	for _, c := range calculables {
		durations = append(durations, c.Rap5f())
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

	return fmt.Sprintf("%.1f", seconds)
}

func (r *raceTimeService) getAttributeFilters() []filter.AttributeId {
	return []filter.AttributeId{
		filter.Turf | filter.Tokyo | filter.Distance1600m | filter.GoodToFirm | filter.Maiden | filter.TwoYearsOld,
		filter.Turf | filter.Tokyo | filter.Distance1600m | filter.GoodToFirm | filter.Maiden | filter.ThreeYearsOld,
		filter.Turf | filter.Tokyo | filter.Distance1600m | filter.GoodToFirm | filter.Maiden | filter.ThreeYearsAndOlder,
		filter.Turf | filter.Tokyo | filter.Distance1600m | filter.GoodToFirm | filter.Maiden | filter.FourYearsAndOlder,
	}
}
