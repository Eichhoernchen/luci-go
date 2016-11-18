// Code generated by protoc-gen-go.
// source: github.com/luci/luci-go/dm/api/service/v1/walk_graph.proto
// DO NOT EDIT!

package dm

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/luci/luci-go/common/proto/google"
import _ "github.com/luci/luci-go/common/proto/google"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Direction indicates that direction of dependencies that the request should
// walk.
type WalkGraphReq_Mode_Direction int32

const (
	WalkGraphReq_Mode_FORWARDS  WalkGraphReq_Mode_Direction = 0
	WalkGraphReq_Mode_BACKWARDS WalkGraphReq_Mode_Direction = 1
	WalkGraphReq_Mode_BOTH      WalkGraphReq_Mode_Direction = 2
)

var WalkGraphReq_Mode_Direction_name = map[int32]string{
	0: "FORWARDS",
	1: "BACKWARDS",
	2: "BOTH",
}
var WalkGraphReq_Mode_Direction_value = map[string]int32{
	"FORWARDS":  0,
	"BACKWARDS": 1,
	"BOTH":      2,
}

func (x WalkGraphReq_Mode_Direction) String() string {
	return proto.EnumName(WalkGraphReq_Mode_Direction_name, int32(x))
}
func (WalkGraphReq_Mode_Direction) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor7, []int{0, 0, 0}
}

// WalkGraphReq allows you to walk from one or more Quests through their
// Attempt's forward dependencies.
//
//
// The handler will evaluate all of the queries, executing them in parallel.
// For each attempt or quest produced by the query, it will queue a walk
// operation for that node, respecting the options set (max_depth, etc.).
type WalkGraphReq struct {
	// Optional. See Include.AttemptResult for restrictions.
	Auth *Execution_Auth `protobuf:"bytes,1,opt,name=auth" json:"auth,omitempty"`
	// Query specifies a list of queries to start the graph traversal on. The
	// traversal will occur as a union of the query results. Redundant
	// specification will not cause additional heavy work; every graph node will
	// be processed exactly once, regardless of how many times it appears in the
	// query results. However, redundancy in the queries will cause the server to
	// retrieve and discard more information.
	Query *GraphQuery         `protobuf:"bytes,2,opt,name=query" json:"query,omitempty"`
	Mode  *WalkGraphReq_Mode  `protobuf:"bytes,3,opt,name=mode" json:"mode,omitempty"`
	Limit *WalkGraphReq_Limit `protobuf:"bytes,4,opt,name=limit" json:"limit,omitempty"`
	// Include allows you to add additional information to the returned
	// GraphData which is typically medium-to-large sized.
	Include *WalkGraphReq_Include `protobuf:"bytes,5,opt,name=include" json:"include,omitempty"`
	Exclude *WalkGraphReq_Exclude `protobuf:"bytes,6,opt,name=exclude" json:"exclude,omitempty"`
}

func (m *WalkGraphReq) Reset()                    { *m = WalkGraphReq{} }
func (m *WalkGraphReq) String() string            { return proto.CompactTextString(m) }
func (*WalkGraphReq) ProtoMessage()               {}
func (*WalkGraphReq) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{0} }

func (m *WalkGraphReq) GetAuth() *Execution_Auth {
	if m != nil {
		return m.Auth
	}
	return nil
}

func (m *WalkGraphReq) GetQuery() *GraphQuery {
	if m != nil {
		return m.Query
	}
	return nil
}

func (m *WalkGraphReq) GetMode() *WalkGraphReq_Mode {
	if m != nil {
		return m.Mode
	}
	return nil
}

func (m *WalkGraphReq) GetLimit() *WalkGraphReq_Limit {
	if m != nil {
		return m.Limit
	}
	return nil
}

func (m *WalkGraphReq) GetInclude() *WalkGraphReq_Include {
	if m != nil {
		return m.Include
	}
	return nil
}

func (m *WalkGraphReq) GetExclude() *WalkGraphReq_Exclude {
	if m != nil {
		return m.Exclude
	}
	return nil
}

