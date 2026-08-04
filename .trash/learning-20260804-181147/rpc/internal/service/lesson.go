// Package service 提供 learning 域的业务服务层（在 Model 之上，rpc logic 之下）。
package service

import (
	"context"
	"database/sql"
	"errors"

	"tjxt/apps/learning/rpc/internal/model"
	"tjxt/pkg/xerr"
)

// LessonService 学习记录业务接口
type LessonService interface {
	GrantCourses(ctx context.Context, userID int64, courseIDs []int64) error
	RevokeCourses(ctx context.Context, userID int64, courseIDs []int64) error
	GetLesson(ctx context.Context, userID, courseID int64) (*model.LearningLesson, error)
	ListLessons(ctx context.Context, userID, pageNo, pageSize int64, asc bool) ([]*model.LearningLesson, int64, error)
	CountLessons(ctx context.Context, courseID int64) (int64, error)
	CreatePlan(ctx context.Context, userID, courseID, weekFreq int64) error
	RemoveLesson(ctx context.Context, userID, courseID int64) error
	ValidateLesson(ctx context.Context, userID, courseID int64) (int64, error)
}

type lessonService struct {
	m model.LearningLessonModel
}

// NewLessonService 业务服务构造函数
func NewLessonService(m model.LearningLessonModel) LessonService {
	return &lessonService{m: m}
}

func (s *lessonService) GrantCourses(ctx context.Context, userID int64, courseIDs []int64) error {
	if userID <= 0 || len(courseIDs) == 0 {
		return xerr.BadRequestf("userID/courseIDs 非法")
	}
	return s.m.GrantCourses(ctx, userID, courseIDs)
}

func (s *lessonService) RevokeCourses(ctx context.Context, userID int64, courseIDs []int64) error {
	if userID <= 0 || len(courseIDs) == 0 {
		return xerr.BadRequestf("userID/courseIDs 非法")
	}
	return s.m.RevokeCourses(ctx, userID, courseIDs)
}

func (s *lessonService) GetLesson(ctx context.Context, userID, courseID int64) (*model.LearningLesson, error) {
	if userID <= 0 || courseID <= 0 {
		return nil, xerr.BadRequestf("userID/courseID 非法")
	}
	lesson, err := s.m.FindByUserCourse(ctx, userID, courseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, xerr.NotFound("学习记录不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询学习记录失败")
	}
	return lesson, nil
}

func (s *lessonService) ListLessons(ctx context.Context, userID, pageNo, pageSize int64, asc bool) ([]*model.LearningLesson, int64, error) {
	if userID <= 0 {
		return nil, 0, xerr.BadRequestf("userID 非法")
	}
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	list, total, err := s.m.ListByUser(ctx, userID, pageNo, pageSize, asc)
	if err != nil {
		return nil, 0, xerr.Wrapf(err, xerr.CodeInternal, "查询学习记录列表失败")
	}
	return list, total, nil
}

func (s *lessonService) CountLessons(ctx context.Context, courseID int64) (int64, error) {
	if courseID <= 0 {
		return 0, xerr.BadRequestf("courseID 非法")
	}
	return s.m.CountByCourse(ctx, courseID)
}

func (s *lessonService) CreatePlan(ctx context.Context, userID, courseID, weekFreq int64) error {
	if userID <= 0 || courseID <= 0 || weekFreq <= 0 {
		return xerr.BadRequestf("userID/courseID/weekFreq 非法")
	}
	if err := s.m.UpdatePlan(ctx, userID, courseID, weekFreq); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return xerr.NotFound("学习记录不存在")
		}
		return xerr.Wrapf(err, xerr.CodeInternal, "更新学习计划失败")
	}
	return nil
}

func (s *lessonService) RemoveLesson(ctx context.Context, userID, courseID int64) error {
	if userID <= 0 || courseID <= 0 {
		return xerr.BadRequestf("userID/courseID 非法")
	}
	return s.m.RemoveLesson(ctx, userID, courseID)
}

// ValidateLesson 校验：存在/未过期则返回 lesson.Id；否则报 NotFound/Conflict
func (s *lessonService) ValidateLesson(ctx context.Context, userID, courseID int64) (int64, error) {
	lesson, err := s.GetLesson(ctx, userID, courseID)
	if err != nil {
		return 0, err
	}
	if lesson.Status == model.LessonStatusExpired {
		return 0, xerr.Conflict("课程已失效")
	}
	return lesson.Id, nil
}