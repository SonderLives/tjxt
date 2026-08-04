package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseCataSubjectDraftModel = (*customCourseCataSubjectDraftModel)(nil)

type (
	// CourseCataSubjectDraftModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseCataSubjectDraftModel.
	CourseCataSubjectDraftModel interface {
		courseCataSubjectDraftModel
	}

	customCourseCataSubjectDraftModel struct {
		*defaultCourseCataSubjectDraftModel
	}
)

// NewCourseCataSubjectDraftModel returns a model for the database table.
func NewCourseCataSubjectDraftModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseCataSubjectDraftModel {
	return &customCourseCataSubjectDraftModel{
		defaultCourseCataSubjectDraftModel: newCourseCataSubjectDraftModel(conn, c, opts...),
	}
}
