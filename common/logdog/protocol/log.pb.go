// Code generated by protoc-gen-go.
// source: log.proto
// DO NOT EDIT!

/*
Package protocol is a generated protocol buffer package.

It is generated from these files:
	log.proto
	butler.proto

It has these top-level messages:
	LogStreamDescriptor
	Text
	Binary
	Datagram
	LogEntry
	LogIndex
*/
package protocol

import proto "github.com/golang/protobuf/proto"
import google_protobuf "github.com/luci/luci-go/common/proto/google"
import google_protobuf1 "github.com/luci/luci-go/common/proto/google"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

// The log stream's content type (required).
type LogStreamDescriptor_StreamType int32

const (
	LogStreamDescriptor_TEXT     LogStreamDescriptor_StreamType = 0
	LogStreamDescriptor_BINARY   LogStreamDescriptor_StreamType = 1
	LogStreamDescriptor_DATAGRAM LogStreamDescriptor_StreamType = 2
)

var LogStreamDescriptor_StreamType_name = map[int32]string{
	0: "TEXT",
	1: "BINARY",
	2: "DATAGRAM",
}
var LogStreamDescriptor_StreamType_value = map[string]int32{
	"TEXT":     0,
	"BINARY":   1,
	"DATAGRAM": 2,
}

func (x LogStreamDescriptor_StreamType) String() string {
	return proto.EnumName(LogStreamDescriptor_StreamType_name, int32(x))
}

type LogStreamDescriptor_TextSource_Newline int32

const (
	LogStreamDescriptor_TextSource_NATIVE LogStreamDescriptor_TextSource_Newline = 0
	LogStreamDescriptor_TextSource_LF     LogStreamDescriptor_TextSource_Newline = 1
	LogStreamDescriptor_TextSource_CRLF   LogStreamDescriptor_TextSource_Newline = 2
)

var LogStreamDescriptor_TextSource_Newline_name = map[int32]string{
	0: "NATIVE",
	1: "LF",
	2: "CRLF",
}
var LogStreamDescriptor_TextSource_Newline_value = map[string]int32{
	"NATIVE": 0,
	"LF":     1,
	"CRLF":   2,
}

func (x LogStreamDescriptor_TextSource_Newline) String() string {
	return proto.EnumName(LogStreamDescriptor_TextSource_Newline_name, int32(x))
}

// The source stream encoding.
type LogStreamDescriptor_TextSource_Encoding int32

const (
	LogStreamDescriptor_TextSource_UTF8 LogStreamDescriptor_TextSource_Encoding = 0
)

var LogStreamDescriptor_TextSource_Encoding_name = map[int32]string{
	0: "UTF8",
}
var LogStreamDescriptor_TextSource_Encoding_value = map[string]int32{
	"UTF8": 0,
}

