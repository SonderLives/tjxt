// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"learning/internal/config"
	"learning/internal/model"
	"learning/internal/service"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config       config.Config
	LessonModel  *model.LearningLessonModel
	LessonService service.LessonService
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DBConf.DSN)
	lessonModel := model.NewLearningLessonModel(conn)
	return &ServiceContext{
		Config:        c,
		LessonModel:   lessonModel,
		LessonService: service.NewLessonService(lessonModel),
	}
}