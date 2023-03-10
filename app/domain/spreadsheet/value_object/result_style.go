package value_object

type PlaceColor int

const (
	OtherPlace PlaceColor = iota
	FirstPlace
	SecondPlace
)

type GradeClassColor int

const (
	NonGrade GradeClassColor = iota
	Grade1
	Grade2
	Grade3
)
