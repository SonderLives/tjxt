package model

import (
	"context"
	"fmt"

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
		DeleteByCourseId(ctx context.Context, courseId int64) error
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

func (m *customCourseTeacherDraftModel) ListByCourseId(ctx context.Context, courseId int64) ([]*CourseTeacherDraft, error) {
	var list []*CourseTeacherDraft
	query := fmt.Sprintf("select * from %s where `deleted` = 0 and `course_id` = ? order by `c_index` asc", m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, courseId); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCourseTeacherDraftModel) DeleteByCourseId(ctx context.Context, courseId int64) error {
	query := fmt.Sprintf("delete from %s where `course_id` = ?", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, courseId)
	return err
}
