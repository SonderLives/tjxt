package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseDraftModel = (*customCourseDraftModel)(nil)

type (
	// CourseDraftModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseDraftModel.
	CourseDraftModel interface {
		courseDraftModel
	}

	customCourseDraftModel struct {
		*defaultCourseDraftModel
	}
)

// NewCourseDraftModel returns a model for the database table.
func NewCourseDraftModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseDraftModel {
	return &customCourseDraftModel{
		defaultCourseDraftModel: newCourseDraftModel(conn, c, opts...),
	}
}
