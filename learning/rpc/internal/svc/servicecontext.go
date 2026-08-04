package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"learning/internal/model"
	"learning/internal/service"
	"learning/rpc/internal/config"
)

type ServiceContext struct {
	Config        config.Config
	LessonService service.LessonService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{Config: c, LessonService: service.NewLessonService(model.NewLearningLessonModel(sqlx.NewMysql(c.DB.DataSource)))}
}
