package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseTeacherDraftModel = (*customCourseTeacherDraftModel)(nil)

type (
	// CourseTeacherDraftModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseTeacherDraftModel.
	CourseTeacherDraftModel interface {
		courseTeacherDraftModel
		ListByCourseId(ctx context.Context, courseId int64) ([]*CourseTeacherDraft, error)
		SaveAll(ctx context.Context, courseId int64, list []*CourseTeacherDraft) error
		DeleteByCourseId(ctx context.Context, courseId int64) error
		CountTeacherCourse(ctx context.Context, teacherIds []int64) (map[int64]int64, error)
	}

	customCourseTeacherDraftModel struct {
		*defaultCourseTeacherDraftModel
	}
)

// NewCourseTeacherDraftModel returns a model for the database table.
func NewCourseTeacherDraftModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseTeacherDraftModel {
	return &customCourseTeacherDraftModel{
		defaultCourseTeacherDraftModel: newCourseTeacherDraftModel(conn, c, opts...),
	}
}

// ListByCourseId 根据课程ID查询老师关系草稿
func (m *customCourseTeacherDraftModel) ListByCourseId(ctx context.Context, courseId int64) ([]*CourseTeacherDraft, error) {
	var rows []*CourseTeacherDraft
	query := fmt.Sprintf("select %s from %s where `course_id` = ? and `deleted` = 0", courseTeacherDraftRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &rows, query, courseId)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SaveAll 全量保存课程老师关系草稿
func (m *customCourseTeacherDraftModel) SaveAll(ctx context.Context, courseId int64, list []*CourseTeacherDraft) error {
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

// DeleteByCourseId 删除课程的所有老师关系草稿
func (m *customCourseTeacherDraftModel) DeleteByCourseId(ctx context.Context, courseId int64) error {
	courseTeacherDraftIdKey := fmt.Sprintf("%s%v", cacheCourseTeacherDraftIdPrefix, courseId)
	query := fmt.Sprintf("delete from %s where `course_id` = ?", m.table)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, courseId)
	}, courseTeacherDraftIdKey)
	return err
}

// CountTeacherCourse 统计老师的课程数量
func (m *customCourseTeacherDraftModel) CountTeacherCourse(ctx context.Context, teacherIds []int64) (map[int64]int64, error) {
	if len(teacherIds) == 0 {
		return map[int64]int64{}, nil
	}
	placeholders := make([]string, 0, len(teacherIds))
	args := make([]any, 0, len(teacherIds))
	for _, id := range teacherIds {
		placeholders = append(placeholders, "?")
		args = append(args, id)
	}
	query := fmt.Sprintf("select `teacher_id`, count(1) from %s where `teacher_id` in (%s) and `deleted` = 0 group by `teacher_id`",
		m.table, strings.Join(placeholders, ", "))
	var rows []struct {
		TeacherId int64 `db:"teacher_id"`
		Cnt       int64 `db:"count(1)"`
	}
	err := m.QueryRowsNoCacheCtx(ctx, &rows, query, args...)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]int64, len(rows))
	for _, r := range rows {
		result[r.TeacherId] = r.Cnt
	}
	return result, nil
}
