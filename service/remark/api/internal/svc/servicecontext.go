package svc

import (
	"tjxt/api/remark/internal/config"

	backendsvc "tjxt/api/remark/backend/svc"
)

type ServiceContext struct {
	Config config.Config

	// Backend remark 后端服务上下文（点赞记录 Model）。原 remark RPC 进程折叠后，
	// 业务逻辑包随迁至 api/remark/backend，api logic 直接构造 backend logic
	// 在进程内调用，不再经 gRPC。
	Backend *backendsvc.ServiceContext
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Backend: backendsvc.NewServiceContext(c.Config),
	}
}
