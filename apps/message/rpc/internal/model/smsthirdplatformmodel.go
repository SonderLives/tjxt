package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SmsThirdPlatformModel = (*customSmsThirdPlatformModel)(nil)

type (
	// SmsThirdPlatformModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSmsThirdPlatformModel.
	SmsThirdPlatformModel interface {
		smsThirdPlatformModel
		// FindAll 查询全部短信平台，按 priority 升序（数字越小优先级越高）
		FindAll(ctx context.Context) ([]*SmsThirdPlatform, error)
		// FindByCode 按平台代号查询
		FindByCode(ctx context.Context, code string) (*SmsThirdPlatform, error)
	}

	customSmsThirdPlatformModel struct {
		*defaultSmsThirdPlatformModel
	}
)

// NewSmsThirdPlatformModel returns a model for the database table.
func NewSmsThirdPlatformModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SmsThirdPlatformModel {
	return &customSmsThirdPlatformModel{
		defaultSmsThirdPlatformModel: newSmsThirdPlatformModel(conn, c, opts...),
	}
}

// FindAll 查询全部短信平台，按 priority 升序
func (m *customSmsThirdPlatformModel) FindAll(ctx context.Context) ([]*SmsThirdPlatform, error) {
	var list []*SmsThirdPlatform
	query := fmt.Sprintf("select %s from %s order by `priority` asc", smsThirdPlatformRows, m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query); err != nil {
		return nil, err
	}
	return list, nil
}

// FindByCode 按平台代号查询
func (m *customSmsThirdPlatformModel) FindByCode(ctx context.Context, code string) (*SmsThirdPlatform, error) {
	var resp SmsThirdPlatform
	query := fmt.Sprintf("select %s from %s where `code` = ? limit 1", smsThirdPlatformRows, m.table)
	if err := m.QueryRowNoCacheCtx(ctx, &resp, query, code); err != nil {
		return nil, err
	}
	return &resp, nil
}
