package config

import (
	"github.com/zeromicro/go-zero/rest"

	backendconfig "tjxt/api/search/backend/config"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	// search 后端业务配置（内联展开：DataSource/Cache/CourseRpc/Elasticsearch/RabbitMQ），
	// 与 app/mq/search-mq 消费者进程共用 backend svc
	backendconfig.Config
}
