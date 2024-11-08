package types

type DistanceCategory int

const (
	UndefinedDistanceCategory DistanceCategory = iota
	TurfSprint
	TurfMile
	TurfIntermediate
	TurfLong
	TurfExtended
	DirtSprint
	DirtMile
	DirtIntermediate
	DirtLong
	DirtExtended
	JumpAllDistance
)

// SMILE定義は米国方式
var distanceCategoryMap = map[DistanceCategory]string{
	UndefinedDistanceCategory: "未定義距離",
	TurfSprint:                "芝1000~1300",
	TurfMile:                  "芝1301~1899",
	TurfIntermediate:          "芝1900~2100",
	TurfLong:                  "芝2101~2700",
	TurfExtended:              "芝2701~",
	DirtSprint:                "ダ1000~1300",
	DirtMile:                  "ダ1301~1899",
	DirtIntermediate:          "ダ1900~2100",
	DirtLong:                  "ダ2101~2700",
	DirtExtended:              "ダ2701~", // レースとしては存在しない
	JumpAllDistance:           "障害全距離",
}

func NewDistanceCategory(distance int, courseCategory CourseCategory) DistanceCategory {
	if courseCategory == Jump {
		return JumpAllDistance
	}
	if distance >= 1000 && distance <= 1300 {
		if courseCategory == Turf {
			return TurfSprint
		} else if courseCategory == Dirt {
			return DirtSprint
		}
	} else if distance >= 1301 && distance <= 1899 {
		if courseCategory == Turf {
			return TurfMile
		} else if courseCategory == Dirt {
			return DirtMile
		}
	} else if distance >= 1900 && distance <= 2100 {
		if courseCategory == Turf {
			return TurfIntermediate
		} else if courseCategory == Dirt {
			return DirtIntermediate
		}
	} else if distance >= 2101 && distance <= 2700 {
		if courseCategory == Turf {
			return TurfLong
		} else if courseCategory == Dirt {
			return DirtLong
		}
	} else if distance >= 2701 {
		if courseCategory == Turf {
			return TurfExtended
		} else if courseCategory == Dirt {
			return DirtExtended
		}
	}

	return UndefinedDistanceCategory
}

func (d DistanceCategory) Value() int {
	return int(d)
}

func (d DistanceCategory) String() string {
	distanceCategoryName, _ := distanceCategoryMap[d]
	return distanceCategoryName
}
