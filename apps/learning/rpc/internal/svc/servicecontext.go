package svc

import (
	"fmt"

	"tjxt/apps/learning/rpc/internal/config"
	"tjxt/apps/learning/rpc/internal/model"
	"tjxt/apps/learning/rpc/internal/service"
	"tjxt/pkg/mq"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	LearningLessonModel model.LearningLessonModel
	LearningService     service.LearningService

	MQClient *mq.Client // 可为 nil：MQ 未就绪时跳过消费者启动
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	lessonModel := model.NewLearningLessonModel(conn, c.Cache)
	svcCtx := &ServiceContext{
		Config:              c,
		LearningLessonModel: lessonModel,
		LearningService:     service.NewLearningService(lessonModel),
	}

	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/", c.RabbitMQ.User, c.RabbitMQ.Pass, c.RabbitMQ.Host, c.RabbitMQ.Port)
	initMQ(svcCtx, dsn)

	return svcCtx
}
