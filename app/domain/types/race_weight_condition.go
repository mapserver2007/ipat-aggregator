package types

type RaceWeightCondition int

const (
	UndefinedRaceWeightCondition RaceWeightCondition = iota
	AgeWeight
	FixedWeight
	SpecialWeight
	HandicapWeight
)

var raceWeightConditionMap = map[RaceWeightCondition]string{
	UndefinedRaceWeightCondition: "未定義重量条件",
	AgeWeight:                    "馬齢",
	FixedWeight:                  "定量",
	SpecialWeight:                "別定",
	HandicapWeight:               "ハンデ",
}

func (r RaceWeightCondition) Value() int {
	return int(r)
}

func (r RaceWeightCondition) String() string {
	raceWeightConditionName, _ := raceWeightConditionMap[r]
	return raceWeightConditionName
}
