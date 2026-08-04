package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PayChannelModel = (*customPayChannelModel)(nil)

type (
	// PayChannelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPayChannelModel.
	PayChannelModel interface {
		payChannelModel

		// FindAllEnabled 列出所有启用状态的渠道（按 channel_priority 升序）
		FindAllEnabled(ctx context.Context) ([]*PayChannel, error)
		// FindByCode 按渠道编码查询
		FindByCode(ctx context.Context, code string) (*PayChannel, error)
		// PageList 分页查询支付渠道
		PageList(ctx context.Context, name, channelCode string, status int64, offset, limit int64) ([]*PayChannel, int64, error)
	}

	customPayChannelModel struct {
		*defaultPayChannelModel
	}
)

// NewPayChannelModel returns a model for the database table.
func NewPayChannelModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PayChannelModel {
	return &customPayChannelModel{
		defaultPayChannelModel: newPayChannelModel(conn, c, opts...),
	}
}

// FindAllEnabled 列出所有启用状态的渠道
func (m *customPayChannelModel) FindAllEnabled(ctx context.Context) ([]*PayChannel, error) {
	var list []*PayChannel
	query := fmt.Sprintf("select %s from %s where `status` = 1 order by `channel_priority` asc", payChannelRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &list, query)
	return list, err
}

// FindByCode 按渠道编码查询
func (m *customPayChannelModel) FindByCode(ctx context.Context, code string) (*PayChannel, error) {
	var resp PayChannel
	query := fmt.Sprintf("select %s from %s where `channel_code` = ? limit 1", payChannelRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, code)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// PageList 分页查询支付渠道
func (m *customPayChannelModel) PageList(ctx context.Context, name, channelCode string, status int64, offset, limit int64) ([]*PayChannel, int64, error) {
	var (
		cond  []string
		args  []any
		list  []*PayChannel
		total int64
	)
	if name != "" {
		cond = append(cond, "`name` like ?")
		args = append(args, "%"+name+"%")
	}
	if channelCode != "" {
		cond = append(cond, "`channel_code` = ?")
		args = append(args, channelCode)
	}
	if status > 0 {
		cond = append(cond, "`status` = ?")
		args = append(args, status)
	}
	where := ""
	if len(cond) > 0 {
		where = " where " + strings.Join(cond, " and ")
	}

	countSQL := fmt.Sprintf("select count(1) from %s%s", m.table, where)
	if err := m.QueryRowNoCacheCtx(ctx, &total, countSQL, args...); err != nil {
		return nil, 0, err
	}

	listSQL := fmt.Sprintf("select %s from %s%s order by `channel_priority` asc limit ? offset ?", payChannelRows, m.table, where)
	args = append(args, limit, offset)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, listSQL, args...); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}