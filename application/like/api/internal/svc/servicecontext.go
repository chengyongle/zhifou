package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"zhifou/application/like/api/internal/config"
	"zhifou/application/like/rpc/like"
	"zhifou/pkg/interceptors"
)

type ServiceContext struct {
	Config  config.Config
	LikeRPC like.Like
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	likeRPC := zrpc.MustNewClient(c.LikeRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	return &ServiceContext{
		Config:  c,
		LikeRPC: like.NewLike(likeRPC),
	}
}
