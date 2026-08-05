package svc

import (
	"tjxt/apps/auth/api/internal/config"
	authclient "tjxt/apps/auth/rpc/client/auth"
	userclient "tjxt/apps/user/rpc/client/user"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	AuthRpc authclient.Auth
	UserRpc userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		AuthRpc: authclient.NewAuth(zrpc.MustNewClient(c.AuthRpc)),
		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}