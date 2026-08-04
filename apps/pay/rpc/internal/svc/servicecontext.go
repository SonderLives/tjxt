package svc

import (
	"tjxt/apps/pay/rpc/internal/config"
	"tjxt/apps/pay/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	PayChannelModel model.PayChannelModel
	PayOrderModel   model.PayOrderModel
	RefundOrderModel model.RefundOrderModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config: c,

		PayChannelModel:  model.NewPayChannelModel(conn, c.Cache),
		PayOrderModel:    model.NewPayOrderModel(conn, c.Cache),
		RefundOrderModel: model.NewRefundOrderModel(conn, c.Cache),
	}
}