package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseCatalogueModel = (*customCourseCatalogueModel)(nil)

type (
	// CourseCatalogueModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseCatalogueModel.
	CourseCatalogueModel interface {
		courseCatalogueModel
	}

	customCourseCatalogueModel struct {
		*defaultCourseCatalogueModel
	}
)

// NewCourseCatalogueModel returns a model for the database table.
func NewCourseCatalogueModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseCatalogueModel {
	return &customCourseCatalogueModel{
		defaultCourseCatalogueModel: newCourseCatalogueModel(conn, c, opts...),
	}
}
