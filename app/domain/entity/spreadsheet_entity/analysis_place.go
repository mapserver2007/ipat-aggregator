package spreadsheet_entity

import (
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"math"
)

type AnalysisPlace struct {
	rateData       *PlaceRateData
	rateStyle      *PlaceRateStyle
	hitCountData   *PlaceHitCountData
	unHitCountData *PlaceUnHitCountData
}

func NewAnalysisPlace(
	rateData *PlaceRateData,
	rateStyle *PlaceRateStyle,
	hitCountData *PlaceHitCountData,
	unHitCountData *PlaceUnHitCountData,
) *AnalysisPlace {
	return &AnalysisPlace{
		rateData:       rateData,
		rateStyle:      rateStyle,
		hitCountData:   hitCountData,
		unHitCountData: unHitCountData,
	}
}

func (a *AnalysisPlace) RateData() *PlaceRateData {
	return a.rateData
}

func (a *AnalysisPlace) RateStyle() *PlaceRateStyle {
	return a.rateStyle
}

func (a *AnalysisPlace) HitCountData() *PlaceHitCountData {
	return a.hitCountData
}

func (a *AnalysisPlace) UnHitCountData() *PlaceUnHitCountData {
	return a.unHitCountData
}

type PlaceRateData struct {
	raceCount      int
	hitRate        float64
	oddsRange1Rate float64
	oddsRange2Rate float64
	oddsRange3Rate float64
	oddsRange4Rate float64
	oddsRange5Rate float64
	oddsRange6Rate float64
	oddsRange7Rate float64
	oddsRange8Rate float64
	filterName     string
}

func NewPlaceRateData(
	hitCountData *PlaceHitCountData,
	unHitCountData *PlaceUnHitCountData,
) *PlaceRateData {
	hitRate := float64(hitCountData.HitCount()) * 100 / float64(hitCountData.RaceCount())
	oddsRange1Rate := float64(hitCountData.OddsRange1Count()) * 100 / float64(hitCountData.OddsRange1Count()+unHitCountData.OddsRange1Count())
	oddsRange2Rate := float64(hitCountData.OddsRange2Count()) * 100 / float64(hitCountData.OddsRange2Count()+unHitCountData.OddsRange2Count())
	oddsRange3Rate := float64(hitCountData.OddsRange3Count()) * 100 / float64(hitCountData.OddsRange3Count()+unHitCountData.OddsRange3Count())
	oddsRange4Rate := float64(hitCountData.OddsRange4Count()) * 100 / float64(hitCountData.OddsRange4Count()+unHitCountData.OddsRange4Count())
	oddsRange5Rate := float64(hitCountData.OddsRange5Count()) * 100 / float64(hitCountData.OddsRange5Count()+unHitCountData.OddsRange5Count())
	oddsRange6Rate := float64(hitCountData.OddsRange6Count()) * 100 / float64(hitCountData.OddsRange6Count()+unHitCountData.OddsRange6Count())
	oddsRange7Rate := float64(hitCountData.OddsRange7Count()) * 100 / float64(hitCountData.OddsRange7Count()+unHitCountData.OddsRange7Count())
	oddsRange8Rate := float64(hitCountData.OddsRange8Count()) * 100 / float64(hitCountData.OddsRange8Count()+unHitCountData.OddsRange8Count())

	return &PlaceRateData{
		raceCount:      hitCountData.RaceCount(),
		hitRate:        hitRate,
		oddsRange1Rate: oddsRange1Rate,
		oddsRange2Rate: oddsRange2Rate,
		oddsRange3Rate: oddsRange3Rate,
		oddsRange4Rate: oddsRange4Rate,
		oddsRange5Rate: oddsRange5Rate,
		oddsRange6Rate: oddsRange6Rate,
		oddsRange7Rate: oddsRange7Rate,
		oddsRange8Rate: oddsRange8Rate,
		filterName:     hitCountData.FilterName(),
	}
}

func (p *PlaceRateData) RaceCount() int {
	return p.raceCount
}

func (p *PlaceRateData) HitRate() float64 {
	return p.hitRate
}

func (p *PlaceRateData) HitRateFormat() string {
	return p.rateFormat(p.hitRate)
}

func (p *PlaceRateData) OddsRange1Rate() float64 {
	return p.oddsRange1Rate
}

func (p *PlaceRateData) OddsRange1RateFormat() string {
	return p.rateFormat(p.oddsRange1Rate)
}

func (p *PlaceRateData) OddsRange2Rate() float64 {
	return p.oddsRange2Rate
}

