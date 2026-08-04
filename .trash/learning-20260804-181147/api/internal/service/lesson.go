package service

import (
	"context"
	"database/sql"

	"tjxt/apps/learning/api/internal/model"

	"github.com/zeromicro/go-zero/core/logx"
)

// LessonService 课程服务接口
type LessonService interface {
	AddUserLessons(ctx context.Context, userID int64, courseIDs []int64) error
	DeleteCourseFromLesson(ctx context.Context, userID int64, courseID int64) error
	GetLesson(ctx context.Context, userID, courseID int64) (*model.LearningLesson, error)
	ListLessons(ctx context.Context, userID, pageNo, pageSize int64) ([]model.LearningLesson, int64, error)
	CreatePlan(ctx context.Context, userID, courseID, frequency int64) error
	CountLessons(ctx context.Context, courseID int64) (int64, error)
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
		logx.Info("skip add user lessons, invalid params")
		return nil
	}

	if err := s.lessonModel.GrantCourses(ctx, userID, courseIDs); err != nil {
		logx.Errorf("grant courses failed, user_id=%d, courses=%v, err=%v", userID, courseIDs, err)
		return err
	}

	logx.Infof("add user lessons success, user_id=%d, courses=%v", userID, courseIDs)
	return nil
}

// DeleteCourseFromLesson 从用户课程中删除课程（退款）
func (s *lessonService) DeleteCourseFromLesson(ctx context.Context, userID int64, courseID int64) error {
	if userID == 0 || courseID == 0 {
		logx.Info("skip delete course from lesson, invalid params")
		return nil
	}

	return s.lessonModel.RevokeCourses(ctx, userID, []int64{courseID})
}

func (s *lessonService) GetLesson(ctx context.Context, userID, courseID int64) (*model.LearningLesson, error) {
	return s.lessonModel.FindByUserCourse(ctx, userID, courseID)
}
func (s *lessonService) ListLessons(ctx context.Context, userID, pageNo, pageSize int64) ([]model.LearningLesson, int64, error) {
	return s.lessonModel.ListByUser(ctx, userID, pageNo, pageSize)
}
func (s *lessonService) CreatePlan(ctx context.Context, userID, courseID, frequency int64) error {
	if frequency <= 0 {
		return sql.ErrNoRows
	}
	return s.lessonModel.UpdatePlan(ctx, userID, courseID, frequency)
}
func (s *lessonService) CountLessons(ctx context.Context, courseID int64) (int64, error) {
	return s.lessonModel.CountByCourse(ctx, courseID)
}
