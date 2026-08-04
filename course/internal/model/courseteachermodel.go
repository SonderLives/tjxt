package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseTeacherModel = (*customCourseTeacherModel)(nil)

type (
	// CourseTeacherModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseTeacherModel.
	CourseTeacherModel interface {
		courseTeacherModel
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
