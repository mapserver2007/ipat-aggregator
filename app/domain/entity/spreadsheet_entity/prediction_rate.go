package spreadsheet_entity

import (
	"fmt"
	"math"
)

type PredictionRateData struct {
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
	oddsRange1Hit  bool
	oddsRange2Hit  bool
	oddsRange3Hit  bool
	oddsRange4Hit  bool
	oddsRange5Hit  bool
	oddsRange6Hit  bool
	oddsRange7Hit  bool
	oddsRange8Hit  bool
	filterName     string
}

func NewPredictionRateData(
	hitCountData *PlaceHitCountData,
	unHitCountData *PlaceUnHitCountData,
	hitSlice []bool,
) *PredictionRateData {
	hitRate := float64(hitCountData.HitCount()) * 100 / float64(hitCountData.RaceCount())
	oddsRange1Rate := float64(hitCountData.OddsRange1Count()) * 100 / float64(hitCountData.OddsRange1Count()+unHitCountData.OddsRange1Count())
	oddsRange2Rate := float64(hitCountData.OddsRange2Count()) * 100 / float64(hitCountData.OddsRange2Count()+unHitCountData.OddsRange2Count())
	oddsRange3Rate := float64(hitCountData.OddsRange3Count()) * 100 / float64(hitCountData.OddsRange3Count()+unHitCountData.OddsRange3Count())
	oddsRange4Rate := float64(hitCountData.OddsRange4Count()) * 100 / float64(hitCountData.OddsRange4Count()+unHitCountData.OddsRange4Count())
	oddsRange5Rate := float64(hitCountData.OddsRange5Count()) * 100 / float64(hitCountData.OddsRange5Count()+unHitCountData.OddsRange5Count())
	oddsRange6Rate := float64(hitCountData.OddsRange6Count()) * 100 / float64(hitCountData.OddsRange6Count()+unHitCountData.OddsRange6Count())
	oddsRange7Rate := float64(hitCountData.OddsRange7Count()) * 100 / float64(hitCountData.OddsRange7Count()+unHitCountData.OddsRange7Count())
	oddsRange8Rate := float64(hitCountData.OddsRange8Count()) * 100 / float64(hitCountData.OddsRange8Count()+unHitCountData.OddsRange8Count())

	return &PredictionRateData{
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
		oddsRange1Hit:  hitSlice[0],
		oddsRange2Hit:  hitSlice[1],
		oddsRange3Hit:  hitSlice[2],
		oddsRange4Hit:  hitSlice[3],
		oddsRange5Hit:  hitSlice[4],
		oddsRange6Hit:  hitSlice[5],
		oddsRange7Hit:  hitSlice[6],
		oddsRange8Hit:  hitSlice[7],
		filterName:     hitCountData.FilterName(),
	}
}

func (p *PredictionRateData) RaceCount() int {
	return p.raceCount
}

func (p *PredictionRateData) HitRate() float64 {
	return p.hitRate
}

func (p *PredictionRateData) HitRateFormat() string {
	return p.rateFormat(p.hitRate)
}

func (p *PredictionRateData) OddsRange1Rate() float64 {
	return p.oddsRange1Rate
}

func (p *PredictionRateData) OddsRange1RateFormat() string {
	return p.rateFormat(p.oddsRange1Rate)
}

func (p *PredictionRateData) OddsRange2Rate() float64 {
	return p.oddsRange2Rate
}

func (p *PredictionRateData) OddsRange2RateFormat() string {
	return p.rateFormat(p.oddsRange2Rate)
}

func (p *PredictionRateData) OddsRange3Rate() float64 {
	return p.oddsRange3Rate
}

func (p *PredictionRateData) OddsRange3RateFormat() string {
	return p.rateFormat(p.oddsRange3Rate)
}

func (p *PredictionRateData) OddsRange4Rate() float64 {
	return p.oddsRange4Rate
}

func (p *PredictionRateData) OddsRange4RateFormat() string {
	return p.rateFormat(p.oddsRange4Rate)
}

func (p *PredictionRateData) OddsRange5Rate() float64 {
	return p.oddsRange5Rate
}

func (p *PredictionRateData) OddsRange5RateFormat() string {
	return p.rateFormat(p.oddsRange5Rate)
}

func (p *PredictionRateData) OddsRange6Rate() float64 {
	return p.oddsRange6Rate
}

func (p *PredictionRateData) OddsRange6RateFormat() string {
	return p.rateFormat(p.oddsRange6Rate)
}

func (p *PredictionRateData) OddsRange7Rate() float64 {
	return p.oddsRange7Rate
}

func (p *PredictionRateData) OddsRange7RateFormat() string {
	return p.rateFormat(p.oddsRange7Rate)
}

func (p *PredictionRateData) OddsRange8Rate() float64 {
	return p.oddsRange8Rate
}

func (p *PredictionRateData) OddsRange8RateFormat() string {
	return p.rateFormat(p.oddsRange8Rate)
}

func (p *PredictionRateData) OddsRange1Hit() bool {
	return p.oddsRange1Hit
}

func (p *PredictionRateData) OddsRange2Hit() bool {
	return p.oddsRange2Hit
}

func (p *PredictionRateData) OddsRange3Hit() bool {
	return p.oddsRange3Hit
}

func (p *PredictionRateData) OddsRange4Hit() bool {
	return p.oddsRange4Hit
}

func (p *PredictionRateData) OddsRange5Hit() bool {
	return p.oddsRange5Hit
}

func (p *PredictionRateData) OddsRange6Hit() bool {
	return p.oddsRange6Hit
}

func (p *PredictionRateData) OddsRange7Hit() bool {
	return p.oddsRange7Hit
}

func (p *PredictionRateData) OddsRange8Hit() bool {
	return p.oddsRange8Hit
}

func (p *PredictionRateData) FilterName() string {
	return p.filterName
}

func (p *PredictionRateData) rateFormat(rate float64) string {
	if math.IsNaN(rate) {
		return "-"
	}
	return fmt.Sprintf("%.2f%%", rate)
}

type PredictionRateStyle struct {
	matchOddsRangeIndex int
}

func NewPredictionRateStyle(
	onOddsRangeSlice []bool,
) *PredictionRateStyle {
	var matchOddsRangeIndex int
	for idx, match := range onOddsRangeSlice {
		if match {
			matchOddsRangeIndex = idx
			break
		}
	}

	return &PredictionRateStyle{
		matchOddsRangeIndex: matchOddsRangeIndex,
	}
}

func (p *PredictionRateStyle) MatchOddsRangeIndex() int {
	return p.matchOddsRangeIndex
}
