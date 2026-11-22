// Package config data 后端业务配置（随原 rpc/data 业务层随迁），
// 进程壳配置（RpcServerConf/RestConf）由宿主进程自行承载。
package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	// Redis 大屏数据（board/top10/today）读写缓存
	Redis redis.RedisConf
}
