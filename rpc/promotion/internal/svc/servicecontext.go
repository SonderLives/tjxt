package svc

import (
	"tjxt/apps/promotion/rpc/internal/config"
	"tjxt/apps/promotion/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	CouponModel     model.CouponModel
	UserCouponModel model.UserCouponModel
	CouponCodeModel model.CouponCodeModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)

	return &ServiceContext{
		Config:          c,
		CouponModel:     model.NewCouponModel(conn, c.Cache),
		UserCouponModel: model.NewUserCouponModel(conn, c.Cache),
		CouponCodeModel: model.NewCouponCodeModel(conn, c.Cache),
	}
}
