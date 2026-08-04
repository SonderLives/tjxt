package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserDetailModel = (*customUserDetailModel)(nil)

type (
	// UserDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserDetailModel.
	UserDetailModel interface {
		userDetailModel
	}

	customUserDetailModel struct {
		*defaultUserDetailModel
	}
)

// NewUserDetailModel returns a model for the database table.
func NewUserDetailModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserDetailModel {
	return &customUserDetailModel{
		defaultUserDetailModel: newUserDetailModel(conn, c, opts...),
	}
}
