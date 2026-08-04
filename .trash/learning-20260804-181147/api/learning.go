// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package main

import (
	"context"
	"flag"
	"fmt"

	"tjxt/apps/learning/api/internal/config"
	"tjxt/apps/learning/api/internal/consumer"
	"tjxt/apps/learning/api/internal/handler"
	"tjxt/apps/learning/api/internal/mq"
	"tjxt/apps/learning/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/learning-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 1. 创建服务上下文
	ctx := svc.NewServiceContext(c)

	// 2. 创建 MQ 客户端
	mqClient := mq.NewClient(c.RabbitMQConf.DSN())
	mqClient.SetPrefetch(10)

	// 3. 注册消费者
	payConsumer := consumer.NewLessonPayConsumer(ctx.LessonService)
	payConsumer.Register(mqClient, c.RabbitMQConf.PayQueue, c.RabbitMQConf.PayExchange, c.RabbitMQConf.PayRoutingKey)

	refundConsumer := consumer.NewLessonRefundConsumer(ctx.LessonService)
	refundConsumer.Register(mqClient, c.RabbitMQConf.RefundQueue, c.RabbitMQConf.RefundExchange, c.RabbitMQConf.RefundRoutingKey)

	// 4. 启动 MQ 消费（在后台运行）
	go func() {
		if err := mqClient.Start(context.Background()); err != nil {
			logx.Errorf("MQ client stopped: %v", err)
		}
	}()

	// 5. 启动 HTTP 服务器
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.RestConf.Host, c.RestConf.Port)
	server.Start()
}
