// Code generated by protoc-gen-go.
// source: github.com/luci/luci-go/logdog/api/logpb/log.proto
// DO NOT EDIT!

package logpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/luci/luci-go/common/proto/google"
import google_protobuf1 "github.com/luci/luci-go/common/proto/google"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// A log stream type.
type StreamType int32

const (
	StreamType_TEXT     StreamType = 0
	StreamType_BINARY   StreamType = 1
	StreamType_DATAGRAM StreamType = 2
)

var StreamType_name = map[int32]string{
	0: "TEXT",
	1: "BINARY",
	2: "DATAGRAM",
}
var StreamType_value = map[string]int32{
	"TEXT":     0,
	"BINARY":   1,
	"DATAGRAM": 2,
}

func (x StreamType) String() string {
	return proto.EnumName(StreamType_name, int32(x))
}
func (StreamType) EnumDescriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

// *
// Log stream descriptor data. This is the full set of information that
// describes a logging stream.
type LogStreamDescriptor struct {
	//
	// The stream's prefix (required).
	//
	// Logs originating from the same Butler instance will share a Prefix.
	//
	// A valid prefix value is a StreamName described in:
	// https://github.com/luci/luci-go/common/logdog/types
	Prefix string `protobuf:"bytes,1,opt,name=prefix" json:"prefix,omitempty"`
	//
	// The log stream's name (required).
	//
	// This is used to uniquely identify a log stream within the scope of its
	// prefix.
	//
	// A valid name value is a StreamName described in:
	// https://github.com/luci/luci-go/common/logdog/types
	Name string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	// The log stream's content type (required).
	StreamType StreamType `protobuf:"varint,3,opt,name=stream_type,json=streamType,enum=logpb.StreamType" json:"stream_type,omitempty"`
	//
	// The stream's content type (required).
	//
	// This must be an HTTP Content-Type value. It is made available to LogDog
	// clients when querying stream metadata. It will also be applied to archived
	// binary log data.
	ContentType string `protobuf:"bytes,4,opt,name=content_type,json=contentType" json:"content_type,omitempty"`
	//
	// The log stream's base timestamp (required).
	//
	// This notes the start time of the log stream. All LogEntries express their
	// timestamp as microsecond offsets from this field.
	Timestamp *google_protobuf.Timestamp `protobuf:"bytes,5,opt,name=timestamp" json:"timestamp,omitempty"`
	//
	// Tag is an arbitrary key/value tag associated with this log stream.
	//
	// LogDog clients can query for log streams based on tag values.
	Tags map[string]string `protobuf:"bytes,6,rep,name=tags" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	//
	// If set, the stream will be joined together during archival to recreate the
	// original stream and made available at <prefix>/+/<name>.ext.
	BinaryFileExt string `protobuf:"bytes,7,opt,name=binary_file_ext,json=binaryFileExt" json:"binary_file_ext,omitempty"`
}

func (m *LogStreamDescriptor) Reset()                    { *m = LogStreamDescriptor{} }
func (m *LogStreamDescriptor) String() string            { return proto.CompactTextString(m) }
func (*LogStreamDescriptor) ProtoMessage()               {}
func (*LogStreamDescriptor) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *LogStreamDescriptor) GetPrefix() string {
	if m != nil {
		return m.Prefix
	}
	return ""
}

func (m *LogStreamDescriptor) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *LogStreamDescriptor) GetStreamType() StreamType {
	if m != nil {
		return m.StreamType
	}
	return StreamType_TEXT
}

func (m *LogStreamDescriptor) GetContentType() string {
	if m != nil {
		return m.ContentType
	}
	return ""
}

func (m *LogStreamDescriptor) GetTimestamp() *google_protobuf.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *LogStreamDescriptor) GetTags() map[string]string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *LogStreamDescriptor) GetBinaryFileExt() string {
	if m != nil {
		return m.BinaryFileExt
	}
	return ""
}

// Text stream content.
type Text struct {
	Lines []*Text_Line `protobuf:"bytes,1,rep,name=lines" json:"lines,omitempty"`
}

func (m *Text) Reset()                    { *m = Text{} }
func (m *Text) String() string            { return proto.CompactTextString(m) }
func (*Text) ProtoMessage()               {}
func (*Text) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *Text) GetLines() []*Text_Line {
	if m != nil {
		return m.Lines
	}
	return nil
}

// Contiguous text lines and their delimiters.
type Text_Line struct {
	// The line's text content, not including its delimiter.
	Value string `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
	//
	// The line's delimiter string.
	//
	// If this is an empty string, this line is continued in the next sequential
	// line, and the line's sequence number does not advance.
	Delimiter string `protobuf:"bytes,2,opt,name=delimiter" json:"delimiter,omitempty"`
}

