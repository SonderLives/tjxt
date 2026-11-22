package svc

import (
	"tjxt/rpc/user/internal/config"
	"tjxt/rpc/user/internal/model"

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
