// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"tjxt/api/media/internal/config"
	mediaclient "tjxt/rpc/media/client/media"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	// MediaRpc 媒资服务的 RPC 客户端
	MediaRpc mediaclient.Media
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		MediaRpc: mediaclient.NewMedia(zrpc.MustNewClient(c.MediaRpc)),
	}
}
