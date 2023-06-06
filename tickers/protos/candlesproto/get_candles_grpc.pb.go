// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.1
// source: get_candles.proto

package candlesproto

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

// StockServiceClient is the client API for StockService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StockServiceClient interface {
	GetPrices(ctx context.Context, in *TickerId, opts ...grpc.CallOption) (*CandlesResponse, error)
}

type stockServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStockServiceClient(cc grpc.ClientConnInterface) StockServiceClient {
	return &stockServiceClient{cc}
}

func (c *stockServiceClient) GetPrices(ctx context.Context, in *TickerId, opts ...grpc.CallOption) (*CandlesResponse, error) {
	out := new(CandlesResponse)
	err := c.cc.Invoke(ctx, "/candlesproto.StockService/GetPrices", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StockServiceServer is the server API for StockService service.
// All implementations must embed UnimplementedStockServiceServer
// for forward compatibility
type StockServiceServer interface {
	GetPrices(context.Context, *TickerId) (*CandlesResponse, error)
	mustEmbedUnimplementedStockServiceServer()
}

// UnimplementedStockServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStockServiceServer struct {
}

func (UnimplementedStockServiceServer) GetPrices(context.Context, *TickerId) (*CandlesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPrices not implemented")
}
func (UnimplementedStockServiceServer) mustEmbedUnimplementedStockServiceServer() {}

// UnsafeStockServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StockServiceServer will
// result in compilation errors.
type UnsafeStockServiceServer interface {
	mustEmbedUnimplementedStockServiceServer()
}

func RegisterStockServiceServer(s grpc.ServiceRegistrar, srv StockServiceServer) {
	s.RegisterService(&StockService_ServiceDesc, srv)
}

func _StockService_GetPrices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TickerId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StockServiceServer).GetPrices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/candlesproto.StockService/GetPrices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StockServiceServer).GetPrices(ctx, req.(*TickerId))
	}
	return interceptor(ctx, in, info, handler)
}

// StockService_ServiceDesc is the grpc.ServiceDesc for StockService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StockService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "candlesproto.StockService",
	HandlerType: (*StockServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPrices",
			Handler:    _StockService_GetPrices_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "get_candles.proto",
}
