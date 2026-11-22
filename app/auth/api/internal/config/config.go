package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	backendconfig "tjxt/service/auth/api/backend/config"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	// auth 后端业务配置（内联展开：DataSource/Cache/Jwt 签发密钥），
	// 与原 rpc/auth 配置合并。Jwt 密钥供 backend SignToken 签发令牌，
	// 与上方 Auth.AccessSecret（rest.WithJwt 校验用）用途独立、互不覆盖。
	backendconfig.Config

	UserRpc zrpc.RpcClientConf
}
