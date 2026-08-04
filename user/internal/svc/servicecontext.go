package svc

import (
	"user/internal/config"
	"user/internal/model"
	"user/internal/service"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	UserModel   *model.UserModel
	DetailModel *model.UserDetailModel
	AuthModel   *model.AuthModel

	AuthService service.AuthService
	UserService service.UserService
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DB.DataSource)

	userModel := model.NewUserModel(conn)
	detailModel := model.NewUserDetailModel(conn)
	authModel := model.NewAuthModel(conn, "tj_auth.login_record", c.RoleTable)

	return &ServiceContext{
		Config:      c,
		UserModel:   userModel,
		DetailModel: detailModel,
		AuthModel:   authModel,
		AuthService: service.NewAuthService(userModel, authModel, c.Auth.AccessSecret, c.Auth.AccessExpire),
		UserService: service.NewUserService(userModel, detailModel, authModel, c.DefaultPassword),
	}
}
