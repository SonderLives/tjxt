// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"fmt"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	RabbitMQConf
	DBConf
	RedisConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	LearningRpc zrpc.RpcClientConf
}

type RedisConf struct {
	Host string
	Type string
	Pass string
}

type DBConf struct {
	// DSN MySQL 连接串，例如
	// root:root@tcp(127.0.0.1:3306)/tj_learning?charset=utf8mb4&parseTime=true&loc=Local
	DataSource string
	DriverName string
}

type RabbitMQConf struct {
	Host             string
	Port             int
	User             string
	Pass             string
	PayQueue         string
	PayExchange      string
	PayRoutingKey    string
	RefundQueue      string
	RefundExchange   string
	RefundRoutingKey string
}

// DSN 返回 amqp 连接串，例如 amqp://guest:guest@127.0.0.1:5672/
func (r RabbitMQConf) DSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", r.User, r.Pass, r.Host, r.Port)
}
