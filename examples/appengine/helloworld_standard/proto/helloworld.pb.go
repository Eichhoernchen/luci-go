// Code generated by protoc-gen-go.
// source: github.com/luci/luci-go/examples/appengine/helloworld_standard/proto/helloworld.proto
// DO NOT EDIT!

/*
Package helloworld is a generated protocol buffer package.

It is generated from these files:
	github.com/luci/luci-go/examples/appengine/helloworld_standard/proto/helloworld.proto

It has these top-level messages:
	HelloRequest
	HelloReply
*/
package helloworld

import prpc "github.com/luci/luci-go/grpc/prpc"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The request message containing the user's name.
type HelloRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *HelloRequest) Reset()                    { *m = HelloRequest{} }
func (m *HelloRequest) String() string            { return proto.CompactTextString(m) }
func (*HelloRequest) ProtoMessage()               {}
func (*HelloRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *HelloRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// The response message containing the greetings
type HelloReply struct {
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *HelloReply) Reset()                    { *m = HelloReply{} }
func (m *HelloReply) String() string            { return proto.CompactTextString(m) }
func (*HelloReply) ProtoMessage()               {}
func (*HelloReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *HelloReply) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*HelloRequest)(nil), "helloworld.HelloRequest")
	proto.RegisterType((*HelloReply)(nil), "helloworld.HelloReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Greeter service

type GreeterClient interface {
	// Sends a greeting
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}
type greeterPRPCClient struct {
	client *prpc.Client
}

func NewGreeterPRPCClient(client *prpc.Client) GreeterClient {
	return &greeterPRPCClient{client}
}

func (c *greeterPRPCClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.client.Call(ctx, "helloworld.Greeter", "SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type greeterClient struct {
	cc *grpc.ClientConn
}

func NewGreeterClient(cc *grpc.ClientConn) GreeterClient {
	return &greeterClient{cc}
}

func (c *greeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := grpc.Invoke(ctx, "/helloworld.Greeter/SayHello", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Greeter service

type GreeterServer interface {
	// Sends a greeting
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
}

func RegisterGreeterServer(s prpc.Registrar, srv GreeterServer) {
	s.RegisterService(&_Greeter_serviceDesc, srv)
}

func _Greeter_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/helloworld.Greeter/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Greeter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "helloworld.Greeter",
	HandlerType: (*GreeterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Greeter_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/luci/luci-go/examples/appengine/helloworld_standard/proto/helloworld.proto",
}

func init() {
	proto.RegisterFile("github.com/luci/luci-go/examples/appengine/helloworld_standard/proto/helloworld.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 190 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x8e, 0x3d, 0xcb, 0xc2, 0x30,
	0x10, 0xc7, 0x9f, 0xc2, 0x83, 0xd5, 0xc3, 0x29, 0x83, 0x14, 0x27, 0xc9, 0x20, 0x2e, 0x36, 0xa0,
	0xbb, 0xab, 0xba, 0x56, 0x9c, 0x25, 0x6d, 0x8f, 0xb4, 0x90, 0x37, 0x93, 0x14, 0xed, 0xb7, 0x17,
	0x83, 0xc5, 0x0e, 0x2e, 0xc7, 0xfd, 0x5f, 0x8e, 0xfb, 0xc1, 0x55, 0xb4, 0xa1, 0xe9, 0xca, 0xbc,
	0x32, 0x8a, 0xc9, 0xae, 0x6a, 0xe3, 0xd8, 0x0a, 0xc3, 0xf0, 0xc9, 0x95, 0x95, 0xe8, 0x19, 0xb7,
	0x16, 0xb5, 0x68, 0x35, 0xb2, 0x06, 0xa5, 0x34, 0x0f, 0xe3, 0x64, 0x7d, 0xf3, 0x81, 0xeb, 0x9a,
	0xbb, 0x9a, 0x59, 0x67, 0x82, 0x19, 0x25, 0x79, 0x34, 0x08, 0x7c, 0x1d, 0x4a, 0x61, 0x7e, 0x7a,
	0xab, 0x02, 0xef, 0x1d, 0xfa, 0x40, 0x08, 0xfc, 0x6b, 0xae, 0x30, 0x4b, 0x56, 0xc9, 0x66, 0x56,
	0xc4, 0x9d, 0xae, 0x01, 0x3e, 0x1d, 0x2b, 0x7b, 0x92, 0x41, 0xaa, 0xd0, 0x7b, 0x2e, 0x86, 0xd2,
	0x20, 0x77, 0x67, 0x48, 0x8f, 0x0e, 0x31, 0xa0, 0x23, 0x07, 0x98, 0x5e, 0x78, 0x1f, 0xaf, 0x48,
	0x96, 0x8f, 0x08, 0xc6, 0xcf, 0x96, 0x8b, 0x1f, 0x89, 0x95, 0x3d, 0xfd, 0x2b, 0x27, 0x91, 0x74,
	0xff, 0x0a, 0x00, 0x00, 0xff, 0xff, 0x2b, 0xad, 0x3e, 0x4f, 0x02, 0x01, 0x00, 0x00,
}
