package spreadsheet_entity

import (
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"math"
	"strconv"
)

type PlaceAllInRateData struct {
	winOdds11HitData *PlaceAllInHitData
	winOdds12HitData *PlaceAllInHitData
	winOdds13HitData *PlaceAllInHitData
	winOdds14HitData *PlaceAllInHitData
	winOdds15HitData *PlaceAllInHitData
	winOdds16HitData *PlaceAllInHitData
	winOdds17HitData *PlaceAllInHitData
	winOdds18HitData *PlaceAllInHitData
	winOdds19HitData *PlaceAllInHitData
	winOdds20HitData *PlaceAllInHitData
	winOdds21HitData *PlaceAllInHitData
	winOdds22HitData *PlaceAllInHitData
	winOdds23HitData *PlaceAllInHitData
	winOdds24HitData *PlaceAllInHitData
	winOdds25HitData *PlaceAllInHitData
	winOdds26HitData *PlaceAllInHitData
	winOdds27HitData *PlaceAllInHitData
	winOdds28HitData *PlaceAllInHitData
	winOdds29HitData *PlaceAllInHitData
	winOdds30HitData *PlaceAllInHitData
	winOdds31HitData *PlaceAllInHitData
	winOdds32HitData *PlaceAllInHitData
	winOdds33HitData *PlaceAllInHitData
	winOdds34HitData *PlaceAllInHitData
	winOdds35HitData *PlaceAllInHitData
	winOdds36HitData *PlaceAllInHitData
	winOdds37HitData *PlaceAllInHitData
	winOdds38HitData *PlaceAllInHitData
	winOdds39HitData *PlaceAllInHitData
	filterName       string
}

func NewPlaceAllInRateData(
	hitCountData *PlaceAllInHitCountData,
	unHitCountData *PlaceAllInUnHitCountData,
) *PlaceAllInRateData {
	return &PlaceAllInRateData{
		winOdds11HitData: NewPlaceAllInHitData(hitCountData.WinOdds11Count(), unHitCountData.WinOdds11Count()),
		winOdds12HitData: NewPlaceAllInHitData(hitCountData.WinOdds12Count(), unHitCountData.WinOdds12Count()),
		winOdds13HitData: NewPlaceAllInHitData(hitCountData.WinOdds13Count(), unHitCountData.WinOdds13Count()),
		winOdds14HitData: NewPlaceAllInHitData(hitCountData.WinOdds14Count(), unHitCountData.WinOdds14Count()),
		winOdds15HitData: NewPlaceAllInHitData(hitCountData.WinOdds15Count(), unHitCountData.WinOdds15Count()),
		winOdds16HitData: NewPlaceAllInHitData(hitCountData.WinOdds16Count(), unHitCountData.WinOdds16Count()),
		winOdds17HitData: NewPlaceAllInHitData(hitCountData.WinOdds17Count(), unHitCountData.WinOdds17Count()),
		winOdds18HitData: NewPlaceAllInHitData(hitCountData.WinOdds18Count(), unHitCountData.WinOdds18Count()),
		winOdds19HitData: NewPlaceAllInHitData(hitCountData.WinOdds19Count(), unHitCountData.WinOdds19Count()),
		winOdds20HitData: NewPlaceAllInHitData(hitCountData.WinOdds20Count(), unHitCountData.WinOdds20Count()),
		winOdds21HitData: NewPlaceAllInHitData(hitCountData.WinOdds21Count(), unHitCountData.WinOdds21Count()),
		winOdds22HitData: NewPlaceAllInHitData(hitCountData.WinOdds22Count(), unHitCountData.WinOdds22Count()),
		winOdds23HitData: NewPlaceAllInHitData(hitCountData.WinOdds23Count(), unHitCountData.WinOdds23Count()),
		winOdds24HitData: NewPlaceAllInHitData(hitCountData.WinOdds24Count(), unHitCountData.WinOdds24Count()),
		winOdds25HitData: NewPlaceAllInHitData(hitCountData.WinOdds25Count(), unHitCountData.WinOdds25Count()),
		winOdds26HitData: NewPlaceAllInHitData(hitCountData.WinOdds26Count(), unHitCountData.WinOdds26Count()),
		winOdds27HitData: NewPlaceAllInHitData(hitCountData.WinOdds27Count(), unHitCountData.WinOdds27Count()),
		winOdds28HitData: NewPlaceAllInHitData(hitCountData.WinOdds28Count(), unHitCountData.WinOdds28Count()),
		winOdds29HitData: NewPlaceAllInHitData(hitCountData.WinOdds29Count(), unHitCountData.WinOdds29Count()),
		winOdds30HitData: NewPlaceAllInHitData(hitCountData.WinOdds30Count(), unHitCountData.WinOdds30Count()),
		winOdds31HitData: NewPlaceAllInHitData(hitCountData.WinOdds31Count(), unHitCountData.WinOdds31Count()),
		winOdds32HitData: NewPlaceAllInHitData(hitCountData.WinOdds32Count(), unHitCountData.WinOdds32Count()),
		winOdds33HitData: NewPlaceAllInHitData(hitCountData.WinOdds33Count(), unHitCountData.WinOdds33Count()),
		winOdds34HitData: NewPlaceAllInHitData(hitCountData.WinOdds34Count(), unHitCountData.WinOdds34Count()),
		winOdds35HitData: NewPlaceAllInHitData(hitCountData.WinOdds35Count(), unHitCountData.WinOdds35Count()),
		winOdds36HitData: NewPlaceAllInHitData(hitCountData.WinOdds36Count(), unHitCountData.WinOdds36Count()),
		winOdds37HitData: NewPlaceAllInHitData(hitCountData.WinOdds37Count(), unHitCountData.WinOdds37Count()),
		winOdds38HitData: NewPlaceAllInHitData(hitCountData.WinOdds38Count(), unHitCountData.WinOdds38Count()),
		winOdds39HitData: NewPlaceAllInHitData(hitCountData.WinOdds39Count(), unHitCountData.WinOdds39Count()),
	}
}

