package value_object

type GradeClass int

const (
	NonGrade       GradeClass = 0
	Grade1         GradeClass = 1
	Grade2         GradeClass = 2
	Grade3         GradeClass = 3
	LocalGrade     GradeClass = 4
	OpenClass      GradeClass = 5
	JumpGrade1     GradeClass = 10
	JumpGrade2     GradeClass = 11
	JumpGrade3     GradeClass = 12
	ListedClass    GradeClass = 15
	Jpn1           GradeClass = 19
	Jpn2           GradeClass = 20
	Jpn3           GradeClass = 21
	Maiden         GradeClass = 31 // 未勝利
	OneWinClass    GradeClass = 32 // 1勝クラス
	TwoWinClass    GradeClass = 33 // 2勝クラス
	ThreeWinClass  GradeClass = 34 // 3勝クラス
	JumpMaiden     GradeClass = 35 // 障害未勝利
	JumpOpenClass  GradeClass = 36 // 障害オープン
	MakeDebut      GradeClass = 37 // 新馬
	AllowanceClass GradeClass = 98 // Class1-3は特別戦、AllowanceClassは非特別戦の条件戦
	NonGradeClass  GradeClass = 99 // リステッド,OP,条件戦をまとめるためのクラス
)

var gradeClassMap = map[GradeClass]string{
	NonGrade:       "なし",
	Grade1:         "G1",
	Grade2:         "G2",
	Grade3:         "G3",
	LocalGrade:     "地方重賞",
	OpenClass:      "OP/L/地方重賞",
	JumpGrade1:     "JG1",
	JumpGrade2:     "JG2",
	JumpGrade3:     "JG3",
	ListedClass:    "L",
	Jpn1:           "Jpn1",
	Jpn2:           "Jpn2",
	Jpn3:           "Jpn3",
	Maiden:         "未勝利",
	MakeDebut:      "新馬",
	OneWinClass:    "1勝クラス",
	TwoWinClass:    "2勝クラス",
	ThreeWinClass:  "3勝クラス",
	JumpMaiden:     "障害未勝利",
	JumpOpenClass:  "障害オープン",
	AllowanceClass: "条件戦",
	NonGradeClass:  "平場",
}

func (g GradeClass) String() string {
	gradeClassName, _ := gradeClassMap[g]
	return gradeClassName
}
