// Code generated by protoc-gen-go.
// source: ensure_quests.proto
// DO NOT EDIT!

package dm

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// EnsureQuestsReq allows either a human or a client to ensure that various
// Quests exist in DM.
type EnsureQuestsReq struct {
	// optional: Only needed if this is being called from an execution.
	Auth *ExecutionAuth `protobuf:"bytes,1,opt,name=auth" json:"auth,omitempty"`
	// ToEnsure is a list of quest descriptors. DM will ensure that the
	// corresponding Quests exist. If they don't, they'll be created.
	ToEnsure []*QuestDescriptor `protobuf:"bytes,2,rep,name=to_ensure" json:"to_ensure,omitempty"`
}

func (m *EnsureQuestsReq) Reset()                    { *m = EnsureQuestsReq{} }
func (m *EnsureQuestsReq) String() string            { return proto.CompactTextString(m) }
func (*EnsureQuestsReq) ProtoMessage()               {}
func (*EnsureQuestsReq) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *EnsureQuestsReq) GetAuth() *ExecutionAuth {
	if m != nil {
		return m.Auth
	}
	return nil
}

func (m *EnsureQuestsReq) GetToEnsure() []*QuestDescriptor {
	if m != nil {
		return m.ToEnsure
	}
	return nil
}

type EnsureQuestsRsp struct {
	// QuestIds is a list that matches 1:1 with the QuestDescriptors listed in the
	// request. These IDs are the canonical ids of the Quests.
	QuestIds []*QuestID `protobuf:"bytes,1,rep,name=quest_ids" json:"quest_ids,omitempty"`
}

func (m *EnsureQuestsRsp) Reset()                    { *m = EnsureQuestsRsp{} }
func (m *EnsureQuestsRsp) String() string            { return proto.CompactTextString(m) }
func (*EnsureQuestsRsp) ProtoMessage()               {}
func (*EnsureQuestsRsp) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

func (m *EnsureQuestsRsp) GetQuestIds() []*QuestID {
	if m != nil {
		return m.QuestIds
	}
	return nil
}

func init() {
	proto.RegisterType((*EnsureQuestsReq)(nil), "dm.EnsureQuestsReq")
	proto.RegisterType((*EnsureQuestsRsp)(nil), "dm.EnsureQuestsRsp")
}

var fileDescriptor3 = []byte{
	// 163 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x4e, 0xcd, 0x2b, 0x2e,
	0x2d, 0x4a, 0x8d, 0x2f, 0x2c, 0x4d, 0x2d, 0x2e, 0x29, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17,
	0x62, 0x4a, 0xc9, 0x95, 0xe2, 0x2e, 0xa9, 0x2c, 0x48, 0x85, 0x0a, 0x28, 0x45, 0x71, 0xf1, 0xbb,
	0x82, 0xd5, 0x05, 0x82, 0x95, 0x05, 0xa5, 0x16, 0x0a, 0xc9, 0x73, 0xb1, 0x24, 0x96, 0x96, 0x64,
	0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x1b, 0x09, 0xea, 0xa5, 0xe4, 0xea, 0xb9, 0x56, 0xa4, 0x26,
	0x97, 0x96, 0x64, 0xe6, 0xe7, 0x39, 0x02, 0x25, 0x84, 0xd4, 0xb8, 0x38, 0x4b, 0xf2, 0xe3, 0x21,
	0xc6, 0x4b, 0x30, 0x29, 0x30, 0x03, 0x55, 0x09, 0x83, 0x54, 0x81, 0x8d, 0x70, 0x49, 0x2d, 0x4e,
	0x2e, 0xca, 0x2c, 0x28, 0xc9, 0x2f, 0x52, 0x32, 0x44, 0x33, 0xbb, 0xb8, 0x40, 0x48, 0x8e, 0x8b,
	0x13, 0xec, 0x9e, 0xf8, 0xcc, 0x94, 0x62, 0xa0, 0x05, 0x20, 0xad, 0xdc, 0x70, 0xad, 0x9e, 0x2e,
	0x49, 0x6c, 0x60, 0x57, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x45, 0xfc, 0x50, 0xe7, 0xbd,
	0x00, 0x00, 0x00,
}