type WalkGraphReq_Mode struct {
	// DFS sets whether this is a Depth-first (ish) or a Breadth-first (ish) load.
	// Since the load operation is multi-threaded, the search order is best
	// effort, but will actually be some hybrid between DFS and BFS. This setting
	// controls the bias direction of the hybrid loading algorithm.
	Dfs       bool                        `protobuf:"varint,1,opt,name=dfs" json:"dfs,omitempty"`
	Direction WalkGraphReq_Mode_Direction `protobuf:"varint,2,opt,name=direction,enum=dm.WalkGraphReq_Mode_Direction" json:"direction,omitempty"`
}

func (m *WalkGraphReq_Mode) Reset()                    { *m = WalkGraphReq_Mode{} }
func (m *WalkGraphReq_Mode) String() string            { return proto.CompactTextString(m) }
func (*WalkGraphReq_Mode) ProtoMessage()               {}
func (*WalkGraphReq_Mode) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{0, 0} }

func (m *WalkGraphReq_Mode) GetDfs() bool {
	if m != nil {
		return m.Dfs
	}
	return false
}

func (m *WalkGraphReq_Mode) GetDirection() WalkGraphReq_Mode_Direction {
	if m != nil {
		return m.Direction
	}
	return WalkGraphReq_Mode_FORWARDS
}

type WalkGraphReq_Limit struct {
	// MaxDepth sets the number of attempts to traverse; 0 means 'immediate'
	// (no dependencies), -1 means 'no limit', and >0 is a limit.
	//
	// Any negative value besides -1 is an error.
	MaxDepth int64 `protobuf:"varint,1,opt,name=max_depth,json=maxDepth" json:"max_depth,omitempty"`
	// MaxTime sets the maximum amount of time that the query processor should
	// take. Application of this deadline is 'best effort', which means the query
	// may take a bit longer than this timeout and still attempt to return data.
	//
	// This is different than the grpc timeout header, which will set a hard
	// deadline for the request.
	MaxTime *google_protobuf.Duration `protobuf:"bytes,2,opt,name=max_time,json=maxTime" json:"max_time,omitempty"`
	// MaxDataSize sets the maximum amount of 'Data' (in bytes) that can be
	// returned, if include.quest_data, include.attempt_data, and/or
	// include.attempt_result are set. If this limit is hit, then the
	// appropriate 'partial' value will be set for that object, but the base
	// object would still be included in the result.
	//
	// If this limit is 0, a default limit of 16MB will be used. If this limit
	// exceeds 30MB, it will be reduced to 30MB.
	MaxDataSize uint32 `protobuf:"varint,3,opt,name=max_data_size,json=maxDataSize" json:"max_data_size,omitempty"`
}

func (m *WalkGraphReq_Limit) Reset()                    { *m = WalkGraphReq_Limit{} }
func (m *WalkGraphReq_Limit) String() string            { return proto.CompactTextString(m) }
func (*WalkGraphReq_Limit) ProtoMessage()               {}
func (*WalkGraphReq_Limit) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{0, 1} }

func (m *WalkGraphReq_Limit) GetMaxDepth() int64 {
	if m != nil {
		return m.MaxDepth
	}
	return 0
}

func (m *WalkGraphReq_Limit) GetMaxTime() *google_protobuf.Duration {
	if m != nil {
		return m.MaxTime
	}
	return nil
}

func (m *WalkGraphReq_Limit) GetMaxDataSize() uint32 {
	if m != nil {
		return m.MaxDataSize
	}
	return 0
}

