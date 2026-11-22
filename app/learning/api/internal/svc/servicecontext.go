// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	courseclient "tjxt/app/course/rpc/client/course"
	"tjxt/app/learning/api/internal/config"
	learningclient "tjxt/app/learning/rpc/client/learning"
	"tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	LearningRpc learningclient.Learning
	CourseRpc   courseclient.Course
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		LearningRpc: learningclient.NewLearning(zrpc.MustNewClient(c.LearningRpc,
			// 将当前登录用户 id 透传到下游 RPC 的 gRPC metadata 中，
			// 供 learning-rpc 侧读取用户态（学习记录/计划等接口依赖）。
			zrpc.WithUnaryClientInterceptor(auth.OutgoingUserIdInterceptor))),
		CourseRpc:   courseclient.NewCourse(zrpc.MustNewClient(c.CourseRpc)),
	}
}
