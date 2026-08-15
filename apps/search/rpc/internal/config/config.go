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
	// 本机 Docker 上的 Elasticsearch 9.4.1（127.0.0.1:9200）。
	// Analyzer：索引分词器。默认 cjk（ES 内置、无需插件，中文按二元分词，
	// 相关性优于 standard 的单字切分）；安装 IK/smartcn 插件后可改为
	// ik_max_word / smartcn 获得更佳中文分词。
	// SearchAnalyzer：检索分词器，缺省与 Analyzer 一致；IK 场景可设为
	// ik_smart（粗粒度）以提升检索性能与精度。
	Elasticsearch struct {
		Addresses     []string
		Username      string
		Password      string
		Analyzer      string
		SearchAnalyzer string
	}

	// RabbitMQ 配置，用于消费课程上下架事件同步 ES 索引
	RabbitMQ struct {
		Host string
		Port int
		User string
		Pass string
	}
}
