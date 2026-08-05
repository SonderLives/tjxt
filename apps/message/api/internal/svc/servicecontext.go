package svc

import (
	"tjxt/apps/message/api/internal/config"
	messageclient "tjxt/apps/message/rpc/message"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	MessageRpc  messageclient.Message
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		MessageRpc: messageclient.NewMessage(zrpc.MustNewClient(c.MessageRpc)),
	}
}
