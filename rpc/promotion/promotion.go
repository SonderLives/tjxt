package main

import (
	"flag"
	"fmt"

	"tjxt/rpc/promotion/internal/config"
	promotionserver "tjxt/rpc/promotion/internal/server/promotion"
	"tjxt/rpc/promotion/internal/svc"
	"tjxt/rpc/promotion/pb"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"tjxt/pkg/confenv"
)

var configFile = flag.String("f", "etc/promotion.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	confenv.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterPromotionServer(grpcServer, promotionserver.NewPromotionServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
