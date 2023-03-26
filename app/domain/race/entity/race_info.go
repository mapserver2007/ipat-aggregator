package entity

type RaceInfo struct {
	races []*Race
}

func NewRaceInfo(
	races []*Race,
) *RaceInfo {
	return &RaceInfo{
		races: races,
	}
}

func (r *RaceInfo) Races() []*Race {
	return r.races
}

func (r *RaceInfo) Get(idx int) *Race {
	return r.races[idx]
}
