// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.4
// source: services/termserv/termserv.proto

package termserv

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

// TermservClient is the client API for Termserv service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TermservClient interface {
	LogoffUser(ctx context.Context, in *LogoffUserRequest, opts ...grpc.CallOption) (*LogoffUserResponse, error)
}

type termservClient struct {
	cc grpc.ClientConnInterface
}

func NewTermservClient(cc grpc.ClientConnInterface) TermservClient {
	return &termservClient{cc}
}

func (c *termservClient) LogoffUser(ctx context.Context, in *LogoffUserRequest, opts ...grpc.CallOption) (*LogoffUserResponse, error) {
	out := new(LogoffUserResponse)
	err := c.cc.Invoke(ctx, "/services.termserv.Termserv/LogoffUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TermservServer is the server API for Termserv service.
// All implementations must embed UnimplementedTermservServer
// for forward compatibility
type TermservServer interface {
	LogoffUser(context.Context, *LogoffUserRequest) (*LogoffUserResponse, error)
	mustEmbedUnimplementedTermservServer()
}

// UnimplementedTermservServer must be embedded to have forward compatible implementations.
type UnimplementedTermservServer struct {
}

func (UnimplementedTermservServer) LogoffUser(context.Context, *LogoffUserRequest) (*LogoffUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LogoffUser not implemented")
}
func (UnimplementedTermservServer) mustEmbedUnimplementedTermservServer() {}

// UnsafeTermservServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TermservServer will
// result in compilation errors.
type UnsafeTermservServer interface {
	mustEmbedUnimplementedTermservServer()
}

func RegisterTermservServer(s grpc.ServiceRegistrar, srv TermservServer) {
	s.RegisterService(&Termserv_ServiceDesc, srv)
}

func _Termserv_LogoffUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoffUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TermservServer).LogoffUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.termserv.Termserv/LogoffUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TermservServer).LogoffUser(ctx, req.(*LogoffUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Termserv_ServiceDesc is the grpc.ServiceDesc for Termserv service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Termserv_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "services.termserv.Termserv",
	HandlerType: (*TermservServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LogoffUser",
			Handler:    _Termserv_LogoffUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/termserv/termserv.proto",
}