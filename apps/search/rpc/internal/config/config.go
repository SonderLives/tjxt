package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource string
	Cache      cache.CacheConf

	// CourseRpc 课程服务客户端，用于上架事件回源课程索引数据
	CourseRpc zrpc.RpcClientConf

	// Elasticsearch 配置
	// 本机 Docker 上的 Elasticsearch 9.4.1（127.0.0.1:9200），未安装
	// analysis-smartcn 插件，Analyzer 固定为 standard（中文按单字切分）。
	Elasticsearch struct {
		Addresses []string
		Username  string
		Password  string
		Analyzer  string
	}

	// RabbitMQ 配置，用于消费课程上下架事件同步 ES 索引
	RabbitMQ struct {
		Host string
		Port int
		User string
		Pass string
	}
}
