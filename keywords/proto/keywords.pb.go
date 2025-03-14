// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.29.1
// source: proto/keywords.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type FindAllReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text string `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
}

func (x *FindAllReq) Reset() {
	*x = FindAllReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_keywords_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindAllReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindAllReq) ProtoMessage() {}

func (x *FindAllReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_keywords_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindAllReq.ProtoReflect.Descriptor instead.
func (*FindAllReq) Descriptor() ([]byte, []int) {
	return file_proto_keywords_proto_rawDescGZIP(), []int{0}
}

func (x *FindAllReq) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

type FindAllRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Words []string `protobuf:"bytes,1,rep,name=words,proto3" json:"words,omitempty"`
}

func (x *FindAllRes) Reset() {
	*x = FindAllRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_keywords_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindAllRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindAllRes) ProtoMessage() {}

func (x *FindAllRes) ProtoReflect() protoreflect.Message {
	mi := &file_proto_keywords_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindAllRes.ProtoReflect.Descriptor instead.
func (*FindAllRes) Descriptor() ([]byte, []int) {
	return file_proto_keywords_proto_rawDescGZIP(), []int{1}
}

func (x *FindAllRes) GetWords() []string {
	if x != nil {
		return x.Words
	}
	return nil
}

var File_proto_keywords_proto protoreflect.FileDescriptor

var file_proto_keywords_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x63, 0x68, 0x61, 0x74, 0x67, 0x70, 0x74, 0x5f,
	0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x22, 0x20, 0x0a, 0x0a, 0x46, 0x69, 0x6e, 0x64,
	0x41, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x22, 0x22, 0x0a, 0x0a, 0x46, 0x69,
	0x6e, 0x64, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x77, 0x6f, 0x72, 0x64,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x32, 0x5a,
	0x0a, 0x0f, 0x43, 0x68, 0x61, 0x74, 0x47, 0x50, 0x54, 0x4b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64,
	0x73, 0x12, 0x47, 0x0a, 0x07, 0x46, 0x69, 0x6e, 0x64, 0x41, 0x6c, 0x6c, 0x12, 0x1c, 0x2e, 0x63,
	0x68, 0x61, 0x74, 0x67, 0x70, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x2e,
	0x46, 0x69, 0x6e, 0x64, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x1a, 0x1c, 0x2e, 0x63, 0x68, 0x61,
	0x74, 0x67, 0x70, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x2e, 0x46, 0x69,
	0x6e, 0x64, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x22, 0x00, 0x42, 0x18, 0x5a, 0x16, 0x63, 0x68,
	0x61, 0x74, 0x67, 0x70, 0x74, 0x2d, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_keywords_proto_rawDescOnce sync.Once
	file_proto_keywords_proto_rawDescData = file_proto_keywords_proto_rawDesc
)

func file_proto_keywords_proto_rawDescGZIP() []byte {
	file_proto_keywords_proto_rawDescOnce.Do(func() {
		file_proto_keywords_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_keywords_proto_rawDescData)
	})
	return file_proto_keywords_proto_rawDescData
}

var file_proto_keywords_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_keywords_proto_goTypes = []interface{}{
	(*FindAllReq)(nil), // 0: chatgpt_keywords.FindAllReq
	(*FindAllRes)(nil), // 1: chatgpt_keywords.FindAllRes
}
var file_proto_keywords_proto_depIdxs = []int32{
	0, // 0: chatgpt_keywords.ChatGPTKeywords.FindAll:input_type -> chatgpt_keywords.FindAllReq
	1, // 1: chatgpt_keywords.ChatGPTKeywords.FindAll:output_type -> chatgpt_keywords.FindAllRes
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_keywords_proto_init() }
func file_proto_keywords_proto_init() {
	if File_proto_keywords_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_keywords_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindAllReq); i {
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
		file_proto_keywords_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindAllRes); i {
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
			RawDescriptor: file_proto_keywords_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_keywords_proto_goTypes,
		DependencyIndexes: file_proto_keywords_proto_depIdxs,
		MessageInfos:      file_proto_keywords_proto_msgTypes,
	}.Build()
	File_proto_keywords_proto = out.File
	file_proto_keywords_proto_rawDesc = nil
	file_proto_keywords_proto_goTypes = nil
	file_proto_keywords_proto_depIdxs = nil
}
