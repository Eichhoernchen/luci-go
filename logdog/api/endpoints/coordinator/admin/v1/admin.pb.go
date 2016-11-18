// Code generated by protoc-gen-go.
// source: github.com/luci/luci-go/logdog/api/endpoints/coordinator/admin/v1/admin.proto
// DO NOT EDIT!

/*
Package logdog is a generated protocol buffer package.

It is generated from these files:
	github.com/luci/luci-go/logdog/api/endpoints/coordinator/admin/v1/admin.proto

It has these top-level messages:
	SetConfigRequest
*/
package logdog

import prpc "github.com/luci/luci-go/grpc/prpc"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/luci/luci-go/common/proto/google"

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

// GlobalConfig is the LogDog Coordinator global configuration.
//
// This is intended to act as an entry point. The majority of the configuration
// will be stored in a "luci-config" service Config protobuf.
type SetConfigRequest struct {
	// ConfigServiceURL is the API URL of the base "luci-config" service. If
	// empty, the defualt service URL will be used.
	ConfigServiceUrl string `protobuf:"bytes,1,opt,name=config_service_url,json=configServiceUrl" json:"config_service_url,omitempty"`
	// ConfigSet is the name of the configuration set to load from.
	ConfigSet string `protobuf:"bytes,2,opt,name=config_set,json=configSet" json:"config_set,omitempty"`
	// ConfigPath is the path of the text-serialized configuration protobuf.
	ConfigPath string `protobuf:"bytes,3,opt,name=config_path,json=configPath" json:"config_path,omitempty"`
	// If not empty, is the service account JSON file data that will be used for
	// Storage access.
	//
	// TODO(dnj): Remove this option once Cloud BigTable has cross-project ACLs.
	StorageServiceAccountJson []byte `protobuf:"bytes,100,opt,name=storage_service_account_json,json=storageServiceAccountJson,proto3" json:"storage_service_account_json,omitempty"`
}

func (m *SetConfigRequest) Reset()                    { *m = SetConfigRequest{} }
func (m *SetConfigRequest) String() string            { return proto.CompactTextString(m) }
func (*SetConfigRequest) ProtoMessage()               {}
func (*SetConfigRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *SetConfigRequest) GetConfigServiceUrl() string {
	if m != nil {
		return m.ConfigServiceUrl
	}
	return ""
}

func (m *SetConfigRequest) GetConfigSet() string {
	if m != nil {
		return m.ConfigSet
	}
	return ""
}

func (m *SetConfigRequest) GetConfigPath() string {
	if m != nil {
		return m.ConfigPath
	}
	return ""
}

func (m *SetConfigRequest) GetStorageServiceAccountJson() []byte {
	if m != nil {
		return m.StorageServiceAccountJson
	}
	return nil
}