func (x LogStreamDescriptor_TextSource_Encoding) String() string {
	return proto.EnumName(LogStreamDescriptor_TextSource_Encoding_name, int32(x))
}

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
	Name       string                          `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	StreamType LogStreamDescriptor_StreamType  `protobuf:"varint,3,opt,name=stream_type,enum=protocol.LogStreamDescriptor_StreamType" json:"stream_type,omitempty"`
	TextSource *LogStreamDescriptor_TextSource `protobuf:"bytes,4,opt,name=text_source" json:"text_source,omitempty"`
	//
	// The stream's content type (required).
	//
	// This must be an HTTP Content-Type value. It is made available to LogDog
	// clients when querying stream metadata. It will also be applied to archived
	// binary log data.
	ContentType string `protobuf:"bytes,5,opt,name=content_type" json:"content_type,omitempty"`
	//
	// The log stream's base timestamp (required).
	//
	// This notes the start time of the log stream. All LogEntries express their
	// timestamp as microsecond offsets from this field.
	Timestamp *google_protobuf.Timestamp `protobuf:"bytes,6,opt,name=timestamp" json:"timestamp,omitempty"`
	// The set of associated log tags.
	Tags []*LogStreamDescriptor_Tag `protobuf:"bytes,7,rep,name=tags" json:"tags,omitempty"`
	//
	// If set, the stream will be joined together during archival to recreate the
	// original stream and made available at <prefix>/+/<name>.ext.
	BinaryFileExt string `protobuf:"bytes,8,opt,name=binary_file_ext" json:"binary_file_ext,omitempty"`
}

func (m *LogStreamDescriptor) Reset()         { *m = LogStreamDescriptor{} }
func (m *LogStreamDescriptor) String() string { return proto.CompactTextString(m) }
func (*LogStreamDescriptor) ProtoMessage()    {}

func (m *LogStreamDescriptor) GetTextSource() *LogStreamDescriptor_TextSource {
	if m != nil {
		return m.TextSource
	}
	return nil
}

func (m *LogStreamDescriptor) GetTimestamp() *google_protobuf.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *LogStreamDescriptor) GetTags() []*LogStreamDescriptor_Tag {
	if m != nil {
		return m.Tags
	}
	return nil
}

//
// Additional information about the source data for a text log stream.
//
// Required if `stream_type` is `TEXT`.
type LogStreamDescriptor_TextSource struct {
	Newline  LogStreamDescriptor_TextSource_Newline  `protobuf:"varint,1,opt,name=newline,enum=protocol.LogStreamDescriptor_TextSource_Newline" json:"newline,omitempty"`
	Encoding LogStreamDescriptor_TextSource_Encoding `protobuf:"varint,2,opt,name=encoding,enum=protocol.LogStreamDescriptor_TextSource_Encoding" json:"encoding,omitempty"`
}

func (m *LogStreamDescriptor_TextSource) Reset()         { *m = LogStreamDescriptor_TextSource{} }
func (m *LogStreamDescriptor_TextSource) String() string { return proto.CompactTextString(m) }
func (*LogStreamDescriptor_TextSource) ProtoMessage()    {}

//
// Tag is an arbitrary key/value tag associated with this log stream.
//
// LogDog clients can query for log streams based on tag values.
type LogStreamDescriptor_Tag struct {
	// The tag key (required).
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	// The tag value.
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *LogStreamDescriptor_Tag) Reset()         { *m = LogStreamDescriptor_Tag{} }
func (m *LogStreamDescriptor_Tag) String() string { return proto.CompactTextString(m) }
func (*LogStreamDescriptor_Tag) ProtoMessage()    {}

// Text stream content.
type Text struct {
	Lines []string `protobuf:"bytes,1,rep,name=lines" json:"lines,omitempty"`
}

func (m *Text) Reset()         { *m = Text{} }
func (m *Text) String() string { return proto.CompactTextString(m) }
func (*Text) ProtoMessage()    {}

// Binary stream content.
type Binary struct {
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *Binary) Reset()         { *m = Binary{} }
func (m *Binary) String() string { return proto.CompactTextString(m) }
func (*Binary) ProtoMessage()    {}

// Datagram stream content type.
type Datagram struct {
	// The size in bytes of the overall datagram.
	Size uint64 `protobuf:"varint,1,opt,name=size" json:"size,omitempty"`
	// This datagram data.
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *Datagram) Reset()         { *m = Datagram{} }
func (m *Datagram) String() string { return proto.CompactTextString(m) }
func (*Datagram) ProtoMessage()    {}

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
	TimeOffset *google_protobuf1.Duration `protobuf:"bytes,1,opt,name=time_offset" json:"time_offset,omitempty"`
	//
	// The message index within the Prefix (required).
	//
	// This is value is unique to this LogEntry across the entire set of entries
	// sharing the stream's Prefix. It is used to designate unambiguous log
	// ordering.
	PrefixIndex uint64 `protobuf:"varint,2,opt,name=prefix_index" json:"prefix_index,omitempty"`
	//
	// The message index within its Stream (required).
	//
	// This value is unique across all entries sharing the same Prefix and Stream
	// Name. It is used to designate unambiguous log ordering within the stream.
	StreamIndex uint64 `protobuf:"varint,3,opt,name=stream_index" json:"stream_index,omitempty"`
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
	// If true, the last content entry in this LogEntry is continued in the first
	// content entry of the stream's next sequential LogEntry.
	//
	// Only Datagrams and Text content types support this field. Binary types are
	// implicitly partial.
	//
	// Note that sparse reassembly is NOT supported. If stream message `N` which
	// has content entries [0, 1, 2] sets partial to `true`, this implies that
	// the last content entry (2) is incomplete. LogEntry `N+1` must begin with
	// the next sequential piece of content entry 2. This can continue across
	// as many sequential LogEntry records are necessary to complete the content.
	//
	// For example, if content entry 2 must be split across multiple LogEntry
	// starting with sequence number `N`, the LogEntry composition must look like:
	// N:   [0, 1, 2]
	// N+1: [2]
	// N+2: [2, 3, ...]
	//
	// As a special case, if the last line of the last Text LogEntry in a stream
	// is marked partial, this means that there is no terminating newline in the
	// stream.
	Partial bool `protobuf:"varint,5,opt,name=partial" json:"partial,omitempty"`
	// Text Stream: Lines of log text.
	Text *Text `protobuf:"bytes,10,opt,name=text" json:"text,omitempty"`
	// Binary stream: data segment.
	Binary *Binary `protobuf:"bytes,11,opt,name=binary" json:"binary,omitempty"`
	// Datagram stream: Datagrams.
	Datagram *Datagram `protobuf:"bytes,12,opt,name=datagram" json:"datagram,omitempty"`
}

func (m *LogEntry) Reset()         { *m = LogEntry{} }
func (m *LogEntry) String() string { return proto.CompactTextString(m) }
func (*LogEntry) ProtoMessage()    {}

func (m *LogEntry) GetTimeOffset() *google_protobuf1.Duration {
	if m != nil {
		return m.TimeOffset
	}
	return nil
}

func (m *LogEntry) GetText() *Text {
	if m != nil {
		return m.Text
	}
	return nil
}

func (m *LogEntry) GetBinary() *Binary {
	if m != nil {
		return m.Binary
	}
	return nil
}

func (m *LogEntry) GetDatagram() *Datagram {
	if m != nil {
		return m.Datagram
	}
	return nil
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
}

func (m *LogIndex) Reset()         { *m = LogIndex{} }
func (m *LogIndex) String() string { return proto.CompactTextString(m) }
func (*LogIndex) ProtoMessage()    {}

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

//
// Entry is a single index entry.
//
// The index is composed of a series of entries, each corresponding to a
// sequential snapshot of of the log stream.
type LogIndex_Entry struct {
	//
	// The sequence number of the first content entry.
	//
	// Text: This is the line index of the first included line. Line indices
	//     begin at zero.
	// Binary: This is the byte offset of the first byte in the included data.
	// Datagram: This is the index of the datagram. The first datagram has index
	//     zero.
	Sequence uint64 `protobuf:"varint,1,opt,name=sequence" json:"sequence,omitempty"`
	//
	// The log index that this entry describes (required).
	//
	// This is used by clients to identify a specific LogEntry within a set of
	// streams sharing a Prefix.
	PrefixIndex uint64 `protobuf:"varint,2,opt,name=prefix_index" json:"prefix_index,omitempty"`
	//
	// The time offset of this log entry (required).
	//
	// This is used by clients to identify a specific LogEntry within a log
	// stream.
	StreamIndex uint64 `protobuf:"varint,3,opt,name=stream_index" json:"stream_index,omitempty"`
	//
	// The time offset of this log entry, in microseconds.
	//
	// This is added to the descriptor's "timestamp" field to identify the
	// specific timestamp of this log. It is used by clients to identify a
	// specific LogEntry by time.
	TimeOffset *google_protobuf1.Duration `protobuf:"bytes,4,opt,name=time_offset" json:"time_offset,omitempty"`
}

func (m *LogIndex_Entry) Reset()         { *m = LogIndex_Entry{} }
func (m *LogIndex_Entry) String() string { return proto.CompactTextString(m) }
func (*LogIndex_Entry) ProtoMessage()    {}

func (m *LogIndex_Entry) GetTimeOffset() *google_protobuf1.Duration {
	if m != nil {
		return m.TimeOffset
	}
	return nil
}

func init() {
	proto.RegisterEnum("protocol.LogStreamDescriptor_StreamType", LogStreamDescriptor_StreamType_name, LogStreamDescriptor_StreamType_value)
	proto.RegisterEnum("protocol.LogStreamDescriptor_TextSource_Newline", LogStreamDescriptor_TextSource_Newline_name, LogStreamDescriptor_TextSource_Newline_value)
	proto.RegisterEnum("protocol.LogStreamDescriptor_TextSource_Encoding", LogStreamDescriptor_TextSource_Encoding_name, LogStreamDescriptor_TextSource_Encoding_value)
}
