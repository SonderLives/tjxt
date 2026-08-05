package svc

import (
	"tjxt/apps/exam/api/internal/config"
	examclient "tjxt/apps/exam/rpc/exam"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	ExamRpc examclient.Exam
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		ExamRpc: examclient.NewExam(zrpc.MustNewClient(c.ExamRpc)),
	}
}
