package event

// CourseEvent 课程上下架事件
// 与 search 服务 ES 索引同步的契约对齐（json 字段名必须匹配）：
// course 服务在课程上架/下架时向 ExchangeCourse(course.events) 发布，
// routing key 为 RoutingKeyCourseUp(course.up) / RoutingKeyCourseDown(course.down)。
type CourseEvent struct {
	CourseID int64 `json:"courseId"`
}
