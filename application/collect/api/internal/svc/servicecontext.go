package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"zhifou/application/collect/api/internal/config"
	"zhifou/application/collect/rpc/collect"
	"zhifou/pkg/interceptors"
)

type ServiceContext struct {
	Config     config.Config
	CollectRPC collect.Collect
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	collectRPC := zrpc.MustNewClient(c.CollectRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	return &ServiceContext{
		Config:     c,
		CollectRPC: collect.NewCollect(collectRPC),
	}
}
