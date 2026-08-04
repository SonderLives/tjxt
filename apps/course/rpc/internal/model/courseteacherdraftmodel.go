package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseTeacherDraftModel = (*customCourseTeacherDraftModel)(nil)

type (
	// CourseTeacherDraftModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseTeacherDraftModel.
	CourseTeacherDraftModel interface {
		courseTeacherDraftModel
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
