package svc

import (
	"fmt"

	"tjxt/apps/pay/rpc/internal/config"
	"tjxt/apps/pay/rpc/internal/gateway"
	"tjxt/apps/pay/rpc/internal/model"
	"tjxt/pkg/mq"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	PayChannelModel  model.PayChannelModel
	PayOrderModel    model.PayOrderModel
	RefundOrderModel model.RefundOrderModel

	// Gateway 支付渠道实现（mock 或真实渠道），logic 层经它与渠道交互
	Gateway    gateway.PaymentGateway
	MQProducer *mq.Producer // 可为 nil：MQ 未就绪时不阻塞启动
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	svcCtx := &ServiceContext{
		Config: c,

		PayChannelModel:  model.NewPayChannelModel(conn, c.Cache),
		PayOrderModel:    model.NewPayOrderModel(conn, c.Cache),
		RefundOrderModel: model.NewRefundOrderModel(conn, c.Cache),
	}

	// 支付渠道实现选择；未识别的 provider 回退 mock 并告警
	switch c.PayProvider {
	case "", "mock":
		svcCtx.Gateway = gateway.MockGateway{}
	default:
		logx.Errorf("unknown pay provider %q, fallback to mock gateway", c.PayProvider)
		svcCtx.Gateway = gateway.MockGateway{}
	}

	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/", c.RabbitMQ.User, c.RabbitMQ.Pass, c.RabbitMQ.Host, c.RabbitMQ.Port)
	if prod, err := mq.NewProducer(dsn); err != nil {
		logx.Errorf("init rabbitmq producer failed, will skip event publish: %v", err)
	} else {
		svcCtx.MQProducer = prod
	}

	return svcCtx
}
