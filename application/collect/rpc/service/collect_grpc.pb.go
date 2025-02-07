// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: collect.proto

package service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CollectClient is the client API for Collect service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CollectClient interface {
	// 收藏
	Collect(ctx context.Context, in *CollectRequest, opts ...grpc.CallOption) (*CollectResponse, error)
	// 取消收藏
	UnCollect(ctx context.Context, in *UnCollectRequest, opts ...grpc.CallOption) (*UnCollectResponse, error)
	// 收藏列表
	CollectList(ctx context.Context, in *CollectListRequest, opts ...grpc.CallOption) (*CollectListResponse, error)
}

type collectClient struct {
	cc grpc.ClientConnInterface
}

func NewCollectClient(cc grpc.ClientConnInterface) CollectClient {
	return &collectClient{cc}
}

func (c *collectClient) Collect(ctx context.Context, in *CollectRequest, opts ...grpc.CallOption) (*CollectResponse, error) {
	out := new(CollectResponse)
	err := c.cc.Invoke(ctx, "/service.Collect/Collect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectClient) UnCollect(ctx context.Context, in *UnCollectRequest, opts ...grpc.CallOption) (*UnCollectResponse, error) {
	out := new(UnCollectResponse)
	err := c.cc.Invoke(ctx, "/service.Collect/UnCollect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectClient) CollectList(ctx context.Context, in *CollectListRequest, opts ...grpc.CallOption) (*CollectListResponse, error) {
	out := new(CollectListResponse)
	err := c.cc.Invoke(ctx, "/service.Collect/CollectList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CollectServer is the server API for Collect service.
// All implementations must embed UnimplementedCollectServer
// for forward compatibility
type CollectServer interface {
	// 收藏
	Collect(context.Context, *CollectRequest) (*CollectResponse, error)
	// 取消收藏
	UnCollect(context.Context, *UnCollectRequest) (*UnCollectResponse, error)
	// 收藏列表
	CollectList(context.Context, *CollectListRequest) (*CollectListResponse, error)
	mustEmbedUnimplementedCollectServer()
}

// UnimplementedCollectServer must be embedded to have forward compatible implementations.
type UnimplementedCollectServer struct {
}

func (UnimplementedCollectServer) Collect(context.Context, *CollectRequest) (*CollectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Collect not implemented")
}
func (UnimplementedCollectServer) UnCollect(context.Context, *UnCollectRequest) (*UnCollectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnCollect not implemented")
}
func (UnimplementedCollectServer) CollectList(context.Context, *CollectListRequest) (*CollectListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectList not implemented")
}
func (UnimplementedCollectServer) mustEmbedUnimplementedCollectServer() {}

// UnsafeCollectServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CollectServer will
// result in compilation errors.
type UnsafeCollectServer interface {
	mustEmbedUnimplementedCollectServer()
}

func RegisterCollectServer(s grpc.ServiceRegistrar, srv CollectServer) {
	s.RegisterService(&Collect_ServiceDesc, srv)
}

func _Collect_Collect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CollectServer).Collect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Collect/Collect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CollectServer).Collect(ctx, req.(*CollectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Collect_UnCollect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnCollectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CollectServer).UnCollect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Collect/UnCollect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CollectServer).UnCollect(ctx, req.(*UnCollectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Collect_CollectList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CollectServer).CollectList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Collect/CollectList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CollectServer).CollectList(ctx, req.(*CollectListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Collect_ServiceDesc is the grpc.ServiceDesc for Collect service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Collect_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.Collect",
	HandlerType: (*CollectServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Collect",
			Handler:    _Collect_Collect_Handler,
		},
		{
			MethodName: "UnCollect",
			Handler:    _Collect_UnCollect_Handler,
		},
		{
			MethodName: "CollectList",
			Handler:    _Collect_CollectList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "collect.proto",
}