func (p *PlaceRateData) OddsRange2RateFormat() string {
	return p.rateFormat(p.oddsRange2Rate)
}

func (p *PlaceRateData) OddsRange3Rate() float64 {
	return p.oddsRange3Rate
}

func (p *PlaceRateData) OddsRange3RateFormat() string {
	return p.rateFormat(p.oddsRange3Rate)
}

func (p *PlaceRateData) OddsRange4Rate() float64 {
	return p.oddsRange4Rate
}

func (p *PlaceRateData) OddsRange4RateFormat() string {
	return p.rateFormat(p.oddsRange4Rate)
}

func (p *PlaceRateData) OddsRange5Rate() float64 {
	return p.oddsRange5Rate
}

func (p *PlaceRateData) OddsRange5RateFormat() string {
	return p.rateFormat(p.oddsRange5Rate)
}

func (p *PlaceRateData) OddsRange6Rate() float64 {
	return p.oddsRange6Rate
}

func (p *PlaceRateData) OddsRange6RateFormat() string {
	return p.rateFormat(p.oddsRange6Rate)
}

func (p *PlaceRateData) OddsRange7Rate() float64 {
	return p.oddsRange7Rate
}

func (p *PlaceRateData) OddsRange7RateFormat() string {
	return p.rateFormat(p.oddsRange7Rate)
}

func (p *PlaceRateData) OddsRange8Rate() float64 {
	return p.oddsRange8Rate
}

func (p *PlaceRateData) OddsRange8RateFormat() string {
	return p.rateFormat(p.oddsRange8Rate)
}

func (p *PlaceRateData) FilterName() string {
	return p.filterName
}

func (p *PlaceRateData) rateFormat(rate float64) string {
	if math.IsNaN(rate) {
		return "-"
	}
	return fmt.Sprintf("%.2f%%", rate)
}

type PlaceRateStyle struct {
	oddsRange1CellColorType types.CellColorType
	oddsRange2CellColorType types.CellColorType
	oddsRange3CellColorType types.CellColorType
	oddsRange4CellColorType types.CellColorType
	oddsRange5CellColorType types.CellColorType
	oddsRange6CellColorType types.CellColorType
	oddsRange7CellColorType types.CellColorType
	oddsRange8CellColorType types.CellColorType
}

func NewPlaceStyle(
	data *PlaceRateData,
) *PlaceRateStyle {
	rateColorTypeFunc := func(rate float64) types.CellColorType {
		if rate >= 75 {
			return types.FirstColor
		} else if rate >= 50 && rate < 75 {
			return types.SecondColor
		} else if rate >= 33 && rate < 50 {
			return types.ThirdColor
		}
		return types.NoneColor
	}

	return &PlaceRateStyle{
		oddsRange1CellColorType: rateColorTypeFunc(data.OddsRange1Rate()),
		oddsRange2CellColorType: rateColorTypeFunc(data.OddsRange2Rate()),
		oddsRange3CellColorType: rateColorTypeFunc(data.OddsRange3Rate()),
		oddsRange4CellColorType: rateColorTypeFunc(data.OddsRange4Rate()),
		oddsRange5CellColorType: rateColorTypeFunc(data.OddsRange5Rate()),
		oddsRange6CellColorType: rateColorTypeFunc(data.OddsRange6Rate()),
		oddsRange7CellColorType: rateColorTypeFunc(data.OddsRange7Rate()),
		oddsRange8CellColorType: rateColorTypeFunc(data.OddsRange8Rate()),
	}
}

func (p *PlaceRateStyle) OddsRange1CellColorType() types.CellColorType {
	return p.oddsRange1CellColorType
}

func (p *PlaceRateStyle) OddsRange2CellColorType() types.CellColorType {
	return p.oddsRange2CellColorType
}

func (p *PlaceRateStyle) OddsRange3CellColorType() types.CellColorType {
	return p.oddsRange3CellColorType
}

func (p *PlaceRateStyle) OddsRange4CellColorType() types.CellColorType {
	return p.oddsRange4CellColorType
}

func (p *PlaceRateStyle) OddsRange5CellColorType() types.CellColorType {
	return p.oddsRange5CellColorType
}

func (p *PlaceRateStyle) OddsRange6CellColorType() types.CellColorType {
	return p.oddsRange6CellColorType
}

func (p *PlaceRateStyle) OddsRange7CellColorType() types.CellColorType {
	return p.oddsRange7CellColorType
}

func (p *PlaceRateStyle) OddsRange8CellColorType() types.CellColorType {
	return p.oddsRange8CellColorType
}

