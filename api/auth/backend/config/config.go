// Package config auth 后端业务配置（随原 rpc/auth 业务层随迁），
// 进程壳配置（RpcServerConf/RestConf）由宿主进程自行承载。
package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
)

type Config struct {
	DataSource string
	Cache      cache.CacheConf

	// Jwt 令牌签发密钥（backend SignToken 用）。
	// 注意与 api 侧 Auth.AccessSecret（rest.WithJwt 校验用）区分：
	// 两者通常配置同一个 TJXT_JWT_SECRET，但用途独立、互不覆盖。
	Jwt struct {
		AccessSecret  string
		AccessExpire  int64 // 秒，访问令牌有效期
		RefreshSecret string
		RefreshExpire int64 // 秒,刷新令牌有效期
	}
}