type WalkGraphReq_Include struct {
	Quest     *WalkGraphReq_Include_Options `protobuf:"bytes,1,opt,name=quest" json:"quest,omitempty"`
	Attempt   *WalkGraphReq_Include_Options `protobuf:"bytes,2,opt,name=attempt" json:"attempt,omitempty"`
	Execution *WalkGraphReq_Include_Options `protobuf:"bytes,3,opt,name=execution" json:"execution,omitempty"`
	// Executions is the number of Executions to include per Attempt. If this
	// is 0, then the execution data will be omitted completely.
	//
	// Executions included are from high ids to low ids. So setting this to `1`
	// would return the LAST execution made for this Attempt.
	NumExecutions uint32 `protobuf:"varint,4,opt,name=num_executions,json=numExecutions" json:"num_executions,omitempty"`
	// FwdDeps instructs WalkGraph to include forward dependency information
	// from the result. This only changes the presence of information in the
	// result; if the query is walking forward attempt dependencies, that will
	// still occur even if this is false.
	FwdDeps bool `protobuf:"varint,5,opt,name=fwd_deps,json=fwdDeps" json:"fwd_deps,omitempty"`
	// BackDeps instructs WalkGraph to include the backwards dependency
	// information. This only changes the presence of information in the result;
	// if the query is walking backward attempt dependencies, that will still
	// occur even if this is false.
	BackDeps bool `protobuf:"varint,6,opt,name=back_deps,json=backDeps" json:"back_deps,omitempty"`
}

func (m *WalkGraphReq_Include) Reset()                    { *m = WalkGraphReq_Include{} }
func (m *WalkGraphReq_Include) String() string            { return proto.CompactTextString(m) }
func (*WalkGraphReq_Include) ProtoMessage()               {}
func (*WalkGraphReq_Include) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{0, 2} }

func (m *WalkGraphReq_Include) GetQuest() *WalkGraphReq_Include_Options {
	if m != nil {
		return m.Quest
	}
	return nil
}

func (m *WalkGraphReq_Include) GetAttempt() *WalkGraphReq_Include_Options {
	if m != nil {
		return m.Attempt
	}
	return nil
}

func (m *WalkGraphReq_Include) GetExecution() *WalkGraphReq_Include_Options {
	if m != nil {
		return m.Execution
	}
	return nil
}

func (m *WalkGraphReq_Include) GetNumExecutions() uint32 {
	if m != nil {
		return m.NumExecutions
	}
	return 0
}

func (m *WalkGraphReq_Include) GetFwdDeps() bool {
	if m != nil {
		return m.FwdDeps
	}
	return false
}

func (m *WalkGraphReq_Include) GetBackDeps() bool {
	if m != nil {
		return m.BackDeps
	}
	return false
}

type WalkGraphReq_Include_Options struct {
	// Fills the 'id' field.
	//
	// If this is false, it will be omitted.
	//
	// Note that there's enough information contextually to derive these ids
	// on the client side, though it can be handy to have the server produce
	// them for you.
	Ids bool `protobuf:"varint,1,opt,name=ids" json:"ids,omitempty"`
	// Instructs the request to include the Data field
	Data bool `protobuf:"varint,2,opt,name=data" json:"data,omitempty"`
	// Instructs finished objects to include the Result field.
	//
	// If the requestor is an execution, the query logic will only include the
	// result if the execution's Attempt depends on it, otherwise it will be
	// blank.
	//
	// If the request's cumulative result data would be more than
	// limit.max_data_size of data, the remaining results will have their
	// Partial.Result set to DATA_SIZE_LIMIT.
	//
	// Has no effect for Quests.
	Result bool `protobuf:"varint,3,opt,name=result" json:"result,omitempty"`
	// If set to true, objects with an abnormal termination will be included.
	Abnormal bool `protobuf:"varint,4,opt,name=abnormal" json:"abnormal,omitempty"`
	// If set to true, expired objects will be included.
	Expired bool `protobuf:"varint,5,opt,name=expired" json:"expired,omitempty"`
}

func (m *WalkGraphReq_Include_Options) Reset()         { *m = WalkGraphReq_Include_Options{} }
func (m *WalkGraphReq_Include_Options) String() string { return proto.CompactTextString(m) }
func (*WalkGraphReq_Include_Options) ProtoMessage()    {}
func (*WalkGraphReq_Include_Options) Descriptor() ([]byte, []int) {
	return fileDescriptor7, []int{0, 2, 0}
}

