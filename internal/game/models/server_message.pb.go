// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.21.1
// source: server_message.proto

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

type ServerMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Types that are assignable to Payload:
	//	*ServerMessage_User_
	//	*ServerMessage_Response_
	//	*ServerMessage_Stats
	//	*ServerMessage_Reports_
	//	*ServerMessage_WorldState
	//	*ServerMessage_StateChange
	Payload isServerMessage_Payload `protobuf_oneof:"payload"`
}

func (x *ServerMessage) Reset() {
	*x = ServerMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerMessage) ProtoMessage() {}

func (x *ServerMessage) ProtoReflect() protoreflect.Message {
	mi := &file_server_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerMessage.ProtoReflect.Descriptor instead.
func (*ServerMessage) Descriptor() ([]byte, []int) {
	return file_server_message_proto_rawDescGZIP(), []int{0}
}

func (x *ServerMessage) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (m *ServerMessage) GetPayload() isServerMessage_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (x *ServerMessage) GetUser() *ServerMessage_User {
	if x, ok := x.GetPayload().(*ServerMessage_User_); ok {
		return x.User
	}
	return nil
}

func (x *ServerMessage) GetResponse() *ServerMessage_Response {
	if x, ok := x.GetPayload().(*ServerMessage_Response_); ok {
		return x.Response
	}
	return nil
}

func (x *ServerMessage) GetStats() *Stats {
	if x, ok := x.GetPayload().(*ServerMessage_Stats); ok {
		return x.Stats
	}
	return nil
}

func (x *ServerMessage) GetReports() *ServerMessage_Reports {
	if x, ok := x.GetPayload().(*ServerMessage_Reports_); ok {
		return x.Reports
	}
	return nil
}

func (x *ServerMessage) GetWorldState() *WorldState {
	if x, ok := x.GetPayload().(*ServerMessage_WorldState); ok {
		return x.WorldState
	}
	return nil
}

func (x *ServerMessage) GetStateChange() *GameStatePatch {
	if x, ok := x.GetPayload().(*ServerMessage_StateChange); ok {
		return x.StateChange
	}
	return nil
}

type isServerMessage_Payload interface {
	isServerMessage_Payload()
}

type ServerMessage_User_ struct {
	User *ServerMessage_User `protobuf:"bytes,3,opt,name=user,proto3,oneof"`
}

type ServerMessage_Response_ struct {
	Response *ServerMessage_Response `protobuf:"bytes,4,opt,name=response,proto3,oneof"`
}

type ServerMessage_Stats struct {
	Stats *Stats `protobuf:"bytes,5,opt,name=stats,proto3,oneof"`
}

type ServerMessage_Reports_ struct {
	Reports *ServerMessage_Reports `protobuf:"bytes,6,opt,name=reports,proto3,oneof"`
}

type ServerMessage_WorldState struct {
	WorldState *WorldState `protobuf:"bytes,7,opt,name=worldState,proto3,oneof"`
}

type ServerMessage_StateChange struct {
	StateChange *GameStatePatch `protobuf:"bytes,9,opt,name=stateChange,proto3,oneof"`
}

func (*ServerMessage_User_) isServerMessage_Payload() {}

func (*ServerMessage_Response_) isServerMessage_Payload() {}

func (*ServerMessage_Stats) isServerMessage_Payload() {}

func (*ServerMessage_Reports_) isServerMessage_Payload() {}

func (*ServerMessage_WorldState) isServerMessage_Payload() {}

func (*ServerMessage_StateChange) isServerMessage_Payload() {}

type ServerMessage_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestId string `protobuf:"bytes,1,opt,name=requestId,proto3" json:"requestId,omitempty"`
	Result    bool   `protobuf:"varint,2,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *ServerMessage_Response) Reset() {
	*x = ServerMessage_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerMessage_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerMessage_Response) ProtoMessage() {}

func (x *ServerMessage_Response) ProtoReflect() protoreflect.Message {
	mi := &file_server_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerMessage_Response.ProtoReflect.Descriptor instead.
func (*ServerMessage_Response) Descriptor() ([]byte, []int) {
	return file_server_message_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ServerMessage_Response) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

func (x *ServerMessage_Response) GetResult() bool {
	if x != nil {
		return x.Result
	}
	return false
}

type ServerMessage_User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
}

func (x *ServerMessage_User) Reset() {
	*x = ServerMessage_User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerMessage_User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerMessage_User) ProtoMessage() {}

func (x *ServerMessage_User) ProtoReflect() protoreflect.Message {
	mi := &file_server_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerMessage_User.ProtoReflect.Descriptor instead.
func (*ServerMessage_User) Descriptor() ([]byte, []int) {
	return file_server_message_proto_rawDescGZIP(), []int{0, 1}
}

func (x *ServerMessage_User) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

type ServerMessage_Reports struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Reports []*Report `protobuf:"bytes,1,rep,name=reports,proto3" json:"reports,omitempty"`
}

func (x *ServerMessage_Reports) Reset() {
	*x = ServerMessage_Reports{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_message_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerMessage_Reports) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerMessage_Reports) ProtoMessage() {}

func (x *ServerMessage_Reports) ProtoReflect() protoreflect.Message {
	mi := &file_server_message_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerMessage_Reports.ProtoReflect.Descriptor instead.
func (*ServerMessage_Reports) Descriptor() ([]byte, []int) {
	return file_server_message_proto_rawDescGZIP(), []int{0, 2}
}

func (x *ServerMessage_Reports) GetReports() []*Report {
	if x != nil {
		return x.Reports
	}
	return nil
}

var File_server_message_proto protoreflect.FileDescriptor

var file_server_message_proto_rawDesc = []byte{
	0x0a, 0x14, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69,
	0x62, 0x65, 0x73, 0x1a, 0x0f, 0x67, 0x61, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x0c, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x0b, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xac, 0x04, 0x0a,
	0x0d, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x35,
	0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x70,
	0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x48, 0x00, 0x52,
	0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x41, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74,
	0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x00, 0x52, 0x08,
	0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74,
	0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74,
	0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x73, 0x48, 0x00, 0x52, 0x05, 0x73,
	0x74, 0x61, 0x74, 0x73, 0x12, 0x3e, 0x0a, 0x07, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69,
	0x62, 0x65, 0x73, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x48, 0x00, 0x52, 0x07, 0x72, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x73, 0x12, 0x39, 0x0a, 0x0a, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61,
	0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x48, 0x00, 0x52, 0x0a, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x3f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62,
	0x65, 0x73, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x50, 0x61, 0x74, 0x63,
	0x68, 0x48, 0x00, 0x52, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x1a, 0x40, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x1a, 0x22, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73,
	0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73,
	0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x1a, 0x38, 0x0a, 0x07, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x73, 0x12, 0x2d, 0x0a, 0x07, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73,
	0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x07, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73,
	0x42, 0x09, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x35, 0x5a, 0x33, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6e, 0x61, 0x74, 0x74, 0x65,
	0x2f, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x2d, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x2f, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_server_message_proto_rawDescOnce sync.Once
	file_server_message_proto_rawDescData = file_server_message_proto_rawDesc
)

func file_server_message_proto_rawDescGZIP() []byte {
	file_server_message_proto_rawDescOnce.Do(func() {
		file_server_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_server_message_proto_rawDescData)
	})
	return file_server_message_proto_rawDescData
}

var file_server_message_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_server_message_proto_goTypes = []interface{}{
	(*ServerMessage)(nil),          // 0: pizzatribes.ServerMessage
	(*ServerMessage_Response)(nil), // 1: pizzatribes.ServerMessage.Response
	(*ServerMessage_User)(nil),     // 2: pizzatribes.ServerMessage.User
	(*ServerMessage_Reports)(nil),  // 3: pizzatribes.ServerMessage.Reports
	(*Stats)(nil),                  // 4: pizzatribes.Stats
	(*WorldState)(nil),             // 5: pizzatribes.WorldState
	(*GameStatePatch)(nil),         // 6: pizzatribes.GameStatePatch
	(*Report)(nil),                 // 7: pizzatribes.Report
}
var file_server_message_proto_depIdxs = []int32{
	2, // 0: pizzatribes.ServerMessage.user:type_name -> pizzatribes.ServerMessage.User
	1, // 1: pizzatribes.ServerMessage.response:type_name -> pizzatribes.ServerMessage.Response
	4, // 2: pizzatribes.ServerMessage.stats:type_name -> pizzatribes.Stats
	3, // 3: pizzatribes.ServerMessage.reports:type_name -> pizzatribes.ServerMessage.Reports
	5, // 4: pizzatribes.ServerMessage.worldState:type_name -> pizzatribes.WorldState
	6, // 5: pizzatribes.ServerMessage.stateChange:type_name -> pizzatribes.GameStatePatch
	7, // 6: pizzatribes.ServerMessage.Reports.reports:type_name -> pizzatribes.Report
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_server_message_proto_init() }
func file_server_message_proto_init() {
	if File_server_message_proto != nil {
		return
	}
	file_gamestate_proto_init()
	file_stats_proto_init()
	file_report_proto_init()
	file_world_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_server_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerMessage); i {
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
		file_server_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerMessage_Response); i {
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
		file_server_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerMessage_User); i {
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
		file_server_message_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerMessage_Reports); i {
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
	file_server_message_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*ServerMessage_User_)(nil),
		(*ServerMessage_Response_)(nil),
		(*ServerMessage_Stats)(nil),
		(*ServerMessage_Reports_)(nil),
		(*ServerMessage_WorldState)(nil),
		(*ServerMessage_StateChange)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_server_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_server_message_proto_goTypes,
		DependencyIndexes: file_server_message_proto_depIdxs,
		MessageInfos:      file_server_message_proto_msgTypes,
	}.Build()
	File_server_message_proto = out.File
	file_server_message_proto_rawDesc = nil
	file_server_message_proto_goTypes = nil
	file_server_message_proto_depIdxs = nil
}