// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	courseclient "tjxt/rpc/course/client/course"
	"tjxt/api/learning/internal/config"
	learningclient "tjxt/rpc/learning/client/learning"

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
		LearningRpc: learningclient.NewLearning(zrpc.MustNewClient(c.LearningRpc)),
		CourseRpc:   courseclient.NewCourse(zrpc.MustNewClient(c.CourseRpc)),
	}
}
