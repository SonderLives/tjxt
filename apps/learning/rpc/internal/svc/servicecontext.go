package svc

import (
	"tjxt/apps/learning/rpc/internal/config"
	"tjxt/apps/learning/rpc/internal/model"
	"tjxt/apps/learning/rpc/internal/service"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	LearningLessonModel model.LearningLessonModel
	LearningService     service.LearningService
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	lessonModel := model.NewLearningLessonModel(conn, c.Cache)
	return &ServiceContext{
		Config:              c,
		LearningLessonModel: lessonModel,
		LearningService:     service.NewLearningService(lessonModel),
	}
}