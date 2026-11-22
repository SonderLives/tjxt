// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
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
