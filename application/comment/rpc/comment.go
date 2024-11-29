package main

import (
	"flag"
	"fmt"

	"zhifou/application/comment/rpc/internal/config"
	"zhifou/application/comment/rpc/internal/server"
	"zhifou/application/comment/rpc/internal/svc"
	"zhifou/application/comment/rpc/service"

	"github.com/zeromicro/go-zero/core/conf"
	zs "github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/comment.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		service.RegisterCommentServer(grpcServer, server.NewCommentServer(ctx))

		if c.Mode == zs.DevMode || c.Mode == zs.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
