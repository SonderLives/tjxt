package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"learning/internal/pkg/idgen"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// LearningLesson 对应 learning_lesson 表（学生课程表）。
type LearningLesson struct {
	Id              int64
	UserId          int64
	CourseId        int64
	Status          int64
	WeekFreq        sql.NullInt64
	PlanStatus      int64
	LearnedSections int64
	LatestSectionId sql.NullInt64
	LatestLearnTime sql.NullTime
	CreateTime      time.Time
	ExpireTime      sql.NullTime
	UpdateTime      time.Time
}

type LearningLessonModel struct {
	conn  sqlx.SqlConn
	table string
}

func NewLearningLessonModel(conn sqlx.SqlConn) *LearningLessonModel {
	return &LearningLessonModel{conn: conn, table: "learning_lesson"}
}

// insertCourseSQL 幂等开课：已存在则仅刷新 update_time，不覆盖学习进度。
const insertCourseSQL = `INSERT INTO %s
	(id, user_id, course_id, status, week_freq, plan_status, learned_sections,
	 latest_section_id, latest_learn_time, create_time, expire_time, update_time)
VALUES (?, ?, ?, 0, NULL, 0, 0, NULL, NULL, NOW(), NULL, NOW())
ON DUPLICATE KEY UPDATE update_time = NOW()`

// GrantCourses 为用户开通所购课程，逐条写入并自动去重。
func (m *LearningLessonModel) GrantCourses(ctx context.Context, userId int64, courseIds []int64) error {
	for _, courseID := range courseIds {
		if _, err := m.conn.ExecCtx(
			ctx,
			fmt.Sprintf(insertCourseSQL, m.table),
			idgen.NextID(), userId, courseID,
		); err != nil {
			return err
		}
	}
	return nil
}

func (m *LearningLessonModel) FindByUserCourse(ctx context.Context, userID, courseID int64) (*LearningLesson, error) {
	var lesson LearningLesson
	err := m.conn.QueryRowCtx(ctx, &lesson, fmt.Sprintf("SELECT id, user_id, course_id, status, week_freq, plan_status, learned_sections, latest_section_id, latest_learn_time, create_time, expire_time, update_time FROM %s WHERE user_id = ? AND course_id = ?", m.table), userID, courseID)
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (m *LearningLessonModel) ListByUser(ctx context.Context, userID, pageNo, pageSize int64) ([]LearningLesson, int64, error) {
	var total int64
	if err := m.conn.QueryRowCtx(ctx, &total, fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE user_id = ? AND status <> 3", m.table), userID); err != nil {
		return nil, 0, err
	}
	var lessons []LearningLesson
	offset := (pageNo - 1) * pageSize
	err := m.conn.QueryRowsCtx(ctx, &lessons, fmt.Sprintf("SELECT id, user_id, course_id, status, week_freq, plan_status, learned_sections, latest_section_id, latest_learn_time, create_time, expire_time, update_time FROM %s WHERE user_id = ? AND status <> 3 ORDER BY create_time DESC LIMIT ? OFFSET ?", m.table), userID, pageSize, offset)
	return lessons, total, err
}

func (m *LearningLessonModel) UpdatePlan(ctx context.Context, userID, courseID int64, frequency int64) error {
	result, err := m.conn.ExecCtx(ctx, fmt.Sprintf("UPDATE %s SET week_freq = ?, plan_status = 1, update_time = NOW() WHERE user_id = ? AND course_id = ? AND status <> 3", m.table), frequency, userID, courseID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (m *LearningLessonModel) RevokeCourses(ctx context.Context, userID int64, courseIDs []int64) error {
	for _, courseID := range courseIDs {
		if _, err := m.conn.ExecCtx(ctx, fmt.Sprintf("UPDATE %s SET status = 3, plan_status = 0, update_time = NOW() WHERE user_id = ? AND course_id = ?", m.table), userID, courseID); err != nil {
			return err
		}
	}
	return nil
}

func (m *LearningLessonModel) CountByCourse(ctx context.Context, courseID int64) (int64, error) {
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE course_id = ? AND status <> 3", m.table), courseID)
	return count, err
}
