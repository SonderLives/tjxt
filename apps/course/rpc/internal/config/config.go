package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource string
	Cache      cache.CacheConf

	// RabbitMQ 配置，用于发布课程上架/下架事件（course.events 交换机），
	// 供 search 服务消费并增量同步 ES 课程索引。MQ 配置缺失或连接不可用时
	// 发布失败仅告警（best-effort），不阻塞课程主流程。
	RabbitMQ struct {
		Host string
		Port int
		User string
		Pass string
	}
}