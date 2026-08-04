package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config learning.rpc 配置
type Config struct {
	zrpc.RpcServerConf
	DataSource string
	Cache      cache.CacheConf

	// RabbitMQ 消费：订单支付/退款事件触达 learning 开通/撤销课程
	RabbitMQ struct {
		Host             string
		Port             int
		User             string
		Pass             string
		PayQueue         string
		RefundQueue      string
		PayExchange      string
		RefundExchange   string
		PayRoutingKey    string
		RefundRoutingKey string
	}
}