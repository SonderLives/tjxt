package logic

import (
	"database/sql"

	"tjxt/apps/learning/rpc/internal/model"
	"tjxt/apps/learning/rpc/pb"
)

const timeLayout = "2006-01-02 15:04:05"

// 学习状态 / 计划状态 对外枚举字符串
const (
	lessonStatusNotStart = "NOT_BEGIN"
	lessonStatusLearning = "LEARNING"
	lessonStatusFinished = "FINISHED"
	lessonStatusExpired  = "EXPIRED"

	planStatusNone    = "NO_PLAN"
	planStatusRunning = "PLAN_RUNNING"
)

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
		CreateTime:      l.CreateTime.Format(timeLayout),
		ExpireTime:      nullTime(l.ExpireTime),
		LatestLearnTime: nullTime(l.LatestLearnTime),
	}
}

func lessonStatusDesc(status int64) string {
	switch status {
	case model.LessonStatusInLearn:
		return lessonStatusLearning
	case model.LessonStatusDone:
		return lessonStatusFinished
	case model.LessonStatusExpired:
		return lessonStatusExpired
	default:
		return lessonStatusNotStart
	}
}

func planStatusDesc(status int64) string {
	if status == model.PlanStatusInPlan {
		return planStatusRunning
	}
	return planStatusNone
}

func nullInt64(v sql.NullInt64) int64 {
	if v.Valid {
		return v.Int64
	}
	return 0
}

func nullTime(t sql.NullTime) string {
	if t.Valid {
		return t.Time.Format(timeLayout)
	}
	return ""
}

func calcPages(total, pageSize int64) int64 {
	if pageSize <= 0 {
		return 0
	}
	return (total + pageSize - 1) / pageSize
}
