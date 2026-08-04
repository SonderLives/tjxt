package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource string
	Cache      cache.CacheConf

	// RabbitMQ 用于支付成功/退款事件发布到 learning
	RabbitMQ struct {
		Host             string
		Port             int
		User             string
		Pass             string
		Exchange         string
		PayRoutingKey    string
		RefundRoutingKey string
	}

	// PayRpc 调用 pay 服务（支付/退款下单、关单、支付/退款结果查询）
	PayRpc zrpc.RpcClientConf
}