package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseContentModel = (*customCourseContentModel)(nil)

type (
	// CourseContentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseContentModel.
	CourseContentModel interface {
		courseContentModel
		FindById(ctx context.Context, id int64) (*CourseContent, error)
	}

	customCourseContentModel struct {
		*defaultCourseContentModel
	}
)

// NewCourseContentModel returns a model for the database table.
func NewCourseContentModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseContentModel {
	return &customCourseContentModel{
		defaultCourseContentModel: newCourseContentModel(conn, c, opts...),
	}
}

// FindById 根据ID查询
func (m *customCourseContentModel) FindById(ctx context.Context, id int64) (*CourseContent, error) {
	return m.FindOne(ctx, id)
}
