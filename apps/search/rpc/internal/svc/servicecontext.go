package svc

import (
	"tjxt/apps/search/rpc/internal/config"
	"tjxt/apps/search/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config         config.Config
	InterestsModel model.InterestsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:         c,
		InterestsModel: model.NewInterestsModel(conn, c.Cache),
	}
}