func (m *WalkGraphReq_Include_Options) GetIds() bool {
	if m != nil {
		return m.Ids
	}
	return false
}

func (m *WalkGraphReq_Include_Options) GetData() bool {
	if m != nil {
		return m.Data
	}
	return false
}

func (m *WalkGraphReq_Include_Options) GetResult() bool {
	if m != nil {
		return m.Result
	}
	return false
}

func (m *WalkGraphReq_Include_Options) GetAbnormal() bool {
	if m != nil {
		return m.Abnormal
	}
	return false
}

func (m *WalkGraphReq_Include_Options) GetExpired() bool {
	if m != nil {
		return m.Expired
	}
	return false
}

type WalkGraphReq_Exclude struct {
	// Do not include data from the following quests in the response.
	Quests []string `protobuf:"bytes,1,rep,name=quests" json:"quests,omitempty"`
	// Do not include data from the following attempts in the response.
	Attempts *AttemptList `protobuf:"bytes,2,opt,name=attempts" json:"attempts,omitempty"`
}

func (m *WalkGraphReq_Exclude) Reset()                    { *m = WalkGraphReq_Exclude{} }
func (m *WalkGraphReq_Exclude) String() string            { return proto.CompactTextString(m) }
func (*WalkGraphReq_Exclude) ProtoMessage()               {}
func (*WalkGraphReq_Exclude) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{0, 3} }

func (m *WalkGraphReq_Exclude) GetQuests() []string {
	if m != nil {
		return m.Quests
	}
	return nil
}

func (m *WalkGraphReq_Exclude) GetAttempts() *AttemptList {
	if m != nil {
		return m.Attempts
	}
	return nil
}

func init() {
	proto.RegisterType((*WalkGraphReq)(nil), "dm.WalkGraphReq")
	proto.RegisterType((*WalkGraphReq_Mode)(nil), "dm.WalkGraphReq.Mode")
	proto.RegisterType((*WalkGraphReq_Limit)(nil), "dm.WalkGraphReq.Limit")
	proto.RegisterType((*WalkGraphReq_Include)(nil), "dm.WalkGraphReq.Include")
	proto.RegisterType((*WalkGraphReq_Include_Options)(nil), "dm.WalkGraphReq.Include.Options")
	proto.RegisterType((*WalkGraphReq_Exclude)(nil), "dm.WalkGraphReq.Exclude")
	proto.RegisterEnum("dm.WalkGraphReq_Mode_Direction", WalkGraphReq_Mode_Direction_name, WalkGraphReq_Mode_Direction_value)
}

func init() {
	proto.RegisterFile("github.com/luci/luci-go/dm/api/service/v1/walk_graph.proto", fileDescriptor7)
}