func (p *PlaceAllInRateData) WinOdds11HitData() *PlaceAllInHitData {
	return p.winOdds11HitData
}

func (p *PlaceAllInRateData) WinOdds12HitData() *PlaceAllInHitData {
	return p.winOdds12HitData
}

func (p *PlaceAllInRateData) WinOdds13HitData() *PlaceAllInHitData {
	return p.winOdds13HitData
}

func (p *PlaceAllInRateData) WinOdds14HitData() *PlaceAllInHitData {
	return p.winOdds14HitData
}

func (p *PlaceAllInRateData) WinOdds15HitData() *PlaceAllInHitData {
	return p.winOdds15HitData
}

func (p *PlaceAllInRateData) WinOdds16HitData() *PlaceAllInHitData {
	return p.winOdds16HitData
}

func (p *PlaceAllInRateData) WinOdds17HitData() *PlaceAllInHitData {
	return p.winOdds17HitData
}

func (p *PlaceAllInRateData) WinOdds18HitData() *PlaceAllInHitData {
	return p.winOdds18HitData
}

func (p *PlaceAllInRateData) WinOdds19HitData() *PlaceAllInHitData {
	return p.winOdds19HitData
}

func (p *PlaceAllInRateData) WinOdds20HitData() *PlaceAllInHitData {
	return p.winOdds20HitData
}

func (p *PlaceAllInRateData) WinOdds21HitData() *PlaceAllInHitData {
	return p.winOdds21HitData
}

func (p *PlaceAllInRateData) WinOdds22HitData() *PlaceAllInHitData {
	return p.winOdds22HitData
}

func (p *PlaceAllInRateData) WinOdds23HitData() *PlaceAllInHitData {
	return p.winOdds23HitData
}

func (p *PlaceAllInRateData) WinOdds24HitData() *PlaceAllInHitData {
	return p.winOdds24HitData
}

func (p *PlaceAllInRateData) WinOdds25HitData() *PlaceAllInHitData {
	return p.winOdds25HitData
}

func (p *PlaceAllInRateData) WinOdds26HitData() *PlaceAllInHitData {
	return p.winOdds26HitData
}

func (p *PlaceAllInRateData) WinOdds27HitData() *PlaceAllInHitData {
	return p.winOdds27HitData
}

func (p *PlaceAllInRateData) WinOdds28HitData() *PlaceAllInHitData {
	return p.winOdds28HitData
}

func (p *PlaceAllInRateData) WinOdds29HitData() *PlaceAllInHitData {
	return p.winOdds29HitData
}

func (p *PlaceAllInRateData) WinOdds30HitData() *PlaceAllInHitData {
	return p.winOdds30HitData
}

func (p *PlaceAllInRateData) WinOdds31HitData() *PlaceAllInHitData {
	return p.winOdds31HitData
}

func (p *PlaceAllInRateData) WinOdds32HitData() *PlaceAllInHitData {
	return p.winOdds32HitData
}

func (p *PlaceAllInRateData) WinOdds33HitData() *PlaceAllInHitData {
	return p.winOdds33HitData
}

