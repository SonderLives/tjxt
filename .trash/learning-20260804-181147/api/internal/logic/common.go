package logic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"tjxt/apps/learning/api/internal/model"
	"tjxt/apps/learning/api/internal/types"
)

func currentUserID(ctx context.Context) (int64, error) {
	value := ctx.Value("userId")
	if value == nil {
		return 0, fmt.Errorf("missing authenticated user")
	}
	userID, err := strconv.ParseInt(fmt.Sprint(value), 10, 64)
	if err != nil || userID <= 0 {
		return 0, fmt.Errorf("invalid authenticated user")
	}
	return userID, nil
}

func normalizePage(pageNo, pageSize int64) (int64, int64) {
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return pageNo, pageSize
}

func lessonResponse(lesson *model.LearningLesson) types.Lesson {
	response := types.Lesson{Id: lesson.Id, CourseId: lesson.CourseId, Status: lesson.Status, PlanStatus: lesson.PlanStatus, LearnedSections: lesson.LearnedSections, CreateTime: lesson.CreateTime.Format(time.RFC3339)}
	if lesson.WeekFreq.Valid {
		response.WeekFreq = lesson.WeekFreq.Int64
	}
	if lesson.ExpireTime.Valid {
		response.ExpireTime = lesson.ExpireTime.Time.Format(time.RFC3339)
	}
	return response
}

func success(data any) *types.Result { return &types.Result{Code: 200, Msg: "OK", Data: data} }
