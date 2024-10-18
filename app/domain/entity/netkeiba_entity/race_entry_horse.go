package netkeiba_entity

type RaceEntryHorse struct {
	horseId       string
	horseName     string
	bracketNumber int
	horseNumber   int
	jockeyId      int
}

func NewRaceEntryHorse(
	horseId string,
	horseName string,
	bracketNumber int,
	horseNumber int,
	jockeyId int,
) *RaceEntryHorse {
	return &RaceEntryHorse{
		horseId:       horseId,
		horseName:     horseName,
		bracketNumber: bracketNumber,
		horseNumber:   horseNumber,
		jockeyId:      jockeyId,
	}
}

func (r *RaceEntryHorse) HorseId() string {
	return r.horseId
}

func (r *RaceEntryHorse) HorseName() string {
	return r.horseName
}

func (r *RaceEntryHorse) BracketNumber() int {
	return r.bracketNumber
}

func (r *RaceEntryHorse) HorseNumber() int {
	return r.horseNumber
}

func (r *RaceEntryHorse) JockeyId() int {
	return r.jockeyId
}
