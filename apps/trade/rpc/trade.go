package main

import (
	"context"
	"flag"
	"fmt"

	"tjxt/apps/trade/rpc/internal/config"
	"tjxt/apps/trade/rpc/internal/server"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"tjxt/pkg/confenv"
)

var configFile = flag.String("f", "etc/trade.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	confenv.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterTradeServer(grpcServer, server.NewTradeServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// MQ 消费者（pay.success → 订单回写已支付并转发 order.pay）以 goroutine 运行，
	// 失败仅告警，不阻塞服务启动
	if ctx.MQClient != nil {
		go func() {
			if err := ctx.MQClient.Start(context.Background()); err != nil {
				logx.Errorf("trade mq consumer stopped: %v", err)
			}
		}()
	}

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
