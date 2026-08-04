// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package main

import (
	"flag"
	"fmt"

	"common/mq"
	"trade/internal/config"
	"trade/internal/handler"
	"trade/internal/service"
	"trade/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/trade-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 初始化 MQ 生产者，失败仅告警不阻塞启动（未就绪时相关发布会返回依赖不可用）
	var publisher service.EventPublisher
	producer, err := mq.NewProducer(c.RabbitMQConf.DSN())
	if err != nil {
		logx.Errorf("init mq producer failed: %v", err)
	} else {
		publisher = service.NewMQEventPublisher(producer, c.Exchange, c.PayRoutingKey, c.RefundRoutingKey)
		defer producer.Close()
	}

	ctx := svc.NewServiceContext(c, publisher)

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.RestConf.Host, c.RestConf.Port)
	server.Start()
}
