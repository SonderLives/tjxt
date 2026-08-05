package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PublicNoticeModel = (*customPublicNoticeModel)(nil)

type (
	// PublicNoticeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPublicNoticeModel.
	PublicNoticeModel interface {
		publicNoticeModel
	}

	customPublicNoticeModel struct {
		*defaultPublicNoticeModel
	}
)

// NewPublicNoticeModel returns a model for the database table.
func NewPublicNoticeModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PublicNoticeModel {
	return &customPublicNoticeModel{
		defaultPublicNoticeModel: newPublicNoticeModel(conn, c, opts...),
	}
}
