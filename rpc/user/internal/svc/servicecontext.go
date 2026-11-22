package svc

import (
	"tjxt/apps/user/rpc/internal/config"
	"tjxt/apps/user/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	UserModel       model.UserModel
	UserDetailModel model.UserDetailModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config: c,

		UserModel:       model.NewUserModel(conn, c.Cache),
		UserDetailModel: model.NewUserDetailModel(conn, c.Cache),
	}
}
