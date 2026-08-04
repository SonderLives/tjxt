package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseCatalogueDraftModel = (*customCourseCatalogueDraftModel)(nil)

type (
	// CourseCatalogueDraftModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseCatalogueDraftModel.
	CourseCatalogueDraftModel interface {
		courseCatalogueDraftModel
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
