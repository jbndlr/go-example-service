// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// ExampleClient is the client API for Example service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExampleClient interface {
	Info(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*InfoMessage, error)
}

type exampleClient struct {
	cc grpc.ClientConnInterface
}

func NewExampleClient(cc grpc.ClientConnInterface) ExampleClient {
	return &exampleClient{cc}
}

func (c *exampleClient) Info(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*InfoMessage, error) {
	out := new(InfoMessage)
	err := c.cc.Invoke(ctx, "/grpc.Example/Info", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExampleServer is the server API for Example service.
// All implementations must embed UnimplementedExampleServer
// for forward compatibility
type ExampleServer interface {
	Info(context.Context, *Empty) (*InfoMessage, error)
	mustEmbedUnimplementedExampleServer()
}

// UnimplementedExampleServer must be embedded to have forward compatible implementations.
type UnimplementedExampleServer struct {
}

func (UnimplementedExampleServer) Info(context.Context, *Empty) (*InfoMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Info not implemented")
}
func (UnimplementedExampleServer) mustEmbedUnimplementedExampleServer() {}

// UnsafeExampleServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExampleServer will
// result in compilation errors.
type UnsafeExampleServer interface {
	mustEmbedUnimplementedExampleServer()
}

func RegisterExampleServer(s grpc.ServiceRegistrar, srv ExampleServer) {
	s.RegisterService(&_Example_serviceDesc, srv)
}

func _Example_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExampleServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Example/Info",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExampleServer).Info(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Example_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Example",
	HandlerType: (*ExampleServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Info",
			Handler:    _Example_Info_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
