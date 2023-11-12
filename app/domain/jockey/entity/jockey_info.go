package entity

type JockeyInfo struct {
	jockeys []*Jockey
}

func NewJockeyInfo(
	jockeys []*Jockey,
) *JockeyInfo {
	return &JockeyInfo{
		jockeys: jockeys,
	}
}

func (j *JockeyInfo) Jockeys() []*Jockey {
	return j.jockeys
}
