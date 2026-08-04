package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseTeacherModel = (*customCourseTeacherModel)(nil)

type (
	// CourseTeacherModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseTeacherModel.
	CourseTeacherModel interface {
		courseTeacherModel
		ListByCourseId(ctx context.Context, courseId int64) ([]*CourseTeacher, error)
		CountTeacherCourse(ctx context.Context, teacherIds []int64) (map[int64]int64, error)
		SaveAll(ctx context.Context, courseId int64, list []*CourseTeacher) error
		DeleteByCourseId(ctx context.Context, courseId int64) error
	}

	customCourseTeacherModel struct {
		*defaultCourseTeacherModel
	}
)

// NewCourseTeacherModel returns a model for the database table.
func NewCourseTeacherModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseTeacherModel {
	return &customCourseTeacherModel{
		defaultCourseTeacherModel: newCourseTeacherModel(conn, c, opts...),
	}
}

// ListByCourseId 根据课程ID查询老师关系
func (m *customCourseTeacherModel) ListByCourseId(ctx context.Context, courseId int64) ([]*CourseTeacher, error) {
	var rows []*CourseTeacher
	query := fmt.Sprintf("select %s from %s where `course_id` = ? and `deleted` = 0", courseTeacherRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &rows, query, courseId)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// CountTeacherCourse 统计老师的课程数量
func (m *customCourseTeacherModel) CountTeacherCourse(ctx context.Context, teacherIds []int64) (map[int64]int64, error) {
	if len(teacherIds) == 0 {
		return map[int64]int64{}, nil
	}
	ph := make([]string, 0, len(teacherIds))
	args := make([]any, 0, len(teacherIds))
	for _, id := range teacherIds {
		ph = append(ph, "?")
		args = append(args, id)
	}
	query := fmt.Sprintf("select `teacher_id`, count(1) from %s where `teacher_id` in (%s) and `deleted` = 0 group by `teacher_id`", m.table, strings.Join(ph, ","))
	var rows []struct {
		TeacherId int64 `db:"teacher_id"`
		Count     int64 `db:"count"`
	}
	err := m.QueryRowsNoCacheCtx(ctx, &rows, query, args...)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]int64)
	for _, row := range rows {
		result[row.TeacherId] = row.Count
	}
	return result, nil
}

// SaveAll 全量保存课程老师关系
func (m *customCourseTeacherModel) SaveAll(ctx context.Context, courseId int64, list []*CourseTeacher) error {
	// 先删除旧数据
	if err := m.DeleteByCourseId(ctx, courseId); err != nil {
		return err
	}
	// 批量插入新数据
	for _, c := range list {
		if _, err := m.Insert(ctx, c); err != nil {
			return err
		}
	}
	return nil
}

// DeleteByCourseId 删除课程的所有老师关系
func (m *customCourseTeacherModel) DeleteByCourseId(ctx context.Context, courseId int64) error {
	courseTeacherIdKey := fmt.Sprintf("%s%v", cacheCourseTeacherIdPrefix, courseId)
	query := fmt.Sprintf("delete from %s where `course_id` = ?", m.table)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, courseId)
	}, courseTeacherIdKey)
	return err
}
