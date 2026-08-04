// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"tjxt/apps/trade/api/internal/config"
	tradeclient "tjxt/apps/trade/rpc/trade"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	// TradeRpc 交易域 RPC 客户端
	TradeRpc tradeclient.Trade
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		TradeRpc: tradeclient.NewTrade(zrpc.MustNewClient(c.TradeRpc)),
	}
}