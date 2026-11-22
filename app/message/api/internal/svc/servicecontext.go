package svc

import (
	"tjxt/service/message/api/internal/config"

	backendsvc "tjxt/service/message/api/backend/svc"
)

type ServiceContext struct {
	Config config.Config

	// Backend message 后端服务上下文（通知模板/任务/收件箱/公告 Model 等）。
	// 原 message RPC 进程折叠后，业务逻辑包随迁至 api/message/backend，
	// api logic 直接构造 backend logic 在进程内调用，不再经 gRPC。
	Backend *backendsvc.ServiceContext
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Backend: backendsvc.NewServiceContext(c.Config),
	}
}
