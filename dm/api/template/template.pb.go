// Code generated by protoc-gen-go.
// source: github.com/luci/luci-go/dm/api/template/template.proto
// DO NOT EDIT!

/*
Package dmTemplate is a generated protocol buffer package.

It is generated from these files:
	github.com/luci/luci-go/dm/api/template/template.proto

It has these top-level messages:
	File
*/
package dmTemplate

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import templateproto "github.com/luci/luci-go/common/data/text/templateproto"
import dm1 "github.com/luci/luci-go/dm/api/service/v1"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// File represents a file full of DM template definitions.
type File struct {
	Template map[string]*File_Template `protobuf:"bytes,1,rep,name=template" json:"template,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *File) Reset()                    { *m = File{} }
func (m *File) String() string            { return proto.CompactTextString(m) }
func (*File) ProtoMessage()               {}
func (*File) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *File) GetTemplate() map[string]*File_Template {
	if m != nil {
		return m.Template
	}
	return nil
}

// Template defines a single template.
type File_Template struct {
	DistributorConfigName string                       `protobuf:"bytes,1,opt,name=distributor_config_name,json=distributorConfigName" json:"distributor_config_name,omitempty"`
	Parameters            *templateproto.File_Template `protobuf:"bytes,2,opt,name=parameters" json:"parameters,omitempty"`
	DistributorParameters *templateproto.File_Template `protobuf:"bytes,3,opt,name=distributor_parameters,json=distributorParameters" json:"distributor_parameters,omitempty"`
	Meta                  *dm1.Quest_Desc_Meta         `protobuf:"bytes,4,opt,name=meta" json:"meta,omitempty"`
}

func (m *File_Template) Reset()                    { *m = File_Template{} }
func (m *File_Template) String() string            { return proto.CompactTextString(m) }
func (*File_Template) ProtoMessage()               {}
func (*File_Template) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

func (m *File_Template) GetDistributorConfigName() string {
	if m != nil {
		return m.DistributorConfigName
	}
	return ""
}

func (m *File_Template) GetParameters() *templateproto.File_Template {
	if m != nil {
		return m.Parameters
	}
	return nil
}

func (m *File_Template) GetDistributorParameters() *templateproto.File_Template {
	if m != nil {
		return m.DistributorParameters
	}
	return nil
}

func (m *File_Template) GetMeta() *dm1.Quest_Desc_Meta {
	if m != nil {
		return m.Meta
	}
	return nil
}

func init() {
	proto.RegisterType((*File)(nil), "dmTemplate.File")
	proto.RegisterType((*File_Template)(nil), "dmTemplate.File.Template")
}

func init() {
	proto.RegisterFile("github.com/luci/luci-go/dm/api/template/template.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 327 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x90, 0x5f, 0x4b, 0xc3, 0x30,
	0x14, 0xc5, 0xe9, 0x3a, 0x65, 0xde, 0x21, 0x48, 0x44, 0x9d, 0x45, 0x64, 0xf8, 0xe2, 0x5e, 0x4c,
	0x70, 0xc2, 0x90, 0xe1, 0x9b, 0xce, 0x37, 0x45, 0xab, 0xf8, 0x3a, 0xb2, 0xf6, 0xda, 0x05, 0x9b,
	0xa6, 0xa4, 0xb7, 0xc3, 0x7d, 0x16, 0xbf, 0xab, 0x48, 0xbb, 0x7f, 0x9d, 0x30, 0xf4, 0xa5, 0xdc,
	0xe6, 0x9c, 0xfc, 0x72, 0xce, 0x85, 0x5e, 0xa4, 0x68, 0x9c, 0x8f, 0x78, 0x60, 0xb4, 0x88, 0xf3,
	0x40, 0x95, 0x9f, 0x8b, 0xc8, 0x88, 0x50, 0x0b, 0x99, 0x2a, 0x41, 0xa8, 0xd3, 0x58, 0x12, 0x2e,
	0x07, 0x9e, 0x5a, 0x43, 0x86, 0x41, 0xa8, 0x5f, 0xe7, 0x27, 0xde, 0x60, 0x13, 0x23, 0x30, 0x5a,
	0x9b, 0x44, 0x84, 0x92, 0xa4, 0x20, 0xfc, 0xa4, 0x25, 0xa4, 0x64, 0xfc, 0x42, 0x7a, 0xfd, 0x3f,
	0xa2, 0x64, 0x68, 0x27, 0x2a, 0x40, 0x31, 0xb9, 0x14, 0x91, 0x95, 0xe9, 0x78, 0x58, 0x70, 0x67,
	0x77, 0xcf, 0xbe, 0x5c, 0xa8, 0xdf, 0xab, 0x18, 0x59, 0x1f, 0x1a, 0x0b, 0x6c, 0xcb, 0x69, 0xbb,
	0x9d, 0x66, 0xf7, 0x94, 0xaf, 0xa2, 0xf2, 0xc2, 0xc3, 0x17, 0x7f, 0x83, 0x84, 0xec, 0xd4, 0x5f,
	0xfa, 0xbd, 0x6f, 0x07, 0x1a, 0x0b, 0x8d, 0xf5, 0xe0, 0x28, 0x54, 0x19, 0x59, 0x35, 0xca, 0xc9,
	0xd8, 0x61, 0x60, 0x92, 0x77, 0x15, 0x0d, 0x13, 0xa9, 0x0b, 0xae, 0xd3, 0xd9, 0xf1, 0x0f, 0x2a,
	0xf2, 0x6d, 0xa9, 0x3e, 0x4a, 0x8d, 0xec, 0x06, 0x20, 0x95, 0x56, 0x6a, 0x24, 0xb4, 0x59, 0xab,
	0xd6, 0x76, 0x3a, 0xcd, 0xee, 0x09, 0x5f, 0x2b, 0xbe, 0x9e, 0xc2, 0xaf, 0xf8, 0xd9, 0x0b, 0x1c,
	0x56, 0x5f, 0xad, 0x90, 0xdc, 0x7f, 0x90, 0xaa, 0x91, 0x9e, 0x56, 0xd0, 0x73, 0xa8, 0x6b, 0x24,
	0xd9, 0xaa, 0x97, 0x88, 0x7d, 0x1e, 0x6a, 0xfe, 0x9c, 0x63, 0x46, 0xfc, 0x0e, 0xb3, 0x80, 0x3f,
	0x20, 0x49, 0xbf, 0x34, 0x78, 0x6f, 0xb0, 0xbb, 0xb6, 0x1b, 0xb6, 0x07, 0xee, 0x07, 0x4e, 0xe7,
	0x85, 0x8b, 0x91, 0x09, 0xd8, 0x9a, 0xc8, 0x38, 0xc7, 0x79, 0xb3, 0xe3, 0x8d, 0xcb, 0xf5, 0x67,
	0xbe, 0x7e, 0xed, 0xda, 0x19, 0x6d, 0x97, 0x61, 0xaf, 0x7e, 0x02, 0x00, 0x00, 0xff, 0xff, 0x3d,
	0x8e, 0xec, 0x55, 0x6d, 0x02, 0x00, 0x00,
}
