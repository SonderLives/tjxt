package svc

import (
	"tjxt/apps/remark/api/internal/config"
	remarkclient "tjxt/apps/remark/rpc/client/remark"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	RemarkRpc remarkclient.Remark
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		RemarkRpc: remarkclient.NewRemark(zrpc.MustNewClient(c.RemarkRpc)),
	}
}