package raw_entity

type RawJockeyNetkeiba struct {
	id   int
	name string
}

func NewRawJockeyNetkeiba(
	id int,
	name string,
) *RawJockeyNetkeiba {
	return &RawJockeyNetkeiba{
		id:   id,
		name: name,
	}
}

func (r *RawJockeyNetkeiba) Id() int {
	return r.id
}

func (r *RawJockeyNetkeiba) Name() string {
	return r.name
}
