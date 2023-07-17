package value_object

type GradeClass int

const (
	NonGrade       GradeClass = 0
	Grade1         GradeClass = 1
	Grade2         GradeClass = 2
	Grade3         GradeClass = 3
	OpenClass      GradeClass = 5
	JumpGrade1     GradeClass = 10
	JumpGrade2     GradeClass = 11
	JumpGrade3     GradeClass = 12
	ListedClass    GradeClass = 15
	Jpn1           GradeClass = 19
	Jpn2           GradeClass = 20
	Jpn3           GradeClass = 21
	AllowanceClass GradeClass = 98 // Class1-3は特別戦、AllowanceClassは非特別戦の条件戦
	NonGradeClass  GradeClass = 99 // リステッド,OP,条件戦をまとめるためのクラス
)

var gradeClassMap = map[GradeClass]string{
	NonGrade:       "なし",
	Grade1:         "G1",
	Grade2:         "G2",
	Grade3:         "G3",
	OpenClass:      "OP",
	JumpGrade1:     "JG1",
	JumpGrade2:     "JG2",
	JumpGrade3:     "JG3",
	ListedClass:    "L",
	Jpn1:           "Jpn1",
	Jpn2:           "Jpn2",
	Jpn3:           "Jpn3",
	AllowanceClass: "条件戦",
	NonGradeClass:  "平場",
}

func (g GradeClass) String() string {
	return convertToGradeClassName(g)
}

func convertToGradeClassName(g GradeClass) string {
	gradeClassName, _ := gradeClassMap[g]
	return gradeClassName
}
