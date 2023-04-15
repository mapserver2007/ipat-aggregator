package entity

type Jockey struct {
	jockeyName string
}

func NewJockey(
	jockeyName string,
) *Jockey {
	return &Jockey{
		jockeyName: jockeyName,
	}
}

func (j *Jockey) JockeyName() string {
	return j.jockeyName
}
