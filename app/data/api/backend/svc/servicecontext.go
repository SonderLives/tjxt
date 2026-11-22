package svc

import (
	"tjxt/app/data/api/backend/config"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config config.Config
	Rds    *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Rds:    redis.MustNewRedis(c.Redis),
	}
}
