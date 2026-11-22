package main

import (
	"flag"
	"fmt"

	"tjxt/app/course/rpc/internal/config"
	courseserver "tjxt/app/course/rpc/internal/server/course"
	"tjxt/app/course/rpc/internal/svc"
	"tjxt/app/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"tjxt/pkg/confenv"
)

var configFile = flag.String("f", "etc/course.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	confenv.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterCourseServer(grpcServer, courseserver.NewCourseServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
