// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"fmt"

	"github.com/zeromicro/go-zero/rest"
)

// Config 交易中心服务配置
type Config struct {
	rest.RestConf
	RabbitMQConf
	DBConf
	RedisConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	// 下游服务（课程/用户）内部 HTTP 地址，用于跨服务数据获取
	CourseService HttpServiceConf
	UserService   HttpServiceConf
}

// HttpServiceConf 内部服务 HTTP 客户端配置
type HttpServiceConf struct {
	BaseURL string `json:",omitempty"` // 例如 http://127.0.0.1:8810
	Timeout int64  `json:",omitempty"` // 超时毫秒数，默认 3000
}

// RedisConf 缓存配置
type RedisConf struct {
	Host string
	Type string
	Pass string
}

// DBConf 数据库配置
type DBConf struct {
	// DataSource MySQL 连接串，例如
	// root:root@tcp(127.0.0.1:3306)/tj_trade?charset=utf8mb4&parseTime=true&loc=Local
	DataSource string
	DriverName string
}

// RabbitMQConf MQ 配置
type RabbitMQConf struct {
	Host string
	Port int
	User string
	Pass string
	// 事件交换机与路由键，与 learning 消费端契约一致
	Exchange         string
	PayRoutingKey    string
	RefundRoutingKey string
}

// DSN 返回 amqp 连接串，例如 amqp://guest:guest@127.0.0.1:5672/
func (r RabbitMQConf) DSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", r.User, r.Pass, r.Host, r.Port)
}
