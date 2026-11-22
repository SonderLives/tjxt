package svc

import (
	"tjxt/app/exam/api/internal/config"

	backendsvc "tjxt/app/exam/api/backend/svc"
)

type ServiceContext struct {
	Config config.Config

	// Backend exam 后端服务上下文（题目 Model 等）。原 exam RPC 进程折叠后，
	// 业务逻辑包随迁至 api/exam/backend，api logic 直接构造 backend logic
	// 在进程内调用，不再经 gRPC。
	Backend *backendsvc.ServiceContext
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Backend: backendsvc.NewServiceContext(c.Config),
	}
}
