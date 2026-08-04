package svc

import (
	"fmt"

	"tjxt/apps/trade/rpc/internal/config"
	"tjxt/apps/trade/rpc/internal/model"
	"tjxt/pkg/mq"

	payclient "tjxt/apps/pay/rpc/pay"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	CartModel        model.CartModel
	OrderModel       model.OrderModel
	OrderDetailModel model.OrderDetailModel
	RefundApplyModel model.RefundApplyModel

	PayRpc      payclient.Pay
	MQProducer  *mq.Producer // 可为 nil：MQ 未就绪时不阻塞启动
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)

	svcCtx := &ServiceContext{
		Config:           c,
		CartModel:        model.NewCartModel(conn, c.Cache),
		OrderModel:       model.NewOrderModel(conn, c.Cache),
		OrderDetailModel: model.NewOrderDetailModel(conn, c.Cache),
		RefundApplyModel: model.NewRefundApplyModel(conn, c.Cache),
		PayRpc:           payclient.NewPay(zrpc.MustNewClient(c.PayRpc)),
	}

	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/", c.RabbitMQ.User, c.RabbitMQ.Pass, c.RabbitMQ.Host, c.RabbitMQ.Port)
	if prod, err := mq.NewProducer(dsn); err != nil {
		logx.Errorf("init rabbitmq producer failed, will skip event publish: %v", err)
	} else {
		svcCtx.MQProducer = prod
	}

	return svcCtx
}