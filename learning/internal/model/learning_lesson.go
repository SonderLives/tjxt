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