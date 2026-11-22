package svc

import (
	"tjxt/app/search/api/internal/config"

	backendsvc "tjxt/app/search/api/backend/svc"
)

type ServiceContext struct {
	Config config.Config

	// Backend search 后端服务上下文（ES 客户端/兴趣 Model/course 客户端/MQ），
	// api logic 直调 backend 包，不再走 search RPC
	Backend *backendsvc.ServiceContext
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Backend: backendsvc.NewServiceContext(c.Config),
	}
}
