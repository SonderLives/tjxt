package config

import "github.com/zeromicro/go-zero/rest"

// Config 课程中心服务配置
type Config struct {
	rest.RestConf
	// 用于校验 JWT 的密钥与有效期（秒），与本仓其它服务保持一致
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	// 数据库连接（tj_course 库）
	DB struct {
		DataSource string
		DriverName string
	}
	// 用户详情表（tj_user 库，同实例跨库查询），用于老师信息与操作人姓名
	UserDetailTable string `json:",optional"`
	// 免费试看时长（分钟）
	TrailerDuration int64 `json:",optional"`
}
