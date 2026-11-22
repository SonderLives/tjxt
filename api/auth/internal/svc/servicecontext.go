package svc

import (
	"tjxt/api/auth/internal/config"
	authclient "tjxt/rpc/auth/client/auth"
	userclient "tjxt/rpc/user/client/user"

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
