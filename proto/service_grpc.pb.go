// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: proto/service.proto

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

// DiscoveryClient is the client API for Discovery service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DiscoveryClient interface {
	GetInfo(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*InfoResponse, error)
	Connect(ctx context.Context, opts ...grpc.CallOption) (Discovery_ConnectClient, error)
	ConnectBack(ctx context.Context, opts ...grpc.CallOption) (Discovery_ConnectBackClient, error)
	// rpc Sync(stream SyncMessage) returns (Close);   // alternative for join
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PongResponse, error)
	Call(ctx context.Context, in *CallRequest, opts ...grpc.CallOption) (*CallResponse, error)
}

type discoveryClient struct {
	cc grpc.ClientConnInterface
}

func NewDiscoveryClient(cc grpc.ClientConnInterface) DiscoveryClient {
	return &discoveryClient{cc}
}

func (c *discoveryClient) GetInfo(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*InfoResponse, error) {
	out := new(InfoResponse)
	err := c.cc.Invoke(ctx, "/proto.Discovery/GetInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discoveryClient) Connect(ctx context.Context, opts ...grpc.CallOption) (Discovery_ConnectClient, error) {
	stream, err := c.cc.NewStream(ctx, &Discovery_ServiceDesc.Streams[0], "/proto.Discovery/Connect", opts...)
	if err != nil {
		return nil, err
	}
	x := &discoveryConnectClient{stream}
	return x, nil
}

type Discovery_ConnectClient interface {
	Send(*ConnectMessage) error
	CloseAndRecv() (*Close, error)
	grpc.ClientStream
}

type discoveryConnectClient struct {
	grpc.ClientStream
}

func (x *discoveryConnectClient) Send(m *ConnectMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *discoveryConnectClient) CloseAndRecv() (*Close, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Close)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *discoveryClient) ConnectBack(ctx context.Context, opts ...grpc.CallOption) (Discovery_ConnectBackClient, error) {
	stream, err := c.cc.NewStream(ctx, &Discovery_ServiceDesc.Streams[1], "/proto.Discovery/ConnectBack", opts...)
	if err != nil {
		return nil, err
	}
	x := &discoveryConnectBackClient{stream}
	return x, nil
}

type Discovery_ConnectBackClient interface {
	Send(*ConnectBackMessage) error
	CloseAndRecv() (*Close, error)
	grpc.ClientStream
}

type discoveryConnectBackClient struct {
	grpc.ClientStream
}

func (x *discoveryConnectBackClient) Send(m *ConnectBackMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *discoveryConnectBackClient) CloseAndRecv() (*Close, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Close)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *discoveryClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PongResponse, error) {
	out := new(PongResponse)
	err := c.cc.Invoke(ctx, "/proto.Discovery/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discoveryClient) Call(ctx context.Context, in *CallRequest, opts ...grpc.CallOption) (*CallResponse, error) {
	out := new(CallResponse)
	err := c.cc.Invoke(ctx, "/proto.Discovery/Call", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DiscoveryServer is the server API for Discovery service.
// All implementations should embed UnimplementedDiscoveryServer
// for forward compatibility
type DiscoveryServer interface {
	GetInfo(context.Context, *EmptyRequest) (*InfoResponse, error)
	Connect(Discovery_ConnectServer) error
	ConnectBack(Discovery_ConnectBackServer) error
	// rpc Sync(stream SyncMessage) returns (Close);   // alternative for join
	Ping(context.Context, *PingRequest) (*PongResponse, error)
	Call(context.Context, *CallRequest) (*CallResponse, error)
}

// UnimplementedDiscoveryServer should be embedded to have forward compatible implementations.
type UnimplementedDiscoveryServer struct {
}

func (UnimplementedDiscoveryServer) GetInfo(context.Context, *EmptyRequest) (*InfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInfo not implemented")
}
func (UnimplementedDiscoveryServer) Connect(Discovery_ConnectServer) error {
	return status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedDiscoveryServer) ConnectBack(Discovery_ConnectBackServer) error {
	return status.Errorf(codes.Unimplemented, "method ConnectBack not implemented")
}
func (UnimplementedDiscoveryServer) Ping(context.Context, *PingRequest) (*PongResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedDiscoveryServer) Call(context.Context, *CallRequest) (*CallResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Call not implemented")
}

// UnsafeDiscoveryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DiscoveryServer will
// result in compilation errors.
type UnsafeDiscoveryServer interface {
	mustEmbedUnimplementedDiscoveryServer()
}

func RegisterDiscoveryServer(s grpc.ServiceRegistrar, srv DiscoveryServer) {
	s.RegisterService(&Discovery_ServiceDesc, srv)
}

func _Discovery_GetInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscoveryServer).GetInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Discovery/GetInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscoveryServer).GetInfo(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Discovery_Connect_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DiscoveryServer).Connect(&discoveryConnectServer{stream})
}

