package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NoticeTaskModel = (*customNoticeTaskModel)(nil)

type (
	// NoticeTaskModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNoticeTaskModel.
	NoticeTaskModel interface {
		noticeTaskModel
		// FindCount 统计通知任务数量
		FindCount(ctx context.Context) (int64, error)
		// FindList 分页查询通知任务，按 update_time 倒序
		FindList(ctx context.Context, offset, limit int64) ([]*NoticeTask, error)
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

// FindCount 统计通知任务数量
func (m *customNoticeTaskModel) FindCount(ctx context.Context) (int64, error) {
	var total int64
	query := fmt.Sprintf("select count(1) from %s", m.table)
	if err := m.QueryRowNoCacheCtx(ctx, &total, query); err != nil {
		return 0, err
	}
	return total, nil
}

// FindList 分页查询通知任务，按 update_time 倒序
func (m *customNoticeTaskModel) FindList(ctx context.Context, offset, limit int64) ([]*NoticeTask, error) {
	var list []*NoticeTask
	query := fmt.Sprintf("select %s from %s order by `update_time` desc limit ? offset ?", noticeTaskRows, m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, limit, offset); err != nil {
		return nil, err
	}
	return list, nil
}