func (p *PlaceAllInRateData) WinOdds34HitData() *PlaceAllInHitData {
	return p.winOdds34HitData
}

func (p *PlaceAllInRateData) WinOdds35HitData() *PlaceAllInHitData {
	return p.winOdds35HitData
}

func (p *PlaceAllInRateData) WinOdds36HitData() *PlaceAllInHitData {
	return p.winOdds36HitData
}

func (p *PlaceAllInRateData) WinOdds37HitData() *PlaceAllInHitData {
	return p.winOdds37HitData
}

func (p *PlaceAllInRateData) WinOdds38HitData() *PlaceAllInHitData {
	return p.winOdds38HitData
}

func (p *PlaceAllInRateData) WinOdds39HitData() *PlaceAllInHitData {
	return p.winOdds39HitData
}

type PlaceAllInHitData struct {
	hitCount   int
	unHitCount int
	hitRate    float64
}

func NewPlaceAllInHitData(
	hitCount int,
	unHitCount int,
) *PlaceAllInHitData {
	return &PlaceAllInHitData{
		hitCount:   hitCount,
		unHitCount: unHitCount,
		hitRate:    float64(hitCount) * 100 / float64(hitCount+unHitCount),
	}
}

func (p *PlaceAllInHitData) HitCount() int {
	return p.hitCount
}

func (p *PlaceAllInHitData) UnHitCount() int {
	return p.unHitCount
}

func (p *PlaceAllInHitData) HitRate() float64 {
	return p.hitRate
}

func (p *PlaceAllInHitData) HitRateFormat() string {
	return p.rateFormat(p.hitRate)
}

func (p *PlaceAllInHitData) rateFormat(rate float64) string {
	if math.IsNaN(rate) {
		return "-"
	}
	return fmt.Sprintf("%s%%", strconv.Itoa(int(math.Round(rate))))
}

type PlaceAllInRateStyle struct {
	winOdds11CellColorType types.CellColorType
	winOdds12CellColorType types.CellColorType
	winOdds13CellColorType types.CellColorType
	winOdds14CellColorType types.CellColorType
	winOdds15CellColorType types.CellColorType
	winOdds16CellColorType types.CellColorType
	winOdds17CellColorType types.CellColorType
	winOdds18CellColorType types.CellColorType
	winOdds19CellColorType types.CellColorType
	winOdds20CellColorType types.CellColorType
	winOdds21CellColorType types.CellColorType
	winOdds22CellColorType types.CellColorType
	winOdds23CellColorType types.CellColorType
	winOdds24CellColorType types.CellColorType
	winOdds25CellColorType types.CellColorType
	winOdds26CellColorType types.CellColorType
	winOdds27CellColorType types.CellColorType
	winOdds28CellColorType types.CellColorType
	winOdds29CellColorType types.CellColorType
	winOdds30CellColorType types.CellColorType
	winOdds31CellColorType types.CellColorType
	winOdds32CellColorType types.CellColorType
	winOdds33CellColorType types.CellColorType
	winOdds34CellColorType types.CellColorType
	winOdds35CellColorType types.CellColorType
	winOdds36CellColorType types.CellColorType
	winOdds37CellColorType types.CellColorType
	winOdds38CellColorType types.CellColorType
	winOdds39CellColorType types.CellColorType
}

