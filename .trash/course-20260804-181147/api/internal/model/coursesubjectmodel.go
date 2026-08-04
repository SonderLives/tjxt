package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseSubjectModel = (*customCourseSubjectModel)(nil)

type (
	// CourseSubjectModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseSubjectModel.
	CourseSubjectModel interface {
		courseSubjectModel
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
