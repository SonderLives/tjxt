package svc

import (
	"tjxt/apps/auth/rpc/internal/config"
	"tjxt/apps/auth/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config             config.Config
	AccountRoleModel   model.AccountRoleModel
	LoginRecordModel   model.LoginRecordModel
	MenuModel          model.MenuModel
	PrivilegeModel     model.PrivilegeModel
	RoleMenuModel      model.RoleMenuModel
	RoleModel          model.RoleModel
	RolePrivilegeModel model.RolePrivilegeModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:             c,
		AccountRoleModel:   model.NewAccountRoleModel(conn, c.Cache),
		LoginRecordModel:   model.NewLoginRecordModel(conn, c.Cache),
		MenuModel:          model.NewMenuModel(conn, c.Cache),
		PrivilegeModel:     model.NewPrivilegeModel(conn, c.Cache),
		RoleMenuModel:      model.NewRoleMenuModel(conn, c.Cache),
		RoleModel:          model.NewRoleModel(conn, c.Cache),
		RolePrivilegeModel: model.NewRolePrivilegeModel(conn, c.Cache),
	}
}
