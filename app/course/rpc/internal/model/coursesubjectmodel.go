package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseSubjectModel = (*customCourseSubjectModel)(nil)

type (
	// CourseSubjectModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseSubjectModel.
	CourseSubjectModel interface {
		courseSubjectModel
		ListByCourseId(ctx context.Context, courseId int64) ([]*CourseSubject, error)
		DeleteByCourseId(ctx context.Context, courseId int64) error
	}
	customCourseSubjectModel struct {
		*defaultCourseSubjectModel
	}
)

// NewCourseSubjectModel returns a model for the database table.
func NewCourseSubjectModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseSubjectModel {
	return &customCourseSubjectModel{
		defaultCourseSubjectModel: newCourseSubjectModel(conn, c, opts...),
	}
}

func (m *customCourseSubjectModel) ListByCourseId(ctx context.Context, courseId int64) ([]*CourseSubject, error) {
	var list []*CourseSubject
	query := fmt.Sprintf("select * from %s where `course_id` = ?", m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, courseId); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCourseSubjectModel) DeleteByCourseId(ctx context.Context, courseId int64) error {
	query := fmt.Sprintf("delete from %s where `course_id` = ?", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, courseId)
	return err
}
