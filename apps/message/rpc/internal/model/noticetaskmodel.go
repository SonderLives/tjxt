package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NoticeTaskModel = (*customNoticeTaskModel)(nil)

type (
	// NoticeTaskModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNoticeTaskModel.
	NoticeTaskModel interface {
		noticeTaskModel
	}

	customNoticeTaskModel struct {
		*defaultNoticeTaskModel
	}
)

// NewNoticeTaskModel returns a model for the database table.
func NewNoticeTaskModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) NoticeTaskModel {
	return &customNoticeTaskModel{
		defaultNoticeTaskModel: newNoticeTaskModel(conn, c, opts...),
	}
}
