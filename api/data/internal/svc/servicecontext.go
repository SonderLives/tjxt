package svc

import (
	"tjxt/api/data/internal/config"

	backendsvc "tjxt/api/data/backend/svc"
)

type ServiceContext struct {
	Config config.Config

	// Backend data 后端服务上下文（Redis 客户端），
	// api logic 直调 backend 包，不再走 data RPC
	Backend *backendsvc.ServiceContext
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Backend: backendsvc.NewServiceContext(c.Config),
	}
}
