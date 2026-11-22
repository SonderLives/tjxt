package svc

import (
	"tjxt/api/exam/internal/config"
	examclient "tjxt/rpc/exam/client/exam"
	"tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	ExamRpc examclient.Exam
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		ExamRpc: examclient.NewExam(zrpc.MustNewClient(c.ExamRpc,
			zrpc.WithUnaryClientInterceptor(auth.OutgoingUserIdInterceptor))),
	}
}
