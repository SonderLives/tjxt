package learninglogic

import (
	"database/sql"

	"tjxt/app/learning/rpc/internal/model"
	"tjxt/app/learning/rpc/pb"
	"tjxt/pkg/timex"
)

// nullInt64 解包 sql.NullInt64：无效值返回 0
func nullInt64(v sql.NullInt64) int64 {
	if v.Valid {
		return v.Int64
	}
	return 0
}

// toLessonVO 将 model.LearningLesson 转为 pb 视图对象。
// 课程维度的字段（名称/封面/价格/小节数/小节名）由 API 层经 CourseRpc 回填，
// RPC 层只填充 learning_lesson 自身字段。
func toLessonVO(l *model.LearningLesson) *pb.LearningLessonVO {
	if l == nil {
		return nil
	}
	return &pb.LearningLessonVO{
		Id:              l.Id,
		CourseId:        l.CourseId,
		LearnedSections: int32(l.LearnedSections),
		Status:          lessonStatusDesc(l.Status),
		PlanStatus:      planStatusDesc(l.PlanStatus),
		WeekFreq:        int32(nullInt64(l.WeekFreq)),
		LatestSectionId: nullInt64(l.LatestSectionId),
		CreateTime:      timex.Format(l.CreateTime),
		ExpireTime:      timex.FormatNullTime(l.ExpireTime),
		LatestLearnTime: timex.FormatNullTime(l.LatestLearnTime),
	}
}
