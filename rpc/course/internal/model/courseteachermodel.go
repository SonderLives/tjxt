package model

import (
	"context"
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
		ListByTeacherIds(ctx context.Context, teacherIds []int64) ([]*CourseTeacher, error)
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

func (m *customCourseTeacherModel) ListByCourseId(ctx context.Context, courseId int64) ([]*CourseTeacher, error) {
	var list []*CourseTeacher
	query := fmt.Sprintf("select * from %s where `deleted` = 0 and `course_id` = ? order by `c_index` asc", m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, courseId); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCourseTeacherModel) ListByTeacherIds(ctx context.Context, teacherIds []int64) ([]*CourseTeacher, error) {
	if len(teacherIds) == 0 {
		return []*CourseTeacher{}, nil
	}
	var list []*CourseTeacher
	ph := strings.TrimSuffix(strings.Repeat("?,", len(teacherIds)), ",")
	args := make([]any, 0, len(teacherIds))
	for _, id := range teacherIds {
		args = append(args, id)
	}
	query := fmt.Sprintf("select * from %s where `deleted` = 0 and `teacher_id` in (%s)", m.table, ph)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCourseTeacherModel) DeleteByCourseId(ctx context.Context, courseId int64) error {
	query := fmt.Sprintf("delete from %s where `course_id` = ?", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, courseId)
	return err
}
