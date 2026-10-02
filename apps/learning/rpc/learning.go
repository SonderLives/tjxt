package main

import (
	"context"
	"flag"
	"fmt"

	"tjxt/apps/learning/rpc/internal/config"
	"tjxt/apps/learning/rpc/internal/server"
	"tjxt/apps/learning/rpc/internal/svc"
	"tjxt/apps/learning/rpc/pb"

	"tjxt/pkg/confenv"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/learning.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	confenv.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterLearningServer(grpcServer, server.NewLearningServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// MQ 消费者（order.pay/order.refund → 开课/撤课）以 goroutine 运行，
	// 失败仅告警，不阻塞服务启动
	if ctx.MQClient != nil {
		go func() {
			if err := ctx.MQClient.Start(context.Background()); err != nil {
				logx.Errorf("learning mq consumer stopped: %v", err)
			}
		}()
	}

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
