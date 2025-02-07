// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.2
// Source: collect.proto

package collect

import (
	"context"

	"zhifou/application/collect/rpc/service"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CollectItem         = service.CollectItem
	CollectListRequest  = service.CollectListRequest
	CollectListResponse = service.CollectListResponse
	CollectRequest      = service.CollectRequest
	CollectResponse     = service.CollectResponse
	UnCollectRequest    = service.UnCollectRequest
	UnCollectResponse   = service.UnCollectResponse

	Collect interface {
		// 收藏
		Collect(ctx context.Context, in *CollectRequest, opts ...grpc.CallOption) (*CollectResponse, error)
		// 取消收藏
		UnCollect(ctx context.Context, in *UnCollectRequest, opts ...grpc.CallOption) (*UnCollectResponse, error)
		// 收藏列表
		CollectList(ctx context.Context, in *CollectListRequest, opts ...grpc.CallOption) (*CollectListResponse, error)
	}

	defaultCollect struct {
		cli zrpc.Client
	}
)

func NewCollect(cli zrpc.Client) Collect {
	return &defaultCollect{
		cli: cli,
	}
}

// 收藏
func (m *defaultCollect) Collect(ctx context.Context, in *CollectRequest, opts ...grpc.CallOption) (*CollectResponse, error) {
	client := service.NewCollectClient(m.cli.Conn())
	return client.Collect(ctx, in, opts...)
}

// 取消收藏
func (m *defaultCollect) UnCollect(ctx context.Context, in *UnCollectRequest, opts ...grpc.CallOption) (*UnCollectResponse, error) {
	client := service.NewCollectClient(m.cli.Conn())
	return client.UnCollect(ctx, in, opts...)
}

// 收藏列表
func (m *defaultCollect) CollectList(ctx context.Context, in *CollectListRequest, opts ...grpc.CallOption) (*CollectListResponse, error) {
	client := service.NewCollectClient(m.cli.Conn())
	return client.CollectList(ctx, in, opts...)
}
