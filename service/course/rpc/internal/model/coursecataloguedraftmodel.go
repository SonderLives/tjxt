package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseCatalogueDraftModel = (*customCourseCatalogueDraftModel)(nil)

type (
	// CourseCatalogueDraftModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseCatalogueDraftModel.
	CourseCatalogueDraftModel interface {
		courseCatalogueDraftModel
		ListByCourseId(ctx context.Context, courseId int64) ([]*CourseCatalogueDraft, error)
		DeleteByCourseId(ctx context.Context, courseId int64) error
	}
	customCourseCatalogueDraftModel struct {
		*defaultCourseCatalogueDraftModel
	}
)

// NewCourseCatalogueDraftModel returns a model for the database table.
func NewCourseCatalogueDraftModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseCatalogueDraftModel {
	return &customCourseCatalogueDraftModel{
		defaultCourseCatalogueDraftModel: newCourseCatalogueDraftModel(conn, c, opts...),
	}
}

func (m *customCourseCatalogueDraftModel) ListByCourseId(ctx context.Context, courseId int64) ([]*CourseCatalogueDraft, error) {
	var list []*CourseCatalogueDraft
	query := fmt.Sprintf("select * from %s where `deleted` = 0 and `course_id` = ? order by `c_index` asc", m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, courseId); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCourseCatalogueDraftModel) DeleteByCourseId(ctx context.Context, courseId int64) error {
	query := fmt.Sprintf("delete from %s where `course_id` = ?", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, courseId)
	return err
}
