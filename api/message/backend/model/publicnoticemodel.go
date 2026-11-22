package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PublicNoticeModel = (*customPublicNoticeModel)(nil)

type (
	// PublicNoticeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPublicNoticeModel.
	PublicNoticeModel interface {
		publicNoticeModel
		// FindCount 统计公告数量
		FindCount(ctx context.Context) (int64, error)
		// FindList 分页查询公告，按 push_time 倒序
		FindList(ctx context.Context, offset, limit int64) ([]*PublicNotice, error)
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

// FindCount 统计公告数量
func (m *customPublicNoticeModel) FindCount(ctx context.Context) (int64, error) {
	var total int64
	query := fmt.Sprintf("select count(1) from %s", m.table)
	if err := m.QueryRowNoCacheCtx(ctx, &total, query); err != nil {
		return 0, err
	}
	return total, nil
}

// FindList 分页查询公告，按 push_time 倒序
func (m *customPublicNoticeModel) FindList(ctx context.Context, offset, limit int64) ([]*PublicNotice, error) {
	var list []*PublicNotice
	query := fmt.Sprintf("select %s from %s order by `push_time` desc limit ? offset ?", publicNoticeRows, m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, limit, offset); err != nil {
		return nil, err
	}
	return list, nil
}
