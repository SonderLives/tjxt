package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource string
	Cache      cache.CacheConf
	Jwt        struct {
		AccessSecret  string
		AccessExpire  int64 // 秒，访问令牌有效期
		RefreshSecret string
		RefreshExpire int64 // 秒,刷新令牌有效期
	}
}
