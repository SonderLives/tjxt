package svc

import (
	"tjxt/apps/exam/rpc/internal/config"
	"tjxt/apps/exam/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config              config.Config
	QuestionBizModel    model.QuestionBizModel
	QuestionDetailModel model.QuestionDetailModel
	QuestionModel       model.QuestionModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:              c,
		QuestionBizModel:    model.NewQuestionBizModel(conn, c.Cache),
		QuestionDetailModel: model.NewQuestionDetailModel(conn, c.Cache),
		QuestionModel:       model.NewQuestionModel(conn, c.Cache),
	}
}
