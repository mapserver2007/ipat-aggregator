package types

type Organizer int

const (
	UnknownOrganizer Organizer = iota
	JRA
	NAR
	OverseaOrganizer
)

func NewOrganizer(value int) Organizer {
	switch value {
	case 1:
		return JRA
	case 2:
		return NAR
	case 3:
		return OverseaOrganizer
	}
	return UnknownOrganizer
}

func (o Organizer) Value() int {
	return int(o)
}
