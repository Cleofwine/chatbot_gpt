// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.29.1
// source: proto/chatgpt_data.proto

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

// 消息定义
type Record struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Account         string   `protobuf:"bytes,2,opt,name=account,proto3" json:"account,omitempty"`
	GroupId         string   `protobuf:"bytes,3,opt,name=groupId,json=group_id,proto3" json:"groupId,omitempty"`
	UserMsg         string   `protobuf:"bytes,4,opt,name=userMsg,json=user_msg,proto3" json:"userMsg,omitempty"`
	UserMsgTokens   int32    `protobuf:"varint,5,opt,name=userMsgTokens,json=user_msg_tokens,proto3" json:"userMsgTokens,omitempty"`
	UserMsgKeywords []string `protobuf:"bytes,6,rep,name=userMsgKeywords,json=user_msg_keywords,proto3" json:"userMsgKeywords,omitempty"`
	AiMsg           string   `protobuf:"bytes,7,opt,name=aiMsg,json=ai_msg,proto3" json:"aiMsg,omitempty"`
	AiMsgTokens     int32    `protobuf:"varint,8,opt,name=aiMsgTokens,json=ai_msg_tokens,proto3" json:"aiMsgTokens,omitempty"`
	ReqTokens       int32    `protobuf:"varint,9,opt,name=reqTokens,json=req_tokens,proto3" json:"reqTokens,omitempty"`
	CreateAt        int64    `protobuf:"varint,10,opt,name=createAt,json=create_at,proto3" json:"createAt,omitempty"`
	EnterpriseId    string   `protobuf:"bytes,11,opt,name=enterpriseId,json=enterprise_id,proto3" json:"enterpriseId,omitempty"`
	Endpoint        int32    `protobuf:"varint,12,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	EndpointAccount string   `protobuf:"bytes,13,opt,name=endpointAccount,json=endpoint_account,proto3" json:"endpointAccount,omitempty"`
}

func (x *Record) Reset() {
	*x = Record{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chatgpt_data_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Record) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Record) ProtoMessage() {}

func (x *Record) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chatgpt_data_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Record.ProtoReflect.Descriptor instead.
func (*Record) Descriptor() ([]byte, []int) {
	return file_proto_chatgpt_data_proto_rawDescGZIP(), []int{0}
}

func (x *Record) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Record) GetAccount() string {
	if x != nil {
		return x.Account
	}
	return ""
}

func (x *Record) GetGroupId() string {
	if x != nil {
		return x.GroupId
	}
	return ""
}

func (x *Record) GetUserMsg() string {
	if x != nil {
		return x.UserMsg
	}
	return ""
}

func (x *Record) GetUserMsgTokens() int32 {
	if x != nil {
		return x.UserMsgTokens
	}
	return 0
}

func (x *Record) GetUserMsgKeywords() []string {
	if x != nil {
		return x.UserMsgKeywords
	}
	return nil
}

func (x *Record) GetAiMsg() string {
	if x != nil {
		return x.AiMsg
	}
	return ""
}

func (x *Record) GetAiMsgTokens() int32 {
	if x != nil {
		return x.AiMsgTokens
	}
	return 0
}

func (x *Record) GetReqTokens() int32 {
	if x != nil {
		return x.ReqTokens
	}
	return 0
}

func (x *Record) GetCreateAt() int64 {
	if x != nil {
		return x.CreateAt
	}
	return 0
}

func (x *Record) GetEnterpriseId() string {
	if x != nil {
		return x.EnterpriseId
	}
	return ""
}

func (x *Record) GetEndpoint() int32 {
	if x != nil {
		return x.Endpoint
	}
	return 0
}

func (x *Record) GetEndpointAccount() string {
	if x != nil {
		return x.EndpointAccount
	}
	return ""
}

// 服务定义
type RecordRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *RecordRes) Reset() {
	*x = RecordRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chatgpt_data_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecordRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecordRes) ProtoMessage() {}

func (x *RecordRes) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chatgpt_data_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecordRes.ProtoReflect.Descriptor instead.
func (*RecordRes) Descriptor() ([]byte, []int) {
	return file_proto_chatgpt_data_proto_rawDescGZIP(), []int{1}
}

func (x *RecordRes) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

var File_proto_chatgpt_data_proto protoreflect.FileDescriptor

var file_proto_chatgpt_data_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x67, 0x70, 0x74, 0x5f,
	0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x63, 0x68, 0x61, 0x74,
	0x67, 0x70, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x22, 0x9f, 0x03, 0x0a, 0x06, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x19, 0x0a,
	0x07, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x12, 0x19, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x4d, 0x73, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x6d, 0x73, 0x67, 0x12, 0x26, 0x0a, 0x0d, 0x75, 0x73, 0x65, 0x72, 0x4d, 0x73, 0x67, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0f, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x6d, 0x73, 0x67, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x12, 0x2a, 0x0a, 0x0f, 0x75,
	0x73, 0x65, 0x72, 0x4d, 0x73, 0x67, 0x4b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x06,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x11, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x6d, 0x73, 0x67, 0x5f, 0x6b,
	0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x15, 0x0a, 0x05, 0x61, 0x69, 0x4d, 0x73, 0x67,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x69, 0x5f, 0x6d, 0x73, 0x67, 0x12, 0x22,
	0x0a, 0x0b, 0x61, 0x69, 0x4d, 0x73, 0x67, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0d, 0x61, 0x69, 0x5f, 0x6d, 0x73, 0x67, 0x5f, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x73, 0x12, 0x1d, 0x0a, 0x09, 0x72, 0x65, 0x71, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x72, 0x65, 0x71, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x73, 0x12, 0x1b, 0x0a, 0x08, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x74, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x61, 0x74, 0x12, 0x23,
	0x0a, 0x0c, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x49, 0x64, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65,
	0x5f, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18,
	0x0c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12,
	0x29, 0x0a, 0x0f, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69,
	0x6e, 0x74, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x1b, 0x0a, 0x09, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x32, 0x4b, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x74, 0x47,
	0x50, 0x54, 0x44, 0x61, 0x74, 0x61, 0x12, 0x3c, 0x0a, 0x09, 0x41, 0x64, 0x64, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x12, 0x14, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x67, 0x70, 0x74, 0x5f, 0x64, 0x61,
	0x74, 0x61, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x1a, 0x17, 0x2e, 0x63, 0x68, 0x61, 0x74,
	0x67, 0x70, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52,
	0x65, 0x73, 0x22, 0x00, 0x42, 0x14, 0x5a, 0x12, 0x63, 0x68, 0x61, 0x74, 0x67, 0x70, 0x74, 0x2d,
	0x64, 0x61, 0x74, 0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_chatgpt_data_proto_rawDescOnce sync.Once
	file_proto_chatgpt_data_proto_rawDescData = file_proto_chatgpt_data_proto_rawDesc
)

func file_proto_chatgpt_data_proto_rawDescGZIP() []byte {
	file_proto_chatgpt_data_proto_rawDescOnce.Do(func() {
		file_proto_chatgpt_data_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_chatgpt_data_proto_rawDescData)
	})
	return file_proto_chatgpt_data_proto_rawDescData
}

var file_proto_chatgpt_data_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_chatgpt_data_proto_goTypes = []interface{}{
	(*Record)(nil),    // 0: chatgpt_data.Record
	(*RecordRes)(nil), // 1: chatgpt_data.RecordRes
}
var file_proto_chatgpt_data_proto_depIdxs = []int32{
	0, // 0: chatgpt_data.ChatGPTData.AddRecord:input_type -> chatgpt_data.Record
	1, // 1: chatgpt_data.ChatGPTData.AddRecord:output_type -> chatgpt_data.RecordRes
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_chatgpt_data_proto_init() }
func file_proto_chatgpt_data_proto_init() {
	if File_proto_chatgpt_data_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_chatgpt_data_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Record); i {
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
		file_proto_chatgpt_data_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecordRes); i {
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
			RawDescriptor: file_proto_chatgpt_data_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_chatgpt_data_proto_goTypes,
		DependencyIndexes: file_proto_chatgpt_data_proto_depIdxs,
		MessageInfos:      file_proto_chatgpt_data_proto_msgTypes,
	}.Build()
	File_proto_chatgpt_data_proto = out.File
	file_proto_chatgpt_data_proto_rawDesc = nil
	file_proto_chatgpt_data_proto_goTypes = nil
	file_proto_chatgpt_data_proto_depIdxs = nil
}
