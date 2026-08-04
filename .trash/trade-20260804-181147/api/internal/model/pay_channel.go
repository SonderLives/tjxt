package model

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// PayChannel 支付渠道（pay_channel 表，本地冗余用于渠道展示）
type PayChannel struct {
	Id              int64
	Name            string
	ChannelCode     string
	ChannelPriority int64
	ChannelIcon     string
	Status          int64
	Creater         int64
	Updater         int64
	CreateTime      time.Time
	UpdateTime      time.Time
}

// PayChannelModel 支付渠道数据访问
type PayChannelModel struct {
	conn  sqlx.SqlConn
	table string
}

// NewPayChannelModel 创建支付渠道数据访问对象
func NewPayChannelModel(conn sqlx.SqlConn) *PayChannelModel {
	return &PayChannelModel{conn: conn, table: "pay_channel"}
}

// ListEnabled 查询启用中的支付渠道。
func (m *PayChannelModel) ListEnabled(ctx context.Context) ([]PayChannel, error) {
	var list []PayChannel
	err := m.conn.QueryRowsCtx(ctx, &list, fmt.Sprintf(
		`SELECT id, name, channel_code, channel_priority, channel_icon, status, creater, updater,
		 create_time, update_time FROM %s WHERE status = 1 ORDER BY channel_priority ASC`, m.table))
	return list, err
}

// FindByCode 按渠道编码查询支付渠道。
func (m *PayChannelModel) FindByCode(ctx context.Context, code string) (*PayChannel, error) {
	var p PayChannel
	err := m.conn.QueryRowCtx(ctx, &p, fmt.Sprintf(
		`SELECT id, name, channel_code, channel_priority, channel_icon, status, creater, updater,
		 create_time, update_time FROM %s WHERE channel_code = ? AND status = 1 LIMIT 1`, m.table), code)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
