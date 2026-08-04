// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"tjxt/apps/learning/api/internal/config"
	"tjxt/apps/learning/api/internal/model"
	"tjxt/apps/learning/api/internal/service"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	LessonModel   *model.LearningLessonModel
	LessonService service.LessonService
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DBConf.DataSource)
	lessonModel := model.NewLearningLessonModel(conn)
	return &ServiceContext{
		Config:        c,
		LessonModel:   lessonModel,
		LessonService: service.NewLessonService(lessonModel),
	}
}
