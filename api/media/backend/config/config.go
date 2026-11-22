// Package config media 后端业务配置（随原 rpc/media 业务层随迁），
// 进程壳配置（RpcServerConf/RestConf）由宿主进程自行承载。
package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
)

type Config struct {
	DataSource string
	Cache      cache.CacheConf

	// MinIO 对象存储（S3 兼容协议：RustFS / MinIO / COS / OSS）。
	// Endpoint 为空时回退 mock 实现，不阻塞本地启动。
	MinIO struct {
		Endpoint  string
		AccessKey string
		SecretKey string
		Bucket    string
		Secure    bool
	}
}
