// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	courseclient "tjxt/apps/course/rpc/course"
	"tjxt/apps/learning/api/internal/config"
	learningclient "tjxt/apps/learning/rpc/learning"

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