func (m *Text_Line) Reset()                    { *m = Text_Line{} }
func (m *Text_Line) String() string            { return proto.CompactTextString(m) }
func (*Text_Line) ProtoMessage()               {}
func (*Text_Line) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1, 0} }

func (m *Text_Line) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *Text_Line) GetDelimiter() string {
	if m != nil {
		return m.Delimiter
	}
	return ""
}

// Binary stream content.
type Binary struct {
	// The byte offset in the stream of the first byte of data.
	Offset uint64 `protobuf:"varint,1,opt,name=offset" json:"offset,omitempty"`
	// The binary stream's data.
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *Binary) Reset()                    { *m = Binary{} }
func (m *Binary) String() string            { return proto.CompactTextString(m) }
func (*Binary) ProtoMessage()               {}
func (*Binary) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *Binary) GetOffset() uint64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *Binary) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

// Datagram stream content type.
type Datagram struct {
	// This datagram data.
	Data    []byte            `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Partial *Datagram_Partial `protobuf:"bytes,2,opt,name=partial" json:"partial,omitempty"`
}

func (m *Datagram) Reset()                    { *m = Datagram{} }
func (m *Datagram) String() string            { return proto.CompactTextString(m) }
func (*Datagram) ProtoMessage()               {}
func (*Datagram) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func (m *Datagram) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Datagram) GetPartial() *Datagram_Partial {
	if m != nil {
		return m.Partial
	}
	return nil
}

//
// If this is not a partial datagram, this field will include reassembly and
// state details for the full datagram.
type Datagram_Partial struct {
	//
	// The index, starting with zero, of this datagram fragment in the full
	// datagram.
	Index uint32 `protobuf:"varint,1,opt,name=index" json:"index,omitempty"`
	// The size of the full datagram
	Size uint64 `protobuf:"varint,2,opt,name=size" json:"size,omitempty"`
	// If true, this is the last partial datagram in the overall datagram.
	Last bool `protobuf:"varint,3,opt,name=last" json:"last,omitempty"`
}

func (m *Datagram_Partial) Reset()                    { *m = Datagram_Partial{} }
func (m *Datagram_Partial) String() string            { return proto.CompactTextString(m) }
func (*Datagram_Partial) ProtoMessage()               {}
func (*Datagram_Partial) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3, 0} }

func (m *Datagram_Partial) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *Datagram_Partial) GetSize() uint64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *Datagram_Partial) GetLast() bool {
	if m != nil {
		return m.Last
	}
	return false
}

// *
// An individual log entry.
//
// This contains the superset of transmissible log data. Its content fields
// should be interpreted in the context of the log stream's content type.
type LogEntry struct {
	//
	// The stream time offset for this content.
	//
	// This offset is added to the log stream's base "timestamp" to resolve the
	// timestamp for this specific Content.
	TimeOffset *google_protobuf1.Duration `protobuf:"bytes,1,opt,name=time_offset,json=timeOffset" json:"time_offset,omitempty"`
	//
	// The message index within the Prefix (required).
	//
	// This is value is unique to this LogEntry across the entire set of entries
	// sharing the stream's Prefix. It is used to designate unambiguous log
	// ordering.
	PrefixIndex uint64 `protobuf:"varint,2,opt,name=prefix_index,json=prefixIndex" json:"prefix_index,omitempty"`
	//
	// The message index within its Stream (required).
	//
	// This value is unique across all entries sharing the same Prefix and Stream
	// Name. It is used to designate unambiguous log ordering within the stream.
	StreamIndex uint64 `protobuf:"varint,3,opt,name=stream_index,json=streamIndex" json:"stream_index,omitempty"`
	//
	// The sequence number of the first content entry in this LogEntry.
	//
	// Text: This is the line index of the first included line. Line indices begin
	//     at zero.
	// Binary: This is the byte offset of the first byte in the included data.
	// Datagram: This is the index of the datagram. The first datagram has index
	//     zero.
	Sequence uint64 `protobuf:"varint,4,opt,name=sequence" json:"sequence,omitempty"`
	//
	// The content of the message. The field that is populated here must
	// match the log's `stream_type`.
	//
	// Types that are valid to be assigned to Content:
	//	*LogEntry_Text
	//	*LogEntry_Binary
	//	*LogEntry_Datagram
	Content isLogEntry_Content `protobuf_oneof:"content"`
}

func (m *LogEntry) Reset()                    { *m = LogEntry{} }
func (m *LogEntry) String() string            { return proto.CompactTextString(m) }
func (*LogEntry) ProtoMessage()               {}
func (*LogEntry) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

type isLogEntry_Content interface {
	isLogEntry_Content()
}

type LogEntry_Text struct {
	Text *Text `protobuf:"bytes,10,opt,name=text,oneof"`
}
type LogEntry_Binary struct {
	Binary *Binary `protobuf:"bytes,11,opt,name=binary,oneof"`
}
type LogEntry_Datagram struct {
	Datagram *Datagram `protobuf:"bytes,12,opt,name=datagram,oneof"`
}

func (*LogEntry_Text) isLogEntry_Content()     {}
func (*LogEntry_Binary) isLogEntry_Content()   {}
func (*LogEntry_Datagram) isLogEntry_Content() {}

func (m *LogEntry) GetContent() isLogEntry_Content {
	if m != nil {
		return m.Content
	}
	return nil
}

func (m *LogEntry) GetTimeOffset() *google_protobuf1.Duration {
	if m != nil {
		return m.TimeOffset
	}
	return nil
}

func (m *LogEntry) GetPrefixIndex() uint64 {
	if m != nil {
		return m.PrefixIndex
	}
	return 0
}

func (m *LogEntry) GetStreamIndex() uint64 {
	if m != nil {
		return m.StreamIndex
	}
	return 0
}

func (m *LogEntry) GetSequence() uint64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *LogEntry) GetText() *Text {
	if x, ok := m.GetContent().(*LogEntry_Text); ok {
		return x.Text
	}
	return nil
}

func (m *LogEntry) GetBinary() *Binary {
	if x, ok := m.GetContent().(*LogEntry_Binary); ok {
		return x.Binary
	}
	return nil
}

func (m *LogEntry) GetDatagram() *Datagram {
	if x, ok := m.GetContent().(*LogEntry_Datagram); ok {
		return x.Datagram
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*LogEntry) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _LogEntry_OneofMarshaler, _LogEntry_OneofUnmarshaler, _LogEntry_OneofSizer, []interface{}{
		(*LogEntry_Text)(nil),
		(*LogEntry_Binary)(nil),
		(*LogEntry_Datagram)(nil),
	}
}

func _LogEntry_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*LogEntry)
	// content
	switch x := m.Content.(type) {
	case *LogEntry_Text:
		b.EncodeVarint(10<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Text); err != nil {
			return err
		}
	case *LogEntry_Binary:
		b.EncodeVarint(11<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Binary); err != nil {
			return err
		}
	case *LogEntry_Datagram:
		b.EncodeVarint(12<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Datagram); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("LogEntry.Content has unexpected type %T", x)
	}
	return nil
}

func _LogEntry_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*LogEntry)
	switch tag {
	case 10: // content.text
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Text)
		err := b.DecodeMessage(msg)
		m.Content = &LogEntry_Text{msg}
		return true, err
	case 11: // content.binary
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Binary)
		err := b.DecodeMessage(msg)
		m.Content = &LogEntry_Binary{msg}
		return true, err
	case 12: // content.datagram
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Datagram)
		err := b.DecodeMessage(msg)
		m.Content = &LogEntry_Datagram{msg}
		return true, err
	default:
		return false, nil
	}
}

func _LogEntry_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*LogEntry)
	// content
	switch x := m.Content.(type) {
	case *LogEntry_Text:
		s := proto.Size(x.Text)
		n += proto.SizeVarint(10<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *LogEntry_Binary:
		s := proto.Size(x.Binary)
		n += proto.SizeVarint(11<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *LogEntry_Datagram:
		s := proto.Size(x.Datagram)
		n += proto.SizeVarint(12<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// *
// LogIndex is an index into an at-rest log storage.
//
// The log stream and log index are generated by the Archivist during archival.
//
// An archived log stream is a series of contiguous LogEntry frames. The index
// maps a log's logical logation in its stream, prefix, and timeline to its
// frame's binary offset in the archived log stream blob.
type LogIndex struct {
	//
	// The LogStreamDescriptor for this log stream (required).
	//
	// The index stores the stream's LogStreamDescriptor so that a client can
	// know the full set of log metadata by downloading its index.
	Desc *LogStreamDescriptor `protobuf:"bytes,1,opt,name=desc" json:"desc,omitempty"`
	//
	// A series of ascending-ordered Entry messages representing snapshots of an
	// archived log stream.
	//
	// Within this set of Entry messages, the "offset", "prefix_index",
	// "stream_index", and "time_offset" fields will be ascending.
	//
	// The frequency of Entry messages is not defined; it is up to the Archivist
	// process to choose a frequency.
	Entries []*LogIndex_Entry `protobuf:"bytes,2,rep,name=entries" json:"entries,omitempty"`
	// *
	// The last prefix index in the log stream.
	//
	// This is optional. If zero, there is either no information about the last
	// prefix index, or there are zero entries in the prefix.
	LastPrefixIndex uint64 `protobuf:"varint,3,opt,name=last_prefix_index,json=lastPrefixIndex" json:"last_prefix_index,omitempty"`
	// *
	// The last stream index in the log stream.
	//
	// This is optional. If zero, there is either no information about the last
	// stream index, or there are zero entries in the stream.
	LastStreamIndex uint64 `protobuf:"varint,4,opt,name=last_stream_index,json=lastStreamIndex" json:"last_stream_index,omitempty"`
	// *
	// The number of log entries in the stream.
	//
	// This is optional. If zero, there is either no information about the number
	// of log entries, or there are zero entries in the stream.
	LogEntryCount uint64 `protobuf:"varint,5,opt,name=log_entry_count,json=logEntryCount" json:"log_entry_count,omitempty"`
}

func (m *LogIndex) Reset()                    { *m = LogIndex{} }
func (m *LogIndex) String() string            { return proto.CompactTextString(m) }
func (*LogIndex) ProtoMessage()               {}
func (*LogIndex) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{5} }

func (m *LogIndex) GetDesc() *LogStreamDescriptor {
	if m != nil {
		return m.Desc
	}
	return nil
}

func (m *LogIndex) GetEntries() []*LogIndex_Entry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func (m *LogIndex) GetLastPrefixIndex() uint64 {
	if m != nil {
		return m.LastPrefixIndex
	}
	return 0
}

func (m *LogIndex) GetLastStreamIndex() uint64 {
	if m != nil {
		return m.LastStreamIndex
	}
	return 0
}

func (m *LogIndex) GetLogEntryCount() uint64 {
	if m != nil {
		return m.LogEntryCount
	}
	return 0
}

//
// Entry is a single index entry.
//
// The index is composed of a series of entries, each corresponding to a
// sequential snapshot of of the log stream.
type LogIndex_Entry struct {
	//
	// The byte offset in the emitted log stream of the RecordIO entry for the
	// LogEntry corresponding to this Entry.
	Offset uint64 `protobuf:"varint,1,opt,name=offset" json:"offset,omitempty"`
	//
	// The sequence number of the first content entry.
	//
	// Text: This is the line index of the first included line. Line indices
	//     begin at zero.
	// Binary: This is the byte offset of the first byte in the included data.
	// Datagram: This is the index of the datagram. The first datagram has index
	//     zero.
	Sequence uint64 `protobuf:"varint,2,opt,name=sequence" json:"sequence,omitempty"`
	//
	// The log index that this entry describes (required).
	//
	// This is used by clients to identify a specific LogEntry within a set of
	// streams sharing a Prefix.
	PrefixIndex uint64 `protobuf:"varint,3,opt,name=prefix_index,json=prefixIndex" json:"prefix_index,omitempty"`
	//
	// The time offset of this log entry (required).
	//
	// This is used by clients to identify a specific LogEntry within a log
	// stream.
	StreamIndex uint64 `protobuf:"varint,4,opt,name=stream_index,json=streamIndex" json:"stream_index,omitempty"`
	//
	// The time offset of this log entry, in microseconds.
	//
	// This is added to the descriptor's "timestamp" field to identify the
	// specific timestamp of this log. It is used by clients to identify a
	// specific LogEntry by time.
	TimeOffset *google_protobuf1.Duration `protobuf:"bytes,5,opt,name=time_offset,json=timeOffset" json:"time_offset,omitempty"`
}

func (m *LogIndex_Entry) Reset()                    { *m = LogIndex_Entry{} }
func (m *LogIndex_Entry) String() string            { return proto.CompactTextString(m) }
func (*LogIndex_Entry) ProtoMessage()               {}
func (*LogIndex_Entry) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{5, 0} }

func (m *LogIndex_Entry) GetOffset() uint64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *LogIndex_Entry) GetSequence() uint64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *LogIndex_Entry) GetPrefixIndex() uint64 {
	if m != nil {
		return m.PrefixIndex
	}
	return 0
}

func (m *LogIndex_Entry) GetStreamIndex() uint64 {
	if m != nil {
		return m.StreamIndex
	}
	return 0
}

func (m *LogIndex_Entry) GetTimeOffset() *google_protobuf1.Duration {
	if m != nil {
		return m.TimeOffset
	}
	return nil
}

func init() {
	proto.RegisterType((*LogStreamDescriptor)(nil), "logpb.LogStreamDescriptor")
	proto.RegisterType((*Text)(nil), "logpb.Text")
	proto.RegisterType((*Text_Line)(nil), "logpb.Text.Line")
	proto.RegisterType((*Binary)(nil), "logpb.Binary")
	proto.RegisterType((*Datagram)(nil), "logpb.Datagram")
	proto.RegisterType((*Datagram_Partial)(nil), "logpb.Datagram.Partial")
	proto.RegisterType((*LogEntry)(nil), "logpb.LogEntry")
	proto.RegisterType((*LogIndex)(nil), "logpb.LogIndex")
	proto.RegisterType((*LogIndex_Entry)(nil), "logpb.LogIndex.Entry")
	proto.RegisterEnum("logpb.StreamType", StreamType_name, StreamType_value)
}

func init() { proto.RegisterFile("github.com/luci/luci-go/logdog/api/logpb/log.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 782 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x54, 0xdb, 0x6e, 0xdb, 0x46,
	0x10, 0x15, 0x25, 0xea, 0x36, 0x94, 0x2a, 0x79, 0x7b, 0x63, 0x89, 0xa2, 0x95, 0x84, 0xc2, 0x15,
	0x0c, 0x98, 0x6a, 0xd5, 0x02, 0x31, 0xfc, 0x26, 0x47, 0x8e, 0x6d, 0xc0, 0x49, 0x8c, 0xb5, 0x1e,
	0x92, 0x27, 0x62, 0x25, 0xad, 0x98, 0x4d, 0x28, 0x2e, 0x43, 0xae, 0x02, 0x29, 0x9f, 0x92, 0x5f,
	0x08, 0x90, 0x3f, 0xc8, 0x27, 0xe5, 0x1f, 0x82, 0xbd, 0xe8, 0xe2, 0x5b, 0x80, 0xbc, 0x10, 0xb3,
	0x33, 0x67, 0x66, 0x67, 0xce, 0x9c, 0x25, 0xf4, 0x43, 0x26, 0x5e, 0x2d, 0xc6, 0xfe, 0x84, 0xcf,
	0x7b, 0xd1, 0x62, 0xc2, 0xd4, 0xe7, 0x30, 0xe4, 0xbd, 0x88, 0x87, 0x53, 0x1e, 0xf6, 0x48, 0xc2,
	0xa4, 0x99, 0x8c, 0xe5, 0xd7, 0x4f, 0x52, 0x2e, 0x38, 0x2a, 0x2a, 0x87, 0xf7, 0x67, 0xc8, 0x79,
	0x18, 0xd1, 0x9e, 0x72, 0x8e, 0x17, 0xb3, 0x9e, 0x60, 0x73, 0x9a, 0x09, 0x32, 0x4f, 0x34, 0xce,
	0xfb, 0xe3, 0x36, 0x60, 0xba, 0x48, 0x89, 0x60, 0x3c, 0xd6, 0xf1, 0xce, 0x97, 0x3c, 0xfc, 0x78,
	0xc9, 0xc3, 0x6b, 0x91, 0x52, 0x32, 0x1f, 0xd2, 0x6c, 0x92, 0xb2, 0x44, 0xf0, 0x14, 0xfd, 0x02,
	0xa5, 0x24, 0xa5, 0x33, 0xb6, 0x74, 0xad, 0x96, 0xd5, 0xad, 0x62, 0x73, 0x42, 0x08, 0xec, 0x98,
	0xcc, 0xa9, 0x9b, 0x57, 0x5e, 0x65, 0xa3, 0x3e, 0x38, 0x99, 0xca, 0x0f, 0xc4, 0x2a, 0xa1, 0x6e,
	0xa1, 0x65, 0x75, 0x7f, 0xe8, 0xef, 0xf9, 0xaa, 0x43, 0x5f, 0x57, 0x1e, 0xad, 0x12, 0x8a, 0x21,
	0xdb, 0xd8, 0xa8, 0x0d, 0xb5, 0x09, 0x8f, 0x05, 0x8d, 0x85, 0x4e, 0xb2, 0x55, 0x3d, 0xc7, 0xf8,
	0x14, 0xe4, 0x08, 0xaa, 0x9b, 0x69, 0xdc, 0x62, 0xcb, 0xea, 0x3a, 0x7d, 0xcf, 0xd7, 0xe3, 0xf8,
	0xeb, 0x71, 0xfc, 0xd1, 0x1a, 0x81, 0xb7, 0x60, 0x74, 0x04, 0xb6, 0x20, 0x61, 0xe6, 0x96, 0x5a,
	0x85, 0xae, 0xd3, 0xff, 0xcb, 0x74, 0x72, 0xcf, 0x98, 0xfe, 0x88, 0x84, 0xd9, 0x69, 0x2c, 0xd2,
	0x15, 0x56, 0x19, 0x68, 0x1f, 0x1a, 0x63, 0x16, 0x93, 0x74, 0x15, 0xcc, 0x58, 0x44, 0x03, 0xba,
	0x14, 0x6e, 0x59, 0x75, 0x56, 0xd7, 0xee, 0x27, 0x2c, 0xa2, 0xa7, 0x4b, 0xe1, 0x3d, 0x82, 0xea,
	0x26, 0x15, 0x35, 0xa1, 0xf0, 0x86, 0xae, 0x0c, 0x51, 0xd2, 0x44, 0x3f, 0x41, 0xf1, 0x1d, 0x89,
	0x16, 0x6b, 0x9a, 0xf4, 0xe1, 0x38, 0x7f, 0x64, 0x75, 0x5e, 0x83, 0x3d, 0xa2, 0x4b, 0x81, 0xf6,
	0xa1, 0x18, 0xb1, 0x98, 0x66, 0xae, 0xa5, 0x7a, 0x6c, 0x9a, 0x1e, 0x65, 0xcc, 0xbf, 0x64, 0x31,
	0xc5, 0x3a, 0xec, 0x1d, 0x83, 0x2d, 0x8f, 0xdb, 0x8a, 0xd6, 0x4e, 0x45, 0xf4, 0x3b, 0x54, 0xa7,
	0x34, 0x62, 0x73, 0x26, 0x68, 0x6a, 0xee, 0xda, 0x3a, 0x3a, 0xff, 0x43, 0xe9, 0x44, 0x75, 0x2d,
	0xb7, 0xc9, 0x67, 0xb3, 0x8c, 0x0a, 0x95, 0x6e, 0x63, 0x73, 0x92, 0xdb, 0x9c, 0x12, 0x41, 0x54,
	0x6a, 0x0d, 0x2b, 0xbb, 0xf3, 0xc1, 0x82, 0xca, 0x90, 0x08, 0x12, 0xa6, 0x64, 0xbe, 0x01, 0x58,
	0x5b, 0x00, 0xfa, 0x17, 0xca, 0x09, 0x49, 0x05, 0x23, 0x91, 0xca, 0x73, 0xfa, 0xbf, 0x9a, 0xe6,
	0xd7, 0x59, 0xfe, 0x95, 0x0e, 0xe3, 0x35, 0xce, 0x3b, 0x83, 0xb2, 0xf1, 0xc9, 0x41, 0x58, 0x3c,
	0xa5, 0x5a, 0x57, 0x75, 0xac, 0x0f, 0xf2, 0x9e, 0x8c, 0xbd, 0xd7, 0x7c, 0xd9, 0x58, 0xd9, 0xd2,
	0x17, 0x91, 0x4c, 0x28, 0x3d, 0x55, 0xb0, 0xb2, 0x3b, 0x9f, 0xf2, 0x50, 0xb9, 0xe4, 0xa1, 0xe6,
	0xfd, 0x18, 0x1c, 0xb9, 0xf3, 0x60, 0x67, 0x34, 0xa7, 0xff, 0xdb, 0x1d, 0x89, 0x0c, 0x8d, 0xe2,
	0x31, 0x48, 0xf4, 0x73, 0x3d, 0x79, 0x1b, 0x6a, 0x5a, 0xd1, 0x81, 0xee, 0x46, 0x5f, 0xec, 0x68,
	0xdf, 0x85, 0xea, 0xa9, 0x0d, 0x35, 0x23, 0x6b, 0x0d, 0x29, 0x68, 0x88, 0xf6, 0x69, 0x88, 0x07,
	0x95, 0x8c, 0xbe, 0x5d, 0xd0, 0x78, 0xa2, 0x15, 0x6c, 0xe3, 0xcd, 0x19, 0xb5, 0xc1, 0x16, 0x52,
	0x3f, 0xa0, 0xda, 0x72, 0x76, 0x16, 0x7c, 0x9e, 0xc3, 0x2a, 0x84, 0xfe, 0x86, 0x92, 0x96, 0x95,
	0xeb, 0x28, 0x50, 0xdd, 0x80, 0xf4, 0xd6, 0xce, 0x73, 0xd8, 0x84, 0xd1, 0x21, 0x54, 0xa6, 0x86,
	0x5c, 0xb7, 0xa6, 0xa0, 0x8d, 0x5b, 0x9c, 0x9f, 0xe7, 0xf0, 0x06, 0x72, 0x52, 0x85, 0xb2, 0x79,
	0x48, 0x9d, 0x8f, 0x05, 0x45, 0x98, 0x6e, 0xd7, 0x07, 0x7b, 0x4a, 0xb3, 0x89, 0x61, 0xca, 0x7b,
	0xf8, 0x5d, 0x60, 0x85, 0x43, 0x3d, 0x28, 0xd3, 0x58, 0xa4, 0x8c, 0x66, 0x6e, 0x5e, 0xc9, 0xf4,
	0xe7, 0x6d, 0x8a, 0xaa, 0xe8, 0xeb, 0xb7, 0xb3, 0x46, 0xa1, 0x03, 0xd8, 0x93, 0x6b, 0x0a, 0x6e,
	0x50, 0xab, 0x79, 0x6b, 0xc8, 0xc0, 0xd5, 0x0e, 0xbd, 0x6b, 0xec, 0x0d, 0x8e, 0xed, 0x2d, 0xf6,
	0x7a, 0x87, 0xe7, 0x7d, 0x68, 0x44, 0x3c, 0x0c, 0xe4, 0x35, 0xab, 0x60, 0xc2, 0x17, 0xb1, 0x50,
	0x3f, 0x04, 0x1b, 0xd7, 0x23, 0x23, 0x86, 0xc7, 0xd2, 0xe9, 0x7d, 0xb6, 0xa0, 0xa8, 0xb5, 0xf1,
	0x90, 0xe2, 0x77, 0x37, 0x96, 0xbf, 0xb3, 0xb1, 0xda, 0x3d, 0x8d, 0x7f, 0x53, 0x13, 0xf6, 0x5d,
	0x4d, 0xdc, 0x52, 0x65, 0xf1, 0x3b, 0x54, 0x79, 0xf0, 0x0f, 0xc0, 0xf6, 0x7f, 0x89, 0x2a, 0x60,
	0x8f, 0x4e, 0x5f, 0x8c, 0x9a, 0x39, 0x04, 0x50, 0x3a, 0xb9, 0x78, 0x36, 0xc0, 0x2f, 0x9b, 0x16,
	0xaa, 0x41, 0x65, 0x38, 0x18, 0x0d, 0xce, 0xf0, 0xe0, 0x69, 0x33, 0x3f, 0x2e, 0xa9, 0x82, 0xff,
	0x7d, 0x0d, 0x00, 0x00, 0xff, 0xff, 0xfb, 0x12, 0x93, 0xb0, 0x44, 0x06, 0x00, 0x00,
}
