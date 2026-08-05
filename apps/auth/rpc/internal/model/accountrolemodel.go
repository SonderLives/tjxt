package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AccountRoleModel = (*customAccountRoleModel)(nil)

type (
	// AccountRoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAccountRoleModel.
	AccountRoleModel interface {
		accountRoleModel
	}

	customAccountRoleModel struct {
		*defaultAccountRoleModel
	}
)

// NewAccountRoleModel returns a model for the database table.
func NewAccountRoleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AccountRoleModel {
	return &customAccountRoleModel{
		defaultAccountRoleModel: newAccountRoleModel(conn, c, opts...),
	}
}
