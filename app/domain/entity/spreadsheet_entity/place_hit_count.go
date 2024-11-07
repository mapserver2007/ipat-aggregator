package spreadsheet_entity

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
	oddsRange9Count int
}

func NewPlaceHitCountData(
	oddsRangeCountSlice []int,
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
		oddsRange9Count: oddsRangeCountSlice[8],
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

func (p *PlaceHitCountData) OddsRange9Count() int {
	return p.oddsRange9Count
}
