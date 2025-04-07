package types

type RaceAgeCondition int

const (
	UnknownRaceAgeCondition RaceAgeCondition = iota
	TwoYearsOld
	ThreeYearsOld
	ThreeYearsAndOlder
	FourYearsAndOlder
)

var raceAgeConditionMap = map[RaceAgeCondition]string{
	UnknownRaceAgeCondition: "未定義年齡条件",
	TwoYearsOld:             "2歳",
	ThreeYearsOld:           "3歳",
	ThreeYearsAndOlder:      "3歳上",
	FourYearsAndOlder:       "4歳上",
}

func (r RaceAgeCondition) Value() int {
	return int(r)
}

func (r RaceAgeCondition) String() string {
	raceAgeConditionName, _ := raceAgeConditionMap[r]
	return raceAgeConditionName
}
