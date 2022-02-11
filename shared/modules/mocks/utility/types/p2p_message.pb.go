// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.3
// source: p2p_message.proto

package types

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PocketTopic int32

const (
	PocketTopic_CONSENSUS PocketTopic = 0
	PocketTopic_P2P       PocketTopic = 1
)

// Enum value maps for PocketTopic.
var (
	PocketTopic_name = map[int32]string{
		0: "CONSENSUS",
		1: "P2P",
	}
	PocketTopic_value = map[string]int32{
		"CONSENSUS": 0,
		"P2P":       1,
	}
)

func (x PocketTopic) Enum() *PocketTopic {
	p := new(PocketTopic)
	*p = x
	return p
}

func (x PocketTopic) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PocketTopic) Descriptor() protoreflect.EnumDescriptor {
	return file_p2p_message_proto_enumTypes[0].Descriptor()
}

func (PocketTopic) Type() protoreflect.EnumType {
	return &file_p2p_message_proto_enumTypes[0]
}

func (x PocketTopic) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PocketTopic.Descriptor instead.
func (PocketTopic) EnumDescriptor() ([]byte, []int) {
	return file_p2p_message_proto_rawDescGZIP(), []int{0}
}

type NetworkMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic string     `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	Data  *anypb.Any `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *NetworkMessage) Reset() {
	*x = NetworkMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_p2p_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NetworkMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetworkMessage) ProtoMessage() {}

func (x *NetworkMessage) ProtoReflect() protoreflect.Message {
	mi := &file_p2p_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetworkMessage.ProtoReflect.Descriptor instead.
func (*NetworkMessage) Descriptor() ([]byte, []int) {
	return file_p2p_message_proto_rawDescGZIP(), []int{0}
}

func (x *NetworkMessage) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *NetworkMessage) GetData() *anypb.Any {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_p2p_message_proto protoreflect.FileDescriptor

var file_p2p_message_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x32, 0x70, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x07, 0x75, 0x74, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x1a, 0x19, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x51, 0x0a, 0x0f, 0x4e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f,
	0x70, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63,
	0x12, 0x28, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x41, 0x6e, 0x79, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x2a, 0x25, 0x0a, 0x0b, 0x50, 0x6f,
	0x63, 0x6b, 0x65, 0x74, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x0d, 0x0a, 0x09, 0x43, 0x4f, 0x4e,
	0x53, 0x45, 0x4e, 0x53, 0x55, 0x53, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x50, 0x32, 0x50, 0x10,
	0x01, 0x42, 0x11, 0x5a, 0x0f, 0x2f, 0x75, 0x74, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x2f, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_p2p_message_proto_rawDescOnce sync.Once
	file_p2p_message_proto_rawDescData = file_p2p_message_proto_rawDesc
)

func file_p2p_message_proto_rawDescGZIP() []byte {
	file_p2p_message_proto_rawDescOnce.Do(func() {
		file_p2p_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_p2p_message_proto_rawDescData)
	})
	return file_p2p_message_proto_rawDescData
}

var file_p2p_message_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_p2p_message_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_p2p_message_proto_goTypes = []interface{}{
	(PocketTopic)(0),        // 0: utility.PocketTopic
	(*NetworkMessage)(nil), // 1: utility.NetworkMessage
	(*anypb.Any)(nil),       // 2: google.protobuf.Any
}
var file_p2p_message_proto_depIdxs = []int32{
	2, // 0: utility.NetworkMessage.data:type_name -> google.protobuf.Any
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_p2p_message_proto_init() }
func file_p2p_message_proto_init() {
	if File_p2p_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_p2p_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NetworkMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_p2p_message_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_p2p_message_proto_goTypes,
		DependencyIndexes: file_p2p_message_proto_depIdxs,
		EnumInfos:         file_p2p_message_proto_enumTypes,
		MessageInfos:      file_p2p_message_proto_msgTypes,
	}.Build()
	File_p2p_message_proto = out.File
	file_p2p_message_proto_rawDesc = nil
	file_p2p_message_proto_goTypes = nil
	file_p2p_message_proto_depIdxs = nil
}
