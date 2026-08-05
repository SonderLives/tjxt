package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PrivilegeModel = (*customPrivilegeModel)(nil)

type (
	// PrivilegeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPrivilegeModel.
	PrivilegeModel interface {
		privilegeModel
	}

	customPrivilegeModel struct {
		*defaultPrivilegeModel
	}
)

// NewPrivilegeModel returns a model for the database table.
func NewPrivilegeModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PrivilegeModel {
	return &customPrivilegeModel{
		defaultPrivilegeModel: newPrivilegeModel(conn, c, opts...),
	}
}
