package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseContentDraftModel = (*customCourseContentDraftModel)(nil)

type (
	// CourseContentDraftModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseContentDraftModel.
	CourseContentDraftModel interface {
		courseContentDraftModel
	}

	customCourseContentDraftModel struct {
		*defaultCourseContentDraftModel
	}
)

// NewCourseContentDraftModel returns a model for the database table.
func NewCourseContentDraftModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseContentDraftModel {
	return &customCourseContentDraftModel{
		defaultCourseContentDraftModel: newCourseContentDraftModel(conn, c, opts...),
	}
}