type PlaceHitCountData struct {
	raceCount       int
	hitCount        int
	oddsRange1Count int
	oddsRange2Count int
	oddsRange3Count int
	oddsRange4Count int
	oddsRange5Count int
	oddsRange6Count int
	oddsRange7Count int
	oddsRange8Count int
	filterName      string
}

func NewPlaceHitCountData(
	oddsRangeCountSlice []int,
	filter filter.Id,
	raceCount int,
) *PlaceHitCountData {
	hitCount := 0
	for _, n := range oddsRangeCountSlice {
		hitCount += n
	}

	return &PlaceHitCountData{
		raceCount:       raceCount,
		hitCount:        hitCount,
		oddsRange1Count: oddsRangeCountSlice[0],
		oddsRange2Count: oddsRangeCountSlice[1],
		oddsRange3Count: oddsRangeCountSlice[2],
		oddsRange4Count: oddsRangeCountSlice[3],
		oddsRange5Count: oddsRangeCountSlice[4],
		oddsRange6Count: oddsRangeCountSlice[5],
		oddsRange7Count: oddsRangeCountSlice[6],
		oddsRange8Count: oddsRangeCountSlice[7],
		filterName:      filter.String(),
	}
}

func (p *PlaceHitCountData) RaceCount() int {
	return p.raceCount
}

func (p *PlaceHitCountData) HitCount() int {
	return p.hitCount
}

func (p *PlaceHitCountData) OddsRange1Count() int {
	return p.oddsRange1Count
}

func (p *PlaceHitCountData) OddsRange2Count() int {
	return p.oddsRange2Count
}

func (p *PlaceHitCountData) OddsRange3Count() int {
	return p.oddsRange3Count
}

func (p *PlaceHitCountData) OddsRange4Count() int {
	return p.oddsRange4Count
}

func (p *PlaceHitCountData) OddsRange5Count() int {
	return p.oddsRange5Count
}

func (p *PlaceHitCountData) OddsRange6Count() int {
	return p.oddsRange6Count
}

func (p *PlaceHitCountData) OddsRange7Count() int {
	return p.oddsRange7Count
}

func (p *PlaceHitCountData) OddsRange8Count() int {
	return p.oddsRange8Count
}

func (p *PlaceHitCountData) FilterName() string {
	return p.filterName
}

type PlaceUnHitCountData struct {
	raceCount       int
	unHitCount      int
	oddsRange1Count int
	oddsRange2Count int
	oddsRange3Count int
	oddsRange4Count int
	oddsRange5Count int
	oddsRange6Count int
	oddsRange7Count int
	oddsRange8Count int
	filterName      string
}

func NewPlaceUnHitCountData(
	oddsRangeCountSlice []int,
	filter filter.Id,
	raceCount int,
) *PlaceUnHitCountData {
	unHitCount := 0
	for _, n := range oddsRangeCountSlice {
		unHitCount += n
	}

	return &PlaceUnHitCountData{
		raceCount:       raceCount,
		unHitCount:      unHitCount,
		oddsRange1Count: oddsRangeCountSlice[0],
		oddsRange2Count: oddsRangeCountSlice[1],
		oddsRange3Count: oddsRangeCountSlice[2],
		oddsRange4Count: oddsRangeCountSlice[3],
		oddsRange5Count: oddsRangeCountSlice[4],
		oddsRange6Count: oddsRangeCountSlice[5],
		oddsRange7Count: oddsRangeCountSlice[6],
		oddsRange8Count: oddsRangeCountSlice[7],
		filterName:      filter.String(),
	}
}

func (p *PlaceUnHitCountData) RaceCount() int {
	return p.raceCount
}

func (p *PlaceUnHitCountData) UnHitCount() int {
	return p.unHitCount
}

func (p *PlaceUnHitCountData) OddsRange1Count() int {
	return p.oddsRange1Count
}

func (p *PlaceUnHitCountData) OddsRange2Count() int {
	return p.oddsRange2Count
}

func (p *PlaceUnHitCountData) OddsRange3Count() int {
	return p.oddsRange3Count
}

func (p *PlaceUnHitCountData) OddsRange4Count() int {
	return p.oddsRange4Count
}

func (p *PlaceUnHitCountData) OddsRange5Count() int {
	return p.oddsRange5Count
}

func (p *PlaceUnHitCountData) OddsRange6Count() int {
	return p.oddsRange6Count
}

func (p *PlaceUnHitCountData) OddsRange7Count() int {
	return p.oddsRange7Count
}

func (p *PlaceUnHitCountData) OddsRange8Count() int {
	return p.oddsRange8Count
}

func (p *PlaceUnHitCountData) FilterName() string {
	return p.filterName
}
