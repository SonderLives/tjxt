package learninglogic

import (
	"tjxt/app/learning/rpc/internal/model"
)

// 学习状态 / 计划状态 对外枚举字符串
const (
	lessonStatusNotStart = "NOT_BEGIN"
	lessonStatusLearning = "LEARNING"
	lessonStatusFinished = "FINISHED"
	lessonStatusExpired  = "EXPIRED"

	planStatusNone    = "NO_PLAN"
	planStatusRunning = "PLAN_RUNNING"
)

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
