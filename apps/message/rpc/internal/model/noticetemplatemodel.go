package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NoticeTemplateModel = (*customNoticeTemplateModel)(nil)

type (
	// NoticeTemplateModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNoticeTemplateModel.
	NoticeTemplateModel interface {
		noticeTemplateModel
	}

	customNoticeTemplateModel struct {
		*defaultNoticeTemplateModel
	}
)

// NewNoticeTemplateModel returns a model for the database table.
func NewNoticeTemplateModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) NoticeTemplateModel {
	return &customNoticeTemplateModel{
		defaultNoticeTemplateModel: newNoticeTemplateModel(conn, c, opts...),
	}
}
