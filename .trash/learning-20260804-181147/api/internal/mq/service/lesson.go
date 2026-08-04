package service

import (
	"context"

	"tjxt/apps/learning/api/internal/model"

	"github.com/zeromicro/go-zero/core/logx"
)

// LessonService 课程服务接口
type LessonService interface {
	AddUserLessons(ctx context.Context, userID int64, courseIDs []int64) error
	DeleteCourseFromLesson(ctx context.Context, userID int64, courseID int64) error
}

// lessonService 课程服务实现
type lessonService struct {
	lessonModel *model.LearningLessonModel
}

// NewLessonService 创建课程服务
func NewLessonService(lessonModel *model.LearningLessonModel) LessonService {
	return &lessonService{
		lessonModel: lessonModel,
	}
}

// AddUserLessons 为用户添加课程（开课）
func (s *lessonService) AddUserLessons(ctx context.Context, userID int64, courseIDs []int64) error {
	if userID == 0 || len(courseIDs) == 0 {
		logx.Info("接受课程添加事件，但参数无效")
		return nil
	}

	if err := s.lessonModel.GrantCourses(ctx, userID, courseIDs); err != nil {
		logx.Errorf("授予课程失败, user_id=%d, courses=%v, err=%v", userID, courseIDs, err)
		return err
	}

	logx.Infof("授予课程成功, user_id=%d, courses=%v", userID, courseIDs)
	return nil
}

// DeleteCourseFromLesson 从用户课程中删除课程（退款）
func (s *lessonService) DeleteCourseFromLesson(ctx context.Context, userID int64, courseID int64) error {
	if userID == 0 || courseID == 0 {
		logx.Info("skip delete course from lesson, invalid params")
		return nil
	}

	// 这里可以实现软删除或硬删除逻辑
	// 简化版：标记状态为失效
	// 实际项目中可能需要更复杂的逻辑
	logx.Infof("delete course from lesson, user_id=%d, course_id=%d", userID, courseID)
	// TODO: 实现具体的删除逻辑
	return nil
}
