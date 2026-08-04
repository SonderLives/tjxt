package model

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// CourseTeacherDraft 课程老师关系草稿表（course_teacher_draft 表）
type CourseTeacherDraft struct {
	Id        int64
	CourseId  int64
	TeacherId int64
	IsShow    int64
	CIndex    int64
	DepId     int64
	CreateTime time.Time
	UpdateTime time.Time
	Creater   int64
	Updater   int64
	Deleted   int64
}

// CourseTeacherDraftModel 课程老师关系草稿数据访问
type CourseTeacherDraftModel struct {
	conn  sqlx.SqlConn
	table string
}

// NewCourseTeacherDraftModel 创建课程老师关系草稿数据访问对象
func NewCourseTeacherDraftModel(conn sqlx.SqlConn) *CourseTeacherDraftModel {
	return &CourseTeacherDraftModel{conn: conn, table: "course_teacher_draft"}
}

const courseTeacherDraftColumns = `id, course_id, teacher_id, is_show, c_index, dep_id, create_time,
	update_time, creater, updater, deleted`

// ListByCourseId 按课程查询草稿老师关系（按 c_index 排序）。
func (m *CourseTeacherDraftModel) ListByCourseId(ctx context.Context, courseId int64) ([]*CourseTeacherDraft, error) {
	var rows []*CourseTeacherDraft
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE course_id = ? AND deleted = 0 ORDER BY c_index`,
		courseTeacherDraftColumns, m.table), courseId)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// DeleteByCourseId 逻辑删除草稿老师关系。
func (m *CourseTeacherDraftModel) DeleteByCourseId(ctx context.Context, courseId int64) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET deleted = 1, update_time = NOW() WHERE course_id = ? AND deleted = 0`, m.table), courseId)
	return err
}

// ReplaceAll 全删重插。
func (m *CourseTeacherDraftModel) ReplaceAll(ctx context.Context, courseId int64, list []*CourseTeacherDraft) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := session.ExecCtx(ctx, fmt.Sprintf(
			`UPDATE %s SET deleted = 1, update_time = NOW() WHERE course_id = ? AND deleted = 0`,
			m.table), courseId); err != nil {
			return err
		}
		for _, c := range list {
			if _, err := session.ExecCtx(ctx, fmt.Sprintf(
				`INSERT INTO %s (id, course_id, teacher_id, is_show, c_index, dep_id, create_time,
					update_time, creater, updater) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW(), ?, ?)`,
				m.table),
				c.Id, c.CourseId, c.TeacherId, c.IsShow, c.CIndex, c.DepId, c.Creater, c.Updater); err != nil {
				return err
			}
		}
		return nil
	})
}

// CountTeacherCourse 统计老师名下草稿课程数量。
func (m *CourseTeacherDraftModel) CountTeacherCourse(ctx context.Context, teacherIds []int64) (map[int64]int64, error) {
	result := make(map[int64]int64)
	if len(teacherIds) == 0 {
		return result, nil
	}
	ph := make([]string, 0, len(teacherIds))
	args := make([]any, 0, len(teacherIds))
	for _, id := range teacherIds {
		ph = append(ph, "?")
		args = append(args, id)
	}
	var rows []struct {
		TeacherId int64 `db:"teacher_id"`
		Num       int64 `db:"num"`
	}
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT teacher_id, COUNT(*) AS num FROM %s WHERE deleted = 0 AND teacher_id IN (%s) GROUP BY teacher_id`,
		m.table, strings.Join(ph, ",")), args...)
	if err != nil {
		return result, err
	}
	for _, r := range rows {
		result[r.TeacherId] = r.Num
	}
	return result, nil
}