type Discovery_ConnectServer interface {
	SendAndClose(*Close) error
	Recv() (*ConnectMessage, error)
	grpc.ServerStream
}

type discoveryConnectServer struct {
	grpc.ServerStream
}

func (x *discoveryConnectServer) SendAndClose(m *Close) error {
	return x.ServerStream.SendMsg(m)
}

func (x *discoveryConnectServer) Recv() (*ConnectMessage, error) {
	m := new(ConnectMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Discovery_ConnectBack_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DiscoveryServer).ConnectBack(&discoveryConnectBackServer{stream})
}

type Discovery_ConnectBackServer interface {
	SendAndClose(*Close) error
	Recv() (*ConnectBackMessage, error)
	grpc.ServerStream
}

type discoveryConnectBackServer struct {
	grpc.ServerStream
}

func (x *discoveryConnectBackServer) SendAndClose(m *Close) error {
	return x.ServerStream.SendMsg(m)
}

func (x *discoveryConnectBackServer) Recv() (*ConnectBackMessage, error) {
	m := new(ConnectBackMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Discovery_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscoveryServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Discovery/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscoveryServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Discovery_Call_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CallRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscoveryServer).Call(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Discovery/Call",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscoveryServer).Call(ctx, req.(*CallRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Discovery_ServiceDesc is the grpc.ServiceDesc for Discovery service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Discovery_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Discovery",
	HandlerType: (*DiscoveryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetInfo",
			Handler:    _Discovery_GetInfo_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Discovery_Ping_Handler,
		},
		{
			MethodName: "Call",
			Handler:    _Discovery_Call_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Connect",
			Handler:       _Discovery_Connect_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "ConnectBack",
			Handler:       _Discovery_ConnectBack_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "proto/service.proto",
}

// ProxyManagerClient is the client API for ProxyManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProxyManagerClient interface {
	SendPayload(ctx context.Context, in *RequestPayload, opts ...grpc.CallOption) (*ResponsePayload, error)
}

type proxyManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewProxyManagerClient(cc grpc.ClientConnInterface) ProxyManagerClient {
	return &proxyManagerClient{cc}
}

func (c *proxyManagerClient) SendPayload(ctx context.Context, in *RequestPayload, opts ...grpc.CallOption) (*ResponsePayload, error) {
	out := new(ResponsePayload)
	err := c.cc.Invoke(ctx, "/proto.ProxyManager/SendPayload", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProxyManagerServer is the server API for ProxyManager service.
// All implementations should embed UnimplementedProxyManagerServer
// for forward compatibility
type ProxyManagerServer interface {
	SendPayload(context.Context, *RequestPayload) (*ResponsePayload, error)
}

// UnimplementedProxyManagerServer should be embedded to have forward compatible implementations.
type UnimplementedProxyManagerServer struct {
}

func (UnimplementedProxyManagerServer) SendPayload(context.Context, *RequestPayload) (*ResponsePayload, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendPayload not implemented")
}

// UnsafeProxyManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProxyManagerServer will
// result in compilation errors.
type UnsafeProxyManagerServer interface {
	mustEmbedUnimplementedProxyManagerServer()
}

func RegisterProxyManagerServer(s grpc.ServiceRegistrar, srv ProxyManagerServer) {
	s.RegisterService(&ProxyManager_ServiceDesc, srv)
}

func _ProxyManager_SendPayload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestPayload)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyManagerServer).SendPayload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ProxyManager/SendPayload",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyManagerServer).SendPayload(ctx, req.(*RequestPayload))
	}
	return interceptor(ctx, in, info, handler)
}

// ProxyManager_ServiceDesc is the grpc.ServiceDesc for ProxyManager service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProxyManager_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ProxyManager",
	HandlerType: (*ProxyManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendPayload",
			Handler:    _ProxyManager_SendPayload_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/service.proto",
}
