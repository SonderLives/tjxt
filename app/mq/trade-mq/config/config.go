package config

import (
	"tjxt/app/trade/api/backend/config"

	"github.com/zeromicro/go-zero/core/logx"
)

// Config trade-mq 消费进程配置：由原 rpc/trade/internal/config 裁剪而来，
// 去除 RPC 服务端字段（ListenOn/Etcd/Prometheus/Telemetry）。
// 业务字段（DataSource/Cache/RabbitMQ/下游客户端）内嵌 backend/config，
// 与 backend/svc.ServiceContext 的整体构造保持一致。
type Config struct {
	Log logx.LogConf

	config.Config
}
