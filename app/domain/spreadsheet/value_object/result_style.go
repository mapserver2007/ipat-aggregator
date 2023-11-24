package value_object

type PlaceColor int

const (
	OtherPlace PlaceColor = iota
	FirstPlace
	SecondPlace
	ThirdPlace
)

type PopularColor int

const (
	OtherPopular PopularColor = iota
	FirstPopular
	SecondPopular
	ThirdPopular
)

type GradeClassColor int

const (
	NonGrade GradeClassColor = iota
	Grade1
	Grade2
	Grade3
)
