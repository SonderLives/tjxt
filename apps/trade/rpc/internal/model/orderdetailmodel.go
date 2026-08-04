package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderDetailModel = (*customOrderDetailModel)(nil)

type (
	// OrderDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderDetailModel.
	OrderDetailModel interface {
		orderDetailModel
	}

	customOrderDetailModel struct {
		*defaultOrderDetailModel
	}
)

// NewOrderDetailModel returns a model for the database table.
func NewOrderDetailModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderDetailModel {
	return &customOrderDetailModel{
		defaultOrderDetailModel: newOrderDetailModel(conn, c, opts...),
	}
}
