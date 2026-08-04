package svc

import (
	"tjxt/apps/learning/rpc/internal/config"
	"tjxt/apps/learning/rpc/internal/model"
	"tjxt/apps/learning/rpc/internal/service"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	LessonService service.LessonService
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DB.DataSource)
	lessonModel := model.NewLearningLessonModel(conn, c.Cache)
	return &ServiceContext{
		Config:        c,
		LessonService: service.NewLessonService(lessonModel),
	}
}