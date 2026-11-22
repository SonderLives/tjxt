package config

import (
	"github.com/zeromicro/go-zero/rest"

	backendconfig "tjxt/api/remark/backend/config"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	// remark 后端业务配置（原 rpc/remark 配置随包并入：DataSource/Cache），
	// 与 backend/svc.ServiceContext 构造共用
	backendconfig.Config
}
