package courselogic

// 分类状态
const (
	CategoryStatusNormal  = 1
	CategoryStatusDisable = 2
)

// 课程状态（与 course / course_draft 表注释一致）
const (
	CourseStatusPending   = 1 // 待上架
	CourseStatusUpShelf   = 2 // 已上架
	CourseStatusDownShelf = 3 // 下架
	CourseStatusFinished  = 4 // 已完结
)

// 目录类型
const (
	CatalogueTypeChapter  = 1 // 章
	CatalogueTypeSection  = 2 // 节
	CatalogueTypePractice = 3 // 测试/练习
)

// courseStatusDesc 课程状态中文描述
func courseStatusDesc(status int64) string {
	switch status {
	case CourseStatusPending:
		return "待上架"
	case CourseStatusUpShelf:
		return "已上架"
	case CourseStatusDownShelf:
		return "已下架"
	case CourseStatusFinished:
		return "已完结"
	default:
		return ""
	}
}

// categoryStatusDesc 分类状态中文描述
func categoryStatusDesc(status int64) string {
	switch status {
	case CategoryStatusNormal:
		return "正常"
	case CategoryStatusDisable:
		return "禁用"
	default:
		return ""
	}
}
