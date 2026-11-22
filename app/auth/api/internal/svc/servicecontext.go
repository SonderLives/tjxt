package svc

import (
	backendsvc "tjxt/app/auth/api/backend/svc"
	"tjxt/app/auth/api/internal/config"
	userclient "tjxt/app/user/rpc/client/user"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	// Backend auth 后端服务上下文（Model 全套 + Jwt 签发配置），
	// api logic 直调 backend 包，不再走 auth RPC
	Backend *backendsvc.ServiceContext

	UserRpc userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Backend: backendsvc.NewServiceContext(c.Config),
		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
