package spreadsheet_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"

type PlaceAllInUnHitCountData struct {
	raceCount      int
	hitCount       int
	winOdds11Count int
	winOdds12Count int
	winOdds13Count int
	winOdds14Count int
	winOdds15Count int
	winOdds16Count int
	winOdds17Count int
	winOdds18Count int
	winOdds19Count int
	winOdds20Count int
	winOdds21Count int
	winOdds22Count int
	winOdds23Count int
	winOdds24Count int
	winOdds25Count int
	winOdds26Count int
	winOdds27Count int
	winOdds28Count int
	winOdds29Count int
	winOdds30Count int
	winOdds31Count int
	winOdds32Count int
	winOdds33Count int
	winOdds34Count int
	winOdds35Count int
	winOdds36Count int
	winOdds37Count int
	winOdds38Count int
	winOdds39Count int
	filterName     string
}

func NewPlaceAllInUnHitCountData(
	winOddsCountSlice []int,
	filter filter.Id,
	raceCount int,
) *PlaceAllInUnHitCountData {
	hitCount := 0
	for _, n := range winOddsCountSlice {
		hitCount += n
	}

	return &PlaceAllInUnHitCountData{
		raceCount:      raceCount,
		hitCount:       hitCount,
		winOdds11Count: winOddsCountSlice[0],
		winOdds12Count: winOddsCountSlice[1],
		winOdds13Count: winOddsCountSlice[2],
		winOdds14Count: winOddsCountSlice[3],
		winOdds15Count: winOddsCountSlice[4],
		winOdds16Count: winOddsCountSlice[5],
		winOdds17Count: winOddsCountSlice[6],
		winOdds18Count: winOddsCountSlice[7],
		winOdds19Count: winOddsCountSlice[8],
		winOdds20Count: winOddsCountSlice[9],
		winOdds21Count: winOddsCountSlice[10],
		winOdds22Count: winOddsCountSlice[11],
		winOdds23Count: winOddsCountSlice[12],
		winOdds24Count: winOddsCountSlice[13],
		winOdds25Count: winOddsCountSlice[14],
		winOdds26Count: winOddsCountSlice[15],
		winOdds27Count: winOddsCountSlice[16],
		winOdds28Count: winOddsCountSlice[17],
		winOdds29Count: winOddsCountSlice[18],
		winOdds30Count: winOddsCountSlice[19],
		winOdds31Count: winOddsCountSlice[20],
		winOdds32Count: winOddsCountSlice[21],
		winOdds33Count: winOddsCountSlice[22],
		winOdds34Count: winOddsCountSlice[23],
		winOdds35Count: winOddsCountSlice[24],
		winOdds36Count: winOddsCountSlice[25],
		winOdds37Count: winOddsCountSlice[26],
		winOdds38Count: winOddsCountSlice[27],
		winOdds39Count: winOddsCountSlice[28],
		filterName:     filter.String(),
	}
}

func (p *PlaceAllInUnHitCountData) RaceCount() int {
	return p.raceCount
}

func (p *PlaceAllInUnHitCountData) HitCount() int {
	return p.hitCount
}

func (p *PlaceAllInUnHitCountData) WinOdds11Count() int {
	return p.winOdds11Count
}

func (p *PlaceAllInUnHitCountData) WinOdds12Count() int {
	return p.winOdds12Count
}

func (p *PlaceAllInUnHitCountData) WinOdds13Count() int {
	return p.winOdds13Count
}

func (p *PlaceAllInUnHitCountData) WinOdds14Count() int {
	return p.winOdds14Count
}

func (p *PlaceAllInUnHitCountData) WinOdds15Count() int {
	return p.winOdds15Count
}

func (p *PlaceAllInUnHitCountData) WinOdds16Count() int {
	return p.winOdds16Count
}

func (p *PlaceAllInUnHitCountData) WinOdds17Count() int {
	return p.winOdds17Count
}

func (p *PlaceAllInUnHitCountData) WinOdds18Count() int {
	return p.winOdds18Count
}

func (p *PlaceAllInUnHitCountData) WinOdds19Count() int {
	return p.winOdds19Count
}

func (p *PlaceAllInUnHitCountData) WinOdds20Count() int {
	return p.winOdds20Count
}

func (p *PlaceAllInUnHitCountData) WinOdds21Count() int {
	return p.winOdds21Count
}

func (p *PlaceAllInUnHitCountData) WinOdds22Count() int {
	return p.winOdds22Count
}

func (p *PlaceAllInUnHitCountData) WinOdds23Count() int {
	return p.winOdds23Count
}

func (p *PlaceAllInUnHitCountData) WinOdds24Count() int {
	return p.winOdds24Count
}

func (p *PlaceAllInUnHitCountData) WinOdds25Count() int {
	return p.winOdds25Count
}

func (p *PlaceAllInUnHitCountData) WinOdds26Count() int {
	return p.winOdds26Count
}

func (p *PlaceAllInUnHitCountData) WinOdds27Count() int {
	return p.winOdds27Count
}

func (p *PlaceAllInUnHitCountData) WinOdds28Count() int {
	return p.winOdds28Count
}

func (p *PlaceAllInUnHitCountData) WinOdds29Count() int {
	return p.winOdds29Count
}

func (p *PlaceAllInUnHitCountData) WinOdds30Count() int {
	return p.winOdds30Count
}

func (p *PlaceAllInUnHitCountData) WinOdds31Count() int {
	return p.winOdds31Count
}

func (p *PlaceAllInUnHitCountData) WinOdds32Count() int {
	return p.winOdds32Count
}

func (p *PlaceAllInUnHitCountData) WinOdds33Count() int {
	return p.winOdds33Count
}

func (p *PlaceAllInUnHitCountData) WinOdds34Count() int {
	return p.winOdds34Count
}

func (p *PlaceAllInUnHitCountData) WinOdds35Count() int {
	return p.winOdds35Count
}

func (p *PlaceAllInUnHitCountData) WinOdds36Count() int {
	return p.winOdds36Count
}

func (p *PlaceAllInUnHitCountData) WinOdds37Count() int {
	return p.winOdds37Count
}

func (p *PlaceAllInUnHitCountData) WinOdds38Count() int {
	return p.winOdds38Count
}

func (p *PlaceAllInUnHitCountData) WinOdds39Count() int {
	return p.winOdds39Count
}
