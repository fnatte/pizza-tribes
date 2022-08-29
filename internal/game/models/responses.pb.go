// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.21.1
// source: responses.proto

package models

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type ApiUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username   string                      `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Ambassador *ApiUserResponse_Ambassador `protobuf:"bytes,2,opt,name=ambassador,proto3" json:"ambassador,omitempty"`
}

func (x *ApiUserResponse) Reset() {
	*x = ApiUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_responses_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApiUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApiUserResponse) ProtoMessage() {}

func (x *ApiUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_responses_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApiUserResponse.ProtoReflect.Descriptor instead.
func (*ApiUserResponse) Descriptor() ([]byte, []int) {
	return file_responses_proto_rawDescGZIP(), []int{0}
}

func (x *ApiUserResponse) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *ApiUserResponse) GetAmbassador() *ApiUserResponse_Ambassador {
	if x != nil {
		return x.Ambassador
	}
	return nil
}

type ApiUserResponse_Ambassador struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Appearance *MouseAppearance `protobuf:"bytes,1,opt,name=appearance,proto3" json:"appearance,omitempty"`
}

func (x *ApiUserResponse_Ambassador) Reset() {
	*x = ApiUserResponse_Ambassador{}
	if protoimpl.UnsafeEnabled {
		mi := &file_responses_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApiUserResponse_Ambassador) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApiUserResponse_Ambassador) ProtoMessage() {}

func (x *ApiUserResponse_Ambassador) ProtoReflect() protoreflect.Message {
	mi := &file_responses_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApiUserResponse_Ambassador.ProtoReflect.Descriptor instead.
func (*ApiUserResponse_Ambassador) Descriptor() ([]byte, []int) {
	return file_responses_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ApiUserResponse_Ambassador) GetAppearance() *MouseAppearance {
	if x != nil {
		return x.Appearance
	}
	return nil
}

var File_responses_proto protoreflect.FileDescriptor

var file_responses_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0b, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x1a, 0x10,
	0x61, 0x70, 0x70, 0x65, 0x61, 0x72, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xc2, 0x01, 0x0a, 0x0f, 0x41, 0x70, 0x69, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x47, 0x0a, 0x0a, 0x61, 0x6d, 0x62, 0x61, 0x73, 0x73, 0x61, 0x64, 0x6f, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62,
	0x65, 0x73, 0x2e, 0x41, 0x70, 0x69, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x2e, 0x41, 0x6d, 0x62, 0x61, 0x73, 0x73, 0x61, 0x64, 0x6f, 0x72, 0x52, 0x0a, 0x61,
	0x6d, 0x62, 0x61, 0x73, 0x73, 0x61, 0x64, 0x6f, 0x72, 0x1a, 0x4a, 0x0a, 0x0a, 0x41, 0x6d, 0x62,
	0x61, 0x73, 0x73, 0x61, 0x64, 0x6f, 0x72, 0x12, 0x3c, 0x0a, 0x0a, 0x61, 0x70, 0x70, 0x65, 0x61,
	0x72, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x70, 0x69,
	0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x4d, 0x6f, 0x75, 0x73, 0x65, 0x41,
	0x70, 0x70, 0x65, 0x61, 0x72, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x0a, 0x61, 0x70, 0x70, 0x65, 0x61,
	0x72, 0x61, 0x6e, 0x63, 0x65, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6e, 0x61, 0x74, 0x74, 0x65, 0x2f, 0x70, 0x69, 0x7a, 0x7a, 0x61,
	0x2d, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x67, 0x61, 0x6d, 0x65, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_responses_proto_rawDescOnce sync.Once
	file_responses_proto_rawDescData = file_responses_proto_rawDesc
)

func file_responses_proto_rawDescGZIP() []byte {
	file_responses_proto_rawDescOnce.Do(func() {
		file_responses_proto_rawDescData = protoimpl.X.CompressGZIP(file_responses_proto_rawDescData)
	})
	return file_responses_proto_rawDescData
}

var file_responses_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_responses_proto_goTypes = []interface{}{
	(*ApiUserResponse)(nil),            // 0: pizzatribes.ApiUserResponse
	(*ApiUserResponse_Ambassador)(nil), // 1: pizzatribes.ApiUserResponse.Ambassador
	(*MouseAppearance)(nil),            // 2: pizzatribes.MouseAppearance
}
var file_responses_proto_depIdxs = []int32{
	1, // 0: pizzatribes.ApiUserResponse.ambassador:type_name -> pizzatribes.ApiUserResponse.Ambassador
	2, // 1: pizzatribes.ApiUserResponse.Ambassador.appearance:type_name -> pizzatribes.MouseAppearance
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_responses_proto_init() }
func file_responses_proto_init() {
	if File_responses_proto != nil {
		return
	}
	file_appearance_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_responses_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApiUserResponse); i {
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
		file_responses_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApiUserResponse_Ambassador); i {
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
			RawDescriptor: file_responses_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_responses_proto_goTypes,
		DependencyIndexes: file_responses_proto_depIdxs,
		MessageInfos:      file_responses_proto_msgTypes,
	}.Build()
	File_responses_proto = out.File
	file_responses_proto_rawDesc = nil
	file_responses_proto_goTypes = nil
	file_responses_proto_depIdxs = nil
}