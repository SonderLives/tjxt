package svc

import (
	"tjxt/apps/data/api/data/internal/config"
	dataclient "tjxt/apps/data/rpc/data/client/data"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	DataRpc dataclient.Data
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		DataRpc: dataclient.NewData(zrpc.MustNewClient(c.DataRpc)),
	}
}