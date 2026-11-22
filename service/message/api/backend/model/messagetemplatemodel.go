package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MessageTemplateModel = (*customMessageTemplateModel)(nil)

type (
	// MessageTemplateModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMessageTemplateModel.
	MessageTemplateModel interface {
		messageTemplateModel
		// FindCount 统计短信模板数量
		FindCount(ctx context.Context) (int64, error)
		// FindList 分页查询短信模板，按 update_time 倒序
		FindList(ctx context.Context, offset, limit int64) ([]*MessageTemplate, error)
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

// FindCount 统计短信模板数量
func (m *customMessageTemplateModel) FindCount(ctx context.Context) (int64, error) {
	var total int64
	query := fmt.Sprintf("select count(1) from %s", m.table)
	if err := m.QueryRowNoCacheCtx(ctx, &total, query); err != nil {
		return 0, err
	}
	return total, nil
}

// FindList 分页查询短信模板，按 update_time 倒序
func (m *customMessageTemplateModel) FindList(ctx context.Context, offset, limit int64) ([]*MessageTemplate, error) {
	var list []*MessageTemplate
	query := fmt.Sprintf("select %s from %s order by `update_time` desc limit ? offset ?", messageTemplateRows, m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, limit, offset); err != nil {
		return nil, err
	}
	return list, nil
}
