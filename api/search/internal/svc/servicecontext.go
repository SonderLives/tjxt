package svc

import (
	"tjxt/api/search/internal/config"
	searchclient "tjxt/rpc/search/client/search"

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
