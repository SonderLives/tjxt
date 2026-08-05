// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"tjxt/apps/promotion/api/internal/config"
	promotionclient "tjxt/apps/promotion/rpc/promotion"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	// PromotionRpc 促销域 RPC 客户端
	PromotionRpc promotionclient.Promotion
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		PromotionRpc: promotionclient.NewPromotion(zrpc.MustNewClient(c.PromotionRpc)),
	}
}
