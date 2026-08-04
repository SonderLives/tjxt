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

	// LearningRpc 学习域 RPC 客户端（通过 etcd 服务发现）
	LearningRpc zrpc.RpcClientConf

	// CourseRpc 用于补全 LearningLessonVO 中 course_name/cover 等课程侧字段
	CourseRpc zrpc.RpcClientConf
}