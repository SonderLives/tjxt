package svc

import (
	"tjxt/api/remark/internal/config"
	remarkclient "tjxt/rpc/remark/client/remark"

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
