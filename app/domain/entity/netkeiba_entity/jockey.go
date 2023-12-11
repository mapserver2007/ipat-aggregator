package netkeiba_entity

type Jockey struct {
	id   int
	name string
}

func NewJockey(
	id int,
	name string,
) *Jockey {
	return &Jockey{
		id:   id,
		name: name,
	}
}

func (j *Jockey) Id() int {
	return j.id
}

func (j *Jockey) Name() string {
	return j.name
}
