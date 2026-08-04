package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseModel = (*customCourseModel)(nil)

type (
	// CourseModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseModel.
	CourseModel interface {
		courseModel
	}

	customCourseModel struct {
		*defaultCourseModel
	}
)

// NewCourseModel returns a model for the database table.
func NewCourseModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseModel {
	return &customCourseModel{
		defaultCourseModel: newCourseModel(conn, c, opts...),
	}
}
