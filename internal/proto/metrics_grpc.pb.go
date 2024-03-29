// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: proto/metrics.proto

package proto

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

const (
	Metrics_UpdateMetric_FullMethodName  = "/proto.Metrics/UpdateMetric"
	Metrics_UpdatesMetric_FullMethodName = "/proto.Metrics/UpdatesMetric"
	Metrics_ValueMetric_FullMethodName   = "/proto.Metrics/ValueMetric"
	Metrics_ListMetrics_FullMethodName   = "/proto.Metrics/ListMetrics"
)

// MetricsClient is the client API for Metrics service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricsClient interface {
	UpdateMetric(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error)
	UpdatesMetric(ctx context.Context, in *UpdatesRequest, opts ...grpc.CallOption) (*UpdatesResponse, error)
	ValueMetric(ctx context.Context, in *ValueRequest, opts ...grpc.CallOption) (*ValueResponse, error)
	ListMetrics(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
}

type metricsClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricsClient(cc grpc.ClientConnInterface) MetricsClient {
	return &metricsClient{cc}
}

func (c *metricsClient) UpdateMetric(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error) {
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, Metrics_UpdateMetric_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) UpdatesMetric(ctx context.Context, in *UpdatesRequest, opts ...grpc.CallOption) (*UpdatesResponse, error) {
	out := new(UpdatesResponse)
	err := c.cc.Invoke(ctx, Metrics_UpdatesMetric_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) ValueMetric(ctx context.Context, in *ValueRequest, opts ...grpc.CallOption) (*ValueResponse, error) {
	out := new(ValueResponse)
	err := c.cc.Invoke(ctx, Metrics_ValueMetric_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) ListMetrics(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, Metrics_ListMetrics_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricsServer is the server API for Metrics service.
// All implementations must embed UnimplementedMetricsServer
// for forward compatibility
type MetricsServer interface {
	UpdateMetric(context.Context, *UpdateRequest) (*UpdateResponse, error)
	UpdatesMetric(context.Context, *UpdatesRequest) (*UpdatesResponse, error)
	ValueMetric(context.Context, *ValueRequest) (*ValueResponse, error)
	ListMetrics(context.Context, *ListRequest) (*ListResponse, error)
	mustEmbedUnimplementedMetricsServer()
}

// UnimplementedMetricsServer must be embedded to have forward compatible implementations.
type UnimplementedMetricsServer struct {
}

func (UnimplementedMetricsServer) UpdateMetric(context.Context, *UpdateRequest) (*UpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMetric not implemented")
}
func (UnimplementedMetricsServer) UpdatesMetric(context.Context, *UpdatesRequest) (*UpdatesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatesMetric not implemented")
}
func (UnimplementedMetricsServer) ValueMetric(context.Context, *ValueRequest) (*ValueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValueMetric not implemented")
}
func (UnimplementedMetricsServer) ListMetrics(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMetrics not implemented")
}
func (UnimplementedMetricsServer) mustEmbedUnimplementedMetricsServer() {}

// UnsafeMetricsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsServer will
// result in compilation errors.
type UnsafeMetricsServer interface {
	mustEmbedUnimplementedMetricsServer()
}

func RegisterMetricsServer(s grpc.ServiceRegistrar, srv MetricsServer) {
	s.RegisterService(&Metrics_ServiceDesc, srv)
}

func _Metrics_UpdateMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).UpdateMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metrics_UpdateMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).UpdateMetric(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_UpdatesMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).UpdatesMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metrics_UpdatesMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).UpdatesMetric(ctx, req.(*UpdatesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_ValueMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).ValueMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metrics_ValueMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).ValueMetric(ctx, req.(*ValueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_ListMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).ListMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metrics_ListMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).ListMetrics(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Metrics_ServiceDesc is the grpc.ServiceDesc for Metrics service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Metrics_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Metrics",
	HandlerType: (*MetricsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateMetric",
			Handler:    _Metrics_UpdateMetric_Handler,
		},
		{
			MethodName: "UpdatesMetric",
			Handler:    _Metrics_UpdatesMetric_Handler,
		},
		{
			MethodName: "ValueMetric",
			Handler:    _Metrics_ValueMetric_Handler,
		},
		{
			MethodName: "ListMetrics",
			Handler:    _Metrics_ListMetrics_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/metrics.proto",
}