func init() {
	proto.RegisterType((*SetConfigRequest)(nil), "logdog.SetConfigRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Admin service

type AdminClient interface {
	// SetConfig loads the supplied configuration into a config.GlobalConfig
	// instance.
	SetConfig(ctx context.Context, in *SetConfigRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
}
type adminPRPCClient struct {
	client *prpc.Client
}

func NewAdminPRPCClient(client *prpc.Client) AdminClient {
	return &adminPRPCClient{client}
}

func (c *adminPRPCClient) SetConfig(ctx context.Context, in *SetConfigRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := c.client.Call(ctx, "logdog.Admin", "SetConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type adminClient struct {
	cc *grpc.ClientConn
}

func NewAdminClient(cc *grpc.ClientConn) AdminClient {
	return &adminClient{cc}
}

func (c *adminClient) SetConfig(ctx context.Context, in *SetConfigRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/logdog.Admin/SetConfig", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Admin service

type AdminServer interface {
	// SetConfig loads the supplied configuration into a config.GlobalConfig
	// instance.
	SetConfig(context.Context, *SetConfigRequest) (*google_protobuf.Empty, error)
}

func RegisterAdminServer(s prpc.Registrar, srv AdminServer) {
	s.RegisterService(&_Admin_serviceDesc, srv)
}

func _Admin_SetConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).SetConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/logdog.Admin/SetConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).SetConfig(ctx, req.(*SetConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Admin_serviceDesc = grpc.ServiceDesc{
	ServiceName: "logdog.Admin",
	HandlerType: (*AdminServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetConfig",
			Handler:    _Admin_SetConfig_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/luci/luci-go/logdog/api/endpoints/coordinator/admin/v1/admin.proto",
}

func init() {
	proto.RegisterFile("github.com/luci/luci-go/logdog/api/endpoints/coordinator/admin/v1/admin.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 277 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x5c, 0x90, 0xc1, 0x4b, 0xc3, 0x30,
	0x14, 0xc6, 0xa9, 0xe2, 0x60, 0xd1, 0xc3, 0xc8, 0x41, 0xea, 0x54, 0x1c, 0x9e, 0x76, 0xd0, 0x04,
	0xf5, 0x2c, 0x32, 0x44, 0x0f, 0x82, 0x20, 0x1d, 0x9e, 0x4b, 0x9a, 0x66, 0x69, 0xa4, 0xcd, 0xab,
	0xc9, 0xcb, 0xc0, 0x3f, 0xcf, 0xff, 0x4c, 0x96, 0x74, 0x3d, 0x78, 0x09, 0x79, 0xef, 0xfb, 0x3d,
	0x3e, 0xbe, 0x8f, 0xbc, 0x6b, 0x83, 0x4d, 0xa8, 0x98, 0x84, 0x8e, 0xb7, 0x41, 0x9a, 0xf8, 0xdc,
	0x6a, 0xe0, 0x2d, 0xe8, 0x1a, 0x34, 0x17, 0xbd, 0xe1, 0xca, 0xd6, 0x3d, 0x18, 0x8b, 0x9e, 0x4b,
	0x00, 0x57, 0x1b, 0x2b, 0x10, 0x1c, 0x17, 0x75, 0x67, 0x2c, 0xdf, 0xde, 0xa5, 0x0f, 0xeb, 0x1d,
	0x20, 0xd0, 0x49, 0x3a, 0x9b, 0x9f, 0x6b, 0x00, 0xdd, 0x2a, 0x1e, 0xb7, 0x55, 0xd8, 0x70, 0xd5,
	0xf5, 0xf8, 0x93, 0xa0, 0xeb, 0xdf, 0x8c, 0xcc, 0xd6, 0x0a, 0x9f, 0xc1, 0x6e, 0x8c, 0x2e, 0xd4,
	0x77, 0x50, 0x1e, 0xe9, 0x0d, 0xa1, 0x32, 0x2e, 0x4a, 0xaf, 0xdc, 0xd6, 0x48, 0x55, 0x06, 0xd7,
	0xe6, 0xd9, 0x22, 0x5b, 0x4e, 0x8b, 0x59, 0x52, 0xd6, 0x49, 0xf8, 0x74, 0x2d, 0xbd, 0x24, 0x64,
	0xa4, 0x31, 0x3f, 0x88, 0xd4, 0x74, 0x4f, 0x21, 0xbd, 0x22, 0xc7, 0x83, 0xdc, 0x0b, 0x6c, 0xf2,
	0xc3, 0xa8, 0x0f, 0x17, 0x1f, 0x02, 0x1b, 0xfa, 0x44, 0x2e, 0x3c, 0x82, 0x13, 0x5a, 0x8d, 0x76,
	0x42, 0x4a, 0x08, 0x16, 0xcb, 0x2f, 0x0f, 0x36, 0xaf, 0x17, 0xd9, 0xf2, 0xa4, 0x38, 0x1b, 0x98,
	0xc1, 0x78, 0x95, 0x88, 0x37, 0x0f, 0xf6, 0xfe, 0x95, 0x1c, 0xad, 0x76, 0xb9, 0xe9, 0x23, 0x99,
	0x8e, 0x59, 0x68, 0xce, 0x52, 0x7e, 0xf6, 0x3f, 0xde, 0xfc, 0x94, 0xa5, 0x46, 0xd8, 0xbe, 0x11,
	0xf6, 0xb2, 0x6b, 0xa4, 0x9a, 0xc4, 0xf9, 0xe1, 0x2f, 0x00, 0x00, 0xff, 0xff, 0x39, 0x46, 0x12,
	0xba, 0x88, 0x01, 0x00, 0x00,
}
