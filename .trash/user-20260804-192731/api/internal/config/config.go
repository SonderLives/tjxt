package config

import "github.com/zeromicro/go-zero/rest"

// Config 用户中心服务配置
type Config struct {
	rest.RestConf

	// 用于签发/校验 JWT 的密钥与有效期（秒）
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	// 数据库连接（tj_user 库）
	DB struct {
		DataSource string
		DriverName string
	}

	// 员工角色名查询（auth 库 role 表，跨库同实例查询）
	RoleTable string `json:",omitempty"`

	// 新建用户时的默认密码（含重置密码）
	DefaultPassword string `json:",omitempty"`

	// Redis 缓存配置（可选，留空则不使用缓存）
	Cache struct {
		Host string
		Port int
	}

	// 内部接口 Token 白名单
	Internal struct {
		AccessTokens []string `json:",omitempty"`
	}
}
