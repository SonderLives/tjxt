package svc

import (
	"fmt"

	"tjxt/apps/trade/rpc/internal/config"
	"tjxt/apps/trade/rpc/internal/model"
	"tjxt/pkg/mq"

	courseclient "tjxt/apps/course/rpc/course"
	payclient "tjxt/apps/pay/rpc/pay"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	// DB 原始连接：跨模型事务（如下单的 order+order_detail 原子写入）用它开事务，
	// 新插入的行尚无缓存条目，事务内直接 Exec、无需额外失效缓存
	DB sqlx.SqlConn

	CartModel        model.CartModel
	OrderModel       model.OrderModel
	OrderDetailModel model.OrderDetailModel
	RefundApplyModel model.RefundApplyModel

	PayRpc      payclient.Pay
	CourseRpc   courseclient.Course
	MQProducer  *mq.Producer // 可为 nil：MQ 未就绪时不阻塞启动
	MQClient    *mq.Client   // 可为 nil：MQ 未就绪时跳过消费者启动
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)

	svcCtx := &ServiceContext{
		Config:           c,
		DB:               conn,
		CartModel:        model.NewCartModel(conn, c.Cache),
		OrderModel:       model.NewOrderModel(conn, c.Cache),
		OrderDetailModel: model.NewOrderDetailModel(conn, c.Cache),
		RefundApplyModel: model.NewRefundApplyModel(conn, c.Cache),
		PayRpc:           payclient.NewPay(zrpc.MustNewClient(c.PayRpc)),
		CourseRpc:        courseclient.NewCourse(zrpc.MustNewClient(c.CourseRpc)),
	}

	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/", c.RabbitMQ.User, c.RabbitMQ.Pass, c.RabbitMQ.Host, c.RabbitMQ.Port)
	if prod, err := mq.NewProducer(dsn); err != nil {
		logx.Errorf("init rabbitmq producer failed, will skip event publish: %v", err)
	} else {
		svcCtx.MQProducer = prod
	}

	initMQ(svcCtx, dsn)

	return svcCtx
}