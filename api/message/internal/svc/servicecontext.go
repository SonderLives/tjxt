package svc

import (
	"tjxt/api/message/internal/config"
	messageclient "tjxt/rpc/message/client/message"
	"tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	MessageRpc messageclient.Message
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		MessageRpc: messageclient.NewMessage(zrpc.MustNewClient(c.MessageRpc,
			zrpc.WithUnaryClientInterceptor(auth.OutgoingUserIdInterceptor))),
	}
}
