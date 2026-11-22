// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"tjxt/api/user/internal/config"
	userclient "tjxt/rpc/user/client/user"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	// UserRpc 用户服务的 RPC 客户端
	UserRpc userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
