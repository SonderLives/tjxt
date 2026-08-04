package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseContentModel = (*customCourseContentModel)(nil)

type (
	// CourseContentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseContentModel.
	CourseContentModel interface {
		courseContentModel
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
