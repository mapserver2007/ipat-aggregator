package spreadsheet_entity

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
}

func NewPlaceUnHitCountData(
	oddsRangeCountSlice []int,
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