func NewPlaceAllInRateStyle(
	data *PlaceAllInRateData,
) *PlaceAllInRateStyle {
	rateColorTypeFunc := func(rate float64) types.CellColorType {
		if rate >= 90 {
			return types.FirstColor
		} else if rate >= 85 && rate < 90 {
			return types.SecondColor
		} else if rate >= 80 && rate < 85 {
			return types.ThirdColor
		}
		return types.NoneColor
	}

	return &PlaceAllInRateStyle{
		winOdds11CellColorType: rateColorTypeFunc(data.WinOdds11HitData().HitRate()),
		winOdds12CellColorType: rateColorTypeFunc(data.WinOdds12HitData().HitRate()),
		winOdds13CellColorType: rateColorTypeFunc(data.WinOdds13HitData().HitRate()),
		winOdds14CellColorType: rateColorTypeFunc(data.WinOdds14HitData().HitRate()),
		winOdds15CellColorType: rateColorTypeFunc(data.WinOdds15HitData().HitRate()),
		winOdds16CellColorType: rateColorTypeFunc(data.WinOdds16HitData().HitRate()),
		winOdds17CellColorType: rateColorTypeFunc(data.WinOdds17HitData().HitRate()),
		winOdds18CellColorType: rateColorTypeFunc(data.WinOdds18HitData().HitRate()),
		winOdds19CellColorType: rateColorTypeFunc(data.WinOdds19HitData().HitRate()),
		winOdds20CellColorType: rateColorTypeFunc(data.WinOdds20HitData().HitRate()),
		winOdds21CellColorType: rateColorTypeFunc(data.WinOdds21HitData().HitRate()),
		winOdds22CellColorType: rateColorTypeFunc(data.WinOdds22HitData().HitRate()),
		winOdds23CellColorType: rateColorTypeFunc(data.WinOdds23HitData().HitRate()),
		winOdds24CellColorType: rateColorTypeFunc(data.WinOdds24HitData().HitRate()),
		winOdds25CellColorType: rateColorTypeFunc(data.WinOdds25HitData().HitRate()),
		winOdds26CellColorType: rateColorTypeFunc(data.WinOdds26HitData().HitRate()),
		winOdds27CellColorType: rateColorTypeFunc(data.WinOdds27HitData().HitRate()),
		winOdds28CellColorType: rateColorTypeFunc(data.WinOdds28HitData().HitRate()),
		winOdds29CellColorType: rateColorTypeFunc(data.WinOdds29HitData().HitRate()),
		winOdds30CellColorType: rateColorTypeFunc(data.WinOdds30HitData().HitRate()),
		winOdds31CellColorType: rateColorTypeFunc(data.WinOdds31HitData().HitRate()),
		winOdds32CellColorType: rateColorTypeFunc(data.WinOdds32HitData().HitRate()),
		winOdds33CellColorType: rateColorTypeFunc(data.WinOdds33HitData().HitRate()),
		winOdds34CellColorType: rateColorTypeFunc(data.WinOdds34HitData().HitRate()),
		winOdds35CellColorType: rateColorTypeFunc(data.WinOdds35HitData().HitRate()),
		winOdds36CellColorType: rateColorTypeFunc(data.WinOdds36HitData().HitRate()),
		winOdds37CellColorType: rateColorTypeFunc(data.WinOdds37HitData().HitRate()),
		winOdds38CellColorType: rateColorTypeFunc(data.WinOdds38HitData().HitRate()),
		winOdds39CellColorType: rateColorTypeFunc(data.WinOdds39HitData().HitRate()),
	}
}

func (p *PlaceAllInRateStyle) WinOdds11CellColorType() types.CellColorType {
	return p.winOdds11CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds12CellColorType() types.CellColorType {
	return p.winOdds12CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds13CellColorType() types.CellColorType {
	return p.winOdds13CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds14CellColorType() types.CellColorType {
	return p.winOdds14CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds15CellColorType() types.CellColorType {
	return p.winOdds15CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds16CellColorType() types.CellColorType {
	return p.winOdds16CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds17CellColorType() types.CellColorType {
	return p.winOdds17CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds18CellColorType() types.CellColorType {
	return p.winOdds18CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds19CellColorType() types.CellColorType {
	return p.winOdds19CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds20CellColorType() types.CellColorType {
	return p.winOdds20CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds21CellColorType() types.CellColorType {
	return p.winOdds21CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds22CellColorType() types.CellColorType {
	return p.winOdds22CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds23CellColorType() types.CellColorType {
	return p.winOdds23CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds24CellColorType() types.CellColorType {
	return p.winOdds24CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds25CellColorType() types.CellColorType {
	return p.winOdds25CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds26CellColorType() types.CellColorType {
	return p.winOdds26CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds27CellColorType() types.CellColorType {
	return p.winOdds27CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds28CellColorType() types.CellColorType {
	return p.winOdds28CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds29CellColorType() types.CellColorType {
	return p.winOdds29CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds30CellColorType() types.CellColorType {
	return p.winOdds30CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds31CellColorType() types.CellColorType {
	return p.winOdds31CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds32CellColorType() types.CellColorType {
	return p.winOdds32CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds33CellColorType() types.CellColorType {
	return p.winOdds33CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds34CellColorType() types.CellColorType {
	return p.winOdds34CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds35CellColorType() types.CellColorType {
	return p.winOdds35CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds36CellColorType() types.CellColorType {
	return p.winOdds36CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds37CellColorType() types.CellColorType {
	return p.winOdds37CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds38CellColorType() types.CellColorType {
	return p.winOdds38CellColorType
}

func (p *PlaceAllInRateStyle) WinOdds39CellColorType() types.CellColorType {
	return p.winOdds39CellColorType
}
