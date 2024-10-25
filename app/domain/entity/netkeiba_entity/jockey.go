package netkeiba_entity

type Jockey struct {
	id   string
	name string
}

func NewJockey(
	id string,
	name string,
) *Jockey {
	return &Jockey{
		id:   id,
		name: name,
	}
}

func (j *Jockey) Id() string {
	return j.id
}

func (j *Jockey) Name() string {
	return j.name
}
