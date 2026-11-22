// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"tjxt/api/media/internal/config"

	backendsvc "tjxt/api/media/backend/svc"
)

type ServiceContext struct {
	Config config.Config

	// Backend media 后端服务上下文（模型 + 对象存储客户端），
	// api logic 直调 backend 包，不再走 media RPC
	Backend *backendsvc.ServiceContext
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Backend: backendsvc.NewServiceContext(c.Config),
	}
}
