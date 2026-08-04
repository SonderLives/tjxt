package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RefundApplyModel = (*customRefundApplyModel)(nil)

type (
	// RefundApplyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRefundApplyModel.
	RefundApplyModel interface {
		refundApplyModel
	}

	customRefundApplyModel struct {
		*defaultRefundApplyModel
	}
)

// NewRefundApplyModel returns a model for the database table.
func NewRefundApplyModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RefundApplyModel {
	return &customRefundApplyModel{
		defaultRefundApplyModel: newRefundApplyModel(conn, c, opts...),
	}
}
