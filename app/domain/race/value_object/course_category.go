package value_object

type CourseCategory int

const (
	NonCourseCategory CourseCategory = iota
	Turf
	Dirt
	Jump
)

var courseCategoryMap = map[CourseCategory]string{
	NonCourseCategory: "不明",
	Turf:              "芝",
	Dirt:              "ダート",
	Jump:              "障害",
}

func NewCourseCategory(name string) CourseCategory {
	var courseCategory CourseCategory
	switch name {
	case "芝":
		courseCategory = Turf
	case "ダ":
		courseCategory = Dirt
	case "障":
		courseCategory = Jump
	}

	return courseCategory
}

func (c CourseCategory) String() string {
	return convertToCourseCategoryName(c)
}

func convertToCourseCategoryName(c CourseCategory) string {
	courseCategoryName, _ := courseCategoryMap[c]
	return courseCategoryName
}
