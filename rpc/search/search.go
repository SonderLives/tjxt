package main

import (
	"context"
	"flag"
	"fmt"

	"tjxt/rpc/search/internal/config"
	searchserver "tjxt/rpc/search/internal/server/search"
	"tjxt/rpc/search/internal/svc"
	"tjxt/rpc/search/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"tjxt/pkg/confenv"
)

var configFile = flag.String("f", "etc/search.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	confenv.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterSearchServer(grpcServer, searchserver.NewSearchServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// MQ 消费者（course.up/down → ES 索引同步）以 goroutine 运行，
	// 失败仅告警，不阻塞服务启动
	if ctx.MQClient != nil {
		go func() {
			if err := ctx.MQClient.Start(context.Background()); err != nil {
				logx.Errorf("search mq consumer stopped: %v", err)
			}
		}()
	}

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
