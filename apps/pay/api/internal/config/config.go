// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	// PayNotifySecret 支付回调签名密钥：sign = hex(hmac_sha256(secret, "payOrderNo=<no>"))
	PayNotifySecret string

	// MockPayEnabled 模拟支付接口开关：生产必须置 false（否则登录用户可绕过真实付款）
	MockPayEnabled bool

	// PayRpc 支付 RPC 客户端配置（通过 etcd 服务发现）
	PayRpc zrpc.RpcClientConf
}
