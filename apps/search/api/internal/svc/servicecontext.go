package svc

import (
	"tjxt/apps/search/api/internal/config"
	searchclient "tjxt/apps/search/rpc/search"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	SearchRpc searchclient.Search
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		SearchRpc: searchclient.NewSearch(zrpc.MustNewClient(c.SearchRpc)),
	}
}
