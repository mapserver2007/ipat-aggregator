package value_object

type GradeClass int

const (
	NonGrade       GradeClass = 0
	Grade1                    = 1
	Grade2                    = 2
	Grade3                    = 3
	OpenClass                 = 5
	JumpGrade1                = 10
	JumpGrade2                = 11
	JumpGrade3                = 12
	ListedClass               = 15
	Jpn1                      = 19
	Jpn2                      = 20
	Jpn3                      = 21
	AllowanceClass            = 98 // Class1-3は特別戦、AllowanceClassは非特別戦の条件戦
	NonGradeClass             = 99 // リステッド,OP,条件戦をまとめるためのクラス
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
