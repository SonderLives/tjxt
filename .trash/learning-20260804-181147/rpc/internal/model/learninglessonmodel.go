package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"tjxt/pkg/utils/idgen"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// 业务状态常量
const (
	LessonStatusNotStart  int64 = 0 // 未开始
	LessonStatusInLearn   int64 = 1 // 学习中
	LessonStatusDone      int64 = 2 // 完成
	LessonStatusExpired   int64 = 3 // 失效
	PlanStatusNone        int64 = 0 // 无计划
	PlanStatusInPlan      int64 = 1 // 计划中
)

var _ LearningLessonModel = (*customLearningLessonModel)(nil)

type (
	// LearningLessonModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLearningLessonModel.
	LearningLessonModel interface {
		learningLessonModel

		// GrantCourses 为用户开通课程（幂等，存在则不动）
		GrantCourses(ctx context.Context, userID int64, courseIDs []int64) error
		// RevokeCourses 撤销课程（status 置 3 = 失效）
		RevokeCourses(ctx context.Context, userID int64, courseIDs []int64) error
		// FindByUserCourse 查询用户在指定课程的学习记录
		FindByUserCourse(ctx context.Context, userID, courseID int64) (*LearningLesson, error)
		// ListByUser 分页查询用户的学习记录
		ListByUser(ctx context.Context, userID, pageNo, pageSize int64, asc bool) ([]*LearningLesson, int64, error)
		// CountByCourse 该课程的注册人数
		CountByCourse(ctx context.Context, courseID int64) (int64, error)
		// UpdatePlan 设置/更新学习计划（每周章节数）
		UpdatePlan(ctx context.Context, userID, courseID, weekFreq int64) error
		// RemoveLesson 删除该用户对这门课的学习记录（软删 = status 置 3）
		RemoveLesson(ctx context.Context, userID, courseID int64) error
	}

	customLearningLessonModel struct {
		*defaultLearningLessonModel
	}
)

// NewLearningLessonModel returns a model for the database table.
func NewLearningLessonModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LearningLessonModel {
	return &customLearningLessonModel{
		defaultLearningLessonModel: newLearningLessonModel(conn, c, opts...),
	}
}

// insertCourseSQL 幂等开课：已存在则仅刷新 update_time，不覆盖学习进度
const insertCourseSQL = `INSERT INTO %s
	(id, user_id, course_id, status, week_freq, plan_status, learned_sections,
	 latest_section_id, latest_learn_time, create_time, expire_time, update_time)
VALUES (?, ?, ?, 0, NULL, 0, 0, NULL, NULL, NOW(), NULL, NOW())
ON DUPLICATE KEY UPDATE update_time = NOW()`

func (m *customLearningLessonModel) GrantCourses(ctx context.Context, userID int64, courseIDs []int64) error {
	for _, cid := range courseIDs {
		if _, err := m.ExecNoCacheCtx(ctx, fmt.Sprintf(insertCourseSQL, m.table), idgen.NextID(), userID, cid); err != nil {
			return err
		}
	}
	return nil
}

func (m *customLearningLessonModel) RevokeCourses(ctx context.Context, userID int64, courseIDs []int64) error {
	for _, cid := range courseIDs {
		if _, err := m.ExecNoCacheCtx(ctx,
			fmt.Sprintf("UPDATE %s SET status = ?, plan_status = ?, update_time = NOW() WHERE user_id = ? AND course_id = ?", m.table),
			LessonStatusExpired, PlanStatusNone, userID, cid); err != nil {
			return err
		}
	}
	return nil
}

func (m *customLearningLessonModel) FindByUserCourse(ctx context.Context, userID, courseID int64) (*LearningLesson, error) {
	return m.FindOneByUserIdCourseId(ctx, userID, courseID)
}

func (m *customLearningLessonModel) ListByUser(ctx context.Context, userID, pageNo, pageSize int64, asc bool) ([]*LearningLesson, int64, error) {
	var total int64
	if err := m.QueryRowNoCacheCtx(ctx, &total,
		fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE user_id = ? AND status <> ?", m.table), userID, LessonStatusExpired); err != nil {
		return nil, 0, err
	}

	offset := (pageNo - 1) * pageSize
	order := "DESC"
	if asc {
		order = "ASC"
	}
	var list []*LearningLesson
	err := m.QueryRowsNoCacheCtx(ctx, &list,
		fmt.Sprintf("SELECT %s FROM %s WHERE user_id = ? AND status <> ? ORDER BY create_time %s LIMIT ? OFFSET ?",
			learningLessonRows, m.table, order),
		userID, LessonStatusExpired, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (m *customLearningLessonModel) CountByCourse(ctx context.Context, courseID int64) (int64, error) {
	var count int64
	err := m.QueryRowNoCacheCtx(ctx, &count,
		fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE course_id = ? AND status <> ?", m.table),
		courseID, LessonStatusExpired)
	return count, err
}

func (m *customLearningLessonModel) UpdatePlan(ctx context.Context, userID, courseID, weekFreq int64) error {
	res, err := m.ExecNoCacheCtx(ctx,
		fmt.Sprintf("UPDATE %s SET week_freq = ?, plan_status = ?, update_time = NOW() WHERE user_id = ? AND course_id = ? AND status <> ?", m.table),
		sql.NullInt64{Int64: weekFreq, Valid: weekFreq > 0}, PlanStatusInPlan, userID, courseID, LessonStatusExpired)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (m *customLearningLessonModel) RemoveLesson(ctx context.Context, userID, courseID int64) error {
	_, err := m.ExecNoCacheCtx(ctx,
		fmt.Sprintf("UPDATE %s SET status = ?, update_time = NOW() WHERE user_id = ? AND course_id = ?", m.table),
		LessonStatusExpired, userID, courseID)
	return err
}

// 占位避免误删（learningLessonRows 由 _gen.go 提供）
var _ = strings.TrimRight
var _ = time.Now