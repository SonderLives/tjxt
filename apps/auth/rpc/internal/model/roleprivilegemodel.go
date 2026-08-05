package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RolePrivilegeModel = (*customRolePrivilegeModel)(nil)

type (
	// RolePrivilegeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRolePrivilegeModel.
	RolePrivilegeModel interface {
		rolePrivilegeModel
	}

	customRolePrivilegeModel struct {
		*defaultRolePrivilegeModel
	}
)

// NewRolePrivilegeModel returns a model for the database table.
func NewRolePrivilegeModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RolePrivilegeModel {
	return &customRolePrivilegeModel{
		defaultRolePrivilegeModel: newRolePrivilegeModel(conn, c, opts...),
	}
}
