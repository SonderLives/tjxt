package main

import (
	"context"
	"flag"
	"fmt"

	backendsvc "tjxt/app/trade/api/backend/svc"
	"tjxt/app/mq/trade-mq/config"
	"tjxt/pkg/confenv"
	"tjxt/pkg/mq"

	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/trade-mq.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	confenv.MustLoad(*configFile, &c)
	logx.Must(logx.SetUp(c.Log))

	ctx := backendsvc.NewServiceContext(c.Config)

	// MQ 消费者（pay.success → 订单回写已支付并转发 order.pay）。
	// Start 阻塞消费，连接断开自动重连；失败仅告警。
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/", c.RabbitMQ.User, c.RabbitMQ.Pass, c.RabbitMQ.Host, c.RabbitMQ.Port)
	client := mq.NewClient(dsn)
	registerMQ(client, ctx)

	fmt.Printf("Starting trade-mq consumer...\n")
	if err := client.Start(context.Background()); err != nil {
		logx.Errorf("trade mq consumer stopped: %v", err)
	}
}
