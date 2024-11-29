package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"zhifou/application/follow/api/internal/config"
	"zhifou/application/follow/rpc/follow"
	"zhifou/pkg/interceptors"
)

type ServiceContext struct {
	Config    config.Config
	FollowRPC follow.Follow
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	followRPC := zrpc.MustNewClient(c.FollowRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	return &ServiceContext{
		Config:    c,
		FollowRPC: follow.NewFollow(followRPC),
	}
}
