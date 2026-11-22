// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"tjxt/app/course/api/internal/config"
	courseclient "tjxt/app/course/rpc/client/course"
	learningclient "tjxt/app/learning/rpc/client/learning"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	// CourseRpc 课程域 RPC 客户端
	CourseRpc   courseclient.Course
	LearningRpc learningclient.Learning
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		CourseRpc:   courseclient.NewCourse(zrpc.MustNewClient(c.CourseRpc)),
		LearningRpc: learningclient.NewLearning(zrpc.MustNewClient(c.LearningRpc)),
	}
}
