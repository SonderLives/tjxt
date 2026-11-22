// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	backendsvc "tjxt/service/trade/api/backend/svc"
	"tjxt/service/trade/api/internal/config"
)

type ServiceContext struct {
	Config config.Config

	// Backend 交易域本地 backend 服务上下文（DB/model/下游客户端/MQ producer）。
	// 原 trade RPC 进程折叠后，业务逻辑包随迁至 api/trade/backend，
	// api logic 直接构造 backend logic 在进程内调用，不再经 gRPC。
	Backend *backendsvc.ServiceContext
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Backend: backendsvc.NewServiceContext(c.Config),
	}
}