var fileDescriptor7 = []byte{
	// 659 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x52, 0xdd, 0x6e, 0x13, 0x3d,
	0x10, 0xfd, 0xd2, 0xfc, 0xac, 0x33, 0x6d, 0xfa, 0x55, 0x96, 0xa8, 0xb6, 0x8b, 0x44, 0xab, 0x0a,
	0x50, 0x11, 0xb0, 0x2b, 0xc2, 0xcf, 0x45, 0x11, 0x48, 0x29, 0x29, 0x3f, 0xa2, 0x50, 0xe1, 0x56,
	0xea, 0x65, 0xe4, 0xac, 0x9d, 0xc4, 0xea, 0x3a, 0xbb, 0x5d, 0x7b, 0x9b, 0xb4, 0x12, 0xbc, 0x00,
	0x4f, 0x80, 0x78, 0x59, 0x64, 0xaf, 0x93, 0x22, 0x4a, 0x55, 0x7a, 0x63, 0x79, 0x66, 0xce, 0x99,
	0x19, 0x1f, 0x1f, 0xd8, 0x1e, 0x0a, 0x3d, 0x2a, 0xfa, 0x61, 0x9c, 0xca, 0x28, 0x29, 0x62, 0x61,
	0x8f, 0xc7, 0xc3, 0x34, 0x62, 0x32, 0xa2, 0x99, 0x88, 0x14, 0xcf, 0x4f, 0x45, 0xcc, 0xa3, 0xd3,
	0x27, 0xd1, 0x84, 0x26, 0xc7, 0xbd, 0x61, 0x4e, 0xb3, 0x51, 0x98, 0xe5, 0xa9, 0x4e, 0xf1, 0x02,
	0x93, 0xc1, 0x9d, 0x61, 0x9a, 0x0e, 0x13, 0x1e, 0xd9, 0x4c, 0xbf, 0x18, 0x44, 0xac, 0xc8, 0xa9,
	0x16, 0xe9, 0xb8, 0xc4, 0x04, 0xeb, 0x7f, 0xd6, 0xb5, 0x90, 0x5c, 0x69, 0x2a, 0x33, 0x07, 0xb8,
	0xc1, 0x02, 0x76, 0x76, 0x8f, 0x51, 0x4d, 0x1d, 0xf7, 0xe5, 0x4d, 0xb9, 0x27, 0x05, 0xcf, 0xcf,
	0x1c, 0xf9, 0xf9, 0xbf, 0x93, 0xf5, 0x59, 0xc6, 0x55, 0x49, 0xdb, 0xfc, 0x81, 0x60, 0xe9, 0x88,
	0x26, 0xc7, 0xef, 0x4c, 0x43, 0xc2, 0x4f, 0xf0, 0x7d, 0xa8, 0xd1, 0x42, 0x8f, 0xfc, 0xca, 0x46,
	0x65, 0x6b, 0xb1, 0x8d, 0x43, 0x26, 0xc3, 0xdd, 0x29, 0x8f, 0x0b, 0x2b, 0x42, 0xa7, 0xd0, 0x23,
	0x62, 0xeb, 0xf8, 0x2e, 0xd4, 0xed, 0x78, 0x7f, 0xc1, 0x02, 0x97, 0x0d, 0xd0, 0x36, 0xf9, 0x62,
	0xb2, 0xa4, 0x2c, 0xe2, 0x07, 0x50, 0x93, 0x29, 0xe3, 0x7e, 0xd5, 0x82, 0x6e, 0x19, 0xd0, 0xef,
	0xd3, 0xc2, 0x4f, 0x29, 0xe3, 0xc4, 0x42, 0xf0, 0x23, 0xa8, 0x27, 0x42, 0x0a, 0xed, 0xd7, 0x2c,
	0x76, 0xf5, 0x12, 0x76, 0xcf, 0x54, 0x49, 0x09, 0xc2, 0x6d, 0xf0, 0xc4, 0x38, 0x4e, 0x0a, 0xc6,
	0xfd, 0xba, 0xc5, 0xfb, 0x97, 0xf0, 0x1f, 0xca, 0x3a, 0x99, 0x01, 0x0d, 0x87, 0x4f, 0x4b, 0x4e,
	0xe3, 0x0a, 0xce, 0xee, 0xd4, 0x71, 0x1c, 0x30, 0xf8, 0x5e, 0x81, 0x9a, 0x59, 0x12, 0xaf, 0x40,
	0x95, 0x0d, 0x94, 0x95, 0x05, 0x11, 0x73, 0xc5, 0xaf, 0xa0, 0xc9, 0x44, 0xce, 0x63, 0xa3, 0x8c,
	0x55, 0x61, 0xb9, 0xbd, 0xfe, 0xd7, 0x07, 0x86, 0xdd, 0x19, 0x8c, 0x5c, 0x30, 0x36, 0xdb, 0xd0,
	0x9c, 0xe7, 0xf1, 0x12, 0xa0, 0xb7, 0xfb, 0xe4, 0xa8, 0x43, 0xba, 0x07, 0x2b, 0xff, 0xe1, 0x16,
	0x34, 0x77, 0x3a, 0x6f, 0x3e, 0x96, 0x61, 0x05, 0x23, 0xa8, 0xed, 0xec, 0x1f, 0xbe, 0x5f, 0x59,
	0x08, 0xbe, 0x41, 0xdd, 0xaa, 0x80, 0x6f, 0x43, 0x53, 0xd2, 0x69, 0x8f, 0xf1, 0xcc, 0x7d, 0x55,
	0x95, 0x20, 0x49, 0xa7, 0x5d, 0x13, 0xe3, 0x67, 0x60, 0xee, 0x3d, 0x63, 0x4d, 0xf7, 0x3b, 0x6b,
	0x61, 0xe9, 0xdb, 0x70, 0xe6, 0xdb, 0xb0, 0xeb, 0x7c, 0x4d, 0x3c, 0x49, 0xa7, 0x87, 0x42, 0x72,
	0xbc, 0x09, 0x2d, 0xdb, 0x92, 0x6a, 0xda, 0x53, 0xe2, 0xbc, 0xfc, 0xb3, 0x16, 0x59, 0x34, 0x6d,
	0xa9, 0xa6, 0x07, 0xe2, 0x9c, 0x07, 0x3f, 0xab, 0xe0, 0x39, 0x59, 0xf1, 0x0b, 0x6b, 0x00, 0xa5,
	0x9d, 0x53, 0x36, 0xae, 0xd2, 0x3f, 0xdc, 0xcf, 0xcc, 0x20, 0x45, 0x4a, 0x38, 0xde, 0x06, 0x8f,
	0x6a, 0xcd, 0x65, 0xa6, 0xdd, 0x72, 0xd7, 0x33, 0x67, 0x04, 0xfc, 0x1a, 0x9a, 0x7c, 0x66, 0x46,
	0xe7, 0xa9, 0xeb, 0xd9, 0x17, 0x14, 0x7c, 0x0f, 0x96, 0xc7, 0x85, 0xec, 0xcd, 0x13, 0xca, 0x9a,
	0xad, 0x45, 0x5a, 0xe3, 0x42, 0xce, 0x5d, 0xae, 0xf0, 0x1a, 0xa0, 0xc1, 0x84, 0x19, 0x75, 0x95,
	0x75, 0x17, 0x22, 0xde, 0x60, 0xc2, 0xba, 0x3c, 0x53, 0x46, 0xf8, 0x3e, 0x8d, 0x8f, 0xcb, 0x5a,
	0xc3, 0xd6, 0x90, 0x49, 0x98, 0x62, 0xf0, 0x15, 0x3c, 0x37, 0xd4, 0xd8, 0x45, 0xb0, 0xb9, 0x5d,
	0x04, 0x53, 0x18, 0x43, 0xcd, 0x68, 0x6b, 0x1f, 0x8d, 0x88, 0xbd, 0xe3, 0x55, 0x68, 0xe4, 0x5c,
	0x15, 0x89, 0xb6, 0x8f, 0x41, 0xc4, 0x45, 0x38, 0x00, 0x44, 0xfb, 0xe3, 0x34, 0x97, 0x34, 0xb1,
	0x1b, 0x22, 0x32, 0x8f, 0xb1, 0x6f, 0x5c, 0x9c, 0x89, 0x9c, 0xb3, 0xd9, 0x6e, 0x2e, 0x0c, 0x3e,
	0x83, 0xe7, 0xfc, 0x6b, 0x1a, 0x5b, 0xb5, 0xcd, 0x06, 0xd5, 0xad, 0x26, 0x71, 0x11, 0x7e, 0x08,
	0xc8, 0x69, 0xa9, 0x9c, 0xfa, 0xff, 0x1b, 0xfd, 0x3a, 0x65, 0x6e, 0x4f, 0x28, 0x4d, 0xe6, 0x80,
	0x7e, 0xc3, 0xba, 0xe5, 0xe9, 0xaf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xcd, 0xfd, 0x53, 0x55,
	0x05, 0x00, 0x00,
}
