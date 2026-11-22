package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NoticeTemplateModel = (*customNoticeTemplateModel)(nil)

type (
	// NoticeTemplateModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNoticeTemplateModel.
	NoticeTemplateModel interface {
		noticeTemplateModel
		// FindCount 统计通知模板数量
		FindCount(ctx context.Context) (int64, error)
		// FindList 分页查询通知模板，按 update_time 倒序
		FindList(ctx context.Context, offset, limit int64) ([]*NoticeTemplate, error)
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

// FindCount 统计通知模板数量
func (m *customNoticeTemplateModel) FindCount(ctx context.Context) (int64, error) {
	var total int64
	query := fmt.Sprintf("select count(1) from %s", m.table)
	if err := m.QueryRowNoCacheCtx(ctx, &total, query); err != nil {
		return 0, err
	}
	return total, nil
}

// FindList 分页查询通知模板，按 update_time 倒序
func (m *customNoticeTemplateModel) FindList(ctx context.Context, offset, limit int64) ([]*NoticeTemplate, error) {
	var list []*NoticeTemplate
	query := fmt.Sprintf("select %s from %s order by `update_time` desc limit ? offset ?", noticeTemplateRows, m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, limit, offset); err != nil {
		return nil, err
	}
	return list, nil
}
