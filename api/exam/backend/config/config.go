// Package config exam 后端业务配置（原 rpc/exam/internal/config 去除 RPC 服务端
// 字段后随包迁入；api 进程通过配置内嵌复用）。
package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
)

type Config struct {
	DataSource string
	Cache      cache.CacheConf
}
