// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"tjxt/api/trade/internal/config"
	tradeclient "tjxt/rpc/trade/client/trade"
	"tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	// TradeRpc 交易域 RPC 客户端
	TradeRpc tradeclient.Trade
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		TradeRpc: tradeclient.NewTrade(zrpc.MustNewClient(c.TradeRpc,
			// 将当前登录用户 id 透传到下游 RPC 的 gRPC metadata 中，
			// 供 trade-rpc 侧读取用户态（购物车/订单等接口依赖）。
			zrpc.WithUnaryClientInterceptor(auth.OutgoingUserIdInterceptor))),
	}
}
