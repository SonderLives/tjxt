package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserInboxModel = (*customUserInboxModel)(nil)

type (
	// UserInboxModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserInboxModel.
	UserInboxModel interface {
		userInboxModel
	}

	customUserInboxModel struct {
		*defaultUserInboxModel
	}
)

// NewUserInboxModel returns a model for the database table.
func NewUserInboxModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserInboxModel {
	return &customUserInboxModel{
		defaultUserInboxModel: newUserInboxModel(conn, c, opts...),
	}
}
