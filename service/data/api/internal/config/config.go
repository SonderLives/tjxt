package config

import (
	"github.com/zeromicro/go-zero/rest"

	backendconfig "tjxt/service/data/api/backend/config"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	// data 后端业务配置（内联展开：Redis 大屏缓存），
	// 与原 rpc/data 配置合并（rpc 侧 Redis 挂在 RpcServerConf 下，此处改为
	// backend/config 的独立 redis.RedisConf，去除无用的 RedisKeyConf.Key）
	backendconfig.Config
}
