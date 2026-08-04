// Package service learning 业务层（高于 model 的薄封装，承载跨多张表的协作逻辑）。
package service

import (
	"context"
	"database/sql"
	"errors"

	"tjxt/apps/learning/rpc/internal/model"
	"tjxt/pkg/xerr"
)

// LearningService 学习域业务接口。
type LearningService interface {
	GrantCourses(ctx context.Context, userID int64, courseIDs []int64) error
	RevokeCourses(ctx context.Context, userID int64, courseIDs []int64) error
	RemoveLesson(ctx context.Context, userID, courseID int64) error
	GetLesson(ctx context.Context, userID, courseID int64) (*model.LearningLesson, error)
	ListLessons(ctx context.Context, userID, pageNo, pageSize int64, asc bool) ([]*model.LearningLesson, int64, error)
	CountLessons(ctx context.Context, courseID int64) (int64, error)
	CreatePlan(ctx context.Context, userID, courseID, weekFreq int64) error
	ValidateLesson(ctx context.Context, userID, courseID int64) (int64, error)
	CurrentLesson(ctx context.Context, userID int64) (*model.LearningLesson, error)
	CommitRecord(ctx context.Context, userID, lessonID, sectionID, moment, duration int64, commitTime string) error
	ListLessonPlans(ctx context.Context, userID, pageNo, pageSize int64, asc bool) ([]*model.LearningLesson, int64, error)
}

type learningService struct {
	m model.LearningLessonModel
}

// NewLearningService 业务服务构造函数。
func NewLearningService(m model.LearningLessonModel) LearningService { return &learningService{m: m} }

func (s *learningService) GrantCourses(ctx context.Context, userID int64, courseIDs []int64) error {
	if userID <= 0 || len(courseIDs) == 0 {
		return xerr.BadRequestf("userID/courseIDs 非法")
	}
	return s.m.GrantCourses(ctx, userID, courseIDs)
}

func (s *learningService) RevokeCourses(ctx context.Context, userID int64, courseIDs []int64) error {
	if userID <= 0 || len(courseIDs) == 0 {
		return xerr.BadRequestf("userID/courseIDs 非法")
	}
	return s.m.RevokeCourses(ctx, userID, courseIDs)
}

func (s *learningService) RemoveLesson(ctx context.Context, userID, courseID int64) error {
	if userID <= 0 || courseID <= 0 {
		return xerr.BadRequestf("userID/courseID 非法")
	}
	return s.m.RemoveLesson(ctx, userID, courseID)
}

func (s *learningService) GetLesson(ctx context.Context, userID, courseID int64) (*model.LearningLesson, error) {
	if userID <= 0 || courseID <= 0 {
		return nil, xerr.BadRequestf("userID/courseID 非法")
	}
	lesson, err := s.m.FindByUserCourse(ctx, userID, courseID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, xerr.NotFound("学习记录不存在")
	}
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询学习记录失败")
	}
	return lesson, nil
}

func (s *learningService) ListLessons(ctx context.Context, userID, pageNo, pageSize int64, asc bool) ([]*model.LearningLesson, int64, error) {
	if userID <= 0 {
		return nil, 0, xerr.BadRequestf("userID 非法")
	}
	return s.m.ListByUser(ctx, userID, pageNo, pageSize, asc)
}

func (s *learningService) CountLessons(ctx context.Context, courseID int64) (int64, error) {
	if courseID <= 0 {
		return 0, xerr.BadRequestf("courseID 非法")
	}
	return s.m.CountByCourse(ctx, courseID)
}

func (s *learningService) CreatePlan(ctx context.Context, userID, courseID, weekFreq int64) error {
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

func (s *learningService) ValidateLesson(ctx context.Context, userID, courseID int64) (int64, error) {
	lesson, err := s.GetLesson(ctx, userID, courseID)
	if err != nil {
		return 0, err
	}
	if lesson.Status == model.LessonStatusExpired {
		return 0, xerr.Conflict("课程已失效")
	}
	return lesson.Id, nil
}

// CurrentLesson 返回当前正在学习的课程（最近学习的一条）
func (s *learningService) CurrentLesson(ctx context.Context, userID int64) (*model.LearningLesson, error) {
	if userID <= 0 {
		return nil, xerr.BadRequestf("userID 非法")
	}
	return s.m.FindLatestLearnedByUser(ctx, userID)
}

// CommitRecord 提交学习记录：更新 lesson 的最新进度
func (s *learningService) CommitRecord(ctx context.Context, userID, lessonID, sectionID, moment, duration int64, commitTime string) error {
	if lessonID <= 0 || sectionID <= 0 {
		return xerr.BadRequestf("lessonID/sectionID 非法")
	}
	return s.m.UpdateLatestLearn(ctx, lessonID, sectionID, moment, duration)
}

func (s *learningService) ListLessonPlans(ctx context.Context, userID, pageNo, pageSize int64, asc bool) ([]*model.LearningLesson, int64, error) {
	if userID <= 0 {
		return nil, 0, xerr.BadRequestf("userID 非法")
	}
	return s.m.ListPlansByUser(ctx, userID, pageNo, pageSize, asc)
}