package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource  string
	Cache       cache.CacheConf
	TablePrefix string

	// PayProvider 支付渠道实现：mock（默认）| 预留 wechat/alipay 等真实渠道
	PayProvider string

	// RabbitMQ 用于支付成功事件发布到 trade（pay.exchange / pay.success）
	RabbitMQ struct {
		Host string
		Port int
		User string
		Pass string
	}
}
