// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"tjxt/service/trade/api/backend/config"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	// Backend 交易域 backend 业务配置（原 rpc/trade 配置随包并入：
	// DataSource/Cache/RabbitMQ 与下游 zrpc 客户端 PayRpc/CourseRpc/PromotionRpc，
	// 字段名保持原习惯，供 backend/svc.ServiceContext 整体构造）。
	config.Config
}
