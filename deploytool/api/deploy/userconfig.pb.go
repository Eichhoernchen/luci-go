// Code generated by protoc-gen-go.
// source: github.com/luci/luci-go/deploytool/api/deploy/userconfig.proto
// DO NOT EDIT!

package deploy

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// *
// User configuration file.
//
// This is where user preferences will be written. Currently, no user
// preferences are defined.
type UserConfig struct {
	// *
	// Defines a local path override to a repository URL.
	//
	// When checkout is run with the "--local" flag, repositories whose URLs match
	// the key, the value URL will be used instead.
	//
	// The key is the repository URL to override, and the value is the override
	// URL, typically "file:///...".
	SourceOverride map[string]*Source `protobuf:"bytes,1,rep,name=source_override,json=sourceOverride" json:"source_override,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *UserConfig) Reset()                    { *m = UserConfig{} }
func (m *UserConfig) String() string            { return proto.CompactTextString(m) }
func (*UserConfig) ProtoMessage()               {}
func (*UserConfig) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{0} }

func (m *UserConfig) GetSourceOverride() map[string]*Source {
	if m != nil {
		return m.SourceOverride
	}
	return nil
}

func init() {
	proto.RegisterType((*UserConfig)(nil), "deploy.UserConfig")
}

func init() {
	proto.RegisterFile("github.com/luci/luci-go/deploytool/api/deploy/userconfig.proto", fileDescriptor4)
}

var fileDescriptor4 = []byte{
	// 207 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xb2, 0x4b, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0xcf, 0x29, 0x4d, 0xce, 0x04, 0x13, 0xba, 0xe9, 0xf9,
	0xfa, 0x29, 0xa9, 0x05, 0x39, 0xf9, 0x95, 0x25, 0xf9, 0xf9, 0x39, 0xfa, 0x89, 0x05, 0x99, 0x50,
	0xae, 0x7e, 0x69, 0x71, 0x6a, 0x51, 0x72, 0x7e, 0x5e, 0x5a, 0x66, 0xba, 0x5e, 0x41, 0x51, 0x7e,
	0x49, 0xbe, 0x10, 0x1b, 0x44, 0x42, 0xca, 0x8a, 0x34, 0x73, 0x90, 0xcd, 0x50, 0xda, 0xc0, 0xc8,
	0xc5, 0x15, 0x5a, 0x9c, 0x5a, 0xe4, 0x0c, 0x16, 0x14, 0xf2, 0xe7, 0xe2, 0x2f, 0xce, 0x2f, 0x2d,
	0x4a, 0x4e, 0x8d, 0xcf, 0x2f, 0x4b, 0x2d, 0x2a, 0xca, 0x4c, 0x49, 0x95, 0x60, 0x54, 0x60, 0xd6,
	0xe0, 0x36, 0x52, 0xd3, 0x83, 0xe8, 0xd6, 0x43, 0x28, 0xd6, 0x0b, 0x06, 0xab, 0xf4, 0x87, 0x2a,
	0x74, 0xcd, 0x2b, 0x29, 0xaa, 0x0c, 0xe2, 0x2b, 0x46, 0x11, 0x94, 0x0a, 0xe4, 0x12, 0xc6, 0xa2,
	0x4c, 0x48, 0x80, 0x8b, 0x39, 0x3b, 0xb5, 0x52, 0x82, 0x51, 0x81, 0x51, 0x83, 0x33, 0x08, 0xc4,
	0x14, 0x52, 0xe1, 0x62, 0x2d, 0x4b, 0xcc, 0x29, 0x4d, 0x95, 0x60, 0x52, 0x60, 0xd4, 0xe0, 0x36,
	0xe2, 0x83, 0xd9, 0x07, 0xd1, 0x1d, 0x04, 0x91, 0xb4, 0x62, 0xb2, 0x60, 0x4c, 0x62, 0x03, 0xbb,
	0xdc, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x2b, 0x77, 0x39, 0x22, 0x3f, 0x01, 0x00, 0x00,
}
