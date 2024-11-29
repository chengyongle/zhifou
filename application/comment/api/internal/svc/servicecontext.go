package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"zhifou/application/comment/api/internal/config"
	"zhifou/application/comment/rpc/comment"
	"zhifou/pkg/interceptors"
)

type ServiceContext struct {
	Config     config.Config
	CommentRPC comment.Comment
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	commentRPC := zrpc.MustNewClient(c.CommentRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))

	return &ServiceContext{
		Config:     c,
		CommentRPC: comment.NewComment(commentRPC),
	}
}
