// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"github.com/zeromicro/go-zero/rest"

	backendconfig "tjxt/service/media/api/backend/config"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	// media 后端业务配置（内联展开：DataSource/Cache/MinIO 对象存储），
	// 与原 rpc/media 配置合并。MinIO（S3 兼容：RustFS/MinIO/COS/OSS）为
	// backend 特有键，与 api 侧 Auth.AccessSecret（JWT 校验）命名避让、互不覆盖。
	backendconfig.Config
}
