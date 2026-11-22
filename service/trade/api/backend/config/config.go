package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config 交易 backend 业务配置（原 rpc/trade/internal/config 去除 RPC 服务端
// 字段后随包迁入；api 进程与 trade-mq 进程分别通过各自配置内嵌/展开复用）。
type Config struct {
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

	// CourseRpc 调用 course 服务（购物车/订单需要课程名、价格、封面）
	CourseRpc zrpc.RpcClientConf

	// PromotionRpc 调用 promotion 服务（优惠券可用方案、折扣计算、核销与退还）
	PromotionRpc zrpc.RpcClientConf
}
