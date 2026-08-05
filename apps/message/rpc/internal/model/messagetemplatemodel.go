package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MessageTemplateModel = (*customMessageTemplateModel)(nil)

type (
	// MessageTemplateModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMessageTemplateModel.
	MessageTemplateModel interface {
		messageTemplateModel
	}

	customMessageTemplateModel struct {
		*defaultMessageTemplateModel
	}
)

// NewMessageTemplateModel returns a model for the database table.
func NewMessageTemplateModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) MessageTemplateModel {
	return &customMessageTemplateModel{
		defaultMessageTemplateModel: newMessageTemplateModel(conn, c, opts...),
	}
}
