package config

import (
	"github.com/zeromicro/go-zero/core/service"

	backendconfig "tjxt/api/search/backend/config"
)

// Config search-mq 消费者进程配置：进程壳（Name/Log/Prometheus/Telemetry）
// + search 后端业务配置（DataSource/Cache/CourseRpc/Elasticsearch/RabbitMQ）。
type Config struct {
	service.ServiceConf
	backendconfig.Config
}
