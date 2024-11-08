package types

type RaceSexCondition int

const (
	UndefinedRaceSexCondition RaceSexCondition = iota
	NoRaceSexCondition
	FillyAndMareLimited
)

var raceSexConditionMap = map[RaceSexCondition]string{
	UndefinedRaceSexCondition: "未定義性別条件",
	NoRaceSexCondition:        "混合",
	FillyAndMareLimited:       "牝馬限定",
}

func (r RaceSexCondition) Value() int {
	return int(r)
}

func (r RaceSexCondition) String() string {
	raceSexConditionName, _ := raceSexConditionMap[r]
	return raceSexConditionName
}
