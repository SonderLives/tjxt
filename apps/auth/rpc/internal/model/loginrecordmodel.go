package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LoginRecordModel = (*customLoginRecordModel)(nil)

type (
	// LoginRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLoginRecordModel.
	LoginRecordModel interface {
		loginRecordModel
	}

	customLoginRecordModel struct {
		*defaultLoginRecordModel
	}
)

// NewLoginRecordModel returns a model for the database table.
func NewLoginRecordModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LoginRecordModel {
	return &customLoginRecordModel{
		defaultLoginRecordModel: newLoginRecordModel(conn, c, opts...),
	}
}
