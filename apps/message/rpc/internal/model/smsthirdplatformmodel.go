package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SmsThirdPlatformModel = (*customSmsThirdPlatformModel)(nil)

type (
	// SmsThirdPlatformModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSmsThirdPlatformModel.
	SmsThirdPlatformModel interface {
		smsThirdPlatformModel
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
