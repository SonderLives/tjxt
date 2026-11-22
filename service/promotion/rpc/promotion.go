package main

import (
	"flag"
	"fmt"

	"tjxt/service/promotion/rpc/internal/config"
	promotionserver "tjxt/service/promotion/rpc/internal/server/promotion"
	"tjxt/service/promotion/rpc/internal/svc"
	"tjxt/service/promotion/rpc/pb"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"tjxt/common/confenv"
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
